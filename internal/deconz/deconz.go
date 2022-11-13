package deconz

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/davidborzek/deconz-exporter/internal/metrics"
)

type Client interface {
	CollectMetrics() error
}

type clientImpl struct {
	url        string
	apiKey     string
	httpClient *http.Client
}

func New(url string, apiKey string) Client {
	return &clientImpl{
		url:    url,
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *clientImpl) getSensors() (*GetSensorsResponse, error) {
	u, err := url.Parse(
		fmt.Sprintf("%s/api/%s/sensors", c.url, c.apiKey),
	)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode >= 400 {
		return nil, fmt.Errorf("received erroneous status code %d", res.StatusCode)
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var sensors GetSensorsResponse
	if err := json.Unmarshal(body, &sensors); err != nil {
		return nil, err
	}

	return &sensors, nil
}

func (c *clientImpl) setMetrics(id string, sensor Sensor) {
	for key, state := range sensor.State {
		if value, ok := state.(float64); ok {
			metrics.Sensor.
				WithLabelValues(
					id,
					sensor.Type,
					key,
					sensor.Manufacturername,
					sensor.Modelid,
					sensor.Name,
				).
				Set(value)
		}
	}
}

func (c *clientImpl) CollectMetrics() error {
	sensors, err := c.getSensors()
	if err != nil {
		return err
	}

	for id, sensor := range *sensors {
		c.setMetrics(id, sensor)
	}

	return nil
}
