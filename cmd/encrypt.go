package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultpipe/internal/envfile"
)

func init() {
	encryptCmd := newEncryptCmd()
	RootCmd.AddCommand(encryptCmd)
}

func newEncryptCmd() *cobra.Command {
	var (
		inputFile  string
		outputFile string
		decrypt    bool
	)

	cmd := &cobra.Command{
		Use:   "encrypt",
		Short: "Encrypt or decrypt values in a .env file",
		Long: `Encrypt reads a .env file, encrypts every value with AES-256-GCM
using a passphrase sourced from the VAULTPIPE_PASSPHRASE environment variable,
and writes the result to the output file.

Pass --decrypt to reverse the operation.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runEncrypt(inputFile, outputFile, decrypt)
		},
	}

	cmd.Flags().StringVarP(&inputFile, "input", "i", ".env", "source .env file")
	cmd.Flags().StringVarP(&outputFile, "output", "o", "", "destination file (defaults to input file)")
	cmd.Flags().BoolVar(&decrypt, "decrypt", false, "decrypt values instead of encrypting")

	return cmd
}

func runEncrypt(inputFile, outputFile string, decrypt bool) error {
	passphrase := os.Getenv("VAULTPIPE_PASSPHRASE")
	if passphrase == "" {
		return fmt.Errorf("VAULTPIPE_PASSPHRASE environment variable is not set")
	}

	secrets, err := envfile.Parse(inputFile)
	if err != nil {
		return fmt.Errorf("parse %q: %w", inputFile, err)
	}

	var result map[string]string
	if decrypt {
		result, err = envfile.Decrypt(secrets, passphrase)
		if err != nil {
			return fmt.Errorf("decrypt: %w", err)
		}
	} else {
		result, err = envfile.Encrypt(secrets, passphrase)
		if err != nil {
			return fmt.Errorf("encrypt: %w", err)
		}
	}

	dest := outputFile
	if dest == "" {
		dest = inputFile
	}

	if err := envfile.Write(dest, result); err != nil {
		return fmt.Errorf("write %q: %w", dest, err)
	}

	action := "Encrypted"
	if decrypt {
		action = "Decrypted"
	}
	fmt.Printf("%s %d key(s) → %s\n", action, len(result), dest)
	return nil
}
