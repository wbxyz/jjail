package jjutil

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const AgentBookmark = "agent"

func FindExecutable(name string) string {
	if path, err := exec.LookPath(name); err == nil {
		return path
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return name
	}

	pathEnv := os.Getenv("PATH")
	for _, dir := range strings.Split(pathEnv, string(os.PathListSeparator)) {
		if strings.HasPrefix(dir, "~/") {
			dir = filepath.Join(home, dir[2:])
		} else if dir == "~" {
			dir = home
		}

		fullPath := filepath.Join(dir, name)
		if stat, err := os.Stat(fullPath); err == nil && !stat.IsDir() && stat.Mode()&0111 != 0 {
			return fullPath
		}
	}

	return name
}

// NewJJCmd creates an exec.Cmd for jj, explicitly setting the working directory.
func NewJJCmd(args ...string) *exec.Cmd {
	fullArgs := append([]string{"--no-pager", "--color=never"}, args...)
	cmd := exec.Command(FindExecutable("jj"), fullArgs...)
	cmd.Env = os.Environ()
	if cwd, err := os.Getwd(); err == nil {
		cmd.Dir = cwd
	}
	return cmd
}

// RunJJ executes a jj command and pipes output directly to the agent's stdout/stderr
func RunJJ(args ...string) {
	cmd := NewJJCmd(args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		os.Exit(1) // jj already printed the error to stderr
	}
}
