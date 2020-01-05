//+build linux

package canbus

import (
	"encoding/binary"

	"golang.org/x/sys/unix"
)

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
	frame := Frame{id: id, extended: extended, dlc: dlc, data: data}.buildRaw()
	return frame
}

// FrameFromRaw build a can.Frame out of a raw byte array
func FrameFromRaw(raw [16]byte) Frame {
	arbitrationField := binary.LittleEndian.Uint32(raw[:4])
	extendedFormat := arbitrationField&unix.CAN_EFF_FLAG == unix.CAN_EFF_FLAG
	id := arbitrationField
	if extendedFormat {
		id &= unix.CAN_EFF_MASK
	} else {
		id &= unix.CAN_SFF_MASK
	}
	dlc := raw[4] & 0xF // 4 lower bits are used for length
	frame := Frame{id: id, extended: extendedFormat, dlc: dlc, raw: raw}
	copy(frame.data[:], raw[8:])
	return frame
}

// buildRaw builds the Raw Portion of a Frame
func (frame Frame) buildRaw() Frame {
	frame.raw = [16]byte{}
	arbitrationField := frame.id
	if !frame.extended {
		arbitrationField = arbitrationField & unix.CAN_SFF_MASK
	} else {
		arbitrationField = arbitrationField | unix.CAN_EFF_FLAG
	}
	binary.LittleEndian.PutUint32(frame.raw[:4], arbitrationField)
	frame.raw[4] = frame.dlc
	copy(frame.raw[8:], frame.data[:])
	return frame
}

// Setters

// SetID sets the ID
func (frame Frame) SetID(id uint32) Frame {
	frame.id = id
	return frame.buildRaw()
}

// SetDLC sets the DLC field of a frame
func (frame Frame) SetDLC(dlc byte) Frame {
	frame.dlc = dlc
	return frame.buildRaw()
}

// SetExtended sets the Extended bit of a frame
func (frame Frame) SetExtended(extended bool) Frame {
	frame.extended = extended
	return frame.buildRaw()
}

// SetData sets the data of a frame
func (frame Frame) SetData(data [8]byte) Frame {
	frame.data = data
	return frame.buildRaw()
}

// Getters

// ID of CanFrame
func (frame Frame) ID() uint32 {
	return frame.id
}

// DLC represents the length of the CAN data field
func (frame Frame) DLC() byte {
	return frame.dlc
}

// Extended Frame Format?
func (frame Frame) Extended() bool {
	return frame.extended
}

// Data Field of Frame
func (frame Frame) Data() [8]byte {
	return frame.data
}

// Raw Frame Field
func (frame Frame) Raw() [16]byte {
	return frame.raw
}
