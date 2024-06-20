package main

import (
	"errors"
	"fmt"
	"log"
	"reflect"
)

type Numbers[T ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~float32 | ~float64] []T

func (n *Numbers[T]) Sum() (T, error) {
	if len(*n) == 0 {
		return 0, errors.New("Slice is empty")
	}

	sum := (*n)[0]
	for i := 1; i < len(*n); i++ {
		sum += (*n)[i]
	}
	return sum, nil
}

func (n *Numbers[T]) Product() (T, error) {
	if len(*n) == 0 {
		return 0, errors.New("Slice is empty")
	}

	product := (*n)[0]
	for i := 1; i < len(*n); i++ {
		product *= (*n)[i]
	}
	return product, nil
}

func (n *Numbers[T]) Equal(other Numbers[T]) bool {
	return reflect.DeepEqual(n, other)
}

func (n *Numbers[T]) FindIndex(value T) (int, error) {
	for i, v := range *n {
		if v == value {
			return i, nil
		}
	}
	return -1, errors.New("Value not found")
}

func (n *Numbers[T]) RemoveByValue(value T) {
	for i, v := range *n {
		if v == value {
			*n = append((*n)[:i], (*n)[i+1:]...)
			return
		}
	}
}

func (n *Numbers[T]) RemoveByIndex(index int) {
	*n = append((*n)[:index], (*n)[index+1:]...)
}

func main() {
	nums := Numbers[int]{1, 2, 3, 4, 5}

	sum, err := nums.Sum()
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
	fmt.Println("Sum:", sum)

	product, err := nums.Product()
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
	fmt.Println("Product:", product)

	equals := nums.Equal(Numbers[int]{1, 2, 3, 4, 5})
	fmt.Println("Are equal:", equals)

	index, err := nums.FindIndex(3)
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
	fmt.Println("Index of 3:", index)

	nums.RemoveByValue(3)
	fmt.Println("After removing 3:", nums)

	nums.RemoveByIndex(1)
	fmt.Println("After removing element at index 1:", nums)
}
