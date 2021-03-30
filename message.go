package hzsqlcl

import (
	"github.com/gcla/gowid"
	"github.com/gcla/gowid/widgets/text"
)

func CreateMessage(msg string, color string) *text.Content {
	txtSegment := text.StyledContent(msg, gowid.MakePaletteRef(color))
	return text.NewContent([]text.ContentSegment{
		txtSegment,
	})
}

func CreateHintMessage(msg string) *text.Content {
	return CreateMessage(msg, "hint")
}

func CreateErrorMessage(msg string) *text.Content {
	return CreateMessage(msg, "error")
}
