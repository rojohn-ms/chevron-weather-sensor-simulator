# chevron-weather-sensor-simulator

The simulator is a simple Go web server that exposes 2 types of endpoints based on the parameters passed on startup. <br/>
To start the simulator, you run the following:
```cmd
$ ./chevron-weather-sensor-simulator --device=http://host:8080/deviceType/deviceTypeID1 ... --device=http://host:8080/deviceType/deviceTypeIDN
```

The server will then expose N device endpoints, one for each of the `--device` parameters. <br/>

** Note: this is built using Golang v1.22.3

## Building
### Building locally
If you want to build the code locally after some changes, run the following:
```cmd
$ go build -o bin/weather-sim ./cmd/app
```

## Running
To run the server, you can build using the steps above, then you can run using the following:
```cmd
$ ./bin/weather-sim --device=http://host:8080/deviceType/deviceTypeID1 ... --device=http://host:8080/deviceType/deviceTypeIDN
```