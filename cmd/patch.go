package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"vaultpipe/internal/envfile"
)

func init() {
	rootCmd.AddCommand(newPatchCmd())
}

func newPatchCmd() *cobra.Command {
	var (
		filePath      string
		setFlags      []string
		deleteFlags   []string
		renameFlags   []string
		ignoreMissing bool
		dryRun        bool
	)

	cmd := &cobra.Command{
		Use:   "patch",
		Short: "Apply set/delete/rename operations to an env file",
		Example: `  vaultpipe patch --file .env --set KEY=value --delete OLD_KEY --rename OLD=NEW`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runPatch(filePath, setFlags, deleteFlags, renameFlags, ignoreMissing, dryRun)
		},
	}

	cmd.Flags().StringVarP(&filePath, "file", "f", ".env", "path to the env file")
	cmd.Flags().StringArrayVar(&setFlags, "set", nil, "set KEY=VALUE (repeatable)")
	cmd.Flags().StringArrayVar(&deleteFlags, "delete", nil, "delete KEY (repeatable)")
	cmd.Flags().StringArrayVar(&renameFlags, "rename", nil, "rename OLD=NEW (repeatable)")
	cmd.Flags().BoolVar(&ignoreMissing, "ignore-missing", true, "silently skip missing keys in delete/rename")
	cmd.Flags().BoolVar(&dryRun, "dry- result without writing")
	return cmd
}

func runPatch(filePath string, sets, deletes, renames []string, ignoreMissing, dryRun bool) error {
	env, err := envfile.Parse(filePath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("parse %s: %w", filePath, err)
if env == nil {
		env = map[string]string{}
	}

	var rules []envfile.PatchRule
	for _, s := range sets {
		k, v, ok := splitKV(s)
		if !ok {
			return fmt.Errorf("--set %q: expected KEY=VALUE", s)
		}
		rules = append(rules, envfile.PatchRule{Op: envfile.PatchSet, Key: k, Value: v})
	}
	for _, d := range deletes {
		rules = append(rules, envfile.PatchRule{Op: envfile.PatchDelete, Key: d})
	}
	for _, r := range renames {
		old, nw, ok := splitKV(r)
		if !ok {
			return fmt.Errorf("--rename %q: expected OLD=NEW", r)
		}
		rules = append(rules, envfile.PatchRule{Op: envfile.PatchRename, Key: old, NewKey: nw})
	}

	opts := envfile.PatchOptions{IgnoreMissing: ignoreMissing}
	patched, err := envfile.Patch(env, rules, opts)
	if err != nil {
		return fmt.Errorf("patch: %w", err)
	}

	if dryRun {
		for k, v := range patched {
			fmt.Printf("%s=%s\n", k, v)
		}
		return nil
	}
	return envfile.Write(filePath, patched)
}

func splitKV(s string) (string, string, bool) {
	for i := 0; i < len(s); i++ {
		if s[i] == '=' {
			return s[:i], s[i+1:], true
		}
	}
	return "", "", false
}
