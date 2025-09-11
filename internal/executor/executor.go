package executor

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"syscall"
)

type ExecutorOpts func(cmd *exec.Cmd)

type executor struct {
	opts []ExecutorOpts
}

func NewExecutor(opts ...ExecutorOpts) executor {
	return executor{
		opts: opts,
	}
}

func ExecWithDir(dir string) ExecutorOpts {
	return func(cmd *exec.Cmd) {
		cmd.Dir = dir
	}
}

func ExecWithOutput(stdout io.Writer) ExecutorOpts {
	return func(cmd *exec.Cmd) {
		cmd.Stdout = stdout
	}
}

func ExecWithErrorOutput(stderr io.Writer) ExecutorOpts {
	return func(cmd *exec.Cmd) {
		cmd.Stderr = stderr
	}
}

func ExecWithUserCredential(uid, gid uint32) ExecutorOpts {
	return func(cmd *exec.Cmd) {
		cmd.SysProcAttr = &syscall.SysProcAttr{
			Credential: &syscall.Credential{
				Uid: uid,
				Gid: gid,
			},
		}
	}
}

func ExecWithEnvs(envsMapper map[string]string) ExecutorOpts {
	return func(cmd *exec.Cmd) {
		envsSlice := make([]string, 0, len(envsMapper))
		for name, val := range envsMapper {
			envsSlice = append(envsSlice, fmt.Sprintf("%s=%s", name, val))
		}
		cmd.Env = append(os.Environ(), envsSlice...)
	}
}

func ExecWithCurrentUser() ExecutorOpts {
	return ExecWithUserCredential(
		uint32(syscall.Getuid()),
		uint32(syscall.Getgid()),
	)
}

func (e executor) ExecCommand(
	ctx context.Context,
	name string,
	args ...string,
) error {
	cmd := exec.CommandContext(ctx, name, args...)

	opts := make([]ExecutorOpts, 0, len(e.opts))
	opts = append(opts, e.opts...)
	opts = append(opts, GetExecutorOpts(ctx)...)

	for _, opt := range opts {
		if opt != nil {
			opt(cmd)
		}
	}

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%w: %w", NewExecError(name, args), err)
	}

	return nil
}
