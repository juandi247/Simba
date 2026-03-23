package coreraft

type AppendEntriesMessage struct {
	Sender   int
	Receiver int
	Term     uint64

	CommitIndex int
	LastApplied int //todo: check if this one is worth it or not. dont thinkso but ok

	NextIndex  int
	MatchIndex int
	// ONLY used for simulator
	DeliveryTick int64
}

type AppendEntriesResponse struct {
	Term    int
	Success bool
}

type RequestVoteMessage struct {
	Sender   int
	Receiver int
	Term     uint64

	LastLogIndex  int
	LastTermIndex int
	// ONLY used for simulator
	DeliveryTick int64
}

type RequestVoteResponse struct {
	Term        int
	VoteGranted bool
}
