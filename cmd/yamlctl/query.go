package main

import (
	"bytes"
	"io"

	"github.com/AkihiroSuda/yamlctl/pkg/yamlutil"
	"github.com/spf13/cobra"
)

func newQueryCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "query YAMLPATH [FILE]",
		Short:   "Query",
		Example: "  yamlctl query $.foo.bar input.yaml",
		Args:    cobra.MaximumNArgs(2),
		RunE:    queryAction,

		DisableFlagsInUseLine: true,
	}
	return cmd
}

func queryAction(cmd *cobra.Command, args []string) error {
	query := args[0]
	var inName string
	if len(args) == 2 {
		inName = args[1]
	}
	in, err := readInput(cmd, inName)
	if err != nil {
		return err
	}
	res, err := yamlutil.Query(in, query)
	if err != nil {
		return err
	}
	res = append(res, byte('\n'))
	_, err = io.Copy(cmd.OutOrStdout(), bytes.NewReader(res))
	return err
}
