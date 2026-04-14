package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"vaultpipe/internal/envfile"
)

func init() {
	normalizeCmd := newNormalizeCmd()
	rootCmd.AddCommand(normalizeCmd)
}

func newNormalizeCmd() *cobra.Command {
	var (
		input              string
		output             string
		uppercase          bool
		trim               bool
		replaceHyphens     bool
		collapseUnderscores bool
	)

	cmd := &cobra.Command{
		Use:   "normalize",
		Short: "Normalize env file keys and values",
		Long:  `Normalize applies formatting rules to keys and values in an env file, such as uppercasing keys, trimming whitespace, and replacing hyphens.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runNormalize(input, output, envfile.NormalizeOptions{
				UppercaseKeys:       uppercase,
				TrimValues:          trim,
				ReplaceHyphens:      replaceHyphens,
				CollapseUnderscores: collapseUnderscores,
			})
		},
	}

	cmd.Flags().StringVarP(&input, "input", "i", ".env", "Input env file")
	cmd.Flags().StringVarP(&output, "output", "o", "", "Output file (defaults to input)")
	cmd.Flags().BoolVar(&uppercase, "uppercase", true, "Uppercase all keys")
	cmd.Flags().BoolVar(&trim, "trim", true, "Trim whitespace from values")
	cmd.Flags().BoolVar(&replaceHyphens, "replace-hyphens", true, "Replace hyphens with underscores in keys")
	cmd.Flags().BoolVar(&collapseUnderscores, "collapse-underscores", false, "Collapse consecutive underscores in keys")

	return cmd
}

func runNormalize(input, output string, opts envfile.NormalizeOptions) error {
	env, err := envfile.Parse(input)
	if err != nil {
		return fmt.Errorf("parse %q: %w", input, err)
	}

	normalized, err := envfile.Normalize(env, opts)
	if err != nil {
		return fmt.Errorf("normalize: %w", err)
	}

	dest := output
	if dest == "" {
		dest = input
	}

	if err := envfile.Write(dest, normalized); err != nil {
		return fmt.Errorf("write %q: %w", dest, err)
	}

	fmt.Fprintf(os.Stdout, "normalized %d keys → %s\n", len(normalized), dest)
	return nil
}
