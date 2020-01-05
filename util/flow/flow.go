// Package flow implements a multithreaded flowgraph to rapidly process data through different Routines.
// This package was inspired by gnuradio.
package flow

import (
	"cant/can/flowroutines"
	"fmt"
	"sync"
)

// Routine interface. Each data type implementing this interface can be used as routine
type Routine interface {
	FlowRun(<-chan flowroutines.FlowData, chan<- flowroutines.FlowData, chan bool, *sync.WaitGroup)
}

// Flow Represents a flowgraph
type Flow struct {
	exit          map[Routine]chan bool
	exitMultiplex []chan bool
	started       bool
	stopped       bool
	wg            *sync.WaitGroup
	out           map[Routine]chan flowroutines.FlowData
	in            map[Routine]chan flowroutines.FlowData
	routes        map[chan flowroutines.FlowData][]chan flowroutines.FlowData
	routines      []Routine
}

// New Flow. Each flow can be used exactly once.
func New() *Flow {
	return &Flow{
		exit:          make(map[Routine]chan bool),
		exitMultiplex: []chan bool{},
		started:       false,
		stopped:       false,
		in:            make(map[Routine]chan flowroutines.FlowData),
		out:           make(map[Routine]chan flowroutines.FlowData),
		routes:        make(map[chan flowroutines.FlowData][]chan flowroutines.FlowData),
		wg:            &sync.WaitGroup{},
	}
}

// Connect a flow.Routine src with a destination
func (flow *Flow) Connect(src Routine, dest Routine) error {
	if flow.started || flow.stopped {
		return fmt.Errorf("can not invoke connect, flow has already run")
	}
	for _, routine := range []Routine{src, dest} {
		if _, ok := flow.out[routine]; !ok {
			flow.out[routine] = make(chan flowroutines.FlowData, 1000)
		}
		if _, ok := flow.in[routine]; !ok {
			flow.in[routine] = make(chan flowroutines.FlowData, 1000)
		}
	}
	for _, route := range flow.routes[flow.out[src]] {
		if route == flow.in[dest] {
			return fmt.Errorf("already connected")
		}
	}
	flow.routes[flow.out[src]] = append(flow.routes[flow.out[src]], flow.in[dest])
	srcPresent := false
	destPresent := false
	for _, routine := range flow.routines {
		if routine == src {
			srcPresent = true
		}
		if routine == dest {
			destPresent = true
		}
	}
	if !srcPresent {
		flow.routines = append(flow.routines, src)
		flow.exit[src] = make(chan bool, 2)
	}
	if !destPresent {
		flow.routines = append(flow.routines, dest)
		flow.exit[dest] = make(chan bool, 2)
	}
	return nil
}

// multiplexing channels
func (flow *Flow) multiplex() {
	for k, v := range flow.routes {
		out := k
		ins := v
		if len(ins) == 0 {
			continue
		}
		exit := make(chan bool, 2)
		flow.exitMultiplex = append(flow.exitMultiplex, exit)
		flow.wg.Add(1)
		go func() {
			defer flow.wg.Done()
			for {
				// https://stackoverflow.com/a/11121616
				select {
				case data := <-out:
					for i := range ins {
						ins[i] <- data
					}
					continue
				default:
				}
				select {
				case data := <-out:
					for i := range ins {
						ins[i] <- data
					}
					continue
				case <-exit:
				}
				break
			}
		}()
	}
}

// Start the flow
func (flow *Flow) Start() error {
	if flow.started || flow.stopped {
		return fmt.Errorf("can not start flow since it started/finished before")
	}
	for i := range flow.exit {
		go func(exit chan bool) {
			<-exit
			flow.Stop()
		}(flow.exit[i])
	}
	flow.multiplex()
	for i := range flow.routines {
		flow.wg.Add(1)
		go flow.routines[i].FlowRun(flow.in[flow.routines[i]], flow.out[flow.routines[i]], flow.exit[flow.routines[i]], flow.wg)
	}
	flow.started = true
	return nil
}

// Join the flow
func (flow *Flow) Join() error {
	if flow.started && !flow.stopped {
		flow.wg.Wait()
		return nil
	}
	return fmt.Errorf("can not join flow since it's not running")
}

// Stop the flow
func (flow *Flow) Stop() error {
	if !flow.started || flow.stopped {
		return fmt.Errorf("can not stop flow since it's not running")
	}
	for i := range flow.exitMultiplex {
		flow.exitMultiplex[i] <- true
	}
	for i := range flow.exit {
		flow.exit[i] <- true
	}
	flow.wg.Wait()
	flow.stopped = true
	return nil
}
