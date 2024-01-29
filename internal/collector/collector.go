package collector

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/davidborzek/deconz-exporter/pkg/deconz"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	iso8601 = "2006-01-02T15:04:05"
)

type DeconzCollector struct {
	client deconz.Client
}

func NewDeconzCollector(client deconz.Client) *DeconzCollector {
	return &DeconzCollector{
		client: client,
	}
}

func (c *DeconzCollector) Describe(_ chan<- *prometheus.Desc) {}

func (c *DeconzCollector) Collect(ch chan<- prometheus.Metric) {
	ctx := context.Background()

	sensors, err := c.client.GetSensors(ctx)
	if err != nil {
		log.WithError(err).Error("failed to get sensors")

		errorCounterMetric.Inc()
		errorCounterMetric.Collect(ch)

		return
	}

	for id, sensor := range sensors {
		c.setMetrics(ch, id, sensor)
	}

	sensorMetric.Collect(ch)
}

func (c *DeconzCollector) setMetrics(ch chan<- prometheus.Metric, id string, sensor deconz.Sensor) {
	for key, state := range sensor.State {
		switch v := state.(type) {
		case float64:
			c.setMetric(ch, id, key, sensor, v)
		case bool:
			c.setBoolMetric(ch, id, key, sensor, v)
		case string:
			c.setStringMetric(ch, id, key, sensor, v)
		}
	}
}

func (c *DeconzCollector) setMetric(
	ch chan<- prometheus.Metric,
	id string,
	key string,
	sensor deconz.Sensor,
	value float64,
) {
	sensorMetric.WithLabelValues(
		id,
		sensor.Type,
		key,
		sensor.Manufacturername,
		sensor.Modelid,
		sensor.Name,
	).Set(value)
}

func (c *DeconzCollector) setBoolMetric(
	ch chan<- prometheus.Metric,
	id string,
	key string,
	sensor deconz.Sensor,
	value bool,
) {
	var f float64
	if value {
		f = 1
	}

	c.setMetric(ch, id, key, sensor, f)
}

func (c *DeconzCollector) setStringMetric(
	ch chan<- prometheus.Metric,
	id string,
	key string,
	sensor deconz.Sensor,
	value string,
) {
	if date, err := time.Parse(iso8601, value); err == nil {
		c.setMetric(ch, id, key, sensor, float64(date.Unix()))
		return
	}
}
