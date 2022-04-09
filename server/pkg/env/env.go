/*
Package env provides the environment variables that are used in this service.
*/
package env

import (
	"fmt"
	"strconv"

	"github.com/normanjaeckel/fao-strafrecht/server/pkg/deps"
)

// Environment provides all environment variables that are used in this module.
type Environment struct {
	Host string
	Port int
}

// Parse creates the Environment struct with all environment variables retrieved
// from the given function or with default value.
func Parse(fn deps.GetEnvFunc) (Environment, error) {
	f := envOrDefault(fn)
	e := Environment{}

	e.Host = f("FAO_STRAFRECHT_HOST", "")
	port := f("FAO_STRAFRECHT_PORT", "8000")
	portInt, err := strconv.Atoi(port)
	if err != nil {
		return Environment{}, fmt.Errorf("invalid environment variable FAO_STRAFRECHT_PORT, it should be an integer, got %q", port)
	}
	if portInt <= 0 {
		return Environment{}, fmt.Errorf("invalid environment variable FAO_STRAFRECHT_PORT, it should be positiv, got %d", portInt)
	}
	e.Port = portInt

	return e, nil
}

func envOrDefault(fn deps.GetEnvFunc) func(key string, defaultValue string) string {
	f := func(key string, defaultValue string) string {
		v := fn(key)
		if v == "" {
			return defaultValue
		}
		return v
	}
	return f
}
