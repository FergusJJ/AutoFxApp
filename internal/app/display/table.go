package display

import (
	"fmt"

	"github.com/alexeyco/simpletable"
)

var tableStyle = simpletable.StyleUnicode

func DrawTable(data [][]interface{}) {
	table := simpletable.New()
	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "#"}, //PosMaintRptID
			{Align: simpletable.AlignCenter, Text: "SYMBOL_ID"},
			{Align: simpletable.AlignCenter, Text: "VOLUME"},
			{Align: simpletable.AlignCenter, Text: "CURRENT_PRICE"},
			{Align: simpletable.AlignCenter, Text: "QTTY"},
			{Align: simpletable.AlignCenter, Text: "VOLUME"},
		},
	}

	subtotal := 0
	for _, row := range data {
		r := []*simpletable.Cell{
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("%d", row[0].(int))},
			{Text: row[1].(string)},
			{Text: row[2].(string)},
			{Text: row[3].(string)},
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("%d", row[4])},
		}

		table.Body.Cells = append(table.Body.Cells, r)
		subtotal += row[4].(int)
	}

	table.Footer = &simpletable.Footer{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignRight, Span: 4, Text: "Subtotal"},
			{Align: simpletable.AlignRight, Text: fmt.Sprint("#")},
		},
	}
	// table.Footer = &simpletable.Footer{
	// 	Cells: []*simpletable.Cell{
	// 		// 	{Text: "Balance"},
	// 		{Align: simpletable.AlignCenter, Span: 1, Text: "#"},
	// 	},
	// }
	table.Println()
}
