package simulators

import (
	"context"
	"log"
	"math/rand"
	"time"
)

type (
	// SimpleSensorSim is a simple sensor simulator.
	SimpleSensorSim struct {
		// Sensor Id
		SensorID string

		// Delay between each data point
		DelayMin uint32
		DelayMax uint32
		// Randomize delay between data points if true,
		// otherwise DelayMin will be set as fixed delay
		Randomize bool

		// Done the sensor
		Done chan bool

		// Check if it's running
		IsRunning bool

		// Sensor Alias, to be used in DDATA, instead of name
		Alias uint64

		// Check if it's already assigned to a device,
		// it's only allowed to be be assigned to one device
		IsAssigned *bool

		// Context
		ctx        context.Context
		cancelFunc context.CancelFunc
	}
)

// NewSimpleSensorSim creates a new SimpleSensorSim
func NewSimpleSensorSim(
	id string,
	delayMin uint32,
	delayMax uint32,
	randomize bool,
) *SimpleSensorSim {
	rand.Seed(time.Now().UnixNano()) // nolint:gosec // not used for crypto
	isAssigned := false
	alias := 100 + uint64(rand.Int63n(10_000)) // nolint:gosec // not used for crypto
	return &SimpleSensorSim{
		SensorID:   id,
		IsRunning:  false,
		IsAssigned: &isAssigned,
		Done:       make(chan bool, 1),
		DelayMin:   delayMin,
		DelayMax:   delayMax,
		Randomize:  randomize,
		Alias:      alias,
	}
}

// Run runs the weather simulator.
func (s *SimpleSensorSim) Run() error {
	if s.IsRunning {
		log.Printf("Senor Id '%s': Already running 🔔\n", s.SensorID)
		return nil
	}

	s.IsRunning = true
	if s.DelayMin <= 0 {
		s.DelayMin = 1
	} else if s.DelayMin >= s.DelayMax && s.Randomize {
		s.DelayMax = s.DelayMin
	}

	ctx, cancelFunc := context.WithCancel(context.Background())
	s.ctx, s.cancelFunc = ctx, cancelFunc

	go func() {
		delay := s.DelayMin
		log.Printf("Senor Id '%s': Started running 🔔\n", s.SensorID)
		for {
			select {
			case <-s.ctx.Done():
				log.Printf("Senor Id '%s': Got shutdown signal 🔔\n", s.SensorID)
				s.IsRunning = false
				s.Done <- true
				return

			case <-time.After(time.Duration(delay) * time.Second):
				log.Printf("Senor Id '%s': Next value\n", s.SensorID)
			}
		}
	}()

	return nil
}

// Stop stops the weather simulator.
func (s *SimpleSensorSim) Stop() error {
	if s.IsRunning {
		s.cancelFunc()
		<-s.Done
		log.Printf("Senor Id '%s': Stopped 🔔\n", s.SensorID)
	}

	return nil
}
