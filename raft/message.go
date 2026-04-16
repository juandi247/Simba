package raft

type MessageType int

const (
	MsgAppendEntries MessageType = iota
	MsgAppendEntriesResponse
	MsgRequestVote
	MsgRequestVoteResponse
	MsgLeaderTimeout
	MsgHeartbeatTimeout
	MsgElectionTimeout
	MsgNewEntry
)

type Message interface {
	GetType() MessageType
	GetReceiver() int
}

type AppendEntries struct {
	Sender     int
	Receiver   int
	Term       uint64
	LogEntries []LogBase

	CommitIndex int

	PrevLogIndex int
	PrevLogTerm  int
}

type AppendEntriesResponse struct {
	Sender   int
	Receiver int
	Term     int
	Success  bool
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

type ElectionTimeout struct {
}

func (n *Node) ProcessMessage(message Message) []Message {

	switch m := message.(type) {

	//LEADER METHODS
	case NewEntry:
		return n.handleLeaderEntry(m)
	case AppendEntriesResponse:
		return n.HandleAppendEntriesResponse(m)
	case RequestVoteResponse:
		return n.VoteReceived(m)
	case LeaderTimeout:
		//heartbeats
		return n.buildAppendEntriesMessages(nil)

	//CANDIDATE METHODS
	case ElectionTimeout:
		if n.Role != CANDIDATE {
			return n.StartElection()
		}
		return nil

	//FOLLOWER METHODS
	case AppendEntries:
		return n.handleAppendEntries(m)
	case RequestVote:
		return n.handleRequestVote(m)
	case HeartbeatTimeout:
		return n.StartElection()
	default:
		panic("assertion -> a message with unknown type received")
	}

}

/*NOTE: This are the implementations of the interface for all the posible incomming messages*/

func (m AppendEntries) GetType() MessageType {
	return MsgAppendEntries
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

/*
This messages are not used for comunication with the other nodes.
they represent messages that are also procesed in the same single thread but
being events inside the node. Therefore they DONT have term
*/

func (m NewEntry) GetType() MessageType {
	return MsgNewEntry
}

//timeouts
func (m LeaderTimeout) GetType() MessageType {
	return MsgLeaderTimeout
}

func (m HeartbeatTimeout) GetType() MessageType {
	return MsgHeartbeatTimeout
}

func (m ElectionTimeout) GetType() MessageType {
	return MsgElectionTimeout
}

//workaround to get the receiver so that the message queue can deliver the message
func (m AppendEntries) GetReceiver() int {
	return m.GetReceiver()
}

func (m AppendEntriesResponse) GetReceiver() int {
	return m.GetReceiver()
}

func (m RequestVote) GetReceiver() int {
	return m.GetReceiver()
}

func (m RequestVoteResponse) GetReceiver() int {
	return m.GetReceiver()
}

func (m NewEntry) GetReceiver() int {
	return m.GetReceiver()
}

func (m LeaderTimeout) GetReceiver() int {
	return m.GetReceiver()
}
func (m HeartbeatTimeout) GetReceiver() int {
	return m.GetReceiver()
}
func (m ElectionTimeout) GetReceiver() int {
	return m.GetReceiver()
}
