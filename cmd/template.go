package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourorg/vaultpipe/internal/envfile"
)

var (
	tplSrcFile  string
	tplDstFile  string
	tplVarFiles []string
	tplMissing  string
)

func init() {
	templateCmd := &cobra.Command{
		Use:   "template",
		Short: "Render a .env template by substituting ${KEY} placeholders",
		RunE:  runTemplate,
	}

	templateCmd.Flags().StringVarP(&tplSrcFile, "src", "s", "", "Source template file (required)")
	templateCmd.Flags().StringVarP(&tplDstFile, "dst", "d", "", "Destination output file (required)")
	templateCmd.Flags().StringArrayVarP(&tplVarFiles, "vars", "v", nil, "One or more .env files supplying variable values")
	templateCmd.Flags().StringVar(&tplMissing, "missing-key", "error", "Behaviour for missing keys: error|keep|empty")

	_ = templateCmd.MarkFlagRequired("src")
	_ = templateCmd.MarkFlagRequired("dst")

	rootCmd.AddCommand(templateCmd)
}

func runTemplate(cmd *cobra.Command, _ []string) error {
	vars := map[string]string{}

	// Merge all supplied var files (later files win).
	for _, vf := range tplVarFiles {
		parsed, err := envfile.Parse(vf)
		if err != nil {
			return fmt.Errorf("reading var file %s: %w", vf, err)
		}
		for k, v := range parsed {
			vars[k] = v
		}
	}

	// Also honour variables already set in the process environment.
	for _, kv := range os.Environ() {
		for i := 0; i < len(kv); i++ {
			if kv[i] == '=' {
				k, v := kv[:i], kv[i+1:]
				if _, exists := vars[k]; !exists {
					vars[k] = v
				}
				break
			}
		}
	}

	opts := envfile.TemplateOptions{MissingKey: tplMissing}
	if err := envfile.RenderFile(tplSrcFile, tplDstFile, vars, opts); err != nil {
		return err
	}

	cmd.Printf("rendered %s → %s\n", tplSrcFile, tplDstFile)
	return nil
}
