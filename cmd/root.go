package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	vaultAddr  string
	vaultToken string
	vaultPath  string
	envFile    string
	dryRun     bool
)

var rootCmd = &cobra.Command{
	Use:   "vaultpipe",
	Short: "Sync secrets from HashiCorp Vault into local .env files",
	Long: `vaultpipe pulls secrets from a HashiCorp Vault KV path
and writes them into a local .env file with diff-aware updates.

Set VAULT_ADDR and VAULT_TOKEN environment variables or use flags.`,
	RunE: runSync,
}

func init() {
	rootCmd.PersistentFlags().StringVar(&vaultAddr, "vault-addr", "", "Vault server address (overrides VAULT_ADDR)")
	rootCmd.PersistentFlags().StringVar(&vaultToken, "vault-token", "", "Vault token (overrides VAULT_TOKEN)")
	rootCmd.PersistentFlags().StringVar(&vaultPath, "path", "", "Vault KV secret path (required)")
	rootCmd.PersistentFlags().StringVar(&envFile, "env-file", ".env", "Target .env file path")
	rootCmd.PersistentFlags().BoolVar(&dryRun, "dry-run", false, "Preview changes without writing to disk")

	_ = rootCmd.MarkPersistentFlagRequired("path")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
