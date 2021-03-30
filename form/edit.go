package form

import (
	"fmt"

	"github.com/gcla/gowid"
	"github.com/gcla/gowid/widgets/columns"
	"github.com/gcla/gowid/widgets/edit"
	"github.com/gcla/gowid/widgets/text"
)

type LabeledEdit struct {
	gowid.IWidget
	target *string
}

func NewLabeledEdit(target *string, label string) *LabeledEdit {
	labelWidget := text.New(label)
	widget := &LabeledEdit{target: target}
	editWidget := edit.New()
	editWidget.OnTextSet(gowid.WidgetCallback{fmt.Sprintf("edit%s", label), func(app gowid.IApp, w gowid.IWidget) {
		edt := w.(*edit.Widget)
		*widget.target = edt.Text()
	}})
	widget.IWidget = columns.NewFixed(labelWidget, editWidget)
	return widget
}
