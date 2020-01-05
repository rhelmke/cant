package fakeview58880

import (
	"encoding/binary"
	"time"

	"cant/can/flowroutines"
	"cant/can/flowroutines/filter/manipulate/fakeview59136"
	"cant/util/globals"
)

func GetName() string {
	return "SetNumericValue Response Filter"
}

func SupportedPGNs() []uint32 {
	return []uint32{58880}
}

func UniqIdentifier() int {
	return 7
}

var count = 0
var last flowroutines.FlowData
var targetOID = uint16(318)
var srcAddr = 200
var dstAddr = 135

func Run(data flowroutines.FlowData) (flowroutines.FlowData, bool) {
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
	// Request filter is not active
	if !fakeview59136.ValueChanActive {
		return data, true
	}
	// for testing purposes: set raw numeric value to 0x00 0x00 0x00
	replace := <-fakeview59136.ValueChan
	for i := 0; i < len(replace); i++ {
		content[4+i] = replace[i]
	}
	data.CanFrame = data.CanFrame.SetData(content)
	logdata := globals.Livelog.NewLogData()
	logdata.Identifier = "FakeView58880"
	logdata.Timestamp = time.Now().UnixNano()
	logdata.Msg = "Replaced Data in Set Numeric Value response"
	globals.Livelog.Send <- logdata
	globals.Statistics.AddManipulated <- uint64(1)
	return data, true // manipulated frame, continue?
}
