package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultpipe/internal/envfile"
)

func init() {
	var strategy string
	var suffix string
	var output string

	mergeCmd := &cobra.Command{
		Use:   "merge <base-file> <incoming-file>",
		Short: "Merge two .env files with configurable conflict resolution",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runMerge(args[0], args[1], strategy, suffix, output)
		},
	}

	mergeCmd.Flags().StringVarP(&strategy, "strategy", "s", "keep",
		"Conflict strategy: keep, overwrite, both")
	mergeCmd.Flags().StringVar(&suffix, "suffix", "_NEW",
		"Suffix for incoming key when strategy=both")
	mergeCmd.Flags().StringVarP(&output, "output", "o", "",
		"Output file (defaults to base file)")

	rootCmd.AddCommand(mergeCmd)
}

func runMerge(baseFile, incomingFile, strategy, suffix, output string) error {
	base, err := envfile.Parse(baseFile)
	if err != nil {
		return fmt.Errorf("parsing base file: %w", err)
	}

	incoming, err := envfile.Parse(incomingFile)
	if err != nil {
		return fmt.Errorf("parsing incoming file: %w", err)
	}

	opts := envfile.MergeOptions{ConflictSuffix: suffix}
	switch strategy {
	case "overwrite":
		opts.Strategy = envfile.StrategyOverwrite
	case "both":
		opts.Strategy = envfile.StrategyKeepBoth
	default:
		opts.Strategy = envfile.StrategyKeepExisting
	}

	merged := envfile.MergeWith(base, incoming, opts)

	dest := baseFile
	if output != "" {
		dest = output
	}

	if err := envfile.Write(dest, merged, os.FileMode(0600)); err != nil {
		return fmt.Errorf("writing merged file: %w", err)
	}

	fmt.Printf("merged %d keys into %s\n", len(merged), dest)
	return nil
}
