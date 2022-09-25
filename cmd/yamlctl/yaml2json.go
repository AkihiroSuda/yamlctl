package main

import (
	"bytes"
	"io"

	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
)

func newYAML2JSONCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "yaml2json [FILE]",
		Short:   "YAML2JSON",
		Example: "  yamlctl yaml2json input.yaml",
		Args:    cobra.MaximumNArgs(1),
		RunE:    yaml2jsonAction,

		DisableFlagsInUseLine: true,
	}
	return cmd
}

func yaml2jsonAction(cmd *cobra.Command, args []string) error {
	var inName string
	if len(args) == 1 {
		inName = args[0]
	}
	in, err := readInput(cmd, inName)
	if err != nil {
		return err
	}
	res, err := yaml.YAMLToJSON(in)
	if err != nil {
		return err
	}
	_, err = io.Copy(cmd.OutOrStdout(), bytes.NewReader(res))
	return err
}
