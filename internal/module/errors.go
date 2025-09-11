//revive:disable:max-public-structs // Errors dir with error wrapper.
package module

import (
	"errors"
	"fmt"
	"path/filepath"
)

var (
	errCmdDirNotFound    = errors.New("module targets dir not found")
	errNothingToBuild    = errors.New("build targets not found in module")
	errMakefilePathIsDir = errors.New("makefile path is directory")
)

var (
	_ interface{ UserMessage() string } = (*UpdateModuleDepsError)(nil)
	_ interface{ UserMessage() string } = (*TargetsNotFoundError)(nil)
	_ interface{ UserMessage() string } = (*PathNoGoModuleError)(nil)
	_ interface{ UserMessage() string } = (*BuildTargetError)(nil)
	_ interface{ UserMessage() string } = (*LintModuleError)(nil)
	_ interface{ UserMessage() string } = (*TestModuleError)(nil)
	_ interface{ UserMessage() string } = (*AccessMakefileError)(nil)
	_ interface{ UserMessage() string } = (*MakeTargetError)(nil)
)

type UpdateModuleDepsError struct {
	modulePath string
	reason     error
}

func (e *UpdateModuleDepsError) Error() string {
	return fmt.Sprintf(
		"update module deps from %q: %s",
		filepath.Join(e.modulePath, "go.mod"), e.reason.Error(),
	)
}

func (e *UpdateModuleDepsError) Unwrap() error {
	return e.reason
}

func (e *UpdateModuleDepsError) UserMessage() string {
	return fmt.Sprintf(
		"Failed to update module deps from %q.",
		filepath.Join(e.modulePath, "go.mod"),
	)
}

type TargetsNotFoundError struct {
	modulePath string
	reason     error
}

func (e *TargetsNotFoundError) Error() string {
	return fmt.Sprintf("build targets not found in %q: %s", e.modulePath, e.reason.Error())
}

func (e *TargetsNotFoundError) Unwrap() error {
	return e.reason
}

func (e *TargetsNotFoundError) UserMessage() string {
	return fmt.Sprintf("Targets for build not found in %q.", e.modulePath)
}

type PathNoGoModuleError struct {
	modulePath string
}

func (e *PathNoGoModuleError) Error() string {
	return fmt.Sprintf("path %q not contains go module", e.modulePath)
}

func (e *PathNoGoModuleError) UserMessage() string {
	return fmt.Sprintf("Path %q not contains go module.", e.modulePath)
}

type BuildTargetError struct {
	modulePath string
	targetName string
	reason     error
}

func (e *BuildTargetError) Error() string {
	return fmt.Sprintf("build module target %q: %s", e.targetName, e.reason.Error())
}

func (e *BuildTargetError) Unwrap() error {
	return e.reason
}

func (e *BuildTargetError) UserMessage() string {
	return fmt.Sprintf("Failed to build target %q in %q.", e.targetName, e.modulePath)
}

type LintModuleError struct {
	modulePath string
	configFile string
	reason     error
}

func (e *LintModuleError) Error() string {
	return fmt.Sprintf("lint module %q with config %q: %s", e.modulePath, e.configFile, e.reason.Error())
}

func (e *LintModuleError) Unwrap() error {
	return e.reason
}

func (e *LintModuleError) UserMessage() string {
	return fmt.Sprintf("Failed to lint module %q with config %q.", e.modulePath, e.configFile)
}

type TestModuleError struct {
	modulePath string
	reason     error
}

func (e *TestModuleError) Error() string {
	return fmt.Sprintf("test module %q: %s", e.modulePath, e.reason.Error())
}

func (e *TestModuleError) Unwrap() error {
	return e.reason
}

func (e *TestModuleError) UserMessage() string {
	return fmt.Sprintf("Failed to test module %q.", e.modulePath)
}

type AccessMakefileError struct {
	makefile string
	reason   error
}

func (e *AccessMakefileError) Error() string {
	return fmt.Sprintf(
		"access file %q: %s", e.makefile, e.reason.Error(),
	)
}

func (e *AccessMakefileError) Unwrap() error {
	return e.reason
}

func (e *AccessMakefileError) UserMessage() string {
	return fmt.Sprintf("Failed to acess makefile %q.", e.makefile)
}

type MakeTargetError struct {
	modulePath string
	makefile   string
	target     string
	reason     error
}

func (e *MakeTargetError) Error() string {
	return fmt.Sprintf(
		"run target %q from file %q for module %q: %s",
		e.target, e.makefile, e.modulePath, e.reason.Error(),
	)
}

func (e *MakeTargetError) Unwrap() error {
	return e.reason
}

func (e *MakeTargetError) UserMessage() string {
	return fmt.Sprintf(
		"Make target %q from %q for module %q finished with a non-zero status.",
		e.target, e.makefile, e.modulePath,
	)
}
