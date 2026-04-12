package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultpipe/internal/envfile"
)

var convertCmd = &cobra.Command{
	Use:   "convert [file]",
	Short: "Convert a .env file to another format (dotenv, yaml, toml)",
	Args:  cobra.ExactArgs(1),
	RunE:  runConvert,
}

func init() {
	convertCmd.Flags().StringP("format", "f", "dotenv", "Output format: dotenv | yaml | toml")
	convertCmd.Flags().StringP("prefix", "p", "", "Prefix to prepend to every key")
	convertCmd.Flags().StringP("output", "o", "", "Write output to file instead of stdout")
	convertCmd.Flags().Bool("no-sort", false, "Disable alphabetical key sorting")
	rootCmd.AddCommand(convertCmd)
}

func runConvert(cmd *cobra.Command, args []string) error {
	src := args[0]

	formatStr, _ := cmd.Flags().GetString("format")
	prefix, _ := cmd.Flags().GetString("prefix")
	outputFile, _ := cmd.Flags().GetString("output")
	noSort, _ := cmd.Flags().GetBool("no-sort")

	env, err := envfile.Parse(src)
	if err != nil {
		return fmt.Errorf("convert: parse %q: %w", src, err)
	}

	opts := envfile.ConvertOptions{
		Format: envfile.ConvertFormat(formatStr),
		Sort:   !noSort,
		Prefix: prefix,
	}

	out, err := envfile.Convert(env, opts)
	if err != nil {
		return fmt.Errorf("convert: %w", err)
	}

	if outputFile != "" {
		if err := os.WriteFile(outputFile, []byte(out), 0o644); err != nil {
			return fmt.Errorf("convert: write %q: %w", outputFile, err)
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Written to %s\n", outputFile)
		return nil
	}

	fmt.Fprint(cmd.OutOrStdout(), out)
	return nil
}
