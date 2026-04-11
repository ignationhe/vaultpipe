package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"vaultpipe/internal/envfile"
)

func init() {
	var fileA, fileB string
	cmd := &cobra.Command{
		Use:   "compare",
		Short: "Compare two .env files and show differences",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCompare(fileA, fileB)
		},
	}
	cmd.Flags().StringVarP(&fileA, "file-a", "a", ".env", "first .env file")
	cmd.Flags().StringVarP(&fileB, "file-b", "b", ".env.new", "second .env file")
	rootCmd.AddCommand(cmd)
}

func runCompare(pathA, pathB string) error {
	a, err := envfile.Parse(pathA)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("reading %s: %w", pathA, err)
	}
	b, err := envfile.Parse(pathB)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("reading %s: %w", pathB, err)
	}
	r := envfile.Compare(a, b)
	if !r.HasDifferences() {
		fmt.Println("No differences found.")
		return nil
	}
	for _, k := range r.SortedKeys() {
		if v, ok := r.OnlyInA[k]; ok {
			fmt.Printf("- %s=%s\n", k, v)
			continue
		}
		if v, ok := r.OnlyInB[k]; ok {
			fmt.Printf("+ %s=%s\n", k, v)
			continue
		}
		if pair, ok := r.Different[k]; ok {
			fmt.Printf("~ %s: %s -> %s\n", k, pair[0], pair[1])
		}
	}
	return nil
}
