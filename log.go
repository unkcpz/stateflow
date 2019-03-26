package gflow

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
)

// Audit implement from SciPipe

var (
	// Trace is a log hander for extremely detailed level logs.
	Trace *log.Logger
	// Audit is a log handler for audit level logs
	Audit     *log.Logger
	logExists bool
)

// InitLog initiates logging handlers
func InitLog(
	traceHandle io.Writer,
	auditHandle io.Writer) {

	if !logExists {
		Trace = log.New(traceHandle,
			"TRACE  ",
			log.Ldate|log.Ltime|log.Lshortfile)

		Audit = log.New(auditHandle,
			"AUDIT  ",
			log.Ldate|log.Ltime)

		logExists = true
	}
}

// InitLogAudit initiate logging with level=AUDIT
func InitLogAudit() {
	InitLog(
		ioutil.Discard,
		os.Stdout,
	)
}

// auditLogPattern decides the format of Audit log messages across package
const auditLogPattern = "| %-16s | %s\n"

// LogAuditf logs a pretty printed log message with the AUDIT log level, where
// componentName is a name of a process, task, workflow or similar that
// generates the message, while message and values are formatted in the manner
// of fmt.Printf
func LogAuditf(componentName string, message string, values ...interface{}) {
	Audit.Printf(auditLogPattern, componentName, fmt.Sprintf(message, values...))
}
