package app

import "github.com/gosuri/uilive"

type ScreenWriter struct {
	Writer *uilive.Writer
}

func NewScreenWriter() *ScreenWriter {
	w1 := uilive.New()

	return &ScreenWriter{
		Writer: w1,
	}
}
