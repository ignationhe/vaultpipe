package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultpipe/internal/envfile"
)

var (
	auditEnvFile    string
	auditOutputFile string
	auditFormat     string
	auditRedact     bool
	auditUnchanged  bool
)

func init() {
	auditCmd := newAuditCmd()
	rootCmd.AddCommand(auditCmd)
}

func newAuditCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "audit",
		Short: "Audit changes between two .env files",
		Long: `Compare two .env files and produce an audit log of added, removed,
updated, and optionally unchanged keys. Output can be written as plain
text or JSON and optionally saved to a file.`,
		Example: `  vaultpipe audit --from .env.before --to .env --format json
  vaultpipe audit --from .env.before --to .env --output audit.log`,
		RunE: runAudit,
	}

	cmd.Flags().StringVar(&auditEnvFile, "from", "", "path to the previous .env file (required)")
	cmd.Flags().StringVar(&auditOutputFile, "output", "", "path to write the audit log (default: stdout)")
	cmd.Flags().StringVar(&auditFormat, "format", "text", "output format: text or json")
	cmd.Flags().BoolVar(&auditRedact, "redact", true, "redact sensitive values in the audit log")
	cmd.Flags().BoolVar(&auditUnchanged, "unchanged", false, "include unchanged keys in the audit log")
	cmd.Flags().StringVar(&syncEnvFile, "to", ".env", "path to the current .env file")

	_ = cmd.MarkFlagRequired("from")

	return cmd
}

func runAudit(cmd *cobra.Command, args []string) error {
	before, err := envfile.Parse(auditEnvFile)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("reading --from file: %w", err)
	}
	if before == nil {
		before = map[string]string{}
	}

	after, err := envfile.Parse(syncEnvFile)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("reading --to file: %w", err)
	}
	if after == nil {
		after = map[string]string{}
	}

	diffs := envfile.Diff(before, after)

	opts := envfile.DefaultAuditOptions()
	opts.IncludeUnchanged = auditUnchanged
	opts.Redact = auditRedact

	log, err := envfile.Audit(diffs, opts)
	if err != nil {
		return fmt.Errorf("generating audit log: %w", err)
	}

	if auditOutputFile != "" {
		if err := envfile.WriteAuditLogToFile(log, auditOutputFile, auditFormat); err != nil {
			return fmt.Errorf("writing audit log to file: %w", err)
		}
		fmt.Fprintf(cmd.OutOrStdout(), "audit log written to %s\n", auditOutputFile)
		return nil
	}

	return envfile.WriteAuditLog(log, cmd.OutOrStdout(), auditFormat)
}
