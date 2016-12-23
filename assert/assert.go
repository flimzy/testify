package assert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"golang.org/x/net/html"

	"github.com/PuerkitoBio/goquery"
	"github.com/davecgh/go-spew/spew"
	"github.com/pkg/errors"
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

	errorTrace := strings.Join(assert.CallerInfo(), "\n\t\t\t")
	msg := fmt.Sprintf("%s\n\tError Trace:\t%s\n\tError:%s\n",
		getWhitespaceString(),
		errorTrace,
		indentMessageLines(failureMessage, 2),
	)
	if len(diff) > 0 {
		msg = msg + fmt.Sprintf("\tDiff:\n\r\t%s\n",
			indentMessageLines(diff, 3),
		)
	}
	if len(message) > 0 {
		msg = msg + fmt.Sprintf("\tMessages:\t%s\n",
			message,
		)
	}

	t.Errorf(msg)
	return false
}

// Fail reports a failure through
func Fail(t TestingT, failureMessage string, msgAndArgs ...interface{}) bool {
	message := messageFromMsgAndArgs(msgAndArgs...)

	errorTrace := strings.Join(assert.CallerInfo(), "\n\t\t\t")
	if len(message) > 0 {
		t.Errorf("%s\tError Trace:\t%s\n"+
			"\tError:%s\n"+
			"\tMessages:\t%s\n",
			getWhitespaceString(),
			errorTrace,
			indentMessageLines(failureMessage, 2),
			message)
	} else {
		t.Errorf("%s\tError Trace:\t%s\n"+
			"\tError:%s\n",
			getWhitespaceString(),
			errorTrace,
			indentMessageLines(failureMessage, 2))
	}

	return false
}

// var capRE = regexp.MustCompile("cap=[0-9]+\\)")
// var capRepl = "cap=X"
// var addRE = regexp.MustCompile("\\(0x[0-9a-f]{6,10}\\)")
// var addRepl = "(0xXXXXXXXXXX)"

func diff(expected, actual string) string {
	if !strings.HasSuffix(expected, "\n") {
		expected = expected + "\n"
	}
	if !strings.HasSuffix(actual, "\n") {
		actual = actual + "\n"
	}
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

	// expString = capRE.ReplaceAllString(expString, capRepl)
	// actString = capRE.ReplaceAllString(actString, capRepl)
	// expString = addRE.ReplaceAllString(expString, addRepl)
	// actString = addRE.ReplaceAllString(actString, addRepl)

	return diff(expString, actString)
}

// DeepEqual asserts that two objects are deeply equal.
func DeepEqual(t TestingT, expected, actual interface{}, msgAndArgs ...interface{}) bool {
	if reflect.DeepEqual(expected, actual) {
		return true
	}
	return FailDiff(t, "Structs differ", interfaceDiff(expected, actual), msgAndArgs...)
}

// DeepEqual asserts that two objects are deeply equal.
func (a *Assertions) DeepEqual(expected, actual interface{}, msgAndArgs ...interface{}) bool {
	return DeepEqual(a.t, expected, actual, msgAndArgs...)
}

// DeepEqualJSON marshals the expected and actual interfaces to JSON, then
// unmarshals before doing a reflect.DeepEqual check on them. If they are
// unequal, a diff of their respective JSON representations is produced as
// output.
func DeepEqualJSON(t TestingT, expected, actual interface{}, msgAndArgs ...interface{}) bool {
	expectedJSON := marshalJSON(t, expected, msgAndArgs...)
	actualJSON := marshalJSON(t, actual, msgAndArgs...)
	var e, a interface{}
	json.Unmarshal(expectedJSON, &e)
	json.Unmarshal(actualJSON, &a)
	if reflect.DeepEqual(e, a) {
		return true
	}
	return FailDiff(t, "JSON representations differ", diff(string(expectedJSON), string(actualJSON)), msgAndArgs...)
}

// DeepEqualJSON marshals the expected and actual interfaces to JSON, then
// unmarshals before doing a reflect.DeepEqual check on them. If they are
// unequal, a diff of their respective JSON representations is produced as
// output.
func (a *Assertions) DeepEqualJSON(expected, actual interface{}, msgAndArgs ...interface{}) bool {
	return DeepEqualJSON(a.t, expected, actual, msgAndArgs...)
}

func marshalJSON(t TestingT, i interface{}, msgAndArgs ...interface{}) []byte {
	output, err := json.MarshalIndent(i, "", "    ")
	if err != nil {
		Fail(t, fmt.Sprintf("Error marshaling JSON: %s\n", err), msgAndArgs...)
	}
	return output
}

// MarshalsToJSON asserts that the actual interface{} marshals to the expected
// JSON.
func MarshalsToJSON(t TestingT, expected []byte, actual interface{}, msgAndArgs ...interface{}) bool {
	actualJSON := marshalJSON(t, actual, msgAndArgs...)
	var e, a interface{}
	if err := json.Unmarshal(expected, &e); err != nil {
		return Fail(t, "Error unmarshaling expected JSON string", msgAndArgs...)
	}
	json.Unmarshal(actualJSON, &a)
	if reflect.DeepEqual(e, a) {
		return true
	}
	return FailDiff(t, "JSON representations differ", diff(string(expected), string(actualJSON)), msgAndArgs...)
}

// MarshalsToJSON asserts that the actual interface{} marshals to the expected
// JSON.
func (a *Assertions) MarshalsToJSON(expected []byte, actual interface{}, msgAndArgs ...interface{}) bool {
	return MarshalsToJSON(a.t, expected, actual, msgAndArgs...)
}

// LinesEqual asserts that the two strings are equal, or shows a line-by-line
// diff of their differences.
func LinesEqual(t TestingT, expected, actual string, msgAndArgs ...interface{}) bool {
	if expected == actual {
		return true
	}
	return FailDiff(t, "Strings differ", diff(expected, actual), msgAndArgs...)
}

// LinesEqual asserts that the two strings are equal, or shows a line-by-line
// diff of their differences.
func (a *Assertions) LinesEqual(expected, actual string, msgAndArgs ...interface{}) bool {
	return LinesEqual(a.t, expected, actual, msgAndArgs...)
}

// HTMLEqual asserts that the two arguments represent equivalent HTML. Accepts
// strings, byte arrays, *html.Node objects, or goquery selection.
func HTMLEqual(t TestingT, expected, actual interface{}, msgAndArgs ...interface{}) bool {
	expDoc, err := toHTMLNode(expected)
	if err != nil {
		t.Errorf("invalid expected document: %s", err)
		t.FailNow()
	}
	actDoc, err := toHTMLNode(actual)
	if err != nil {
		t.Errorf("invalid actual document: %s", err)
		t.FailNow()
	}
	if !reflect.DeepEqual(expDoc, actDoc) {
		expBuf := new(bytes.Buffer)
		html.Render(expBuf, expDoc)
		actBuf := new(bytes.Buffer)
		html.Render(actBuf, actDoc)
		return FailDiff(t, "HTML differs", diff(expBuf.String(), actBuf.String()), msgAndArgs...)
	}
	return true
}

func toHTMLNode(i interface{}) (*html.Node, error) {
	switch i.(type) {
	case *html.Node:
		return i.(*html.Node), nil
	case string:
		r := strings.NewReader(i.(string))
		return html.Parse(r)
	case []byte:
		r := bytes.NewReader(i.([]byte))
		return html.Parse(r)
	case *goquery.Selection:
		str, err := goquery.OuterHtml(i.(*goquery.Selection))
		if err != nil {
			return nil, errors.Wrap(err, "failed to get outer html")
		}
		r := strings.NewReader(str)
		return html.Parse(r)
	}
	return nil, errors.Errorf("unknown type: %T", i)
}
