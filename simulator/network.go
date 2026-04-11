package simulator

import (
	"simba/adapters"
	raft "simba/raft"
)

type SimNetwork struct {
	messageQueue *messageQueue
	messageInbox *messageInbox
	FuzzyConfig  FuzzyConfig
	TimeAdapter  adapters.TimeAdapter
}

type SimMessage struct {
	DeliveryTick int64
	Message      raft.Message
}

type messageQueue struct {
	queue       []SimMessage
	size        uint64
	copyCounter uint64
}

type messageInbox struct {
	inbox []raft.Message
	size  uint64
}

func (s *SimNetwork) SendMessage(messages []raft.Message) {
	for _, message := range messages {
		var delayTicks int64
		var lost bool
		if isNetworkMessage(message) {
			lost, delayTicks = s.FuzzyConfig.RandomizeNetwork()
		} else {
			lost, delayTicks = false, 1
		}
		//TODO: there should be a tracker or something for the later UI that indicates that a message was LOST
		if !lost {
			simMessage := SimMessage{
				DeliveryTick: s.TimeAdapter.Now() + delayTicks,
				Message:      message,
			}

			if s.messageQueue.size >= maxQueueSize {
				panic("MESSAGE QUEUE is FULL")
			}
			s.messageQueue.size++
			s.messageQueue.queue[s.messageQueue.size] = simMessage
		}
	}
}


