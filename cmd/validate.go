package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/your-org/vaultpipe/internal/envfile"
)

func init() {
	validateCmd.Flags().StringP("file", "f", ".env", "Path to the .env file to validate")
	validateCmd.Flags().Bool("require-uppercase", false, "Require all keys to be UPPER_CASE")
	validateCmd.Flags().Bool("forbid-empty", false, "Reject keys with empty values")
	RootCmd.AddCommand(validateCmd)
}

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate the syntax and conventions of a .env file",
	Long: `Parses a .env file and checks every key against naming rules.
Optional flags enforce stricter conventions such as UPPER_CASE keys
or non-empty values.`,
	RunE: runValidate,
}

func runValidate(cmd *cobra.Command, _ []string) error {
	filePath, err := cmd.Flags().GetString("file")
	if err != nil {
		return err
	}
	requireUpper, _ := cmd.Flags().GetBool("require-uppercase")
	forbidEmpty, _ := cmd.Flags().GetBool("forbid-empty")

	env, err := envfile.Parse(filePath)
	if err != nil {
		// Provide a more actionable message when the file does not exist.
		if os.IsNotExist(err) {
			return fmt.Errorf("file not found: %s", filePath)
		}
		return fmt.Errorf("failed to parse %s: %w", filePath, err)
	}

	opts := envfile.ValidateOptions{
		RequireUppercase: requireUpper,
		ForbidEmpty:      forbidEmpty,
	}

	result := envfile.Validate(env, opts)
	if result.HasErrors() {
		fmt.Fprintf(os.Stderr, "validation failed for %s:\n", filePath)
		for _, e := range result.Errors {
			fmt.Fprintf(os.Stderr, "  - %s\n", e.Error())
		}
		return fmt.Errorf("%d validation error(s) found", len(result.Errors))
	}

	fmt.Printf("✓ %s is valid (%d keys checked)\n", filePath, len(env))
	return nil
}
