package components

import (
	"github.com/gcla/gowid"
	"github.com/gcla/gowid/widgets/button"
	"github.com/gcla/gowid/widgets/hpadding"
	"github.com/gcla/gowid/widgets/styled"
)

func MakeStylishButton(btn *button.Widget) *gowid.ContainerWidget {
	backgroundStyle := gowid.MakePaletteRef("background")
	buttonStyle := gowid.MakePaletteRef("button")
	return &gowid.ContainerWidget{
		hpadding.New(
			styled.NewExt(btn, backgroundStyle, buttonStyle),
			gowid.HAlignMiddle{},
			gowid.RenderFixed{},
		),
		gowid.RenderWithWeight{W: 1},
	}
}
