package util

import (
	"fmt"
	"runtime"
)

func WrapErrorWithDetails(err error) error {
	// from svenwltr @ https://stackoverflow.com/a/38551362/17836168
	pc, _, _, ok := runtime.Caller(1)
	if ok {
		details := runtime.FuncForPC(pc)
		return fmt.Errorf("%s: error: %w", details, err)
	}
	return err
}

func LogDetailedError(err error) {
	pc, file, no, ok := runtime.Caller(1)
	if ok {
		details := runtime.FuncForPC(pc)
		fmt.Printf("%s#%d: %s: error: %s\n", file, no, details, err.Error())
	} else {
		fmt.Printf("LogDetailedError: Error in logger! Cannot log!")
	}
}
