package raft

func (n *Node) AppendToLog(entry NewEntry, term int) error {

	ls := n.Log.Size
	if n.Log.Size >= MaxLogSize {
		panic("the log reached the limit, ending programm")
	}

	ls++
	logEntry := &LogBase{
		Term:  term,
		Entry: entry.Command,
	}
	n.Log.LogArr[ls] = logEntry
	return nil
}

func buildTempLog(entry NewEntry, term int) []LogBase {
	arr := make([]LogBase, 1, 1)

	logEntry := LogBase{
		Term:  int(term),
		Entry: entry.Command,
	}
	arr[0] = logEntry

	return arr
}

func (n *Node) buildAppendEntriesMessages(lb []LogBase) []Message {
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
