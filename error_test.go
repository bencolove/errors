package errors

import (
	"bytes"
	stderrors "errors"
	"strings"
	"testing"
)

func TestWrap(t *testing.T) {
	stdroot := stderrors.New("root std error")

	// wrap another error
	wrapString := "it is a Wrapped error"
	err := Wrap(wrapString, stdroot)

	if !stderrors.Is(err, stdroot) {
		t.Errorf("%v should be %v", err, stdroot)
	}

	var foundErr *stackerror
	if !stderrors.As(err, &foundErr) {
		t.Errorf("%v should be of type Stackerror", err)
	}

	if !strings.HasSuffix(foundErr.Message, wrapString) {
		t.Errorf("Error message should endswith: %s", wrapString)
	}
}

func errorFn1() error {
	return New("root cause")
}

func errorFn2() error {
	m := "second cause"
	err := errorFn1()

	if err != nil {
		secErr := Wrap(m, err)
		return secErr
	}

	return err
}
func TestStacktrace(t *testing.T) {
	err := errorFn2()

	var buf bytes.Buffer

	WriteStacktrace(&buf, err)

	t.Logf("%s", buf.String())
}
