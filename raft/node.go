package raft



type Role int
const (
	FOLLOWER Role = iota
	CANDIDATE
	LEADER
)

const TotalNodesNumber uint8 = 5

var Quorum uint8

const MaxLogSize int = 10000

func init() {
	Quorum = (TotalNodesNumber / 2) + 1
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
	FriendNodesId [TotalNodesNumber - 1]int //this would be probaly a list or pool of ip addreses, and in the simulator just the ID to send
	Role          Role
	Leader        int
	NumberOfVotes int
	CommitIndex   uint64

	//key= id , value = the index
	NextIndex       map[int]int
	MatchIndex      map[int]int
	Timeout         uint32
	LeaderHeartbeat uint32
	ElectionTimeout uint32

	SimulatorFields *SimulatorFields
}

type Log struct {
	Size    int
	LogArr []*LogBase
}

type LogBase struct {
	Term  int
	Entry string
}

type SimulatorFields struct {
	Timeoutcounter         uint32 
	LeaderHeartbeatCounter uint32 
	Alive                  bool
	ComeBackToLiveTick     int64
	ElectionTimeoutCounter uint32
}


