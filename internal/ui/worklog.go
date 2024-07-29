package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

// TODO: Remove sample data and use real pending worklogs.
func NewWorklogTable(_ fyne.Window) fyne.CanvasObject {
	t := widget.NewTableWithHeaders(
		func() (int, int) { return 10, 4 },
		func() fyne.CanvasObject {
			return widget.NewLabel("Cell 000, 000")
		},
		func(id widget.TableCellID, cell fyne.CanvasObject) {
			label := cell.(*widget.Label)
			switch id.Col {
			case 0:
				label.SetText("A longer cell")
			default:
				label.SetText(fmt.Sprintf("Cell %d, %d", id.Row+1, id.Col+1))
			}
		},
	)
	t.SetColumnWidth(0, 102)
	t.SetRowHeight(2, 50)

	// TODO: binding json
	t.CreateHeader = func() fyne.CanvasObject {
		return widget.NewLabel("")
	}
	t.UpdateHeader = func(id widget.TableCellID, template fyne.CanvasObject) {
		label := template.(*widget.Label)
		switch id.Col {
		case 0:
			label.SetText("A longer Header")
		default:
			label.SetText(fmt.Sprintf("Header %d, %d", id.Row+1, id.Col+1))
		}
	}
	return t
}
