package utils

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

func NewDataList(bindings binding.ExternalStringList) *widget.List {
	list := widget.NewListWithData(
		bindings,
		func() fyne.CanvasObject {
			option := widget.NewLabel("")
			check := widget.NewCheck("", func(checked bool) {
				fmt.Println("checked", checked)
				fmt.Println("wch--- data: %+v\n", option.Text)
			})
			return container.NewHBox(check, option)
		},
		func(i binding.DataItem, item fyne.CanvasObject) {
			info := i.(binding.String)
			item.(*fyne.Container).Objects[1].(*widget.Label).Bind(info)
		},
	)
	return list
}
