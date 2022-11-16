package deconz_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/davidborzek/deconz-exporter/internal/deconz"
	"github.com/davidborzek/deconz-exporter/internal/metrics"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
)

var (
	testSensor = deconz.Sensor{
		Manufacturername: "TestManufacturer",
		Modelid:          "TestModel",
		Name:             "Test",
		Type:             "TestType",
		State: map[string]interface{}{
			"numeric": 123,
		},
	}

	secondTestSensor = deconz.Sensor{
		Manufacturername: "TestManufacturer2",
		Modelid:          "TestModel2",
		Name:             "Test2",
		Type:             "TestType2",
		State: map[string]interface{}{
			"bool": true,
		},
	}
)

func resetPrometheus() {
	metrics.Sensor.Reset()
}

func TestCollectMetricsSucceedsForNoSensors(t *testing.T) {
	defer resetPrometheus()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "{}")
	}))
	defer srv.Close()

	d := deconz.New(srv.URL, "")

	err := d.CollectMetrics()

	assert.Equal(t, 0, testutil.CollectAndCount(metrics.Sensor))
	assert.Nil(t, err)
}

func TestCollectMetricsSucceedsAndSetsMetrics(t *testing.T) {
	defer resetPrometheus()

	sensorId := "1"
	secondSensorId := "2"

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		res, err := json.Marshal(deconz.GetSensorsResponse{
			sensorId:       testSensor,
			secondSensorId: secondTestSensor,
		})
		if err != nil {
			panic(err)
		}
		w.Write(res)
	}))
	defer srv.Close()

	d := deconz.New(srv.URL, "someKey")
	err := d.CollectMetrics()

	assert.Nil(t, err)
	assert.Equal(t, 2, testutil.CollectAndCount(metrics.Sensor))

	numericMetric := testutil.ToFloat64(metrics.Sensor.WithLabelValues(
		sensorId,
		testSensor.Type,
		"numeric",
		testSensor.Manufacturername,
		testSensor.Modelid,
		testSensor.Name,
	))
	assert.Equal(t, float64(123), numericMetric)

	boolMetric := testutil.ToFloat64(metrics.Sensor.WithLabelValues(
		secondSensorId,
		secondTestSensor.Type,
		"bool",
		secondTestSensor.Manufacturername,
		secondTestSensor.Modelid,
		secondTestSensor.Name,
	))
	assert.Equal(t, float64(1), boolMetric)
}

func TestCollectMetricsReturnsErrorForErroneousStatusCode(t *testing.T) {
	defer resetPrometheus()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	d := deconz.New(srv.URL, "")

	err := d.CollectMetrics()

	assert.Equal(t, 0, testutil.CollectAndCount(metrics.Sensor))
	assert.NotNil(t, err)
}

func TestCollectMetricsReturnsErrorForClientError(t *testing.T) {
	defer resetPrometheus()
	d := deconz.New("UNKNOWN-URL", "")

	err := d.CollectMetrics()

	assert.Equal(t, 0, testutil.CollectAndCount(metrics.Sensor))
	assert.NotNil(t, err)
}

func TestCollectMetricsReturnsErrorForUnmarshallingError(t *testing.T) {
	defer resetPrometheus()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "BROKEN")
	}))
	defer srv.Close()

	d := deconz.New(srv.URL, "")

	err := d.CollectMetrics()

	assert.Equal(t, 0, testutil.CollectAndCount(metrics.Sensor))
	assert.NotNil(t, err)
}
