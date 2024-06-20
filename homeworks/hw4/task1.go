package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func readInput(inputChan chan string, errChain chan error) {
	reader := bufio.NewReader(os.Stdin)
	for {
		input, err := reader.ReadString('\n')
		if err != nil {
			errChain <- err
		}

		inputChan <- input
	}
}

func writeToFile(inputChan chan string, errChain chan error) {
	file, err := os.Create("./output.txt")
	if err != nil {
		errChain <- err
	}
	defer func() {
		if err = file.Close(); err != nil {
			errChain <- err
		}
	}()

	for {
		select {
		case input := <-inputChan:
			if _, err := file.WriteString(input); err != nil {
				errChain <- err
			}
		}
	}
}

func main() {
	fmt.Println("Старт работы программы")

	inputChan := make(chan string)
	defer close(inputChan)

	errChain := make(chan error)
	defer close(errChain)

	go func() {
		for {
			select {
			case err := <-errChain:
				log.Fatalf("Произошла ошибка: %s", err)
			}
		}
	}()

	go readInput(inputChan, errChain)
	go writeToFile(inputChan, errChain)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	fmt.Println("Завершение работы программы")

	os.Exit(0)
}
