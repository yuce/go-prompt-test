package hzsqlcl

import (
	"github.com/gcla/gowid"
	"github.com/gcla/gowid/widgets/pile"
)

type ResizeablePileWidget struct {
	*pile.Widget
	offset int
}

func NewResizeablePile(widgets []gowid.IContainerWidget) *ResizeablePileWidget {
	res := &ResizeablePileWidget{}
	res.Widget = pile.New(widgets)
	return res
}
