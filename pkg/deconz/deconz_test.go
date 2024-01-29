package deconz_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/davidborzek/deconz-exporter/pkg/deconz"
	"github.com/stretchr/testify/assert"
)

const (
	apiKey = "1234567"
)

var (
	ctx = context.Background()

	sensors = map[string]deconz.Sensor{
		"1": {
			Config: map[string]interface{}{
				"testConfig": "testValue",
			},
			Ep:               1,
			Etag:             "testEtag",
			Lastseen:         "testLastSeen",
			Manufacturername: "testManufacturerName",
			Modelid:          "testModelId",
			Name:             "testName",
			State: map[string]interface{}{
				"testState": "testValue",
			},
			Swversion: "testSwVersion",
			Type:      "testType",
			Uniqueid:  "testUniqueId",
		},
		"2": {
			Config: map[string]interface{}{
				"testConfig2": "testValue2",
			},
			Ep:               2,
			Etag:             "testEtag2",
			Lastseen:         "testLastSeen2",
			Manufacturername: "testManufacturerName2",
			Modelid:          "testModelId2",
			Name:             "testName2",
			State: map[string]interface{}{
				"testState2": "testValue2",
			},
			Swversion: "testSwVersion2",
			Type:      "testType2",
			Uniqueid:  "testUniqueId2",
		},
	}
)

func TestGetSensors(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(sensors)

		assert.Equal(t, fmt.Sprintf("/api/%s/sensors", apiKey), r.URL.Path)
	}))
	defer srv.Close()

	client := deconz.NewClient(srv.URL, apiKey)

	res, err := client.GetSensors(ctx)

	assert.Nil(t, err)
	assert.Equal(t, sensors, res)
}

func TestGetSensorsFailsWithErroneousStatusCode(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	client := deconz.NewClient(srv.URL, apiKey)

	res, err := client.GetSensors(ctx)

	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "500")
	assert.Nil(t, res)
}

func TestGetSensorsFailsWithErroneousResponseBody(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("test"))
	}))
	defer srv.Close()

	client := deconz.NewClient(srv.URL, apiKey)

	res, err := client.GetSensors(ctx)

	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "invalid character")
	assert.Nil(t, res)
}
