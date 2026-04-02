package raft

func newMessages() []Message {
	return make([]Message, 0, MaxNumberMessages)
}

func (n *Node) RoleTransition(targetRole Role) {

	switch targetRole {
	case FOLLOWER:
	case LEADER:
	case CANDIDATE:
		if n.Role == LEADER {
			// assertion!!
			panic("A Leader can NOT start an election, since he is already the leader")
		}
		n.Role = CANDIDATE
	}
}
