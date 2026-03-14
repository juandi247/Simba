package simulator

import (
	coreraft "simba/coreRaft"
)
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



func (st *SimTime) DetermineTimeouts(n *coreraft.Node, TransportAdapter coreraft.TransportAdapter){
/*
	   si el tiempo del lider llego a cero o es menor a cero, reiniciar y enviar el heartbeat
	*/
	if n.Role == coreraft.LEADER && n.LeaderHeartbeatCounter <= 0 {
		// send a heartbeat message to everyone
		// Node.sendHeartbeat(transportAdapter)


		//restart the time counter for the heartbeat itself
		n.LeaderHeartbeatCounter = n.LeaderHeartbeat

	}


	if n.Role == coreraft.FOLLOWER && n.Timeoutcounter<=0{

		//Node.StartElection()
		//ChangeRole to CANDIDATE. 
		//this will probalby call another function, but this should be on core raft
		n.Timeoutcounter = n.Timeout

	}

}
