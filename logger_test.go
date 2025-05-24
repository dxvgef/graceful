package graceful

import (
	"io"
	"log"
)

type testLogger struct {
	stdLog *log.Logger
}

func testNewLogger(out io.Writer, prefix string, flag int) *testLogger {
	l := log.New(out, prefix, flag)
	return &testLogger{
		stdLog: l,
	}
}

func (l *testLogger) Output(v ...interface{}) {
	l.stdLog.Println(v...)
}
