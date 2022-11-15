package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

var ErrorCounter = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name:      "scrape_errors_total",
		Namespace: "deconz",
		Help:      "Total errors during scraping metrics from deconz",
	},
)

var Sensor = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name:      "sensor_state",
		Namespace: "deconz",
		Help:      "Sensor state value",
	},
	[]string{"sensor", "type", "state", "manufacturername", "modelid", "name"},
)

func Init() {
	prometheus.MustRegister(ErrorCounter)
	prometheus.MustRegister(Sensor)

	log.Info("prometheus metrics successfully registered")
}
