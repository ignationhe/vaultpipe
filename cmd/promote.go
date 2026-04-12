package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultpipe/internal/envfile"
)

func init() {
	promoteCmd := newPromoteCmd()
	rootCmd.AddCommand(promoteCmd)
}

func newPromoteCmd() *cobra.Command {
	var (
		keys      []string
		overwrite bool
		fromEnv   string
		toEnv     string
	)

	cmd := &cobra.Command{
		Use:   "promote <src-file> <dst-file>",
		Short: "Promote secrets from one env file to another",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runPromote(args[0], args[1], fromEnv, toEnv, keys, overwrite)
		},
	}

	cmd.Flags().StringSliceVarP(&keys, "keys", "k", nil, "Comma-separated keys to promote (default: all)")
	cmd.Flags().BoolVarP(&overwrite, "overwrite", "o", false, "Overwrite existing keys in destination")
	cmd.Flags().StringVar(&fromEnv, "from", "", "Source environment label (informational)")
	cmd.Flags().StringVar(&toEnv, "to", "", "Destination environment label (informational)")

	return cmd
}

func runPromote(srcPath, dstPath, fromEnv, toEnv string, keys []string, overwrite bool) error {
	src, err := envfile.Parse(srcPath)
	if err != nil {
		return fmt.Errorf("reading source %q: %w", srcPath, err)
	}

	dst, err := envfile.Parse(dstPath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("reading destination %q: %w", dstPath, err)
	}
	if dst == nil {
		dst = map[string]string{}
	}

	opts := envfile.DefaultPromoteOptions()
	opts.Keys = keys
	opts.Overwrite = overwrite
	opts.FromEnv = fromEnv
	opts.ToEnv = toEnv

	out, result, err := envfile.Promote(src, dst, opts)
	if err != nil {
		return err
	}

	if err := envfile.Write(dstPath, out); err != nil {
		return fmt.Errorf("writing destination %q: %w", dstPath, err)
	}

	if len(result.Promoted) > 0 {
		fmt.Printf("promoted:    %s\n", strings.Join(result.Promoted, ", "))
	}
	if len(result.Overwritten) > 0 {
		fmt.Printf("overwritten: %s\n", strings.Join(result.Overwritten, ", "))
	}
	if len(result.Skipped) > 0 {
		fmt.Printf("skipped:     %s\n", strings.Join(result.Skipped, ", "))
	}

	return nil
}
