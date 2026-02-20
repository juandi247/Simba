package coreraft


type Node struct {
	Role              string
	Term              int
	Leader            string
	VotedFor          []string
	Log               []string
	CommitIndex       int
	Heartbeat_Timeout int
}

type Message struct {
	Sender        string
	Receiver      string
	Term          int
	TypeOfMessage int
	LogIndex      int
}

