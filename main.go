package main

import (
	coreraft "simba/coreRaft"
	"simba/reality"
	"simba/simulator"
)

const matrixMode bool = true
const NodesNumber uint8 = 5
const SEED = 12345

// By DEFAULT LOW but this should ve changed for the simulations, and for runtime too(?)
const fuzzyLevel simulator.FuzzyLevel = simulator.LOW

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

		fuzzyConfig := simulator.FuzzyConfiguration(SEED, fuzzyLevel)

		runner = &simulator.SimulationRunner{
			NumberOfNodes:      NodesNumber,
			TimeAdapt:          timeAdapter,
			NetworkAdapt:       transportAdapter,
			FuzzyProbabilities: &fuzzyConfig,
		}

	} else {
		transportAdapter = &reality.RealNetwork{}
		timeAdapter = &reality.PhysicalTime{}

		// runner= some runner

	}

	runner.Start()
}
