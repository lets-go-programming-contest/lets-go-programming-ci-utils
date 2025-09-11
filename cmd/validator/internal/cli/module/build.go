package module

import (
	"context"

	"github.com/lets-go-programming-contest/lets-go-programming-ci-utils/internal/config"
	"github.com/lets-go-programming-contest/lets-go-programming-ci-utils/internal/module"
	service "github.com/lets-go-programming-contest/lets-go-programming-ci-utils/internal/service/module"
	"github.com/spf13/cobra"
)

const makeFileBuildTarget = "build"

var buildModeGetterFunc = func(config config.Config) config.Mode {
	return config.BuildMode.Mode
}

var buildFuncMapper = map[config.Mode]runEFunc{
	config.SkipMode:    runSkip(),
	config.CommonMode:  runMakeFromCommon(makeFileBuildTarget),
	config.StudentMode: runMakeFromStudent(makeFileBuildTarget),
	config.DefaultMode: runDefaultBuildCmd,
}

func newBuildCmd(opts ...func(cmd *cobra.Command)) *cobra.Command {
	buildCmd := &cobra.Command{
		Use:   "build",
		Short: "Build student task",
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
		return err
	}

	if err := srv.RunBuildModuleTargets(context.Background(), outputDir); err != nil {
		return err
	}

	return nil
}
