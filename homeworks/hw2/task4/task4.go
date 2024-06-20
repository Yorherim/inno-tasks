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
type ObjectStat struct {
	SumAllResultsObject    int
	SumAllStudentsInObject int
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

	objGradeResult := make(map[int]map[int]ObjectStat)

	for _, r := range exam.Results {
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

		if _, ok := objGradeResult[obj.Id]; !ok {
			objGradeResult[obj.Id] = make(map[int]ObjectStat)
		}

		if _, ok := objGradeResult[obj.Id][student.Grade]; !ok {
			objGradeResult[obj.Id][student.Grade] = ObjectStat{}
		}

		currentObjStat := objGradeResult[obj.Id][student.Grade]
		currentObjStat.SumAllResultsObject += r.Result
		currentObjStat.SumAllStudentsInObject++
		objGradeResult[obj.Id][student.Grade] = currentObjStat
	}

	for objectId, grades := range objGradeResult {
		obj, ok := mapObjects[objectId]
		if !ok {
			fmt.Printf("Object %v not found in mapObjects", objectId)
			return
		}

		const countRepeat = 20
		var sumResults, sumStudents int

		fmt.Printf("%14s\n", strings.Repeat("-", countRepeat))
		fmt.Printf(" %-9s | %-5s\n", obj.Name, "Mean")
		fmt.Printf("%14s\n", strings.Repeat("-", countRepeat))

		for key, objStat := range grades {
			mean := float64(objStat.SumAllResultsObject) / float64(objStat.SumAllStudentsInObject)

			fmt.Printf(" %-3d grade | %.1f\n", key, mean)

			sumResults += objStat.SumAllResultsObject
			sumStudents += objStat.SumAllStudentsInObject
		}

		allMean := float64(sumResults) / float64(sumStudents)

		fmt.Printf("%14s\n", strings.Repeat("-", countRepeat))
		fmt.Printf(" %-9s | %-5.1f\n", "mean", allMean)
		fmt.Printf("%14s\n", strings.Repeat("-", countRepeat))
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
