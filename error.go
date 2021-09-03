package errors

import (
	"bytes"
	stderrors "errors"
	"fmt"
	"io"
	"strings"
)

var (
	EMPTY_TRACE     = []callframe{}
	FILETER_WRAP_FN = "github.com/bencolove/errors.Wrap"
	FILETER_NEW_FN  = "github.com/bencolove/errors.New"
)

type stackerror struct {
	Message string
	Cause   error
	Traces  []callframe
}

func (e *stackerror) Error() string {
	return e.Message
}

func (e *stackerror) File() string {
	return e.Traces[0].file
}

func (e *stackerror) Line() int {
	return e.Traces[0].line
}

func (e *stackerror) Unwrap() error {
	return e.Cause
}

// // collect
// func (e *stackerror) WriteStacktrace(writer io.Writer) {

// 	var output = fmt.Fprintf

// 	if e.err == nil {
// 		return
// 	}

// 	var (
// 		loop error = e
// 	)

// 	for loop != nil {

// 		// test it against stack_error
// 		serr, ok := loop.(*stackerror)
// 		if ok {
// 			// write cause and stracktraces
// 			// cause
// 			output(writer, "caused by %T: %s\n", serr, serr.s)
// 			// stacktraces
// 			for _, frame := range e.st {
// 				output(writer, "    %s(%s:%d)", frame.file, frame.fn, frame.line)
// 			}

// 		} else {
// 			// fmt.Fprintf(&buf, " ... : %s\n", loop)
// 			output(writer, "caused by %T: %s\n", loop, loop.Error())
// 		}

// 		loop = stderrors.Unwrap(loop)
// 	}
// }

// func (e *stackerror) PrintStacktrace() {
// 	var buf bytes.Buffer
// 	e.WriteStacktrace(&buf)
// 	fmt.Println(buf.String())
// }

// Wrap works like fmt.Errorf() to wrap an existing error as cause making it `errors.Stackerror`
func Wrap(message string, cause error) error {
	frames := GetCallframes(2)

	// filter internal functions calls: function name startsWith 'github.com/bencolove/errors.Wrap'

	var filtered []callframe
	for _, frame := range frames {
		if strings.Compare(frame.fn, FILETER_WRAP_FN) != 0 &&
			strings.Compare(frame.fn, FILETER_NEW_FN) != 0 {
			filtered = append(filtered, frame)
		}
	}

	return &stackerror{
		Message: message,
		Cause:   cause,
		Traces:  filtered,
	}
}

// WriteStacktrace will try to write the chain of caused errors in the reversed order of calling frames
// if an error in the chain is of type `errors.*stackerror`, it will also write the including stacktraces
// otherwise it treat other error like a string (error message)
func WriteStacktrace(writer io.Writer, e error) {

	var output = fmt.Fprintf

	if e == nil {
		return
	}

	var (
		loop error = e
	)

	for loop != nil {

		// test it against stack_error
		serr, ok := loop.(*stackerror)
		if ok {
			// write cause and stracktraces
			// cause
			output(writer, "caused by %T: %s\n", serr, serr.Message)
			// stacktraces
			for _, frame := range serr.Traces {
				output(writer, "    %s:%d(%s)\n", frame.file, frame.line, frame.fn)
			}

		} else {
			// fmt.Fprintf(&buf, " ... : %s\n", loop)
			output(writer, "caused by %T: %s\n", loop, loop.Error())
		}

		loop = stderrors.Unwrap(loop)
	}
}

// PrintStacktrace outputs the stacktraces it finds on the error and print to stdout console by fmt
func PrintStacktrace(e error) {
	var buf bytes.Buffer
	WriteStacktrace(&buf, e)
	fmt.Println(buf.String())
}
