package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const agentBookmark = "gemini"

func findExecutable(name string) string {
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

func usage() {
	usageText := `Usage: jj-agent <command> [args...]
Allowed commands:
  log                     - List the mutable subtree
  new [base_rev]          - Create a new change on top of base_rev (defaults to @)
  describe <rev> <msg>    - Update the description of a change
  rebase <src> <dest>     - Rebase <src> onto <dest>
  squash <src> [into_rev] - Squash <src> into <into_rev> (defaults to parent)
  split <rev>             - Split a change into two
  duplicate <rev>         - Duplicate a change`
	fmt.Fprintln(os.Stderr, usageText)
	os.Exit(1)
}

// newJJCmd creates an exec.Cmd for jj, explicitly setting the working directory.
func newJJCmd(args ...string) *exec.Cmd {
	fullArgs := append([]string{"--no-pager", "--color=never"}, args...)
	cmd := exec.Command(findExecutable("jj"), fullArgs...)
	cmd.Env = os.Environ()
	if cwd, err := os.Getwd(); err == nil {
		cmd.Dir = cwd
	}
	return cmd
}

// runJJ executes a jj command and pipes output directly to the agent's stdout/stderr
func runJJ(args ...string) {
	cmd := newJJCmd(args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		os.Exit(1) // jj already printed the error to stderr
	}
}

// validateRevs enforces the sandbox boundaries for the given revisions.
func validateRevs(revs ...string) {
	for _, rev := range revs {
		if strings.TrimSpace(rev) == "" {
			continue
		}

		// 1. Ensure the revision exists and is valid
		checkCmd := newJJCmd("log", "-r", rev, "--no-graph", "-T", "")
		if err := checkCmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: Revision '%s' is invalid or does not exist.\n", rev)
			os.Exit(1)
		}

		// 2. Sandbox Check: (Target) ~ (Bookmark::)
		// If the result contains any commits, the target is out of bounds.
		query := fmt.Sprintf("(%s) ~ (%s::)", rev, agentBookmark)
		sandboxCmd := newJJCmd("log", "-r", query, "--no-graph", "-T", "commit_id")

		out, err := sandboxCmd.CombinedOutput()
		if err != nil {
			// If jj fails to parse the query itself, assume unsafe
			fmt.Fprintf(os.Stderr, "Error validating sandbox boundaries for '%s': %v\n", rev, err)
			os.Exit(1)
		}

		if len(bytes.TrimSpace(out)) > 0 {
			fmt.Fprintf(os.Stderr, "Error: Sandbox violation! Revision '%s' falls outside the '%s' subtree.\n", rev, agentBookmark)
			os.Exit(1)
		}
	}
}

func main() {
	if len(os.Args) < 2 {
		usage()
	}

	command := os.Args[1]
	args := os.Args[2:]

	switch command {
	case "log", "list":
		fmt.Printf("--- Viewing allowed subtree: %s:: ---\n", agentBookmark)
		runJJ("log", "-r", fmt.Sprintf("%s::", agentBookmark))

	case "new":
		target := "@"
		if len(args) > 0 {
			target = args[0]
		}
		validateRevs(target)
		fmt.Printf("Creating new change on top of %s...\n", target)
		runJJ("new", target)

	case "describe":
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "Error: Missing revision or description message.")
			usage()
		}
		target := args[0]
		msg := args[1]
		validateRevs(target)
		fmt.Printf("Updating description for %s...\n", target)
		runJJ("describe", target, "-m", msg)

	case "rebase":
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "Error: Missing source or destination.")
			usage()
		}
		src := args[0]
		dest := args[1]
		validateRevs(src, dest)
		fmt.Printf("Rebasing %s onto %s...\n", src, dest)
		runJJ("rebase", "-s", src, "-d", dest)

	case "squash":
		src := "@"
		if len(args) > 0 {
			src = args[0]
		}
		if len(args) > 1 {
			dest := args[1]
			validateRevs(src, dest)
			fmt.Printf("Squashing %s into %s...\n", src, dest)
			runJJ("squash", "-r", src, "--into", dest)
		} else {
			validateRevs(src)
			fmt.Printf("Squashing %s into its immediate parent...\n", src)
			runJJ("squash", "-r", src)
		}

	case "split":
		if len(args) < 1 {
			usage()
		}
		target := args[0]
		validateRevs(target)
		fmt.Printf("Splitting change %s...\n", target)
		runJJ("split", "-r", target)

	case "duplicate":
		if len(args) < 1 {
			usage()
		}
		target := args[0]
		validateRevs(target)
		fmt.Printf("Duplicating change %s...\n", target)
		runJJ("duplicate", "-r", target)

	default:
		fmt.Fprintf(os.Stderr, "Error: Unknown or unauthorized command '%s'.\n", command)
		usage()
	}
}
