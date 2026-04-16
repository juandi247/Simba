package raft

import (
	"cmp"
	"fmt"
	"slices"
)

const MaxNumberMessages = TotalNodesNumber + 25

func newMessages() []Message {
	return make([]Message, 0, MaxNumberMessages)
}

func (n *Node) RoleTransition(targetRole Role) {

	switch targetRole {
	case FOLLOWER:
		n.Role = FOLLOWER
	case LEADER:
		if n.Role == FOLLOWER {
			panic("a follower can not become a leader, without being a candidae first")
		}
		n.Role = LEADER
	case CANDIDATE:
		if n.Role == LEADER {
			// assertion!!
			panic("A Leader can NOT start an election, since he is already the leader")
		}
		n.Role = CANDIDATE
	}
}

func buildTempLog(entry NewEntry, term int) []LogBase {
	arr := make([]LogBase, 1)

	logEntry := LogBase{
		Term:  int(term),
		Entry: entry.Command,
	}
	arr[0] = logEntry

	return arr
}

type MapNode struct {
	key int
	val int
}

func createSliceFromMap(mapita map[int]int) []MapNode {
	slice := make([]MapNode, 0, len(mapita))
	for k, v := range mapita {
		slice = append(slice, MapNode{
			key: k,
			val: v,
		})
	}

	slices.SortStableFunc(slice, func(a, b MapNode) int {
		return cmp.Compare(b.val, a.val)
	})
	return slice
}

func checkEntryQuorum(n *Node) {

	arr := createSliceFromMap(n.MatchIndex)
	newCommitIndex, err := setUpFinalValueForCommitIndex(arr, int(n.CommitIndex))
	if err != nil {
		return
		//we just ignore because we dont need to update the commitIndex
	}
	n.CommitIndex = uint64(newCommitIndex)
	//TODO: io function to Write the entries from the logcurrent index to commitIndex
}

func setUpFinalValueForCommitIndex(arr []MapNode, currCommitIndex int) (int, error) {
	quorumIndex := Quorum - 1
	if quorumIndex < 0 {
		panic("assertion for quorum, it should be bigger than 1, because the minimum defined is 3 or 5")
	}

	if arr[quorumIndex].val <= currCommitIndex {
		return 0, fmt.Errorf("No need to update anything, commitIndex up to date")
	}

	return arr[quorumIndex].val, nil

}
