package main

import (
	"fmt"
	"simba/adapters"
	// "simba/reality"
	"simba/simulator"
)

const matrixMode bool = true
const SEED = 12345

// By DEFAULT LOW but this should ve changed for the simulations, and for runtime too(?)
const fuzzyLevel simulator.FuzzyLevel = simulator.LOW

func main() {
	var runner adapters.Runner

	if matrixMode {

		fuzzyConfig := simulator.FuzzyConfiguration(SEED, fuzzyLevel)

/*	server := NewServer(s.Port, s.IsHttps)
	go func(){

		err:= server.StartServer()

		if err!=nil{
			log.Fatal("the server failed: ", err)
		}
		log.Println("server started correctlz")
	}() */

		runner = &simulator.SimulationRunner{
			Time:               &simulator.SimTime{},
			Network:            &simulator.SimNetwork{},
			FuzzyProbabilities: fuzzyConfig,
			Port: "8080",
			IsHttps: false,
		}

	} else {
		// transportAdapter:= &reality.RealNetwork{}
		// timeAdapter := &reality.PhysicalTime{}

		// runner= some runner

	}
	fmt.Println("starting program")
	runner.Start()
}





func coso(i, j int) int{

return i+j
}
