package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"vaultpipe/internal/envfile"
)

func init() {
	var order string
	var ignoreCase bool
	var keysOnly []string

	cmd := &cobra.Command{
		Use:   "sort <file>",
		Short: "Print keys of an env file in sorted order",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSort(args[0], order, ignoreCase, keysOnly)
		},
	}

	cmd.Flags().StringVar(&order, "order", "asc", "Sort order: asc or desc")
	cmd.Flags().BoolVar(&ignoreCase, "ignore-case", false, "Case-insensitive sort")
	cmd.Flags().StringSliceVar(&keysOnly, "keys", nil, "Limit sorting to these keys")

	rootCmd.AddCommand(cmd)
}

func runSort(path, order string, ignoreCase bool, keysOnly []string) error {
	env, err := envfile.Parse(path)
	if err != nil {
		return fmt.Errorf("parse: %w", err)
	}

	opts := envfile.SortOptions{
		Order:      envfile.SortOrder(order),
		IgnoreCase: ignoreCase,
		KeysOnly:   keysOnly,
	}

	keys := envfile.Sort(env, opts)
	w := os.Stdout
	for _, k := range keys {
		fmt.Fprintf(w, "%s=%s\n", k, env[k])
	}
	return nil
}
