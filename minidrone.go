package minidrone

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"sync"
	"time"

	"tinygo.org/x/bluetooth"
)

type Minidrone struct {
	device                     *bluetooth.Device
	commandService             *bluetooth.DeviceService
	commandCharacteristic      *bluetooth.DeviceCharacteristic
	pcmdCharacteristic         *bluetooth.DeviceCharacteristic
	notificationService        *bluetooth.DeviceService
	flightStatusCharacteristic *bluetooth.DeviceCharacteristic

	buf       []byte
	stepsfa0a uint16
	stepsfa0b uint16
	pcmdMutex sync.Mutex
	Flying    bool
	Pcmd      Pcmd
	pcmddata  []byte
	shutdown  chan bool

	pilotingStateHandler func(state, substate int)
}

var (
	// BLE services
	droneCommandServiceUUID      = bluetooth.NewUUID([16]byte{0x9a, 0x66, 0xfa, 0x00, 0x08, 0x00, 0x91, 0x91, 0x11, 0xe4, 0x01, 0x2d, 0x15, 0x40, 0xcb, 0x8e})
	droneNotificationServiceUUID = bluetooth.NewUUID([16]byte{0x9a, 0x66, 0xfb, 0x00, 0x08, 0x00, 0x91, 0x91, 0x11, 0xe4, 0x01, 0x2d, 0x15, 0x40, 0xcb, 0x8e})

	// send characteristics
	pcmdCharacteristicUUID     = bluetooth.NewUUID([16]byte{0x9a, 0x66, 0xfa, 0x0a, 0x08, 0x00, 0x91, 0x91, 0x11, 0xe4, 0x01, 0x2d, 0x15, 0x40, 0xcb, 0x8e})
	commandCharacteristicUUID  = bluetooth.NewUUID([16]byte{0x9a, 0x66, 0xfa, 0x0b, 0x08, 0x00, 0x91, 0x91, 0x11, 0xe4, 0x01, 0x2d, 0x15, 0x40, 0xcb, 0x8e})
	priorityCharacteristicUUID = bluetooth.NewUUID([16]byte{0x9a, 0x66, 0xfa, 0x0c, 0x08, 0x00, 0x91, 0x91, 0x11, 0xe4, 0x01, 0x2d, 0x15, 0x40, 0xcb, 0x8e})

	// receive characteristics
	flightStatusCharacteristicUUID = bluetooth.NewUUID([16]byte{0x9a, 0x66, 0xfb, 0x0e, 0x08, 0x00, 0x91, 0x91, 0x11, 0xe4, 0x01, 0x2d, 0x15, 0x40, 0xcb, 0x8e})
	batteryCharacteristicUUID      = bluetooth.NewUUID([16]byte{0x9a, 0x66, 0xfb, 0x0f, 0x08, 0x00, 0x91, 0x91, 0x11, 0xe4, 0x01, 0x2d, 0x15, 0x40, 0xcb, 0x8e})
)

const (
	// piloting states
	PilotingStateFlatTrimChanged    = 0
	PilotingStateFlyingStateChanged = 1

	// flying states
	FlyingStateLanded    = 0
	FlyingStateTakeoff   = 1
	FlyingStateHovering  = 2
	FlyingStateFlying    = 3
	FlyingStateLanding   = 4
	FlyingStateEmergency = 5
	FlyingStateRolling   = 6

	// Battery event
	Battery = "battery"

	// FlightStatus event
	FlightStatus = "flightstatus"

	// Takeoff event
	Takeoff = "takeoff"

	// Hovering event
	Hovering = "hovering"

	// Flying event
	Flying = "flying"

	// Landing event
	Landing = "landing"

	// Landed event
	Landed = "landed"

	// Emergency event
	Emergency = "emergency"

	// Rolling event
	Rolling = "rolling"

	// FlatTrimChange event
	FlatTrimChange = "flattrimchange"

	// LightFixed mode for LightControl
	LightFixed = 0

	// LightBlinked mode for LightControl
	LightBlinked = 1

	// LightOscillated mode for LightControl
	LightOscillated = 3

	// ClawOpen mode for ClawControl
	ClawOpen = 0

	// ClawClosed mode for ClawControl
	ClawClosed = 1
)

// Pcmd is the Parrot Command structure for flight control
type Pcmd struct {
	Flag  int
	Roll  int
	Pitch int
	Yaw   int
	Gaz   int
	Psi   float32
}

func NewMinidrone(dev *bluetooth.Device) *Minidrone {
	n := &Minidrone{
		device: dev,
		Pcmd: Pcmd{
			Flag:  0,
			Roll:  0,
			Pitch: 0,
			Yaw:   0,
			Gaz:   0,
			Psi:   0,
		},
		pcmddata: make([]byte, 19),
		shutdown: make(chan bool),
		buf:      make([]byte, 255),
	}

	return n
}

func (m *Minidrone) PilotingStateChange(handler func(state, substate int)) {
	m.pilotingStateHandler = handler
}

func (m *Minidrone) Start() (err error) {
	if debug {
		println("drone: Start")
	}
	srvcs, err := m.device.DiscoverServices([]bluetooth.UUID{
		droneCommandServiceUUID,
		droneNotificationServiceUUID,
	})
	switch {
	case err != nil:
		return err
	case len(srvcs) == 0:
		return errors.New("could not find drone services")
	}

	m.commandService = &srvcs[0]
	m.notificationService = &srvcs[1]
	if debug {
		println("found drone command service", m.commandService.UUID().String())
		println("found drone notify service", m.notificationService.UUID().String())
	}

	chars, err := m.commandService.DiscoverCharacteristics([]bluetooth.UUID{
		commandCharacteristicUUID,
		pcmdCharacteristicUUID,
	})
	switch {
	case err != nil:
		return err
	case len(chars) == 0:
		return errors.New("could not find drone command characteristics")
	}

	if debug {
		println("found drone command characteristics", chars[0].UUID().String(), chars[1].UUID().String())
	}
	m.commandCharacteristic = &chars[0]
	m.pcmdCharacteristic = &chars[1]

	chars, err = m.notificationService.DiscoverCharacteristics([]bluetooth.UUID{
		flightStatusCharacteristicUUID,
	})
	switch {
	case err != nil:
		return err
	case len(chars) == 0:
		return errors.New("could not find drone notify characteristics")
	}

	if debug {
		println("found drone notify characteristics", chars[0].UUID().String())
	}
	m.flightStatusCharacteristic = &chars[0]

	err = m.Init()
	if err != nil {
		if debug {
			println("init error", err.Error())
		}
		return err
	}

	if debug {
		println("drone init complete")
	}
	m.FlatTrim()
	m.StartPcmd()
	m.FlatTrim()

	return
}

// Halt stops minidrone driver (void)
func (m *Minidrone) Halt() (err error) {
	m.Land()

	m.shutdown <- true
	time.Sleep(500 * time.Millisecond)
	return
}

// Init initializes the BLE insterfaces used by the Minidrone
func (m *Minidrone) Init() (err error) {
	if debug {
		println("init")
	}
	err = m.GenerateAllStates()
	if err != nil {
		println(err.Error())
		return
	}

	if debug {
		println("enabling pcmd notifications")
	}

	// if you do not enable these notifications, then you cannot send commands to the drone.
	err = m.flightStatusCharacteristic.EnableNotifications(func(buf []byte) {
		m.processFlightStatus(buf)
	})

	// TODO: subscribe to battery notifications

	return
}

func (m *Minidrone) Disconnect() {
	m.device.Disconnect()
}

// GenerateAllStates sets up all the default states aka settings on the drone
func (m *Minidrone) GenerateAllStates() (err error) {
	m.stepsfa0b++
	buf := []byte{0x04, byte(m.stepsfa0b) & 0xff, 0x00, 0x04, 0x01, 0x00, 0x32, 0x30, 0x31, 0x34, 0x2D, 0x31, 0x30, 0x2D, 0x32, 0x38, 0x00}
	_, err = m.commandCharacteristic.WriteWithoutResponse(buf)

	return err
}

// TakeOff tells the Minidrone to takeoff
func (m *Minidrone) TakeOff() (err error) {
	m.stepsfa0b++
	buf := []byte{0x02, byte(m.stepsfa0b) & 0xff, 0x02, 0x00, 0x01, 0x00}
	_, err = m.commandCharacteristic.WriteWithoutResponse(buf)

	return err
}

// Land tells the Minidrone to land
func (m *Minidrone) Land() (err error) {
	m.stepsfa0b++
	buf := []byte{0x02, byte(m.stepsfa0b) & 0xff, 0x02, 0x00, 0x03, 0x00}
	_, err = m.commandCharacteristic.WriteWithoutResponse(buf)

	return err
}

// FlatTrim calibrates the Minidrone to use its current position as being level
func (m *Minidrone) FlatTrim() (err error) {
	m.stepsfa0b++
	buf := []byte{0x02, byte(m.stepsfa0b) & 0xff, 0x02, 0x00, 0x00, 0x00}
	_, err = m.commandCharacteristic.WriteWithoutResponse(buf)

	return err
}

// Emergency sets the Minidrone into emergency mode
func (m *Minidrone) Emergency() (err error) {
	m.stepsfa0b++
	buf := []byte{0x02, byte(m.stepsfa0b) & 0xff, 0x02, 0x00, 0x04, 0x00}
	_, err = m.commandCharacteristic.WriteWithoutResponse(buf)

	return err
}

// StartPcmd starts the continuous Pcmd communication with the Minidrone
func (m *Minidrone) StartPcmd() {
	go func() {
		// wait a little bit so that there is enough time to get some ACKs
		time.Sleep(500 * time.Millisecond)
		for {
			select {
			case <-m.shutdown:
				return
			default:
			}

			m.generatePcmd()
			_, err := m.pcmdCharacteristic.WriteWithoutResponse(m.pcmddata)
			if err != nil {
				fmt.Println("pcmd write error:", err)
			}
			time.Sleep(50 * time.Millisecond)
		}
	}()
}

// Up tells the drone to ascend. Pass in an int from 0-100.
func (m *Minidrone) Up(val int) error {
	m.pcmdMutex.Lock()
	defer m.pcmdMutex.Unlock()

	m.Pcmd.Flag = 1
	m.Pcmd.Gaz = validatePitch(val)
	return nil
}

// Down tells the drone to descend. Pass in an int from 0-100.
func (m *Minidrone) Down(val int) error {
	m.pcmdMutex.Lock()
	defer m.pcmdMutex.Unlock()

	m.Pcmd.Flag = 1
	m.Pcmd.Gaz = validatePitch(val) * -1
	return nil
}

// Forward tells the drone to go forward. Pass in an int from 0-100.
func (m *Minidrone) Forward(val int) error {
	m.pcmdMutex.Lock()
	defer m.pcmdMutex.Unlock()

	m.Pcmd.Flag = 1
	m.Pcmd.Pitch = validatePitch(val)
	return nil
}

// Backward tells drone to go in reverse. Pass in an int from 0-100.
func (m *Minidrone) Backward(val int) error {
	m.pcmdMutex.Lock()
	defer m.pcmdMutex.Unlock()

	m.Pcmd.Flag = 1
	m.Pcmd.Pitch = validatePitch(val) * -1
	return nil
}

// Right tells drone to go right. Pass in an int from 0-100.
func (m *Minidrone) Right(val int) error {
	m.pcmdMutex.Lock()
	defer m.pcmdMutex.Unlock()

	m.Pcmd.Flag = 1
	m.Pcmd.Roll = validatePitch(val)
	return nil
}

// Left tells drone to go left. Pass in an int from 0-100.
func (m *Minidrone) Left(val int) error {
	m.pcmdMutex.Lock()
	defer m.pcmdMutex.Unlock()

	m.Pcmd.Flag = 1
	m.Pcmd.Roll = validatePitch(val) * -1
	return nil
}

// Clockwise tells drone to rotate in a clockwise direction. Pass in an int from 0-100.
func (m *Minidrone) Clockwise(val int) error {
	m.pcmdMutex.Lock()
	defer m.pcmdMutex.Unlock()

	m.Pcmd.Flag = 1
	m.Pcmd.Yaw = validatePitch(val)
	return nil
}

// CounterClockwise tells drone to rotate in a counter-clockwise direction.
// Pass in an int from 0-100.
func (m *Minidrone) CounterClockwise(val int) error {
	m.pcmdMutex.Lock()
	defer m.pcmdMutex.Unlock()

	m.Pcmd.Flag = 1
	m.Pcmd.Yaw = validatePitch(val) * -1
	return nil
}

// Hover tells the drone to stop moving in any direction and simply hover in place
func (m *Minidrone) Hover() error {
	m.pcmdMutex.Lock()
	defer m.pcmdMutex.Unlock()

	m.Pcmd = Pcmd{
		Flag:  0,
		Roll:  0,
		Pitch: 0,
		Yaw:   0,
		Gaz:   0,
		Psi:   0,
	}

	return nil
}

// FrontFlip tells the drone to perform a front flip
func (m *Minidrone) FrontFlip() error {
	_, err := m.commandCharacteristic.WriteWithoutResponse(m.generateAnimation(0))

	return err
}

// BackFlip tells the drone to perform a backflip
func (m *Minidrone) BackFlip() error {
	_, err := m.commandCharacteristic.WriteWithoutResponse(m.generateAnimation(1))

	return err
}

// RightFlip tells the drone to perform a flip to the right
func (m *Minidrone) RightFlip() error {
	_, err := m.commandCharacteristic.WriteWithoutResponse(m.generateAnimation(2))

	return err
}

// LeftFlip tells the drone to perform a flip to the left
func (m *Minidrone) LeftFlip() error {
	_, err := m.commandCharacteristic.WriteWithoutResponse(m.generateAnimation(3))

	return err
}

func (m *Minidrone) generateAnimation(anim int) []byte {
	m.stepsfa0b++
	return []byte{0x02, byte(m.stepsfa0b) & 0xff, 0x02, 0x04, 0x00, 0x00, byte(anim), 0x00, 0x00, 0x00}
}

func FlyingState(state int) string {
	switch state {
	case FlyingStateLanded:
		return Landed
	case FlyingStateTakeoff:
		return Takeoff
	case FlyingStateHovering:
		return Hovering
	case FlyingStateFlying:
		return Flying
	case FlyingStateLanding:
		return Landing
	case FlyingStateEmergency:
		return Emergency
	case FlyingStateRolling:
		return Rolling
	}

	return "unknown"
}

func (m *Minidrone) generatePcmd() {
	m.pcmdMutex.Lock()
	defer m.pcmdMutex.Unlock()

	m.stepsfa0a++
	m.pcmddata[0] = 0x02
	m.pcmddata[1] = byte(m.stepsfa0a)
	m.pcmddata[2] = 0x02
	m.pcmddata[3] = 0x00
	m.pcmddata[4] = 0x02
	m.pcmddata[5] = 0x00
	m.pcmddata[6] = byte(m.Pcmd.Flag)
	m.pcmddata[7] = byte(m.Pcmd.Roll)
	m.pcmddata[8] = byte(m.Pcmd.Pitch)
	m.pcmddata[9] = byte(m.Pcmd.Yaw)
	m.pcmddata[10] = byte(m.Pcmd.Gaz)
	binary.LittleEndian.PutUint32(m.buf[11:], math.Float32bits(m.Pcmd.Psi))
	m.pcmddata[15] = 0x00
	m.pcmddata[16] = 0x00
	m.pcmddata[17] = 0x00
	m.pcmddata[18] = 0x00

	return
}

func (m *Minidrone) processFlightStatus(data []byte) {
	if len(data) < 5 {
		// ignore, just a sync
		return
	}

	switch data[4] {
	case PilotingStateFlatTrimChanged:
		if debug {
			println("flatTrimChanged")
		}

		if m.pilotingStateHandler != nil {
			m.pilotingStateHandler(int(data[4]), 0)
		}

	case PilotingStateFlyingStateChanged:
		switch data[6] {
		case FlyingStateLanded:
			if m.Flying {
				m.Flying = false
				if debug {
					println("flyingStateLanded")
				}
			}

		case FlyingStateTakeoff:
			if debug {
				println("flyingStateTakeoff")
			}

		case FlyingStateHovering:
			if !m.Flying {
				m.Flying = true
				if debug {
					println("flyingStateHovering")
				}
			}

		case FlyingStateFlying:
			if !m.Flying {
				m.Flying = true
				if debug {
					println("flyingStateFlying")
				}
			}

		case FlyingStateLanding:
			if debug {
				println("flyingStateLanding")
			}

		case FlyingStateEmergency:
			if debug {
				println("flyingStateEmergency")
			}

		case FlyingStateRolling:
			if debug {
				println("flyingStateRolling")
			}

		}

		if m.pilotingStateHandler != nil {
			m.pilotingStateHandler(int(data[4]), int(data[6]))
		}
	}
}

// ValidatePitch helps validate pitch values such as those created by
// a joystick to values between 0-100 that are required as
// params to Parrot Minidrone PCMDs
func ValidatePitch(data float64, offset float64) int {
	value := math.Abs(data) / offset
	if value >= 0.1 {
		if value <= 1.0 {
			return int((float64(int(value*100)) / 100) * 100)
		}
		return 100
	}
	return 0
}

func validatePitch(val int) int {
	if val > 100 {
		return 100
	} else if val < 0 {
		return 0
	}

	return val
}
