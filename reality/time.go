package reality

import "time"

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