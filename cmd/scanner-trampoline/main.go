package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbletea"
	"github.com/dmorn/scanner-trampoline/cli/setup"
)

func main() {
	if err := tea.NewProgram(setup.New()).Start(); err != nil {
		fmt.Printf("could not start program: %s\n", err)
		os.Exit(1)
	}
}
