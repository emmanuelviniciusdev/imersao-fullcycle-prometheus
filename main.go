package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"math/rand"
	"net/http"
	"time"
)

var onlineUsersGauge = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "golang_app_online_users",
	Help: "Golang App Online Users",
	ConstLabels: map[string]string{
		"foo": "bar",
	},
})

var totalHTTPRequestsCertainEndpointCounterVec = prometheus.NewCounterVec(prometheus.CounterOpts{
	Name: "golang_app_total_http_requests_certain_endpoint",
	Help: "Total number of HTTP requests made to a certain endpoint...",
}, []string{})

var durationHTTPRequestCertainEndpointHistogramVec = prometheus.NewHistogramVec(prometheus.HistogramOpts{
	Name: "golang_app_duration_http_request_certain_endpoint",
	Help: "Duration of HTTP requests made to a certain endpoint...",
}, []string{"handler"})

func main() {
	r := prometheus.NewRegistry()

	r.MustRegister(onlineUsersGauge)
	r.MustRegister(totalHTTPRequestsCertainEndpointCounterVec)
	r.MustRegister(durationHTTPRequestCertainEndpointHistogramVec)

	go func() {
		for {
			onlineUsersGauge.Set(float64(rand.Intn(1_000_000)))
		}
	}()

	certainEndpointHandlerFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
		w.WriteHeader(http.StatusNoContent)
	})

	anotherCertainEndpointHandlerFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Duration(rand.Intn(5)) * time.Second)
		w.WriteHeader(http.StatusNoContent)
	})

	certainEndpointHandler := promhttp.InstrumentHandlerDuration(
		durationHTTPRequestCertainEndpointHistogramVec.MustCurryWith(
			prometheus.Labels{"handler": "certainEndpointHandlerFunc"}),
		promhttp.InstrumentHandlerCounter(totalHTTPRequestsCertainEndpointCounterVec, certainEndpointHandlerFunc),
	)

	anotherCertainEndpointHandler := promhttp.InstrumentHandlerDuration(
		durationHTTPRequestCertainEndpointHistogramVec.MustCurryWith(
			prometheus.Labels{"handler": "anotherCertainEndpointHandlerFunc"}),
		promhttp.InstrumentHandlerCounter(totalHTTPRequestsCertainEndpointCounterVec, anotherCertainEndpointHandlerFunc),
	)

	http.Handle("/metrics", promhttp.HandlerFor(r, promhttp.HandlerOpts{}))
	http.Handle("/certain-endpoint", certainEndpointHandler)
	http.Handle("/another-certain-endpoint", anotherCertainEndpointHandler)

	log.Fatal(http.ListenAndServe(":8181", nil))
}
