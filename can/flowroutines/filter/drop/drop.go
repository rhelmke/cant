package drop

import (
	"cant/can/flowroutines"
	"cant/util/globals"
	"time"
)

func GetName() string {
	return "Drop"
}

func UniqIdentifier() int {
	return 1
}

func SupportedPGNs() []uint32 {
	return []uint32{0} // all
}

func Run(data flowroutines.FlowData) (flowroutines.FlowData, bool) {
	logdata := globals.Livelog.NewLogData()
	logdata.Identifier = "Dropper"
	logdata.Timestamp = time.Now().UnixNano()
	logdata.Msg = "Dropped Frame"
	globals.Livelog.Send <- logdata
	globals.Statistics.AddBlocked <- uint64(1)
	return data, false // manipulated frame, continue?
}
