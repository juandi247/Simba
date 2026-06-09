package simulator

import (
	"fmt"
	"math/rand"
	"simba/raft"
)


func buildFriendsIds(numberOfFriends int, currId int) []int{
	rta:= []int{}
	for i:=1; i<=numberOfFriends; i++{
		if i!=currId{
		rta = append(rta, i)
		}

	}
	return rta
}


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
	t:= uint32(MinFollowerTimeout + rng.Intn(MaxFollowerTimeout-MinFollowerTimeout+1))
	fmt.Println("timeout generado es: ", t) 
	return t
}




