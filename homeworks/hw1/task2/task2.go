package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

func main() {
	rootCmd.AddCommand(createCmd, readCmd, deleteCmd)
	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
