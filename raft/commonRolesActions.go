package raft

import "fmt"

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
	fmt.Println("printing myself", n.Id)
	messages := newMessages()
	messages = append(messages, HeartbeatTimeout{
		Receiver: n.Id,
	})
	return messages
}

func (n *Node) TriggerHeartbeat() []Message {
	messages := newMessages()
	messages = append(messages, LeaderTimeout{
		Receiver: n.Id,
	})
	return messages
}

func (n *Node) TriggerElectionTimeout() []Message {
	messages := newMessages()
	messages = append(messages, HeartbeatTimeout{
		Receiver: n.Id,
	})
	return messages
}
