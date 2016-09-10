package require

import (
	"github.com/stretchr/testify/require"

	"github.com/flimzy/testify/assert"
)

// TestingT is an interface wrapper around *testing.T
type TestingT interface {
	Errorf(format string, args ...interface{})
	FailNow()
}

// Assertions provides assertion methods around the TestingT interface.
type Assertions struct {
	*require.Assertions
	t TestingT
}

// New makes a new Assertions object for the specified TestingT.
func New(t TestingT) *Assertions {
	require := require.New(t)
	return &Assertions{
		Assertions: require,
		t:          t,
	}
}

// DeepEqual asserts that two objects are deeply equal.
func DeepEqual(t TestingT, expected, actual interface{}, msgAndArgs ...interface{}) {
	if !assert.DeepEqual(t, expected, actual, msgAndArgs...) {
		t.FailNow()
	}
}

// DeepEqual asserts that two objects are deeply equal.
func (a *Assertions) DeepEqual(expected, actual interface{}, msgAndArgs ...interface{}) {
	DeepEqual(a.t, expected, actual, msgAndArgs...)
}
