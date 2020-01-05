package flowroutines

import (
	"cant/can/canbus"
	"cant/can/isobus"
)

// FlowData ...
type FlowData struct {
	IsobusFrame  isobus.Frame
	CanFrame     canbus.Frame
	SrcInterface int
	DstInterface int
}

type Endpoint interface {
	DefaultRoute(Endpoint)
	GetDefaultRoute() int
	Sleep()
	Wake()
	IsAsleep() bool
}
