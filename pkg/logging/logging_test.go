package logging

import (
	"errors"
	"testing"
)

// TestTrace calls LogTrace(msg string)
func TestTrace(t *testing.T) {
	l := NewLogger()
	l.LogTrace("Hello Trace")
	if l.HadError() {
		t.Fatalf(`LogTrace(string) had an error. HadError was %v, want false`, l.HadError())
	}
}

// TestTrace calls LogTrace(msg string)
func TestInfo(t *testing.T) {
	l := NewLogger()
	l.LogInfo("Hello Info")
	if l.HadError() {
		t.Fatalf(`LogInfo(string) had an error. HadError was %v, want false`, l.HadError())
	}
}

// TestTrace calls LogTrace(msg string)
func TestTraceError(t *testing.T) {
	l := NewLogger()
	e := errors.New("Hello error")
	l.LogError(e, "")
	if l.HadError() {
		t.Fatalf(`LogError(error string) had an error. HadError was %v, want false`, l.HadError())
	}
}
