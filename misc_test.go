package flowmat

import (
  "testing"
  "io/ioutil"
  "os"
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
