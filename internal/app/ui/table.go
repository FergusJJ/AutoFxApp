package ui

import "github.com/rivo/tview"

//should hold the data that will be displayed within the table, as well as the table itself
//on set new data, should push to table and call rerender
type PositionTable struct {
	TableView *tview.Table
}

func newTable() *PositionTable {
	table := tview.NewTable()

	return &PositionTable{
		TableView: table,
	}
}
