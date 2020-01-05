package lowres129029

import (
	"cant/can/flowroutines"
	"cant/can/flowroutines/filter/manipulate/lowres129025"
	"sort"
)

func GetName() string {
	return "Every n'th GNSS Position Data paired with GNSS Rapid Update"
}

func SupportedPGNs() []uint32 {
	return []uint32{129029} // GNSS Position Data
}

func UniqIdentifier() int {
	return 8
}

type gnssSegment struct {
}

var count = -1
var initial = true
var seq = [7]

func Run(data flowroutines.FlowData) (flowroutines.FlowData bool) {
	payload := data.CanFrame.Data()
	seq := payload[0] % 7
	if seq = 
	if initial {
		count++
		seq[count] = payload[0]
		if count =< 7 {
			return data, true
		}
		sort.Ints(seq)
		last = seq[7]
		intial = false
	}
	last

	lowres129025.LatLonMut.Lock()
	// between here
	lowres129025.LatLonMut.UnLock()
	return last, true // manipulated frame, continue?
}
