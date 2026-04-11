package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultpipe/internal/envfile"
)

var auditCmd = &cobra.Command{
	Use:   "audit <before.env> <after.env>",
	Short: "Show an audit log of changes between two env files",
	Args:  cobra.ExactArgs(2),
	RunE:  runAudit,
}

func init() {
	auditCmd.Flags().StringP("format", "f", "text", "Output format: text or json")
	auditCmd.Flags().StringP("output", "o", "", "Write audit log to file (appends)")
	auditCmd.Flags().Bool("include-unchanged", false, "Include unchanged keys in the audit log")
	auditCmd.Flags().Bool("no-redact", false, "Do not redact secret values in the audit log")
	auditCmd.Flags().StringP("source", "s", "", "Label the source of the changes (e.g. vault path)")
	rootCmd.AddCommand(auditCmd)
}

func runAudit(cmd *cobra.Command, args []string) error {
	beforePath, afterPath := args[0], args[1]

	before, err := envfile.Parse(beforePath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("reading before file: %w", err)
	}
	after, err := envfile.Parse(afterPath)
	if err != nil {
		return fmt.Errorf("reading after file: %w", err)
	}

	diffs := envfile.Diff(before, after)

	opts := envfile.DefaultAuditOptions()
	if v, _ := cmd.Flags().GetBool("include-unchanged"); v {
		opts.IncludeUnchanged = true
	}
	if v, _ := cmd.Flags().GetBool("no-redact"); v {
		opts.RedactValues = false
	}
	if v, _ := cmd.Flags().GetString("source"); v != "" {
		opts.Source = v
	}

	log := envfile.Audit(diffs, opts)

	fmtStr, _ := cmd.Flags().GetString("format")
	fmt := envfile.AuditFormatText
	if fmtStr == "json" {
		fmt = envfile.AuditFormatJSON
	}

	outPath, _ := cmd.Flags().GetString("output")
	if outPath != "" {
		if err := envfile.WriteAuditLogToFile(log, outPath, fmt); err != nil {
			return err
		}
		cmd.Printf("Audit log written to %s\n", outPath)
		cmd.Printf("%s\n", log.Summary())
		return nil
	}

	return envfile.WriteAuditLog(log, cmd.OutOrStdout(), fmt)
}
