package backend

import (
	"fmt"
	"os"

	"github.com/go-git/go-git/v6"
)

func IsGitRepo(localPath string) bool {
	_, err := git.PlainOpen(localPath)
	return err == nil
}

func CloneGitRepo(remoteURL, localPath string) error {
	_, err := git.PlainClone(localPath, &git.CloneOptions{
		URL:      remoteURL,
		Progress: os.Stdout,
	})

	return err
}

func PullGitRepo(localPath string) error {
	repo, err := git.PlainOpen(localPath)
	if err != nil {
		return err
	}

	w, err := repo.Worktree()
	if err != nil {
		return err
	}

	err = w.Pull(&git.PullOptions{RemoteName: "origin"})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		return err
	}

	return nil
}

func PushGitRepo(localPath string) error {
	repo, err := git.PlainOpen(localPath)
	if err != nil {
		return err
	}

	err = repo.Push(&git.PushOptions{})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		return err
	}

	return nil
}

func AddAndCommitGitRepo(localPath, message string) error {
	repo, err := git.PlainOpen(localPath)

	if err != nil {
		return err
	}

	if repo == nil {
		return fmt.Errorf("repository not found at %s", localPath)
	}

	w, err := repo.Worktree()

	if err != nil {
		return err
	}

	if w == nil {
		return fmt.Errorf("worktree not found at %s", localPath)
	}

	status, err := w.Status()

	if err != nil {
		return err
	}

	if status.IsClean() {
		return nil
	}

	_, err = w.Add(".")
	if err != nil {
		return err
	}

	_, err = w.Commit(message, &git.CommitOptions{})
	if err != nil {
		return fmt.Errorf("could not commit changes: %w", err)
	}

	return nil
}

func IsCleanGitRepo(localPath string) (bool, error) {
	repo, err := git.PlainOpen(localPath)
	if err != nil {
		return false, err
	}

	w, err := repo.Worktree()
	if err != nil {
		return false, err
	}

	status, err := w.Status()
	if err != nil {
		return false, err
	}

	return status.IsClean(), nil
}

func InitGitRepo(localPath string) error {
	_, err := git.PlainInit(localPath, false)
	return err
}
