package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	WorkerPoolSize = 3
	WorkerTimeout  = 2 * time.Second

	DefaultExtensionForFiles = ".jpg"
)

var urls = []string{
	"https://via.placeholder.com/150/92c952",
	"https://via.placeholder.com/150/771796",
	"https://via.placeholder.com/150/24f355",
	"https://via.placeholder.com/150/d32776",
	"https://via.placeholder.com/150/f66b97",
}

func getFileNameFromUrl(url string) string {
	segments := strings.Split(url, "/")
	fileName := segments[len(segments)-1]

	if !strings.Contains(fileName, ".") {
		return fileName + DefaultExtensionForFiles
	}

	return fileName
}

func downloadFile(url string, wg *sync.WaitGroup, results chan<- string) error {
	defer wg.Done()

	client := http.Client{
		Timeout: WorkerTimeout,
	}

	resp, err := client.Get(url)
	if err != nil {
		results <- fmt.Sprintf("Failed to download %s: %v", url, err)
		return err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Println("Error closing resp.Body")
		}
	}()

	file, err := os.Create("files/file_" + getFileNameFromUrl(url))
	if err != nil {
		results <- fmt.Sprintf("Failed to create file for %s: %v", url, err)
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Println("Error closing file")
		}
	}()

	if _, err = io.Copy(file, resp.Body); err != nil {
		results <- fmt.Sprintf("Failed to write file for %s: %v", url, err)
		return err
	}

	results <- fmt.Sprintf("File downloaded from url %s", url)
	return nil
}

func worker(jobs <-chan string, results chan<- string, wg *sync.WaitGroup) {
	for url := range jobs {
		if err := downloadFile(url, wg, results); err != nil {
			log.Println("File download error", err)
			continue
		}
	}
}

func main() {
	if err := os.MkdirAll("files", 0755); err != nil {
		log.Fatalf("Failed to create directory for files: %v", err)
	}

	jobs := make(chan string, 3)
	results := make(chan string, 3)

	wg := &sync.WaitGroup{}
	wg.Add(len(urls))

	for w := 0; w < WorkerPoolSize; w++ {
		go worker(jobs, results, wg)
	}

	for _, url := range urls {
		jobs <- url
	}
	close(jobs)

	go func() {
		wg.Wait()
		close(results)
	}()

	for result := range results {
		fmt.Println(result)
	}
}
