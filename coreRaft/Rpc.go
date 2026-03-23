package coreraft

import "fmt"

/*
Remote procedure calls definition of AppendEntries and RequestVOte
*/
func (n *Node) AppendEntries(request AppendEntriesMessage) AppendEntriesResponse {
	if n.VotedFor != 0 {
		// means that we already voted for someone
		return AppendEntriesResponse{
			Term:    int(n.CurrentTerm),
			Success: false,
		}
	}

	n.VotedFor = request.Sender

	//todo: check for the lastlogindex comparation between sender and receiver (?)
	//todo: also check the message term, and react it.

	return AppendEntriesResponse{
		Term:    int(n.CurrentTerm),
		Success: true,
	}
}

/*
Definition of the builder of messages that uses the interface to send the message.
The interface then will depending on the behavior use GRPC(real worldd) or use the simple appendentries method (simulator)
*/
func (n *Node) sendAppendEntries(receiver int) {

	message := AppendEntriesMessage{
		Sender:      n.Id,
		Receiver:    receiver,
		Term:        n.CurrentTerm,
		CommitIndex: int(n.CommitIndex),

		NextIndex:  n.NextIndex[receiver],
		MatchIndex: n.MatchIndex[receiver],
	}
	response := n.Adapters.TransportAdapter.AppendEntries(message)
	fmt.Println(response)
}

func (n *Node) sendRequestVote(receiver int) {

	message := RequestVoteMessage{
		Sender:        n.Id,
		Receiver:      receiver,
		Term:          n.CurrentTerm,
		LastLogIndex:  n.Log.Size,
		LastTermIndex: n.Log.LogBase[n.Log.Size].Term,
	}
	response := n.Adapters.TransportAdapter.RequestVote(message)

	if !response.VoteGranted {
		return
	}

	n.NumberOfVotes.Add(1)
	if n.NumberOfVotes.Load() > int32(NodesNumber) {
		panic("we have more votes than actual number of nodes, there is some race condition or something weird")
	}

	if uint8(n.NumberOfVotes.Load()) >= Quorum {
		n.RoleTransition(LEADER)
	}

}
