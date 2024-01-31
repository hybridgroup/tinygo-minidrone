package main

import (
	"machine"
	"time"

	minidrone "github.com/hybridgroup/tinygo-minidrone"
	"tinygo.org/x/bluetooth"
)

// replace this with the MAC address of the Parrot Minidrone you want to connect to.
var DeviceAddress string

var (
	// buttons
	b1   = machine.D7
	b2   = machine.D9
	b3   = machine.D10
	b4   = machine.D11
	bjoy = machine.D12

	// joystick
	stickX = machine.ADC{machine.A0}
	stickY = machine.ADC{machine.A1}

	xPos                                     uint16
	yPos                                     uint16
	b1push, b2push, b3push, b4push, bjoypush bool
	leftX, leftY, rightX, rightY             int
	droneconnected                           bool

	adapter = bluetooth.DefaultAdapter
	device  *bluetooth.Device
	ch      = make(chan bluetooth.ScanResult, 1)

	drone *minidrone.Minidrone
	speed = 20
)

func main() {
	machine.I2C0.Configure(machine.I2CConfig{})

	initPins()

	go handleDisplay()
	time.Sleep(3 * time.Second)

	must("enable BLE interface", adapter.Enable())

	println("start scan...")
	must("start scan", adapter.Scan(scanHandler))

	var err error
	select {
	case result := <-ch:
		device, err = adapter.Connect(result.Address, bluetooth.ConnectionParams{})
		must("connect to peripheral device", err)

		println("connected to ", result.Address.String())
		droneconnected = true
	}

	defer device.Disconnect()

	drone = minidrone.NewMinidrone(device)
	must("drone start", drone.Start())

	go readControls()
	controlDrone()
}

func scanHandler(a *bluetooth.Adapter, d bluetooth.ScanResult) {
	println("device:", d.Address.String(), d.RSSI, d.LocalName())
	if d.Address.String() == DeviceAddress {
		a.StopScan()
		ch <- d
	}
}

func controlDrone() {
	for {
		rightStick := getRightStick()

		switch {
		case rightStick.y+detente < center:
			drone.Backward(speed)
		case rightStick.y-detente > center:
			drone.Forward(speed)
		default:
			drone.Forward(0)
		}

		switch {
		case rightStick.x-detente > center:
			drone.Right(speed)
		case rightStick.x+detente < center:
			drone.Left(speed)
		default:
			drone.Right(0)
		}

		leftStick := getLeftStick()

		switch {
		case leftStick.y+detente < center:
			drone.Down(speed)
		case leftStick.y-detente > center:
			drone.Up(speed)
		default:
			drone.Up(0)
		}

		switch {
		case leftStick.x-detente > center:
			drone.Clockwise(speed)
		case leftStick.x+detente < center:
			drone.CounterClockwise(speed)
		default:
			drone.Clockwise(0)
		}

		time.Sleep(100 * time.Millisecond)
	}
}

func must(action string, err error) {
	if err != nil {
		for {
			println("failed to " + action + ": " + err.Error())
			time.Sleep(time.Second)
		}
	}
}
