package raft


func (n *Node) AppendToLog(entry NewEntry) []Message {

	ls := n.Log.Size
	if n.Log.Size >= MaxLogSize {
		panic("the log reached the limit, ending programm")
	}

	ls++

	logEntry := LogBase{
		Term:  int(n.CurrentTerm),
		Entry: entry.Command,
	}
	n.Log.LogArr[ls] = &logEntry

	msgArr := newMessages()
	for _, receiverId := range n.FriendNodesId {
		msgArr = append(msgArr, AppendEntries{
			Sender:      n.Id,
			Receiver:    receiverId,
			//TODO: add the log entries
			//LogEntries: array of new entries,

			Term:        n.CurrentTerm,
			CommitIndex: int(n.CommitIndex),

			NextIndex:  n.NextIndex[receiverId],
			MatchIndex: n.MatchIndex[receiverId],
		})
	}

	return msgArr
}



