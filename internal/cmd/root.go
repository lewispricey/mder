package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yourname/mded/internal/model"
)

var (
	parsedMode model.Mode
	parsedPath string
)

func NewRootCmd() *cobra.Command {
	var viewMode, editMode bool
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

			parsedMode = model.EditMode
			if viewMode {
				parsedMode = model.ViewMode
			}
			parsedPath = path

			if _, err := os.Stat(path); err != nil {
				return fmt.Errorf("file %q does not exist", path)
			}
			return nil
		},
	}

	cmd.Flags().BoolVar(&viewMode, "view", false, "open in view-only mode")
	cmd.Flags().BoolVar(&editMode, "edit", false, "open in edit mode")
	return cmd
}

func ParseCLI() (model.Mode, string, error) {
	root := NewRootCmd()
	root.SetArgs(os.Args[1:])
	if err := root.Execute(); err != nil {
		return 0, "", err
	}
	return parsedMode, parsedPath, nil
}
