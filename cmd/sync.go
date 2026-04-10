package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/your-org/vaultpipe/internal/sync"
	"github.com/your-org/vaultpipe/internal/vault"
)

func runSync(cmd *cobra.Command, args []string) error {
	cfg, err := vault.ConfigFromEnv()
	if err != nil {
		return fmt.Errorf("vault config: %w", err)
	}

	// CLI flags override environment variables
	if vaultAddr != "" {
		cfg.Address = vaultAddr
	}
	if vaultToken != "" {
		cfg.Token = vaultToken
	}

	if cfg.Address == "" {
		return fmt.Errorf("vault address is required: set VAULT_ADDR or use --vault-addr")
	}
	if cfg.Token == "" {
		return fmt.Errorf("vault token is required: set VAULT_TOKEN or use --vault-token")
	}

	client, err := vault.NewClientFromConfig(cfg)
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}

	syncer := sync.New(client, sync.Options{
		SecretPath: vaultPath,
		EnvFile:    envFile,
		DryRun:     dryRun,
	})

	result, err := syncer.Sync(cmd.Context())
	if err != nil {
		return fmt.Errorf("sync failed: %w", err)
	}

	if dryRun {
		fmt.Fprintln(os.Stdout, "[dry-run] no changes written")
	}

	if !result.Changed {
		fmt.Fprintln(os.Stdout, "no changes detected")
		return nil
	}

	fmt.Fprintf(os.Stdout, "synced %d secret(s) to %s\n", result.Count, envFile)
	for _, d := range result.Diffs {
		fmt.Fprintf(os.Stdout, "  %s %s\n", d.Status, d.Key)
	}

	return nil
}
