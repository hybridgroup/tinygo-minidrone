// takeoff is a tinygo example that connects to a Parrot minidrone and
// causes it to takeoff and land.
//
// You can run this example either on your computer or on a microcontroller with Bluetooth support.
//
// On your computer:
// go run ./examples/takeoff 4C:D2:6C:17:82:6E
//
// On a microcontroller with Bluetooth support:
// tinygo flash -target=nano-rp2040 -ldflags="-X main.DeviceAddress=4C:D2:6C:17:82:6E" ./examples/takeoff
package main

import (
	"time"

	minidrone "github.com/hybridgroup/tinygo-minidrone"
	"tinygo.org/x/bluetooth"
)

var (
	adapter = bluetooth.DefaultAdapter
	device  bluetooth.Device
	ch      = make(chan bluetooth.ScanResult, 1)

	drone *minidrone.Minidrone
)

func main() {
	wait()

	println("enabling...")

	must("enable BLE interface", adapter.Enable())

	println("start scan...")

	must("start scan", adapter.Scan(scanHandler))

	var err error
	select {
	case result := <-ch:
		device, err = adapter.Connect(result.Address, bluetooth.ConnectionParams{})
		must("connect to peripheral device", err)

		println("connected to ", result.Address.String())
	}

	defer device.Disconnect()

	drone = minidrone.NewMinidrone(&device)
	err = drone.Start()
	if err != nil {
		failMessage(err.Error())
	}

	time.Sleep(3 * time.Second)

	println("takeoff")
	err = drone.TakeOff()
	if err != nil {
		failMessage(err.Error())
	}
	time.Sleep(10 * time.Second)

	println("land")
	err = drone.Land()
	if err != nil {
		failMessage(err.Error())
	}

	drone.Halt()

	done()
}

func scanHandler(a *bluetooth.Adapter, d bluetooth.ScanResult) {
	println("device:", d.Address.String(), d.RSSI, d.LocalName())
	if d.Address.String() == connectAddress() {
		a.StopScan()
		ch <- d
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
