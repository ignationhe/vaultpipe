package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultpipe/internal/envfile"
)

func init() {
	rootCmd.AddCommand(newScopeCmd())
}

func newScopeCmd() *cobra.Command {
	var (
		input  string
		output string
		scopes []string
		target string
		sep    string
		keep   bool
	)

	cmd := &cobra.Command{
		Use:   "scope",
		Short: "Resolve a scoped env file for a target environment",
		Long: `Reads a flat .env file containing scope-prefixed keys (e.g. PROD__DB_URL)
and emits a resolved file containing only keys for the target scope,
with global (unscoped) keys included as defaults.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runScope(input, output, scopes, target, sep, keep)
		},
	}

	cmd.Flags().StringVarP(&input, "input", "i", ".env", "Source .env file")
	cmd.Flags().StringVarP(&output, "output", "o", "", "Output file (default: stdout)")
	cmd.Flags().StringSliceVarP(&scopes, "scopes", "s", nil, "Ordered list of scopes (e.g. base,staging,prod)")
	cmd.Flags().StringVarP(&target, "target", "t", "", "Target scope to resolve")
	cmd.Flags().StringVar(&sep, "sep", "__", "Separator between scope prefix and key")
	cmd.Flags().BoolVar(&keep, "keep-prefix", false, "Keep scope prefix in output keys")
	_ = cmd.MarkFlagRequired("target")
	_ = cmd.MarkFlagRequired("scopes")

	return cmd
}

func runScope(input, output string, scopes []string, target, sep string, keepPrefix bool) error {
	env, err := envfile.Parse(input)
	if err != nil {
		return fmt.Errorf("parse %q: %w", input, err)
	}

	opts := envfile.DefaultScopeOptions()
	opts.Scopes = scopes
	opts.TargetScope = target
	opts.Separator = sep
	opts.StripPrefix = !keepPrefix

	resolved, err := envfile.Scope(env, opts)
	if err != nil {
		return fmt.Errorf("scope: %w", err)
	}

	if output == "" {
		for k, v := range resolved {
			if strings.ContainsAny(v, " \t") {
				fmt.Fprintf(os.Stdout, "%s=%q\n", k, v)
			} else {
				fmt.Fprintf(os.Stdout, "%s=%s\n", k, v)
			}
		}
		return nil
	}

	if err := envfile.Write(output, resolved); err != nil {
		return fmt.Errorf("write %q: %w", output, err)
	}
	fmt.Fprintf(os.Stderr, "wrote %d keys to %s\n", len(resolved), output)
	return nil
}
