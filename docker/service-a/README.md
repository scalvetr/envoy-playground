# Service A

## Build

```shell
go build
```

## Run
```shell

export SERVICE_HOST="0.0.0.0"
export SERVICE_PORT=8080

export METRICS_HOST="0.0.0.0"
export METRICS_PORT=8081

export DOWNSTREAM_SERVICE="http://service-b/"

go run main.go
```