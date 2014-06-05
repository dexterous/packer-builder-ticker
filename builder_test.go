package main

import (
	"fmt"
	"github.com/mitchellh/packer/packer"
	"testing"
	"time"
)

func TestBuilder_ImplementsBuilder(t *testing.T) {
	var raw interface{} = &Builder{}

	if _, ok := raw.(packer.Builder); !ok {
		t.Error("must implement Builder")
	}
}

func TestBuilder_Prepare_DefaultsConfig(t *testing.T) {
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
	var builder Builder
	var ui = newTestUi(t)

	ui.ShouldSay("Running for 2 second(s), ticking every 1 second(s)...")
	ui.ShouldSay("Building... 1")
	ui.ShouldSay("Done! Stopping...")
	ui.ShouldSay("Stopped!")

	builder.Prepare(&map[string]interface{}{"duration": float64(2)})
	builder.Run(&ui, nil, nil)

	ui.Verify()
}

func TestBuilder_Cancel_SaysCancelling(t *testing.T) {
	var builder Builder
	var ui = newTestUi(t)
	var semaphore = make(chan int, 1)

	ui.ShouldSay("Running for 5 second(s), ticking every 1 second(s)...")
	ui.ShouldSay("Building... 1")
	ui.ShouldSay("Cancelling...")
	ui.ShouldSay("Cancelled! Stopping...")
	ui.ShouldSay("Stopped!")

	builder.Prepare(&map[string]interface{}{"duration": float64(5)})
	go func() {
		builder.Run(&ui, nil, nil)
		semaphore <- 1
	}()

	time.AfterFunc(1*time.Second+1*time.Millisecond, builder.Cancel)
	<-semaphore

	ui.Verify()
}

func TestBuilder_Cancel_DoesNotSayCancellingIfDone(t *testing.T) {
	var builder Builder
	var ui = newTestUi(t)

	ui.ShouldSay("Running for 2 second(s), ticking every 1 second(s)...")
	ui.ShouldSay("Building... 1")
	ui.ShouldSay("Done! Stopping...")
	ui.ShouldSay("Stopped!")
	ui.ShouldNotSay("Cancelling...")
	ui.ShouldNotSay("Cancelled! Stopping...")

	builder.Prepare(&map[string]interface{}{"duration": float64(2)})
	builder.Run(&ui, nil, nil)

	builder.Cancel()

	ui.Verify()
}


type result struct {
	found    bool
	expected bool
}

func (r *result) invalid() bool {
	return (r.expected && !r.found) || (!r.expected && r.found)
}

func (r *result) foundString() string {
	if r.found {
		return "was said"
	} else {
		return "was not said"
	}
}

func (r *result) expectedString() string {
	if r.expected {
		return "should be said"
	} else {
		return "should not be said"
	}
}

func (r *result) String() string {
	return fmt.Sprintf("%s, %s", r.expectedString(), r.foundString())
}

type TestUi struct {
	conversation map[string]*result
	*testing.T
	packer.Ui
}

func newTestUi(t *testing.T) TestUi {
	return TestUi{T: t, conversation: map[string]*result{}}
}

func (u *TestUi) Say(actual string) {
	u.Logf("Said %s", actual)
	if result, present := u.conversation[actual]; present && !result.found {
		result.found = true
	}
}

func (u *TestUi) ShouldSay(said string) {
	u.conversation[said] = &result{expected: true}
}

func (u *TestUi) ShouldNotSay(said string) {
	u.conversation[said] = &result{expected: false}
}

func (u *TestUi) Verify() {
	for expected, result := range u.conversation {
		u.Logf("Verifying '%s' %s", expected, result)
		if result.invalid() {
			u.Errorf("'%s' %s", expected, result)
		}
	}
}
