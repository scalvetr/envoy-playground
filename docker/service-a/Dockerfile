FROM golang:1.17-alpine AS build-env

# All these steps will be cached
RUN mkdir /app
WORKDIR /app

COPY go.mod .
COPY go.sum .

# Get dependancies - will also be cached if we won't change mod/sum
RUN go mod download

# COPY the source code as the last step
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /go/bin/app

FROM scratch

COPY --from=build-env /go/bin/app /go/bin/app

# Run
ENV SERVICE_HOST="0.0.0.0"
ENV SERVICE_PORT=8080

ENV METRICS_HOST="0.0.0.0"
ENV METRICS_PORT=8081

ENV DOWNSTREAM_SERVICE="http://service-b/"

ENTRYPOINT ["/go/bin/app"]