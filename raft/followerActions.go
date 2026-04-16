package raft

func (n *Node) handleRequestVote(msg RequestVote) []Message {
	messages := newMessages()
	messages = append(messages, RequestVoteResponse{
		Term:        int(n.CurrentTerm),
		VoteGranted: false,
	})

	if msg.Term <= n.CurrentTerm {
		return messages
	}

	/*TODO: Check for this voted for, because on the jump to state, the voted for can be restarted (?)
	this can be, if its a new term, then i wont have voted for anyone, the check of who i voted for should be onlz on if both match(?)
	*/
	if n.VotedFor != 0 {
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

	n.VotedFor = msg.Sender
	messages[0] = RequestVoteResponse{
		Term:        int(n.CurrentTerm),
		VoteGranted: true,
	}
	return messages
}

func (n *Node) handleAppendEntries(message AppendEntries) []Message {
	//TODO: check this for a middeware easier
	if message.Term < n.CurrentTerm {
		return nil
	}

	if message.Term > n.CurrentTerm {
		n.RoleTransition(FOLLOWER)
	}

	var success = false
	lastEntry := n.Log.Size

	if message.PrevLogIndex <= lastEntry {
		if message.PrevLogTerm == n.Log.LogArr[message.PrevLogIndex].Term {
			success = true
		}
	}

	messages := newMessages()
	messages = append(messages, AppendEntriesResponse{
		Term:    int(n.CurrentTerm),
		Success: success,
	})

	return messages
}
