package util

import (
	"fmt"
	"runtime"
)

// GetDetailedProgramLocation TODO
func _() string {
	// from svenwltr @ https://stackoverflow.com/a/38551362/17836168
	pc, file, no, ok := runtime.Caller(1)
	if ok {
		details := runtime.FuncForPC(pc)
		return fmt.Sprintf("%s#%d: %s", file, no, details)
	}
	return ""
}

func LogDetailedError(err error) {
	pc, file, no, ok := runtime.Caller(1)
	if ok {
		details := runtime.FuncForPC(pc)
		fmt.Printf("%s#%d: %s: Error: %s\n", file, no, details, err.Error())
	} else {
		fmt.Printf("LogDetailedError: Error in logger! Cannot log!")
	}
}
