package display

import (
	"github.com/alexeyco/simpletable"
)

var tableStyle = simpletable.StyleUnicode

func DrawTable(data [][]interface{}) {
	table := simpletable.New()
	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "ID"},
			{Align: simpletable.AlignCenter, Text: "SYMBOL"},
			{Align: simpletable.AlignCenter, Text: "VOLUME"},
			{Align: simpletable.AlignCenter, Text: "DIRECTION"},
			{Align: simpletable.AlignCenter, Text: "ENTRY PRICE"},
			{Align: simpletable.AlignCenter, Text: "CURRENT EXCHANGE RATE"},
			{Align: simpletable.AlignCenter, Text: "SWAP"},
			{Align: simpletable.AlignCenter, Text: "COMMISSIONS"},
		},
	}
	table.SetStyle(tableStyle)
	for _, row := range data {
		r := []*simpletable.Cell{
			{Align: simpletable.AlignRight, Text: row[0].(string)},
			{Text: row[1].(string)},
			{Text: row[2].(string)},
			{Text: row[3].(string)},
			{Text: row[4].(string)},
			{Text: row[5].(string)},
		}
		table.Body.Cells = append(table.Body.Cells, r)
	}
	// table.Footer = &simpletable.Footer{
	// 	Cells: []*simpletable.Cell{
	// 		// 	{Text: "Balance"},
	// 		{Align: simpletable.AlignCenter, Span: 1, Text: "#"},
	// 	},
	// }
	table.Println()
}
