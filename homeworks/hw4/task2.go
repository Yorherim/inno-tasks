package main

import (
	"fmt"
	"sync"
)

func main() {
	nums := []int{2, 3, 4, 5, 6, 7, 8, 9, 10}
	primeNums := make([]int, 0)
	compositeNums := make([]int, 0)

	primeChan := goNums(nums, true)
	compositeChan := goNums(nums, false)

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case num, ok := <-primeChan:
				if !ok {
					return
				}
				primeNums = append(primeNums, num)
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case num, ok := <-compositeChan:
				if !ok {
					return
				}
				compositeNums = append(compositeNums, num)
			}
		}
	}()

	wg.Wait()

	fmt.Println(primeNums)
	fmt.Println(compositeNums)
}

func goNums(nums []int, prime bool) chan int {
	ch := make(chan int)
	go func() {
		defer close(ch)

		for _, num := range nums {
			if isPrime(num) == prime {
				ch <- num
			}
		}
	}()
	return ch
}

func isPrime(num int) bool {
	if num < 2 {
		return false
	}
	for i := 2; i*i <= num; i++ {
		if num%i == 0 {
			return false
		}
	}
	return true
}
