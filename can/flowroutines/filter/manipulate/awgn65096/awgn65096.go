package awgn65096

import (
	"cant/can/flowroutines"
	"cant/util/globals"
	"encoding/binary"
	"math/rand"
	"time"
)

func GetName() string {
	return "Additive White Gaussian Noise for PGN 65096"
}

func SupportedPGNs() []uint32 {
	return []uint32{65096} // Wheel-based Speed and Distance
}

func UniqIdentifier() int {
	return 2
}

func Run(data flowroutines.FlowData) (flowroutines.FlowData, bool) {
	logdata := globals.Livelog.NewLogData()
	logdata.Identifier = "AWGN 65096"
	logdata.Timestamp = time.Now().UnixNano()
	logdata.Msg = "Manipulated Frame with PGN 65096 using AWGN"
	globals.Livelog.Send <- logdata
	content := data.CanFrame.Data()
	contentSpeed := binary.LittleEndian.Uint16(content[:2])
	kmh := float64(contentSpeed) * 0.001 * 0.001 * 60.0 * 60.0
	manipulatedKmh := float64(kmh + rand.NormFloat64())
	if manipulatedKmh < 0 {
		manipulatedKmh = 0
	}
	binary.LittleEndian.PutUint16(content[:2], uint16(manipulatedKmh*1000.0*1000.0/60.0/60.0))
	data.CanFrame = data.CanFrame.SetData(content)
	globals.Statistics.AddManipulated <- uint64(1)
	return data, true // manipulated frame, continue?
}
