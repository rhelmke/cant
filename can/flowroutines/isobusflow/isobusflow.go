package isobusflow

import (
	// "fmt"
	"cant/can/flowroutines"
	"cant/can/isobus"
	"cant/models/pgn"
	"cant/util/globals"
	"sync"
)

// Identifier is a type implementing the flow.Routine interface
type Identifier struct {
	// wow, i didn't know this. If you want to return multiple instances of a basic structure,
	// you need to feed the compiler some trash data. Otherwise it optimizes the process
	// by returning a pointer to a constant struct
	trash byte
}

// CreateIdentifier creates a flow.Routine for the flow package
func CreateIdentifier() *Identifier {
	return &Identifier{}
}

// FlowRun implements the flow.Routine interface for the FlowIdentifier
func (id *Identifier) FlowRun(in <-chan flowroutines.FlowData, out chan<- flowroutines.FlowData, exit chan bool, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		// https://stackoverflow.com/a/11121616
		select {
		case data := <-in:
			data.IsobusFrame = isobus.NewFrame(data.CanFrame)
			//globals.Livelog.Send(fmt.Sprintf("identifier[%p]", id), fmt.Sprintf("Calculated isobus frame out of canbus frame (%p) => {PGN: %d, EDP: %t, DP: %t, PF: %d, PS: %d}", &canFrame, isobusFrame.PGN, isobusFrame.EDP, isobusFrame.DP, isobusFrame.PF, isobusFrame.PS))
			//result := false
			//name := "UNKNOWN"
			if _, ok := pgn.PGNCache[data.IsobusFrame.PGN]; ok {
				//name = pgnResult.Name
				//result = true
				globals.Statistics.AddIdentified <- uint64(1)
			} else {
				globals.Statistics.AddUnknown <- uint64(1)
			}
			//globals.Livelog.Send(fmt.Sprintf("identifier[%p]", id), fmt.Sprintf("Matched PGN of canbus frame (%p) against database => {Identified: %t, Name: %s}", &canFrame, result, name))
			out <- data
			continue
		default:
		}
		select {
		case data := <-in:
			data.IsobusFrame = isobus.NewFrame(data.CanFrame)
			//globals.Livelog.Send(fmt.Sprintf("identifier[%p]", id), fmt.Sprintf("Calculated isobus frame out of canbus frame (%p) => {PGN: %d, EDP: %t, DP: %t, PF: %d, PS: %d}", &canFrame, isobusFrame.PGN, isobusFrame.EDP, isobusFrame.DP, isobusFrame.PF, isobusFrame.PS))
			//result := false
			//name := "UNKNOWN"
			if _, ok := pgn.PGNCache[data.IsobusFrame.PGN]; ok {
				//name = pgnResult.Name
				//result = true
				globals.Statistics.AddIdentified <- uint64(1)
			} else {
				globals.Statistics.AddUnknown <- uint64(1)
			}
			//globals.Livelog.Send(fmt.Sprintf("identifier[%p]", id), fmt.Sprintf("Matched PGN of canbus frame (%p) against database => {Identified: %t, Name: %s}", &canFrame, result, name))
			out <- data
			continue
		case <-exit:
		}
		break
	}
}
