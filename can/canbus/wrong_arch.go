//+build !linux darwin windows

package canbus

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

// Socket isnt used here
type Socket struct {
	fd    int
	Iface *net.Interface
	In    *bufio.Reader
	Out   *bufio.Writer
}

// NewSocket creates a new can Socket
func NewSocket(iface string) (*Socket, error) {
	return nil, fmt.Errorf("can't doesn't support build architectures other than linux")
}

// Close socket
func (can *Socket) Close() error {
	return fmt.Errorf("can't doesn't support build architectures other than linux")
}

// io.Reader interface
func (can *Socket) Read(data []byte) (int, error) {
	return 0, fmt.Errorf("can't doesn't support build architectures other than linux")
}

// io.Writer interface
func (can *Socket) Write(data []byte) (int, error) {
	return 0, fmt.Errorf("can't doesn't support build architectures other than linux")
}

// GetFrame gets the next Frame from the socket
func (can *Socket) GetFrame() (Frame, error) {
	return Frame{}, fmt.Errorf("can't doesn't support build architectures other than linux")
}

// SendFrame sends a frame
func (can *Socket) SendFrame(frame Frame) error {
	return fmt.Errorf("can't doesn't support build architectures other than linux")
}

// Frame represents a CAN frame
type Frame struct {
	id       uint32
	extended bool
	dlc      byte
	data     [8]byte
	raw      [16]byte
}

// NewFrame ...
func NewFrame(id uint32, extended bool, dlc byte, data [8]byte) Frame {
	fmt.Fprintf(os.Stderr, "can't doesn't support build architectures other than linux")
	os.Exit(1)
	return Frame{}
}

// FrameFromRaw build a can.Frame out of a raw byte array
func FrameFromRaw(raw [16]byte) Frame {
	fmt.Fprintf(os.Stderr, "can't doesn't support build architectures other than linux")
	os.Exit(1)
	return Frame{}
}

// buildRaw builds the Raw Portion of a Frame
func (frame Frame) buildRaw() Frame {
	fmt.Fprintf(os.Stderr, "can't doesn't support build architectures other than linux")
	os.Exit(1)
	return Frame{}
}

// Setters

// SetID sets the ID
func (frame Frame) SetID(id uint32) Frame {
	fmt.Fprintf(os.Stderr, "can't doesn't support build architectures other than linux")
	os.Exit(1)
	return Frame{}
}

// SetDLC sets the DLC field of a frame
func (frame Frame) SetDLC(dlc byte) Frame {
	fmt.Fprintf(os.Stderr, "can't doesn't support build architectures other than linux")
	os.Exit(1)
	return Frame{}
}

// SetExtended sets the Extended bit of a frame
func (frame Frame) SetExtended(extended bool) Frame {
	fmt.Fprintf(os.Stderr, "can't doesn't support build architectures other than linux")
	os.Exit(1)
	return Frame{}
}

// SetData sets the data of a frame
func (frame Frame) SetData(data [8]byte) Frame {
	fmt.Fprintf(os.Stderr, "can't doesn't support build architectures other than linux")
	os.Exit(1)
	return Frame{}
}

// Getters

// ID of CanFrame
func (frame Frame) ID() uint32 {
	fmt.Fprintf(os.Stderr, "can't doesn't support build architectures other than linux")
	os.Exit(1)
	return 0
}

// DLC represents the length of the CAN data field
func (frame Frame) DLC() byte {
	fmt.Fprintf(os.Stderr, "can't doesn't support build architectures other than linux")
	os.Exit(1)
	return 0
}

// Extended Frame Format?
func (frame Frame) Extended() bool {
	fmt.Fprintf(os.Stderr, "can't doesn't support build architectures other than linux")
	os.Exit(1)
	return false
}

// Data Field of Frame
func (frame Frame) Data() [8]byte {
	fmt.Fprintf(os.Stderr, "can't doesn't support build architectures other than linux")
	os.Exit(1)
	return [8]byte{}
}

// Raw Frame Field
func (frame Frame) Raw() [16]byte {
	fmt.Fprintf(os.Stderr, "can't doesn't support build architectures other than linux")
	os.Exit(1)
	return [16]byte{}
}
