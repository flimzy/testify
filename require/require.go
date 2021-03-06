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

// DeepEqualJSON marshals the expected and actual interfaces to JSON, then
// unmarshals before doing a reflect.DeepEqual check on them. If they are
// unequal, a diff of their respective JSON representations is produced as
// output.
func DeepEqualJSON(t TestingT, expected, actual interface{}, msgAndArgs ...interface{}) {
	if !assert.DeepEqualJSON(t, expected, actual, msgAndArgs...) {
		t.FailNow()
	}
}

// DeepEqualJSON marshals the expected and actual interfaces to JSON, then
// unmarshals before doing a reflect.DeepEqual check on them. If they are
// unequal, a diff of their respective JSON representations is produced as
// output.
func (a *Assertions) DeepEqualJSON(expected, actual interface{}, msgAndArgs ...interface{}) {
	DeepEqualJSON(a.t, expected, actual, msgAndArgs...)
}

// MarshalsToJSON asserts that the actual interface{} marshals to the expected
// JSON.
func MarshalsToJSON(t TestingT, expected []byte, actual interface{}, msgAndArgs ...interface{}) {
	if !assert.MarshalsToJSON(t, expected, actual, msgAndArgs...) {
		t.FailNow()
	}
}

// MarshalsToJSON asserts that the actual interface{} marshals to the expected
// JSON.
func (a *Assertions) MarshalsToJSON(expected []byte, actual interface{}, msgAndArgs ...interface{}) {
	MarshalsToJSON(a.t, expected, actual, msgAndArgs...)
}

// LinesEqual asserts that the two strings are equal, or shows a line-by-line
// diff of their differences.
func LinesEqual(t TestingT, expected, actual string, msgAndArgs ...interface{}) {
	if !assert.LinesEqual(t, expected, actual, msgAndArgs...) {
		t.FailNow()
	}
}

// LinesEqual asserts that the two strings are equal, or shows a line-by-line
// diff of their differences.
func (a *Assertions) LinesEqual(expected, actual string, msgAndArgs ...interface{}) {
	LinesEqual(a.t, expected, actual, msgAndArgs...)
}

// HTMLEqual asserts that the two arguments represent equivalent HTML. Accepts
// strings, byte arrays, *html.Node objects, or goquery selection.
func HTMLEqual(t TestingT, expected, actual interface{}, msgAndArgs ...interface{}) {
	if !assert.HTMLEqual(t, expected, actual, msgAndArgs...) {
		t.FailNow()
	}
}

// HTMLEqual asserts that the two arguments represent equivalent HTML. Accepts
// strings, byte arrays, *html.Node objects, or goquery selection.
func (a *Assertions) HTMLEqual(expected, actual interface{}, msgAndArgs ...interface{}) {
	HTMLEqual(a.t, expected, actual, msgAndArgs...)
}
