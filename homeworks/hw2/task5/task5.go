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

var cacheObjects Cache[int, Object]
var cacheStudents Cache[int, Student]

func main() {
	exam, err := parseJson("dz3.json")
	if err != nil {
		log.Fatalf("error parse json: %s", err)
	}
	if len(exam.Students) == 0 {
		fmt.Println("Файл не содержит данных")
		return
	}

	cacheObjects.Init()
	for _, obj := range exam.Objects {
		cacheObjects.Set(obj.Id, obj)
	}

	cacheStudents.Init()
	for _, s := range exam.Students {
		cacheStudents.Set(s.Id, s)
	}

	fmt.Printf("%-44s\n", strings.Repeat("-", 44))
	fmt.Printf(" %-12s | %-4s | %-9s | %-9s \n", "Student name", "Grade", "Object", "Result")
	fmt.Printf("%-44s\n", strings.Repeat("-", 44))

	for _, r := range exam.Results {
		obj, ok := cacheObjects.Get(r.ObjectId)
		if !ok {
			fmt.Printf("Object %v not found in mapObjects", r.ObjectId)
			return
		}

		student, ok := cacheStudents.Get(r.StudentId)
		if !ok {
			fmt.Printf("Student %v not found in mapStudents", r.StudentId)
			return
		}

		fmt.Printf(" %-12s | %-5v | %-9s | %-9v \n", student.Name, student.Grade, obj.Name, r.Result)
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

type Cache[K comparable, V any] struct {
	m map[K]V
}

func (c *Cache[K, V]) Init() {
	c.m = make(map[K]V)
}

func (c *Cache[K, V]) Set(key K, value V) {
	c.m[key] = value
}

func (c *Cache[K, V]) Get(key K) (V, bool) {
	k, ok := c.m[key]
	return k, ok
}
