package simulator

import (
	"container/heap"
	"fmt"
	raft "simba/raft"
	"time"
)

type SimulationRunner struct {
	Time               *SimTime
	Network            *SimNetwork
	FuzzyProbabilities FuzzyConfig
	Port               string
	IsHttps            bool
	LeaderId 	int
}

func (s *SimulationRunner) Start() {

	// Config for the simulated Time struct
	s.Time.Tick = 0

	// Config for the simulated Network struct
	s.Network.TimeAdapter = s.Time

	pq:= make(PriorityQueue, 0)
	heap.Init(&pq)
	s.Network.messageQueue = &pq

	s.Network.FuzzyConfig = s.FuzzyProbabilities

	//This is all intiial configuration preivous to the FOR loop that ocntains the running engine steps
	nodeList := initializeNodes(s.FuzzyProbabilities)


	//requests:= GenerateRequests(s.FuzzyProbabilities.rand)

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
		if s.Network.messageQueue.Len() > 0 {
			readMessagesToInbox(s.Network, nodeList)
		}
	/*	if s.Network.messageInbox.size > 0 {
			shuffleInbox(s.FuzzyProbabilities.rand, s.Network)
			deliverInboxMessages(s.Network, nodeList)
		}
*/
		fmt.Printf("Tick %v completed. \n", s.Time.Now())
		time.Sleep(10* time.Millisecond)


		if s.Time.Now()>=40{
			time.Sleep(1000*time.Millisecond)
		}
		
	}
}

func (s *SimulationRunner) Stop() {
}

func initializeNextIndex(nodesNumber, id int) map[int]int{
	mapita:= make(map[int]int, nodesNumber-1)

	for i:=1; i<=nodesNumber; i++{
		if i==id{
			continue
		}
		mapita[i]=1
	}

	return mapita
}

func initializeMatchIndex(nodesNumber, id int) map[int]int{
	mapita:= make(map[int]int, nodesNumber-1)

	for i:=1; i<=nodesNumber; i++{
		if i==id{
			continue
		}
		mapita[i]=0
	}

	return mapita
}





func initializeNodes(fuzzyProbabilites FuzzyConfig) []*raft.Node {
	nodeList := make([]*raft.Node, raft.TotalNodesNumber)

	for i := 1; i <= int(raft.TotalNodesNumber); i++ {

		 timeout := generateFollowerTimeout(fuzzyProbabilites.rand)

		if i==1{
			timeout=3
		}
		nodeList[i-1] = &raft.Node{
			Id:            i,
			FriendNodesId: buildFriendsIds(int(raft.TotalNodesNumber), i),
			Role:          raft.FOLLOWER,
			CurrentTerm:   0,
			//leader NOT USED because all will start as candidates. so this will be null for now (or cero)
			Leader:   0,
			VotedFor: make(map[int]int),
			Log: raft.Log{
				Size:   0,
				LogArr: make([]*raft.LogBase, raft.MaxLogSize),
			},
			CommitIndex:     0,
			NextIndex: initializeNextIndex(int(raft.TotalNodesNumber), i),
			MatchIndex: initializeMatchIndex(int(raft.TotalNodesNumber), i),
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
		//shouldCrash, comeBackToLiveTick := fuzzyProbabilites.determineCrashingProbabily()
		shouldCrash:=false
		if !shouldCrash {
			continue
		}
		node.SimulatorFields.Alive = false
		fmt.Println("crashed node: ", node.Id)
		//TODO: uncomment thisss
		//node.SimulatorFields.ComeBackToLiveTick = currentTick + comeBackToLiveTick
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

func readMessagesToInbox(sn *SimNetwork, nodeList []*raft.Node) {

	if sn.messageQueue.Len()<=0{
		panic("wtf this hsuold be bigger than cero")
	}

	for i:=0; i<sn.messageQueue.Len(); i++{
		msg:= sn.messageQueue.Peek()
		if msg.DeliveryTick > int(sn.TimeAdapter.Now()) {
			return
		}
		msg = sn.messageQueue.Pop().(*SimMessage)		//this should changeb
		receiverNodeID := 0

		if _, isEntry := msg.Message.(raft.NewEntry); isEntry{
			leaderId:= checkLeader(nodeList)
			if leaderId==0{
			fmt.Println("there is no current LEADER to respond this message")
				continue
			}
			receiverNodeID = leaderId

		}else{
			receiverNodeID= msg.Message.GetReceiver() 
		}
		fmt.Printf("Message para deliverear: %T\n", msg.Message)
		
		if receiverNodeID <0{
			panic("menor a cero wwtf")
		}
		node:= nodeList[receiverNodeID-1]
		
		responseMessages:=node.ProcessMessage(msg.Message)
		sn.SendMessage(responseMessages)
	}

}

func checkLeader(nodeList []*raft.Node) int{
	leader:=0
	maxTerm:=0
	for _ , n :=range nodeList{
		if n.Role == raft.LEADER && n.CurrentTerm > uint64(maxTerm){
			leader= n.Id
		} 
	}
	return leader
}


/*
ACA ya se habran reducido los tiks por nodo. por lo tanto lo unico seria validar el teimpo no?
*/
func handleTimeout(nodeList []*raft.Node, sm *SimNetwork) {

	for _, node := range nodeList {
		if !node.SimulatorFields.Alive {
			continue
		}

		switch node.Role {
		case raft.LEADER:
			if node.SimulatorFields.LeaderHeartbeatCounter <= 0 {
				fmt.Println("we reached a timeout leader")
				msg := node.TriggerHeartbeat()
				sm.SendMessage(msg)
			}
		case raft.FOLLOWER:
			if node.SimulatorFields.Timeoutcounter <= 0 {
				fmt.Printf("we reached a timeout follower id: %v, this should trigger a election \n", node.Id)
				timeoutMessages := node.TriggerTimeout()
				sm.SendMessage(timeoutMessages)
			}
		case raft.CANDIDATE:
			if node.SimulatorFields.ElectionTimeoutCounter <= 0 {
				fmt.Println("we reached a timeout candidate")
				timeoutMessage := node.TriggerElectionTimeout()
				sm.SendMessage(timeoutMessage)
			}
		}
	}
}
