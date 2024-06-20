package main

import (
	"reflect"
)

func IsEqualArrays[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}

	mapA := make(map[T]struct{})
	mapB := make(map[T]struct{})

	for _, val := range a {
		mapA[val] = struct{}{}
	}

	for _, val := range b {
		mapB[val] = struct{}{}
	}

	return reflect.DeepEqual(mapA, mapB)
}

func main() {
	arr1 := []int{1, 2, 3, 4, 5}
	arr2 := []int{3, 5, 1, 2, 4}

	result := IsEqualArrays(arr1, arr2)
	println(result)
}
