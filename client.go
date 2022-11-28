// Package s7client implements Siemens s7 client.
package s7client

import (
	"bytes"
	"encoding/binary"
	"errors"
	"net"
	"time"
)

// Errors:
var (
	ErrUpgradeConn   = errors.New("upgrade connection error")
	ErrNegotiatePDU  = errors.New("negotiate pdu error")
	ErrNotConnected  = errors.New("not connected error")
	ErrShortResponse = errors.New("short response error")
	ErrRead          = errors.New("read error")
	ErrShortPayload  = errors.New("short payload error")
	ErrInvalidIndex  = errors.New("invalid index error")
	ErrInvalidLength = errors.New("invalid length error")
)

// s7 Parameters
const (
	readResHeaderLen = 25
	stringHeaderLen  = 1
)

const defaultResBufSize = 512

// Client defines the behaviors of a Siemens s7 client.
type Client interface {
	// Connect uses net.DialTimeout to establish an underlying TCP connection with the s7 server.
	Connect() error

	// SetDeadline sets the underlying TCP connection's deadline. Returns a s7client.ErrNotconnected if the client is not connected.
	SetDeadline(t time.Time) error

	// Read reads data from a data block of a s7 device and writes it to the provided payload. Returns the read-byte count and a s7client.ErrNotconnected if the client is not connected to the server.
	Read(p []byte, dataBlockNum uint16, addr uint32, count uint16) (n int, err error)

	// ReadErr parses and returns the Modbus read error of the provided payload. Returns a modbusclient.ErrShortResponse if the payload is short.
	ReadErr(p []byte) error

	// Bool parses and returns a bool value fron the provided payload. Returns a s7client.ErrShortResponse if the payload is short.
	Bool(p []byte, offset int, index int) (bool, error)

	// Uint8 parses and returns a uint8 value from the provided payload. Returns a s7client.ErrShortResponse if the payload is short.
	Uint8(p []byte, offset int) (byte, error)

	// Int8 parses and returns an int8 value from the provided payload. Returns a s7client.ErrShortResponse if the payload is short.
	Int8(p []byte, offset int) (int8, error)

	// Uint16 parses and returns an uint16 value from the provided payload. Returns a s7client.ErrShortResponse if the payload is short.
	Uint16(p []byte, offset int) (uint16, error)

	// Int16 parses and returns an int16 value from the provided payload. Returns a s7client.ErrShortResponse if the payload is short.
	Int16(p []byte, offset int) (int16, error)

	// Uint32 parses and returns an uint32 value from the provided payload. Returns a s7client.ErrShortResponse if the payload is short.
	Uint32(p []byte, offset int) (uint32, error)

	// Int32 parses and returns an int32 value from the provided payload. Returns a s7client.ErrShortResponse if the payload is short.
	Int32(p []byte, offset int) (int32, error)

	// Float32 parses and returns a float32 value from the provided payload. Returns a s7client.ErrShortResponse if the payload is short.
	Float32(p []byte, offset int) (float32, error)

	// String parses and returns a string value from the provided payload. Returns a s7client.ErrShortResponse if the payload is short.
	String(p []byte, offset int, length int) (string, error)

	// Close closes the underlying TCP connection. Returns a s7client.ErrNotconnected if the client is not connected to the server.
	Close() error
}

type client struct {
	Addr        string
	Rack        uint16
	Slot        uint16
	ConnTimeout time.Duration
	isoConnReq  []byte
	pduNegReq   []byte
	conn        net.Conn
	resBuf      []byte
}

// NewClient creates and returns a new Siemens s7 Client.
func NewClient(addr string, rack uint16, slot uint16, connTimeout time.Duration) Client {
	return &client{
		Addr:        addr,
		Rack:        rack,
		Slot:        slot,
		ConnTimeout: connTimeout,
		isoConnReq:  makeISOConnReq(rack, slot),
		pduNegReq:   makePDUNegReq(),
		resBuf:      make([]byte, defaultResBufSize),
	}
}

func (c *client) Connect() error {
	if err := c.connect(); err != nil {
		return err
	}

	if err := c.upgradeConn(); err != nil {
		return err
	}

	if err := c.negotiatePDU(); err != nil {
		return err
	}

	return nil
}

func (c *client) connect() error {
	conn, err := net.DialTimeout("tcp4", c.Addr, c.ConnTimeout)
	if err != nil {
		return err
	}
	c.conn = conn
	return nil
}

func (c *client) upgradeConn() error {
	if err := c.conn.SetDeadline(time.Now().Add(c.ConnTimeout)); err != nil {
		return err
	}

	_, err := c.conn.Write(c.isoConnReq)
	if err != nil {
		return err
	}

	n, err := c.conn.Read(c.resBuf)
	if err != nil {
		return err
	}
	if n != 22 {
		return ErrShortResponse
	}
	if c.resBuf[5] != 0xD0 {
		return ErrUpgradeConn
	}
	return nil
}

func makeISOConnReq(rack uint16, slot uint16) []byte {
	tsap := (0x01 << 8) + (rack << 5) + slot
	tsapHigh := byte((tsap >> 8) & 0xFF)
	tsapLow := byte(tsap & 0xFF)
	return []byte{
		0x03, 0x00, 0x00, 0x16,
		0x11, 0xE0, 0x00, 0x00,
		0x00, 0x01, 0x00, 0xC0,
		0x01, 0x0A, 0xC1, 0x02,
		0x01, 0x00, 0xC2, 0x02,
		tsapHigh, tsapLow,
	}
}

func (c *client) negotiatePDU() error {
	if err := c.conn.SetDeadline(time.Now().Add(c.ConnTimeout)); err != nil {
		return err
	}

	_, err := c.conn.Write(c.pduNegReq)
	if err != nil {
		return err
	}

	n, err := c.conn.Read(c.resBuf)
	if err != nil {
		return err
	}
	if n != 27 {
		return ErrShortResponse
	}
	if c.resBuf[17] != 0x00 {
		return ErrNegotiatePDU
	}
	if c.resBuf[18] != 0x00 {
		return ErrNegotiatePDU
	}
	return nil
}

func makePDUNegReq() []byte {
	return []byte{
		0x03, 0x00, 0x00, 0x19,
		0x02, 0xF0, 0x80, 0x32,
		0x01, 0x00, 0x00, 0x04,
		0x00, 0x00, 0x08, 0x00,
		0x00, 0xF0, 0x00, 0x00,
		0x01, 0x00, 0x01, 0x01,
		0xE0,
	}
}

func (c *client) SetDeadline(t time.Time) error {
	if c.conn == nil {
		return ErrNotConnected
	}

	return c.conn.SetDeadline(t)
}

func (c *client) Read(p []byte, dataBlockNum uint16, addr uint32, count uint16) (int, error) {
	if c.conn == nil {
		return 0, ErrNotConnected
	}

	req := makeReadReq(dataBlockNum, addr, count)
	if _, err := c.conn.Write(req); err != nil {
		return 0, err
	}
	return c.conn.Read(p)
}

func makeReadReq(dataBlockNum uint16, addr uint32, count uint16) []byte {
	countHigh := byte((count >> 8) & 0xFF)
	countLow := byte(count & 0xFF)
	dataBlockNumHigh := byte((dataBlockNum >> 8) & 0xFF)
	dataBlockNumLow := byte(dataBlockNum & 0xFF)
	return []byte{
		0x03, 0x00, 0x00, 0x1F,
		0x02, 0xF0, 0x80, 0x32,
		0x01, 0x00, 0x00, 0x05,
		0x00, 0x00, 0x0E, 0x00,
		0x00, 0x04, 0x01, 0x12,
		0x0A, 0x10, 0x02, countHigh,
		countLow, dataBlockNumHigh, dataBlockNumLow, 0x84,
		0x00, 0x00, 0x00,
	}
}

func (c *client) ReadErr(p []byte) error {
	if len(p) < readResHeaderLen {
		return ErrShortResponse
	}

	if p[21] != 0xFF {
		return ErrRead
	}
	return nil
}

func (c *client) Bool(p []byte, offset int, index int) (bool, error) {
	offset += readResHeaderLen
	if len(p) < offset+1 {
		return false, ErrShortPayload
	}

	if index < 0 || index > 7 {
		return false, ErrInvalidIndex
	}

	mask := byte(1 << index)
	v := p[offset]&mask != 0
	return v, nil
}

func (c *client) Uint8(p []byte, offset int) (byte, error) {
	offset += readResHeaderLen
	if len(p) < offset+1 {
		return 0, ErrShortPayload
	}

	v := p[offset]
	return v, nil
}

func (c *client) Int8(p []byte, offset int) (int8, error) {
	offset += readResHeaderLen
	if len(p) < offset+1 {
		return 0, ErrShortPayload
	}

	v := int8(p[offset])
	return v, nil
}

func (c *client) Uint16(p []byte, offset int) (uint16, error) {
	offset += readResHeaderLen
	if len(p) < offset+2 {
		return 0, ErrShortPayload
	}

	return binary.BigEndian.Uint16(p[offset : offset+2]), nil
}

func (c *client) Int16(p []byte, offset int) (int16, error) {
	offset += readResHeaderLen
	if len(p) < offset+2 {
		return 0, ErrShortPayload
	}

	r := bytes.NewReader(p[offset : offset+2])
	var v int16
	if err := binary.Read(r, binary.BigEndian, &v); err != nil {
		return 0, err
	}
	return v, nil
}

func (c *client) Uint32(p []byte, offset int) (uint32, error) {
	offset += readResHeaderLen
	if len(p) < offset+4 {
		return 0, ErrShortPayload
	}

	return binary.BigEndian.Uint32(p[offset : offset+4]), nil
}

func (c *client) Int32(p []byte, offset int) (int32, error) {
	offset += readResHeaderLen
	if len(p) < offset+4 {
		return 0, ErrShortPayload
	}

	r := bytes.NewReader(p[offset : offset+4])
	var v int32
	if err := binary.Read(r, binary.BigEndian, &v); err != nil {
		return 0, err
	}
	return v, nil
}

func (c *client) Float32(p []byte, offset int) (float32, error) {
	offset += readResHeaderLen
	if len(p) < offset+4 {
		return 0, ErrShortPayload
	}

	r := bytes.NewReader(p[offset : offset+4])
	var v float32
	if err := binary.Read(r, binary.BigEndian, &v); err != nil {
		return 0, err
	}
	return v, nil
}

func (c *client) String(p []byte, offset int, length int) (string, error) {
	offset += readResHeaderLen + stringHeaderLen
	if len(p) < offset+length {
		return "", ErrShortPayload
	}

	if length <= 0 {
		return "", ErrInvalidLength
	}

	v := string(p[offset : offset+length])
	return v, nil
}

func (c *client) Close() error {
	if c.conn == nil {
		return ErrNotConnected
	}

	return c.conn.Close()
}
