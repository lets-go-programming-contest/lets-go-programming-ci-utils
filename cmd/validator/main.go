package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/lets-go-programming-contest/lets-go-programming-ci-utils/cmd/validator/internal/cli/module"
	"github.com/lets-go-programming-contest/lets-go-programming-ci-utils/cmd/validator/internal/cli/sanity"
	"github.com/spf13/cobra"
)

const (
	cliConfigFlag      = "config"
	cliBaseRevFlag     = "base"
	cliTargetRevFlag   = "target"
	cliStudentNameFlag = "student"
	cliTaskNameFlag    = "task"
)

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	//nolint:exhaustruct // Set defaults values for another fields.
	rootCmd := &cobra.Command{
		Use:           "lgp_validator",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	rootCmd.AddCommand(module.NewModuleCmd(
		initConfigFlag,
		initStudentFlag,
		initTaskFlag,
	))
	rootCmd.AddCommand(sanity.NewSanityCmd(
		initBaseRevFlag,
		initTargetRevFlag,
	))

	initConfigFlag(rootCmd)

	var errUserMessage interface {
		UserMessage() string
	}

	err := rootCmd.Execute()

	switch {
	case err == nil:
		fmt.Println("The execution finished without errors.\tOK.")
	case errors.As(err, &errUserMessage):
		fmt.Fprintln(os.Stderr, errUserMessage.UserMessage())
		fmt.Fprintf(os.Stderr, "error trace: %s\n", err.Error())
		os.Exit(1)
	default:
		fmt.Fprintln(os.Stderr, "The execution finished with internal errors.")
		fmt.Fprintln(os.Stderr, "Please refer to the practices for details.")

		panic(err)
	}
}

func initConfigFlag(cmd *cobra.Command) {
	cmd.PersistentFlags().String(cliConfigFlag, "config.yaml", "Path to task configuration")
	panicIfErr(cmd.MarkPersistentFlagFilename(cliConfigFlag))
}

func initStudentFlag(cmd *cobra.Command) {
	cmd.PersistentFlags().String(cliStudentNameFlag, "", "Student name")
	panicIfErr(cmd.MarkPersistentFlagRequired(cliStudentNameFlag))
}

func initTaskFlag(cmd *cobra.Command) {
	cmd.PersistentFlags().String(cliTaskNameFlag, "", "Task name")
	panicIfErr(cmd.MarkPersistentFlagRequired(cliTaskNameFlag))
}

func initBaseRevFlag(cmd *cobra.Command) {
	cmd.PersistentFlags().String(cliBaseRevFlag, "origin/main", "Base rev for diff calculation")
}

func initTargetRevFlag(cmd *cobra.Command) {
	cmd.PersistentFlags().String(cliTargetRevFlag, "HEAD", "Target rev for diff calculation")
}
