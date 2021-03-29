package main

import (
	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()
	form := tview.NewForm().
		AddInputField("Name", "", 20, nil, nil).
		AddInputField("Field 1 Name", "", 20, nil, nil).
		AddDropDown("Field 1 Type", []string{"String", "Int", "Float", "Boolean"}, 0, nil).
		//AddCheckbox("Age 18+", false, nil).
		//AddPasswordField("Password", "", 10, '*', nil).
		AddInputField("Option Key 1", "", 20, nil, nil).
		AddInputField("Option Value 1", "", 20, nil, nil).
		AddButton("Add Field", nil).
		AddButton("Add Option", nil).
		AddButton("Save", nil).
		AddButton("Cancel", func() {
			app.Stop()
		})
	form.SetBorder(true).SetTitle(" Create Mapping ").SetTitleAlign(tview.AlignCenter)
	if err := app.SetRoot(form, true).SetFocus(form).Run(); err != nil {
		panic(err)
	}
}
