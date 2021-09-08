package env

import (
	"strings"

	"github.com/urfave/cli/v2"
)

const dev = "dev"

var (
	e = dev

	FlagEnv = cli.StringFlag{
		Name:        "e",
		Aliases:     []string{"env"},
		Usage:       "run mode, load the config files in 'conf/<env>/*'",
		Required:    true,
		Value:       dev,
		DefaultText: dev,
		Destination: &e,
	}
)

func WithEnv(ee string) {
	ee = strings.TrimSpace(ee)
	switch ee {
	case "", dev:
		e = dev
	default:
		e = ee
	}
}

func Name() string { return e }

func IsDev() bool          { return e == dev }
func Is(ee string) bool    { return e == ee }
func Equal(ee string) bool { return e == ee }
