package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	exporter "noufel/custom-node-exporter/exporter"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	var (
		promPort = flag.Int("prom.port", 9150, "port to expose prometheus metrics")
	)
	flag.Parse()

	// called on each collector.Collect which is called on each prometheus interval default 15 sec
	getMetrics := func() (exporter.Metrics,error){
		return exporter.GetMetrics()
	}

	// make prometheus aware of our collectors
	collector := exporter.NewCollector(getMetrics)

	registry := prometheus.NewRegistry()
	registry.MustRegister(collector)

	go func() {
		for {
			collector.Update()
			time.Sleep(time.Second)
		}
	}()

	mux := http.NewServeMux()
	promHandler := promhttp.HandlerFor(registry ,promhttp.HandlerOpts{})
	mux.Handle("/metrics",promHandler)

	// start listening for HTTP connnections
	port := fmt.Sprintf(":%d",*promPort)
	log.Printf("starting exporter on %q/metrics",port)
	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatalf("cannot start exporter : %s ", err)
	}
}