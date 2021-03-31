package components

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
	name   string
	target *string
}

func NewRadioGroup(name string, target *string, items ...string) *RadioGroup {
	// TODO: set selected by name
	widget := &RadioGroup{
		name:   name,
		target: target,
	}
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

func NewLabeledRadioGroup(name string, target *string, label string, items ...string) *LabeledRadioGroup {
	labelWidget := text.New(label)
	widget := &LabeledRadioGroup{target: target}
	radioGroupWidget := NewRadioGroup(name, target, items...)
	widget.IWidget = columns.NewFixed(labelWidget, radioGroupWidget)
	return widget
}

type LabeledEditRadioGroup struct {
	gowid.IWidget
	editTarget *string
	rgTarget   *string
}

//func NewLabaledEditRadioGroup(editTarget *string, rgTarget *string, label string)
