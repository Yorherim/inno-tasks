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
	RunE: func(cmd *cobra.Command, args []string) error {
		file := args[0]

		if _, err := os.Create(file); err != nil {
			fmt.Println("Ошибка при создании файла:", err)
			return err
		}

		fmt.Println("Файл", file, "успешно создан")
		return nil
	},
}

var readCmd = &cobra.Command{
	Use:   "read [file]",
	Short: "Читает содержимое файла",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		file := args[0]

		data, err := os.ReadFile(file)
		if err != nil {
			fmt.Println("Ошибка при чтении файла:", err)
			return err
		}

		fmt.Printf("Содержимое файла %s:\n%s", file, string(data))
		return nil
	},
}

var deleteCmd = &cobra.Command{
	Use:   "delete [file]",
	Short: "Удаляет файл",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		file := args[0]

		if err := os.Remove(file); err != nil {
			fmt.Println("Ошибка при удалении файла:", err)
			return err
		}

		fmt.Println("Файл", file, "успешно удален")
		return nil
	},
}
