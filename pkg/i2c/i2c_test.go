package i2c

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteReg(t *testing.T) {
	buffer := &fileBuffer{*bytes.NewBuffer([]byte{})}
	b := I2C{
		io: buffer,
	}

	n, err := b.WriteReg(1, 2, 3, 4)
	assert.NoError(t, err)
	assert.Equal(t, 4, n)
	assert.Equal(t, []byte{1, 2, 3, 4}, buffer.Bytes())
}

func TestReadReg(t *testing.T) {
	buffer := &fileBuffer{*bytes.NewBuffer([]byte{4, 3, 2})}
	b := I2C{
		io: buffer,
	}

	bytes := make([]byte, 3)
	n, err := b.ReadReg(1, bytes)

	assert.Equal(t, 3, n)
	assert.NoError(t, err)
	assert.Equal(t, []byte{1}, buffer.Bytes())
	assert.Equal(t, []byte{4, 3, 2}, bytes)
}

type fileBuffer struct {
	bytes.Buffer
}

func (b *fileBuffer) Close() error {
	return nil
}
