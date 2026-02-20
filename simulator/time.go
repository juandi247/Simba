package simulator


type SimTime struct {
	Tick int64
}

func (st *SimTime) Now() int64 {
	return st.Tick
}
func (st *SimTime) Advance(delta int64) {
	st.Tick+= delta
}
func (st *SimTime) Sleep(int64) {
}