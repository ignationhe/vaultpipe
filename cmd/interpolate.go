package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultpipe/internal/envfile"
)

var (
	interpolateInput      string
	interpolateOutput     string
	interpolateFallback   bool
	interpolateErrMissing bool
	interpolateDefault    string
)

func init() {
	interpolateCmd := &cobra.Command{
		Use:   "interpolate",
		Short: "Resolve variable references within a .env file",
		Long:  "Reads a .env file and expands \${VAR} / $VAR references using values from the same file or the OS environment.",
		RunE:  runInterpolate,
	}

	interpolateCmd.Flags().StringVarP(&interpolateInput, "input", "i", ".env", "Input .env file")
	interpolateCmd.Flags().StringVarP(&interpolateOutput, "output", "o", "", "Output file (defaults to input file)")
	interpolateCmd.Flags().BoolVar(&interpolateFallback, "fallback-env", true, "Fall back to OS environment for missing vars")
	interpolateCmd.Flags().BoolVar(&interpolateErrMissing, "error-missing", false, "Return error when a referenced variable is not found")
	interpolateCmd.Flags().StringVar(&interpolateDefault, "default", "", "Default value for missing variables")

	rootCmd.AddCommand(interpolateCmd)
}

func runInterpolate(cmd *cobra.Command, args []string) error {
	env, err := envfile.Parse(interpolateInput)
	if err != nil {
		return fmt.Errorf("parse %q: %w", interpolateInput, err)
	}

	opts := envfile.InterpolateOptions{
		FallbackToEnv:  interpolateFallback,
		ErrorOnMissing: interpolateErrMissing,
		DefaultValue:   interpolateDefault,
	}

	resolved, err := envfile.Interpolate(env, opts)
	if err != nil {
		return fmt.Errorf("interpolate: %w", err)
	}

	dest := interpolateOutput
	if dest == "" {
		dest = interpolateInput
	}

	if err := envfile.Write(dest, resolved); err != nil {
		return fmt.Errorf("write %q: %w", dest, err)
	}

	fmt.Fprintf(os.Stdout, "interpolated %d keys → %s\n", len(resolved), dest)
	return nil
}
