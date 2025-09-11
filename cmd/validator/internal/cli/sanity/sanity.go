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

const cliToEnvFileFlag = "to-env-file"

const (
	envStudentNameKey = "STUDENT_NAME"
	envTaskNameKey    = "TASK_NAME"
	envFilePerms      = 0o644
)

var errCreateService = errors.New("create service")

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}

	return v
}

func getToEnvFileFlag(flagset *pflag.FlagSet) (string, bool) {
	if flagset.Changed(cliToEnvFileFlag) {
		return must(flagset.GetString(cliToEnvFileFlag)), true
	}

	return "", false
}

func getBaseRev(flagSet *pflag.FlagSet) string {
	return must(flagSet.GetString("base"))
}

func getTargetRev(flagSet *pflag.FlagSet) string {
	return must(flagSet.GetString("target"))
}

func NewSanityCmd(opts ...func(cmd *cobra.Command)) *cobra.Command {
	//nolint:exhaustruct // Set defaults values for another fields.
	sanityCmd := &cobra.Command{
		Use:   "sanity",
		Short: "Set of commands for validating student tasks",
	}

	sanityCmd.AddCommand(NewSanityFilesCmd())
	sanityCmd.AddCommand(NewSanityStudentsCmd(initToNamedEnvFlag))
	sanityCmd.AddCommand(NewSanityTasksCmd(initToNamedEnvFlag))

	for _, opt := range opts {
		if opt != nil {
			opt(sanityCmd)
		}
	}

	return sanityCmd
}

func initToNamedEnvFlag(cmd *cobra.Command) {
	cmd.Flags().String(cliToEnvFileFlag, "", "The name of the environment variable file in which the value will be stored")
}

func NewSanityFilesCmd() *cobra.Command {
	//nolint:exhaustruct // Set defaults values for another fields.
	return &cobra.Command{
		Use:   "files",
		Short: "Checks changed files",
		Long: `Checks that there have been no changes to the root 
			files of the repository in the current changes.`,
		RunE: runSanityFiles,
	}
}

func NewSanityStudentsCmd(opts ...func(cmd *cobra.Command)) *cobra.Command {
	//nolint:exhaustruct // Set defaults values for another fields.
	studentsCmd := &cobra.Command{
		Use:   "students",
		Short: "Checks students",
		Long:  `Checks that the current changes have been modified for only one student.`,
		RunE:  runSanityStudents,
	}

	for _, opt := range opts {
		if opt != nil {
			opt(studentsCmd)
		}
	}

	return studentsCmd
}

func NewSanityTasksCmd(opts ...func(cmd *cobra.Command)) *cobra.Command {
	//nolint:exhaustruct // Set defaults values for another fields.
	tasksCmd := &cobra.Command{
		Use:   "tasks",
		Short: "Checks tasks",
		Long:  `Checks that the current changes have been modified for only one task.`,
		RunE:  runSanityTasks,
	}

	for _, opt := range opts {
		if opt != nil {
			opt(tasksCmd)
		}
	}

	return tasksCmd
}

func runSanityFiles(cmd *cobra.Command, _ []string) error {
	srv, err := sanity.NewService(
		envs.GetPwd(),
		getBaseRev(cmd.Flags()),
		getTargetRev(cmd.Flags()),
	)
	if err != nil {
		return fmt.Errorf("%w: %w", errCreateService, err)
	}

	if err := srv.RunSanityTaskFiles(context.Background()); err != nil {
		return fmt.Errorf("sanity task files: %w", err)
	}

	fmt.Println("Files from commits are accepted for automatic review!")

	return nil
}

func appendValueIntoEnvFile(filename string, key string, value string) error {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_RDWR, envFilePerms)
	if err != nil {
		return fmt.Errorf("open env file %q: %w", filename, err)
	}
	defer file.Close()

	if _, err := fmt.Fprintf(file, "%s=%s\n", key, value); err != nil {
		return fmt.Errorf("write to env file: %q: %w", filename, err)
	}

	fmt.Printf("The result was recorded in the environment file %q as %q.\n", filename, key)

	return nil
}

func runSanityStudents(cmd *cobra.Command, _ []string) error {
	srv, err := sanity.NewService(
		envs.GetPwd(),
		getBaseRev(cmd.Flags()),
		getTargetRev(cmd.Flags()),
	)
	if err != nil {
		return fmt.Errorf("%w: %w", errCreateService, err)
	}

	studentName, err := srv.RunSanityStudents(context.Background())
	if err != nil {
		return fmt.Errorf("sanity students: %w", err)
	}

	fmt.Printf("Tasks by %s are accepted for automatic review!\n", studentName)

	if envFileName, ok := getToEnvFileFlag(cmd.Flags()); ok {
		if err := appendValueIntoEnvFile(envFileName, envStudentNameKey, studentName); err != nil {
			return err
		}
	}

	return nil
}

func runSanityTasks(cmd *cobra.Command, _ []string) error {
	srv, err := sanity.NewService(
		envs.GetPwd(),
		getBaseRev(cmd.Flags()),
		getTargetRev(cmd.Flags()),
	)
	if err != nil {
		return fmt.Errorf("%w: %w", errCreateService, err)
	}

	taskName, err := srv.RunSanityTasks(context.Background())
	if err != nil {
		return fmt.Errorf("sanity tasks: %w", err)
	}

	fmt.Printf("Task %s accepted for automatic review!\n", taskName)

	if envFileName, ok := getToEnvFileFlag(cmd.Flags()); ok {
		if err := appendValueIntoEnvFile(envFileName, envTaskNameKey, taskName); err != nil {
			return err
		}
	}

	return nil
}
