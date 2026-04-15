package raft

import "math"

//arbitrary value
const MaxNumberMessages = TotalNodesNumber + 25

func (n *Node) ProcessMessage(message Message) []Message {

	switch m := message.(type) {

	//LEADER METHODS
	case NewEntry:
		return n.handleLeaderEntry(m)
	case AppendEntriesResponse:
		return n.HandleAppendEntriesResponse(m)
	case RequestVoteResponse:
		return n.handleRequestVoteResponse(m)
	case LeaderTimeout:
		return n.handleLeaderTimeout()

	//CANDIDATE METHODS
	case ElectionTimeout:
		return n.handleElectionTimeout()

	//FOLLOWER METHODS
	case AppendEntries:
		return n.handleAppendEntries(m)
	case RequestVote:
	case HeartbeatTimeout:
		return n.handleHeartbeatTimeout()
	default:
		panic("assertion -> a message with unknown type received")
	}

	return nil
}

func (n *Node) handleLeaderEntry(entry NewEntry) []Message {
	err := n.AppendToLog(entry, int(n.CurrentTerm))

	if err != nil {
		panic("some error ocurred inside the append To log")
	}

	tmpLog := buildTempLog(entry, int(n.CurrentTerm))
	messages := n.buildAppendEntriesMessages(tmpLog)
	return messages
}

func (n *Node) HandleAppendEntriesResponse(msg AppendEntriesResponse) []Message {

	if uint64(msg.Term) < n.CurrentTerm {
		return nil //ignore it
	}
	followerId := msg.Sender

	if !msg.Success {
		n.NextIndex[followerId]--
		return nil
	}

	/*		match index should be updated to the currentValue -1
			nextIndex ++
	*/

	/* 
	if its succesfull, then the matchIndex updates 
*/

	n.MatchIndex[followerId] = n.NextIndex[followerId] -1
	n.NextIndex[followerId]++

	n.checkQuorum()

	return []Message{}
}


func (n *Node) checkQuorum(){
	
	minVal:=math.MaxInt
	quorumCounter:=0
	/*
this will check the quorum for the smallest value, overtime the rest of the values will be checked 
todo: there is the posiblity of executing the for loop, until we reach a nonquorum state. 
using some double for loop (mmmm) or something, to detect matches of three pairs. 


Leetcode type problem: check the biggest posible value, that matches quorum. 

commitIndex= 5
in this example if we have 3,6, 8,9,7

then 3 is ignored, then 6 is the minValue. 
if we reach 9, then there is the count of 3, quorum reached and ends there. 
So we could have searched for this minimum value that completes quorum, that will be 7. so the commitIndex gets updated faster.


posible solutions? a maxHeap of size of the quorum? take the arr[i] that is equal to the quorum -1.
then if the quorum is 3. and we order this in a max heap it will be like: 9,8,7,6,3, so the quorum is 3 
arr[quorum-1] = arr[2] = 7 

could be an aproach(?)
or just a for loop that does this updating the commitIndex and thats it. 

we coudl make a benchmark test for both aproaches. im not sure about it
*/
	for _,value:=range n.MatchIndex{
		if value<int(n.CommitIndex){
			continue
		}

		 if value<minVal{
			minVal=value
		}

		if value>=minVal{
			quorumCounter++
		}

		if quorumCounter>=int(Quorum){
			n.CommitIndex= uint64(minVal)
		}

	}

	/*
the values should be higher than the commitIndex to be part of a valid quorum, because if they are lower
means that its alreadz commited 


Then we evaluate the values and take the most close one to the commitIndex 
commitINdex= 5

[1,2,8,8,7,6]

the ones from 6,7,8 are valid



Traverse the entire matchIndex
Add to a variable of minValue that is always more than the commitIndex 

Everytime we see a value that is bigger or euqal to it, we will have a counter of ++
if we find a smaller value than this (alwazs keeping the rule of being bigger than the commit index, we update it, but keep the ++)


*/


}



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

func (n *Node) handleRequestVoteResponse(msg RequestVoteResponse) []Message {
	return n.VoteReceived(msg)
}

func (n *Node) handleHeartbeatTimeout() []Message {
	return n.StartElection()
}

func (n *Node) handleElectionTimeout() []Message {
	if n.Role != CANDIDATE {
		return n.StartElection()
	}
	return nil
}

func (n *Node) handleLeaderTimeout() []Message {
	//nil because its just send the heartbeats
	return n.buildAppendEntriesMessages(nil)
}
