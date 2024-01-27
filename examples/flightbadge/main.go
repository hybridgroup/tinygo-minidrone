// flightbadge is a tinygo example that connects to a Parrot Mambo drone and
// allows you to control it with the Adafruit PyBadge.
//
// You can run this example with the following command:
// tinygo flash -target=pybadge -ldflags="-X main.DeviceAddress=4C:D2:6C:17:82:6E" ./examples/flightbadge
package main

import (
	"strconv"
	"time"

	minidrone "github.com/hybridgroup/tinygo-minidrone"
	"tinygo.org/x/bluetooth"
)

const speed = 20

var (
	adapter = bluetooth.DefaultAdapter
	device  bluetooth.Device
	ch      = make(chan bluetooth.ScanResult, 1)

	DeviceAddress string

	drone *minidrone.Minidrone

	droneconnected bool
	takeoff        bool
	direction      int
	shifted        bool
)

func main() {
	setupDisplay()
	time.Sleep(3 * time.Second)

	terminalOutput("enable bluetooth adapter...")
	must("enable BLE interface", adapter.Enable())

	terminalOutput("start scan...")
	must("start scan", adapter.Scan(scanHandler))

	var err error
	select {
	case result := <-ch:
		terminalOutput("connecting to " + result.Address.String() + "...")

		device, err = adapter.Connect(result.Address, bluetooth.ConnectionParams{})
		must("connect to peripheral device", err)

		terminalOutput("connected to " + result.Address.String())
		droneconnected = true
	}

	defer device.Disconnect()

	terminalOutput("staring drone...")

	drone = minidrone.NewMinidrone(&device)
	must("drone start", drone.Start())

	terminalOutput("ready")

	go readControls()
	controlDrone()
}

func scanHandler(a *bluetooth.Adapter, d bluetooth.ScanResult) {
	terminalOutput("device: " + d.Address.String() + " " + strconv.Itoa(int(d.RSSI)) + " " + d.LocalName())
	if d.Address.String() == DeviceAddress {
		a.StopScan()
		ch <- d
	}
}

func controlDrone() {
	for {
		if !droneconnected {
			time.Sleep(100 * time.Millisecond)
			continue
		}

		switch direction {
		case directionForward:
			drone.Forward(speed)
		case directionBackward:
			drone.Backward(speed)
		default:
			drone.Forward(0)
		}

		switch direction {
		case directionLeft:
			drone.Left(speed)
		case directionRight:
			drone.Right(speed)
		default:
			drone.Right(0)
		}

		switch direction {
		case directionUp:
			drone.Up(speed)
		case directionDown:
			drone.Down(speed)
		default:
			drone.Up(0)
		}

		switch direction {
		case directionTurnLeft:
			drone.CounterClockwise(speed)
		case directionTurnRight:
			drone.Clockwise(speed)
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
