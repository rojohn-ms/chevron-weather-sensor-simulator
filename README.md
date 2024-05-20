# chevron-weather-sensor-simulator

The simulator is a simple Go web server that exposes 2 types of endpoints based on the parameters passed on startup. <br/>
To start the simulator, you run the following:
```cmd
$ ./chevron-weather-sensor-simulator --device=http://host:8080/deviceType/deviceTypeID1 ... --device=http://host:8080/deviceType/deviceTypeIDN
```

The server will then expose N device endpoints, one for each of the `--device` parameters. <br/>

** Note: this is built using Golang v1.22.3

## Endpoints
### Discovery
The discovery endpoint returns a newline separated list of all the "devices" that have been created.

URL
```REST
GET /discovery
```

Response
```text
http://host:port/deviceType/deviceTypeID1
http://host:port/deviceType/deviceTypeID2
...
http://host:port/deviceType/deviceTypeIDN
```

### Device
The device endpoint returns a single float64 values for the device.

URL
```REST
GET /deviceType/deviceTypeID1
```

Response
```text
31.2739485
```

## Building
### Building locally
If you want to build the code locally after some changes, use mage
```cmd
$ mage ci
```

## Running
To run the server, you can build using the steps above, then you can run using the following:
```cmd
$ ./chevron-weather-sensor-simulator --device=http://host:8080/deviceType/deviceTypeID1 ... --device=http://host:8080/deviceType/deviceTypeIDN
```