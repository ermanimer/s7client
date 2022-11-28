package s7client

import (
	"bytes"
	"encoding/binary"
	"errors"
	"testing"
)

func TestErrShortPayload(t *testing.T) {
	c := &client{}

	var p []byte

	_, err := c.Bool(p, 0, 1)
	if !errors.Is(err, ErrShortPayload) {
		t.Error("error is not ErrShortPayload")
	}

	_, err = c.Uint8(p, 0)
	if !errors.Is(err, ErrShortPayload) {
		t.Error("error is not ErrShortPayload")
	}

	_, err = c.Int8(p, 0)
	if !errors.Is(err, ErrShortPayload) {
		t.Error("error is not ErrShortPayload")
	}

	_, err = c.Uint16(p, 0)
	if !errors.Is(err, ErrShortPayload) {
		t.Error("error is not ErrShortPayload")
	}

	_, err = c.Int16(p, 0)
	if !errors.Is(err, ErrShortPayload) {
		t.Error("error is not ErrShortPayload")
	}

	_, err = c.Uint32(p, 0)
	if !errors.Is(err, ErrShortPayload) {
		t.Error("error is not ErrShortPayload")
	}

	_, err = c.Int32(p, 0)
	if !errors.Is(err, ErrShortPayload) {
		t.Error("error is not ErrShortPayload")
	}

	_, err = c.Float32(p, 0)
	if !errors.Is(err, ErrShortPayload) {
		t.Error("error is not ErrShortPayload")
	}

	_, err = c.String(p, 0, 1)
	if !errors.Is(err, ErrShortPayload) {
		t.Error("error is not ErrShortPayload")
	}
}

func TestErrInvalidIndex(t *testing.T) {
	c := &client{}

	p := make([]byte, readResHeaderLen+1)

	_, err := c.Bool(p, 0, -1)
	if !errors.Is(err, ErrInvalidIndex) {
		t.Error("error is not ErrInvalidIndex")
	}
}

func TestErrInvalidLength(t *testing.T) {
	c := &client{}

	p := make([]byte, readResHeaderLen+1)

	_, err := c.String(p, 0, -1)
	if !errors.Is(err, ErrInvalidLength) {
		t.Error("errors is not ErrInvalidLength")
	}
}

func TestBool(t *testing.T) {
	c := &client{}

	var expectedByte byte = 1 // 0000001
	expected := true
	p := make([]byte, readResHeaderLen+1)
	p[readResHeaderLen] = expectedByte

	v, err := c.Bool(p, 0, 0)
	if err != nil {
		t.Error(err)
	}
	if v != expected {
		t.Error("value is not equal to expected", v, expected)
	}
}

func TestUint8(t *testing.T) {
	c := &client{}

	var expected uint8 = 1
	p := make([]byte, readResHeaderLen+1)
	p[readResHeaderLen] = byte(expected)

	v, err := c.Uint8(p, 0)
	if err != nil {
		t.Error(err)
	}
	if v != expected {
		t.Error("value is not equal to expected", v, expected)
	}
}

func TestInt8(t *testing.T) {
	c := &client{}

	var expected int8 = 1
	p := make([]byte, readResHeaderLen+1)
	p[readResHeaderLen] = byte(expected)

	v, err := c.Int8(p, 0)
	if err != nil {
		t.Error(err)
	}
	if v != expected {
		t.Error("value is not equal to expected", v, expected)
	}
}

func TestUint16(t *testing.T) {
	c := &client{}

	var expected uint16 = 1
	p := make([]byte, readResHeaderLen+2)
	binary.BigEndian.PutUint16(p[readResHeaderLen:], expected)

	v, err := c.Uint16(p, 0)
	if err != nil {
		t.Error(err)
	}
	if v != expected {
		t.Error("value is not equal to expected", v, expected)
	}
}

func TestInt16(t *testing.T) {
	c := &client{}

	var expected int16 = 1
	p := make([]byte, readResHeaderLen)
	w := bytes.NewBuffer(p)
	binary.Write(w, binary.BigEndian, expected)
	p = w.Bytes()

	v, err := c.Int16(p, 0)
	if err != nil {
		t.Error(err)
	}
	if v != expected {
		t.Error("value is not equal to expected", v, expected)
	}
}

func TestUint32(t *testing.T) {
	c := &client{}

	var expected uint32 = 1
	p := make([]byte, readResHeaderLen+4)
	binary.BigEndian.PutUint32(p[readResHeaderLen:], expected)

	v, err := c.Uint32(p, 0)
	if err != nil {
		t.Error(err)
	}
	if v != expected {
		t.Error("value is not equal to expected", v, expected)
	}
}

func TestInt32(t *testing.T) {
	c := &client{}

	var expected int32 = 1
	p := make([]byte, readResHeaderLen)
	w := bytes.NewBuffer(p)
	binary.Write(w, binary.BigEndian, expected)
	p = w.Bytes()

	v, err := c.Int32(p, 0)
	if err != nil {
		t.Error(err)
	}
	if v != expected {
		t.Error("value is not equal to expected", v, expected)
	}
}

func TestFloat32(t *testing.T) {
	c := &client{}

	var expected float32 = 1
	p := make([]byte, readResHeaderLen)
	w := bytes.NewBuffer(p)
	binary.Write(w, binary.BigEndian, expected)
	p = w.Bytes()

	v, err := c.Float32(p, 0)
	if err != nil {
		t.Error(err)
	}
	if v != expected {
		t.Error("value is not equal to expected", v, expected)
	}
}

func TestString(t *testing.T) {
	c := &client{}

	expected := "a"
	p := make([]byte, readResHeaderLen+stringHeaderLen)
	p = append(p, []byte(expected)...)

	v, err := c.String(p, 0, len(expected))
	if err != nil {
		t.Error(err)
	}
	if v != expected {
		t.Error("value is not equal to expected", v, expected)
	}
}
