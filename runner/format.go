package runner

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
)

var (
	goImportsPath = getGoImportsPath() // path to `goimports`
)

// Gets the path for the `goimports` executable,
// checking inside $GOPATH/bin if $GOBIN is unset.
func getGoImportsPath() string {
	var goBinPath string
	if v, ok := os.LookupEnv("GOBIN"); ok {
		goBinPath = v
	} else {
		goBinPath = path.Join(os.Getenv("GOPATH"), "bin")
	}

	return path.Join(goBinPath, "goimports")
}

// FormatGo runs `goimports` on a Go program as a string
// and returns the formatted result as well as an error containing
// the command's STDERR if a the command exited with a non-zero code.
func FormatGo(text string) (string, error) {
	cmd := exec.Command(goImportsPath)
	var stdIn, stdOut, stdErr bytes.Buffer
	cmd.Stdin, cmd.Stdout, cmd.Stderr = &stdIn, &stdOut, &stdErr

	_, err := stdIn.WriteString(text)
	if err != nil {
		log.Fatalf("Failed to write to goimports STDIN: %v\n", err)
	}

	err = cmd.Run()
	if _, ok := err.(*exec.ExitError); ok {
		return "", fmt.Errorf("%v - %s", err, stdErr.String())
	}

	if err != nil {
		log.Fatalf("Failed to run `%s`: %v\n", cmd, err)
	}

	return stdOut.String(), nil
}
