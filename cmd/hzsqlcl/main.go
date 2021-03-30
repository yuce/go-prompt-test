package main

import (
	"fmt"
	"github.com/gcla/gowid"
	"github.com/gcla/gowid/widgets/edit"
	"github.com/gcla/gowid/widgets/fill"
	"github.com/gcla/gowid/widgets/framed"
	"github.com/gcla/gowid/widgets/holder"
	"github.com/gcla/gowid/widgets/pile"
	"github.com/gcla/gowid/widgets/styled"
	"github.com/gcla/gowid/widgets/text"
	"github.com/gcla/gowid/widgets/vpadding"
	"github.com/gdamore/tcell"
	hz "github.com/hazelcast/hazelcast-go-client/v4/hazelcast"
	"github.com/hazelcast/hazelcast-go-client/v4/hazelcast/property"
	"github.com/hazelcast/hazelcast-go-client/v4/hazelcast/sql"
	"log"
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

type EditBox struct {
	gowid.IWidget
	resultWidget *text.Widget
}

func NewEditBox(resultWidget *text.Widget) *EditBox {
	editWidget := edit.New(edit.Options{Caption: "SQL> "})
	return &EditBox{
		IWidget:      editWidget,
		resultWidget: resultWidget,
	}
}

func handleSqlResult(result sql.Result) string {
	var res string
	rows := result.Rows()

	counter := 0

	for rows.HasNext() {
		row := rows.Next()
		rowMetadata := row.Metadata()
		columnCount := rowMetadata.ColumnCount()
		// print column names once
		if counter == 0 {
			for i := 0; i < columnCount; i++ {
				res += " | " + rowMetadata.Column(i).Name()
			}
			res += "\n"
		}
		counter++
		for i := 0; i < columnCount; i++ {
			// column := rowMetadata.Column(i)
			// column.Type()
			res += fmt.Sprintf("   %v", row.ValueAtIndex(i))
			res += "   "

		}
		res += "\n"
	}
	return res
}

func (w *EditBox) UserInput(ev interface{}, size gowid.IRenderSize, focus gowid.Selector, app gowid.IApp) bool {
	res := true
	if evk, ok := ev.(*tcell.EventKey); ok {
		switch evk.Key() {
		case tcell.KeyEnter:
			t := w.IWidget.(*edit.Widget).Text()
			result, err := client.ExecuteSQL(t)
			var responseMessage string
			if err != nil {
				responseMessage = fmt.Sprintf("%s", err)
				w.resultWidget.SetContent(app, createErrorMessage(responseMessage))
			} else {
				responseMessage = handleSqlResult(result)
				w.resultWidget.SetContent(app, createHintMessage(responseMessage))
			}

			//w.IWidget = text.New(fmt.Sprintf("Executed SQL: %s.", t))
			w.IWidget.(*edit.Widget).SetText("", app)
		default:
			res = w.IWidget.UserInput(ev, size, focus, app)
		}
	}
	return res
}

type ResizeablePileWidget struct {
	*pile.Widget
	offset int
}

func NewResizeablePile(widgets []gowid.IContainerWidget) *ResizeablePileWidget {
	res := &ResizeablePileWidget{}
	res.Widget = pile.New(widgets)
	return res
}

// for testing
func populateMap(client *hz.Client) {
	someMap, _ := client.GetMap("someMap")

	_ = someMap.Clear()

	_, _ = someMap.Put(1, "hi")
	_, _ = someMap.Put(2, "hi2")
	_, _ = someMap.Put(3, "hi3")
	_, _ = someMap.Put(4, "hi4")
	_, _ = someMap.Put(5, "hi5")
	_, _ = someMap.Put(6, "hi6")
}

var client *hz.Client

func main() {
	// connect the client
	cb := hz.NewClientConfigBuilder()
	cb.Cluster().SetName("jet")
	cb.SetProperty(property.LoggingLevel, "error")

	config, err := cb.Config()
	if err != nil {
		log.Fatal(err)
	}
	client, err = hz.StartNewClientWithConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	populateMap(client)

	palette := gowid.Palette{
		"banner": gowid.MakePaletteEntry(gowid.NewUrwidColor("light gray"), gowid.ColorDefault),
		"error":  gowid.MakePaletteEntry(gowid.ColorRed, gowid.ColorDefault),
		"line":   gowid.MakeStyledPaletteEntry(gowid.NewUrwidColor("black"), gowid.NewUrwidColor("light gray"), gowid.StyleBold),
	}
	hline := styled.New(fill.New('-'), gowid.MakePaletteRef("line"))
	txt := text.NewFromContentExt(createHintMessage("Hit tab to auto-complete"),
		text.Options{
			Align: gowid.HAlignLeft{},
		},
	)
	resultWidget := text.NewFromContentExt(createHintMessage(""),
		text.Options{
			Align: gowid.HAlignLeft{},
		},
	)
	editBox := NewEditBox(resultWidget)
	flow := gowid.RenderFlow{}
	pilew := NewResizeablePile([]gowid.IContainerWidget{
		&gowid.ContainerWidget{IWidget: resultWidget, D: gowid.RenderWithWeight{2}},
		//&gowid.ContainerWidget{IWidget: divider.NewBlank(), D: flow},
		&gowid.ContainerWidget{vpadding.New(
			pile.New([]gowid.IContainerWidget{
				&gowid.ContainerWidget{IWidget: hline, D: gowid.RenderWithUnits{U: 1}},
				&gowid.ContainerWidget{IWidget: editBox, D: flow},
				&gowid.ContainerWidget{IWidget: txt, D: flow},
			}),
			gowid.VAlignBottom{}, flow,
		), flow},
		//&gowid.ContainerWidget{IWidget: editBox, D: flow},
		//&gowid.ContainerWidget{IWidget: txt, D: flow},
	})
	/*
		view := vpadding.New(
			pile.New([]gowid.IContainerWidget{
				&gowid.ContainerWidget{IWidget: resultWidget, D: flow},
				&gowid.ContainerWidget{IWidget: divider.NewBlank(), D: flow},
				&gowid.ContainerWidget{IWidget: editBox, D: flow},
				&gowid.ContainerWidget{IWidget: txt, D: flow},
			}),
			gowid.VAlignBottom{}, flow,
		)
	*/
	tw := text.New(" localhost:5701 ")
	//twi := styled.New(tw, gowid.MakePaletteRef("invred"))
	twp := holder.New(tw)
	view := framed.New(pilew, framed.Options{
		Frame:       framed.UnicodeFrame,
		TitleWidget: twp,
	})
	app, _ := gowid.NewApp(gowid.AppArgs{
		View:    view,
		Palette: palette,
	})

	/*
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
	*/

	app.SimpleMainLoop()
}
