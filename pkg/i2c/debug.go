package i2c

import "fmt"

// Debug indicates if the debug mode is enabled. In debug mode, all incoming and outgoing bytes are written to stdout.
var Debug = false

func debugOut(bytes []byte) {
	debug("<", bytes)
}

func debugIn(bytes []byte) {
	debug(">", bytes)
}

func debug(dir string, bytes []byte) {
	if !Debug {
		return
	}

	fmt.Printf("i2c %s % x\n", dir, bytes)
}

func debugError(err error) {
	if !Debug {
		return
	}
	if err == nil {
		return
	}

	fmt.Printf("i2c ! %v\n", err)
}
