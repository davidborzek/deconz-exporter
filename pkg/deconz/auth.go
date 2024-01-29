package deconz

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

var authHttpClient = &http.Client{
	Timeout: 10 * time.Second,
}

type AuthRequest struct {
	Username   string `json:"username,omitempty"`
	Devicetype string `json:"devicetype"`
}
type AuthSuccess struct {
	Username string `json:"username"`
}

type AuthSuccessResponse struct {
	Success AuthSuccess `json:"success"`
}

type AuthError struct {
	Address     string `json:"address"`
	Description string `json:"description"`
	Type        int    `json:"type"`
}

type AuthErrorResponse struct {
	Error AuthError `json:"error"`
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

	res, err := authHttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		return nil, unmarshalAuthError(res)
	}

	return unmarshalAuthResponse(res)
}

func unmarshalAuthError(res *http.Response) error {
	var data []AuthErrorResponse
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return err
	}

	if len(data) == 0 {
		return errors.New("unknown authentication failure")
	}

	return errors.New(data[0].Error.Description)
}

func unmarshalAuthResponse(res *http.Response) (*AuthSuccessResponse, error) {
	var data []AuthSuccessResponse
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return nil, errors.New("unknown authentication failure")
	}

	return &data[0], nil
}
