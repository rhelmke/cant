package isobus

import (
	"cant/can/canbus"
)

type Frame struct {
	EDP bool
	DP  bool
	PF  uint32
	PS  uint32
	PGN uint32
}

func NewFrame(cf canbus.Frame) Frame {
	id := cf.ID()
	frame := Frame{
		EDP: ((id >> 25) & 0x01) == 0x01,
		DP:  ((id >> 24) & 0x01) == 0x01,
		PF:  (id >> 16) & 0xff,
		PS:  (id >> 8) & 0xff,
		PGN: uint32(0),
	}
	edpUint := uint32(0)
	if frame.EDP {
		edpUint = 1
	}
	dpUint := uint32(0)
	if frame.DP {
		dpUint = 1
	}
	frame.PGN += edpUint<<17 + dpUint<<16 + frame.PF<<8
	if frame.PF >= 240 {
		frame.PGN += frame.PS
	}
	return frame
}
