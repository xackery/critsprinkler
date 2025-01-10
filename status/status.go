package status

import "fmt"

var (
	status string
)

// Set sets the status
func Setf(format string, a ...interface{}) {
	status = fmt.Sprintf(format, a...)
}

// Set sets the status
func Set(line string) {
	status = line
}

// Status returns the current status
func String() string {
	return status
}
