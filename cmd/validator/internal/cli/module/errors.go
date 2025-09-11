package module

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	baseErrors "github.com/lets-go-programming-contest/lets-go-programming-ci-utils/internal/errors"
	"github.com/lets-go-programming-contest/lets-go-programming-ci-utils/internal/module"
)

func processErr(
	student string,
	task string,
	err error,
) error {
	var (
		errBuildModuleTarget         *module.BuildTargetError
		errRunMakeTargetFromMakefile *module.RunMakeTargetError
	)

	switch {
	case err == nil:
		return nil
	case errors.Is(err, module.ErrPathNotModule):
		fmt.Fprintf(os.Stderr, "Dir %q is not Go module.",
			filepath.Join(student, task))

		return baseErrors.ValidationError(err)

	case errors.Is(err, module.ErrUpdateModuleDeps):
		fmt.Fprintln(os.Stderr, "Failed to update module deps.")

		return baseErrors.ValidationError(err)

	case errors.Is(err, module.ErrNothingToBuild),
		errors.Is(err, module.ErrCmdDirNotFound):
		fmt.Fprintf(os.Stderr, "Targets for build not found in %q.\n",
			filepath.Join(student, task, "cmd"))

		return baseErrors.ValidationError(err)

	case errors.As(err, &errBuildModuleTarget):
		fmt.Fprintf(os.Stderr, "Failed to build target %q in %q.\n",
			errBuildModuleTarget.TargetName, filepath.Join(student, task))

		return baseErrors.ValidationError(err)

	case errors.Is(err, module.ErrLintModule):
		fmt.Fprintf(os.Stderr, "Found linters issues in %q.\n",
			filepath.Join(student, task))

		return baseErrors.ValidationError(err)
	case errors.Is(err, module.ErrTestModule):
		fmt.Fprintf(os.Stderr, "Tests failed to module %q.\n",
			filepath.Join(student, task))

		return baseErrors.ValidationError(err)

	case errors.Is(err, module.ErrMakefilePathIsDir),
		errors.Is(err, module.ErrGetMakefile):

		var fromStudentErr errMakeTargetFromStudentFiles
		if errors.As(err, &fromStudentErr) {
			fmt.Fprintf(os.Stderr, `The target could not be completedso the 
				Makefile from %q could not be accessed for execution.\n`,
				filepath.Join(student, task))

			return baseErrors.ValidationError(err)
		}

		return baseErrors.InternalError(err)

	case errors.As(err, &errRunMakeTargetFromMakefile):
		fmt.Fprintf(os.Stderr, "Make target %q from %q finished with a non-zero status.\n",
			errRunMakeTargetFromMakefile.Target,
			errRunMakeTargetFromMakefile.MakefilePath,
		)

		return baseErrors.ValidationError(err)
	default:
		return baseErrors.InternalError(err)
	}
}
