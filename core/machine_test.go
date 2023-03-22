package core

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMachine(t *testing.T) {
	s := "(+ 1 2)"

	stream := NewByteStream([]byte{})
	reader := Reader(stream)
	machine := NewMachine(reader)

	stream.AppendData([]byte(s)).Mark()
	// _, err := h.reader.Next(c, true)
	// h.stream.Reset(c)
	res, err := machine.Eval(context.Background())
	assert.Nil(t, err)
	rs, _ := res.AsString()
	assert.Equal(t, "3", rs)
}
