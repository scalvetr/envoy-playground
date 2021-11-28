package main

import (
	"context"
	"fmt"
	"github.com/hellofresh/health-go/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

var ERROR_RATE int

var (
	requestCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "service_b_requests_total",
			Help: "The total number of requests received by Service B.",
		},
		[]string{"status"},
	)
)

func buildHealth() *health.Health {
	h, _ := health.New(health.WithChecks(health.Config{
		Name:      "status",
		Timeout:   time.Second * 5,
		SkipOnErr: true,
		Check: func(ctx context.Context) error {
			// rabbitmq health check implementation goes here
			return nil
		}},
	))
	return h
}

func handler(w http.ResponseWriter, r *http.Request) {
	if rand.Intn(100) >= ERROR_RATE {
		fmt.Fprintln(w, "Service B: Yay! nounce", rand.Uint32())
		requestCounter.WithLabelValues("2xx").Inc()
		log.Println("HTTP 200", r.Method, r.URL, r.RemoteAddr)
	} else {
		http.Error(w, fmt.Sprintf("Service B: Ooops... nounce %v", rand.Uint32()),
			http.StatusInternalServerError)
		requestCounter.WithLabelValues("5xx").Inc()
		log.Println("HTTP 500", r.Method, r.URL, r.RemoteAddr)
	}
}

func main() {
	log.Println("Starting service-a")
	log.Println(os.Environ())
	n, err := strconv.Atoi(os.Getenv("ERROR_RATE"))
	if err != nil {
		log.Fatal("Can not parse ERROR_RATE env")
	}
	ERROR_RATE = n

	prometheus.MustRegister(requestCounter)
	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(":"+os.Getenv("METRICS_PORT"), nil)
	log.Println("Metrics server started")

	health := buildHealth()
	http.Handle("/health/readiness", health.Handler())
	http.Handle("/health/liveness", health.Handler())
	http.HandleFunc("/", handler)
	log.Println("Handlers registered!!")
	log.Fatal(http.ListenAndServe(":"+os.Getenv("SERVICE_PORT"), nil))
}
