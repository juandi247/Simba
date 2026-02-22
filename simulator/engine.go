package simulator

import (
	"fmt"
	coreraft "simba/coreRaft"
)


const TickFrequency=1
const maxTicks= 1000000
const maxQueueSize=100000
const maxInboxSize=10000

type SimulationRunner struct {
	NumberOfNodes uint8
	TimeAdapt       coreraft.TimeAdapter
	NetworkAdapt	   coreraft.TransportAdapter
	FuzzyProbabilities *FuzzyConfig
}

type messageQueue struct{
	queue []coreraft.Message
	size uint64
	copyCounter uint64
}


type messageInbox struct{
	inbox []coreraft.Message
	size uint64
}






func (s *SimulationRunner) Start() {

	nodeList:= make([]*coreraft.Node, s.NumberOfNodes)

	/* 1. create the nodes according to the number of nodes
	assign them the values of id, term on 0, index on cero, everythingg on cero AND TIMING
	Timing must be random for ticks  */
	for  i:=1; i<=int(s.NumberOfNodes); i++{
		// randomHeartBeatTimeout:=  FUNCTIONNN!!
		
		nodeList[i-1]= &coreraft.Node{
			Id: i, 
			FriendNodesId: make([]int, s.NumberOfNodes-1),
			Role: coreraft.FOLLOWER ,
			Term: 0,
			// Leader: , no leader because the comprobation of who is the leader, will be Id== LEader, so at the begining there is no elader
			VotedFor: make([]string, s.NumberOfNodes),
			Log: make([]string, coreraft.MaxLogSize),
			CommitIndex: 0,
			HeartbeatTimeout:2 , //this shuold be random
			LeaderHeartbeatTime: 2, // this also reandom but less than the hearbeat timeout 
			Alive: true,
		}
	}


	messageQueue:= &messageQueue{
		queue: make([]coreraft.Message, maxQueueSize),
		size: 0,
		copyCounter: 0,
	}
	
	
	messageInbox:= &messageInbox{
		inbox: make([]coreraft.Message, maxInboxSize),
		size: 0,
	}

	/* 2. Start the for loop containing the engine
	Inside the engine:
		-Update tick
		s.timing.Advance(1)

		- Update tick on Nodes for the heartbeat (or that could go on raft logic, depends)

	*/


	// MAIN LOOP OF TICKS
	for s.TimeAdapt.Now() <= maxTicks{
		// advance 1 tick
		s.TimeAdapt.Advance(1)


		// TODO: Here should be the calculation of probability of each node. Here is the crashing node fuzzy.

		// Advance the time of the nodes before any processing of messages
		for _ , node:= range nodeList{
			if node.Id == node.Leader{
				node.LeaderHeartbeatTime--
			}else{
				node.HeartbeatTimeout--
			}
		} 
		

		/* 
		 TODO:  ALSO HERE should be the logic of reading the tick of the nodes that are CRASHED, and if the current tick is 
		 TODO: equal to the node tick to come back alive, we mark it as ALIVE=TRUE and restart the values of LeaderHEarbeat and HEARBEAT (this can be using the random because the flow is determnisitic)
		*/




			/* 
			Steps: 
			Check with a for loop each one of the values inside the queue.
			If the message has a <=Tick with the currentTick, we should add it to the inbox.
			Update the inbox.
			Decrease the size of the queue

			Else we copy the value that was bigger than current, on the messageQueeue, with the copyCounter INdex
			Update copycounter
			*/
		

		for _, msg:= range messageQueue.queue{
			//! ASSERTION
			if msg.DeliveryTick<s.TimeAdapt.Now() && TickFrequency==1{
				panic("We found a message that has a lower Tick than the current, with a 1 tickfrecuency, there was something wrong")
			}

			
			if msg.DeliveryTick<=s.TimeAdapt.Now(){
				// append the message to the inbox. (its really adding because we are preallocating everyything so its not an appnend!!!)
				messageInbox.inbox[messageInbox.size]= msg
				messageInbox.size++

				// decrease the value of the size of the messagequeue, because we are deleting the value (but already allocated with size the queue)
				messageQueue.size--
				continue
			}

			// if reached this point, the value of the tick of the message was bigger than the current, so we move it to the place within the copy counter
			messageQueue.queue[messageQueue.copyCounter]=msg
			messageQueue.copyCounter++
		}




		// DELIVER THE MESSAGES!
		if messageInbox.size>0{
			// todo: easier to use a MAP instead of a nested for loop in this case, but meh later
			for _, node:= range nodeList{
				//!Validation for crashed nodes, if its crashed, means that he SHOULD not receive or process mesages
				if !node.Alive{
					continue
				}
				//? Running NON FAULTY NODES
				for i:=uint64(0); i<=messageInbox.size ; i++{
					if node.Id == messageInbox.inbox[i].Receiver{
						// here the RAFT would READ it
						// function of that
						messagesToBroadcast, numberOfMessages:= node.Step(messageInbox.inbox[i])
						if numberOfMessages>0{
							// FUZZING FOR EACH MESSAGE AND APPEND THEM
							fmt.Println(messagesToBroadcast)					
						}
					}
				} 
			}

			// RESET the inbox, so that is "empty"
			messageInbox.size=0
		}

		
		/* In this step the nodes should read the timeouts, and define them.
		 in order to execute some things, leader election. ETC */



		
	}

}



func (s *SimulationRunner) Stop() {
}