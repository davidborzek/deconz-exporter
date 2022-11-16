package server

import (
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/davidborzek/deconz-exporter/internal/deconz"
	"github.com/davidborzek/deconz-exporter/internal/handler"
	"github.com/davidborzek/deconz-exporter/internal/metrics"
	"github.com/urfave/cli/v2"

	log "github.com/sirupsen/logrus"
)

var (
	flags = []cli.Flag{
		&cli.StringFlag{
			Name:     "deconz-url",
			Usage:    "The URL of the deCONZ server",
			Required: true,
			EnvVars:  []string{"DECONZ_URL"},
		},
		&cli.StringFlag{
			Name:     "deconz-key",
			Required: true,
			Usage:    "The api key (also called username) of deCONZ server. See the docs on how to generate it",
			EnvVars:  []string{"DECONZ_API_KEY"},
		},
		&cli.StringFlag{
			Name:    "port",
			Value:   "8080",
			Usage:   "The port of deCONZ exporter server",
			EnvVars: []string{"DECONZ_EXPORTER_PORT"},
		},
		&cli.StringFlag{
			Name:    "host",
			Usage:   "The host of deCONZ exporter server",
			EnvVars: []string{"DECONZ_EXPORTER_HOST"},
		},
		&cli.StringFlag{
			Name:    "auth-token",
			Usage:   "Optional auth token for deCONZ exporter server. If no token is set authentication is disabled.",
			EnvVars: []string{"DECONZ_EXPORTER_AUTH_TOKEN"},
		},
		&cli.StringFlag{
			Name:    "log-level",
			Usage:   "Log level",
			Value:   "info",
			EnvVars: []string{"DECONZ_EXPORTER_LOG_LEVEL"},
		},
	}

	Cmd = &cli.Command{
		Name:   "server",
		Usage:  "Starts the deCONZ exporter server",
		Action: run,
		Flags:  flags,
	}
)

func parseLogLevel(level string) log.Level {
	switch strings.ToLower(level) {
	case "debug":
		return log.DebugLevel
	case "info":
		return log.InfoLevel
	case "warning":
		return log.WarnLevel
	case "error":
		return log.ErrorLevel
	case "fatal":
		return log.FatalLevel
	}

	log.WithField("level", level).
		Warn("invalid log level provided - falling back to 'info'")

	return log.InfoLevel
}

func run(ctx *cli.Context) error {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	log.SetLevel(
		parseLogLevel(ctx.String("log-level")),
	)

	log.WithField("pid", os.Getpid()).
		Infof("deCONZ prometheus exporter started")

	metrics.Init()

	d := deconz.New(
		ctx.String("deconz-url"),
		ctx.String("deconz-key"),
	)

	token := ctx.String("auth-token")
	if len(token) > 0 {
		log.Info("authentication is enabled")
	}

	h := handler.New(d, token)

	addr := net.JoinHostPort(ctx.String("host"),
		ctx.String("port"))
	log.WithField("addr", addr).
		Infof("starting the http server")

	return http.ListenAndServe(addr, h)
}
