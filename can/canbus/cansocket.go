//+build linux

package canbus

import (
	"bufio"
	"fmt"
	"net"

	"cant/util/globals"
	"cant/util/networking"

	"golang.org/x/sys/unix"
)

// Socket implements a CAN socket
type Socket struct {
	fd    int
	Iface *net.Interface
	In    *bufio.Reader
	Out   *bufio.Writer
}

// New CAN Socket
func NewSocket(iface string) (*Socket, error) {
	can := Socket{}
	fd, err := unix.Socket(unix.AF_CAN, unix.SOCK_RAW, unix.CAN_RAW)
	if err != nil {
		return nil, err
	}
	can.fd = fd
	ifaces, err := networking.GetCANInterfaces()
	if err != nil {
		return nil, err
	}
	if len(ifaces) == 0 {
		return nil, fmt.Errorf("No CAN interfaces found")
	}
	pos := -1
	for i := range ifaces {
		if ifaces[i].Name == iface {
			pos = i
		}
	}
	if pos == -1 {
		return nil, fmt.Errorf("'%s' is not a CAN interface", iface)
	}
	can.Iface, _ = net.InterfaceByName(iface)
	can.In = bufio.NewReaderSize(&can, 16000)
	can.Out = bufio.NewWriterSize(&can, 16000)
	globals.Statistics.AddInterface <- can.Iface.Name
	return &can, unix.Bind(can.fd, &unix.SockaddrCAN{Ifindex: can.Iface.Index})
}

// Close socket
func (can *Socket) Close() error {
	return unix.Close(can.fd)
}

// io.Reader interface
func (can *Socket) Read(data []byte) (int, error) {
	return unix.Read(can.fd, data)
}

// io.Writer interface
func (can *Socket) Write(data []byte) (int, error) {
	return unix.Write(can.fd, data)
}

// GetFrame gets the next Frame from the socket
func (can *Socket) GetFrame() (Frame, error) {
	raw := [16]byte{}
	if n, err := can.In.Read(raw[:]); err != nil || n != 16 {
		return Frame{}, err
	}
	frame := FrameFromRaw(raw)
	return frame, nil
}

// SendFrame sends a frame
func (can *Socket) SendFrame(frame Frame) error {
	if n, err := can.Out.Write(frame.raw[:]); err != nil || n != 16 {
		return err
	}
	return can.Out.Flush()
}
