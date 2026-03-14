package coreraft

type Role int

const (
	FOLLOWER Role = iota
	CANDIDATE
	LEADER
)

const MaxLogSize = 10000

type Node struct {
	Id            int   //this would be probably ip addres, or something, not sure
	FriendNodesId []int //this would be probaly a list or pool of ip addreses, and in the simulator just the ID to send
	Role          Role
	Term          uint64
	Leader        int
	VotedFor      []string
	Log           []string //this is in memory LOG, this shuold be also by the network simlator or fuzzer modified time (BECAUSE its IO undeterminsitc)
	CommitIndex   uint64

	Timeout        uint32
	Timeoutcounter uint32 //simulator

	LeaderHeartbeat        uint32
	LeaderHeartbeatCounter uint32 //simulator
	// SIMULATOR ONLY
	Alive              bool
	ComeBackToLiveTick int64
}



/*This is the procesing of messages */
func (n *Node) ProcessMesage(msg Message, transportAdapter TransportAdapter, timeAdapter TimeAdapter) {
	/*
		the logic of the transport of messages or delivery is inside the interface, so it gets used inside the core wihtout knowing th implmenetation
		this is just for testing (for now) but i think this could be the best option
	*/
	switch msg.MessageType {
	case HEARTBEAT:
		if msg.Term < n.Term{
			// means we have an OLD message, in this case we can send a message to the node, saying hey the current term now is this one
			//n.SendMesage("hey you need to update your term", "hits is the leader for the term")
			return
		}


		//here should be the logic to reestart the timers
		timeAdapter.RestartTimeoutTimer(n, )



	case ELECTION:
	case APPEND:
	case ACK:
		transportAdapter.SendMessage(Message{})

	}
}

func (n *Node) Tick(timeAdapter TimeAdapter, TransportAdapter TransportAdapter) {
	timeAdapter.DetermineTimeouts(n, TransportAdapter)
}

func (n *Node) SendMesage(messageType MessageType, receiver, LogIndex int, TransportAdapter TransportAdapter) {

	message := Message{
		Sender:      n.Id,
		Receiver:    receiver,
		Term:        n.Term,
		MessageType: messageType,
		LogIndex:    uint64(LogIndex),
		//mesage delivery tick will be modified by the implmentation (in the case of sim) on the real its NOT used
	}
	TransportAdapter.SendMessage(message)
}



// func (n *Node) 
