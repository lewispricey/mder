package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	viewMode bool
	editMode bool
)

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "mded [flags] <file>",
		Short:         "Markdown editor",
		SilenceUsage:  true,
		SilenceErrors: true,
		Args:          cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			path := args[0]

			if viewMode && editMode {
				return fmt.Errorf("cannot set both --view and --edit")
			}

			mode := "edit"
			if viewMode {
				mode = "view"
			}

			if _, err := os.Stat(path); err != nil {
				return fmt.Errorf("file %q does not exist", path)
			}

			fmt.Printf("Opening: %s in %s mode\n", path, mode)
			return nil
		},
	}

	cmd.Flags().BoolVar(&viewMode, "view", false, "open in view-only mode")
	cmd.Flags().BoolVar(&editMode, "edit", false, "open in edit mode")
	return cmd
}

func Execute() {
	if err := NewRootCmd().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
