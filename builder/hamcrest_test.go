package builder

import (
	"fmt"
	"reflect"
	"testing"
)

type Matcher interface {
	matches(interface{}) bool
	description(interface{}) string
}

type Expectation struct {
	actual interface{}
	*testing.T
}

func expec(t *testing.T) func(interface{}) *Expectation {
	return func(actual interface{}) *Expectation {
		return &Expectation{T: t, actual: actual}
	}
}

func (e *Expectation) to(matcher Matcher) {
	if !matcher.matches(e.actual) {
		e.Errorf(matcher.description(e.actual))
	}
}

func be(m Matcher) Matcher {
	return m
}

type InstanceOf struct {
	expected reflect.Type
}

func instanceOf(templateInstance interface{}) *InstanceOf {
	return &InstanceOf{reflect.TypeOf(templateInstance).Elem()}
}

func (m *InstanceOf) matches(actual interface{}) bool {
	return reflect.TypeOf(actual).Implements(m.expected)
}

func (m *InstanceOf) description(actual interface{}) string {
	return fmt.Sprintf("Expected %#v to be %v but found %[1]T", actual, m.expected)
}
