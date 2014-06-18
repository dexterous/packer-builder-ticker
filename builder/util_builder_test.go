package builder

import (
	"github.com/mitchellh/packer/packer"
	"testing"
)

type TestUi struct {
	conversation map[string]bool
	*testing.T
	packer.Ui
}

func newTestUi(t *testing.T) TestUi {
	return TestUi{T: t, conversation: map[string]bool{}}
}

func (u *TestUi) Say(actual string) {
	u.Logf("Said %s", actual)
	u.conversation[actual] = true
}

func (u *TestUi) shouldHaveSaid(dialog string) {
	u.verify(dialog, true, "should be said, but was not said")
}

func (u *TestUi) shouldNotHaveSaid(dialog string) {
	u.verify(dialog, false, "should not be said, but was said")
}

func (u *TestUi) verify(dialog string, expected bool, message string) {
	u.Logf("Verifying '%s'", dialog)
	if u.conversation[dialog] != expected {
		u.Errorf("'%s' %s", dialog, message)
	}
}
