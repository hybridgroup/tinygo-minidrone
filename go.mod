module github.com/hybridgroup/tinygo-minidrone

go 1.20

replace tinygo.org/x/bluetooth => /home/ron/Development/tinygo/bluetooth

require (
	tinygo.org/x/bluetooth v0.8.1-0.20240122175515-17114fd460ad
	tinygo.org/x/drivers v0.26.1-0.20240117074700-3c5e17423a16
	tinygo.org/x/tinydraw v0.0.0-20200416172542-c30d6d84353c
	tinygo.org/x/tinyfont v0.4.0
	tinygo.org/x/tinyterm v0.3.1-0.20231207163921-6842651de7e1
)

require (
	github.com/go-ole/go-ole v1.2.6 // indirect
	github.com/godbus/dbus/v5 v5.1.0 // indirect
	github.com/google/shlex v0.0.0-20191202100458-e7afc7fbc510 // indirect
	github.com/saltosystems/winrt-go v0.0.0-20230921082907-2ab5b7d431e1 // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	github.com/tinygo-org/cbgo v0.0.4 // indirect
	golang.org/x/sys v0.11.0 // indirect
)
