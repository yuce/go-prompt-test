package hzsqlcl

import (
	"fmt"
	"hzsqlcl/components"
	"hzsqlcl/form"
	"math/rand"

	"github.com/gcla/gowid/widgets/cellmod"
	"github.com/gcla/gowid/widgets/shadow"

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

	buttons := []*button.Widget{}
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
	colsW := []gowid.IContainerWidget{}
	for _, btn := range buttons {
		colsW = append(colsW, components.MakeStylishButton(btn))
	}
	return columns.New(colsW)
}

func (wiz *Wizard) widgetForCurrentPage() gowid.IWidget {
	borderStyle := gowid.MakePaletteRef("border")
	backgroundStyle := gowid.MakePaletteRef("background")

	page := wiz.pages[wiz.currentPage]
	flow := gowid.RenderFlow{}
	hline := styled.New(fill.New('-'), borderStyle)
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
	styledFrame := shadow.New(cellmod.Opaque(styled.New(frame, backgroundStyle)), 1)
	return styledFrame
}

func (wiz *Wizard) gotoNextPage(app gowid.IApp) {
	if wiz.currentPage < len(wiz.pages)-1 && wiz.currentHolderWidget != nil {
		wiz.currentPage++
		wiz.currentHolderWidget.SetSubWidget(wiz.widgetForCurrentPage(), app)
	}
}

const (
	MappingName                  = "mappingName"
	MappingType                  = "mappingType"
	SerializationType            = "serializationType"
	ConnectionAddress            = "connectionAddress"
	MappingTypeKafka             = "Kafka"
	MappingTypeFile              = "File"
	MappingTypeIMap              = "IMap"
	MappingSerializationJson     = "json"
	MappingSerializationAvro     = "avro"
	MappingSerializationPortable = "portable"
)

type SourceNameAndTypePage struct {
	gowid.IWidget
	mappingName string
	mappingType string
	editName    *edit.Widget
}

func NewSourceNameAndTypePage() *SourceNameAndTypePage {
	page := &SourceNameAndTypePage{
		mappingType: MappingTypeKafka,
	}
	nameWidget := form.NewLabeledEdit(&page.mappingName, "Mapping Name: ")
	typeGroup := form.NewLabeledRadioGroup(&page.mappingType, "Mapping Type: ", MappingTypeKafka, MappingTypeFile)
	page.IWidget = pile.NewFixed(nameWidget, typeGroup)
	return page
}

func (p SourceNameAndTypePage) PageName() string {
	return "Source"
}

func (p SourceNameAndTypePage) UpdateState(state map[string]interface{}) {
	state[MappingName] = p.mappingName
	state[MappingType] = p.mappingType
}

func (p SourceNameAndTypePage) ExtraButtons() []*button.Widget {
	return nil
}

type FieldsPage struct {
	gowid.IWidget
	fields []form.FieldFormState
	pageName string
}

func NewFieldsPage(header string, pageName string) *FieldsPage {
	widget := &FieldsPage{pageName: pageName}
	widget.IWidget = holder.New(text.New(header))
	return widget
}

func (p FieldsPage) PageName() string {
	return p.pageName
}

func (p FieldsPage) UpdateState(state map[string]interface{}) {
	for _, field := range p.fields {
		state[fmt.Sprintf("Field_%s", field.FieldName)] = field.FieldType
	}
}

func (p *FieldsPage) ExtraButtons() []*button.Widget {
	fieldTypes := []string{"VARCHAR", "INT"}
	addFieldBtn := button.New(text.New("Add Column"))
	addFieldBtn.OnClick(gowid.WidgetCallback{"cbAddField", func(app gowid.IApp, w gowid.IWidget) {
		frm := form.NewFormContainer("Add Column", form.NewFieldForm(fieldTypes...), nil, func(app gowid.IApp, state interface{}) {
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

type SerializationPage struct {
	gowid.IWidget
	serializationType string
	pageName string
}

func NewSerializationPage(pageName string) *SerializationPage {
	widget := &SerializationPage{
		serializationType: MappingSerializationJson,
		pageName: pageName,
	}

	serializationGroup := form.NewLabeledRadioGroup(&widget.serializationType, "Serialization Type: ", MappingSerializationJson, MappingSerializationAvro, MappingSerializationPortable)
	widget.IWidget = pile.NewFixed(serializationGroup)

	return widget
}

func (p SerializationPage) PageName() string {
	return p.pageName
}

func (p SerializationPage) ExtraButtons() []*button.Widget {
	return nil
}

func (p SerializationPage) UpdateState(state map[string]interface{}) {
	state[fmt.Sprintf("Option_%s", "value_format")] = p.serializationType
	//state[SerializationType] = p.serializationType
}

type SourceOptionsPage struct {
	gowid.IWidget
	connectionAddress string
}

func NewSourceOptionsPage() *SourceOptionsPage {
	widget := &SourceOptionsPage{
		connectionAddress: "127.0.0.1:9092",
	}
	widget.IWidget = form.NewLabeledEdit(&widget.connectionAddress, "Connection Address: ")
	return widget
}

func (p SourceOptionsPage) PageName() string {
	return "Additional Options"
}

func (p SourceOptionsPage) UpdateState(state map[string]interface{}) {
	state[fmt.Sprintf("Option_%s", "bootstrap.server")] = p.connectionAddress
}

func (p *SourceOptionsPage) ExtraButtons() []*button.Widget {
	return nil
}

type SinkOptionsPage struct {
	gowid.IWidget
	connectionAddress string
}

func NewSinkOptionsPage() *SinkOptionsPage {
	widget := &SinkOptionsPage{
		connectionAddress: "127.0.0.1:9092",
	}

	widget.IWidget = form.NewLabeledEdit(&widget.connectionAddress, "Connection Address: ")
	return widget
}

func (p SinkOptionsPage) PageName() string {
	return "Additional Options"
}

func (p SinkOptionsPage) UpdateState(state map[string]interface{}) {
	state[fmt.Sprintf("Option_%s", "bootstrap.server")] = p.connectionAddress
	state[fmt.Sprintf("Option_%s", "key_format")] = "int"
	randomInt := rand.Intn(100)
	state[fmt.Sprintf("Option_Int_%s", "valuePortableFactoryId")] = randomInt
	randomInt++
	state[fmt.Sprintf("Option_Int_%s", "valuePortableClassId")] = randomInt
}

func (p *SinkOptionsPage) ExtraButtons() []*button.Widget {
	return nil
}

// Original general purpose OptionsPage
//
//type OptionsPage struct {
//	gowid.IWidget
//	options []form.OptionFormState
//}
//
//func NewOptionsPage() *OptionsPage {
//	widget := &OptionsPage{}
//	widget.IWidget = holder.New(text.New("Click Add Option button to add options."))
//
//	return widget
//}
//
//func (p OptionsPage) PageName() string {
//	return "Options"
//}
//
//func (p OptionsPage) UpdateState(state map[string]interface{}) {
//	for _, option := range p.options {
//		state[fmt.Sprintf("Option_%s", option.OptionName)] = option.OptionValue
//	}
//}
//
//func (p *OptionsPage) ExtraButtons() []*button.Widget {
//	addOptionBtn := button.New(text.New("Add Option"))
//	addOptionBtn.OnClick(gowid.WidgetCallback{"cbAddOption", func(app gowid.IApp, w gowid.IWidget) {
//		frm := form.NewFormContainer("Add Option", form.NewOptionForm(), nil, func(app gowid.IApp, state interface{}) {
//			option := state.(form.OptionFormState)
//			p.options = append(p.options, option)
//			hl := p.IWidget.(*holder.Widget)
//			widgets := []gowid.IWidget{}
//			for _, f := range p.options {
//				txtOptionName := text.New(f.OptionName)
//				txtOptionType := text.New(f.OptionValue)
//				widgets = append(widgets, txtOptionName, txtOptionType)
//			}
//			grd := grid.New(widgets, 20, 3, 1, gowid.HAlignMiddle{})
//			hl.SetSubWidget(grd, app)
//		})
//		frm.Open(app.SubWidget().(*holder.Widget), gowid.RenderWithRatio{R: 0.5}, app)
//	}})
//	return []*button.Widget{addOptionBtn}
//}

// Create mapping for sink wizard
type SinkNameAndTypePage struct {
	gowid.IWidget
	mappingName string
	mappingType string
	editName    *edit.Widget
}

func NewSinkNameAndTypePage() *SinkNameAndTypePage {
	page := &SinkNameAndTypePage{
		mappingType: MappingTypeIMap,
	}
	nameWidget := form.NewLabeledEdit(&page.mappingName, "Mapping Name: ")
	typeGroup := form.NewLabeledRadioGroup(&page.mappingType, "Mapping Type: ",MappingTypeIMap, MappingTypeKafka)
	page.IWidget = pile.NewFixed(nameWidget, typeGroup)
	return page
}

func (p SinkNameAndTypePage) PageName() string {
	return "Sink"
}

func (p SinkNameAndTypePage) UpdateState(state map[string]interface{}) {
	state[MappingName] = p.mappingName
	state[MappingType] = p.mappingType
}

func (p SinkNameAndTypePage) ExtraButtons() []*button.Widget {
	return nil
}
