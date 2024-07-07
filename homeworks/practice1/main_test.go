package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"
)

// Сценарий 1: Успешная запись
func TestSuccessfulMessageWrite(t *testing.T) {
	validTokens := []string{"token1", "token2"}
	tokenValidator := NewSimpleTokenValidator(validTokens)
	numMessages := 5 // Кол-во сообщений для теста

	cache := NewSimpleMessageCache()
	worker := NewWorker(cache, 1*time.Millisecond, 10*time.Millisecond, 3)

	messageChannels := make([]chan UserMessage, 0, 1)
	messageChannel := make(chan UserMessage, numMessages)
	messageChannels = append(messageChannels, messageChannel)

	// Генерация сообщений для теста
	for i := 0; i < numMessages; i++ {
		msg := UserMessage{
			Token:  "token1", // Правильный токен для успешной записи
			FileID: "test_file",
			Data:   fmt.Sprintf("Test message %d", i),
		}
		messageChannel <- msg
	}
	close(messageChannel)

	messageQueue := NewSimpleMessageQueue(tokenValidator, numMessages)
	messageQueue.GetMessages(messageChannels)
	cache.WriteMessages(messageQueue.queue)

	// ждем какое-то время, пока не прочитаются сообщения и не запишутся в кеш
	// в горутине в cache.WriteMessages
	time.Sleep(10 * time.Millisecond)

	if len(cache.data["test_file"]) != numMessages {
		t.Errorf("Expected %d messages in cache, got %d", numMessages, len(cache.data["test_file"]))
	}

	for w := 1; w <= 3; w++ {
		worker.Run(context.Background())
	}

	// ждем какое-то время, пока не сработают воркеры
	time.Sleep(50 * time.Millisecond)

	// Убедимся, что данные в кеше были очищены
	if len(cache.data["test_file"]) != 0 {
		t.Errorf("Expected cache to be empty, but found %d messages", len(cache.data["test_file"]))
	}

	// Проверяем, что сообщение было записано в файл
	files, err := os.ReadDir("files")
	if err != nil {
		t.Fatalf("Error reading files directory: %v", err)
	}

	if len(files) == 0 {
		t.Fatalf("No files were written by the worker")
	}

	foundTestFile := false
	for _, file := range files {
		if file.Name() == "test_file.txt" {
			foundTestFile = true

			// Проверка, что файл не пустой
			fileInfo, err := file.Info()
			if err != nil {
				t.Fatalf("Error getting file info for test_file: %v", err)
			}
			if fileInfo.Size() == 0 {
				t.Fatalf("File test_file is empty")
			}

			break
		}
	}

	if !foundTestFile {
		t.Fatalf("File test_file.txt not found in directory")
	}

	if err := os.Remove("files/test_file.txt"); err != nil {
		t.Fatalf("Error removing test_file.txt: %v", err)
	}
}

// Сценарий 2: Неверный токен
func TestInvalidToken(t *testing.T) {
	validTokens := []string{"token1"}
	tokenValidator := NewSimpleTokenValidator(validTokens)
	numMessages := 5

	cache := NewSimpleMessageCache()
	//worker := NewWorker(cache, 1*time.Second, 3*time.Second, 3)

	messageChannels := make([]chan UserMessage, 0, 1)
	messageChannel := make(chan UserMessage, numMessages)
	messageChannels = append(messageChannels, messageChannel)

	// Генерация сообщений для теста
	for i := 0; i < numMessages; i++ {
		msg := UserMessage{
			Token:  fmt.Sprintf("token%d", i),
			FileID: "test_file",
			Data:   fmt.Sprintf("Test message %d", i),
		}
		messageChannel <- msg
	}
	close(messageChannel)

	messageQueue := NewSimpleMessageQueue(tokenValidator, numMessages)
	messageQueue.GetMessages(messageChannels)
	cache.WriteMessages(messageQueue.queue)

	// ждем какое-то время, пока не прочитаются сообщения и не запишутся в кеш
	// в горутине в cache.WriteMessages
	time.Sleep(100 * time.Millisecond)

	if len(cache.data["test_file"]) != 1 {
		t.Errorf("Expected 1 message in cache, got %d", len(cache.data["test_file"]))
	}
}

// Сценарий 3: Остановка приложения (Graceful Shutdown)
func TestGracefulShutdown(t *testing.T) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	validTokens := []string{"token1", "token2"}
	tokenValidator := NewSimpleTokenValidator(validTokens)
	numMessages := 5 // Кол-во сообщений для теста

	cache := NewSimpleMessageCache()
	worker := NewWorker(cache, 5*time.Second, 3*time.Second, 3)

	messageChannels := make([]chan UserMessage, 0, 1)
	messageChannel := make(chan UserMessage, numMessages)
	messageChannels = append(messageChannels, messageChannel)

	// Генерация сообщений для теста
	for i := 0; i < numMessages; i++ {
		msg := UserMessage{
			Token:  "token1", // Правильный токен для успешной записи
			FileID: "test_file",
			Data:   fmt.Sprintf("Test message %d", i),
		}
		messageChannel <- msg
	}
	close(messageChannel)

	messageQueue := NewSimpleMessageQueue(tokenValidator, numMessages)
	messageQueue.GetMessages(messageChannels)
	cache.WriteMessages(messageQueue.queue)

	for w := 1; w <= 3; w++ {
		worker.Run(ctx)
	}

	time.Sleep(100 * time.Millisecond)

	stop()
	<-ctx.Done()
	log.Println("shutting down gracefully")
	time.Sleep(1 * time.Second)

	// Проверяем, что кеш теперь пустой
	for key := range cache.data {
		if len(cache.data[key]) != 0 {
			t.Errorf("Expected cache for file %s to be empty, but found %d messages", key, len(cache.data[key]))
		}
	}

	// Проверяем, что файл был создан
	files, err := os.ReadDir("files")
	if err != nil {
		t.Fatalf("Error reading files directory: %v", err)
	}

	foundTestFile := false
	for _, file := range files {
		if file.Name() == "test_file.txt" {
			foundTestFile = true

			// Проверка, что файл не пустой
			fileInfo, err := file.Info()
			if err != nil {
				t.Fatalf("Error getting file info for test_file: %v", err)
			}
			if fileInfo.Size() == 0 {
				t.Fatalf("File test_file is empty")
			}

			break
		}
	}

	if !foundTestFile {
		t.Fatalf("File test_file.txt not found in directory")
	}

	if err := os.Remove("files/test_file.txt"); err != nil {
		t.Fatalf("Error removing test_file.txt: %v", err)
	}
}

// Сценарий 4: Высокая нагрузка
func TestHighLoadWorkers(t *testing.T) {
	numMessages := 1000
	numUsers := 5

	cache := NewSimpleMessageCache()
	worker := NewWorker(cache, 5*time.Second, 3*time.Second, 3)

	workerPoolSize := worker.GetCorrectWorkerPoolSize(numUsers, numMessages, 3)

	if workerPoolSize <= 3 {
		t.Fatalf("workerPoolSize must be greater than 3, got")
	}
}

// Cценарий 5: Файл с одновременной записью
func TestConcurrentWritingToFile(t *testing.T) {
	validTokens := []string{"token1", "token2"}
	tokenValidator := NewSimpleTokenValidator(validTokens)
	numMessages := 5
	numUsers := 5

	cache := NewSimpleMessageCache()
	worker := NewWorker(cache, 1*time.Millisecond, 10*time.Millisecond, 3)

	messageChannels := make([]chan UserMessage, 0, numUsers)
	for i := 0; i < numUsers; i++ {
		messageChannel := make(chan UserMessage, numMessages)
		messageChannels = append(messageChannels, messageChannel)
	}

	// Генерация сообщений для теста
	for user, ch := range messageChannels {
		for i := 0; i < numMessages; i++ {
			msg := UserMessage{
				Token:  "token1", // Правильный токен для успешной записи
				FileID: "test_file",
				Data:   fmt.Sprintf("Test message %d from user %d", i, user),
			}
			ch <- msg
		}
		close(ch)
	}

	messageQueue := NewSimpleMessageQueue(tokenValidator, numMessages)
	messageQueue.GetMessages(messageChannels)
	cache.WriteMessages(messageQueue.queue)

	// ждем какое-то время, пока не прочитаются сообщения и не запишутся в кеш
	// в горутине в cache.WriteMessages
	time.Sleep(10 * time.Millisecond)

	if len(cache.data["test_file"]) != numMessages*numUsers {
		t.Errorf("Expected %d messages in cache, got %d", numMessages, len(cache.data["test_file"]))
	}

	for w := 1; w <= 3; w++ {
		worker.Run(context.Background())
	}

	// ждем какое-то время, пока не сработают воркеры
	time.Sleep(50 * time.Millisecond)

	// Убедимся, что данные в кеше были очищены
	if len(cache.data["test_file"]) != 0 {
		t.Errorf("Expected cache to be empty, but found %d messages", len(cache.data["test_file"]))
	}

	// Проверяем, что сообщение было записано в файл
	files, err := os.ReadDir("files")
	if err != nil {
		t.Fatalf("Error reading files directory: %v", err)
	}

	if len(files) == 0 {
		t.Fatalf("No files were written by the worker")
	}

	foundTestFile := false
	for _, file := range files {
		if file.Name() == "test_file.txt" {
			foundTestFile = true

			file, err := os.Open("files/test_file.txt")
			if err != nil {
				t.Fatalf("error open file")
			}
			defer func() {
				if err = file.Close(); err != nil {
					t.Fatalf("error close file")
				}
				if err := os.Remove("files/test_file.txt"); err != nil {
					t.Fatalf("Error removing test_file.txt: %v", err)
				}
			}()

			scanner := bufio.NewScanner(file)
			lineCount := 0
			for scanner.Scan() {
				lineCount++
			}

			if err := scanner.Err(); err != nil {
				t.Fatalf("error scna file")
			}

			if lineCount != numMessages*numUsers {
				t.Fatalf("Expected %d lines in file, got %d", numMessages*numUsers, lineCount)
			}

			break
		}
	}

	if !foundTestFile {
		t.Fatalf("File test_file.txt not found in directory")
	}
}

// Cценарий 6: Сбой работы воркера
func TestRetryWriteFile(t *testing.T) {
	validTokens := []string{"token1", "token2"}
	tokenValidator := NewSimpleTokenValidator(validTokens)
	numMessages := 5
	numUsers := 5

	cache := NewSimpleMessageCache()

	messageChannels := make([]chan UserMessage, 0, numUsers)
	for i := 0; i < numUsers; i++ {
		messageChannel := make(chan UserMessage, numMessages)
		messageChannels = append(messageChannels, messageChannel)
	}

	// Генерация сообщений для теста
	for user, ch := range messageChannels {
		for i := 0; i < numMessages; i++ {
			msg := UserMessage{
				Token:  "token1",
				FileID: ">>>",
				Data:   fmt.Sprintf("Test message %d from user %d", i, user),
			}
			ch <- msg
		}
		close(ch)
	}

	messageQueue := NewSimpleMessageQueue(tokenValidator, numMessages)
	messageQueue.GetMessages(messageChannels)
	cache.WriteMessages(messageQueue.queue)

	// ждем какое-то время, пока не прочитаются сообщения и не запишутся в кеш
	// в горутине в cache.WriteMessages
	time.Sleep(100 * time.Millisecond)

	if err := cache.RetryWriteToFiles(context.Background(), 10*time.Millisecond, 3); err != nil {
		if err.Error() != "reached max retries, operation is not available" {
			t.Errorf("Expected error reached max retries, got %v", err)
		}
	}

	// Убедимся, что данные в кеше остались
	if len(cache.data[">>>"]) != numUsers*numMessages {
		t.Errorf("Expected %d messages in cache, got %d", numUsers*numMessages, len(cache.data[">>>"]))
	}
}
