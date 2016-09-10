package assert

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/pmezard/go-difflib/difflib"
	"github.com/stretchr/testify/assert"
)

// TestingT is an interface wrapper around *testing.T
type TestingT interface {
	Errorf(format string, args ...interface{})
	FailNow()
}

// Assertions provides assertion methods around the TestingT interface.
type Assertions struct {
	*assert.Assertions
	t TestingT
}

// New makes a new Assertions object for the specified TestingT.
func New(t TestingT) *Assertions {
	return &Assertions{
		Assertions: assert.New(t),
		t:          t,
	}
}

// FailDiff reports a failure through, including a contextual diff
func FailDiff(t TestingT, failureMessage, diff string, msgAndArgs ...interface{}) bool {
	if diff == "" {
		return Fail(t, failureMessage, msgAndArgs...)
	}
	message := messageFromMsgAndArgs(msgAndArgs...)

	errorTrace := strings.Join(assert.CallerInfo(), "\n\r\t\t\t")
	msg := fmt.Sprintf("\r%s\r\tError Trace:\t%s\n"+
		"\r\tError:%s\n",
		getWhitespaceString(),
		errorTrace,
		indentMessageLines(failureMessage, 2),
	)
	if len(diff) > 0 {
		msg = msg + fmt.Sprintf("\r\tDiff:%s\n\r",
			indentMessageLines(diff, 2),
		)
	}
	if len(message) > 0 {
		msg = msg + fmt.Sprintf("\r\tMessages:\t%s\n\r",
			message,
		)
	}

	t.Errorf(msg)
	return false
}

// Fail reports a failure through, with a
func Fail(t TestingT, failureMessage string, msgAndArgs ...interface{}) bool {
	return assert.Fail(t, failureMessage, msgAndArgs...)
}

// DeepEqual asserts that two objects are deeply equal.
func DeepEqual(t TestingT, expected, actual interface{}, msgAndArgs ...interface{}) bool {
	if reflect.DeepEqual(expected, actual) {
		return true
	}
	return FailDiff(t, "Structs differ", interfaceDiff(expected, actual), msgAndArgs...)
}

// DeepEqual asserts that two objects are deeply equal.
func (a *Assertions) DeepEqual(t TestingT, expected, actual interface{}, msgAndArgs ...interface{}) bool {
	return DeepEqual(a.t, expected, actual, msgAndArgs...)
}

var capRE = regexp.MustCompile("cap=[0-9]+\\)")
var capRepl = "cap=X"
var addRE = regexp.MustCompile("\\(0x[0-9a-f]{6,10}\\)")
var addRepl = "(0xXXXXXXXXXX)"

func diff(expected, actual string) string {
	udiff := difflib.UnifiedDiff{
		A:        strings.SplitAfter(expected, "\n"),
		FromFile: "expected",
		B:        strings.SplitAfter(actual, "\n"),
		ToFile:   "actual",
		Context:  2,
	}
	diff, err := difflib.GetUnifiedDiffString(udiff)
	if err != nil {
		panic("Error producing diff: " + err.Error())
	}
	return diff
}

func interfaceDiff(expected, actual interface{}) string {
	scs := spew.ConfigState{
		Indent:         "  ",
		DisableMethods: true,
		SortKeys:       true,
	}
	expString := scs.Sdump(expected)
	actString := scs.Sdump(actual)

	expString = capRE.ReplaceAllString(expString, capRepl)
	actString = capRE.ReplaceAllString(actString, capRepl)
	expString = addRE.ReplaceAllString(expString, addRepl)
	actString = addRE.ReplaceAllString(actString, addRepl)

	return diff(expString, actString)
}