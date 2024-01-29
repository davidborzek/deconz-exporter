package collector_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/davidborzek/deconz-exporter/internal/collector"
	"github.com/davidborzek/deconz-exporter/pkg/deconz"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

var (
	sensors = map[string]deconz.Sensor{
		"1": {
			Manufacturername: "TestManufacturer",
			Modelid:          "TestModel",
			Name:             "Test",
			Type:             "TestType",
			State: map[string]interface{}{
				"numeric":  123,
				"datetime": "2023-09-18T17:43:12",
			},
		},
		"2": {
			Manufacturername: "TestManufacturer2",
			Modelid:          "TestModel2",
			Name:             "Test2",
			Type:             "TestType2",
			State: map[string]interface{}{
				"bool": true,
			},
		},
	}
)

func TestCollectMetrics(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(sensors)
	}))
	defer srv.Close()

	client := deconz.NewClient(srv.URL, "1234567")
	dc := collector.NewDeconzCollector(client)

	const expected = `
	# HELP deconz_sensor_state Sensor state value
	# TYPE deconz_sensor_state gauge
	deconz_sensor_state{manufacturername="TestManufacturer",modelid="TestModel",name="Test",sensor="1",state="datetime",type="TestType"} 1.695058992e+09
	deconz_sensor_state{manufacturername="TestManufacturer",modelid="TestModel",name="Test",sensor="1",state="numeric",type="TestType"} 123
	deconz_sensor_state{manufacturername="TestManufacturer2",modelid="TestModel2",name="Test2",sensor="2",state="bool",type="TestType2"} 1
	`

	if err := testutil.CollectAndCompare(dc, strings.NewReader(expected)); err != nil {
		t.Errorf("unexpected collecting result:\n%s", err)
	}
}

func TestCollectMetricsErrorCounter(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	client := deconz.NewClient(srv.URL, "1234567")
	dc := collector.NewDeconzCollector(client)

	const expected = `
	# HELP deconz_scrape_errors_total Total errors during scraping metrics from deconz
	# TYPE deconz_scrape_errors_total counter
	deconz_scrape_errors_total 1
	`

	if err := testutil.CollectAndCompare(dc, strings.NewReader(expected)); err != nil {
		t.Errorf("unexpected collecting result:\n%s", err)
	}
}
