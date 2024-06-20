# Инструкция по запуску программы:

в `main.go` файле должна быть задействована строка `task1.Task1()` из пакета
`github.com/Yorherim/inno-tasks/homeworks/hw1/task1`  
Запуск приложения:
- без флагов
```
go run main.go
```

- флаг `shuffle` для рандомного порядка вопросов
```
go run main.go -shuffle
```

- флаг `file` для указания файла формата `.csv` c вопросами
```
go run main.go -file=someTitle.csv
```