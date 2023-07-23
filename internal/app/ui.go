package app

import (
	"fmt"
	"io"
	"time"

	"github.com/alexeyco/simpletable"
	"github.com/fatih/color"
	"github.com/gosuri/uilive"
)

type screenWriter struct {
	writer      *uilive.Writer
	tableWriter io.Writer
	tableData   []interface{}
	messages    []string
	maxLines    int
}

/*



 */

func NewScreenWriter(maxLines int) *screenWriter {
	if maxLines < 1 {
		maxLines = 1
	}
	w1 := uilive.New()
	w2 := w1.Newline()

	return &screenWriter{
		messages:    []string{},
		maxLines:    maxLines,
		writer:      w1,
		tableWriter: w2,
		tableData:   nil,
	}
}

// have this run inside loop
func (sw *screenWriter) Write(newMessage string) {
	sw.writer.Start()

	if newMessage != "" {
		sw.updateMessages(newMessage)
	}

	if sw.tableData != nil {
		fmt.Fprint(sw.tableWriter, sw.getTable())
	}
	for _, message := range sw.messages {
		if message == "" {
			continue
		}
		color.New(color.FgYellow).Fprintln(sw.writer, message)
	}
	sw.writer.Flush()
	sw.writer.Stop()
}

func (sw *screenWriter) updateMessages(message string) {
	message = formatMessage(message)
	if len(sw.messages) == sw.maxLines {
		sw.messages = sw.messages[1:]
	}
	sw.messages = append(sw.messages, message)
}

func formatMessage(message string) string {
	currentTime := time.Now()
	return fmt.Sprintf("%s %s", currentTime.Format("15:04:05"), message)
}

func (sw *screenWriter) getTable() string {
	table := simpletable.New()
	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "#"},
			{Align: simpletable.AlignCenter, Text: "Symbol"},
			{Align: simpletable.AlignCenter, Text: "Volume"},
			{Align: simpletable.AlignCenter, Text: "Direction"},
			{Align: simpletable.AlignCenter, Text: "EntryPrice"},
			{Align: simpletable.AlignCenter, Text: "CurrentPrice"},
			{Align: simpletable.AlignCenter, Text: "GrossProfitEUR"},
			{Align: simpletable.AlignCenter, Text: "NetProfitEUR"},
		},
	}
	/*
		Range through data, append all data to cols
	*/

	/*
		Add any footers if necessary
	*/

	table.SetStyle(simpletable.StyleCompactLite)
	return table.String()
}
