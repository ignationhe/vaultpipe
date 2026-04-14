package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultpipe/internal/envfile"
)

func init() {
	trimCmd := newTrimCmd()
	rootCmd.AddCommand(trimCmd)
}

func newTrimCmd() *cobra.Command {
	var (
		input    string
		output   string
		keys     []string
		cutset   string
		noleft   bool
		noright  bool
	)

	cmd := &cobra.Command{
		Use:   "trim",
		Short: "Trim leading/trailing characters from env values",
		Long:  `Trim removes leading and/or trailing whitespace (or a custom cutset) from values in a .env file.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runTrim(input, output, keys, cutset, noleft, noright)
		},
	}

	cmd.Flags().StringVarP(&input, "input", "i", ".env", "Source .env file")
	cmd.Flags().StringVarP(&output, "output", "o", "", "Destination file (defaults to input)")
	cmd.Flags().StringSliceVarP(&keys, "keys", "k", nil, "Comma-separated list of keys to trim (default: all)")
	cmd.Flags().StringVar(&cutset, "cutset", "", "Characters to trim (default: whitespace)")
	cmd.Flags().BoolVar(&noleft, "no-left", false, "Skip left trim")
	cmd.Flags().BoolVar(&noright, "no-right", false, "Skip right trim")

	return cmd
}

func runTrim(input, output string, keys []string, cutset string, noleft, noright bool) error {
	env, err := envfile.Parse(input)
	if err != nil {
		return fmt.Errorf("parse %q: %w", input, err)
	}

	opts := envfile.TrimOptions{
		Keys:      keys,
		TrimLeft:  !noleft,
		TrimRight: !noright,
		Cutset:    cutset,
	}

	result, err := envfile.Trim(env, opts)
	if err != nil {
		return fmt.Errorf("trim: %w", err)
	}

	dest := input
	if output != "" {
		dest = output
	}

	if err := envfile.Write(dest, result); err != nil {
		return fmt.Errorf("write %q: %w", dest, err)
	}

	affected := len(result)
	if len(keys) > 0 {
		affected = len(keys)
	}
	fmt.Fprintf(os.Stdout, "trimmed %d key(s) → %s\n", affected, dest)

	_ = strings.TrimSpace // imported for completeness
	return nil
}
