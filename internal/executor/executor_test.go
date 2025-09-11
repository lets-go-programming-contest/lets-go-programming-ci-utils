package executor_test

import (
	"context"
	"os"
	"testing"

	"github.com/lets-go-programming-contest/lets-go-programming-ci-utils/internal/executor"
	"github.com/stretchr/testify/require"
)

type exitMode string

const (
	okExitMode   exitMode = "ok"
	failExitMode exitMode = "fail"
)

//revive:disable:deep-exit // Use this test function for check exec wrapper. Exit statuses is required.
func TestHelperProcess(*testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}

	args := os.Args[3:]
	switch exitMode(args[0]) {
	case okExitMode:
		os.Stdout.WriteString("ok\n")
		os.Exit(0)
	case failExitMode:
		os.Stderr.WriteString("exit 42\n")
		os.Exit(42)
	default:
		os.Stderr.WriteString("unexpected exit\n")
		os.Exit(1)
	}
}

var casesExecCommand = []struct {
	name string
	mode exitMode
}{
	{
		name: "valid exec",
		mode: okExitMode,
	},
	{
		name: "invalid exec",
		mode: failExitMode,
	},
}

func TestExecCommand(t *testing.T) {
	t.Parallel()

	for _, tt := range casesExecCommand {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			testExecutor := executor.NewExecutor()
			annotatedCtx := executor.WithExecutorOpts(
				context.TODO(),
				executor.ExecWithEnvs(map[string]string{
					"GO_WANT_HELPER_PROCESS": "1",
				}),
			)
			args := []string{"-test.run=TestHelperProcess", "--", string(tt.mode)}

			err := testExecutor.ExecCommand(
				annotatedCtx,
				os.Args[0],
				args...,
			)

			if tt.mode == okExitMode {
				require.NoError(t, err)
			} else {
				var execError *executor.ExecError

				require.ErrorAs(t, err, &execError)
				require.Equal(t, os.Args[0], execError.CommandName)
				require.Equal(t, args, execError.CommandArgs)
			}
		})
	}
}
