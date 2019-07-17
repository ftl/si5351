package i2c

import (
	"fmt"
	"io"
	"os"
	"syscall"
)

// I2C bus.
type I2C struct {
	address uint8
	bus     int
	io      io.ReadWriteCloser
	err     error
}

// Open a new I2C connection to the device with the given address on the given I2C bus.
func Open(address uint8, bus int) (*I2C, error) {
	filename := fmt.Sprintf("/dev/i2c-%d", bus)
	file, err := os.OpenFile(filename, os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}

	if err := initCommunication(file.Fd(), address); err != nil {
		return nil, err
	}
	result := &I2C{
		address: address,
		bus:     bus,
		io:      file,
	}

	return result, nil
}

func initCommunication(filedescriptor uintptr, address uint8) error {
	const i2cSlave = 0x0703
	_, _, err := syscall.Syscall6(syscall.SYS_IOCTL, filedescriptor, i2cSlave, uintptr(address), 0, 0, 0)
	if err != 0 {
		return err
	}
	return nil
}

// Address of the device this instance is communicating with.
func (b *I2C) Address() uint8 {
	return b.address
}

// Bus on which this instance is communicating.
func (b *I2C) Bus() int {
	return b.bus
}

// Err returns the first error that happened when using this instance or nil if everything is fine.
// If Err returns an error != nil, all further invocations of Read and Write will also fail with this error. This instance
// should not be used anymore.
func (b *I2C) Err() error {
	return b.err
}

// Read from the bus.
func (b *I2C) Read(p []byte) (int, error) {
	if b.err != nil {
		debugError(b.err)
		return 0, b.err
	}
	n, err := b.io.Read(p)
	if err != nil {
		b.err = err
		debugError(err)
	}
	debugIn(p)
	return n, err
}

// Write to the bus.
func (b *I2C) Write(p []byte) (int, error) {
	if b.err != nil {
		debugError(b.err)
		return 0, b.err
	}
	n, err := b.io.Write(p)
	if err != nil {
		b.err = err
		debugError(err)
	}
	debugOut(p)
	return n, err
}

// Close the I2C communication
func (b *I2C) Close() error {
	return b.io.Close()
}

// ReadReg reads len(p) bytes from the given register (and the following if len(p) > 1).
func (b *I2C) ReadReg(reg uint8, p []byte) (int, error) {
	buf := make([]byte, 1)
	n := 0
	for i := range p {
		if _, err := b.Write([]byte{reg}); err != nil {
			return 0, err
		}
		if _, err := b.Read(buf); err != nil {
			p[i] = buf[0]
			n++
		}
	}
	return n, b.err
}

// WriteReg writes the given bytes to the given register (and the following if there is more than one byte given).
func (b *I2C) WriteReg(reg uint8, values ...byte) (int, error) {
	return b.Write(append([]byte{reg}, values...))
}

// RegWriter returns a writer that writes to the given register.
func (b *I2C) RegWriter(reg uint8) io.Writer {
	return &regWriter{reg: reg, bus: b}
}

type regWriter struct {
	reg uint8
	bus *I2C
}

func (w *regWriter) Write(p []byte) (int, error) {
	return w.bus.WriteReg(w.reg, p...)
}
