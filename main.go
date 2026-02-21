package main

import (
	coreraft "simba/coreRaft"
	"simba/reality"
	"simba/simulator"
)

var matrixMode bool = true
var NodesNumber uint8 = 5

func main() {
	var runner coreraft.Runner
	var transportAdapter coreraft.TransportAdapter
	var timeAdapter coreraft.TimeAdapter

	if matrixMode {
		transportAdapter = &simulator.SimulatedNetwork{
			MessageQueue: make([]coreraft.Message, 100),
			CurrentTick:  0,
		}

		timeAdapter = &simulator.SimTime{
			Tick: 0,
		}
		runner = &simulator.SimulationRunner{
			NumberOfNodes: NodesNumber,
			TimeAdapt: timeAdapter,
			NetworkAdapt: transportAdapter,
		}

	} else {
		transportAdapter = &reality.RealNetwork{}
		timeAdapter = &reality.PhysicalTime{}

		// runner= some runner

	}

	runner.Start()
}


