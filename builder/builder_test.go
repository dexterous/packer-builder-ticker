package builder

import (
	"github.com/mitchellh/packer/packer"
	"testing"
	"time"
)

func TestBuilder_ImplementsBuilder(t *testing.T) {
	t.Parallel()
	var raw interface{} = &Builder{}

	if _, ok := raw.(packer.Builder); !ok {
		t.Error("must implement Builder")
	}
}

func TestBuilder_Prepare_DefaultsConfig(t *testing.T) {
	t.Parallel()
	var builder Builder

	builder.Prepare()

	if builder.config.period != 1 {
		t.Errorf("Period defaulted to %d", builder.config.period)
	}

	if builder.config.duration != 5 {
		t.Errorf("Period defaulted to %d", builder.config.duration)
	}
}

func TestBuilder_Prepare_SetsConfig(t *testing.T) {
	t.Parallel()
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

func TestBuilder_Run_SaysRunning(t *testing.T) {
	t.Parallel()
	var builder Builder
	var ui = newTestUi(t)

	builder.Prepare(&map[string]interface{}{"duration": float64(2)})
	builder.Run(&ui, nil, nil)

	ui.shouldHaveSaid("Running for 2 second(s), ticking every 1 second(s)...")
	ui.shouldHaveSaid("Building... 1")
	ui.shouldHaveSaid("Done! Stopping...")
	ui.shouldHaveSaid("Stopped!")
}

func TestBuilder_Cancel_SaysCancelling(t *testing.T) {
	t.Parallel()
	var builder Builder
	var ui = newTestUi(t)
	var semaphore = make(chan int, 1)

	builder.Prepare(&map[string]interface{}{"duration": float64(5)})
	semaphore <- 1
	go func() {
		builder.Run(&ui, nil, nil)
		<-semaphore
	}()

	time.AfterFunc(1*time.Second+1*time.Millisecond, builder.Cancel)
	semaphore <- 1

	ui.shouldHaveSaid("Running for 5 second(s), ticking every 1 second(s)...")
	ui.shouldHaveSaid("Building... 1")
	ui.shouldHaveSaid("Cancelling...")
	ui.shouldHaveSaid("Cancelled! Stopping...")
	ui.shouldHaveSaid("Stopped!")
}

func TestBuilder_Cancel_DoesNotSayCancellingIfDone(t *testing.T) {
	t.Parallel()
	var builder Builder
	var ui = newTestUi(t)

	builder.Prepare(&map[string]interface{}{"duration": float64(2)})
	builder.Run(&ui, nil, nil)

	builder.Cancel()

	ui.shouldHaveSaid("Running for 2 second(s), ticking every 1 second(s)...")
	ui.shouldHaveSaid("Building... 1")
	ui.shouldHaveSaid("Done! Stopping...")
	ui.shouldHaveSaid("Stopped!")
	ui.shouldNotHaveSaid("Cancelling...")
	ui.shouldNotHaveSaid("Cancelled! Stopping...")
}
