package main

import (
	"fmt"
	"time"

	"github.com/alexeyco/simpletable"
	"github.com/gosuri/uilive"
)

//maybe add lines to a slice, and on each print, iterate and write each one out
//could try to use the table inside of the writer, or could just append new writer for each position
//each writer should probably track the last thing they wrote in order to prevent lines that aren't updated from being deleted

//func (w *ScreenWriter) Write(items [int]string){] //if items[0] is empty then do not overwrite w.Writers[0] or something

func main() {
	var data = [][]interface{}{
		{1, "Newton G. Goetz", 532.7},
		{2, "Rebecca R. Edney", 1423.25},
		{3, "John R. Jackson", 7526.12},
		{4, "Ron J. Gomes", 123.84},
		{5, "Penny R. Lewis", 3221.11},
	}
	table := getDummyTable(data)

	writer := NewScreenWriter()
	writer.Writer.Start()
	for i := 0; i <= 50; i++ {
		fmt.Fprintf(writer.Writer, table.String())
		fmt.Fprintf(writer.Writer.Newline(), "\n")
		fmt.Fprintf(writer.Writer.Newline(), "this is line 2\n")
		time.Sleep(time.Millisecond * 50)
	}
	writer.Writer.Stop()
	time.Sleep(5 * time.Second)
}

type ScreenWriter struct {
	Writer *uilive.Writer
}

func NewScreenWriter() *ScreenWriter {
	w1 := uilive.New()

	return &ScreenWriter{
		Writer: w1,
	}
}

func getDummyTable(data [][]interface{}) *simpletable.Table {
	table := simpletable.New()

	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "#"},
			{Align: simpletable.AlignCenter, Text: "NAME"},
			{Align: simpletable.AlignCenter, Text: "TAX"},
		},
	}

	subtotal := float64(0)
	for _, row := range data {
		r := []*simpletable.Cell{
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("%d", row[0].(int))},
			{Text: row[1].(string)},
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("$ %.2f", row[2].(float64))},
		}

		table.Body.Cells = append(table.Body.Cells, r)
		subtotal += row[2].(float64)
	}

	table.Footer = &simpletable.Footer{
		Cells: []*simpletable.Cell{
			{},
			{Align: simpletable.AlignRight, Text: "Subtotal"},
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("$ %.2f", subtotal)},
		},
	}

	table.SetStyle(simpletable.StyleCompactLite)
	return table
}
