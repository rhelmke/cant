package awgn60928dummy

import (
	"encoding/binary"
	"math/rand"
	"time"

	"cant/can/flowroutines"
	"cant/util/globals"
)

func GetName() string {
	return "Additive White Gaussian Noise for PGN 60928 (Dummy)"
}

func SupportedPGNs() []uint32 {
	return []uint32{60928}
}

func UniqIdentifier() int {
	return 3
}

func Run(data flowroutines.FlowData) (flowroutines.FlowData, bool) {
	logdata := globals.Livelog.NewLogData()
	logdata.Identifier = "AWGN 60928"
	logdata.Timestamp = time.Now().UnixNano()
	logdata.Msg = "Manipulated Frame with PGN 60928 using AWGN"
	globals.Livelog.Send <- logdata
	content := data.CanFrame.Data()
	binary.LittleEndian.PutUint16(content[:2], uint16(float64(0x7FFF)+rand.NormFloat64()))
	binary.LittleEndian.PutUint16(content[2:4], uint16(float64(0x7FFF)+rand.NormFloat64()))
	binary.LittleEndian.PutUint16(content[4:6], uint16(float64(0x7FFF)+rand.NormFloat64()))
	binary.LittleEndian.PutUint16(content[6:], uint16(float64(0x7FFF)+rand.NormFloat64()))
	data.CanFrame = data.CanFrame.SetData(content)
	globals.Statistics.AddManipulated <- uint64(1)
	return data, true // manipulated frame, continue?
}
