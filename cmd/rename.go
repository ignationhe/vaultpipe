package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"vaultpipe/internal/envfile"
)

func init() {
	rootCmd.AddCommand(newRenameCmd())
}

func newRenameCmd() *cobra.Command {
	var (
		input         string
		output        string
		rules         []string
		keepOriginal  bool
		errorMissing  bool
		patternMode   bool
	)

	cmd := &cobra.Command{
		Use:   "rename",
		Short: "Rename keys in an env file using exact or pattern-based rules",
		Example: `  vaultpipe rename --input .env --rule OLD_KEY=NEW_KEY
  vaultpipe rename --input .env --rule '^OLD_(.*)=LEGACY_$1' --pattern`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRename(input, output, rules, keepOriginal, errorMissing, patternMode)
		},
	}

	cmd.Flags().StringVarP(&input, "input", "i", ".env", "source env file")
	cmd.Flags().StringVarP(&output, "output", "o", "", "destination file (defaults to input)")
	cmd.Flags().StringArrayVarP(&rules, "rule", "r", nil, "rename rule as FROM=TO (repeatable)")
	cmd.Flags().BoolVar(&keepOriginal, "keep-original", false, "preserve the original key after renaming")
	cmd.Flags().BoolVar(&errorMissing, "error-missing", false, "return error if a source key is absent")
	cmd.Flags().BoolVar(&patternMode, "pattern", false, "treat FROM as a regular expression")
	_ = cmd.MarkFlagRequired("rule")
	return cmd
}

func runRename(input, output string, rawRules []string, keepOriginal, errorMissing, patternMode bool) error {
	env, err := envfile.Parse(input)
	if err != nil {
		return fmt.Errorf("parse: %w", err)
	}

	var ruleList []envfile.RenameRule
	for _, r := range rawRules {
		parts := strings.SplitN(r, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid rule %q: expected FROM=TO", r)
		}
		ruleList = append(ruleList, envfile.RenameRule{
			From:    parts[0],
			To:      parts[1],
			Pattern: patternMode,
		})
	}

	opts := envfile.DefaultRenameOptions()
	opts.Rules = ruleList
	opts.DeleteOriginal = !keepOriginal
	opts.SkipMissing = !errorMissing

	renamed, err := envfile.Rename(env, opts)
	if err != nil {
		return err
	}

	dest := output
	if dest == "" {
		dest = input
	}

	if err := envfile.Write(renamed, dest); err != nil {
		return fmt.Errorf("write: %w", err)
	}

	fmt.Fprintf(os.Stdout, "renamed %d rule(s) → %s\n", len(ruleList), dest)
	return nil
}
