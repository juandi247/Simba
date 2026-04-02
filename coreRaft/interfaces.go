package coreraft

type MessageType int

const (
	MsgAppendEntries MessageType = iota
	MsgAppendEntriesResponse
	MsgRequestVote
	MsgRequestVoteResponse
	MsgLeaderTimeout
	MsgFollowerHeartbeatTimeout

	MsgNewEntry
)

type Message interface {
	GetType() MessageType
	/* bool represents here if the message contains a Term, help us ignore the validation of term on some messages
	example: the new entry does not have a term, because comes directly from the client
	*/
	GetTerm() (int, bool)
}

type AppendEntries struct {
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

type RequestVote struct {
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

type NewEntry struct {
	Command string
}

/*
Implement the message interface for all the messages to be received
*/
func (m AppendEntries) GetType() MessageType {
	return MsgAppendEntries
}
func (m AppendEntries) GetTerm() (int, bool) {
	return int(m.Term), true
}

func (m AppendEntriesResponse) GetType() MessageType {
	return MsgAppendEntriesResponse
}
func (m RequestVote) GetType() MessageType {
	return MsgRequestVote
}
func (m RequestVoteResponse) GetType() MessageType {
	return MsgRequestVoteResponse
}

func (m NewEntry) GetType() MessageType {
	return MsgNewEntry
}

//TODO: here are missing the timeouts messages
