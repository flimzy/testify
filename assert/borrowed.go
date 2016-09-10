package assert

import (
	"bufio"
	"bytes"
	"fmt"
	"runtime"
	"strings"
)

/*
Functions in this file are unexported methods borrowed heavily or entirely from
https://github.com/stretchr/testify
*/

// getWhitespaceString returns a string that is long enough to overwrite the default
// output from the go testing framework.
func getWhitespaceString() string {

	_, file, line, ok := runtime.Caller(1)
	if !ok {
		return ""
	}
	parts := strings.Split(file, "/")
	file = parts[len(parts)-1]

	return strings.Repeat(" ", len(fmt.Sprintf("%s:%d:      ", file, line)))

}

func messageFromMsgAndArgs(msgAndArgs ...interface{}) string {
	if len(msgAndArgs) == 0 || msgAndArgs == nil {
		return ""
	}
	if len(msgAndArgs) == 1 {
		return msgAndArgs[0].(string)
	}
	if len(msgAndArgs) > 1 {
		return fmt.Sprintf(msgAndArgs[0].(string), msgAndArgs[1:]...)
	}
	return ""
}

// Indents all lines of the message by appending a number of tabs to each line, in an output format compatible with Go's
// test printing (see inner comment for specifics)
func indentMessageLines(message string, tabs int) string {
	outBuf := new(bytes.Buffer)

	for i, scanner := 0, bufio.NewScanner(strings.NewReader(message)); scanner.Scan(); i++ {
		if i != 0 {
			outBuf.WriteRune('\n')
		}
		for ii := 0; ii < tabs; ii++ {
			outBuf.WriteRune('\t')
			// Bizarrely, all lines except the first need one fewer tabs prepended, so deliberately advance the counter
			// by 1 prematurely.
			if ii == 0 && i > 0 {
				ii++
			}
		}
		outBuf.WriteString(scanner.Text())
	}

	return outBuf.String()
}
