package main

import (
	"errors"
	"fmt"
	"sort"
)

func findIntersection(slices ...[]int) ([]int, error) {
	if len(slices) == 0 {
		return []int{}, errors.New("передано ноль слайсов")
	}
	if len(slices[0]) == 0 {
		return []int{}, nil
	}

	intersection := make(map[int]struct{}, len(slices[0]))

	for _, num := range slices[0] {
		intersection[num] = struct{}{}
	}

	for _, slice := range slices[1:] {
		tempMap := make(map[int]struct{}, len(intersection))

		for _, num := range slice {
			if _, ok := intersection[num]; ok {
				tempMap[num] = struct{}{}
			}
		}

		intersection = tempMap
	}

	result := make([]int, 0, len(intersection))

	for num := range intersection {
		result = append(result, num)
	}

	sort.Ints(result)
	return result, nil
}

func main() {
	fmt.Println(findIntersection([]int{1, 2, 3, 4}, []int{3, 2}))
	fmt.Println(findIntersection([]int{1, 2, 3, 2}))
	fmt.Println(findIntersection([]int{1, 2, 3, 4}, []int{3, 2}, []int{}))
}
