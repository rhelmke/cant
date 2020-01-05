package limit65096

import (
	"cant/can/flowroutines"
	"cant/util/globals"
	"encoding/binary"
	"math/rand"
	"time"
)

func GetName() string {
	return "Limit Wheel-based Speed to 14 km/h"
}

func SupportedPGNs() []uint32 {
	return []uint32{65096} // Wheel-based Speed and Distance
}

func UniqIdentifier() int {
	return 4
}

func Run(data flowroutines.FlowData) (flowroutines.FlowData, bool) {
	content := data.CanFrame.Data()
	contentSpeed := binary.LittleEndian.Uint16(content[:2])
	kmh := float64(contentSpeed) * 0.001 * 0.001 * 60.0 * 60.0
	if kmh <= 14.0 {
		return data, true
	}
	manipulatedKmh := float64(14.0 + rand.NormFloat64()*0.0325)
	binary.LittleEndian.PutUint16(content[:2], uint16(manipulatedKmh*1000.0*1000.0/60.0/60.0))
	data.CanFrame = data.CanFrame.SetData(content)
	globals.Statistics.AddManipulated <- uint64(1)
	logdata := globals.Livelog.NewLogData()
	logdata.Identifier = "Speed Limiter"
	logdata.Timestamp = time.Now().UnixNano()
	logdata.Msg = "Manipulated Frame with PGN 65096 by capping the speed at 14 kmh"
	globals.Livelog.Send <- logdata
	return data, true // manipulated frame, continue?
}
