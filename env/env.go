package env

import (
	"strings"
)

const dev = "dev"

var e = dev

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
