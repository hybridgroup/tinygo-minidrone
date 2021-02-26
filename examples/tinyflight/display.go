package main

import (
	"machine"

	"image/color"
	"strconv"
	"time"

	"tinygo.org/x/drivers/ssd1306"
	"tinygo.org/x/tinydraw"
	"tinygo.org/x/tinyfont"
	"tinygo.org/x/tinyfont/freemono"
)

func handleDisplay() {
	println("init display")

	display := ssd1306.NewI2C(machine.I2C0)
	display.Configure(ssd1306.Config{
		Address: ssd1306.Address_128_32,
		Width:   128,
		Height:  64,
	})

	display.ClearDisplay()

	black := color.RGBA{1, 1, 1, 255}

	for {
		display.ClearBuffer()

		if !droneconnected {
			msg := []byte("connecting")
			tinyfont.WriteLine(&display, &freemono.Bold9pt7b, 10, 20, string(msg), black)
		} else {
			x := strconv.Itoa(int(xPos))
			y := strconv.Itoa(int(yPos))
			msg := []byte("x: " + x)
			tinyfont.WriteLine(&display, &freemono.Bold9pt7b, 10, 20, string(msg), black)

			msg2 := []byte("y: " + y)
			tinyfont.WriteLine(&display, &freemono.Bold9pt7b, 10, 40, string(msg2), black)

			var radius int16 = 4
			if b1push {
				tinydraw.FilledCircle(&display, 16+32*0, 64-radius-1, radius, black)
			} else {
				tinydraw.Circle(&display, 16+32*0, 64-radius-1, radius, black)
			}
			if b2push {
				tinydraw.FilledCircle(&display, 16+32*1, 64-radius-1, radius, black)
			} else {
				tinydraw.Circle(&display, 16+32*1, 64-radius-1, radius, black)
			}
			if b3push {
				tinydraw.FilledCircle(&display, 16+32*2, 64-radius-1, radius, black)
			} else {
				tinydraw.Circle(&display, 16+32*2, 64-radius-1, radius, black)
			}
			if b4push {
				if drone.Flying {
					tinydraw.Rectangle(&display, 16+32*3, 64-radius-1, radius, radius, black)
				} else {
					tinydraw.FilledCircle(&display, 16+32*3, 64-radius-1, radius, black)
				}
			} else {
				tinydraw.Circle(&display, 16+32*3, 64-radius-1, radius, black)
			}
		}

		display.Display()

		time.Sleep(200 * time.Millisecond)
	}
}
