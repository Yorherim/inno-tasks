package main

import (
	"errors"
	"fmt"
	"log"
)

func main() {
	ch1, err := chIntInit(0, 5)
	if err != nil {
		log.Fatalf("Ошибка: %s", err)
	}

	ch2, err := chIntInit(5, 10)
	if err != nil {
		log.Fatalf("Ошибка: %s", err)
	}

	merged := mergeChannels(ch1, ch2)

	for v := range merged {
		fmt.Println(v)
	}
}

func mergeChannels[T any](ch1, ch2 <-chan T) <-chan T {
	outCh := make(chan T, len(ch1)+len(ch2))
	defer close(outCh)

	for v := range ch1 {
		outCh <- v
	}
	for v := range ch2 {
		outCh <- v
	}

	return outCh
}

func chIntInit(startIndex, endIndex int) (<-chan int, error) {
	if startIndex > endIndex {
		return nil, errors.New("начальное значение больше конечного")
	}

	ch := make(chan int, endIndex-startIndex)
	defer close(ch)

	for i := startIndex; i < endIndex; i++ {
		ch <- i
	}

	return ch, nil
}
