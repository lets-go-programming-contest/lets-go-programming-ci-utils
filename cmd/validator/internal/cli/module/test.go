package module

import (
	"context"

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
		config.DefaultMode: runDefaultBuildCmd,
	}

	testModeGetterFunc = func(config config.Config) config.Mode {
		return config.TestMode.Mode
	}
)

func newTestCmd() *cobra.Command {
	testCmd := &cobra.Command{
		Use:   "test",
		Short: "Test student task",
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
		return err
	}

	if err := srv.RunTestModule(context.Background()); err != nil {
		return err
	}

	return nil
}
