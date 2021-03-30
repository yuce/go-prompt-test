package main

import (
	"fmt"
	"github.com/gcla/gowid"
	"github.com/gcla/gowid/widgets/fill"
	"github.com/gcla/gowid/widgets/framed"
	"github.com/gcla/gowid/widgets/holder"
	"github.com/gcla/gowid/widgets/pile"
	"github.com/gcla/gowid/widgets/styled"
	"github.com/gcla/gowid/widgets/text"
	"github.com/gcla/gowid/widgets/vpadding"
	hz "github.com/hazelcast/hazelcast-go-client/v4/hazelcast"
	"github.com/hazelcast/hazelcast-go-client/v4/hazelcast/sql"
	"hzsqlcl"
	"time"
)

func main() {
	// connect the client
	/*
		cb := hz.NewClientConfigBuilder()
		cb.Cluster().SetName("jet")
		config, err := cb.Config()
		if err != nil {
			log.Fatal(err)
		}
		client, err = hz.StartNewClientWithConfig(config)
		if err != nil {
			log.Fatal(err)
		}
	*/

	//populateMap(client)

	palette := gowid.Palette{
		"hint":  gowid.MakePaletteEntry(gowid.ColorBlack, gowid.NewUrwidColor("light gray")),
		"error": gowid.MakePaletteEntry(gowid.ColorWhite, gowid.ColorRed),
		"line":  gowid.MakeStyledPaletteEntry(gowid.NewUrwidColor("black"), gowid.NewUrwidColor("light gray"), gowid.StyleBold),
	}
	hline := styled.New(fill.New('-'), gowid.MakePaletteRef("line"))
	txt := text.NewFromContentExt(hzsqlcl.CreateHintMessage("Hit tab to auto-complete"),
		text.Options{
			Align: gowid.HAlignLeft{},
		},
	)
	statusBar := styled.New(txt, gowid.MakePaletteRef("hint"))
	resultWidget := text.NewFromContentExt(hzsqlcl.CreateHintMessage(""),
		text.Options{
			Align: gowid.HAlignLeft{},
		},
	)
	editBox := hzsqlcl.NewEditBox(resultWidget)
	flow := gowid.RenderFlow{}
	pilew := hzsqlcl.NewResizeablePile([]gowid.IContainerWidget{
		&gowid.ContainerWidget{IWidget: resultWidget, D: gowid.RenderWithWeight{2}},
		&gowid.ContainerWidget{vpadding.New(
			pile.New([]gowid.IContainerWidget{
				&gowid.ContainerWidget{IWidget: hline, D: gowid.RenderWithUnits{U: 1}},
				&gowid.ContainerWidget{IWidget: editBox, D: flow},
				&gowid.ContainerWidget{IWidget: statusBar, D: flow},
			}),
			gowid.VAlignBottom{}, flow,
		), flow},
	})
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

	go func() {
		for {
			time.Sleep(2 * time.Second)
			txt.SetContent(app, hzsqlcl.CreateHintMessage("create mapping MAPPING_NAME MAPPINT TYPE"))
			app.Redraw()
			time.Sleep(2 * time.Second)
			txt.SetContent(app, hzsqlcl.CreateErrorMessage("ERROR: connection to the server was lost"))
			app.Redraw()
			time.Sleep(2 * time.Second)
			txt.SetContent(app, hzsqlcl.CreateHintMessage(""))
			app.Redraw()
		}
	}()

	app.SimpleMainLoop()
}

/*
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
*/

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
				fmt.Print(" | " + rowMetadata.Column(i).Name())
			}
			fmt.Println()
		}
		counter++
		for i := 0; i < columnCount; i++ {
			fmt.Print("Value: ")
			// column := rowMetadata.Column(i)
			// column.Type()
			fmt.Print(row.ValueAtIndex(i))
			fmt.Print(" ")

		}
		fmt.Println()
	}
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
