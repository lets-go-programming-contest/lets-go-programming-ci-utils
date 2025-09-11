package module

import (
	"context"

	"github.com/lets-go-programming-contest/lets-go-programming-ci-utils/internal/config"
	"github.com/lets-go-programming-contest/lets-go-programming-ci-utils/internal/module"
	service "github.com/lets-go-programming-contest/lets-go-programming-ci-utils/internal/service/module"
	"github.com/spf13/cobra"
)

const (
	makeFileLintTarget = "lint"
	configFileName     = ".golangci.yaml"
)

var (
	lintFuncMapper = map[config.Mode]runEFunc{
		config.SkipMode:    runSkip(),
		config.CommonMode:  runMakeFromCommon(makeFileLintTarget),
		config.StudentMode: runMakeFromStudent(makeFileLintTarget),
		config.DefaultMode: runDefaultBuildCmd,
	}

	lintModeGetterFunc = func(config config.Config) config.Mode {
		return config.LintMode.Mode
	}
)

func newLintCmd() *cobra.Command {
	lintCmd := &cobra.Command{
		Use:   "lint",
		Short: "Lint student task",
		RunE:  selectorRun(lintFuncMapper, lintModeGetterFunc),
	}

	return lintCmd
}

func runDefaultLintCmd(cmd *cobra.Command, _ []string) error {
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

	if err := srv.RunLintModule(context.Background()); err != nil {
		return err
	}

	return nil
}
