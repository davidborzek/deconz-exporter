package deconz_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/davidborzek/deconz-exporter/pkg/deconz"
	"github.com/stretchr/testify/assert"
)

const (
	deviceName = "testDevice"
)

var (
	authSuccessResponse = deconz.AuthSuccessResponse{
		Success: deconz.AuthSuccess{
			Username: "someKey",
		},
	}

	authErrorResponse = deconz.AuthErrorResponse{
		Error: deconz.AuthError{
			Address:     "TestAddress",
			Description: "TestDescription",
			Type:        0,
		},
	}
)

func TestAuthSucceeds(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req, err := io.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}

		var authRequest deconz.AuthRequest
		if err := json.Unmarshal(req, &authRequest); err != nil {
			panic(err)
		}

		assert.Equal(t, deviceName, authRequest.Devicetype)
		assert.Equal(t, "", authRequest.Username)

		res, err := json.Marshal([]deconz.AuthSuccessResponse{
			authSuccessResponse,
		})
		if err != nil {
			panic(err)
		}

		w.Write(res)
	}))
	defer srv.Close()

	res, err := deconz.Auth(srv.URL, deviceName, "")

	assert.Nil(t, err)
	assert.Equal(t, authSuccessResponse, *res)
}

func TestAuthReturnsErrorWhenServerReturnsOKButEmptyArray(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		res, err := json.Marshal([]deconz.AuthSuccessResponse{})
		if err != nil {
			panic(err)
		}

		w.Write(res)
	}))
	defer srv.Close()

	res, err := deconz.Auth(srv.URL, deviceName, "")

	assert.Nil(t, res)
	assert.Errorf(t, err, "unknown authentication failure")
}

func TestAuthReturnsErrorForClientError(t *testing.T) {
	res, err := deconz.Auth("", deviceName, "")

	assert.NotNil(t, err)
	assert.Nil(t, res)
}

func TestAuthReturnsErrorForServerSideError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		res, err := json.Marshal([]deconz.AuthErrorResponse{
			authErrorResponse,
		})
		if err != nil {
			panic(err)
		}

		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
	}))
	defer srv.Close()

	res, err := deconz.Auth(srv.URL, deviceName, "")

	assert.NotNil(t, err)
	assert.Nil(t, res)

	assert.EqualError(t, err, authErrorResponse.Error.Description)
}

func TestAuthServerReturnsEmptyErrorResponse(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		res, err := json.Marshal([]deconz.AuthErrorResponse{})
		if err != nil {
			panic(err)
		}

		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
	}))
	defer srv.Close()

	res, err := deconz.Auth(srv.URL, deviceName, "")

	assert.Nil(t, res)
	assert.Errorf(t, err, "unknown authentication failure")
}
