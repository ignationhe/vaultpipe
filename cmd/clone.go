package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultpipe/internal/envfile"
)

func init() {
	rootCmd.AddCommand(newCloneCmd())
}

func newCloneCmd() *cobra.Command {
	var (
		overwrite  bool
		filterKeys []string
		prefix     string
	)

	cmd := &cobra.Command{
		Use:   "clone <src> <dst>",
		Short: "Clone an env file, optionally filtering or prefixing keys",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runClone(args[0], args[1], overwrite, filterKeys, prefix)
		},
	}

	cmd.Flags().BoolVar(&overwrite, "overwrite", false, "overwrite destination if it exists")
	cmd.Flags().StringSliceVar(&filterKeys, "keys", nil, "comma-separated list of keys to include")
	cmd.Flags().StringVar(&prefix, "prefix", "", "prefix to prepend to every key in the destination")

	return cmd
}

func runClone(src, dst string, overwrite bool, filterKeys []string, prefix string) error {
	opts := envfile.DefaultCloneOptions()
	opts.Overwrite = overwrite
	opts.FilterKeys = filterKeys

	var n int
	var err error

	if prefix != "" {
		prefix = strings.ToUpper(prefix)
		if !strings.HasSuffix(prefix, "_") {
			prefix += "_"
		}
		n, err = envfile.CloneWithPrefix(src, dst, prefix, opts)
	} else {
		n, err = envfile.Clone(src, dst, opts)
	}

	if err != nil {
		return fmt.Errorf("clone failed: %w", err)
	}

	fmt.Printf("cloned %d key(s) from %s → %s\n", n, src, dst)
	return nil
}
