package reality

import (
	coreraft "simba/coreRaft"
)

type RealNetwork struct {
	// NodeIPs map[int]string
}

func (n *RealNetwork) SendMessage(msg coreraft.Message) {
	// targetIP := n.nodeIPs[msg.DestinationID]
	// sendOverNetwork(msg, targetIP)
}

