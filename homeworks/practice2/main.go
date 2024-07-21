package main

import (
	"errors"
	"fmt"
)

func EvalSequence(matrix [][]int, userAnswer []int) int {
	if err := validation(matrix, userAnswer); err != nil {
		fmt.Printf("error: %v\n", err)
		return 0
	}

	maxGrade := calMaxGrade(matrix)
	//fmt.Println("maxGrade", maxGrade)
	userGrade := calcUserGrade(matrix, userAnswer)
	//fmt.Println("userGrade", userGrade)

	percent := userGrade * 100 / maxGrade

	return percent
}

func calcUserGrade(matrix [][]int, userAnswer []int) int {
	var sum int

	for i := 1; i < len(userAnswer); i++ {
		sum += matrix[userAnswer[i-1]][userAnswer[i]]
	}

	return sum
}

func calMaxGrade(matrix [][]int) int {
	var maxGrade int

	for i := range matrix {
		visited := make([]bool, len(matrix))
		path := dfsMaxGrade(i, 0, visited, matrix)

		if path > maxGrade {
			maxGrade = path
		}
	}

	return maxGrade
}

func dfsMaxGrade(vertex int, sum int, visited []bool, adjMatrix [][]int) int {
	visited[vertex] = true
	maxSum := sum

	for i := range adjMatrix {
		if adjMatrix[vertex][i] != 0 && !visited[i] {
			path := dfsMaxGrade(i, sum+adjMatrix[vertex][i], visited, adjMatrix)

			if path > maxSum {
				maxSum = path
			}
		}
	}

	return maxSum
}

func validation(matrix [][]int, userAnswer []int) error {
	matrixLength := len(matrix)
	userAnswerLength := len(userAnswer)
	//fmt.Println(matrixLength, userAnswerLength)

	if matrixLength == 0 {
		return errors.New(ErrorEmptyMatrix)
	}

	if userAnswerLength == 0 {
		return errors.New(ErrorUserAnswer)
	}

	for i := range matrix {
		if len(matrix[i]) != matrixLength {
			return errors.New(ErrorMatrixNotSquare)
		}
	}

	for i := range matrix {
		if matrix[i][i] != 0 {
			return errors.New(ErrorGraphLoop)
		}
	}

	if matrixLength < userAnswerLength {
		return errors.New(ErrorUserAnswersRange)
	}

	for _, node := range userAnswer {
		if node > len(matrix)-1 {
			return errors.New(ErrorUserAnswersIncorrect)
		}
	}

	uniqueAnswers := make(map[int]struct{}, len(userAnswer))
	for i := range userAnswer {
		if _, ok := uniqueAnswers[userAnswer[i]]; ok {
			return errors.New(ErrorUserAnswerNotUnique)
		}
		uniqueAnswers[userAnswer[i]] = struct{}{}
	}

	return nil
}
