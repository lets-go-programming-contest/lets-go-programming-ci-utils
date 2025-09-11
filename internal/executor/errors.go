package executor

import (
	"fmt"
	"strings"
)

type ExecError struct {
	CommandName string
	CommandArgs []string
}

func NewExecError(commandName string, commandArgs []string) error {
	return &ExecError{
		CommandName: commandName,
		CommandArgs: commandArgs,
	}
}

func (e *ExecError) Error() string {
	return fmt.Sprintf(
		"run command %q execution finished with error",
		strings.Join(append([]string{e.CommandName}, e.CommandArgs...), " "),
	)
}
