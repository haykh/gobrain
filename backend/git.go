package backend

import (
	"fmt"
	"os"
	"os/exec"
)

func IsGitRepo(localPath string) bool {
	gitDir := fmt.Sprintf("%s/.git", localPath)
	info, err := os.Stat(gitDir)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func CloneGitRepo(remoteURL, localPath string) error {
	cmd := exec.Command("git", "clone", remoteURL, localPath)
	return cmd.Run()
}

func FetchGitRepo(localPath string) error {
	cmd := exec.Command("git", "-C", localPath, "fetch")
	return cmd.Run()
}

func PushGitRepo(localPath string) error {
	cmd := exec.Command("git", "-C", localPath, "push")
	return cmd.Run()
}

func AddAndCommitGitRepo(localPath, message string) error {
	cmdAdd := exec.Command("git", "-C", localPath, "add", ".")
	if err := cmdAdd.Run(); err != nil {
		return fmt.Errorf("could not add changes: %w", err)
	}

	cmdCommit := exec.Command("git", "-C", localPath, "commit", "-m", message)
	if err := cmdCommit.Run(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			if exitError.ExitCode() == 1 {
				return nil
			}
		}
		return fmt.Errorf("could not commit changes: %w", err)
	}

	return nil
}

func IsCleanGitRepo(localPath string) (bool, error) {
	cmd := exec.Command("git", "-C", localPath, "status", "--porcelain")
	output, err := cmd.Output()
	if err != nil {
		return false, err
	}

	return len(output) == 0, nil
}

func InitGitRepo(localPath string) error {
	cmd := exec.Command("git", "init", localPath)
	return cmd.Run()
}

func CheckGitAhead(localPath string) (bool, error) {
	cmd := exec.Command("git", "-C", localPath, "status", "-s")
	output, err := cmd.Output()
	if err != nil {
		return false, err
	}
	return len(output) > 0, nil
}

func CheckGitBehind(localPath string) (bool, error) {
	if err := FetchGitRepo(localPath); err != nil {
		return false, err
	}
	cmd := exec.Command("git", "-C", localPath, "log", "--oneline", "HEAD..origin/master")
	output, err := cmd.Output()
	if err != nil {
		return false, err
	}
	return len(output) > 0, nil
}
