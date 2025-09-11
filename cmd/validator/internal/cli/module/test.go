package module

import (
	"context"
	"fmt"

	"github.com/lets-go-programming-contest/lets-go-programming-ci-utils/internal/config"
	"github.com/lets-go-programming-contest/lets-go-programming-ci-utils/internal/module"
	service "github.com/lets-go-programming-contest/lets-go-programming-ci-utils/internal/service/module"
	"github.com/spf13/cobra"
)

const (
	makeFileTestTarget = "test"
)

var (
	testFuncMapper = map[config.Mode]runEFunc{
		config.SkipMode:    runSkip(),
		config.CommonMode:  runMakeFromCommon(makeFileTestTarget),
		config.StudentMode: runMakeFromStudent(makeFileTestTarget),
		config.DefaultMode: runDefaultTestCmd,
	}

	testModeGetterFunc = func(cfg config.Config) config.Mode {
		return cfg.TestMode.Mode
	}
)

func newTestCmd() *cobra.Command {
	//nolint:exhaustruct // Set defaults values for another fields.
	testCmd := &cobra.Command{
		Use:   "test",
		Short: "Run tests for current module with using common tests",
		RunE:  selectorRun(testFuncMapper, testModeGetterFunc),
	}

	return testCmd
}

func runDefaultTestCmd(cmd *cobra.Command, _ []string) error {
	var (
		student = getStudentName(cmd.Flags())
		task    = getTaskName(cmd.Flags())
	)

	srv, err := service.NewService(
		student,
		task,
		module.WithTargetsCalculation(),
	)
	if err != nil {
		return fmt.Errorf("create service: %w", err)
	}

	if err := srv.RunTestModule(context.Background()); err != nil {
		return fmt.Errorf("test stage: %w", err)
	}

	return nil
}
