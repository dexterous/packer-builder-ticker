package builder

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
	period   uint
	duration uint
}

func (b *Builder) Prepare(raws ...interface{}) ([]string, error) {
	b.config.period = 1
	b.config.duration = 5

	for _, e := range raws {
		var m map[string]interface{} = *e.(*map[string]interface{})

		if v, ok := m["period"]; ok {
			b.config.period = uint(v.(float64))
		}

		if v, ok := m["duration"]; ok {
			b.config.duration = uint(v.(float64))
		}
	}

	b.cancel = make(chan int, 1)

	return nil, nil
}

func (b *Builder) Run(ui packer.Ui, _ packer.Hook, _ packer.Cache) (packer.Artifact, error) {
	ui.Say(fmt.Sprintf("Running for %d second(s), ticking every %d second(s)...", b.config.duration, b.config.period))

	b.ui = ui

	tick := time.Tick(time.Duration(b.config.period) * time.Second)
	stop := time.After(time.Duration(b.config.duration) * time.Second)
	start := time.Now()

	for !b.done {
		select {
		case <-tick:
			ui.Say(fmt.Sprintf("Building... %d", uint(time.Since(start).Seconds())))
		case <-stop:
			ui.Say("Done! Stopping...")
			b.done = true
		case <-b.cancel:
			ui.Say("Cancelled! Stopping...")
			b.done = true
		}
	}

	ui.Say("Stopped!")
	return nil, nil
}

func (b *Builder) Cancel() {
	if b.done {
		return
	}

	b.ui.Say("Cancelling...")
	b.cancel <- 1
}
