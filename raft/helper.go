package raft

func newMessages() []Message {
	return make([]Message, 0, MaxNumberMessages)
}

func (n *Node) RoleTransition(targetRole Role) {

	switch targetRole {
	case FOLLOWER:
		n.Role= FOLLOWER
	case LEADER:
		if n.Role==FOLLOWER{
			panic("a follower can not become a leader, without being a candidae first")
		}
		n.Role= LEADER
	case CANDIDATE:
		if n.Role == LEADER {
			// assertion!!
			panic("A Leader can NOT start an election, since he is already the leader")
		}
		n.Role = CANDIDATE
	}
}
