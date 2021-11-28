package main

import (
	"context"
	"fmt"
	"github.com/hellofresh/health-go/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

var requestCounter = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "service_a_requests_total",
		Help: "The total number of requests received by Service A.",
	},
	[]string{"status"},
)

func httpGet(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("HTTP %v - %v", resp.StatusCode, string(body))
	}

	return string(body), nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	resp, err := httpGet(os.Getenv("DOWNSTREAM_SERVICE"))
	if err == nil {
		fmt.Fprintln(w, "Service A: downstream responded with:", resp)
		requestCounter.WithLabelValues("2xx").Inc()
		log.Println("HTTP 200", r.Method, r.URL, r.RemoteAddr)
	} else {
		http.Error(w, fmt.Sprintf("Service A: downstream failed with: %v", err.Error()),
			http.StatusInternalServerError)
		requestCounter.WithLabelValues("5xx").Inc()
		log.Println("HTTP 500", r.Method, r.URL, r.RemoteAddr)
	}
}

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

func main() {
	log.Println("Starting service-a")
	log.Println(os.Environ())

	prometheus.MustRegister(requestCounter)
	http.Handle("/metrics", promhttp.Handler())

	go http.ListenAndServe(os.Getenv("METRICS_HOST")+":"+os.Getenv("METRICS_PORT"), nil)
	log.Println("Metrics server started")

	health := buildHealth()
	http.Handle("/health/readiness", health.Handler())
	http.Handle("/health/liveness", health.Handler())
	http.HandleFunc("/", handler)
	log.Println("Handlers registered!!")
	log.Fatal(http.ListenAndServe(
		os.Getenv("SERVICE_HOST")+":"+os.Getenv("SERVICE_PORT"), nil))
}
