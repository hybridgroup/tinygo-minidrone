package main

import (
	"machine"
	"time"
)

type pair struct {
	x int
	y int
}

const (
	// frameX    = 400
	// frameY    = 300
	// frameSize = frameX * frameY * 3
	center  = 32767
	detente = 20000
)

func getLeftStick() pair {
	s := pair{x: 0, y: 0}
	s.x = leftX
	s.y = leftY
	return s
}

func getRightStick() pair {
	s := pair{x: 0, y: 0}
	s.x = rightX
	s.y = rightY
	return s
}

func initPins() {
	// buttons
	b1.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	b2.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	b3.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	b4.Configure(machine.PinConfig{Mode: machine.PinInputPullup})

	// joystick
	machine.InitADC()
	stickX.Configure(machine.ADCConfig{})
	stickY.Configure(machine.ADCConfig{})
}

func readControls() {
	for {
		stickmode := "right"
		b1push = false
		b2push = false
		b3push = false
		b4push = false

		if !b1.Get() {
			b1push = true
			println("takeoff")
			err := drone.TakeOff()
			if err != nil {
				println(err)
			}
		}
		if !b2.Get() {
			b2push = true
			println("land")
			err := drone.Land()
			if err != nil {
				println(err)
			}
		}
		if !b3.Get() {
			b3push = true
			stickmode = "left"
		}
		if !b4.Get() {
			b4push = true
		}

		// read control stick
		xPos = stickX.Get()
		yPos = stickY.Get()
		if stickmode == "right" {
			// set left to center position
			leftX = center
			leftY = center

			// set right x,y to stick values
			rightX = int(xPos)
			rightY = int(yPos)
		} else {
			// set left x,y to stick values
			leftX = int(xPos)
			leftY = int(yPos)

			// set right to center position
			rightX = center
			rightY = center
		}

		time.Sleep(time.Millisecond * 100)
	}
}
