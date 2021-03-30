package hzsqlcl

import (
	"fmt"

	"github.com/gcla/gowid/widgets/holder"

	"github.com/gcla/gowid"
	"github.com/gcla/gowid/widgets/button"
	"github.com/gcla/gowid/widgets/columns"
	"github.com/gcla/gowid/widgets/edit"
	"github.com/gcla/gowid/widgets/fill"
	"github.com/gcla/gowid/widgets/framed"
	"github.com/gcla/gowid/widgets/grid"

	//"github.com/gcla/gowid/widgets/holder"
	"github.com/gcla/gowid/widgets/overlay"
	"github.com/gcla/gowid/widgets/pile"
	"github.com/gcla/gowid/widgets/radio"
	"github.com/gcla/gowid/widgets/styled"
	"github.com/gcla/gowid/widgets/text"
	"github.com/gcla/gowid/widgets/vpadding"
)

type WizardPage interface {
	gowid.IWidget
	PageName() string
}

type Wizard struct {
	//gowid.IWidget
	//dlg         *dialog.Widget
	currentPage         int
	pages               []WizardPage
	currentHolderWidget *holder.Widget
}

func NewWizard(pages []WizardPage) *Wizard {
	if len(pages) == 0 {
		panic("no wizard pages!")
	}
	return &Wizard{pages: pages}
}

func (wiz *Wizard) Open(container gowid.ISettableComposite, width gowid.IWidgetDimension, app gowid.IApp) {
	widget := wiz.widgetForCurrentPage()
	hl := holder.New(widget)
	wiz.currentHolderWidget = hl
	ov := overlay.New(hl, container.SubWidget(),
		gowid.VAlignMiddle{}, gowid.RenderFlow{}, // Intended to mean use as much vertical space as you need
		gowid.HAlignMiddle{}, width)
	container.SetSubWidget(ov, app)
}

func (wiz *Wizard) buttonBarForPage() gowid.IWidget {
	isLastPage := wiz.currentPage == len(wiz.pages)-1
	nextBtn := button.New(text.New("Next"))
	nextBtn.OnClick(gowid.WidgetCallback{"cbNext", func(app gowid.IApp, w gowid.IWidget) {
		wiz.gotoNextPage(app)
	}})

	okBtn := button.New(text.New("OK"))
	okBtn.OnClick(gowid.WidgetCallback{"cbOK", func(app gowid.IApp, w gowid.IWidget) {
		fmt.Println("OK")
	}})
	cancelBtn := button.New(text.New("Cancel"))
	cancelBtn.OnClick(gowid.WidgetCallback{"cbCancel", func(app gowid.IApp, w gowid.IWidget) {
		fmt.Println("CANCEL")
	}})

	buttons := []interface{}{cancelBtn}
	if isLastPage {
		buttons = append(buttons, okBtn)
	} else {
		buttons = append(buttons, nextBtn)
	}
	return columns.NewFixed(buttons...)
}

func (wiz *Wizard) widgetForCurrentPage() gowid.IWidget {
	page := wiz.pages[wiz.currentPage]
	flow := gowid.RenderFlow{}
	hline := styled.New(fill.New(' '), gowid.MakePaletteRef("line"))
	btnBar := wiz.buttonBarForPage()
	pilew := NewResizeablePile([]gowid.IContainerWidget{
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

type NameAndTypePage struct {
	gowid.IWidget
}

func NewNameAndTypePage() *NameAndTypePage {
	txtName := text.New("Mapping Name:")
	editName := edit.New()
	txtType := text.New("Mapping Type:")
	fixed := gowid.RenderFixed{}
	rbgroup := make([]radio.IWidget, 0)
	rb1 := radio.New(&rbgroup)
	rbt1 := text.New(" Kafka ")
	rb2 := radio.New(&rbgroup)
	rbt2 := text.New(" File ")
	c2cols := []gowid.IContainerWidget{
		&gowid.ContainerWidget{rb1, fixed},
		&gowid.ContainerWidget{rbt1, fixed},
		&gowid.ContainerWidget{rb2, fixed},
		&gowid.ContainerWidget{rbt2, fixed},
	}
	cols2 := columns.New(c2cols)
	widgets := []gowid.IWidget{txtName, editName, txtType, cols2}
	grid1 := grid.New(widgets, 20, 3, 1, gowid.HAlignMiddle{})
	return &NameAndTypePage{grid1}
}

func (p NameAndTypePage) PageName() string {
	return "Name and Type"
}

type PageWidget2 struct {
	gowid.IWidget
}

func NewPageWidget2() *PageWidget2 {
	txtName := text.New("XXXXX:")
	editName := edit.New()
	widgets := []gowid.IWidget{txtName, editName}
	grid1 := grid.New(widgets, 20, 3, 1, gowid.HAlignMiddle{})
	return &PageWidget2{grid1}
}

func (p PageWidget2) PageName() string {
	return "Page 2"
}
