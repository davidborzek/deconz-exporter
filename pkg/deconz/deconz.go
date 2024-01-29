package deconz

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type Client interface {
	GetSensors(ctx context.Context) (map[string]Sensor, error)
}

type client struct {
	http   *http.Client
	url    string
	apiKey string
}

type Sensor struct {
	Config           map[string]interface{} `json:"config"`
	Ep               int                    `json:"ep"`
	Etag             string                 `json:"etag"`
	Lastseen         string                 `json:"lastseen"`
	Manufacturername string                 `json:"manufacturername"`
	Modelid          string                 `json:"modelid"`
	Name             string                 `json:"name"`
	State            map[string]interface{} `json:"state"`
	Swversion        string                 `json:"swversion"`
	Type             string                 `json:"type"`
	Uniqueid         string                 `json:"uniqueid"`
}

func NewClient(url string, apiKey string) Client {
	return &client{
		http: &http.Client{
			Timeout: 10 * time.Second,
		},
		url:    url,
		apiKey: apiKey,
	}
}

func NewClientWithHTTPClient(httpClient *http.Client, url string, apiKey string) Client {
	return &client{
		http:   httpClient,
		url:    url,
		apiKey: apiKey,
	}
}

func (c *client) GetSensors(ctx context.Context) (map[string]Sensor, error) {
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

	res, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		return nil, fmt.Errorf("received erroneous status code %d", res.StatusCode)
	}

	var sensors map[string]Sensor
	if err := json.NewDecoder(res.Body).Decode(&sensors); err != nil {
		return nil, err
	}

	return sensors, nil
}
