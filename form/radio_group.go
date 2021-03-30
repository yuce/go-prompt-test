package form

import (
	"fmt"

	"github.com/gcla/gowid"
	"github.com/gcla/gowid/widgets/columns"
	"github.com/gcla/gowid/widgets/pile"
	"github.com/gcla/gowid/widgets/radio"
	"github.com/gcla/gowid/widgets/text"
)

type RadioGroup struct {
	gowid.IWidget
	target *string
}

func NewRadioGroup(target *string, items ...string) *RadioGroup {
	widget := &RadioGroup{target: target}
	rbGroup := []radio.IWidget{}
	rows := []interface{}{}
	for _, name := range items {
		func(name string) {
			rb := radio.New(&rbGroup)
			rb.OnClick(gowid.WidgetCallback{fmt.Sprintf("cbRadio%s", name), func(app gowid.IApp, w gowid.IWidget) {
				*widget.target = name
			}})
			rbt := text.New(fmt.Sprintf(" %s ", name))
			rows = append(rows, columns.NewFixed(rb, rbt))
		}(name)
	}
	widget.IWidget = pile.NewFixed(rows...)
	return widget
}

type LabeledRadioGroup struct {
	gowid.IWidget
	target *string
}

func NewLabeledRadioGroup(target *string, label string, items ...string) *LabeledRadioGroup {
	labelWidget := text.New(label)
	widget := &LabeledRadioGroup{target: target}
	radioGroupWidget := NewRadioGroup(target, items...)
	widget.IWidget = columns.NewFixed(labelWidget, radioGroupWidget)
	return widget
}

type LabeledEditRadioGroup struct {
	gowid.IWidget
	editTarget *string
	rgTarget   *string
}

//func NewLabaledEditRadioGroup(editTarget *string, rgTarget *string, label string)
