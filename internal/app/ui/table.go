package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// should hold the data that will be displayed within the table, as well as the table itself
// on set new data, should push to table and call rerender
type PositionTable struct {
	TableView   *tview.Table
	WrappedView *tview.TextView
	Cols        []string
	Entries     []interface{}
}

// type Entry struct {
// 	Id     string
// 	Entry1 string
// 	Entry2 string
// }

/*
&{OrderID:94878953 ClOrdID:3f357761-3c81-4a50-b777-b3d38fd6bd92 TotNumReports: ExecType:0 OrdStatus:0 Symbol:3 Side:2 TransactTime:20230630-15:49:07.160 AvgPx: OrderQty:12000 LeavesQty:12000 CumQty:0 LastQty: OrdType:1 Price: StopPx: TimeInForce:3 ExpireTime: Text: OrdRejReason: PosMaintRptID:52942663 Designation: MassStatusReqID: AbsoluteTP: RelativeTP: AbsoluteSL: RelativeSL: TrailingSL: TriggerMethodSL: GuaranteedSL:}
*/
type Entry struct {
	OrderID       string
	ClOrdID       string
	ExecType      string
	OrdStatus     string
	Symbol        string
	Side          string
	TransactTime  string
	OrderQty      string
	LeavesQty     string
	PosMaintRptID string
}

func NewTable(cols []string) *PositionTable {
	table := tview.NewTable()
	table.SetBorders(false)
	table.SetTitle("Table")
	table.SetFixed(1, 1)
	table.SetTitleAlign(tview.AlignCenter)
	table.SetEvaluateAllRows(true)
	table.SetBorderPadding(0, 0, 2, 2)

	return &PositionTable{
		TableView: table,
		Cols:      cols,
		Entries:   []interface{}{},
	}
}

func (t *PositionTable) Init() {
	t.TableView.Clear()
	for col, v := range t.Cols {
		t.TableView.SetCell(
			0, col, &tview.TableCell{
				Text:  v,
				Align: tview.AlignCenter,
				Color: tcell.ColorDimGrey,
			},
		)
	}

}

func (t *PositionTable) AddEntry(entry ...interface{}) {
	if len(entry) == 0 {
		return
	}
	for _, v := range entry {
		t.Entries = append(t.Entries, v)
		row := len(t.Entries) //row 0 will just be title

		newEntry, ok := v.(Entry)
		if !ok {
			continue
		}

		/*
		   OrderID:      94878953,
		   				ClOrdID:      "3f357761-3c81-4a50-b777-b3d38fd6bd92",
		   				ExecType:     0,
		   				OrdStatus:    0,
		   				Symbol:       3,
		   				Side:         2,
		   				TransactTime: "20230630-15:49:07.160",
		   				OrderQty:     12000, LeavesQty: 12000,
		   				PosMaintRptID
		*/

		t.TableView.SetCell(
			row,
			0,
			&tview.TableCell{
				Text:  newEntry.OrderID,
				Align: tview.AlignLeft,
			},
		)

		t.TableView.SetCell(
			row,
			1,
			&tview.TableCell{
				Text:  newEntry.ClOrdID,
				Align: tview.AlignLeft,
			},
		)

		t.TableView.SetCell(
			row,
			2,
			&tview.TableCell{
				Text:  newEntry.ExecType,
				Align: tview.AlignLeft,
			},
		)

		t.TableView.SetCell(
			row,
			0,
			&tview.TableCell{
				Text:  newEntry.OrdStatus,
				Align: tview.AlignLeft,
			},
		)

		t.TableView.SetCell(
			row,
			0,
			&tview.TableCell{
				Text:  newEntry.Symbol,
				Align: tview.AlignLeft,
			},
		)

		t.TableView.SetCell(
			row,
			0,
			&tview.TableCell{
				Text:  newEntry.Side,
				Align: tview.AlignLeft,
			},
		)

		t.TableView.SetCell(
			row,
			0,
			&tview.TableCell{
				Text:  newEntry.TransactTime,
				Align: tview.AlignLeft,
			},
		)

		t.TableView.SetCell(
			row,
			0,
			&tview.TableCell{
				Text:  newEntry.OrderQty,
				Align: tview.AlignLeft,
			},
		)
		t.TableView.SetCell(
			row,
			0,
			&tview.TableCell{
				Text:  newEntry.PosMaintRptID,
				Align: tview.AlignLeft,
			},
		)

	}
}

func (t *PositionTable) RemoveEntry(id string) {
	for i, v := range t.Entries {
		entry, ok := v.(Entry)
		if !ok {
			continue
		}
		if id == entry.OrderID {
			leftovers := remove(t.Entries, i)
			t.Entries = []interface{}{}
			t.Init()
			t.AddEntry(leftovers...)

			break
		}
	}
	//update table
}

func remove(slice []interface{}, s int) []interface{} {
	return append(slice[:s], slice[s+1:]...)
}
