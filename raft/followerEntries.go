package raft

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
