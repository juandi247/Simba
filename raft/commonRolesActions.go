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

func (n *Node) TriggerTimeout() []Message {
	messages := newMessages()
	messages = append(messages, LeaderTimeout{})
	return messages
}

func (n *Node) TriggerHeartbeat() []Message {
	messages := newMessages()
	messages = append(messages, HeartbeatTimeout{})
	return messages
}

func (n *Node) TriggerElectionTimeout() []Message {
	messages := newMessages()
	messages = append(messages, HeartbeatTimeout{})
	return messages
}
