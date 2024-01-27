package main

import (
	"image/color"
	"machine"
	"strings"

	"tinygo.org/x/drivers/st7735"
	"tinygo.org/x/tinyfont/proggy"
	"tinygo.org/x/tinyterm"
)

var (
	display = st7735.New(machine.SPI1, machine.TFT_RST, machine.TFT_DC, machine.TFT_CS, machine.TFT_LITE)

	terminal = tinyterm.NewTerminal(&display)

	black = color.RGBA{0, 0, 0, 255}
	white = color.RGBA{255, 255, 255, 255}
	red   = color.RGBA{255, 0, 0, 255}
	blue  = color.RGBA{0, 0, 255, 255}
	green = color.RGBA{0, 255, 0, 255}

	font = &proggy.TinySZ8pt7b
)

var logo = `
  ___ _ _      _   _      
 | __| (_)__ _| |_| |_    
 | _|| | / _\ | ' \  _|   
 |_|_|_|_\__, |_||_\__|   
 | _ ) __|___/| |__ _ ___ 
 | _ \/ _\ / _\ / _\ / -_)
 |___/\__,_\__,_\__, \___|
                |___/     
`

func setupDisplay() {
	machine.SPI1.Configure(machine.SPIConfig{
		SCK:       machine.SPI1_SCK_PIN,
		SDO:       machine.SPI1_SDO_PIN,
		SDI:       machine.SPI1_SDI_PIN,
		Frequency: 8000000,
	})

	display.Configure(st7735.Config{
		Rotation: st7735.ROTATION_90,
	})

	terminal.Configure(&tinyterm.Config{
		Font:              font,
		FontHeight:        10,
		FontOffset:        6,
		UseSoftwareScroll: true,
	})

	display.FillScreen(black)

	showSplash()
}

func showSplash() {
	for _, line := range strings.Split(strings.TrimSuffix(logo, "\n"), "\n") {
		terminal.Write([]byte("\n" + line))
	}
}

func terminalOutput(s string) {
	println(s)
	terminal.Write([]byte("\n" + s))
}
