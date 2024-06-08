package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

const (
	countColumns = 2
	bufSize      = 1024
)

func batchReadRecords(reader *csv.Reader, records [][]string, bufSize int) ([][]string, error) {
	for i := 0; i < bufSize; i++ {
		record, err := reader.Read()
		if err != nil {
			return records, err
		}

		if len(record) != countColumns {
			errMsg := fmt.Sprintf("Строка содержит неправильное количество колонок: %v. "+
				"Нужно %v колонки формата \"Вопрос,Ответ\", а в строке - %v",
				record, countColumns, len(record))
			return records, errors.New(errMsg)
		}

		records = append(records, record)
	}

	return records, nil
}

func main() {
	fileName := flag.String("file", "problems.csv", "CSV файл с вопросами и ответами")
	shuffle := flag.Bool("shuffle", false, "Перемешать вопросы")
	flag.Parse()

	file, err := os.Open(*fileName)
	if err != nil {
		log.Fatalf("Ошибка открытия файла: %s", err)
	}
	defer func() {
		if err = file.Close(); err != nil {
			log.Fatalf("Ошибка закрытия файла: %s", err)
		}
	}()

	reader := csv.NewReader(file)
	records := make([][]string, 0, bufSize)
	readerStdin := bufio.NewReader(os.Stdin)

	correctAnswers, incorrectAnswers := 0, 0
	showStrStartQuestions := false

	for {
		records, errReadRecords := batchReadRecords(reader, records, bufSize)
		if errReadRecords != nil && errReadRecords != io.EOF {
			log.Fatalf("Ошибка считывания записей файла: %s", errReadRecords)
		}

		if len(records) == 0 {
			log.Fatal("̆В файле нет записей")
		}

		if *shuffle {
			source := rand.NewSource(time.Now().UnixNano())
			r := rand.New(source)
			r.Shuffle(len(records), func(i, j int) {
				records[i], records[j] = records[j], records[i]
			})
		}

		if !showStrStartQuestions {
			fmt.Println("===== Вопросы =====")
			showStrStartQuestions = true
		}

		for _, record := range records {
			fmt.Printf("%s: ", record[0])

			answer, err := readerStdin.ReadString('\n')
			if err != nil {
				fmt.Println("Ошибка считывания ответа:", err)
				return
			}

			if strings.EqualFold(strings.TrimSpace(answer), strings.TrimSpace(record[1])) {
				correctAnswers++
			} else {
				incorrectAnswers++
			}
		}

		fmt.Println("\n===== Результат =====")
		fmt.Printf("Правильных ответов: %d\n", correctAnswers)
		fmt.Printf("Неправильных ответов: %d", incorrectAnswers)

		if errReadRecords == io.EOF {
			break
		}

		records = records[:0]
	}
}
