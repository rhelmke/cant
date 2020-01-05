package filter

import (
	"cant/can/flowroutines"
	"cant/can/flowroutines/filter/drop"
	"cant/can/flowroutines/filter/manipulate/awgn60928dummy"
	"cant/can/flowroutines/filter/manipulate/awgn65096"
	"cant/can/flowroutines/filter/manipulate/encrypt"
	"cant/can/flowroutines/filter/manipulate/fakeview58880"
	"cant/can/flowroutines/filter/manipulate/fakeview59136"
	"cant/can/flowroutines/filter/manipulate/limit65096"
	"cant/can/flowroutines/filter/manipulate/lowres129025"
	"cant/can/flowroutines/filter/manipulate/tuxawgn65096"
	"cant/can/flowroutines/filter/manipulate/tuxlimit65096"
	"cant/can/flowroutines/filter/manipulate/tuxlowres129025"
	filterModel "cant/models/filter"

	//"cant/util/globals"
	"sync"
)

// DynamicFilter is a type implementing the flow.Routine interface
type DynamicFilter struct {
	trash int
}

// CreateDynamicFilter creates a flow.Routine for the flow package
func CreateDynamicFilter() *DynamicFilter {
	return &DynamicFilter{}
}

// FlowRun implements the flow.Routine interface for the FlowFilter
func (filter *DynamicFilter) FlowRun(in <-chan flowroutines.FlowData, out chan<- flowroutines.FlowData, exit chan bool, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case data := <-in:
			proceed := true
			result := data
			for _, filter := range filterModel.FilterCache {
				for _, pgn := range filter.For {
					if pgn.PGN == data.IsobusFrame.PGN && pgn.Enabled {
						switch filter.ID {
						case 1:
							result, proceed = drop.Run(result)
						case 2:
							result, proceed = awgn65096.Run(result)
						case 3:
							result, proceed = awgn60928dummy.Run(result)
						case 4:
							result, proceed = limit65096.Run(result)
						case 5:
							result, proceed = lowres129025.Run(result)
						case 6:
							result, proceed = fakeview59136.Run(result)
						case 7:
							result, proceed = fakeview58880.Run(result)
						case 8:
							result, proceed = tuxawgn65096.Run(result)
						case 9:
							result, proceed = tuxlimit65096.Run(result)
						case 10:
							result, proceed = tuxlowres129025.Run(result)
						case 11:
							result, proceed = encrypt.Run(result)
						}
						break
					}
					if !proceed {
						break
					}
				}
			}
			if proceed {
				out <- result
			}
			continue
		default:
		}
		select {
		case data := <-in:
			proceed := true
			result := data
			for _, filter := range filterModel.FilterCache {
				for _, pgn := range filter.For {
					if pgn.PGN == data.IsobusFrame.PGN && pgn.Enabled {
						switch filter.ID {
						case 1:
							result, proceed = drop.Run(result)
						case 2:
							result, proceed = awgn65096.Run(result)
						case 3:
							result, proceed = awgn60928dummy.Run(result)
						case 4:
							result, proceed = limit65096.Run(result)
						case 5:
							result, proceed = lowres129025.Run(result)
						case 6:
							result, proceed = fakeview59136.Run(result)
						case 7:
							result, proceed = fakeview58880.Run(result)
						case 8:
							result, proceed = tuxawgn65096.Run(result)
						case 9:
							result, proceed = tuxlimit65096.Run(result)
						case 10:
							result, proceed = tuxlowres129025.Run(result)
						case 11:
							result, proceed = encrypt.Run(result)
						}
						break
					}
					if !proceed {
						break
					}
				}
			}
			if proceed {
				out <- result
			}
			continue
		case <-exit:
		}
		break
	}
}
