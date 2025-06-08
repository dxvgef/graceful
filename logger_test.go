package graceful

import (
	"testing"
)

type testLogger struct {
	t *testing.T
}

func testNewLogger(t *testing.T) *testLogger {
	return &testLogger{
		t: t,
	}
}

func (tl *testLogger) Output(msg any) {
	tl.t.Log(msg)
}
