package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"vaultpipe/internal/envfile"
)

var maskCmd = &cobra.Command{
	Use:   "mask [file]",
	Short: "Print env file with sensitive values masked",
	Args:  cobra.ExactArgs(1),
	RunE:  runMask,
}

func init() {
	maskCmd.Flags().StringSlice("patterns", nil, "Additional key patterns to mask (comma-separated)")
	maskCmd.Flags().String("placeholder", "****", "Replacement string for masked values")
	maskCmd.Flags().Int("partial", 0, "Reveal first N characters before placeholder")
	maskCmd.Flags().Bool("case-sensitive", false, "Use case-sensitive pattern matching")
	RootCmd.AddCommand(maskCmd)
}

func runMask(cmd *cobra.Command, args []string) error {
	filePath := args[0]

	env, err := envfile.Parse(filePath)
	if err != nil {
		return fmt.Errorf("parse: %w", err)
	}

	opts := envfile.DefaultMaskOptions()

	if p, _ := cmd.Flags().GetString("placeholder"); p != "" {
		opts.Placeholder = p
	}
	if n, _ := cmd.Flags().GetInt("partial"); n > 0 {
		opts.MaskPartial = n
	}
	if cs, _ := cmd.Flags().GetBool("case-sensitive"); cs {
		opts.CaseSensitive = cs
	}
	if extra, _ := cmd.Flags().GetStringSlice("patterns"); len(extra) > 0 {
		opts.Patterns = append(opts.Patterns, extra...)
	}

	masked := envfile.Mask(env, opts)

	for k, v := range masked {
		fmt.Fprintf(os.Stdout, "%s=%s\n", k, quoteIfNeeded(v))
	}
	return nil
}

func quoteIfNeeded(v string) string {
	if strings.ContainsAny(v, " \t\n") {
		return `"` + v + `"`
	}
	return v
}
