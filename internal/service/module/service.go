package module

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/lets-go-programming-contest/lets-go-programming-ci-utils/internal/envs"
	"github.com/lets-go-programming-contest/lets-go-programming-ci-utils/internal/file"
	"github.com/lets-go-programming-contest/lets-go-programming-ci-utils/internal/module"
)

const (
	makefileName           = "Makefile"
	golangciLintConfigName = ".golangci.yaml"
	commonTestsSourcePath  = "tests"
	commonTestsTargetPath  = "tests/ci"
)

var (
	errGetEnviroment  = errors.New("get environment variables")
	ErrRunMakeTarget  = errors.New("run target from makefile")
	ErrBuildTargets   = errors.New("build module targets")
	ErrLintTargets    = errors.New("lint module")
	ErrCopyCommonTest = errors.New("copy test from common dir")
	ErrTestTargets    = errors.New("test module")
)

type serviceModule interface {
	BuildModuleTargets(ctx context.Context, outputDir string) error
	LintModule(ctx context.Context, config string) error
	RunMakeForModule(ctx context.Context, makefilePath string, target string, args ...string) error
	TestModule(ctx context.Context) error
	UpdateModuleDeps(ctx context.Context) error
}

type service struct {
	mod         serviceModule
	studentName string
	taskName    string
}

func NewServiceWithModule(
	studentName string,
	taskName string,
	serviceModule serviceModule,
) (service, error) {
	return service{
		mod:         serviceModule,
		studentName: studentName,
		taskName:    taskName,
	}, nil
}

func NewService(
	studentName string,
	taskName string,
	opts ...module.PlainModuleOpt,
) (service, error) {
	mod, err := module.NewPlainModule(path.Join(studentName, taskName), opts...)
	if err != nil {
		return service{}, fmt.Errorf("get go module: %w", err)
	}

	return NewServiceWithModule(studentName, taskName, mod)
}

func (s service) RunMakeFromCommon(
	ctx context.Context,
	target string,
	args ...string,
) error {
	commonDir, err := envs.GetCommonDirFromEnv()
	if err != nil {
		return errors.Join(errGetEnviroment, err)
	}

	if err := s.mod.RunMakeForModule(
		ctx,
		filepath.Join(commonDir, s.taskName, makefileName),
		target,
		args...,
	); err != nil {
		return fmt.Errorf("run target from common files: %w", err)
	}

	return nil
}

func (s service) RunMakeFromStudent(
	ctx context.Context,
	target string,
	args ...string,
) error {
	if err := s.mod.RunMakeForModule(
		ctx,
		filepath.Join(s.studentName, s.taskName, makefileName),
		target,
		args...,
	); err != nil {
		return fmt.Errorf("run target from studet task files: %w", err)
	}

	return nil
}

func (s service) RunBuildModuleTargets(
	ctx context.Context,
	outputDir string,
) error {
	if err := s.mod.BuildModuleTargets(ctx, outputDir); err != nil {
		return fmt.Errorf("%w: %w", ErrBuildTargets, err)
	}

	return nil
}

func getGolangciLintConfigPaths(
	dirs ...string,
) []string {
	lintConfigFiles := make([]string, 0, len(dirs))

	for _, dir := range dirs {
		configPath := filepath.Join(dir, golangciLintConfigName)

		if _, err := os.Stat(configPath); err == nil {
			lintConfigFiles = append(lintConfigFiles, configPath)
		} else {
			fmt.Printf("Config %q not found.\tSkip.\n", configPath)
		}
	}

	return lintConfigFiles
}

func (s service) RunLintModule(
	ctx context.Context,
) error {
	commonDir, err := envs.GetCommonDirFromEnv()
	if err != nil {
		return fmt.Errorf("%w: %w", errGetEnviroment, err)
	}

	golangciLintConfigs := getGolangciLintConfigPaths(
		filepath.Join(commonDir, s.taskName, golangciLintConfigName),
		filepath.Join(s.studentName, s.taskName, golangciLintConfigName),
	)

	for _, config := range golangciLintConfigs {
		fmt.Printf("Run linters with %q config file.\n", config)

		if err := s.mod.LintModule(ctx, config); err != nil {
			return fmt.Errorf("%w: %w", ErrLintTargets, err)
		}
	}

	return nil
}

func (s service) RunTestModule(
	ctx context.Context,
) error {
	commonDir, err := envs.GetCommonDirFromEnv()
	if err != nil {
		return fmt.Errorf("%w: %w", errGetEnviroment, err)
	}

	err = file.CopyDir(
		filepath.Join(commonDir, s.taskName, commonTestsSourcePath),
		filepath.Join(s.studentName, s.taskName, commonTestsTargetPath),
	)

	switch {
	case err == nil:
	case os.IsNotExist(err):
		fmt.Printf("Common tests for task %q not found.\tSkip.\n", s.taskName)
	default:
		return fmt.Errorf("%w: %w", ErrCopyCommonTest, err)
	}

	if err := s.mod.TestModule(ctx); err != nil {
		return fmt.Errorf("%w: %w", ErrTestTargets, err)
	}

	return nil
}
