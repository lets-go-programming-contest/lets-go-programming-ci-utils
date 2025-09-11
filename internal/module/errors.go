package module

import (
	"errors"
	"fmt"
)

var (
	errRelativePath     = errors.New("path must be absolute")
	errPrepareOutputDir = errors.New("prepare output dir")
	errCalculateAbsPath = errors.New("calculate abs path")

	ErrCmdDirNotFound   = errors.New("module targets dir not found")
	ErrPathNotModule    = errors.New("module path not contains go module")
	ErrNothingToBuild   = errors.New("build targets not found in module")
	ErrUpdateModuleDeps = errors.New("update module deps from go.mod file")

	ErrGetMakefile            = errors.New("makefile not found or cannot be opened")
	ErrMakefilePathIsDir      = errors.New("makefile path is directory")
	ErrExectureMakefileTarget = errors.New("makefile target execution failed")
	ErrBuildModuleTargets     = errors.New("build module targets")
	ErrLintModule             = errors.New("lint module")
	ErrTestModule             = errors.New("test module")
)

type BuildTargetError struct {
	TargetName string
}

func newBuildTargetError(targetName string) error {
	return &BuildTargetError{
		TargetName: targetName,
	}
}

func (e *BuildTargetError) Error() string {
	return fmt.Sprintf("build module target %q", e.TargetName)
}

type LintError struct {
	ConfigPath string
}

func newLintError(configPath string) error {
	return &LintError{
		ConfigPath: configPath,
	}
}

func (e *LintError) Error() string {
	return fmt.Sprintf("lint module with config %q", e.ConfigPath)
}

type RunMakeTargetError struct {
	MakefilePath string
	Target       string
}

func newRunMakeTargetError(makefilePath string, target string) error {
	return &RunMakeTargetError{
		MakefilePath: makefilePath,
		Target:       target,
	}
}

func (e *RunMakeTargetError) Error() string {
	return fmt.Sprintf("run target %q from %q", e.Target, e.MakefilePath)
}
