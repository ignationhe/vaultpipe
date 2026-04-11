package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourorg/vaultpipe/internal/envfile"
)

func init() {
	schemaCmd := newSchemaCmd()
	rootCmd.AddCommand(schemaCmd)
}

func newSchemaCmd() *cobra.Command {
	var envPath string
	var schemaPath string

	cmd := &cobra.Command{
		Use:   "schema",
		Short: "Validate an .env file against a JSON schema",
		Long: `Reads a JSON schema file describing required keys and value patterns,
then validates the given .env file against it. Exits non-zero if violations are found.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSchema(envPath, schemaPath)
		},
	}

	cmd.Flags().StringVarP(&envPath, "env", "e", ".env", "path to the .env file")
	cmd.Flags().StringVarP(&schemaPath, "schema", "s", ".env.schema.json", "path to the JSON schema file")

	return cmd
}

func runSchema(envPath, schemaPath string) error {
	env, err := envfile.Parse(envPath)
	if err != nil {
		return fmt.Errorf("failed to parse env file: %w", err)
	}

	schema, err := envfile.LoadSchema(schemaPath)
	if err != nil {
		return fmt.Errorf("failed to load schema: %w", err)
	}

	violations := envfile.ValidateSchema(env, schema)
	if len(violations) == 0 {
		fmt.Println("schema validation passed: no violations found")
		return nil
	}

	fmt.Fprintf(os.Stderr, "schema validation failed: %d violation(s)\n", len(violations))
	for _, v := range violations {
		fmt.Fprintf(os.Stderr, "  - %s\n", v.Error())
	}
	return fmt.Errorf("schema validation failed with %d violation(s)", len(violations))
}
