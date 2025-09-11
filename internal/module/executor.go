package module

import (
	"context"
	"fmt"
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
		return fmt.Errorf("update module deps: %w", err)
	}

	return nil
}

func (e moduleExecutor) buildTarget(ctx context.Context,
	outputPath string,
	targetName string,
	targetPath string,
) error {
	if err := e.baseExecutor.ExecCommand(
		ctx,
		"go",
		"build",
		"-o", outputPath,
		targetPath,
	); err != nil {
		return fmt.Errorf("build module target %q: %w", targetName, err)
	}

	return nil
}

func (e moduleExecutor) lintModuleFiles(
	ctx context.Context,
	configPath string,
) error {
	if err := e.baseExecutor.ExecCommand(ctx, "golangci-lint", "run", "--config", configPath); err != nil {
		return fmt.Errorf("lint module with using config path %q: %w", configPath, err)
	}

	return nil
}

func (e moduleExecutor) testModuleFiles(ctx context.Context) error {
	if err := e.baseExecutor.ExecCommand(ctx, "go", "test", "-v", "--cover", "./..."); err != nil {
		return fmt.Errorf("test module: %w", err)
	}

	return nil
}

func (e moduleExecutor) runMakeForMakefile(
	ctx context.Context,
	makefilePath string,
	target string,
	args ...string,
) error {
	if err := e.baseExecutor.ExecCommand(
		ctx,
		"make",
		append([]string{"-f", makefilePath, target}, args...)...,
	); err != nil {
		return fmt.Errorf("run target %q from makefile %q: %w", target, makefilePath, err)
	}

	return nil
}
