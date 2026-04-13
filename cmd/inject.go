package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"vaultpipe/internal/envfile"
)

var injectCmd = &cobra.Command{
	Use:   "inject [env-file]",
	Short: "Inject variables from a .env file into the current process environment",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runInject,
}

func init() {
	injectCmd.Flags().String("prefix", "", "Only inject keys with this prefix")
	injectCmd.Flags().Bool("strip-prefix", false, "Strip the prefix from keys before injecting")
	injectCmd.Flags().Bool("overwrite", false, "Overwrite existing environment variables")
	injectCmd.Flags().Bool("dry-run", false, "Print what would be injected without setting variables")
	rootCmd.AddCommand(injectCmd)
}

func runInject(cmd *cobra.Command, args []string) error {
	path := ".env"
	if len(args) == 1 {
		path = args[0]
	}

	prefix, _ := cmd.Flags().GetString("prefix")
	stripPrefix, _ := cmd.Flags().GetBool("strip-prefix")
	overwrite, _ := cmd.Flags().GetBool("overwrite")
	dryRun, _ := cmd.Flags().GetBool("dry-run")

	env, err := envfile.Parse(path)
	if err != nil {
		return fmt.Errorf("inject: parse %q: %w", path, err)
	}

	opts := envfile.DefaultInjectOptions()
	opts.Prefix = prefix
	opts.StripPrefix = stripPrefix
	opts.Overwrite = overwrite
	opts.DryRun = dryRun

	results, err := envfile.Inject(env, opts)
	if err != nil {
		return err
	}

	for _, r := range results {
		switch {
		case r.DryRun:
			fmt.Printf("[dry-run] would inject %s\n", r.Key)
		case r.Skipped:
			fmt.Printf("[skip]    %s (already set)\n", r.Key)
		default:
			fmt.Printf("[inject]  %s\n", r.Key)
		}
	}

	return nil
}
