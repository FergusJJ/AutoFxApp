package main

import (
	"fmt"
	"time"

	"github.com/gosuri/uilive"
)

//maybe add lines to a slice, and on each print, iterate and write each one out
//could try to use the table inside of the writer, or could just append new writer for each position
//each writer should probably track the last thing they wrote in order to prevent lines that aren't updated from being deleted

//func (w *ScreenWriter) Write(items [int]string){] //if items[0] is empty then do not overwrite w.Writers[0] or something

func main() {
	writer := NewScreenWriter()
	writer.Writer.Start()
	for i := 0; i <= 50; i++ {
		fmt.Fprintf(writer.Writer, "this is line 1 %d\n", i)
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
