package components

import (
	"fmt"

	"github.com/gcla/gowid"
	"github.com/gcla/gowid/widgets/edit"
)

type LabeledEdit struct {
	gowid.IWidget
	target *string
	name   string
}

func NewLabeledEdit(target *string, label string) *LabeledEdit {
	widget := &LabeledEdit{target: target}
	editWidget := edit.New(edit.Options{
		Caption: label,
	})
	editWidget.OnTextSet(gowid.WidgetCallback{fmt.Sprintf("edit%s", label), func(app gowid.IApp, w gowid.IWidget) {
		edt := w.(*edit.Widget)
		*widget.target = edt.Text()
	}})
	widget.IWidget = editWidget
	return widget
}

func (e *LabeledEdit) SetText(app gowid.IApp, txt string) {
	e.IWidget.(*edit.Widget).SetText(txt, app)
}

func (e *LabeledEdit) MaybeSetValueFromState(app gowid.IApp, state map[string]interface{}) {
	if v, ok := state[e.name]; ok {
		e.SetText(app, v.(string))
	}
}
