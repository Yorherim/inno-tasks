package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

func main() {
	fileName := flag.String("file", "problems.csv", "CSV файл с вопросами и ответами")
	shuffle := flag.Bool("shuffle", false, "Перемешать вопросы")
	flag.Parse()

	if !strings.Contains(*fileName, ".csv") {
		fmt.Println("Файл должен иметь расширение .csv")
		return
	}

	file, err := os.Open(*fileName)
	if err != nil {
		fmt.Println("Ошибка открытия файла:", err)
		return
	}
	defer file.Close()

	reader := csv.NewReader(bufio.NewReader(file))
	records := make([][]string, 0)

	for {
		record, err := reader.Read()
		if err != nil {
			break
		}

		records = append(records, record)
	}

	if *shuffle {
		source := rand.NewSource(time.Now().UnixNano())
		r := rand.New(source)
		r.Shuffle(len(records), func(i, j int) {
			records[i], records[j] = records[j], records[i]
		})
	}

	correct, incorrect := 0, 0
	readerStdin := bufio.NewReader(os.Stdin)

	fmt.Println("===== Вопросы =====")
	for _, record := range records {
		fmt.Printf("%s: ", record[0])

		answer, err := readerStdin.ReadString('\n')
		if err != nil {
			fmt.Println("Ошибка считывания ответа:", err)
			return
		}

		if strings.EqualFold(strings.TrimSpace(answer), strings.TrimSpace(record[1])) {
			correct++
		} else {
			incorrect++
		}
	}

	fmt.Println("\n===== Результат =====")
	fmt.Printf("Правильных ответов: %d\n", correct)
	fmt.Printf("Неправильных ответов: %d", incorrect)
}
