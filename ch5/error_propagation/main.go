package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime/debug"
)

type MyError struct {
	Inner      error
	Message    string
	StackTrace string
	Misc       map[string]interface{}
}

func (err MyError) Error() string {
	return err.Message
}

func WrapError(err error, messagef string, msgArgs ...interface{}) MyError {
	return MyError{
		Inner:      err,
		Message:    fmt.Sprintf(messagef, msgArgs...),
		StackTrace: string(debug.Stack()),
		Misc:       make(map[string]interface{}),
	}
}

// Low Level Error
type LowLevelErr struct {
	error
}

func isGlobalyExec(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, LowLevelErr{(WrapError(err, err.Error()))}
	}

	return info.Mode().Perm()&0100 == 0100, nil
}

// Intermediate Error
type intermediateErr struct {
	error
}

func runJob(id string) error {
	const jobBinPath = "/bad/job/binary"
	isExecutable, err := isGlobalyExec(jobBinPath)
	if err != nil {
		return intermediateErr{
			WrapError(err, "cannot run job %q: requisite binaries not available", id),
		}
	} else if isExecutable == false {
		return WrapError(nil, "cannot run job %q: requisite binaries are not executable", id)
	}

	return exec.Command(jobBinPath, "--id="+id).Run()
}

// top level
func handleError(key int, err error, message string) {
	log.SetPrefix(fmt.Sprintf("[logID: %v]:", key))
	log.Printf("%#v", err)
	log.Printf("[%v] %v", key, message)
}

func main() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime | log.LUTC)

	err := runJob("1")
	if err != nil {
		msg := "There was an unexpected issue; please report this is a bug"
		if _, ok := err.(intermediateErr); ok {
			msg = err.Error()
		}

		handleError(1, err, msg)
	}
}
