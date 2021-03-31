package components

import (
	"fmt"

	"github.com/gcla/gowid"
	"github.com/gcla/gowid/widgets/edit"
)

type LabeledEdit struct {
	gowid.IWidget
	name   string
	target *string
}

func NewLabeledEdit(name string, target *string, label string) *LabeledEdit {
	value, valueFound := FindGlobalValue(name)
	editOptions := edit.Options{
		Caption: label,
	}
	if valueFound {
		editOptions.Text = value.(string)
	}
	widget := &LabeledEdit{target: target}
	widget.name = name
	editWidget := edit.New(editOptions)
	editWidget.OnTextSet(gowid.WidgetCallback{fmt.Sprintf("edit%s", label), func(app gowid.IApp, w gowid.IWidget) {
		*widget.target = w.(*edit.Widget).Text()
	}})
	widget.IWidget = editWidget
	return widget
}

func (e *LabeledEdit) SetText(app gowid.IApp, txt string) {
	e.IWidget.(*edit.Widget).SetText(txt, app)
}

func (e LabeledEdit) ComponentName() string {
	return e.name
}
