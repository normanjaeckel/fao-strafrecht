/*
Package env provides the environment variables that are used in this service.
These are:

  FAO_STRAFRECHT_HOST
  FAO_STRAFRECHT_PORT
  FAO_STRAFRECHT_DSFILENAME
*/
package env

import (
	"fmt"
	"strconv"
)

const (
	DefaultHost        = ""
	DefaultPort        = "8000"
	DefaultDSFilenname = "ds.jsonl"
)

// Environment provides all environment variables that are used in this module.
type Environment struct {
	vars map[string]string
}

func (e Environment) Host() string {
	return e.vars["FAO_STRAFRECHT_HOST"]
}

func (e Environment) Port() string {
	return e.vars["FAO_STRAFRECHT_PORT"]
}

func (e Environment) DSFilename() string {
	return e.vars["FAO_STRAFRECHT_DSFILENAME"]
}

// Parse creates the Environment struct with all environment variables retrieved
// from the given function or with default value.
func Parse(fn func(key string) string) (Environment, error) {
	e := Environment{
		vars: map[string]string{
			"FAO_STRAFRECHT_HOST":       DefaultHost,
			"FAO_STRAFRECHT_PORT":       DefaultPort,
			"FAO_STRAFRECHT_DSFILENAME": DefaultDSFilenname,
		},
	}

	for k := range e.vars {
		value := fn(k)
		if value != "" {
			e.vars[k] = value
		}
	}

	if err := validatePort(e.Port()); err != nil {
		return Environment{}, fmt.Errorf("invalid environment variable FAO_STRAFRECHT_PORT: %w", err)
	}

	// TODO: Validate FAO_STRAFRECHT_DSFILENAME: https://stackoverflow.com/questions/35231846/golang-check-if-string-is-valid-path

	return e, nil
}

func validatePort(p string) error {
	portInt, err := strconv.Atoi(p)
	if err != nil {
		return fmt.Errorf("port should be an integer, got %q", p)
	}
	if portInt <= 0 {
		return fmt.Errorf("port should be positiv, got %q", p)
	}
	return nil
}
