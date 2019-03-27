package stateflow

import (
	"io/ioutil"
	"os"
	"testing"
)

// All test logs are discard
func TraceTestLogs() {
	InitLog(
		os.Stdout,
		ioutil.Discard,
	)
}

// Test Main
func TestMain(m *testing.M) {
	// InitLogAudit()
	TraceTestLogs()
	m.Run()
}
