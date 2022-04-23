package env_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/normanjaeckel/fao-strafrecht/server/pkg/env"
)

func TestEnv(t *testing.T) {
	hostTestValue := "some_host"
	portTestValue := "8080"
	if err := os.Setenv("FAO_STRAFRECHT_HOST", hostTestValue); err != nil {
		t.Fatalf("setting environment: %v", err)
	}
	if err := os.Setenv("FAO_STRAFRECHT_PORT", fmt.Sprint(portTestValue)); err != nil {
		t.Fatalf("setting environment: %v", err)
	}

	e, err := env.Parse(os.Getenv)
	if err != nil {
		t.Fatalf("getting env vars: %v", err)
	}

	t.Run("test FAO_STRAFRECHT_HOST", func(t *testing.T) {
		expected := hostTestValue
		got := e.Host()
		if got != expected {
			t.Fatalf("retrieving env var: expected %q, got %q", expected, got)

		}
	})

	t.Run("test FAO_STRAFRECHT_PORT", func(t *testing.T) {
		expected := portTestValue
		got := e.Port()
		if got != expected {
			t.Fatalf("retrieving env var: expected %q, got %q", expected, got)

		}
	})
}

func TestEnvWithError(t *testing.T) {
	t.Run("bad port value, used bad string", func(t *testing.T) {
		if err := os.Setenv("FAO_STRAFRECHT_PORT", "bad_value"); err != nil {
			t.Fatalf("setting environment: %v", err)
		}

		_, err := env.Parse(os.Getenv)
		if err == nil {
			t.Fatalf("expecting error, but got nil")
		}
	})
	t.Run("bad port value, used negativ int", func(t *testing.T) {
		if err := os.Setenv("FAO_STRAFRECHT_PORT", "-8000"); err != nil {
			t.Fatalf("setting environment: %v", err)
		}

		_, err := env.Parse(os.Getenv)
		if err == nil {
			t.Fatalf("expecting error, but got nil")
		}
	})
}
