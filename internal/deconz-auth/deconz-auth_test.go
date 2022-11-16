package deconzauth_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	deconzauth "github.com/davidborzek/deconz-exporter/internal/deconz-auth"
	"github.com/stretchr/testify/assert"
)

const (
	deviceName = "testDevice"
)

var (
	authSuccessResponse = deconzauth.AuthSuccessResponse{
		Success: deconzauth.AuthSuccess{
			Username: "someKey",
		},
	}

	authErrorResponse = deconzauth.AuthErrorResponse{
		Error: deconzauth.AuthError{
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

		var authRequest deconzauth.AuthRequest
		if err := json.Unmarshal(req, &authRequest); err != nil {
			panic(err)
		}

		assert.Equal(t, deviceName, authRequest.Devicetype)
		assert.Equal(t, "", authRequest.Username)

		res, err := json.Marshal([]deconzauth.AuthSuccessResponse{
			authSuccessResponse,
		})
		if err != nil {
			panic(err)
		}

		w.Write(res)
	}))
	defer srv.Close()

	res, err := deconzauth.Auth(srv.URL, deviceName, "")

	assert.Nil(t, err)
	assert.Equal(t, authSuccessResponse, *res)
}

func TestAuthReturnsErrorWhenServerReturnsOKButEmptyArray(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		res, err := json.Marshal([]deconzauth.AuthSuccessResponse{})
		if err != nil {
			panic(err)
		}

		w.Write(res)
	}))
	defer srv.Close()

	res, err := deconzauth.Auth(srv.URL, deviceName, "")

	assert.Nil(t, res)
	assert.Errorf(t, err, "unknown authentication failure")
}

func TestAuthReturnsErrorForClientError(t *testing.T) {
	res, err := deconzauth.Auth("", deviceName, "")

	assert.NotNil(t, err)
	assert.Nil(t, res)
}

func TestAuthReturnsErrorForServerSideError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		res, err := json.Marshal([]deconzauth.AuthErrorResponse{
			authErrorResponse,
		})
		if err != nil {
			panic(err)
		}

		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
	}))
	defer srv.Close()

	res, err := deconzauth.Auth(srv.URL, deviceName, "")

	assert.NotNil(t, err)
	assert.Nil(t, res)

	assert.EqualError(t, err, authErrorResponse.Error.Description)
}

func TestAuthServerReturnsEmptyErrorResponse(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		res, err := json.Marshal([]deconzauth.AuthErrorResponse{})
		if err != nil {
			panic(err)
		}

		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
	}))
	defer srv.Close()

	res, err := deconzauth.Auth(srv.URL, deviceName, "")

	assert.Nil(t, res)
	assert.Errorf(t, err, "unknown authentication failure")
}
