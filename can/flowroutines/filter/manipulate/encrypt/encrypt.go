package encrypt

import (
	"cant/can/flowroutines"
	"cant/util/globals"
	"time"
)

func GetName() string {
	return "Encrypt"
}

func UniqIdentifier() int {
	return 11
}

func SupportedPGNs() []uint32 {
	return []uint32{0} // all
}

func Run(data flowroutines.FlowData) (flowroutines.FlowData, bool) {
	logdata := globals.Livelog.NewLogData()
	logdata.Identifier = "Cryptor"
	logdata.Timestamp = time.Now().UnixNano()
	logdata.Msg = "Encrypted Frame"
	globals.Livelog.Send <- logdata
	globals.Statistics.AddManipulated <- uint64(1)
	return data, true // manipulated frame, continue?
}
