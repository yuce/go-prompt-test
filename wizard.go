package hzsqlcl

import (
	"fmt"
	"hzsqlcl/components"
	"hzsqlcl/form"

	"github.com/gcla/gowid/widgets/grid"

	"github.com/gcla/gowid/widgets/holder"

	"github.com/gcla/gowid"
	"github.com/gcla/gowid/widgets/button"
	"github.com/gcla/gowid/widgets/columns"
	"github.com/gcla/gowid/widgets/edit"
	"github.com/gcla/gowid/widgets/fill"
	"github.com/gcla/gowid/widgets/framed"
	"github.com/gcla/gowid/widgets/overlay"
	"github.com/gcla/gowid/widgets/pile"
	"github.com/gcla/gowid/widgets/styled"
	"github.com/gcla/gowid/widgets/text"
	"github.com/gcla/gowid/widgets/vpadding"
)

type WizardPage interface {
	gowid.IWidget
	PageName() string
	UpdateState(state map[string]interface{})
	ExtraButtons() []*button.Widget
}

type WizardState map[string]interface{}

type WizardHandler func(app gowid.IApp, state WizardState)

type Wizard struct {
	pages               []WizardPage
	handler             WizardHandler
	currentPage         int
	currentHolderWidget *holder.Widget
	savedContainer      gowid.ISettableComposite
	savedSubWidget      gowid.IWidget
	state               WizardState
}

func NewWizard(pages []WizardPage, handler WizardHandler) *Wizard {
	if len(pages) == 0 {
		panic("no wizard pages!")
	}
	return &Wizard{
		pages:   pages,
		handler: handler,
	}
}

func (wiz *Wizard) Open(container gowid.ISettableComposite, width gowid.IWidgetDimension, app gowid.IApp) {
	wiz.currentPage = 0
	wiz.state = map[string]interface{}{}
	wiz.currentHolderWidget = holder.New(wiz.widgetForCurrentPage())
	wiz.savedContainer = container
	wiz.savedSubWidget = container.SubWidget()
	ov := overlay.New(wiz.currentHolderWidget, wiz.savedSubWidget,
		gowid.VAlignMiddle{}, gowid.RenderFlow{},
		gowid.HAlignMiddle{}, width)
	container.SetSubWidget(ov, app)
}

func (wiz *Wizard) close(app gowid.IApp) {
	wiz.savedContainer.SetSubWidget(wiz.savedSubWidget, app)
}

func (wiz *Wizard) buttonBarForPage() gowid.IWidget {
	isLastPage := wiz.currentPage == len(wiz.pages)-1
	btnNext := button.New(text.New("Next"))
	btnNext.OnClick(gowid.WidgetCallback{"cbNext", func(app gowid.IApp, w gowid.IWidget) {
		currentPage := wiz.pages[wiz.currentPage]
		currentPage.UpdateState(wiz.state)
		wiz.gotoNextPage(app)
	}})
	btnOk := button.New(text.New("OK"))
	btnOk.OnClick(gowid.WidgetCallback{"cbOK", func(app gowid.IApp, w gowid.IWidget) {
		currentPage := wiz.pages[wiz.currentPage]
		currentPage.UpdateState(wiz.state)
		if wiz.handler != nil {
			wiz.handler(app, wiz.state)
		}
		wiz.close(app)
	}})
	btnCancel := button.New(text.New("Cancel"))
	btnCancel.OnClick(gowid.WidgetCallback{"cbCancel", func(app gowid.IApp, w gowid.IWidget) {
		wiz.close(app)
	}})

	buttons := []interface{}{}
	page := wiz.pages[wiz.currentPage]
	for _, btn := range page.ExtraButtons() {
		buttons = append(buttons, btn)
	}
	buttons = append(buttons, btnCancel)
	if isLastPage {
		buttons = append(buttons, btnOk)
	} else {
		buttons = append(buttons, btnNext)
	}
	return columns.NewFixed(buttons...)
}

func (wiz *Wizard) widgetForCurrentPage() gowid.IWidget {
	page := wiz.pages[wiz.currentPage]
	flow := gowid.RenderFlow{}
	hline := styled.New(fill.New(' '), gowid.MakePaletteRef("line"))
	btnBar := wiz.buttonBarForPage()
	pilew := components.NewResizeablePile([]gowid.IContainerWidget{
		&gowid.ContainerWidget{IWidget: page, D: gowid.RenderWithWeight{2}},
		&gowid.ContainerWidget{vpadding.New(
			pile.New([]gowid.IContainerWidget{
				&gowid.ContainerWidget{IWidget: hline, D: gowid.RenderWithUnits{U: 1}},
				&gowid.ContainerWidget{IWidget: btnBar, D: flow},
			}),
			gowid.VAlignBottom{}, flow,
		), flow},
	})
	frame := framed.New(pilew, framed.Options{
		Frame: framed.UnicodeFrame,
		Title: fmt.Sprintf(" Create Mapping: %s ", page.PageName()),
	})
	return frame
}

func (wiz *Wizard) gotoNextPage(app gowid.IApp) {
	if wiz.currentPage < len(wiz.pages)-1 && wiz.currentHolderWidget != nil {
		wiz.currentPage++
		wiz.currentHolderWidget.SetSubWidget(wiz.widgetForCurrentPage(), app)
	}
}

const (
	MappingName      = "mappingName"
	MappingType      = "mappingType"
	MappingTypeKafka = "Kafka"
	MappingTypeFile  = "File"
)

type NameAndTypePage struct {
	gowid.IWidget
	mappingName string
	mappingType string
	editName    *edit.Widget
}

func NewNameAndTypePage() *NameAndTypePage {
	page := &NameAndTypePage{
		mappingType: MappingTypeKafka,
	}
	nameWidget := form.NewLabeledEdit(&page.mappingName, "Mapping Name: ")
	typeGroup := form.NewLabeledRadioGroup(&page.mappingType, "Mapping Type: ", MappingTypeKafka, MappingTypeFile)
	page.IWidget = pile.NewFixed(nameWidget, typeGroup)
	return page
}

func (p NameAndTypePage) PageName() string {
	return "Source"
}

func (p NameAndTypePage) UpdateState(state map[string]interface{}) {
	state[MappingName] = p.mappingName
	state[MappingType] = p.mappingType
}

func (p NameAndTypePage) ExtraButtons() []*button.Widget {
	return nil
}

type FieldsPage struct {
	gowid.IWidget
	fields []form.FieldFormState
}

func NewFieldsPage() *FieldsPage {
	widget := &FieldsPage{}
	widget.IWidget = holder.New(text.New("Click Add Field button to add fields."))
	return widget
}

func (p FieldsPage) PageName() string {
	return "Fields"
}

func (p FieldsPage) UpdateState(state map[string]interface{}) {
	for _, field := range p.fields {
		state[fmt.Sprintf("Field_%s", field.FieldName)] = field.FieldType
	}
}

func (p *FieldsPage) ExtraButtons() []*button.Widget {
	fieldTypes := []string{"VARCHAR", "INT"}
	addFieldBtn := button.New(text.New("Add Field"))
	addFieldBtn.OnClick(gowid.WidgetCallback{"cbAddField", func(app gowid.IApp, w gowid.IWidget) {
		frm := form.NewFormContainer("Add Field", form.NewFieldForm(fieldTypes...), nil, func(app gowid.IApp, state interface{}) {
			field := state.(form.FieldFormState)
			p.fields = append(p.fields, field)
			hl := p.IWidget.(*holder.Widget)
			widgets := []gowid.IWidget{}
			for _, f := range p.fields {
				txtFieldName := text.New(f.FieldName)
				txtFieldType := text.New(f.FieldType)
				widgets = append(widgets, txtFieldName, txtFieldType)
			}
			grd := grid.New(widgets, 20, 3, 1, gowid.HAlignMiddle{})
			hl.SetSubWidget(grd, app)
		})
		frm.Open(app.SubWidget().(*holder.Widget), gowid.RenderWithRatio{R: 0.5}, app)
	}})
	return []*button.Widget{addFieldBtn}
}

type OptionsPage struct {
	gowid.IWidget
	options []form.OptionFormState
}

func NewOptionsPage() *OptionsPage {
	widget := &OptionsPage{}
	widget.IWidget = holder.New(text.New("Click Add Option button to add fields."))
	return widget
}

func (p OptionsPage) PageName() string {
	return "Options"
}

func (p OptionsPage) UpdateState(state map[string]interface{}) {
	for _, option := range p.options {
		state[fmt.Sprintf("Option_%s", option.OptionName)] = option.OptionValue
	}
}

func (p *OptionsPage) ExtraButtons() []*button.Widget {
	addOptionBtn := button.New(text.New("Add Option"))
	addOptionBtn.OnClick(gowid.WidgetCallback{"cbAddOption", func(app gowid.IApp, w gowid.IWidget) {
		frm := form.NewFormContainer("Add Option", form.NewOptionForm(), nil, func(app gowid.IApp, state interface{}) {
			option := state.(form.OptionFormState)
			p.options = append(p.options, option)
			hl := p.IWidget.(*holder.Widget)
			widgets := []gowid.IWidget{}
			for _, f := range p.options {
				txtOptionName := text.New(f.OptionName)
				txtOptionType := text.New(f.OptionValue)
				widgets = append(widgets, txtOptionName, txtOptionType)
			}
			grd := grid.New(widgets, 20, 3, 1, gowid.HAlignMiddle{})
			hl.SetSubWidget(grd, app)
		})
		frm.Open(app.SubWidget().(*holder.Widget), gowid.RenderWithRatio{R: 0.5}, app)
	}})
	return []*button.Widget{addOptionBtn}
}
