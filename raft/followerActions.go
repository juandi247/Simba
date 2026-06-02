package raft

func (n *Node) handleRequestVote(msg RequestVote) []Message {
	messages := newMessages()
	messages = append(messages, RequestVoteResponse{
		Term:        int(n.CurrentTerm),
		VoteGranted: false,
	})

	if msg.Term < n.CurrentTerm {
		return messages
	}

	n.CurrentTerm  = msg.Term
 	
	_, exists:=n.VotedFor[int(n.CurrentTerm)];

	if exists  {
		return messages
	}


	n.RoleTransition(FOLLOWER)

	currLastIndex := n.Log.Size
	currLastTerm := n.Log.LogArr[currLastIndex].Term

	if msg.LastTermIndex < currLastTerm {
		return messages
	}

	if msg.LastLogIndex < currLastIndex {
		return messages
	}

	n.VotedFor[int(n.CurrentTerm)] = msg.Sender
	messages[0] = RequestVoteResponse{
		Term:        int(n.CurrentTerm),
		VoteGranted: true,
	}
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
		if message.PrevLogTerm == n.Log.LogArr[message.PrevLogIndex].Term {
			success = true
			n.Log.LogArr = append(n.Log.LogArr, message.LogEntries...)

			n.Log.Size+= len(message.LogEntries)
		}
	}

	messages := newMessages()
	messages = append(messages, AppendEntriesResponse{
		Term:    int(n.CurrentTerm),
		Success: success,
		LastLogIndex: lastEntryIndex,
	})

	return messages
}
