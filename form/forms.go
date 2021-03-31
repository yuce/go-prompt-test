package form

import (
	"hzsqlcl/components"

	"github.com/gcla/gowid/widgets/fill"

	"github.com/gcla/gowid/widgets/shadow"

	"github.com/gcla/gowid/widgets/cellmod"

	"github.com/gcla/gowid"
	"github.com/gcla/gowid/widgets/button"
	"github.com/gcla/gowid/widgets/columns"
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
	buttons := []*button.Widget{}
	for _, btn := range f.extraButtons {
		buttons = append(buttons, btn)
	}
	buttons = append(buttons, cancelBtn, okBtn)
	colsW := []gowid.IContainerWidget{}
	for _, btn := range buttons {
		colsW = append(colsW, components.MakeStylishButton(btn))
	}
	return columns.New(colsW)
}

func (f *FormContainer) frame() gowid.IWidget {
	borderStyle := gowid.MakePaletteRef("border")
	backgroundStyle := gowid.MakePaletteRef("background")
	flow := gowid.RenderFlow{}
	hline := styled.New(fill.New('-'), borderStyle)
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
		Style: backgroundStyle,
	})
	styledFrame := shadow.New(cellmod.Opaque(styled.New(frame, backgroundStyle)), 1)
	return styledFrame
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
	//buttonStyle := gowid.MakePaletteEntry(DefaultButtonText, DefaultButton)
	//backgroundStyle := gowid.MakePaletteEntry(DefaultText, DefaultBackground)
	//borderStyle := gowid.MakePaletteEntry(DefaultButton, DefaultBackground)
	widget := &FieldForm{fieldFormState: FieldFormState{FieldType: "VARCHAR"}}
	fieldNameWidget := NewLabeledEdit(&widget.fieldFormState.FieldName, "Column Name: ")
	fieldTypeWidget := NewLabeledRadioGroup(&widget.fieldFormState.FieldType, "Column Type: ", items...)
	pl := pile.NewFixed(fieldNameWidget, fieldTypeWidget)
	//w := hpadding.New(
	//	styled.NewExt(pl, backgroundStyle, buttonStyle),
	//	gowid.HAlignMiddle{},
	//	gowid.RenderFixed{},
	//)
	widget.IWidget = pl
	return widget
}

func (f FieldForm) State() interface{} {
	return f.fieldFormState
}

type OptionFormState struct {
	OptionName  string
	OptionValue string
}

type OptionForm struct {
	gowid.IWidget
	optionFormState OptionFormState
}

func NewOptionForm() *OptionForm {
	widget := &OptionForm{}
	optionNameWidget := NewLabeledEdit(&widget.optionFormState.OptionName, "Option Name: ")
	optionValueWidget := NewLabeledEdit(&widget.optionFormState.OptionValue, "Option Value: ")
	widget.IWidget = pile.NewFixed(optionNameWidget, optionValueWidget)
	return widget
}

func (f OptionForm) State() interface{} {
	return f.optionFormState
}
