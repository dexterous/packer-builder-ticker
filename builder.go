package main

import (
	"fmt"
	"github.com/mitchellh/packer/packer"
	"time"
)

type Builder struct {
	config config
	cancel chan int
	done   bool
	ui     packer.Ui
}

type config struct {
	period   float64
	duration float64
}

func (b *Builder) Prepare(raws ...interface{}) ([]string, error) {
	b.config.period = 1
	b.config.duration = 5

	for _, e := range raws {
		var m map[string]interface{} = *e.(*map[string]interface{})

		if v, ok := m["period"]; ok {
			b.config.period = v.(float64)
		}

		if v, ok := m["duration"]; ok {
			b.config.duration = v.(float64)
		}
	}

	b.cancel = make(chan int, 1)

	return nil, nil
}

func (b *Builder) Run(ui packer.Ui, _ packer.Hook, _ packer.Cache) (packer.Artifact, error) {
	ui.Say(fmt.Sprintf("Running(%3.0f, %3.0f)...", b.config.period, b.config.duration))

	b.ui = ui

	tick := time.Tick(time.Duration(b.config.period) * time.Second)
	stop := time.After(time.Duration(b.config.duration) * time.Second)
	start := time.Now()

	for !b.done {
		select {
		case <-tick:
			ui.Say(fmt.Sprintf("Building... %s", time.Since(start)))
		case <-stop:
			ui.Say("Done! Stopping...")
			b.done = true
		case <-b.cancel:
			ui.Say("Cancelled! Stopping...")
			b.done = true
		}
	}

	return nil, nil
}

func (b *Builder) Cancel() {
	if b.done {
		return
	}

	b.ui.Say("Cancelling...")
	b.cancel <- 1
}
