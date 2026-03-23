package coreraft

import (
	"sync/atomic"
)

type Role int

const (
	FOLLOWER Role = iota
	CANDIDATE
	LEADER
)

const NodesNumber uint8 = 5

var Quorum uint8

const MaxLogSize = 10000

func init() {
	Quorum = (NodesNumber / 2) + 1
}

/*
TODO: wrap the node struct in the simulator as SimNode{ node} so the fileds optionals like timeout counter, or leaderheartbeatcounter, etc, they are only on SIMULATOR
*/
type Node struct {
	//Persistent values
	CurrentTerm uint64
	Log         Log
	VotedFor    int

	Id            int                  //this would be probably ip addres, or something, not sure
	FriendNodesId [NodesNumber - 1]int //this would be probaly a list or pool of ip addreses, and in the simulator just the ID to send
	Role          Role
	Leader        int
	NumberOfVotes atomic.Int32
	CommitIndex   uint64

	//key= id , value = the index
	NextIndex       map[int]int
	MatchIndex      map[int]int
	Timeout         uint32
	LeaderHeartbeat uint32
	Adapters        Adapters
}

type Adapters struct {
	TimeAdapter      TimeAdapter
	TransportAdapter TransportAdapter
	StorageAdapter   StorageAdapter
}

type Log struct {
	Size    int
	LogBase []LogBase
}

type LogBase struct {
	Term  int
	Entry string
}

type SimulatorFields struct {
	Timeoutcounter         uint32 //simulator
	LeaderHeartbeatCounter uint32 //simulator
	Alive                  bool
	ComeBackToLiveTick     int64
}

func (n *Node) StartElection() {
	n.RoleTransition(CANDIDATE)
	n.VotedFor = n.Id

	for _, otherNodesId := range n.FriendNodesId {
		n.sendRequestVote(otherNodesId)
	}

}

func (n *Node) RoleTransition(targetRole Role) {

	switch targetRole {
	case FOLLOWER:
	case LEADER:
	case CANDIDATE:
		if n.Role == LEADER {
			// assertion!!
			panic("A Leader can NOT start an election, since he is already the leader")
		}
		//to cleanup everytime there is a new election, to prevent previous wrong
		n.NumberOfVotes.Store(0)
		n.CurrentTerm++
		n.Role = CANDIDATE
	}

}
