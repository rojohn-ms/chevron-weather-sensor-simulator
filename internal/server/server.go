package server

import (
	"context"
	"flag"
	"fmt"
	"html"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"

	"chevron-weather-sensor-simulator/internal/config"
	"chevron-weather-sensor-simulator/internal/simulators"
)

const (
	addr = ":8080"
)

type (
	// RepeatableFlag is an alias to use repeated flags with flag
	RepeatableFlag []string
)

var (
	_       flag.Value = (*RepeatableFlag)(nil)
	devices RepeatableFlag
)

// Run runs the web server
func Run(shutdown <-chan bool) {
	flag.Var(&devices, "device", "Repeat this flag to add devices to the discovery service")
	flag.Parse()

	// At a minimum, respond on '/'
	if len(devices) == 0 {
		log.Printf("[server.Run] ERROR: Must have at least 1 device defined\n")
		return
	}

	log.Printf("[server.Run] Devices: %d\n", len(devices))

	handler := http.NewServeMux()
	handler.HandleFunc("/discovery", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[server.discovery] Handler entered\n")
		fmt.Fprintf(w, "%s\n", html.EscapeString(devices.String()))
	})

	// Create handler for each endpoint
	for _, devicePath := range devices {
		log.Printf("[server.Run] Creating handler: %s", devicePath)

		// Path is of the form: http://host:port/type/deviceNum
		u, err := url.Parse(devicePath)
		if err != nil {
			log.Printf("[server.Run][url.Parse] ERROR: %s\n", err)
			return
		}

		pathParts := strings.Split(u.Path, "/")
		deviceType := pathParts[1]
		var deviceCfg config.IoTSensorConfig
		switch deviceType {
		case "temperature":
			deviceCfg = config.DefaultTemperatureSensorConfig()

		case "humidity":
			deviceCfg = config.DefaultHumiditySensorConfig()

		default:
			log.Printf(
				"[server.Run] Using temperature configuration since device type '%s' is not recognized.\n",
				deviceType,
			)
			deviceCfg = config.DefaultTemperatureSensorConfig()
		}

		sensor := simulators.NewIoTSensorSim(
			u.Path,
			deviceCfg.Mean,
			deviceCfg.Std,
			deviceCfg.DelayMin,
			deviceCfg.DelayMax,
			deviceCfg.Randomize,
		)

		handler.HandleFunc(u.Path, func(w http.ResponseWriter, r *http.Request) {
			log.Printf("[server.device] Handler entered: %s", u.Path)
			log.Printf("[server.device] Invoking sensor: %s", sensor.SensorID)
			fmt.Fprint(w, sensor.CalculateNextValue().Value)
		})
	}

	httpServer := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	idle := make(chan struct{})
	go func() {
		<-shutdown

		if serr := httpServer.Shutdown(context.Background()); serr != nil {
			log.Printf("[server.Run][http.Shutdown] ERROR: %s\n", serr)
		}

		close(idle)
	}()

	listen, err := net.Listen("tcp", addr) // nolint:gosec // Used just for development
	if err != nil {
		log.Printf("[server.Run][net.Listen] ERROR: %s\n", err)
	}

	log.Printf("[server.Run] Starting device: [%s]\n", addr)
	if err := httpServer.Serve(listen); err != http.ErrServerClosed {
		log.Printf("[server.Run][http.Serve] ERROR: %s\n", err)
	}

	<-idle
}

// String is a method required by flag.Value interface
func (e *RepeatableFlag) String() string {
	result := strings.Join(*e, "\n")
	return result
}

// Set is a method required by flag.Value interface
func (e *RepeatableFlag) Set(value string) error {
	*e = append(*e, value)
	return nil
}
