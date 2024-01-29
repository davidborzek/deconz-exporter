package collector

import "github.com/prometheus/client_golang/prometheus"

var errorCounterMetric = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name:      "scrape_errors_total",
		Namespace: "deconz",
		Help:      "Total errors during scraping metrics from deconz",
	},
)

var sensorMetric = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name:      "sensor_state",
		Namespace: "deconz",
		Help:      "Sensor state value",
	},
	[]string{"sensor", "type", "state", "manufacturername", "modelid", "name"},
)
