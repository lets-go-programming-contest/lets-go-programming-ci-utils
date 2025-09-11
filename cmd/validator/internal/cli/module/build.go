package module

import (
	"context"
	"fmt"

	"github.com/lets-go-programming-contest/lets-go-programming-ci-utils/internal/config"
	"github.com/lets-go-programming-contest/lets-go-programming-ci-utils/internal/module"
	service "github.com/lets-go-programming-contest/lets-go-programming-ci-utils/internal/service/module"
	"github.com/spf13/cobra"
)

const makeFileBuildTarget = "build"

var buildModeGetterFunc = func(cfg config.Config) config.Mode {
	return cfg.BuildMode.Mode
}

var buildFuncMapper = map[config.Mode]runEFunc{
	config.SkipMode:    runSkip(),
	config.CommonMode:  runMakeFromCommon(makeFileBuildTarget),
	config.StudentMode: runMakeFromStudent(makeFileBuildTarget),
	config.DefaultMode: runDefaultBuildCmd,
}

func newBuildCmd(opts ...func(cmd *cobra.Command)) *cobra.Command {
	//nolint:exhaustruct // Set defaults values for another fields.
	buildCmd := &cobra.Command{
		Use:   "build",
		Short: "Build all targets in current module",
		RunE:  selectorRun(buildFuncMapper, buildModeGetterFunc),
	}

	for _, opt := range opts {
		if opt != nil {
			opt(buildCmd)
		}
	}

	return buildCmd
}

func runDefaultBuildCmd(cmd *cobra.Command, _ []string) error {
	var (
		student   = getStudentName(cmd.Flags())
		task      = getTaskName(cmd.Flags())
		outputDir = getOutputDir(cmd.Flags())
	)

	srv, err := service.NewService(
		student,
		task,
		module.WithTargetsCalculation(),
	)
	if err != nil {
		return fmt.Errorf("create service: %w", err)
	}

	if err := srv.RunBuildModuleTargets(context.Background(), outputDir); err != nil {
		return fmt.Errorf("build stage: %w", err)
	}

	return nil
}
