package deconzauth

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

var httpClient = &http.Client{
	Timeout: 10 * time.Second,
}

func unmarshal(body []byte) (*AuthSuccessResponse, error) {
	var r []AuthSuccessResponse
	if err := json.Unmarshal(body, &r); err != nil {
		return nil, err
	}

	if len(r) == 0 {
		return nil, errors.New("unknown authentication failure")
	}

	return &r[0], nil
}

func unmarshalError(body []byte) error {
	var r []AuthErrorResponse
	if err := json.Unmarshal(body, &r); err != nil {
		return err
	}

	if len(r) == 0 {
		return errors.New("unknown authentication failure")
	}

	return errors.New(r[0].Error.Description)
}

func Auth(deconzUrl string, deviceType string, username string) (*AuthSuccessResponse, error) {
	u, err := url.Parse(
		fmt.Sprintf("%s/api", deconzUrl),
	)
	if err != nil {
		return nil, err
	}

	authRequest := AuthRequest{
		Username:   username,
		Devicetype: deviceType,
	}

	reqBody, err := json.Marshal(&authRequest)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, u.String(), bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	res, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode >= 400 {
		return nil, unmarshalError(body)
	}

	return unmarshal(body)
}
