package simulator

import (
	coreraft "simba/coreRaft"
)
type SimNetwork struct {
	messageQueue *messageQueue
	messageInbox *messageInbox
	FuzzyConfig FuzzyConfig
	TimeAdapter coreraft.TimeAdapter
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

func (s *SimNetwork) SendMessage(msg coreraft.Message) {
	lost, delayTicks := s.FuzzyConfig.RandomizeNetwork()
		// if the message is Lost, we simply dont add it to the messagequeue, simlating the LOST ont he network
		// todo: there should be a tracker or something for the later UI that indicates that a message was LOST
		if !lost {
			msg.DeliveryTick = s.TimeAdapter.Now() + delayTicks

			/* is the same as this, but instead of sving the currentick as a vairble on simNetwork,
			 its easier to just call the current time so we dont need to updatemanually */
			// msg.DeliveryTick = currentTick + delayTicks

			if s.messageQueue.size >= maxQueueSize {
				panic("MESSAGE QUEUE is FULL")
			}
			s.messageQueue.size++
			s.messageQueue.queue[s.messageQueue.size] = msg
		}
	
}