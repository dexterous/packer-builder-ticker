package builder

import (
	"fmt"
	"github.com/mitchellh/packer/packer"
	"testing"
)

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
