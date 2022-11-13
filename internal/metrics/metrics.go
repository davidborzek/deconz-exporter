package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

var Sensor = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name:      "sensor_state",
		Namespace: "deconz",
		Help:      "Sensor state value",
	},
	[]string{"sensor", "type", "state"},
)

func Init() {
	prometheus.MustRegister(Sensor)

	log.Info("prometheus metrics successfully registered")
}
