package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

type Exam struct {
	Students []Student `json:"students"`
	Objects  []Object  `json:"objects"`
	Results  []Result  `json:"results"`
}
type Student struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Grade int    `json:"grade"`
}
type Object struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}
type Result struct {
	ObjectId  int `json:"object_id"`
	StudentId int `json:"student_id"`
	Result    int `json:"result"`
}

func main() {
	exam, err := parseJson("dz3.json")
	if err != nil {
		log.Fatalf("error parse json: %s", err)
	}
	if len(exam.Students) == 0 {
		fmt.Println("Файл не содержит данных")
		return
	}

	mapObjects := make(map[int]Object, len(exam.Objects))
	for _, obj := range exam.Objects {
		mapObjects[obj.Id] = obj
	}

	mapStudents := make(map[int]Student, len(exam.Students))
	for _, s := range exam.Students {
		mapStudents[s.Id] = s
	}

	fmt.Printf("%-44s\n", strings.Repeat("-", 44))
	fmt.Printf(" %-12s | %-4s | %-9s | %-9s \n", "Student name", "Grade", "Object", "Result")
	fmt.Printf("%-44s\n", strings.Repeat("-", 44))

	excellentResultsStudents := make(map[string]int)
	filteredResults := Filter(exam.Results, func(r Result) bool {
		if r.Result == 5 {
			excellentResultsStudents[mapStudents[r.StudentId].Name]++
			return true
		}
		return false
	})

	for _, r := range filteredResults {
		obj, ok := mapObjects[r.ObjectId]
		if !ok {
			fmt.Printf("Object %v not found in mapObjects", r.ObjectId)
			return
		}

		student, ok := mapStudents[r.StudentId]
		if !ok {
			fmt.Printf("Student %v not found in mapStudents", r.StudentId)
			return
		}

		if excellentResultsStudents[student.Name] == len(exam.Objects) {
			fmt.Printf(" %-12s | %-5v | %-9s | %-9v \n", student.Name, student.Grade, obj.Name, r.Result)
		}
	}
}

func parseJson(filePath string) (*Exam, error) {
	file, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Ошибка открытия файла", err)
		return nil, err
	}

	var exam Exam
	if err := json.Unmarshal(file, &exam); err != nil {
		fmt.Printf("Ошибка Unmarshal: %s", err)
		return nil, err
	}

	return &exam, nil
}

func Filter[T any](s []T, f func(T) bool) []T {
	var r []T
	for _, v := range s {
		if f(v) {
			r = append(r, v)
		}
	}
	return r
}
