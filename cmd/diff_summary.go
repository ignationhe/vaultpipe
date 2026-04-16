package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"vaultpipe/internal/envfile"
)

var diffSummaryCmd = &cobra.Command{
	Use:   "diff-summary [file-a] [file-b]",
	Short: "Print a human-readable summary of differences between two .env files",
	Args:  cobra.ExactArgs(2),
	RunE:  runDiffSummary,
}

var (
	diffShowValues bool
	diffNoRedact   bool
)

func init() {
	diffSummaryCmd.Flags().BoolVar(&diffShowValues, "show-values", false, "include values in output")
	diffSummaryCmd.Flags().BoolVar(&diffNoRedact, "no-redact", false, "do not redact values")
	rootCmd.AddCommand(diffSummaryCmd)
}

func runDiffSummary(cmd *cobra.Command, args []string) error {
	a, err := envfile.Parse(args[0])
	if err != nil {
		return fmt.Errorf("reading %s: %w", args[0], err)
	}
	b, err := envfile.Parse(args[1])
	if err != nil {
		return fmt.Errorf("reading %s: %w", args[1], err)
	}

	diffs := envfile.Diff(a, b)
	opts := envfile.SummaryOptions{
		ShowValues: diffShowValues,
		Redact:     !diffNoRedact,
	}
	s := envfile.Summarize(diffs, opts)

	for _, l := range s.Lines {
		fmt.Fprintln(os.Stdout, l)
	}
	fmt.Fprintf(os.Stdout, "\nadded: %d  removed: %d  updated: %d\n", s.Added, s.Removed, s.Updated)
	return nil
}
