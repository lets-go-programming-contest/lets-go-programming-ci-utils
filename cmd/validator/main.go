package main

import (
	"fmt"
	"os"

	"github.com/lets-go-programming-contest/lets-go-programming-ci-utils/cmd/validator/internal/cli/sanity"
	"github.com/spf13/cobra"
)

const (
	cliConfigFlag      = "config"
	cliBaseRevFlag     = "base"
	cliTargetRevFlag   = "target"
	cliStudentNameFlag = "student"
	cliTaskNameFlag    = "task"
	cliOutputNameFlag  = "output"
)

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	rootCmd := &cobra.Command{
		Use:           "lgp_validator",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	// rootCmd.AddCommand(module.NewModuleCmd(
	// 	initConfigFlag,
	// 	initStudentFlag,
	// 	initTaskFlag,
	// ))
	rootCmd.AddCommand(sanity.NewSanityCmd(
		initBaseRevFlag,
		initTargetRevFlag,
	))

	initConfigFlag(rootCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
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
