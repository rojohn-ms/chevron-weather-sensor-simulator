// Package config contains configuration information
package config

import (
	"encoding/json"
	"log"
)

type (
	// Config is the configuration
	Config struct {
		EdgeNode EdgeNodeConfig `json:"eon_node"`
	}

	// EdgeNodeConfig is the config for an edge node
	EdgeNodeConfig struct {
		Namespace string         `json:"namespace"`
		GroupID   string         `json:"group_id"`
		NodeID    string         `json:"node_id"`
		Devices   []DeviceConfig `json:"devices"`
	}

	// DeviceConfig is the config for a device
	DeviceConfig struct {
		DeviceID        string            `json:"device_id"`
		StoreAndForward bool              `json:"store_and_forward"`
		TTL             uint32            `json:"time_to_live"`
		Simulators      []IoTSensorConfig `json:"simulators"`
	}

	// IoTSensorConfig is the config for an IoT sensor
	IoTSensorConfig struct {
		SensorID  string  `json:"sensor_id"`
		Mean      float64 `json:"mean"`
		Std       float64 `json:"standard_deviation"`
		DelayMin  uint32  `json:"delay_min"`
		DelayMax  uint32  `json:"delay_max"`
		Randomize bool    `json:"randomize"`
	}
)

// DefaultConfig gets the default config
func DefaultConfig() Config {
	cfg := Config{}

	defaultConfig := []byte(`
	{	
		"eon_node": {
			"namespace": "spBv1.0",
			"group_id": "WeatherSensors",
			"node_id": "SparkplugB",
			"devices": [
				{
					"device_id": "emulatedDevice",
					"store_and_forward": true,
					"time_to_live": 10,
					"simulators": [
						{
							"sensor_id": "Temperature",
							"mean": 30.6,
							"standard_deviation": 3.1,
							"delay_min": 3,
							"delay_max": 6,
							"randomize": true
						}
					]
				},
				{
					"device_id": "anotherEmulatedDevice",
					"store_and_forward": true,
					"time_to_live": 15,
					"simulators": [
						{
							"sensor_id": "Humidity",
							"mean": 40.7,
							"standard_deviation": 2.3,
							"delay_min": 4,
							"delay_max": 10,
							"randomize": false
						}
					]
				}
			]
		}
	}`)

	_ = json.Unmarshal(defaultConfig, &cfg)
	log.Printf("Default config parsed successfully ✅\n")
	return cfg
}

// DefaultTemperatureSensorConfig gets the default temperature sensor config
func DefaultTemperatureSensorConfig() IoTSensorConfig {
	iotSensorCfg := IoTSensorConfig{}

	defaultCfg := []byte(`
	{
		"sensor_id": "Temperature",
		"mean": 30.6,
		"standard_deviation": 3.1,
		"delay_min": 3,
		"delay_max": 6,
		"randomize": true
	}`)

	_ = json.Unmarshal(defaultCfg, &iotSensorCfg)
	log.Printf("Default temperature sensor config parsed successfully ✅\n")
	return iotSensorCfg
}

// DefaultHumiditySensorConfig gets the default humidity sensor config
func DefaultHumiditySensorConfig() IoTSensorConfig {
	iotSensorCfg := IoTSensorConfig{}

	defaultCfg := []byte(`
	{
		"sensor_id": "Humidity",
		"mean": 40.7,
		"standard_deviation": 2.3,
		"delay_min": 4,
		"delay_max": 10,
		"randomize": false
	}`)

	_ = json.Unmarshal(defaultCfg, &iotSensorCfg)
	log.Printf("Default humidity sensor config parsed successfully ✅\n")
	return iotSensorCfg
}
