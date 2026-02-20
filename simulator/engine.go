package simulator

import coreraft "simba/coreRaft"

type SimulationRunner struct {
	NumberOfNodes int
	TimeAdapt       coreraft.TimeAdapter
	NetworkAdapt	   coreraft.TransportAdapter
}

func (s *SimulationRunner) Start() {

	/* 1. create the nodes according to the number of nodes
	assign them the values of id, term on 0, index on cero, everythingg on cero AND TIMING
	Timing must be random for ticks  */

	/* 2. Start the for loop containing the engine
	Inside the engine:
		-Update tick
		s.timing.Advance(1)

		- Update tick on Nodes for the heartbeat (or that could go on raft logic, depends)

	*/

}

func (s *SimulationRunner) Stop() {

}