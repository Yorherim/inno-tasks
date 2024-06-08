package main

import (
	"github.com/spf13/cobra"
	"log"
)

func main() {
	rootCmd.AddCommand(createCmd, readCmd, deleteCmd)
	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Не удалось выполнить команду: %s", err)
	}
}
