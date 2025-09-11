package sanity

import (
	"errors"
	"fmt"
)

var (
	ErrTasksNotRepresented = errors.New("tasks are not represented in commits")
)

type AffectedNoTaskFilesErrorAffect struct {
	Email string
	File  string
	Hash  string
}

type AffectedNoTaskFilesError struct {
	Affected []AffectedNoTaskFilesErrorAffect
}

func NewAffectedNoTaskFilesError(affected ...AffectedNoTaskFilesErrorAffect) error {
	return &AffectedNoTaskFilesError{
		Affected: affected,
	}
}

func (e *AffectedNoTaskFilesError) Error() string {
	return fmt.Sprintf("commits affect %d no task files", len(e.Affected))
}

type MultipleStudentsError struct {
	Names []string
}

func NewMultipleStudentsError(names ...string) error {
	return &MultipleStudentsError{
		Names: names,
	}
}

func (e *MultipleStudentsError) Error() string {
	return fmt.Sprintf("Commits contain changes for %d students", len(e.Names))
}

type MultipleTasksError struct {
	Names []string
}

func NewMultipleTasksError(names ...string) error {
	return &MultipleTasksError{
		Names: names,
	}
}

func (e *MultipleTasksError) Error() string {
	return fmt.Sprintf("Commits contain changes for %d tasks", len(e.Names))
}
