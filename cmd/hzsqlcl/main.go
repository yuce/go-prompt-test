package main

import (
	"fmt"
	ui "github.com/VladimirMarkelov/clui"
	"github.com/gcla/gowid"
	"github.com/gcla/gowid/widgets/divider"
	"github.com/gcla/gowid/widgets/pile"
	"github.com/gcla/gowid/widgets/styled"
	"github.com/gcla/gowid/widgets/text"
	"github.com/gcla/gowid/widgets/vpadding"
	"github.com/gdamore/tcell"
	"log"
	"time"
)

var txt *text.Widget

func unhandled(app gowid.IApp, ev interface{}) bool {
	if evk, ok := ev.(*tcell.EventKey); ok {
		switch evk.Rune() {
		case 'q', 'Q':
			app.Quit()
		default:
			txt.SetText(fmt.Sprintf("hello world - %c", evk.Rune()), app)
		}
	}
	return true
}

func createMessage(msg string, color string) *text.Content {
	txtSegment := text.StyledContent(msg, gowid.MakePaletteRef(color))
	return text.NewContent([]text.ContentSegment{
		txtSegment,
	})
}

func createHintMessage(msg string) *text.Content {
	return createMessage(msg, "banner")
}

func createErrorMessage(msg string) *text.Content {
	return createMessage(msg, "error")
}

func main() {
	palette := gowid.Palette{
		"banner": gowid.MakePaletteEntry(gowid.ColorBlack, gowid.NewUrwidColor("light gray")),
		"streak": gowid.MakePaletteEntry(gowid.ColorBlack, gowid.ColorRed),
		"bg":     gowid.MakePaletteEntry(gowid.ColorBlack, gowid.ColorDarkBlue),
		"error":  gowid.MakePaletteEntry(gowid.ColorWhite, gowid.ColorRed),
	}
	txtSegment := text.StyledContent("Hit tab to auto-complete", gowid.MakePaletteRef("banner"))
	txt := text.NewFromContentExt(
		text.NewContent([]text.ContentSegment{
			txtSegment,
		}), text.Options{
			Align: gowid.HAlignLeft{},
		},
	)
	vert := vpadding.New(txt, gowid.VAlignBottom{}, gowid.RenderFlow{})
	app, _ := gowid.NewApp(gowid.AppArgs{
		View:    vert,
		Palette: palette,
	})

	go func() {
		time.Sleep(2 * time.Second)
		hintMsg := createHintMessage("create mapping MAPPING_NAME MAPPINT TYPE")
		txt.SetContent(app, hintMsg)
		app.Redraw()
		time.Sleep(2 * time.Second)
		errorMsg := createErrorMessage("ERROR: connection to the server was lost")
		txt.SetContent(app, errorMsg)
		app.Redraw()
	}()

	app.SimpleMainLoop()
}

func main4() {
	txt = text.New("hello")
	if app, err := gowid.NewApp(gowid.AppArgs{View: txt}); err != nil {
		log.Fatal(err)
	} else {
		app.MainLoop(gowid.UnhandledInputFunc(unhandled))
	}
}

func main3() {
	txt := text.New("hello, world!")
	if app, err := gowid.NewApp(gowid.AppArgs{View: txt}); err != nil {
		log.Fatal(err)
	} else {
		app.SimpleMainLoop()
	}

}

func main2() {
	palette := gowid.Palette{
		"banner":  gowid.MakePaletteEntry(gowid.ColorWhite, gowid.MakeRGBColor("#60d")),
		"streak":  gowid.MakePaletteEntry(gowid.ColorNone, gowid.MakeRGBColor("#60a")),
		"inside":  gowid.MakePaletteEntry(gowid.ColorNone, gowid.MakeRGBColor("#808")),
		"outside": gowid.MakePaletteEntry(gowid.ColorNone, gowid.MakeRGBColor("#a06")),
		"bg":      gowid.MakePaletteEntry(gowid.ColorNone, gowid.MakeRGBColor("#d06")),
	}

	div := divider.NewBlank()
	outside := styled.New(div, gowid.MakePaletteRef("outside"))
	inside := styled.New(div, gowid.MakePaletteRef("inside"))

	helloworld := styled.New(
		text.NewFromContentExt(
			text.NewContent([]text.ContentSegment{
				text.StyledContent("Now, for something completely different.", gowid.MakePaletteRef("banner")),
			}),
			text.Options{
				Align: gowid.HAlignMiddle{},
			},
		),
		gowid.MakePaletteRef("streak"),
	)

	f := gowid.RenderFlow{}

	view := styled.New(
		vpadding.New(
			pile.New([]gowid.IContainerWidget{
				&gowid.ContainerWidget{IWidget: outside, D: f},
				&gowid.ContainerWidget{IWidget: inside, D: f},
				&gowid.ContainerWidget{IWidget: helloworld, D: f},
				&gowid.ContainerWidget{IWidget: inside, D: f},
				&gowid.ContainerWidget{IWidget: outside, D: f},
			}),
			gowid.VAlignMiddle{},
			f),
		gowid.MakePaletteRef("bg"),
	)

	app, _ := gowid.NewApp(gowid.AppArgs{
		View:    view,
		Palette: &palette,
	})

	app.SimpleMainLoop()
}

func main1() {
	ui.InitLibrary()
	defer ui.DeinitLibrary()

	view1 := ui.AddWindow(0, 0, 10, 7, "Hello, World!")
	btnQuit := ui.CreateButton(view1, 15, 4, "Hi", 1)
	btnQuit.OnClick(func(event ui.Event) {
		go ui.Stop()
	})

	ui.AddWindow(10, 5, 10, 10, "Foo")

	ui.MainLoop()
}

/*
func main() {
	app := tview.NewApplication()
	form := tview.NewForm().
		AddInputField("Name", "", 20, nil, nil).
		AddInputField("Field 1 Name", "", 20, nil, nil).
		AddDropDown("Field 1 Type", []string{"String", "Int", "Float", "Boolean"}, 0, nil).
		//AddCheckbox("Age 18+", false, nil).
		//AddPasswordField("Password", "", 10, '*', nil).
		AddInputField("Option Key 1", "", 20, nil, nil).
		AddInputField("Option Value 1", "", 20, nil, nil).
		AddButton("Add Field", nil).
		AddButton("Add Option", nil).
		AddButton("Save", nil).
		AddButton("Cancel", func() {
			app.Stop()
		})
	form.SetBorder(true).SetTitle(" Create Mapping ").SetTitleAlign(tview.AlignCenter)
	if err := app.SetRoot(form, true).SetFocus(form).Run(); err != nil {
		panic(err)
	}
}
*/
