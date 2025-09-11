package sanity

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/lets-go-programming-contest/lets-go-programming-ci-utils/internal/envs"
	"github.com/lets-go-programming-contest/lets-go-programming-ci-utils/internal/service/sanity"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var (
	errInternalError   = errors.New("run sanity")
	errValidationError = errors.New("validate changes")
)

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}

	return v
}

func getBaseRev(flagSet *pflag.FlagSet) string {
	return must(flagSet.GetString("base"))
}

func getTargetRev(flagSet *pflag.FlagSet) string {
	return must(flagSet.GetString("target"))
}

func NewSanityCmd(opts ...func(cmd *cobra.Command)) *cobra.Command {
	sanityCmd := &cobra.Command{
		Use: "sanity",
	}

	sanityCmd.AddCommand(NewSanityFilesCmd())
	sanityCmd.AddCommand(NewSanityStudentsCmd())
	sanityCmd.AddCommand(NewSanityTasksCmd())

	for _, opt := range opts {
		if opt != nil {
			opt(sanityCmd)
		}
	}

	return sanityCmd
}

func NewSanityFilesCmd() *cobra.Command {
	return &cobra.Command{
		Use:  "files",
		RunE: runSanityFiles,
	}
}

func NewSanityStudentsCmd() *cobra.Command {
	return &cobra.Command{
		Use:  "students",
		RunE: runSanityStudents,
	}
}

func NewSanityTasksCmd() *cobra.Command {
	taskStudentsCmd := &cobra.Command{
		Use:  "tasks",
		RunE: runSanityTasks,
	}

	return taskStudentsCmd
}

func runSanityFiles(cmd *cobra.Command, args []string) error {
	srv, err := sanity.NewService(
		envs.GetPwd(),
		getBaseRev(cmd.Flags()),
		getTargetRev(cmd.Flags()),
	)
	if err != nil {
		return errors.Join(errInternalError, err)
	}

	if err := srv.RunSanityTaskFiles(context.Background()); err != nil {
		var errAffectedNoTaskFiles *sanity.AffectedNoTaskFilesError
		if errors.As(err, &errAffectedNoTaskFiles) {
			fmt.Fprintln(os.Stderr, "The following no tasks files are affected:")
			for _, affect := range errAffectedNoTaskFiles.Affected {
				fmt.Fprintf(os.Stderr, "\t -file %q affected by %q in %q\n", affect.File, affect.Email, affect.Hash)
			}

			return errors.Join(errValidationError, err)
		}

		return errors.Join(errInternalError, err)
	}

	fmt.Println("Files from commits are accepted for automatic review!")
	return nil
}

func runSanityStudents(cmd *cobra.Command, args []string) error {
	srv, err := sanity.NewService(
		envs.GetPwd(),
		getBaseRev(cmd.Flags()),
		getTargetRev(cmd.Flags()),
	)
	if err != nil {
		return errors.Join(errInternalError, err)
	}

	studentName, err := srv.RunSanityStudents(context.Background())
	if err != nil {
		if errors.Is(err, sanity.ErrTasksNotRepresented) {
			fmt.Fprintln(os.Stderr, "Changes for at least one student must be accepted in commits!")
			return errors.Join(errValidationError, err)
		}

		var errMultipleStudents *sanity.MultipleStudentsError
		if errors.As(err, &errMultipleStudents) {
			fmt.Fprintln(os.Stderr, "The following students' tasks were found in the commits:")
			for _, name := range errMultipleStudents.Names {
				fmt.Fprintf(os.Stderr, "\t- %s\n", name)
			}
			fmt.Fprintln(os.Stderr, "However, commits should only contain changes for a one student!")

			return errors.Join(errValidationError, err)
		}

		return errors.Join(errInternalError, err)
	}

	fmt.Printf("Tasks by %s are accepted for automatic review!\n", studentName)

	return nil
}

func runSanityTasks(cmd *cobra.Command, args []string) error {
	srv, err := sanity.NewService(
		envs.GetPwd(),
		getBaseRev(cmd.Flags()),
		getTargetRev(cmd.Flags()),
	)
	if err != nil {
		return errors.Join(errInternalError, err)
	}

	taskName, err := srv.RunSanityTasks(context.Background())
	if err != nil {
		if errors.Is(err, sanity.ErrTasksNotRepresented) {
			fmt.Fprintln(os.Stderr, "Changes for at least one student must be accepted in commits!")
			return errors.Join(errValidationError, err)
		}

		var errMultipleTasks *sanity.MultipleTasksError
		if errors.As(err, &errMultipleTasks) {
			fmt.Fprintln(os.Stderr, "The following tasks were found in the commits:")
			for _, name := range errMultipleTasks.Names {
				fmt.Fprintf(os.Stderr, "\t- %s\n", name)
			}
			fmt.Fprintln(os.Stderr, "However, commits should only contain changes for a single task!")

			return errors.Join(errValidationError, err)
		}

		return errors.Join(errInternalError, err)
	}

	fmt.Printf("Task %s accepted for automatic review!\n", taskName)

	return nil
}
