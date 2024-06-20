package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	inputCh := make(chan string)

	quitCh := make(chan os.Signal, 1)
	signal.Notify(quitCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := readInput(ctx, inputCh); err != nil {
			log.Fatalf("Произошла ошибка при чтении из консоли: %s", err)
		}
	}()
	go func() {
		if err := writeToFile(ctx, inputCh); err != nil {
			log.Fatalf("Произошла ошибка при записи в файл: %s", err)
		}
	}()

	<-quitCh
	cancel()
	fmt.Println("\nЗавершение программы...")
}

func readInput(ctx context.Context, inputCh chan<- string) error {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Введите значения:")
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			if scanner.Scan() {
				inputCh <- scanner.Text()
			} else if scanner.Err() != nil {
				return scanner.Err()
			}
		}
	}
}

func writeToFile(ctx context.Context, inputCh <-chan string) error {
	file, err := os.Create("input.txt")
	if err != nil {
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Println("Ошибка при закрытии файла:", err)
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return nil
		case input, ok := <-inputCh:
			if !ok {
				return nil
			}
			if _, err := file.WriteString(input + "\n"); err != nil {
				return err
			}
		}
	}
}
