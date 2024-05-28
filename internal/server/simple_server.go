package server

import (
	"log"

	"chevron-weather-sensor-simulator/internal/config"
	"chevron-weather-sensor-simulator/internal/simulators"
)

// Run runs the web server
func RunSimple(shutdown <-chan bool) {
	deviceCfg := config.DefaultTemperatureSensorConfig()
	simpleSensorSim := simulators.NewSimpleSensorSim(
		"simple-sim",
		deviceCfg.DelayMin,
		deviceCfg.DelayMax,
		deviceCfg.Randomize,
	)

	idle := make(chan struct{})
	go func() {
		<-shutdown

		if serr := simpleSensorSim.Stop(); serr != nil {
			log.Printf("[server.Run][sim.Stop] ERROR: %s\n", serr)
		}

		log.Printf("[server.Run][sim.Stop] Simulator stopped\n")
		close(idle)
	}()

	log.Printf("[server.Run] Starting device: [%s]\n", addr)
	if err := simpleSensorSim.Run(); err != nil {
		log.Printf("[server.Run][sim.Run] ERROR: %s\n", err)
	}

	<-idle
}
