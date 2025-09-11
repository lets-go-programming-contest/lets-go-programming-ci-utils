package repository

import (
	"fmt"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/storer"
)

type repository struct {
	repo *git.Repository
}

func NewRepository(
	repo *git.Repository,
) repository {
	return repository{
		repo: repo,
	}
}

func (r repository) getCommit(rev string) (*object.Commit, error) {
	hash, err := r.repo.ResolveRevision(plumbing.Revision(rev))
	if err != nil {
		return nil, fmt.Errorf("resolve rev %q: %w", rev, err)
	}

	commit, err := r.repo.CommitObject(*hash)
	if err != nil {
		return nil, fmt.Errorf("get commit with hash %q: %w", hash.String(), err)
	}

	return commit, nil
}

func (repository) getFileName(change *object.Change) string {
	if change.To.Name != "" {
		return change.To.Name
	}

	return change.From.Name
}

type Changes map[string][]FileChanges

type FileChanges struct {
	Hash   string
	Author string
}

func (r repository) getChanges(baseCommit, targetCommit *object.Commit) (Changes, error) {
	iterator, err := r.repo.Log(&git.LogOptions{
		From: targetCommit.Hash,
	})
	if err != nil {
		return nil, fmt.Errorf("create new diff iterator: %w", err)
	}

	fileChanges := make(map[string][]FileChanges)

	if err := iterator.ForEach(func(c *object.Commit) error {
		if c.Hash == baseCommit.Hash {
			return storer.ErrStop
		}

		parent, err := c.Parent(0)
		if err != nil {
			return nil
		}

		parentTree, err := parent.Tree()
		if err != nil {
			return fmt.Errorf("calculate base tree: %w", err)
		}

		iterTree, err := c.Tree()
		if err != nil {
			return fmt.Errorf("calculate target tree: %w", err)
		}

		diff, err := parentTree.Diff(iterTree)
		if err != nil {
			return fmt.Errorf("calculate diff from %q to %q: %w", baseCommit.Hash, targetCommit.Hash, err)
		}

		for _, change := range diff {
			fileName := r.getFileName(change)
			fileChanges[fileName] = append(fileChanges[fileName], FileChanges{
				Author: c.Author.Email,
				Hash:   c.Hash.String(),
			})
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("interate from %q to %q: %w", baseCommit.Hash, targetCommit.Hash, err)
	}

	return fileChanges, nil
}

func (r repository) GetChanges(
	baseRev string,
	targetRev string,
) (map[string][]FileChanges, error) {
	baseCommit, err := r.getCommit(baseRev)
	if err != nil {
		return nil, err
	}

	targetCommit, err := r.getCommit(targetRev)
	if err != nil {
		return nil, err
	}

	if baseCommit.Hash == targetCommit.Hash {
		baseCommit, err = baseCommit.Parent(0)
		if err != nil {
			return nil, fmt.Errorf("recalculate target: %w", err)
		}
	}

	return r.getChanges(baseCommit, targetCommit)
}
