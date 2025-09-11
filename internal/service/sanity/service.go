package sanity

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"

	"github.com/go-git/go-git/v5"
	"github.com/lets-go-programming-contest/lets-go-programming-ci-utils/internal/envs"
	"github.com/lets-go-programming-contest/lets-go-programming-ci-utils/internal/repository"
)

const (
	studentPattern = "[a-z0-9]+\\.[a-z0-9]+"
	taskPattern    = "task-[0-9-]+"
)

var (
	taskFilesPattern = fmt.Sprintf("^%s/%s/.+$", studentPattern, taskPattern)

	studentRegexp   = regexp.MustCompile(fmt.Sprintf("^%s", studentPattern))
	taskRegexp      = regexp.MustCompile(taskPattern)
	taskFilesRegexp = regexp.MustCompile(taskFilesPattern)
)

type repo interface {
	GetChanges(baseRev string, targetRev string) (map[string][]repository.FileChanges, error)
}

type service struct {
	repo      repo
	baseRev   string
	targetRev string
}

func NewServiceWithRepo(repo repo, baseRev, targetRev string) (service, error) {
	return service{
		repo:      repo,
		baseRev:   baseRev,
		targetRev: targetRev,
	}, nil
}

func NewService(repoPath, baseRev, tragetRev string) (service, error) {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return service{}, fmt.Errorf("open repository: %w", err)
	}

	return NewServiceWithRepo(
		repository.NewRepository(repo),
		baseRev,
		tragetRev,
	)
}

func getMaintainers(maintainersFile string) ([]string, error) {
	file, err := os.Open(maintainersFile)
	if err != nil {
		return nil, fmt.Errorf("open maintainers file: %w", err)
	}

	var maintainers []string
	if err := json.NewDecoder(file).Decode(&maintainers); err != nil {
		return nil, fmt.Errorf("unmarshal maintainers: %w", err)
	}

	return maintainers, nil
}

func extractFilesFromChanges(
	changes repository.Changes,
	isTaskFiles bool,
) repository.Changes {
	extracted := make(repository.Changes)
	for file, fileChanges := range changes {
		if taskFilesRegexp.MatchString(file) == isTaskFiles {
			extracted[file] = fileChanges
		}
	}

	return extracted
}

func (s service) RunSanityTaskFiles(ctx context.Context) error {
	changes, err := s.repo.GetChanges(s.baseRev, s.targetRev)
	if err != nil {
		return fmt.Errorf("get changes from %q to %q: %w", s.baseRev, s.targetRev, err)
	}

	commonDir, err := envs.GetCommonDirFromEnv()
	if err != nil {
		return fmt.Errorf("get dir with common files: %w", err)
	}

	maintainers, err := getMaintainers(filepath.Join(commonDir, "MAINTAINERS"))
	if err != nil {
		return fmt.Errorf("get maintainers from common files: %w", err)
	}

	noTaskFiles := extractFilesFromChanges(changes, false)
	affected := make([]AffectedNoTaskFilesErrorAffect, 0, len(noTaskFiles))
	for file, fileChanges := range noTaskFiles {
		for _, change := range fileChanges {
			if !slices.Contains(maintainers, change.Author) {
				affected = append(affected, AffectedNoTaskFilesErrorAffect{
					File:  file,
					Email: change.Author,
					Hash:  change.Hash,
				})
			}
		}
	}

	if len(affected) != 0 {
		return NewAffectedNoTaskFilesError(affected...)
	}

	return nil
}

func extractStudentsFromChanges(changes repository.Changes) []string {
	studentsMap := make(map[string]struct{})

	for file := range extractFilesFromChanges(changes, true) {
		if filepath.Dir(file) != "." {
			if matches := studentRegexp.FindStringSubmatch(file); matches != nil {
				studentsMap[matches[0]] = struct{}{}
			}
		}
	}

	students := make([]string, 0, len(studentsMap))

	for student := range studentsMap {
		students = append(students, student)
	}

	return students
}

func (s service) RunSanityStudents(ctx context.Context) (string, error) {
	changes, err := s.repo.GetChanges(s.baseRev, s.targetRev)
	if err != nil {
		return "", fmt.Errorf("get changes from %q to %q: %w", s.baseRev, s.targetRev, err)
	}

	students := extractStudentsFromChanges(changes)

	switch len(students) {
	case 0:
		return "", ErrTasksNotRepresented
	case 1:
		return students[0], nil
	default:
		return "", NewMultipleStudentsError(students...)
	}
}

func extractTasksFromChanges(changes repository.Changes) []string {
	taskFiles := extractFilesFromChanges(changes, true)
	if len(taskFiles) == 0 {
		return nil
	}

	tasksMap := make(map[string]struct{})

	for file := range taskFiles {
		if mathes := taskRegexp.FindStringSubmatch(file); mathes != nil {
			tasksMap[mathes[0]] = struct{}{}
		}
	}

	tasks := make([]string, 0, len(tasksMap))
	for task := range tasksMap {
		tasks = append(tasks, task)
	}

	return tasks
}

func (s service) RunSanityTasks(ctx context.Context) (string, error) {
	changes, err := s.repo.GetChanges(s.baseRev, s.targetRev)
	if err != nil {
		return "", fmt.Errorf("get changes from %q to %q: %w", s.baseRev, s.targetRev, err)
	}

	tasks := extractTasksFromChanges(changes)

	switch len(tasks) {
	case 0:
		return "", ErrTasksNotRepresented
	case 1:
		return tasks[0], nil
	default:
		return "", NewMultipleTasksError(tasks...)
	}
}
