// Package server contains various servers for running simulators.
package server

import (
	"flag"

	"chevron-weather-sensor-simulator/internal/config"
	"chevron-weather-sensor-simulator/internal/logger"
	"chevron-weather-sensor-simulator/internal/simulators"
)

// RunWeather runs the weather server.
func RunWeather(shutdown <-chan bool) {
	mqServerURLPtr := flag.String("mqServerURL", "", "The server URL for AIO MQ")
	mqTopicPtr := flag.String("mqTopic", "", "The AIO MQ topic to send data to")
	logLvlPtr := flag.String("logLevel", "", "The log level to use. Valid options are [ info, debug, none ]")
	flag.Parse()

	log := logger.New(*logLvlPtr)
	deviceCfg := config.DefaultTemperatureSensorConfig()
	weatherSensorSim := simulators.NewWeatherSensorSim(
		"weather-sim",
		deviceCfg.DelayMin,
		deviceCfg.DelayMax,
		deviceCfg.Randomize,
		*mqServerURLPtr,
		*mqTopicPtr,
		*logLvlPtr,
	)

	idle := make(chan struct{})
	go func() {
		<-shutdown

		if serr := weatherSensorSim.Stop(); serr != nil {
			log.Printf("[server.Run][sim.Stop] ERROR: %s\n", serr)
		}

		log.Printf("[server.Run][sim.Stop] Simulator stopped\n")
		close(idle)
	}()

	log.Printf("[server.Run] Starting weather simulator\n")
	if err := weatherSensorSim.Run(); err != nil {
		log.Printf("[server.Run][sim.Run] ERROR: %s\n", err)
	}

	<-idle
}
