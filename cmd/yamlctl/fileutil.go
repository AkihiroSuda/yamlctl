package main

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

func readInput(cmd *cobra.Command, inName string) ([]byte, error) {
	switch inName {
	case "", "-":
		in, err := io.ReadAll(cmd.InOrStdin())
		if err != nil {
			return nil, fmt.Errorf("failed to read input: %w", err)
		}
		return in, nil
	default:
		return os.ReadFile(inName)
	}
}
