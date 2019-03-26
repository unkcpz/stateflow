package flowmat

import (
  "testing"
  "io/ioutil"
)

// All test logs are discard
func initTestLogs() {
  InitLog(
    ioutil.Discard,
    ioutil.Discard,
  )
}

// Test Main
func TestMain(m *testing.M) {
  // initTestLogs()
  InitLogAudit()
  m.Run()
}
