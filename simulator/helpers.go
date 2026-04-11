package simulator

import (
	"math/rand"
	"simba/raft"
)

func isNetworkMessage(message raft.Message) bool {
	/*
		this messages of Timeouts, come from goroutines in the real life, on the same execute, so they dont pass the simulated network fuzzer
	*/
	if message.GetType() == raft.MsgLeaderTimeout ||
		message.GetType() == raft.MsgHeartbeatTimeout ||
		message.GetType() == raft.MsgLeaderTimeout {
		return false
	}
	return true

}

func generateFollowerTimeout(rng *rand.Rand) uint32 {
	return uint32(MinFollowerTimeout + rng.Intn(MaxFollowerTimeout-MinFollowerTimeout+1))
}
