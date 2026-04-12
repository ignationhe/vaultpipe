package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/your-org/vaultpipe/internal/envfile"
)

func init() {
	transformCmd := newTransformCmd()
	rootCmd.AddCommand(transformCmd)
}

func newTransformCmd() *cobra.Command {
	var (
		inputFile  string
		outputFile string
		operation  string
		keys       []string
		skipErrors bool
	)

	cmd := &cobra.Command{
		Use:   "transform",
		Short: "Apply value transformations (uppercase, lowercase, trim) to an env file",
		Example: `  vaultpipe transform --input .env --op uppercase --keys SECRET_KEY
  vaultpipe transform --input .env --op trim --output .env.out`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runTransform(inputFile, outputFile, operation, keys, skipErrors)
		},
	}

	cmd.Flags().StringVarP(&inputFile, "input", "i", ".env", "Input env file")
	cmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output file (defaults to input file)")
	cmd.Flags().StringVar(&operation, "op", "trim", "Transform operation: uppercase | lowercase | trim")
	cmd.Flags().StringSliceVar(&keys, "keys", nil, "Keys to transform (empty = all keys via wildcard)")
	cmd.Flags().BoolVar(&skipErrors, "skip-errors", false, "Continue on transform errors")

	return cmd
}

func runTransform(input, output, operation string, keys []string, skipErrors bool) error {
	src, err := envfile.Parse(input)
	if err != nil {
		return fmt.Errorf("parse %q: %w", input, err)
	}

	var fn envfile.TransformFunc
	switch strings.ToLower(operation) {
	case "uppercase":
		fn = envfile.UppercaseValues()
	case "lowercase":
		fn = envfile.LowercaseValues()
	case "trim":
		fn = envfile.TrimSpaceValues()
	default:
		return fmt.Errorf("unknown operation %q: choose uppercase, lowercase, or trim", operation)
	}

	opts := envfile.DefaultTransformOptions()
	opts.SkipErrors = skipErrors

	if len(keys) == 0 {
		opts.Rules["*"] = fn
	} else {
		for _, k := range keys {
			opts.Rules[k] = fn
		}
	}

	result, err := envfile.Transform(src, opts)
	if err != nil {
		return fmt.Errorf("transform: %w", err)
	}

	dest := input
	if output != "" {
		dest = output
	}

	if err := envfile.Write(result, dest); err != nil {
		return fmt.Errorf("write %q: %w", dest, err)
	}

	fmt.Fprintf(os.Stdout, "transformed %d keys → %s\n", len(result), dest)
	return nil
}
