package raft

import "fmt"

func (n *Node) StartElection() []Message {
	fmt.Println("Starting election from node: ", n.Id)
	n.RoleTransition(CANDIDATE)
	n.NumberOfVotes = 0
	n.CurrentTerm++
	//NOTE: Votes for itself
	n.VotedFor[int(n.CurrentTerm)] = n.Id

	messages := newMessages()
	lastLogIndex := n.Log.Size
	fmt.Println("curr id", n.Id)
	for _,friendId:=range n.FriendNodesId {
		lastTermIdx:= 0
		if n.Log.LogArr[lastLogIndex]!=nil{
			lastTermIdx =  n.Log.LogArr[lastLogIndex].Term

		}
		fmt.Println("node id to send is: ", friendId)
		messages = append(messages, RequestVote{
			Sender:        n.Id,
			Receiver:      friendId,
			Term:          n.CurrentTerm,
			LastLogIndex:  lastLogIndex,
			LastTermIndex: lastTermIdx,
		})
	}

	/* here should be the command or something, to start the timer.
	i think it could be done as a message again, just put it in the message queue and be processed and thats it
	it would implement the interface from MEsssage so its a valid part, and when processed.
	*/
	return messages
}


func (n *Node) HandleRequestVoteResponse(msg RequestVoteResponse) []Message {
	if !msg.VoteGranted {
		return nil
	}

	n.NumberOfVotes++
	fmt.Println("NODE: ", n.Id, "has: ", n.NumberOfVotes, " Votes")

	if n.NumberOfVotes > int(TotalNodesNumber) {
		panic("we have more votes than actual number of nodes")
	}

	if n.NumberOfVotes < int(Quorum) {
		fmt.Println("havent reached quorum")
		return nil
	}
	
	fmt.Println("reached QUORUM")
	messages := n.BecomeLeader()
	return messages
}


func (n *Node) BecomeLeader() []Message {
	fmt.Println("we are leaders NODE: ", n.Id)
	n.RoleTransition(LEADER)

	/* TODO: here should be the clean up of timers, like the electionTimeout
	also check if we RESTART the values of
		-NextIndex
		-MatchIndex
	cleanupTimersMessages() */

	//send nil because its just a heartbeat with NO data
	messages := n.buildAppendEntries(nil)
	return messages
}
