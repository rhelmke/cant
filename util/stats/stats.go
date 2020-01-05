// Package stats represents basic traffic statistics
package stats

import (
	"encoding/json"
	"sync"
	"time"
)

// Stats holds all statistics
type Stats struct {
	identified     uint64
	unknown        uint64
	blocked        uint64
	manipulated    uint64
	sumTx          uint64
	sumRx          uint64
	ifaces         map[string]*iface
	snap           *Snapshot
	AddInterface   chan string
	AddUnknown     chan uint64
	AddIdentified  chan uint64
	AddManipulated chan uint64
	AddBlocked     chan uint64
	mutex          *sync.Mutex
}

// Snapshot represents copied values of Stats that can be converted to json by Snapshot()
type Snapshot struct {
	Identified  uint64            `json:"identified"`
	Unknown     uint64            `json:"unknown"`
	Blocked     uint64            `json:"blocked"`
	Manipulated uint64            `json:"manipulated"`
	SumTx       uint64            `json:"sum_tx"`
	SumRx       uint64            `json:"sum_rx"`
	Ifaces      map[string]*iface `json:"interfaces"`
	serialized  []byte            `json:"-"`
}

// iface is an internal representation for per-interface statistics
type iface struct {
	In        uint64      `json:"in_pkts"`
	Out       uint64      `json:"out_pkts"`
	AddInPkt  chan uint64 `json:"-"`
	AddOutPkt chan uint64 `json:"-"`
}

// New Stats object
func New() *Stats {
	stats := &Stats{ifaces: make(map[string]*iface),
		snap:           &Snapshot{},
		AddInterface:   make(chan string, 1000),
		AddUnknown:     make(chan uint64, 1000),
		AddIdentified:  make(chan uint64, 1000),
		AddManipulated: make(chan uint64, 1000),
		AddBlocked:     make(chan uint64, 1000),
		mutex:          &sync.Mutex{},
	}
	go stats.updateRoutine()
	return stats
}

func (s *Stats) GetInChannel(iface string) chan<- uint64 {
	if _, ok := s.ifaces[iface]; !ok {
		return nil
	}
	return s.ifaces[iface].AddInPkt
}

func (s *Stats) GetOutChannel(iface string) chan<- uint64 {
	if _, ok := s.ifaces[iface]; !ok {
		return nil
	}
	return s.ifaces[iface].AddOutPkt
}

func (s *Stats) updateRoutine() {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case name := <-s.AddInterface:
			if _, ok := s.ifaces[name]; !ok {
				s.ifaces[name] = &iface{AddInPkt: make(chan uint64, 1000), AddOutPkt: make(chan uint64, 1000)}
				go func() {
					iface := name
					for {
						select {
						case in := <-s.ifaces[iface].AddInPkt:
							s.mutex.Lock()
							s.ifaces[iface].In += in
							s.sumRx += in
							s.mutex.Unlock()
						case out := <-s.ifaces[iface].AddOutPkt:
							s.mutex.Lock()
							s.ifaces[iface].Out += out
							s.sumTx += out
							s.mutex.Unlock()
						}
					}
				}()
			}
		case add := <-s.AddUnknown:
			s.mutex.Lock()
			s.unknown += add
			s.mutex.Unlock()
		case add := <-s.AddIdentified:
			s.mutex.Lock()
			s.identified += add
			s.mutex.Unlock()
		case add := <-s.AddManipulated:
			s.mutex.Lock()
			s.manipulated += add
			s.mutex.Unlock()
		case add := <-s.AddBlocked:
			s.mutex.Lock()
			s.blocked += add
			s.mutex.Unlock()
		case <-ticker.C:
			s.mutex.Lock()
			s.snap.Identified = s.identified
			s.snap.Blocked = s.blocked
			s.snap.Ifaces = s.ifaces
			s.snap.Manipulated = s.manipulated
			s.snap.SumRx = s.sumRx
			s.snap.SumTx = s.sumTx
			s.snap.Unknown = s.unknown
			s.snap.serialized, _ = json.Marshal(s.snap)
			s.mutex.Unlock()
		}
	}
}

// Snapshot returns a snapshot of the current statistics in JSON format
func (s *Stats) Snapshot() []byte {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.snap.serialized
}
