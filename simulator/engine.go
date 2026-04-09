package simulator

import (
	"math/rand"
	raft "simba/raft"
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
	Time               *SimTime
	Network            *SimNetwork
	FuzzyProbabilities FuzzyConfig
}

func generateFollowerTimeout(rng *rand.Rand) uint32 {
	return uint32(MinFollowerTimeout + rng.Intn(MaxFollowerTimeout-MinFollowerTimeout+1))
}

func (s *SimulationRunner) Start() {

	// Config for the simulated Time struct
	s.Time.Tick = 0

	// Config for the simulated Network struct
	s.Network.TimeAdapter = s.Time
	s.Network.messageQueue = &messageQueue{
		queue:       make([]SimMessage, maxQueueSize),
		size:        0,
		copyCounter: 0,
	}

	s.Network.messageInbox = &messageInbox{
		inbox: make([]raft.Message, maxInboxSize),
		size:  0,
	}
	s.Network.FuzzyConfig = s.FuzzyProbabilities
	//This is all intiial configuration preivous to the FOR loop that ocntains the running engine steps
	nodeList := initializeNodes(s.FuzzyProbabilities)

	// Engine Loop of execution
	for s.Time.Now() <= maxTicks {
		// advance 1 tick
		s.Time.Advance(TickFrequency)

		handleNodeCrash(nodeList, s.FuzzyProbabilities, s.Time.Now())

		updateNodeTimers(nodeList)

		handleComeBackToLiveNode(nodeList, s.Time.Now())

		//this is ONLY to read the queue and put the messages into the inbox. No logic of delivering messages to any node here.
		if s.Network.messageQueue.size > 0 {
			readMessagesToInbox(s.Network)

		}
		if s.Network.messageInbox.size > 0 {
			deliverInboxMessages(s.Network, nodeList)
		}

		handleTimeout(nodeList, s.Network)


	}
}

func (s *SimulationRunner) Stop() {
}

func initializeNodes(fuzzyProbabilites FuzzyConfig) []*raft.Node {
	nodeList := make([]*raft.Node, raft.TotalNodesNumber)

	for i := 1; i <= int(raft.TotalNodesNumber); i++ {

		timeout := generateFollowerTimeout(fuzzyProbabilites.rand)

		nodeList[i-1] = &raft.Node{
			Id:            i,
			FriendNodesId: [raft.TotalNodesNumber - 1]int{},
			Role:          raft.FOLLOWER,
			CurrentTerm:   0,
			//leader NOT USED because all will start as candidates. so this will be null for now (or cero)
			Leader:   0,
			VotedFor: 0,
			Log: raft.Log{
				Size:   0,
				LogArr: make([]*raft.LogBase, raft.MaxLogSize),
			},
			CommitIndex:     0,
			Timeout:         timeout,
			LeaderHeartbeat: LeaderHeartbeatFreq,

			SimulatorFields: &raft.SimulatorFields{
				LeaderHeartbeatCounter: LeaderHeartbeatFreq,
				Alive:                  true,
				ComeBackToLiveTick:     0,
				Timeoutcounter:         timeout,
			},
		}
	}

	return nodeList

}

func handleNodeCrash(nodeList []*raft.Node, fuzzyProbabilites FuzzyConfig, currentTick int64) {
	for _, node := range nodeList {
		shouldCrash, comeBackToLiveTick := fuzzyProbabilites.determineCrashingProbabily()
		if !shouldCrash {
			continue
		}
		node.SimulatorFields.Alive = false
		node.SimulatorFields.ComeBackToLiveTick = currentTick + comeBackToLiveTick
	}
}

func updateNodeTimers(nodeList []*raft.Node) {
	for _, node := range nodeList {
		if node.Id == node.Leader {
			node.SimulatorFields.LeaderHeartbeatCounter--
		} else {
			node.SimulatorFields.Timeoutcounter--
		}
	}
}

func handleComeBackToLiveNode(nodeList []*raft.Node, currentTick int64) {

	for _, node := range nodeList {
		if node.SimulatorFields.ComeBackToLiveTick <= currentTick && !node.SimulatorFields.Alive {
			node.SimulatorFields.Alive = true
			//This is to reestart the values of timeouts, so that the node starts Cleanly from scratch.
			node.SimulatorFields.LeaderHeartbeatCounter = node.LeaderHeartbeat
			node.SimulatorFields.Timeoutcounter = node.Timeout
		}

	}
}

func readMessagesToInbox(sn *SimNetwork) {

	for _, msg := range sn.messageQueue.queue {
		// Assertion for the case where tickFreq is only 1. when hacving a TickFReqcuency >=2 This is not valid
		if msg.DeliveryTick < sn.TimeAdapter.Now() && TickFrequency == 1 {
			panic("We found a message that has a lower Tick than the current, with a 1 tickfrecuency, there was something wrong")
		}

		if msg.DeliveryTick <= sn.TimeAdapter.Now() {
			// append the message to the inbox. (its really adding because we are preallocating everyything so its not an appnend!!!)
			sn.messageInbox.inbox[sn.messageInbox.size] = msg.Message
			sn.messageInbox.size++

			// decrease the value of the size of the messagequeue, because we are deleting the value (but already allocated with size the queue)
			sn.messageQueue.size--
			continue
		}

		// if reached this point, the value of the tick of the message was bigger than the current, so we move it to the place within the copy counter
		sn.messageQueue.queue[sn.messageQueue.copyCounter] = msg
		sn.messageQueue.copyCounter++
	}

}

func shuffleInbox(rand *rand.Rand, sn *SimNetwork ){
	rand.Shuffle(int(sn.messageInbox.size), func(i, j int) {
sn.messageInbox.inbox[i], sn.messageInbox.inbox[j]  = sn.messageInbox.inbox[j], sn.messageInbox.inbox[i] 
	})

}


func deliverInboxMessages(sn *SimNetwork, nodeList []*raft.Node) {


	// todo: easier to use a MAP instead of a nested for loop in this case, but meh later
	for _, node := range nodeList {
		if !node.SimulatorFields.Alive {
			continue
		}

		for i := uint64(0); i < sn.messageInbox.size; i++ {
			if node.Id == sn.messageInbox.inbox[i].GetReceiver() {
				messagesToSend := node.ProcessMessage(sn.messageInbox.inbox[i])
				sn.SendMessage(messagesToSend)
			}
		}
	}

	// "empty" the inbox, because all the messages from it were read.
	sn.messageInbox.size = 0
}

/*
ACA ya se habran reducido los tiks por nodo. por lo tanto lo unico seria validar el teimpo no?
*/
func handleTimeout(nodeList []*raft.Node, sm *SimNetwork) {

	for _, node := range nodeList {
		if !node.SimulatorFields.Alive {
			return
		}

		if node.Role == raft.LEADER && int(node.SimulatorFields.LeaderHeartbeatCounter) <= 0 {
			msg := node.TriggerHeartbeat()
			sm.SendMessage(msg)
			continue
		}

		// for both FOLLOWER and CANDIDATE
		if node.SimulatorFields.Timeoutcounter <= 0 {
			timeoutMessage := node.TriggerTimeout()
			//this message will triger a timeout election
			sm.SendMessage(timeoutMessage)
		}
	}
}
