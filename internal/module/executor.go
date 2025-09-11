package module

import (
	"context"
	"fmt"
	"path/filepath"
)

type baseExecutor interface {
	ExecCommand(ctx context.Context, name string, args ...string) error
}

type moduleExecutor struct {
	baseExecutor baseExecutor
}

func newExecutor(baseExecutor baseExecutor) moduleExecutor {
	return moduleExecutor{
		baseExecutor: baseExecutor,
	}
}

func (e moduleExecutor) updateModuleDeps(ctx context.Context) error {
	if err := e.baseExecutor.ExecCommand(ctx, "go", "mod", "tidy"); err != nil {
		return fmt.Errorf("%w: %w", ErrUpdateModuleDeps, err)
	}

	return nil
}

func (e moduleExecutor) buildTarget(ctx context.Context,
	outputPath string,
	targetName string,
	targetPath string,
) error {
	if !filepath.IsAbs(outputPath) {
		return fmt.Errorf("%w: output path %q", errRelativePath, outputPath)
	}

	if !filepath.IsAbs(targetPath) {
		return fmt.Errorf("%w: target path %q", errRelativePath, targetPath)
	}

	if err := e.baseExecutor.ExecCommand(
		ctx,
		"go",
		"build",
		"-o", outputPath,
		targetPath,
	); err != nil {
		return fmt.Errorf("%w: %w", newBuildTargetError(targetName), err)
	}

	return nil
}

func (e moduleExecutor) lintModuleFiles(ctx context.Context, configPath string) error {
	if !filepath.IsAbs(configPath) {
		return fmt.Errorf("config path %q must be absolute", configPath)
	}

	if err := e.baseExecutor.ExecCommand(ctx, "golangci-lint", "run", "--config", configPath); err != nil {
		return fmt.Errorf("%w: %w", newLintError(configPath), err)
	}

	return nil
}

func (e moduleExecutor) testModuleFiles(ctx context.Context) error {
	if err := e.baseExecutor.ExecCommand(ctx, "go", "test", "-v", "--cover", "./..."); err != nil {
		return fmt.Errorf("%w: %w", ErrTestModule, err)
	}

	return nil
}

func (e moduleExecutor) runMakeForMakefile(ctx context.Context, makefilePath string, target string, args ...string) error {
	if !filepath.IsAbs(makefilePath) {
		return fmt.Errorf("%w: makefile path %q", errRelativePath, makefilePath)
	}

	if err := e.baseExecutor.ExecCommand(ctx, "make", append([]string{"-f", makefilePath, target}, args...)...); err != nil {
		return fmt.Errorf("%w: %w", newRunMakeTargetError(makefilePath, target), err)
	}

	return nil
}
