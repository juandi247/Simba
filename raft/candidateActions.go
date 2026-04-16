package raft

func (n *Node) StartElection() []Message {
	n.RoleTransition(CANDIDATE)
	n.NumberOfVotes = 0
	n.CurrentTerm++
	n.VotedFor = n.Id

	messages := newMessages()
	lastLogIndex := n.Log.Size
	for i := 0; i <= int(TotalNodesNumber-1); i++ {
		messages = append(messages, RequestVote{
			Sender:        n.Id,
			Receiver:      n.FriendNodesId[i],
			Term:          n.CurrentTerm,
			LastLogIndex:  lastLogIndex,
			LastTermIndex: n.Log.LogArr[lastLogIndex].Term,
		})
	}

	/* here should be the command or something, to start the timer.
	i think it could be done as a message again, just put it in the message queue and be processed and thats it
	it would implement the interface from MEsssage so its a valid part, and when processed.
	*/
	return messages
}

func (n *Node) BecomeLeader() []Message {
	n.RoleTransition(LEADER)

	/* TODO: here should be the clean up of timers, like the electionTimeout
	also check if we RESTART the values of
		-NextIndex
		-MatchIndex
	cleanupTimersMessages() */

	//send nil because its just a heartbeat with NO data
	messages := n.buildAppendEntriesMessages(nil)
	return messages
}

