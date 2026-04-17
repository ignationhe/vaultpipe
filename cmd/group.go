package cmd

import (
	"fmt"
	"os"
	"sort"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultpipe/internal/envfile"
)

func init() {
	var separator string
	var keepPrefix bool
	var file string

	cmd := &cobra.Command{
		Use:   "group",
		Short: "Group env keys by prefix into named buckets",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGroup(file, separator, keepPrefix)
		},
	}

	cmd.Flags().StringVarP(&file, "file", "f", ".env", "Source .env file")
	cmd.Flags().StringVar(&separator, "sep", "_", "Key separator")
	cmd.Flags().BoolVar(&keepPrefix, "keep-prefix", false, "Retain group prefix in key names")

	rootCmd.AddCommand(cmd)
}

func runGroup(file, separator string, keepPrefix bool) error {
	env, err := envfile.Parse(file)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("file not found: %s", file)
		}
		return err
	}

	opts := envfile.GroupOptions{
		Separator:  separator,
		KeepPrefix: keepPrefix,
	}
	groups := envfile.Group(env, opts)

	// Print groups in sorted order.
	groupNames := make([]string, 0, len(groups))
	for g := range groups {
		groupNames = append(groupNames, g)
	}
	sort.Strings(groupNames)

	for _, g := range groupNames {
		fmt.Printf("[%s]\n", g)
		keys := make([]string, 0, len(groups[g]))
		for k := range groups[g] {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			fmt.Printf("  %s=%s\n", k, groups[g][k])
		}
	}
	return nil
}
