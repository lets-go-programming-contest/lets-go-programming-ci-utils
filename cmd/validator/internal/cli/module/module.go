package module

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/lets-go-programming-contest/lets-go-programming-ci-utils/internal/envs"
	"github.com/lets-go-programming-contest/lets-go-programming-ci-utils/internal/module"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const cliOutputNameFlag = "output"

type (
	errMakeTargetFromCommonFiles  error
	errMakeTargetFromStudentFiles error
)

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}

	return v
}

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

func getStudentName(flagSet *pflag.FlagSet) string {
	return must(flagSet.GetString("student"))
}

func getTaskName(flagSet *pflag.FlagSet) string {
	return must(flagSet.GetString("task"))
}

func getConfigPath(flagSet *pflag.FlagSet) string {
	return must(flagSet.GetString("config"))
}

func getOutputDir(flagSet *pflag.FlagSet) string {
	return must(flagSet.GetString("output"))
}

func runSkip() func(*cobra.Command, []string) error {
	return func(_ *cobra.Command, _ []string) error {
		fmt.Println("Skip mode is set. Step skipped.")

		return nil
	}
}

func NewModuleCmd(opts ...func(cmd *cobra.Command)) *cobra.Command {
	moduleCmd := &cobra.Command{
		Use: "module",
	}

	moduleCmd.AddCommand(newBuildCmd(initOutputFlag))
	moduleCmd.AddCommand(newLintCmd())
	moduleCmd.AddCommand(newTestCmd())

	for _, opt := range opts {
		if opt != nil {
			opt(moduleCmd)
		}
	}

	return moduleCmd
}

func initOutputFlag(cmd *cobra.Command) {
	cmd.Flags().String(cliOutputNameFlag, "bin", "Output dir for binary files")
	panicIfErr(cmd.MarkFlagDirname(cliOutputNameFlag))
}

const makefileName = "Makefile"

func runMakeFromCommon(
	target string,
) func(cmd *cobra.Command, _ []string) error {
	return func(cmd *cobra.Command, _ []string) error {
		commonDirPath, err := envs.GetCommonDirFromEnv()
		if err != nil {
			return err
		}

		var (
			student      = getStudentName(cmd.Flags())
			task         = getTaskName(cmd.Flags())
			makeFilePath = filepath.Join(commonDirPath, task, makefileName)
			ctx          = context.Background()
		)

		if err := runMakeTarget(ctx, student, task, makeFilePath, target); err != nil {
			return processErr(student, task, errMakeTargetFromCommonFiles(err))
		}

		return nil
	}
}

func runMakeFromStudent(
	target string,
) func(cmd *cobra.Command, _ []string) error {
	return func(cmd *cobra.Command, _ []string) error {
		var (
			student      = getStudentName(cmd.Flags())
			task         = getTaskName(cmd.Flags())
			makeFilePath = filepath.Join(student, task, "Makefile")
			ctx          = context.Background()
		)

		if err := runMakeTarget(ctx, student, task, makeFilePath, target); err != nil {
			return processErr(student, task, errMakeTargetFromStudentFiles(err))
		}

		return nil
	}
}

func runMakeTarget(
	ctx context.Context,
	student string,
	task string,
	makefilePath string,
	target string,
) error {
	modulePath := filepath.Join(student, task)

	mod, err := module.NewPlainModule(modulePath)
	if err != nil {
		return fmt.Errorf("get go module: %w", err)
	}

	if err := mod.RunMakeForModule(context.Background(), makefilePath, target); err != nil {
		return fmt.Errorf("run make target: %w", err)
	}

	return nil
}
