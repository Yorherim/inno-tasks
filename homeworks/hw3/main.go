package main

import (
	"fmt"
	"strings"
)

type Formatter interface {
	Format() string
}

// PlainText Структура для обычного текста
type PlainText struct {
	Text string
}

func (p PlainText) Format() string {
	return p.Text
}

// BoldText Структура для жирного текста
type BoldText struct {
	Text string
}

func (b BoldText) Format() string {
	return "**" + b.Text + "**"
}

// CodeText Структура для кода
type CodeText struct {
	Text string
}

func (c CodeText) Format() string {
	return "`" + c.Text + "`"
}

// ItalicText Структура для курсива
type ItalicText struct {
	Text string
}

func (i ItalicText) Format() string {
	return "_" + i.Text + "_"
}

// ChainFormatter Цепочка модификаторов
type ChainFormatter struct {
	formatters []Formatter
}

func (cf *ChainFormatter) AddFormatter(formatter Formatter) {
	cf.formatters = append(cf.formatters, formatter)
}

func (cf *ChainFormatter) FormatChain() string {
	var result strings.Builder

	for _, formatter := range cf.formatters {
		result.WriteString(formatter.Format())
	}

	return result.String()
}

func (cf *ChainFormatter) FormatChainln() string {
	var result strings.Builder

	for _, formatter := range cf.formatters {
		result.WriteString(formatter.Format() + "\n")
	}

	return result.String()
}

func main() {
	// Пример использования
	plain := PlainText{"Пример обычного текста"}
	bold := BoldText{"Пример жирного текста"}
	code := CodeText{"Пример кода"}
	italic := ItalicText{"Пример текста курсивом"}

	chainFormatter := ChainFormatter{}
	chainFormatter.AddFormatter(plain)
	chainFormatter.AddFormatter(bold)
	chainFormatter.AddFormatter(code)
	chainFormatter.AddFormatter(italic)

	fmt.Println(chainFormatter.FormatChain())
	fmt.Println(chainFormatter.FormatChainln())
}
