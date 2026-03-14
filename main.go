package main

import (
	coreraft "simba/coreRaft"
	// "simba/reality"
	"simba/simulator"
)

const matrixMode bool = true
const NodesNumber uint8 = 5
const SEED = 12345

// By DEFAULT LOW but this should ve changed for the simulations, and for runtime too(?)
const fuzzyLevel simulator.FuzzyLevel = simulator.LOW

func main() {
	var runner coreraft.Runner
	

	if matrixMode {
		
		fuzzyConfig := simulator.FuzzyConfiguration(SEED, fuzzyLevel)

		runner = &simulator.SimulationRunner{
			NumberOfNodes:      NodesNumber,
			Time:          &simulator.SimTime{},
			Network:      &simulator.SimNetwork{},
			FuzzyProbabilities: &fuzzyConfig,
		}

	} else {
		// transportAdapter:= &reality.RealNetwork{}
		// timeAdapter := &reality.PhysicalTime{}

		// runner= some runner

	}

	runner.Start()
}
