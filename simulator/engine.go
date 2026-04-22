package simulator

import (
	"fmt"
	"log"
	"math/rand"
	raft "simba/raft"
	"time"
)

type SimulationRunner struct {
	Time               *SimTime
	Network            *SimNetwork
	FuzzyProbabilities FuzzyConfig
	Port               string
	IsHttps            bool
}

func (s *SimulationRunner) Start() {

	server := NewServer(s.Port, s.IsHttps)

	go func(){

		err:= server.StartServer()

		if err!=nil{
			log.Fatal("the server failed: ", err)
		}
		log.Println("server started correctlz")
	}()

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

	fmt.Println("Configuration finished. Starting loop")
	// Engine Loop of execution
	for s.Time.Now() <= maxTicks {
		// advance 1 tick
		s.Time.Advance(TickFrequency)

		crashNodes(nodeList, s.FuzzyProbabilities, s.Time.Now())

		updateNodeTimers(nodeList)

		handleComeBackToLiveNode(nodeList, s.Time.Now())

		handleTimeout(nodeList, s.Network)

		//this is ONLY to read the queue and put the messages into the inbox. No logic of delivering messages to any node here.
		if s.Network.messageQueue.size > 0 {
			readMessagesToInbox(s.Network)

		}
		if s.Network.messageInbox.size > 0 {
			shuffleInbox(s.FuzzyProbabilities.rand, s.Network)
			deliverInboxMessages(s.Network, nodeList)
		}

		fmt.Printf("Tick %v completed. \n", s.Time.Now())
		time.Sleep(1 * time.Second)
		if s.Time.Now() >= 20 {
			panic("finisheddd")
		}

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
				ElectionTimeoutCounter: ElectionTimeout,
			},
		}
	}

	return nodeList

}

func crashNodes(nodeList []*raft.Node, fuzzyProbabilites FuzzyConfig, currentTick int64) {
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
		switch node.Role {
		case raft.LEADER:
			node.SimulatorFields.LeaderHeartbeatCounter--
		case raft.CANDIDATE:
			node.SimulatorFields.ElectionTimeoutCounter--
		case raft.FOLLOWER:
			node.SimulatorFields.Timeoutcounter--
		default:
			panic("a node does not have a valid role")
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

func shuffleInbox(rand *rand.Rand, sn *SimNetwork) {
	rand.Shuffle(int(sn.messageInbox.size), func(i, j int) {
		sn.messageInbox.inbox[i], sn.messageInbox.inbox[j] = sn.messageInbox.inbox[j], sn.messageInbox.inbox[i]
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

		switch node.Role {
		case raft.LEADER:
			if node.SimulatorFields.LeaderHeartbeatCounter <= 0 {
				msg := node.TriggerHeartbeat()
				sm.SendMessage(msg)
			}
		case raft.FOLLOWER:
			if node.SimulatorFields.Timeoutcounter <= 0 {
				timeoutMessage := node.TriggerTimeout()
				sm.SendMessage(timeoutMessage)
			}
		case raft.CANDIDATE:
			if node.SimulatorFields.ElectionTimeoutCounter <= 0 {
				timeoutMessage := node.TriggerElectionTimeout()
				sm.SendMessage(timeoutMessage)
			}
		}
	}
}
