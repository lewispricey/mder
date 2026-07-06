package main

import (
	"fmt"
	"os"

	"charm.land/bubbletea/v2"
	"github.com/yourname/mded/internal/cmd"
	"github.com/yourname/mded/internal/model"
)

func main() {
	mode, path, err := cmd.ParseCLI()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	m := model.New(mode, path)
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
