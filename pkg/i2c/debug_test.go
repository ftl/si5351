package i2c

import (
	"errors"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDebug(t *testing.T) {
	Debug = true
	defer func() {
		Debug = false
	}()

	out, err := grabStdout(func() {
		debugOut([]byte{1, 2, 3, 4, 5})
	})
	assert.NoError(t, err)
	assert.Equal(t, "i2c < 01 02 03 04 05\n", out)

	in, err := grabStdout(func() {
		debugIn([]byte{6, 7, 8, 9, 0})
	})
	assert.NoError(t, err)
	assert.Equal(t, "i2c > 06 07 08 09 00\n", in)

	errOut, err := grabStdout(func() {
		debugError(errors.New("Fail"))
	})
	assert.NoError(t, err)
	assert.Equal(t, "i2c ! Fail\n", errOut)
}

func TestNoDebug(t *testing.T) {
	Debug = false

	out, err := grabStdout(func() {
		debugOut([]byte{1, 2, 3, 4, 5})
	})
	assert.NoError(t, err)
	assert.Equal(t, "", out)

	in, err := grabStdout(func() {
		debugIn([]byte{6, 7, 8, 9, 0})
	})
	assert.NoError(t, err)
	assert.Equal(t, "", in)

	errOut, err := grabStdout(func() {
		debugError(errors.New("Fail"))
	})
	assert.NoError(t, err)
	assert.Equal(t, "", errOut)
}

func grabStdout(f func()) (string, error) {
	in, out, err := os.Pipe()
	if err != nil {
		return "", err
	}

	oldOut := os.Stdout
	os.Stdout = out
	defer func() {
		os.Stdout = oldOut
	}()

	var b []byte
	var readErr error
	waiter := make(chan struct{})
	go func() {
		b, readErr = ioutil.ReadAll(in)
		close(waiter)
	}()

	f()
	out.Close()
	<-waiter

	if readErr != nil {
		return "", readErr
	}

	return string(b), nil
}
