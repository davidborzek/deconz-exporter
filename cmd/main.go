package cmd

import (
	"fmt"
	"strings"

	"github.com/davidborzek/deconz-exporter/cmd/auth"
	"github.com/davidborzek/deconz-exporter/cmd/server"
	"github.com/urfave/cli/v2"
)

func Main(args []string) {
	app := cli.App{
		Name:  "deCONZ Prometheus exporter",
		Usage: "Export metrics from deCONZ to prometheus format",
		Commands: []*cli.Command{
			server.Cmd,
			auth.Cmd,
		},
	}

	// Workaround use `server` as default command.
	if len(args) == 1 || (strings.HasPrefix(args[1], "-") && app.Command(args[1]) == nil) {
		app.Action = server.Cmd.Action
		app.Flags = append(app.Flags, server.Cmd.Flags...)
	}

	if err := app.Run(args); err != nil {
		fmt.Println(err.Error())
	}
}
