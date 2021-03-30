package hzsqlcl

import (
	"github.com/gcla/gowid"
	"github.com/gcla/gowid/widgets/styled"
	"github.com/gcla/gowid/widgets/text"
)

type StatusBar struct {
	gowid.IWidget
	txt *text.Widget
}

func NewStatusBar() *StatusBar {
	txt := text.NewFromContentExt(CreateHintMessage("Hit tab to auto-complete"),
		text.Options{
			Align: gowid.HAlignLeft{},
		},
	)
	bar := styled.New(txt, gowid.MakePaletteRef("hint"))
	return &StatusBar{bar, txt}
}

func (s *StatusBar) SetHint(app gowid.IApp, text string) {
	s.txt.SetContent(app, CreateHintMessage(text))
	// TODO: use a holder and ditch the following
	app.Redraw()
}

func (s *StatusBar) SetError(app gowid.IApp, text string) {
	s.txt.SetContent(app, CreateErrorMessage(text))
	// TODO: use a holder and ditch the following
	app.Redraw()
}

func (s *StatusBar) Clear(app gowid.IApp) {
	s.SetHint(app, "")
}
