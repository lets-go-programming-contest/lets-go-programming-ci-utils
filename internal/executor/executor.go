package executor

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
)

type Opts func(cmd *exec.Cmd)

type executor struct {
	opts []Opts
}

func NewExecutor(opts ...Opts) executor {
	return executor{
		opts: opts,
	}
}

func ExecWithDir(dir string) Opts {
	return func(cmd *exec.Cmd) {
		cmd.Dir = dir
	}
}

func ExecWithOutput(stdout io.Writer) Opts {
	return func(cmd *exec.Cmd) {
		cmd.Stdout = stdout
	}
}

func ExecWithErrorOutput(stderr io.Writer) Opts {
	return func(cmd *exec.Cmd) {
		cmd.Stderr = stderr
	}
}

func ExecWithEnvs(envsMapper map[string]string) Opts {
	return func(cmd *exec.Cmd) {
		envsSlice := make([]string, 0, len(envsMapper))
		for name, val := range envsMapper {
			envsSlice = append(envsSlice, fmt.Sprintf("%s=%s", name, val))
		}

		cmd.Env = append(os.Environ(), envsSlice...)
	}
}

func (e executor) ExecCommand(
	ctx context.Context,
	name string,
	args ...string,
) error {
	cmd := exec.CommandContext(ctx, name, args...)

	opts := make([]Opts, 0, len(e.opts))
	opts = append(opts, e.opts...)
	opts = append(opts, GetExecutorOpts(ctx)...)

	for _, opt := range opts {
		if opt != nil {
			opt(cmd)
		}
	}

	fmt.Printf("⇒ The following output refers to the execution of a command %q in an external process.\n", name)
	defer fmt.Printf("⇒ End of output from a command run in an external process %q.\n", name)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%w: %w", NewExecError(name, args), err)
	}

	return nil
}
