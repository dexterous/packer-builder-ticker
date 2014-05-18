package main

import (
	"github.com/mitchellh/packer/packer"
	"testing"
)

func TestBuilder_ImplementsBuilder(t *testing.T) {
	var raw interface{} = &Builder{}

	if _, ok := raw.(packer.Builder); !ok {
		t.Error("must implement Builder")
	}
}

func TestBuilderPrepare_DefaultsConfig(t *testing.T) {
	var builder Builder

	builder.Prepare()

	if builder.config.period != 1 {
		t.Errorf("Period defaulted to %d", builder.config.period)
	}

	if builder.config.duration != 5 {
		t.Errorf("Period defaulted to %d", builder.config.duration)
	}
}

func TestBuilderPrepare_SetsConfig(t *testing.T) {
	var builder Builder

	builder.Prepare(&map[string]interface{}{
		"period":   float64(10),
		"duration": float64(20),
	})

	if builder.config.period != 10 {
		t.Errorf("Period defaulted to %d", builder.config.period)
	}

	if builder.config.duration != 20 {
		t.Errorf("Period defaulted to %d", builder.config.duration)
	}
}
