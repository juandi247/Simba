package coreraft

type Role int

const (
	FOLLOWER Role = iota
	CANDIDATE
	LEADER
)

type MessageType int

const (
	HEARTBEAT MessageType=iota
	VOTATION //todo: check this if there is only votation, or start election and votation separately
	APPEND
	ACK
)


const MaxLogSize=10000

type Node struct {
	Id int //this would be probably ip addres, or something, not sure
	FriendNodesId []int //this would be probaly a list or pool of ip addreses, and in the simulator just the ID to send
	Role              Role
	Term              uint64
	Leader            int
	VotedFor          []string
	Log               []string
	CommitIndex       uint64		
	HeartbeatTimeout  uint32
	LeaderHeartbeatTime uint32
	// SIMULATOR ONLY
	Alive 			  bool
	ComeBackToLiveTick int64
}

type Message struct {
	Sender        int
	Receiver      int
	Term          uint64
	MessageType MessageType
	LogIndex      uint64
	// ONLY used for simulator
	DeliveryTick  int64
}

// this limit is to have allways the LIMITS defined for the quantity of messages a node produces for a response or in generall
const MaxMessagesToSend= 50
func (n *Node) Step(msg Message) ([]Message, int){
	msgArr:= make([]Message, MaxMessagesToSend)
	size:=0


	// HERE SHOULD BE THE LOGIC OF THE MESSAGES
	switch msg.MessageType{
		case HEARTBEAT: 
		case VOTATION: 
		case APPEND: 
		case ACK: 

	}


	return msgArr, size
}
