package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "filetool",
	Short: "Утилита для основных файловых операций",
}

var createCmd = &cobra.Command{
	Use:   "create [file]",
	Short: "Создает новый файл",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		file := args[0]

		if _, err := os.Create(file); err != nil {
			fmt.Println("Ошибка при создании файла:", err)
			return
		}

		fmt.Println("Файл", file, "успешно создан")
	},
}

var readCmd = &cobra.Command{
	Use:   "read [file]",
	Short: "Читает содержимое файла",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		file := args[0]

		data, err := os.ReadFile(file)
		if err != nil {
			fmt.Println("Ошибка при чтении файла:", err)
			return
		}

		fmt.Printf("Содержимое файла %s:\n%s", file, string(data))
	},
}

var deleteCmd = &cobra.Command{
	Use:   "delete [file]",
	Short: "Удаляет файл",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		file := args[0]

		if err := os.Remove(file); err != nil {
			fmt.Println("Ошибка при удалении файла:", err)
			return
		}

		fmt.Println("Файл", file, "успешно удален")
	},
}
