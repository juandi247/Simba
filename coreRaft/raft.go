package coreraft

import "fmt"

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
todo: wrap the node struct in the simulator as SimNode{ node} so the fileds optionals like timeout counter, or leaderheartbeatcounter, etc, they are only on SIMULATOR
*/
type Node struct {
	Id            int                  //this would be probably ip addres, or something, not sure
	FriendNodesId [NodesNumber - 1]int //this would be probaly a list or pool of ip addreses, and in the simulator just the ID to send
	Role          Role
	Term          uint64
	Leader        int

	VotedFor      int
	NumberOfVotes int
	Log           []string //this is in memory LOG, this shuold be also by the network simlator or fuzzer modified time (BECAUSE its IO undeterminsitc)
	CommitIndex   uint64

	Timeout        uint32
	Timeoutcounter uint32 //simulator

	LeaderHeartbeat        uint32
	LeaderHeartbeatCounter uint32 //simulator
	// SIMULATOR ONLY
	Alive              bool
	ComeBackToLiveTick int64

	// interfaces for implementations depending on simulator or real life
	TimeAdapter      TimeAdapter
	TransportAdapter TransportAdapter
	StorageAdapter   StorageAdapter
}


func (n *Node) ProcessMesage(msg Message) {

	if msg.Term < n.Term {
		fmt.Println("Hey we received an old message.")
		n.SendMesage(UPDATEYOURDATA, msg.Sender, 0)
		return
	}


	if msg.Term > n.Term{
		n.Term= msg.Term
		n.Role = FOLLOWER
		n.VotedFor=0
		n.NumberOfVotes=0
		// todo: here weshould have a reset on timeouts (?)
	}

	switch msg.MessageType {

	case VOTECONFIRMATION:
		if n.Role!=CANDIDATE{
			return
		}

		// todo: check if this is also for a Mutex, because this will be running on waiting for RPCs, so they can be 2 concurrnet updates or ++
		n.NumberOfVotes++
		if n.NumberOfVotes>= int(Quorum){
			n.Role= LEADER
			//todo:  reset timers or start timer for the leader
			for _, receiver:= range n.FriendNodesId{
				n.SendMesage(HEARTBEAT, receiver, 0)
			}
		}

	case REQUESTVOTE:

		if n.VotedFor !=0{
			// means that we already voted for someone
			return
		}


		//todo: check for the lastlogindex comparation between sender and receiver (?)
		n.VotedFor = msg.Sender
		n.SendMesage(VOTECONFIRMATION, msg.Sender, 0)

	case HEARTBEAT:
		//todo: reset the timeout imelmentation (we already checked that the term is valid)
	case APPEND:
	case ACK:
	default:
		fmt.Println("someone is sending random data(????)")
	}
}

func (n *Node) SendMesage(messageType MessageType, receiver, LogIndex int) {

	message := Message{
		Sender:      n.Id,
		Receiver:    receiver,
		Term:        n.Term,
		MessageType: messageType,
		LogIndex:    uint64(LogIndex),
		//mesage delivery tick will be modified by the implmentation (in the case of sim) on the real its NOT used
	}
	n.TransportAdapter.SendMessage(message)
}



func (n *Node) StartElection() {
	if n.Role == LEADER {
		// assertion!!
		panic("A Leader can NOT start an election, since he is already the leader")
	}
	//to cleanup everytime there is a new election, to prevent previous wrong 
	n.NumberOfVotes=0
	n.Term++
	n.Role = CANDIDATE

	
	n.VotedFor = n.Id

	for _, otherNodesId := range n.FriendNodesId {
		n.SendMesage(REQUESTVOTE, otherNodesId, 0)
	}

}
