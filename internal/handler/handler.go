package handler

import (
	"net/http"

	"github.com/davidborzek/deconz-exporter/internal/deconz"
)

type handler struct {
	d             deconz.Client
	expectedToken string
	mux           *http.ServeMux
}

func New(d deconz.Client, authToken string) *handler {
	s := &handler{
		d:             d,
		expectedToken: authToken,
		mux:           http.NewServeMux(),
	}

	s.mux.HandleFunc("/metrics", s.handleMetrics())
	s.mux.HandleFunc("/health", s.handleHealth())

	return s
}

func (s *handler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(rw, r)
}
