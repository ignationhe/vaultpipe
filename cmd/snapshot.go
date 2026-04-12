package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultpipe/internal/envfile"
)

var (
	snapshotDir    string
	snapshotSource string
)

func init() {
	snapshotCmd := newSnapshotCmd()
	rootCmd.AddCommand(snapshotCmd)
}

func newSnapshotCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "snapshot",
		Short: "Save a point-in-time snapshot of an env file",
		Long: `Reads an .env file and writes its current key-value pairs to a
timestamped JSON snapshot file for auditing or rollback purposes.`,
		RunE: runSnapshot,
	}
	cmd.Flags().StringVarP(&snapshotSource, "file", "f", ".env", "env file to snapshot")
	cmd.Flags().StringVarP(&snapshotDir, "dir", "d", "", "directory to store snapshots (default: .vaultpipe/snapshots)")
	return cmd
}

func runSnapshot(cmd *cobra.Command, _ []string) error {
	env, err := envfile.Parse(snapshotSource)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("snapshot: file not found: %s", snapshotSource)
		}
		return fmt.Errorf("snapshot: parse: %w", err)
	}

	path, err := envfile.SaveSnapshot(env, snapshotSource, snapshotDir)
	if err != nil {
		return fmt.Errorf("snapshot: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "snapshot saved: %s (%d keys)\n", path, len(env))
	return nil
}
