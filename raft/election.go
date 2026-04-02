package raft

func (n *Node) StartElection() []Message {
	n.RoleTransition(CANDIDATE)
	n.NumberOfVotes=0
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

func (n *Node) VoteReceived(msg RequestVoteResponse) []Message {
	if !msg.VoteGranted {
		return nil
	}

	n.NumberOfVotes++

	if n.NumberOfVotes > int(TotalNodesNumber) {
		panic("we have more votes than actual number of nodes")
	}

	if n.NumberOfVotes < int(Quorum) {
		return nil
	}

	messages := n.BecomeLeader()
	return messages
}

func (n *Node) BecomeLeader() []Message {
	n.RoleTransition(LEADER)
	messages := newMessages()

	//TODO: here should be the clean up of timers, like the electionTimeout
	//cleanupTimersMessages()

	for _, v := range n.FriendNodesId {
		messages = append(messages, AppendEntries{
			Sender:   n.Id,
			Receiver: v,

			Term:        n.CurrentTerm,
			CommitIndex: int(n.CommitIndex),

			//TODO: check this indexes, because im not sure if they are restarted
			//or what happens on the transition to becoming a leader, if they start from scratch or what.

			NextIndex:  n.NextIndex[v],
			MatchIndex: n.MatchIndex[v],
		})
	}

	return messages
}
