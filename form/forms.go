package form

import (
	"hzsqlcl/components"

	"github.com/gcla/gowid"
	"github.com/gcla/gowid/widgets/button"
	"github.com/gcla/gowid/widgets/columns"
	"github.com/gcla/gowid/widgets/fill"
	"github.com/gcla/gowid/widgets/framed"
	"github.com/gcla/gowid/widgets/overlay"
	"github.com/gcla/gowid/widgets/pile"
	"github.com/gcla/gowid/widgets/styled"
	"github.com/gcla/gowid/widgets/text"
	"github.com/gcla/gowid/widgets/vpadding"
)

type FormFragment interface {
	gowid.IWidget
	State() interface{}
}

type FormContainerHandler func(app gowid.IApp, state interface{})

type FormContainer struct {
	widget         FormFragment
	title          string
	extraButtons   []*button.Widget
	handler        FormContainerHandler
	savedContainer gowid.ISettableComposite
	savedSubWidget gowid.IWidget
}

func NewFormContainer(title string, widget FormFragment, extraButtons []*button.Widget, handler FormContainerHandler) *FormContainer {
	return &FormContainer{
		widget:       widget,
		title:        title,
		extraButtons: extraButtons,
		handler:      handler,
	}
}

func (f *FormContainer) Open(container gowid.ISettableComposite, width gowid.IWidgetDimension, app gowid.IApp) {
	f.savedContainer = container
	f.savedSubWidget = container.SubWidget()
	ov := overlay.New(f.frame(), f.savedSubWidget,
		gowid.VAlignMiddle{}, gowid.RenderFlow{},
		gowid.HAlignMiddle{}, width)
	container.SetSubWidget(ov, app)
}

func (f *FormContainer) close(app gowid.IApp) {
	f.savedContainer.SetSubWidget(f.savedSubWidget, app)
}

func (f *FormContainer) buttonBar() gowid.IWidget {
	okBtn := button.New(text.New("OK"))
	okBtn.OnClick(gowid.WidgetCallback{"cbOK", func(app gowid.IApp, w gowid.IWidget) {
		if f.handler != nil {
			f.handler(app, f.widget.State())
		}
		f.close(app)
	}})
	cancelBtn := button.New(text.New("Cancel"))
	cancelBtn.OnClick(gowid.WidgetCallback{"cbCancel", func(app gowid.IApp, w gowid.IWidget) {
		f.close(app)
	}})
	buttons := []interface{}{}
	for _, btn := range f.extraButtons {
		buttons = append(buttons, btn)
	}
	buttons = append(buttons, cancelBtn, okBtn)
	return columns.NewFixed(buttons...)
}

func (f *FormContainer) frame() gowid.IWidget {
	flow := gowid.RenderFlow{}
	hline := styled.New(fill.New(' '), gowid.MakePaletteRef("line"))
	pilew := components.NewResizeablePile([]gowid.IContainerWidget{
		&gowid.ContainerWidget{IWidget: f.widget, D: gowid.RenderWithWeight{2}},
		&gowid.ContainerWidget{vpadding.New(
			pile.New([]gowid.IContainerWidget{
				&gowid.ContainerWidget{IWidget: hline, D: gowid.RenderWithUnits{U: 1}},
				&gowid.ContainerWidget{IWidget: f.buttonBar(), D: flow},
			}),
			gowid.VAlignBottom{}, flow,
		), flow},
	})
	frame := framed.New(pilew, framed.Options{
		Frame: framed.UnicodeFrame,
		Title: f.title,
	})
	return frame
}

type FieldFormState struct {
	FieldName string
	FieldType string
}

type FieldForm struct {
	gowid.IWidget
	fieldFormState FieldFormState
}

func NewFieldForm(items ...string) *FieldForm {
	widget := &FieldForm{fieldFormState: FieldFormState{FieldType: "VARCHAR"}}
	fieldNameWidget := NewLabeledEdit(&widget.fieldFormState.FieldName, "Field Name:")
	fieldTypeWidget := NewLabeledRadioGroup(&widget.fieldFormState.FieldType, "Field Type:", items...)
	widget.IWidget = pile.NewFixed(fieldNameWidget, fieldTypeWidget)
	return widget
}

func (f FieldForm) State() interface{} {
	return f.fieldFormState
}
