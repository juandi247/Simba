package raft

//arbitrary value
const MaxNumberMessages = TotalNodesNumber + 25

func (n *Node) ProcessMessage(message Message) []Message {

switch m := message.(type) {
	//todo: they need to implement the interface, thats why the error appears
	case NewEntry:
		return n.handleLeaderEntry(m)
	case AppendEntries:
	case AppendEntriesResponse:
	case RequestVote:
	case RequestVoteResponse:
		return n.handleRequestVoteResponse(m)
	case HeartbeatTimeout:
		return n.handleHeartbeatTimeout()
	case LeaderTimeout:
		return n.handleLeaderTimeout()
	case ElectionTimeout: 
		return n.handleElectionTimeout()
	default:
		panic("assertion -> a message with unknown type received")
	}

	return nil
}


func (n *Node) handleLeaderEntry(entry NewEntry) []Message {
	err := n.AppendToLog(entry, int(n.CurrentTerm))

	if err != nil {
		panic("some error ocurred inside the append To log")
	}

	tmpLog := buildTempLog(entry, int(n.CurrentTerm))
	messages := n.buildAppendEntriesMessages(tmpLog)
	return messages
}


func (n *Node) HandleAppendEntries(msg AppendEntries) []Message {
	if msg.Term < n.CurrentTerm {
		/*We ignore (or reject) because that guy has an old term*/
	}

	/*
		logic of reading the Commit INdex, and the index recevied, compare it
		if okaz send succes, if not send the other.
		Also appendLOg etcetc.
	*/
	return nil
}

func (n *Node) HandleAppendEntriesResponse(msg AppendEntriesResponse) [TotalNodesNumber]Message {
	if uint64(msg.Term) < n.CurrentTerm {
		/*We ignore (or reject) because that guy has an old term*/
	}

	/*
		logic of reading the Commit INdex, and the index recevied, compare it
		if okaz send succes, if not send the other.
		Also appendLOg etcetc.
	*/
	return [TotalNodesNumber]Message{}
}

func (n *Node) handleRequestVote(msg RequestVote) []Message{
	return []Message{}
}

func (n *Node) handleRequestVoteResponse(msg RequestVoteResponse) []Message {
	return n.VoteReceived(msg)
}


func (n *Node) handleHeartbeatTimeout() []Message {
	return n.StartElection()
}

func (n *Node) handleElectionTimeout() []Message {
	if n.Role!=CANDIDATE {
		return n.StartElection()
	}
	return nil
}

func (n *Node) handleLeaderTimeout() []Message {
	//nil because its just send the heartbeats 
	return n.buildAppendEntriesMessages(nil)
}



