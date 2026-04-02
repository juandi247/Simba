package coreraft

import "fmt"

func (n *Node) messageFSM(message Message) ([TotalNodesNumber]Message, int) {
	//we define the array of mesages to send with explicit
	messagesToSend := make([]Message, TotalNodesNumber)
	size := 0

	switch m := message.(type) {
		//todo: they need to implement the interface, thats why the error appears
	case NewEntry:
		size = n.handleNewEntry(m, messagesToSend)
	case AppendEntries:
	case AppendEntriesResponse:
	case RequestVote:
	case RequestVoteResponse:
	default:
		panic("assertion -> a message with unknown type received")
	}

	if size > int(TotalNodesNumber) {
		panic("assertion, the size of the messages to send is bigger than the number of nodes itself this should not happen")
	}
	return messagesToSend, size
}

func (n *Node) handleNewEntry(entry NewEntry, responeMessages []Message) int {

	ls := n.Log.Size
	if n.Log.Size >= MaxLogSize {
		panic("the log reached the limit, ending programm")
	}

	ls++

	logEntry := LogBase{
		Term:  int(n.CurrentTerm),
		Entry: entry.Command,
	}
	n.Log.LogArr[ls] = logEntry

	for i, receiverId := range n.FriendNodesId {
		responeMessages[i] = AppendEntries{
			Sender:      n.Id,
			Receiver:    receiverId,
			Term:        n.CurrentTerm,
			CommitIndex: int(n.CommitIndex),

			NextIndex:  n.NextIndex[receiverId],
			MatchIndex: n.MatchIndex[receiverId],
		}
	}

	return int(TotalNodesNumber)
}



func (n *Node) HandleAppendEntries(msg AppendEntries) [TotalNodesNumber]Message {
	if msg.Term < n.CurrentTerm {
		/*We ignore (or reject) because that guy has an old term*/
	}

	/*
		logic of reading the Commit INdex, and the index recevied, compare it
		if okaz send succes, if not send the other.
		Also appendLOg etcetc.
	*/
	return [TotalNodesNumber]Message{}
}



func (n *Node) HandleAppendEntriesResponse(msg AppendEntriesResponse) [TotalNodesNumber]Message {
	if msg.Term < n.CurrentTerm {
		/*We ignore (or reject) because that guy has an old term*/
	}

	/*
		logic of reading the Commit INdex, and the index recevied, compare it
		if okaz send succes, if not send the other.
		Also appendLOg etcetc.
	*/
	return [TotalNodesNumber]Message{}
}




func (n *Node) handleRequestVote(msg RequestVote) {

}




//
// func (n *Node) sendRequestVote(receiver int) {
//
// 	message := RequestVoteMessage{
// 		Sender:        n.Id,
// 		Receiver:      receiver,
// 		Term:          n.CurrentTerm,
// 		LastLogIndex:  n.Log.Size,
// 		LastTermIndex: n.Log.LogBase[n.Log.Size].Term,
// 	}
// 	response := n.Adapters.TransportAdapter.RequestVote(message)
//
// 	if !response.VoteGranted {
// 		return
// 	}
//
// 	n.NumberOfVotes.Add(1)
// 	if n.NumberOfVotes.Load() > int32(NodesNumber) {
// 		panic("we have more votes than actual number of nodes, there is some race condition or something weird")
// 	}
//
// 	if uint8(n.NumberOfVotes.Load()) >= Quorum {
// 		n.RoleTransition(LEADER)
// 	}
//
// }
