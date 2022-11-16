package handler_test

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/davidborzek/deconz-exporter/internal/handler"
	"github.com/davidborzek/deconz-exporter/internal/metrics"
	"github.com/davidborzek/deconz-exporter/mock"
	"github.com/golang/mock/gomock"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
)

const (
	authToken = "someToken"
)

func TestMetricsHandlerReturnsOK(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	req, err := http.NewRequest("GET", "/metrics", nil)
	if err != nil {
		t.Fatal(err)
	}

	deconzMock := mock.NewMockClient(ctrl)
	deconzMock.EXPECT().
		CollectMetrics().
		Times(1)

	rr := httptest.NewRecorder()
	h := handler.New(deconzMock, "")

	h.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestMetricsHandlerReturnsUnauthorizedForEmptyAuthorizationHeader(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	req, err := http.NewRequest("GET", "/metrics", nil)
	if err != nil {
		t.Fatal(err)
	}

	deconzMock := mock.NewMockClient(ctrl)
	deconzMock.EXPECT().
		CollectMetrics().
		Times(0)

	rr := httptest.NewRecorder()
	h := handler.New(deconzMock, authToken)

	h.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestMetricsHandlerReturnsUnauthorizedForInvalidToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	req, err := http.NewRequest("GET", "/metrics", nil)
	req.Header.Add("Authorization", "Bearer invalidToken")
	if err != nil {
		t.Fatal(err)
	}

	deconzMock := mock.NewMockClient(ctrl)
	deconzMock.EXPECT().
		CollectMetrics().
		Times(0)

	rr := httptest.NewRecorder()
	h := handler.New(deconzMock, authToken)

	h.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestMetricsHandlerReturnsOKForValidToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	req, err := http.NewRequest("GET", "/metrics", nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", authToken))
	if err != nil {
		t.Fatal(err)
	}

	deconzMock := mock.NewMockClient(ctrl)
	deconzMock.EXPECT().
		CollectMetrics().
		Times(1)

	rr := httptest.NewRecorder()
	h := handler.New(deconzMock, authToken)

	h.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestMetricsHandlerReturnsOKAndIncremetsErrorCounter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	req, err := http.NewRequest("GET", "/metrics", nil)
	if err != nil {
		t.Fatal(err)
	}

	deconzMock := mock.NewMockClient(ctrl)
	deconzMock.EXPECT().
		CollectMetrics().
		Times(1).
		Return(errors.New("some error"))

	rr := httptest.NewRecorder()
	h := handler.New(deconzMock, "")

	h.ServeHTTP(rr, req)

	assert.Equal(t, float64(1), testutil.ToFloat64(metrics.ErrorCounter))
	assert.Equal(t, http.StatusOK, rr.Code)
}
