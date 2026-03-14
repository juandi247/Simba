package reality

import (
	coreraft "simba/coreRaft"
	"time"
)

type PhysicalTime struct {
}

func (st *PhysicalTime) Now() int64 {
	currTime := time.Now()
	return currTime.Unix()
}
func (st *PhysicalTime) Advance(delta int64) {
}

func (st *PhysicalTime) Sleep(ms int64) {
	time.Sleep(time.Millisecond * time.Duration(ms))
}
func (st *PhysicalTime) DetermineTimeouts(node *coreraft.Node, transportAdapter coreraft.TransportAdapter) {
	// time.Sleep(time.Millisecond * time.Duration(ms))
}