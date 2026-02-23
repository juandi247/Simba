package simulator

import (
	"math/rand"
	coreraft "simba/coreRaft"
)

const TickFrequency = 1
const maxTicks = 1000000
const maxQueueSize = 100000
const maxInboxSize = 10000

const LeaderHeartbeatFreq = 10
const MinFollowerTimeout = 30
const MaxFollowerTimeout = 50

/*
this are some assertions, so that no one starts puting wrong ranges, that will generate numbers outside of the deifned behaviour
*/
func init() {
	if LeaderHeartbeatFreq >= MinFollowerTimeout || LeaderHeartbeatFreq >= MaxFollowerTimeout {
		panic("the leader heart beat MUST be smaller than both minfollower and maxfollower timeouts")
	}

	if MinFollowerTimeout >= MaxFollowerTimeout {
		panic("the minfollowre must be smaller than the max")
	}

	if minCrashNodeDowntime >= maxCrashNodeDowntime {
		panic("the minCrashNodeDowntime must be smaller than the max")
	}

	if minLatencyDelay >= maxCrashNodeDowntime {
		panic("the minCrashNodeDowntime must be smaller than the max")
	}
}

type SimulationRunner struct {
	NumberOfNodes      uint8
	TimeAdapt          coreraft.TimeAdapter
	NetworkAdapt       coreraft.TransportAdapter
	FuzzyProbabilities *FuzzyConfig
}

type messageQueue struct {
	queue       []coreraft.Message
	size        uint64
	copyCounter uint64
}

type messageInbox struct {
	inbox []coreraft.Message
	size  uint64
}

func generateFollowerTimeout(rng *rand.Rand) uint32 {
	return uint32(MinFollowerTimeout + rng.Intn(MaxFollowerTimeout-MinFollowerTimeout+1))
}

func (s *SimulationRunner) Start() {

	nodeList := initializeNodes(s.NumberOfNodes, *s.FuzzyProbabilities)

	messageQueue := &messageQueue{
		queue:       make([]coreraft.Message, maxQueueSize),
		size:        0,
		copyCounter: 0,
	}

	messageInbox := &messageInbox{
		inbox: make([]coreraft.Message, maxInboxSize),
		size:  0,
	}

	// Engine Loop of execution
	for s.TimeAdapt.Now() <= maxTicks {
		// advance 1 tick
		s.TimeAdapt.Advance(1)

		handleNodeCrash(nodeList, *s.FuzzyProbabilities, s.TimeAdapt.Now())

		updateNodeTimers(nodeList)

		handleComeBackToLiveNode(nodeList, s.TimeAdapt.Now())

		readMessagesToInbox(messageQueue, messageInbox, s.TimeAdapt.Now())

		if messageInbox.size > 0 {
			deliverInboxMessageS(messageInbox, *s.FuzzyProbabilities, nodeList, s.TimeAdapt.Now())
		}

		/*
			TODO: Next step where The nodes should read the timeouts, and define them.
			TODO: in order to execute some things, leader election. ETC */

	}
}

func (s *SimulationRunner) Stop() {
}

func initializeNodes(numberOfNodes uint8, fuzzyProbabilites FuzzyConfig) []*coreraft.Node {
	nodeList := make([]*coreraft.Node, numberOfNodes)

	for i := 1; i <= int(numberOfNodes); i++ {

		timeout := generateFollowerTimeout(fuzzyProbabilites.rand)

		nodeList[i-1] = &coreraft.Node{
			Id:            i,
			FriendNodesId: make([]int, numberOfNodes-1),
			Role:          coreraft.FOLLOWER,
			Term:          0,
			// Leader: , no leader because the comprobation of who is the leader, will be Id== LEader, so at the begining there is no elader
			VotedFor:    make([]string, numberOfNodes),
			Log:         make([]string, coreraft.MaxLogSize),
			CommitIndex: 0,

			Timeout:        timeout,
			Timeoutcounter: timeout,

			LeaderHeartbeat:        LeaderHeartbeatFreq,
			LeaderHeartbeatCounter: LeaderHeartbeatFreq,

			Alive:              true,
			ComeBackToLiveTick: 0,
		}
	}

	return nodeList

}

func handleNodeCrash(nodeList []*coreraft.Node, fuzzyProbabilites FuzzyConfig, currentTick int64) {
	for _, node := range nodeList {
		shouldCrash, comeBackToLiveTick := fuzzyProbabilites.determineCrashingProbabily()
		if !shouldCrash {
			continue
		}
		node.Alive = false
		node.ComeBackToLiveTick = currentTick + comeBackToLiveTick
	}
}

func updateNodeTimers(nodeList []*coreraft.Node) {
	for _, node := range nodeList {
		if node.Id == node.Leader {
			node.LeaderHeartbeatCounter--
		} else {
			node.Timeoutcounter--
		}
	}
}

func handleComeBackToLiveNode(nodeList []*coreraft.Node, currentTick int64) {

	for _, node := range nodeList {
		if node.ComeBackToLiveTick <= currentTick && !node.Alive {
			node.Alive = true
		}

		//This is to reestart the values of timeouts, so that the node starts Cleanly from scratch.
		node.LeaderHeartbeatCounter = node.LeaderHeartbeat
		node.Timeoutcounter = node.Timeout

	}
}

func readMessagesToInbox(messageQueue *messageQueue, messageInbox *messageInbox, currentTick int64) {

	for _, msg := range messageQueue.queue {
		// Assertion for the case where tickFreq is only 1. when hacving a TickFReqcuency >=2 This is not valid
		if msg.DeliveryTick < currentTick && TickFrequency == 1 {
			panic("We found a message that has a lower Tick than the current, with a 1 tickfrecuency, there was something wrong")
		}

		if msg.DeliveryTick <= currentTick {
			// append the message to the inbox. (its really adding because we are preallocating everyything so its not an appnend!!!)
			messageInbox.inbox[messageInbox.size] = msg
			messageInbox.size++

			// decrease the value of the size of the messagequeue, because we are deleting the value (but already allocated with size the queue)
			messageQueue.size--
			continue
		}

		// if reached this point, the value of the tick of the message was bigger than the current, so we move it to the place within the copy counter
		messageQueue.queue[messageQueue.copyCounter] = msg
		messageQueue.copyCounter++
	}

}

func deliverInboxMessageS(messageInbox *messageInbox, fuzzyProbabilites FuzzyConfig, nodeList []*coreraft.Node, currentTick int64) {

	// todo: easier to use a MAP instead of a nested for loop in this case, but meh later
	/*The flow would be to read the messages and verify with the Map if the node is alive. This is to mantain DETERMNISM
	 a MAP traversal is NOT Posible, because breaks the determinism of executiong!!! */
	for _, node := range nodeList {
		if !node.Alive {
			continue
		}

		for i := uint64(0); i <= messageInbox.size; i++ {
			if node.Id == messageInbox.inbox[i].Receiver {
				messagesToBroadcast, numberOfMessages := node.Step(messageInbox.inbox[i])
				if numberOfMessages > 0 {
					for _, msg := range messagesToBroadcast {
						lost, delayTicks := fuzzyProbabilites.RandomizeNetwork()
						if !lost {
							msg.DeliveryTick = currentTick + delayTicks
							//TODO: function to send the message (append it to the queue)
						}
					}
				}
			}
		}
	}

	// "empty" the inbox, because all the messages from it were read.
	messageInbox.size = 0
}
