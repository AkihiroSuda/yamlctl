package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/AkihiroSuda/yamlctl/pkg/yamlutil"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func newEditCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "edit [-w] [--bak] [--force] [--set YAMLPATH=VALUE]... [FILE]",
		Short:   "Edit, preserving (most of) the comment lines",
		Example: "  yamlctl edit -w --bak --set $.foo.bar=baz input.yaml",
		Args:    cobra.MaximumNArgs(1),
		RunE:    editAction,

		DisableFlagsInUseLine: true,
	}
	flags := cmd.Flags()
	flags.StringSlice("set", nil, "Set YAMLPATH=VALUE")
	flags.BoolP("write", "w", false, "Write the result to the input file")
	flags.Bool("bak", false, "Create a backup file (Use with -w)")
	flags.Bool("force", false, "Edit forcibly even if unsafe")
	return cmd
}

func editAction(cmd *cobra.Command, args []string) error {
	flags := cmd.Flags()
	setFlags, err := flags.GetStringSlice("set")
	if err != nil {
		return err
	}
	if len(setFlags) == 0 {
		return errors.New("no --set flag was specified")
	}
	wFlag, err := flags.GetBool("write")
	if err != nil {
		return err
	}
	bakFlag, err := flags.GetBool("bak")
	if err != nil {
		return err
	}
	if !wFlag && bakFlag {
		return errors.New("--bak needs to be used in conjunction with -w")
	}
	forceFlag, err := flags.GetBool("force")
	if err != nil {
		return err
	}

	var ops []yamlutil.Op
	for _, f := range setFlags {
		split := strings.SplitN(f, "=", 2)
		if len(split) < 2 {
			return fmt.Errorf("invalid --set flag: %q", f)
		}
		op := yamlutil.Op{
			Type:  yamlutil.OpSet,
			Path:  split[0],
			Value: split[1],
		}
		ops = append(ops, op)
	}

	var inName string
	if len(args) == 1 {
		inName = args[0]
	}
	if wFlag && (inName == "" || inName == "-") {
		return fmt.Errorf("input file has to be specified for -w mode")
	}
	in, err := readInput(cmd, inName)
	if err != nil {
		return err
	}
	if err = yamlutil.Editable(in); err != nil {
		logrus.WithError(err).Debug("the YAML is unsafe to edit")
		if !forceFlag {
			return errors.New("the YAML is unsafe to edit (Hint: set --force to edit forcibly)")
		}
		logrus.Warn("the YAML is unsafe to edit")
	}

	res, err := yamlutil.Apply(in, ops...)
	if err != nil {
		return err
	}

	if wFlag {
		if bakFlag {
			bakNameGlob := inName + ".bak.*"
			existingBakFiles, err := filepath.Glob(bakNameGlob)
			if err != nil {
				return err
			}
			bakName := fmt.Sprintf("%s.bak.%d", inName, len(existingBakFiles))
			logrus.Infof("Creating a backup file %q", bakName)
			if err = os.WriteFile(bakName, in, 0644); err != nil {
				return err
			}
		}
		if err = os.WriteFile(inName, res, 0644); err != nil {
			return err
		}
		return nil
	}
	_, err = io.Copy(cmd.OutOrStdout(), bytes.NewReader(res))
	return err
}
