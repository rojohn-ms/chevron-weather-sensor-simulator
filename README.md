# chevron-weather-sensor-simulator

The simulator is a simple Go simulator that sends data from a weather sensor simulator to AIO MQ. <br/>
To start the simulator, you run the following:
```cmd
$ ./weather-sim --mqServerURL "mqtts://aio-mq-dmqtt-frontend:8883" --mqTopic "foo-topic"
```

** Note: this is built using Golang v1.22.3

## Building
### Building locally
If you want to build the code locally after some changes, run the following:
```cmd
mage ci
```

## Running
To run the server, you can build using the steps above, then you can run using the following:
```cmd
$ ./bin/weather-sim --mqServerURL "mqtts://aio-mq-dmqtt-frontend:8883" --mqTopic "foo-topic"
```