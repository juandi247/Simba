package raft

import "fmt"

func (n *Node) handleRequestVote(msg RequestVote) []Message {
	messages := newMessages()
	messages = append(messages, RequestVoteResponse{
		Sender: n.Id,
		Receiver: msg.Sender,
		Term:        int(n.CurrentTerm),
		VoteGranted: false,
	})

	if msg.Term < n.CurrentTerm {
		fmt.Println("Vote not granted, msg term is smaller than currentTerm")
		return messages
	}

	n.CurrentTerm = msg.Term

	votedFor, exists := n.VotedFor[int(n.CurrentTerm)]

	if exists {
		fmt.Println("Vote not granted, ALREADY voted for: ", votedFor, "in term: ", n.CurrentTerm)
		return messages
	}

	n.RoleTransition(FOLLOWER)
	fmt.Println("transicion bien")

	currLastIndex := n.Log.Size
	//NOTE: this is for validation when the log is cero, because is just starting
	if currLastIndex==0{
		n.VotedFor[int(n.CurrentTerm)] = msg.Sender
		messages[0] = RequestVoteResponse{
			Sender: n.Id,
			Receiver: msg.Sender,
			Term:        int(n.CurrentTerm),
			VoteGranted: true,
		}
		fmt.Println("TERM: ", n.CurrentTerm, "VotedFor: ", msg.Sender)
		return messages
	}
	currLastTerm := n.Log.LogArr[currLastIndex].Term

	if msg.LastTermIndex < currLastTerm {
		fmt.Println("TERM: ", n.CurrentTerm, "RejectedVOte from: ", msg.Sender)
		return messages
	}

	if msg.LastLogIndex < currLastIndex {
		fmt.Println("TERM: ", n.CurrentTerm, "RejectedVOte from: ", msg.Sender)
		return messages
	}

	n.VotedFor[int(n.CurrentTerm)] = msg.Sender
	messages[0] = RequestVoteResponse{
		Sender: n.Id,
		Receiver: msg.Sender,
		Term:        int(n.CurrentTerm),
		VoteGranted: true,
	}
	fmt.Println("TERM: ", n.CurrentTerm, "VotedFor: ", msg.Sender)
	return messages
}

func (n *Node) handleAppendEntries(message AppendEntries) []Message {

	if message.Term < n.CurrentTerm {
		return nil
	}

	if message.Term > n.CurrentTerm {
		n.RoleTransition(FOLLOWER)
	}

	var success = false
	lastEntryIndex := n.Log.Size

	if message.PrevLogIndex <= lastEntryIndex {
		if message.PrevLogIndex>=0 && message.PrevLogTerm == n.Log.LogArr[message.PrevLogIndex].Term {
			success = true
			n.Log.LogArr = append(n.Log.LogArr, message.LogEntries...)
			n.Log.Size += len(message.LogEntries)
		}
	}

	messages := newMessages()
	messages = append(messages, AppendEntriesResponse{
		Sender: n.Id,
		Receiver: message.Sender,
		Term:         int(n.CurrentTerm),
		Success:      success,
		LastLogIndex: lastEntryIndex,
	})

	return messages
}
