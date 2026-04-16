package raft

func (n *Node) handleLeaderEntry(entry NewEntry) []Message {
	err := n.AppendToLog(entry, int(n.CurrentTerm))

	if err != nil {
		panic("some error ocurred inside the append To log")
	}

	tmpLog := buildTempLog(entry, int(n.CurrentTerm))
	messages := n.buildAppendEntries(tmpLog)
	return messages
}

func (n *Node) HandleAppendEntriesResponse(msg AppendEntriesResponse) []Message {

	if uint64(msg.Term) < n.CurrentTerm {
		return nil //ignore it
	}
	followerId := msg.Sender

	if !msg.Success {
		n.NextIndex[followerId]--
		return nil
	}

	n.MatchIndex[followerId] = n.NextIndex[followerId] - 1
	n.NextIndex[followerId]++

	checkEntryQuorum(n)

	return []Message{}
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

func (n *Node) buildAppendEntries(lb []LogBase) []Message {
	if n.Role!=LEADER{
		panic("non leader wants to send a hearbeat or appendentries")
	}
	messages := newMessages()
	for _, followerId := range n.FriendNodesId {
		prevLogIndex := n.NextIndex[followerId] - 1
		prevLogTerm := n.Log.LogArr[prevLogIndex].Term

		messages = append(messages, AppendEntries{
			Sender:   n.Id,
			Receiver: followerId,

			Term:        n.CurrentTerm,
			CommitIndex: int(n.CommitIndex),

			LogEntries: lb,

			PrevLogIndex: prevLogIndex,
			PrevLogTerm:  prevLogTerm,
		})
	}

	return messages

}
