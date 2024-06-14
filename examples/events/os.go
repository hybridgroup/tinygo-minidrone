//go:build !baremetal

package main

import "os"

func connectAddress() string {
	if len(os.Args) < 2 {
		println("usage: events [address]")
		os.Exit(1)
	}

	address := os.Args[1]

	return address
}

// wait on baremetal, proceed immediately on desktop OS.
func wait() {
}

// done just prints a message and allows program to exit.
func done() {
	println("Done.")
}

func failMessage(msg string) {
	println(msg)
	os.Exit(1)
}
