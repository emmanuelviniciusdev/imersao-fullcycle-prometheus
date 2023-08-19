package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"math/rand"
	"net/http"
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
		w.WriteHeader(http.StatusNoContent)
	})

	certainEndpointHandler := promhttp.InstrumentHandlerDuration(
		durationHTTPRequestCertainEndpointHistogramVec.MustCurryWith(
			prometheus.Labels{"handler": "certainEndpointHandlerFunc"}),
		promhttp.InstrumentHandlerCounter(totalHTTPRequestsCertainEndpointCounterVec, certainEndpointHandlerFunc),
	)

	http.Handle("/metrics", promhttp.HandlerFor(r, promhttp.HandlerOpts{}))
	http.Handle("/certain-endpoint", certainEndpointHandler)

	log.Fatal(http.ListenAndServe(":8181", nil))
}
