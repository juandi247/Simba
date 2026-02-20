package simulator

import (
	coreraft "simba/coreRaft"
)
type SimulatedNetwork struct {
	MessageQueue []coreraft.Message
	CurrentTick  int
}

func (s *SimulatedNetwork) SendMessage(msg coreraft.Message) {
	// Logic to append it to the queue

}