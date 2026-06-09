package raft

import "fmt"

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
	LogEntries []*LogBase

	CommitIndex int

	PrevLogIndex int
	PrevLogTerm  int
}

type AppendEntriesResponse struct {
	Sender       int
	Receiver     int
	Term         int
	LastLogIndex int
	Success      bool
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
	Sender      int
	Receiver    int
	Term        int
	VoteGranted bool
}

type NewEntry struct {
	Receiver int
	Command  string
}

type LeaderTimeout struct {
	Receiver int
}

type HeartbeatTimeout struct {
	Receiver int
}

type ElectionTimeout struct {
	Receiver int
}

func (n *Node) ProcessMessage(message Message) []Message {

fmt.Println("---- PROCESING MESSAGE NODE: ", n.Id, "-------")
	defer fmt.Println("----End Procesing Message--- \n ")
//	defer fmt.Println(" ----")
	switch m := message.(type) {

	//LEADER METHODS
	case NewEntry:
		fmt.Println("recevied new Entry")
		return n.handleLeaderEntry(m)
	case AppendEntriesResponse:
		fmt.Println("recevied appendEntries respons")
		return n.HandleAppendEntriesResponse(m)
	case RequestVoteResponse:
		fmt.Println("recevied requestVote respone")
		return n.HandleRequestVoteResponse(m)
	case LeaderTimeout:
		fmt.Println("recevied lead timotu")
		//heartbeats
		return n.buildAppendEntries(nil)

	//CANDIDATE METHODS
	case ElectionTimeout:
		if n.Role != CANDIDATE {
			fmt.Println("received election timeout")
			return n.StartElection()
		}
		return nil

	//FOLLOWER METHODS
	case AppendEntries:
		fmt.Println("recevied appendEntries")
		return n.handleAppendEntries(m)
	case RequestVote:
		fmt.Println("recevied request vote in node ", n.Id, "from: ", m.Sender)
		return n.handleRequestVote(m)
	case HeartbeatTimeout:
		fmt.Println("recevied hb timeout")
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

// timeouts
func (m LeaderTimeout) GetType() MessageType {
	return MsgLeaderTimeout
}

func (m HeartbeatTimeout) GetType() MessageType {
	return MsgHeartbeatTimeout
}

func (m ElectionTimeout) GetType() MessageType {
	return MsgElectionTimeout
}

// workaround to get the receiver so that the message queue can deliver the message
func (m AppendEntries) GetReceiver() int {
	return m.Receiver
}

func (m AppendEntriesResponse) GetReceiver() int {
	return m.Receiver
}

func (m RequestVote) GetReceiver() int {
	return m.Receiver
}

func (m RequestVoteResponse) GetReceiver() int {
	return m.Receiver
}

func (m NewEntry) GetReceiver() int {
	return m.Receiver
}

func (m LeaderTimeout) GetReceiver() int {
	return m.Receiver
}
func (m HeartbeatTimeout) GetReceiver() int {
	return m.Receiver
}
func (m ElectionTimeout) GetReceiver() int {
	return m.Receiver
}
