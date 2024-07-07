package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"math"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// ============== actors ==============

func SimulateUsers(channels []chan UserMessage, numMessages, writeMessagesInterval int) {
	for i := 0; i < numMessages; i++ {
		for user, ch := range channels {
			msg := UserMessage{
				Token:  fmt.Sprintf("token%d", user),
				FileID: fmt.Sprintf("%d", user),
				Data:   fmt.Sprintf("Message %d from user %d", i, user),
			}
			ch <- msg
		}
		time.Sleep(time.Duration(writeMessagesInterval) * time.Millisecond)
	}
}

// ============== queue and token validator ==============

type TokenValidator interface {
	IsValid(token string) bool
}

type SimpleTokenValidator struct {
	ValidTokens map[string]struct{}
}

func NewSimpleTokenValidator(validTokens []string) *SimpleTokenValidator {
	var validTokensMap = make(map[string]struct{})
	for _, token := range validTokens {
		validTokensMap[token] = struct{}{}
	}

	return &SimpleTokenValidator{
		ValidTokens: validTokensMap,
	}
}

func (v *SimpleTokenValidator) IsValid(token string) bool {
	_, ok := v.ValidTokens[token]
	return ok
}

type UserMessage struct {
	Token  string
	FileID string
	Data   string
}

type SimpleMessageQueue struct {
	queue          chan UserMessage
	tokenValidator TokenValidator
}

func NewSimpleMessageQueue(tokenValidator TokenValidator, numMessages int) *SimpleMessageQueue {
	return &SimpleMessageQueue{
		queue:          make(chan UserMessage, numMessages),
		tokenValidator: tokenValidator,
	}
}

func (q *SimpleMessageQueue) GetMessages(messageChannels []chan UserMessage) {
	for _, ch := range messageChannels {
		go func(ch <-chan UserMessage) {
			for {
				msg, ok := <-ch
				if !ok {
					return
				}

				if q.tokenValidator.IsValid(msg.Token) {
					q.queue <- msg
				}
			}
		}(ch)
	}
}

// ============== cache ==============

type MessageCache interface {
	WriteMessages(messages <-chan UserMessage)
	WriteToFiles() error
	RetryWriteToFiles(ctx context.Context, retryTimeInterval time.Duration, retriesWriteToFilesCount uint) error
}

type SimpleMessageCache struct {
	data map[string][]UserMessage
	mu   sync.Mutex
}

func NewSimpleMessageCache() *SimpleMessageCache {
	return &SimpleMessageCache{
		data: make(map[string][]UserMessage),
		mu:   sync.Mutex{},
	}
}

func (c *SimpleMessageCache) WriteMessages(messages <-chan UserMessage) {
	go func() {
		for {
			msg, ok := <-messages
			if !ok {
				return
			}

			c.mu.Lock()
			c.data[msg.FileID] = append(c.data[msg.FileID], msg)
			c.mu.Unlock()
		}
	}()
}

func (c *SimpleMessageCache) WriteToFiles() error {
	c.mu.Lock()

	localMap := c.data
	c.data = make(map[string][]UserMessage)

	c.mu.Unlock()

	if err := os.MkdirAll("files", 0755); err != nil {
		log.Printf("error create dir: %s\n", err)
		return err
	}

	for key, value := range localMap {
		file, err := os.OpenFile("files/"+key+".txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Printf("error opening file: %s\n", err)
			return err
		}
		defer func() {
			if err := file.Close(); err != nil {
				log.Printf("Error closing file: %s\n", err)
			}
		}()

		for _, msg := range value {
			if _, err := file.WriteString(msg.Data + "\n"); err != nil {
				log.Printf("error write string to file: %s\n", err)
				return err
			}
		}
	}

	return nil
}

func (c *SimpleMessageCache) RetryWriteToFiles(ctx context.Context, retryTimeInterval time.Duration,
	retriesWriteToFilesCount uint) error {
	var retryCount uint = 0
	maxRetries := retriesWriteToFilesCount

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			log.Printf("try %d\n", retryCount+1)

			retryCount++
			if retryCount > maxRetries {
				return errors.New("reached max retries, operation is not available")
			}

			if err := c.WriteToFiles(); err != nil {
				log.Printf("Error writing to files: %s\n", err)
				time.Sleep(retryTimeInterval)
				continue
			}

			return nil
		}
	}
}

// ============== worker ==============

type Worker struct {
	cache                                   MessageCache
	interval                                time.Duration
	retriesWriteToFilesCount                uint
	retriesWriteToFilesIntervalMilliseconds time.Duration
}

func NewWorker(cache MessageCache, interval, retriesWriteToFilesIntervalMilliseconds time.Duration,
	retriesWriteToFilesCount uint) *Worker {
	return &Worker{
		cache:                                   cache,
		interval:                                interval,
		retriesWriteToFilesCount:                retriesWriteToFilesCount,
		retriesWriteToFilesIntervalMilliseconds: retriesWriteToFilesIntervalMilliseconds,
	}
}

func (w *Worker) Run(ctx context.Context) {
	ticker := time.NewTicker(w.interval)

	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Println("worker done")
				if err := w.cache.WriteToFiles(); err != nil {
					log.Printf("Error worker done: %s\n", err)
				}
				return
			case <-ticker.C:
				log.Println("Running worker")
				if err := w.cache.RetryWriteToFiles(ctx, w.retriesWriteToFilesIntervalMilliseconds,
					w.retriesWriteToFilesCount); err != nil {
					log.Printf("Error running worker: %s\n", err)
				}
			}
		}
	}()
}

func (w *Worker) GetCorrectWorkerPoolSize(numUsers, numMessages, defaultNum int) int {
	workerPoolSize := defaultNum
	allMessages := numUsers * numMessages
	if allMessages/1000 > workerPoolSize {
		workerPoolSize = int(math.Ceil(float64(allMessages) / 1000))
	}

	return workerPoolSize
}

// ============== main ==============

func main() {
	log.Println("Start app")

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Failed to read config file: %s\n", err)
	}

	workerInterval := viper.GetInt("worker_interval_seconds")
	writeMessagesInterval := viper.GetInt("write_messages_interval_milliseconds")
	numUsers := viper.GetInt("num_users")
	numMessages := viper.GetInt("num_messages")
	retriesWriteToFilesCount := viper.GetUint("retries_write_to_files_count")
	retriesWriteToFilesIntervalMilliseconds := viper.GetInt("retries_write_to_files_interval_milliseconds")

	validTokens := []string{"token1", "token2"}

	tokenValidator := NewSimpleTokenValidator(validTokens)
	messageQueue := NewSimpleMessageQueue(tokenValidator, numMessages)
	cache := NewSimpleMessageCache()
	worker := NewWorker(
		cache,
		time.Duration(workerInterval)*time.Second,
		time.Duration(retriesWriteToFilesIntervalMilliseconds)*time.Millisecond,
		retriesWriteToFilesCount)

	messageChannels := make([]chan UserMessage, 0, numUsers)
	for i := 0; i < numUsers; i++ {
		messageChannels = append(messageChannels, make(chan UserMessage, numMessages))
	}

	go func() {
		SimulateUsers(messageChannels, numMessages, writeMessagesInterval)

		for _, ch := range messageChannels {
			close(ch)
		}
	}()

	messageQueue.GetMessages(messageChannels)
	cache.WriteMessages(messageQueue.queue)

	workerPoolSize := worker.GetCorrectWorkerPoolSize(numUsers, numMessages, 3)
	log.Println("workerPoolSize", workerPoolSize)
	for w := 1; w <= workerPoolSize; w++ {
		worker.Run(ctx)
	}

	<-ctx.Done()

	log.Println("shutting down gracefully")
	time.Sleep(3 * time.Second)
	log.Println("exit app")
}
