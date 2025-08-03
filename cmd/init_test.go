package cmd

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"os"
	"os/exec"
	"testing"
)

func TestRunError(t *testing.T) {
	cmd := exec.Command("/bin/sh", "-c", "fnord")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	_, isExit := err.(*exec.ExitError)
	if isExit {
		fmt.Printf("IS ExitError: exitCode=%d\n", cmd.ProcessState.ExitCode())
		//os.Exit(cmdProcessSate.ExitCode())
	} else {
		fmt.Printf("NOT ExitError: %+v\n", err)
	}
	require.NotNil(t, err)
}
