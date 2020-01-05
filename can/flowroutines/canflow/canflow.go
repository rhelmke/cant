package canflow

import (
	"bufio"
	"net"

	"sync"
	"time"

	"cant/can/canbus"
	"cant/can/flowroutines"
	"cant/util/globals"
)

// Reader is a type implementing the flow.Routine interface. It can be used to create a data source for the flow package
type Reader struct {
	in           *bufio.Reader
	iface        *net.Interface
	defaultRoute int
}

// Writer is a type implementing the flow.Routine interface. It can be used to create a data sink for the flow package
type Writer struct {
	out          *bufio.Writer
	iface        *net.Interface
	defaultRoute int
}

// CreateReader creates a flow.Routine for the flow package
func CreateReader(socket *canbus.Socket) *Reader {
	return &Reader{in: socket.In, iface: socket.Iface}
}

// CreateWriter creates a flow.Routine for the flow package
func CreateWriter(socket *canbus.Socket) *Writer {
	return &Writer{out: socket.Out, iface: socket.Iface}
}

// DefaultRoute creates the default route between a Reader and a Writer. It doesnt matter if you call this on the writer or reader
func (fr *Reader) DefaultRoute(fw *Writer) {
	fr.defaultRoute = fw.iface.Index
	fw.defaultRoute = fr.iface.Index
}

// DefaultRoute creates the default route between a Reader and a Writer. It doesnt matter if you call this on the writer or reader
func (fw *Writer) DefaultRoute(fr *Reader) {
	fw.defaultRoute = fr.iface.Index
	fr.defaultRoute = fw.iface.Index
}

// FlowRun implements the flow.Routine interface for the FlowWriter
func (fw *Writer) FlowRun(in <-chan flowroutines.FlowData, out chan<- flowroutines.FlowData, exit chan bool, wg *sync.WaitGroup) {
	defer wg.Done()
	socketchan := make(chan [16]byte, 1000)
	go func() {
		for {
			select {
			case raw := <-socketchan:
				for i := range raw {
					if err := fw.out.WriteByte(raw[i]); err != nil {
						return
					}
				}
				if err := fw.out.Flush(); err != nil {
					return
				}
			}
		}
	}()
	for {
		// https://stackoverflow.com/a/11121616
		select {
		case data := <-in:
			if data.DstInterface == 0 || data.DstInterface == fw.iface.Index {
				frame := data.CanFrame
				socketchan <- frame.Raw()
				globals.Statistics.GetOutChannel(fw.iface.Name) <- uint64(1)
				logdata := globals.Livelog.NewLogData()
				logdata.Timestamp = time.Now().UnixNano()
				logdata.Identifier = fw.iface.Name
				logdata.Msg = "Sent CAN Frame"
				globals.Livelog.Send <- logdata
			}
			continue
		default:
		}
		select {
		case data := <-in:
			if data.DstInterface == 0 || data.DstInterface == fw.iface.Index {
				frame := data.CanFrame
				socketchan <- frame.Raw()
				globals.Statistics.GetOutChannel(fw.iface.Name) <- uint64(1)
				logdata := globals.Livelog.NewLogData()
				logdata.Timestamp = time.Now().UnixNano()
				logdata.Identifier = fw.iface.Name
				logdata.Msg = "Sent CAN Frame"
				globals.Livelog.Send <- logdata
			}
			continue
		case <-exit:
		}
		break
	}
}

// FlowRun implements the flow.Routine interface for the FlowReader
func (fr *Reader) FlowRun(in <-chan flowroutines.FlowData, out chan<- flowroutines.FlowData, exit chan bool, wg *sync.WaitGroup) {
	defer wg.Done()
	// go does not consider non-blocking io operations as best practice. We need to spin up a new goroutine.
	socketchan := make(chan [16]byte, 1000)
	go func() {
		for {
			raw := [16]byte{}
			for i := range raw {
				var err error
				raw[i], err = fr.in.ReadByte()
				if err != nil {
					return
				}
			}
			socketchan <- raw
		}
	}()
	for {
		// https://stackoverflow.com/a/11121616
		select {
		case raw := <-socketchan:
			frame := canbus.FrameFromRaw(raw)
			logdata := globals.Livelog.NewLogData()
			logdata.Identifier = fr.iface.Name
			logdata.Timestamp = time.Now().UnixNano()
			logdata.Msg = "Received CAN Frame"
			globals.Livelog.Send <- logdata
			out <- flowroutines.FlowData{CanFrame: frame, SrcInterface: fr.iface.Index, DstInterface: fr.defaultRoute}
			globals.Statistics.GetInChannel(fr.iface.Name) <- uint64(1)
			continue
		default:
		}
		select {
		case data := <-in:
			out <- data
			continue
		default:
		}
		select {
		case raw := <-socketchan:
			frame := canbus.FrameFromRaw(raw)
			logdata := globals.Livelog.NewLogData()
			logdata.Identifier = fr.iface.Name
			logdata.Timestamp = time.Now().UnixNano()
			logdata.Msg = "Received CAN Frame"
			globals.Livelog.Send <- logdata
			out <- flowroutines.FlowData{CanFrame: frame, SrcInterface: fr.iface.Index, DstInterface: fr.defaultRoute}
			globals.Statistics.GetInChannel(fr.iface.Name) <- uint64(1)
			continue
		case data := <-in:
			out <- data
			continue
		case <-exit:
		}
		break
	}
}
