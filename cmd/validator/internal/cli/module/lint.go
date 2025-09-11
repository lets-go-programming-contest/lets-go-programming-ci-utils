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
	makeFileLintTarget = "lint"
	configFileName     = ".golangci.yaml"
)

var (
	lintFuncMapper = map[config.Mode]runEFunc{
		config.SkipMode:    runSkip(),
		config.CommonMode:  runMakeFromCommon(makeFileLintTarget),
		config.StudentMode: runMakeFromStudent(makeFileLintTarget),
		config.DefaultMode: runDefaultLintCmd,
	}

	lintModeGetterFunc = func(cfg config.Config) config.Mode {
		return cfg.LintMode.Mode
	}
)

func newLintCmd() *cobra.Command {
	//nolint:exhaustruct // Set defaults values for another fields.
	lintCmd := &cobra.Command{
		Use:   "lint",
		Short: "Lint current module code",
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
		return fmt.Errorf("create service: %w", err)
	}

	if err := srv.RunLintModule(context.Background()); err != nil {
		return fmt.Errorf("lint stage: %w", err)
	}

	return nil
}
