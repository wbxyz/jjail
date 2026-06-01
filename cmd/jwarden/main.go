package main

import (
	"fmt"
	"os"

	"github.com/wbxyz/jjail/internal/jjutil"
)

func usage() {
	usageText := `Usage: jwarden <revision>
Moves the 'agent' bookmark to the specified revision.`
	fmt.Fprintln(os.Stderr, usageText)
	os.Exit(1)
}

func main() {
	if len(os.Args) < 2 {
		usage()
	}

	target := os.Args[1]

	fmt.Printf("Moving '%s' bookmark to %s...\n", jjutil.AgentBookmark, target)
	jjutil.RunJJ("bookmark", "move", jjutil.AgentBookmark, "--allow-backwards", "--to", target)
}
