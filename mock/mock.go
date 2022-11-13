//go:build !prod
// +build !prod

package mock

//go:generate mockgen -package=mock -destination=mock_gen.go github.com/davidborzek/deconz-exporter/internal/deconz Client
