# s7client

s7client is a simple Siemens S7 TCP client.

# Supported Functions

- Read Data Blocks

# Supported Data Types

- bool
- uint8
- int8
- uint16
- int16
- uint32
- int32
- float32
- uint64
- int64
- float64

# Installation

```bash
go get -u github.com/ermanimer/s7client
```

# Methods

- **Connect() error:** Connect uses net.DialTimeout to establish an underlying TCP connection with the s7 server.

- **SetDeadline(t time.Time) error:** SetDeadline sets the underlying TCP connection's deadline. Returns a s7client.ErrNotconnected if the client is not connected.

- **Read(p []byte, unitID byte, addr uint16, count uint16) (n int, err error):** Read reads data from a data block of a s7 device and writes it to the provided payload. Returns the read-byte count and a s7client.ErrNotconnected if the client is not connected to the server.
	
- **ReadErr(p []byte) error:** ReadErr parses and returns the read error of the provided payload. Returns a s7client.ErrShortResponse if the payload is short.

- **Bool(p []byte, offset int, index int) (bool, error):** Bool parses and returns a bool value fron the provided payload. Returns a s7client.ErrShortResponse if the payload is short.

- **Uint8(p []byte, offset int) (byte, error):** Uint8 parses and returns a uint8 value from the provided payload. Returns a s7client.ErrShortResponse if the payload is short.

- **Int8(p []byte, offset int) (int8, error):** Int8 parses and returns an int8 value from the provided payload. Returns a s7client.ErrShortResponse if the payload is short.

- **Uint16(p []byte, offset int) (uint16, error):** Uint16 parses and returns an uint16 value from the provided payload. Returns a s7client.ErrShortResponse if the payload is short.

- **Int16(p []byte, offset int) (int16, error):** Int16 parses and returns an int16 value from the provided payload. Returns a s7client.ErrShortResponse if the payload is short.

- **Uint32(p []byte, offset int) (uint32, error):** Uint32 parses and returns an uint32 value from the provided payload. Returns a s7client.ErrShortResponse if the payload is short.


- **Int32(p []byte, offset int) (int32, error):** Int32 parses and returns an int32 value from the provided payload. Returns a s7client.ErrShortResponse if the payload is short.

- **Float32(p []byte, offset int) (float32, error):** Float32 parses and returns a float32 value from the provided payload. Returns a s7client.ErrShortResponse if the payload is short.

- **String(p []byte, offset int, length int) (string, error):** String parses and returns a string value from the provided payload. Returns a s7client.ErrShortResponse if the payload is short.

- **Close() error:** Close closes the underlying TCP connection. Returns a s7client.ErrNotconnected if the client is not connected to the server.

# Sample Application

The sample application demonstrates reading a sample value from a s7 device.** 

```go
package main

import (
	"log"
	"time"

	"github.com/ermanimer/s7client"
)

// configurations
const (
	address          = "192.168.0.1:102" // address of the device
	rack             = 0                 // rack of the device
	slot             = 0                 // slot of the device
	connTimeout      = 5 * time.Second   // connection timeout
	dataBlockNumber  = 1                 // data block number
	startingAddresss = 0                 // starting address
	count            = 2                 // data count
)

func main() {
	// create client
	client := s7client.NewClient(address, rack, slot, connTimeout)

	// connect
	err := client.Connect()
	if err != nil {
		log.Print(err)
		return
	}
	defer func() {
		err := client.Close()
		if err != nil {
			log.Print(err)
		}
	}()

	// read
	buf := make([]byte, 256)
	n, err := client.Read(buf, dataBlockNumber, startingAddresss, count)
	if err != nil {
		log.Print(err)
		return
	}
	response := buf[:n]

	// check read error
	if err := client.ReadErr(response); err != nil {
		log.Print(err)
		return
	}

	// parse value
	value, err := client.Float32(response, 0)
	if err != nil {
		log.Print(err)
		return
	}
	log.Printf("value: %.2f\n", value)
}
```
