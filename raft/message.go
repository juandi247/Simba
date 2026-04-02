package raft

type MessageType int

const (
	MsgAppendEntries MessageType = iota
	MsgAppendEntriesResponse
	MsgRequestVote
	MsgRequestVoteResponse
	MsgLeaderTimeout
	MsgHeartbeatTimeout
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
	LogEntries []LogBase

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

type LeaderTimeout struct {
}

type HeartbeatTimeout struct {
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
func (m AppendEntriesResponse) GetTerm() (int, bool) {
	return int(m.Term), true
}

func (m RequestVote) GetType() MessageType {
	return MsgRequestVote
}
func (m RequestVote) GetTerm() (int, bool) {
	return int(m.Term), true
}

func (m RequestVoteResponse) GetType() MessageType {
	return MsgRequestVoteResponse
}
func (m RequestVoteResponse) GetTerm() (int, bool) {
	return int(m.Term), true
}

/*
This messages are not used for comunication with the other nodes.
they represent messages that are also procesed in the same single thread but
being events inside the node. Therefore they DONT have term
*/

func (m NewEntry) GetType() MessageType {
	return MsgNewEntry
}
func (m NewEntry) GetTerm() (int, bool) {
	return 0, false
}

//timeouts
func (m LeaderTimeout) GetType() MessageType {
	return MsgLeaderTimeout
}
func (m LeaderTimeout) GetTerm() (int, bool) {
	return 0, false
}

func (m HeartbeatTimeout) GetType() MessageType {
	return MsgHeartbeatTimeout
}
func (m HeartbeatTimeout) GetTerm() (int, bool) {
	return 0, false
}
