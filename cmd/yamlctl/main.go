package main

import (
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func main() {
	if err := newApp().Execute(); err != nil {
		logrus.Fatal(err)
	}
}

func newApp() *cobra.Command {
	var rootCmd = &cobra.Command{
		Use:   "yamlctl",
		Short: "YAML manipulation utility",
		Example: `  Edit:
  $ yamlctl edit -w --bak --set $.foo.bar=baz input.yaml
`,
		Version:       strings.TrimPrefix(version(), "v"),
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	rootCmd.PersistentFlags().Bool("debug", false, "debug mode")
	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		debug, _ := cmd.Flags().GetBool("debug")
		if debug {
			logrus.SetLevel(logrus.DebugLevel)
		}
		return nil
	}
	rootCmd.AddCommand(
		newEditCommand(),
		newEditableCommand(),
		newQueryCommand(),
		newYAML2JSONCommand(),
	)
	return rootCmd
}
