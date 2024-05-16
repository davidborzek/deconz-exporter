package auth

import (
	"bufio"
	"fmt"
	"os"

	"github.com/davidborzek/deconz-exporter/pkg/deconz"
	"github.com/urfave/cli/v2"
)

var (
	flags = []cli.Flag{
		&cli.StringFlag{
			Name:     "url",
			Usage:    "The URL of the deCONZ server",
			Required: true,
			EnvVars:  []string{"DECONZ_URL"},
		},
		&cli.StringFlag{
			Name:    "username",
			Usage:   "(optional) Specify a custom username (api key) (10-40 chars)",
			EnvVars: []string{"DECONZ_USERNAME"},
		},
		&cli.StringFlag{
			Name:    "devicetype",
			Usage:   "Name of the client (0-40 chars)",
			Value:   "deconz-exporter",
			EnvVars: []string{"DECONZ_DEVICETYPE"},
		},
		&cli.StringFlag{
			Name:    "format",
			Usage:   "Output format (json, text)",
			Value:   "text",
			EnvVars: []string{"DECONZ_FORMAT"},
		},
		&cli.BoolFlag{
			Name:    "confirm",
			Usage:   "Skip the confirmation prompt",
			EnvVars: []string{"DECONZ_YES"},
			Aliases: []string{"y"},
		},
	}

	Cmd = &cli.Command{
		Name:   "auth",
		Usage:  "Generates a new deCONZ API key",
		Action: run,
		Flags:  flags,
	}
)

func run(ctx *cli.Context) error {
	format := ctx.String("format")
	if format != "json" && format != "text" {
		fmt.Printf("Invalid format: %s\n", format)
		return nil
	}

	if !ctx.Bool("confirm") {
		fmt.Println("Enable discovery in gateway settings and press enter to continue...")
		input := bufio.NewScanner(os.Stdin)
		input.Scan()
	}

	res, err := deconz.Auth(
		ctx.String("url"),
		ctx.String("devicetype"),
		ctx.String("username"),
	)

	if err != nil {
		fmt.Printf("Failed to acquire an api key from deconz: %s\n", err.Error())
		return nil
	}

	if format == "json" {
		fmt.Printf(`{"api_key": "%s"}`, res.Success.Username)
		return nil
	}

	fmt.Printf("API key: %s\n", res.Success.Username)
	return nil
}
