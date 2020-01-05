package lowres129025

import (
	"time"

	"cant/can/flowroutines"
	"cant/util/globals"
)

func GetName() string {
	return "Every n'th GNSS Position, Rapid Update PGN"
}

func SupportedPGNs() []uint32 {
	return []uint32{129025} // GNSS Position, Rapid Update
}

func UniqIdentifier() int {
	return 5
}

var count = 0
var last flowroutines.FlowData

// var LatLonMut = &sync.Mutex{}
var Lat [4]byte
var Lon [4]byte

func init() {
	// LatLonMut.Lock()
}

func Run(data flowroutines.FlowData) (flowroutines.FlowData, bool) {
	// LatLonMut.Unlock()
	if count%100 == 0 {
		last = data
		latlon := data.CanFrame.Data()
		// LatLonMut.Lock()
		for i := 0; i < len(latlon); i++ {
			if i < 4 {
				Lat[i] = latlon[i]
			} else {
				Lon[i-4] = latlon[i]
			}
		}
		// LatLonMut.Unlock()
		logdata := globals.Livelog.NewLogData()
		logdata.Identifier = "Rapid GNSS Filter"
		logdata.Timestamp = time.Now().UnixNano()
		logdata.Msg = "Got new PGN 129025"
		globals.Livelog.Send <- logdata
	} else {
		globals.Statistics.AddManipulated <- uint64(1)
	}
	count++
	return last, true // manipulated frame, continue?
}
