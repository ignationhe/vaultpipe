package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultpipe/internal/envfile"
)

var (
	tagFile       string
	tagRules      []string
	tagSkipUntagged bool
)

func init() {
	tagCmd := newTagCmd()
	rootCmd.AddCommand(tagCmd)
}

func newTagCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tag",
		Short: "Annotate env keys with user-defined tags",
		Long: `Tag reads a .env file and prints each key alongside any matching tags.

Rules are provided as --rule TAG=PATTERN flags where PATTERN supports a
trailing '*' wildcard (e.g. DB_*). Multiple rules may be specified.`,
		RunE: runTag,
	}
	cmd.Flags().StringVarP(&tagFile, "file", "f", ".env", "path to the .env file")
	cmd.Flags().StringArrayVar(&tagRules, "rule", nil, "tag rule in TAG=PATTERN format (repeatable)")
	cmd.Flags().BoolVar(&tagSkipUntagged, "skip-untagged", false, "omit keys that match no rule")
	return cmd
}

func runTag(cmd *cobra.Command, _ []string) error {
	env, err := envfile.Parse(tagFile)
	if err != nil {
		return fmt.Errorf("tag: parse %q: %w", tagFile, err)
	}

	opts := envfile.DefaultTagOptions()
	opts.SkipUntagged = tagSkipUntagged

	for _, rule := range tagRules {
		parts := strings.SplitN(rule, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("tag: invalid rule %q — expected TAG=PATTERN", rule)
		}
		tag, pattern := parts[0], parts[1]
		opts.Rules[tag] = append(opts.Rules[tag], pattern)
	}

	entries, err := envfile.Tag(env, opts)
	if err != nil {
		return fmt.Errorf("tag: %w", err)
	}

	for _, e := range entries {
		tagStr := "(untagged)"
		if len(e.Tags) > 0 {
			tagStr = strings.Join(e.Tags, ", ")
		}
		fmt.Fprintf(cmd.OutOrStdout(), "%-30s %s\n", e.Key, tagStr)
	}
	return nil
}
