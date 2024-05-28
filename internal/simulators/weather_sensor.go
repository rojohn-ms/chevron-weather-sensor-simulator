package simulators

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/url"
	"time"

	"github.com/eclipse/paho.golang/autopaho"
	"github.com/eclipse/paho.golang/paho"
	"github.com/google/uuid"
)

const (
	// ClientIDPrefix is the prefix that will be added to the generated or provided client ID.
	ClientIDPrefix = "aio-chevron"
	// NetworkClosedError is the error message that is returned when the network connection is closed.
	NetworkClosedError = "use of closed network connection"
)

const (
	// QoS0 at most once delivery.
	QoS0 QoS = 0x00
	// QoS1 at least once delivery.
	QoS1 QoS = 0x01
	// QoS2 exactly once delivery.
	QoS2 QoS = 0x10
)

type (
	// QoS indicates the level of assurance for delivery of an Application Message.
	QoS byte

	// WeatherSensorSim is an simulated sensor for weather data.
	WeatherSensorSim struct {
		// Sensor Id
		SensorID string

		// Delay between each data point
		DelayMin uint32
		DelayMax uint32
		// Randomize delay between data points if true,
		// otherwise DelayMin will be set as fixed delay
		Randomize bool

		// Channel to send data to device
		SensorData chan WeatherSensorData
		// Done the sensor
		Done chan bool

		// Check if it's running
		IsRunning bool

		// Sensor Alias, to be used in DDATA, instead of name
		Alias uint64

		// Check if it's already assigned to a device,
		// it's only allowed to be be assigned to one device
		IsAssigned *bool

		// MQTT information
		mqttServerURL string
		mqttTopic     string
		mqttClient    *paho.Client

		// Context
		ctx        context.Context
		cancelFunc context.CancelFunc
	}

	// WeatherSensorData is the sensor data that will be created.
	WeatherSensorData struct {
		Value     float64   `json:"value"`
		Timestamp time.Time `json:"timestamp"`
		Seq       uint64    `json:"seq"`
	}
)

// NewWeatherSensorSim creates a new WeatherSensorSim
func NewWeatherSensorSim(
	id string,
	delayMin uint32,
	delayMax uint32,
	randomize bool,
	mqttServerURL string,
	mqttTopic string,
) *WeatherSensorSim {
	rand.New(rand.NewSource(time.Now().UnixNano())) // nolint:gosec // not used for crypto
	isAssigned := false
	alias := 100 + uint64(rand.Int63n(10_000)) // nolint:gosec // not used for crypto
	return &WeatherSensorSim{
		SensorID:   id,
		IsRunning:  false,
		IsAssigned: &isAssigned,
		SensorData: make(chan WeatherSensorData),
		Done:       make(chan bool, 1),
		DelayMin:   delayMin,
		DelayMax:   delayMax,
		Randomize:  randomize,
		Alias:      alias,

		mqttServerURL: mqttServerURL,
		mqttTopic:     mqttTopic,
	}
}

// Run runs the weather simulator.
func (s *WeatherSensorSim) Run() error {
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

	u, err := url.Parse(s.mqttServerURL)
	if err != nil {
		return err
	}

	cliCfg := autopaho.ClientConfig{
		ServerUrls:                    []*url.URL{u},
		KeepAlive:                     20,
		CleanStartOnInitialConnection: false,
		SessionExpiryInterval:         60,
		OnConnectionUp: func(cm *autopaho.ConnectionManager, connAck *paho.Connack) {
			log.Printf("MQTT connection established on %s\n", s.mqttServerURL)

			// Subscribing in the OnConnectionUp callback is recommended (ensures the subscription is reestablished if
			// the connection drops)
			_, err := cm.Subscribe(context.Background(), &paho.Subscribe{
				Subscriptions: []paho.SubscribeOptions{
					{Topic: s.mqttTopic, QoS: 1},
				},
			})
			if err != nil {
				log.Printf("Failed to subscribe (%s). This is likely to mean no messages will be received.", err)
			}

			log.Println("MQTT subscription made")
		},
		OnConnectError: func(err error) {
			log.Printf("error whilst attempting connection: %s\n", err)
		},
		ClientConfig: paho.ClientConfig{
			// If you are using QOS 1/2, then it's important to specify a client id (which must be unique)
			ClientID: fmt.Sprintf(
				"%s-%v",
				ClientIDPrefix,
				uuid.New(),
			),
			// OnPublishReceived is a slice of functions that will be called when a message is received.
			// You can write the function(s) yourself or use the supplied Router
			OnPublishReceived: []func(paho.PublishReceived) (bool, error){
				func(pr paho.PublishReceived) (bool, error) {
					log.Printf(
						"Received message on topic %s; body: %s (retain: %t)\n",
						pr.Packet.Topic,
						pr.Packet.Payload,
						pr.Packet.Retain,
					)
					return true, nil
				},
			},
			OnClientError: func(err error) { fmt.Printf("client error: %s\n", err) },
			OnServerDisconnect: func(d *paho.Disconnect) {
				if d.Properties != nil {
					log.Printf("Server requested disconnect: %s\n", d.Properties.ReasonString)
				} else {
					log.Printf("Server requested disconnect; reason code: %d\n", d.ReasonCode)
				}
			},
		},
	}

	ctx, cancelFunc := context.WithCancel(context.Background())
	s.ctx, s.cancelFunc = ctx, cancelFunc

	c, err := autopaho.NewConnection(ctx, cliCfg) // starts process; will reconnect until context cancelled
	if err != nil {
		return err
	}

	// Wait for the connection to come up
	if err = c.AwaitConnection(ctx); err != nil {
		return err
	}

	go func() {
		delay := s.DelayMin
		log.Printf("Senor Id '%s': Started running 🔔\n", s.SensorID)
		// s.SensorData <- s.CalculateNextValue()
		for {
			select {
			case <-s.mqttClient.Done():
				log.Printf("Senor Id '%s': Got shutdown signal 🔔\n", s.SensorID)
				s.IsRunning = false
				s.Done <- true
				return

			case <-time.After(time.Duration(delay) * time.Second):
				if s.Randomize {
					delay = uint32(
						rand.Intn(int(s.DelayMax-s.DelayMin)), // nolint:gosec // not used for crypto
					) + s.DelayMin
				}

				// s.SensorData <- s.CalculateNextValue()
			}
		}
	}()

	return nil
}

// Stop stops the weather simulator.
func (s *WeatherSensorSim) Stop() error {
	if s.IsRunning {
		s.cancelFunc()
		<-s.Done
		log.Printf("Weather Senor Id '%s': Stopped 🔔\n", s.SensorID)
	}

	return nil
}
