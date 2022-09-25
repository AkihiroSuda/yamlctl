package main

import (
	"errors"
	"fmt"

	"github.com/AkihiroSuda/yamlctl/pkg/yamlutil"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func newEditableCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "editable [FILE]",
		Short:   "Check whether safely editable",
		Example: "  yamlctl editable input.yaml",
		Args:    cobra.MaximumNArgs(1),
		RunE:    editableAction,

		DisableFlagsInUseLine: true,
	}
	return cmd
}

func editableAction(cmd *cobra.Command, args []string) error {
	var inName string
	if len(args) == 1 {
		inName = args[0]
	}
	in, err := readInput(cmd, inName)
	if err != nil {
		return err
	}
	if err = yamlutil.Editable(in); err != nil {
		logrus.WithError(err).Debug("the YAML is unsafe to edit")
		fmt.Fprintln(cmd.OutOrStdout(), "false")
		return errors.New("the YAML is unsafe to edit (set --debug to show the reason)")
	}
	fmt.Fprintln(cmd.OutOrStdout(), "true")
	return nil
}
