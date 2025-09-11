package module

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/lets-go-programming-contest/lets-go-programming-ci-utils/internal/executor"
)

type module struct {
	modulePath string
	targets    map[string]string
	executor   moduleExecutor
}

func checkModule(modulePath string) error {
	if _, err := os.Stat(filepath.Join(modulePath, "go.mod")); err != nil {
		return ErrPathNotModule
	}

	return nil
}

func calculateModuleTargets(modulePath string) (map[string]string, error) {
	moduleTargetsSubdir := filepath.Join(modulePath, "cmd")
	entries, err := os.ReadDir(moduleTargetsSubdir)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrCmdDirNotFound, err)
	}

	targets := make(map[string]string)

	for _, entry := range entries {
		if entry.IsDir() {
			absTargetPath, err := filepath.Abs(filepath.Join(modulePath, "cmd", entry.Name()))
			if err != nil {
				return nil, fmt.Errorf("%w: target path %q: %w", errCalculateAbsPath, entry, err)
			}
			targets[entry.Name()] = absTargetPath
		}
	}

	return targets, nil
}

type PlainModuleOpt func(mod *module) error

func WithTargetsCalculation() PlainModuleOpt {
	return func(mod *module) error {
		targets, err := calculateModuleTargets(mod.modulePath)
		if err != nil {
			return err
		}

		mod.targets = targets

		return nil
	}
}

func NewPlainModule(
	modulePath string,
	opts ...PlainModuleOpt,
) (module, error) {
	return NewPlainModuleWithExecutor(
		modulePath,
		executor.NewExecutor(
			executor.ExecWithOutput(os.Stdout),
			executor.ExecWithErrorOutput(os.Stderr),
		),
		opts...,
	)
}

func NewPlainModuleWithExecutor(
	modulePath string,
	baseExecutor baseExecutor,
	opts ...PlainModuleOpt,
) (module, error) {
	absModulePath, err := filepath.Abs(modulePath)
	if err != nil {
		return module{}, fmt.Errorf("%w: module path: %w", errCalculateAbsPath, err)
	}

	if err := checkModule(absModulePath); err != nil {
		return module{}, err
	}

	mod := module{
		modulePath: absModulePath,
		targets:    nil,
		executor:   newExecutor(baseExecutor),
	}

	for _, opt := range opts {
		if opt != nil {
			if err := opt(&mod); err != nil {
				return module{}, err
			}
		}
	}

	return mod, nil
}

func (m module) UpdateModuleDeps(
	ctx context.Context,
) error {
	ctx = executor.WithExecutorOpts(ctx, executor.ExecWithDir(m.modulePath))

	return m.executor.updateModuleDeps(ctx)
}

func (m module) BuildModuleTargets(
	ctx context.Context,
	outputDir string,
) error {
	ctx = executor.WithExecutorOpts(ctx, executor.ExecWithDir(m.modulePath))

	if err := m.executor.updateModuleDeps(ctx); err != nil {
		return err
	}

	if len(m.targets) == 0 {
		return ErrNothingToBuild
	}

	abOutputDir, err := filepath.Abs(outputDir)
	if err != nil {
		return fmt.Errorf("%w: output dir %q: %w", errCalculateAbsPath, outputDir, err)
	}

	if err := os.MkdirAll(abOutputDir, 0o755); err != nil {
		return fmt.Errorf("%w: output dir %q: %w", errPrepareOutputDir, abOutputDir, err)
	}

	for targetName, targetPath := range m.targets {
		if err := m.executor.buildTarget(ctx, abOutputDir, targetName, targetPath); err != nil {
			return err
		}
	}

	return nil
}

func (m module) TestModule(
	ctx context.Context,
) error {
	ctx = executor.WithExecutorOpts(ctx, executor.ExecWithDir(m.modulePath))

	if err := m.executor.updateModuleDeps(ctx); err != nil {
		return err
	}

	return m.executor.testModuleFiles(ctx)
}

func (m module) LintModule(
	ctx context.Context,
	config string,
) error {
	ctx = executor.WithExecutorOpts(ctx, executor.ExecWithDir(m.modulePath))

	if err := m.executor.updateModuleDeps(ctx); err != nil {
		return err
	}

	absConfigPath, err := filepath.Abs(config)
	if err != nil {
		return fmt.Errorf("%w: golangci-lint config %q: %w", errCalculateAbsPath, config, err)
	}

	return m.executor.lintModuleFiles(ctx, absConfigPath)
}

func (m module) RunMakeForModule(
	ctx context.Context,
	makefilePath string,
	target string,
	args ...string,
) error {
	ctx = executor.WithExecutorOpts(ctx, executor.ExecWithDir(m.modulePath))

	absMakeFilePath, err := filepath.Abs(makefilePath)
	if err != nil {
		return fmt.Errorf("%w: makefile %q: %w", err, makefilePath, err)
	}

	info, err := os.Stat(absMakeFilePath)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrGetMakefile, err)
	}

	if info.IsDir() {
		return ErrMakefilePathIsDir
	}

	return m.executor.runMakeForMakefile(
		ctx,
		absMakeFilePath,
		target,
		args...,
	)
}
