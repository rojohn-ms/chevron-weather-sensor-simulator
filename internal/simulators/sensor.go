package simulators

import (
	"log"
	"math"
	"math/rand"
	"time"
)

type (
	// IoTSensorSim is an IoT simulated sensor. Based on the code from
	// https://github.com/amine-amaach/simulators/blob/d4edf5e5b57dbf22af0d28a29e5124317ebdebf6/ioTSensorsMQTT-SpB/internal/simulators/ioTSensorSim.go
	IoTSensorSim struct {
		// Sensor Id
		SensorID string
		// sensor data mean value
		mean float64
		// sensor data standard deviation value
		standardDeviation float64
		// sensor data current value
		currentValue float64

		// Delay between each data point
		DelayMin uint32
		DelayMax uint32
		// Randomize delay between data points if true,
		// otherwise DelayMin will be set as fixed delay
		Randomize bool

		// Channel to send data to device
		SensorData chan SensorData
		// Shutdown the sensor
		Shutdown chan bool

		// Check if it's running
		IsRunning bool

		// Sensor Alias, to be used in DDATA, instead of name
		Alias uint64

		// Check if it's already assigned to a device,
		// it's only allowed to be be assigned to one device
		IsAssigned *bool
	}

	// SensorData is the sensor data that will be created.
	SensorData struct {
		Value     float64   `json:"value"`
		Timestamp time.Time `json:"timestamp"`
		Seq       uint64    `json:"seq"`
	}
)

// NewIoTSensorSim creates a new IoTSensorSim
func NewIoTSensorSim(
	id string,
	mean, standardDeviation float64,
	delayMin uint32,
	delayMax uint32,
	randomize bool,
) *IoTSensorSim {
	rand.Seed(time.Now().UnixNano()) // nolint:gosec // not used for crypto
	isAssigned := false
	alias := 100 + uint64(rand.Int63n(10_000)) // nolint:gosec // not used for crypto
	return &IoTSensorSim{
		SensorID:          id,
		mean:              mean,
		standardDeviation: math.Abs(standardDeviation),
		currentValue:      mean - rand.Float64(), // nolint:gosec // not used for crypto
		IsRunning:         false,
		IsAssigned:        &isAssigned,
		SensorData:        make(chan SensorData),
		Shutdown:          make(chan bool, 1),
		DelayMin:          delayMin,
		DelayMax:          delayMax,
		Randomize:         randomize,
		Alias:             alias,
	}
}

// CalculateNextValue gets the next value
func (s *IoTSensorSim) CalculateNextValue() SensorData {
	// first calculate how much the value will be changed
	valueChange := rand.Float64() * math.Abs(s.standardDeviation) / 10 // nolint:gosec // not used for crypto
	// second decide if the value is increased or decreased
	factor := s.decideFactor()
	// apply valueChange and factor to value and return
	s.currentValue += valueChange * factor
	timestamp := time.Now().Local()
	return SensorData{
		Value:     s.currentValue,
		Timestamp: timestamp,
	}
}

func (s *IoTSensorSim) decideFactor() float64 {
	var (
		continueDirection, changeDirection float64
		distance                           float64 // the distance from the mean.
	)
	// depending on if the current value is smaller or bigger than the mean
	// the direction changes.
	if s.currentValue > s.mean {
		distance = s.currentValue - s.mean
		continueDirection = 1
		changeDirection = -1
	} else {
		distance = s.mean - s.currentValue
		continueDirection = -1
		changeDirection = 1
	}
	// the chance is calculated by taking half of the standardDeviation
	// and subtracting the distance divided by 50. This is done because
	// chance with a distance of zero would mean a 50/50 chance for the
	// randomValue to be higher or lower.
	// The division by 50 was found by empiric testing different values.
	chance := (s.standardDeviation / 2) - (distance / 50)
	randomValue := s.standardDeviation * rand.Float64() // nolint:gosec // not used for crypto
	// if the random value is smaller than the chance we continue in the
	// current direction if not we change the direction.
	if randomValue < chance {
		return continueDirection
	}
	return changeDirection
}

// Run runs the sensor simulator.
func (s *IoTSensorSim) Run() {
	if s.IsRunning {
		log.Printf("Senor Id '%s': Already running 🔔\n", s.SensorID)
		return
	}

	s.IsRunning = true
	if s.DelayMin <= 0 {
		s.DelayMin = 1
	} else if s.DelayMin >= s.DelayMax && s.Randomize {
		s.DelayMax = s.DelayMin
	}

	go func() {
		delay := s.DelayMin
		log.Printf("Senor Id '%s': Started running 🔔\n", s.SensorID)
		s.SensorData <- s.CalculateNextValue()
		for {
			select {
			case _, open := <-s.Shutdown:
				log.Printf("Senor Id '%s': Got shutdown signal 🔔\n", s.SensorID)
				s.IsRunning = false
				if open {
					// Send signal to publisher to shutdown
					s.Shutdown <- true
				}
				return

			case <-time.After(time.Duration(delay) * time.Second):
				if s.Randomize {
					delay = uint32(
						rand.Intn(int(s.DelayMax-s.DelayMin)), // nolint:gosec // not used for crypto
					) + s.DelayMin
				}

				s.SensorData <- s.CalculateNextValue()
			}
		}
	}()
}
