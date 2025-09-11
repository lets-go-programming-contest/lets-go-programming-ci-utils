package sanity

import (
	"errors"
	"fmt"
	"strings"
)

var (
	_ interface{ UserMessage() string } = (*AffectedNoTaskFilesError)(nil)
	_ interface{ UserMessage() string } = (*NoStudentsFoundError)(nil)
	_ interface{ UserMessage() string } = (*MultipleStudentsError)(nil)
	_ interface{ UserMessage() string } = (*NoTasksFoundError)(nil)
	_ interface{ UserMessage() string } = (*MultipleTasksError)(nil)
)

var ErrTasksNotRepresented = errors.New("tasks are not represented in commits")

type affect struct {
	Email string
	File  string
	Hash  string
}

type AffectedNoTaskFilesError struct {
	affects []affect
}

func NewAffectedNoTaskFilesError(affected ...affect) error {
	return &AffectedNoTaskFilesError{
		affects: affected,
	}
}

func (e *AffectedNoTaskFilesError) Error() string {
	return fmt.Sprintf("commits affect %d no task files", len(e.affects))
}

func (e *AffectedNoTaskFilesError) UserMessage() string {
	var builder strings.Builder

	builder.WriteString("The following no tasks files are affected:\n")

	for _, affect := range e.affects {
		builder.WriteString(fmt.Sprintf("\t -file %q affected by %q in %q\n",
			affect.File, affect.Email, affect.Hash),
		)
	}

	return builder.String()
}

type NoStudentsFoundError struct{}

func (*NoStudentsFoundError) Error() string {
	return "no students found"
}

func (*NoStudentsFoundError) UserMessage() string {
	return "Changes for at least one student must be accepted in commits!"
}

type MultipleStudentsError struct {
	names []string
}

func (e *MultipleStudentsError) Error() string {
	return fmt.Sprintf("commits contain changes for %d students", len(e.names))
}

func (e *MultipleStudentsError) UserMessage() string {
	var builder strings.Builder

	builder.WriteString("The following students' tasks were found in the commits:\n")

	for _, name := range e.names {
		builder.WriteString(fmt.Sprintf("\t- %s\n", name))
	}

	builder.WriteString("However, commits should only contain changes for one student!")

	return builder.String()
}

type NoTasksFoundError struct{}

func (*NoTasksFoundError) Error() string {
	return "no tasks found"
}

func (*NoTasksFoundError) UserMessage() string {
	return "Changes for at least one task must be accepted in commits!"
}

type MultipleTasksError struct {
	names []string
}

func NewMultipleTasksError(names ...string) error {
	return &MultipleTasksError{
		names: names,
	}
}

func (e *MultipleTasksError) Error() string {
	return fmt.Sprintf("сommits contain changes for %d tasks", len(e.names))
}

func (e *MultipleTasksError) UserMessage() string {
	var builder strings.Builder

	builder.WriteString("The following tasks were found in the commits:\n")

	for _, name := range e.names {
		builder.WriteString(fmt.Sprintf("\t- %s\n", name))
	}

	builder.WriteString("However, commits should only contain changes for a single task!")

	return builder.String()
}
