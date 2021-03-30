package main

import (
	"bytes"
	"fmt"
	"strings"

	"hzsqlcl"

	"github.com/gcla/gowid"
	"github.com/gcla/gowid/widgets/fill"
	"github.com/gcla/gowid/widgets/framed"
	"github.com/gcla/gowid/widgets/holder"
	"github.com/gcla/gowid/widgets/pile"
	"github.com/gcla/gowid/widgets/styled"
	"github.com/gcla/gowid/widgets/text"
	"github.com/gcla/gowid/widgets/vpadding"
	hz "github.com/hazelcast/hazelcast-go-client/v4/hazelcast"
	"github.com/hazelcast/hazelcast-go-client/v4/hazelcast/property"
	"github.com/hazelcast/hazelcast-go-client/v4/hazelcast/sql"
	"github.com/olekukonko/tablewriter"
	log "github.com/sirupsen/logrus"
)

func createApp(statusBar *hzsqlcl.StatusBar) (*gowid.App, error) {
	var viewHolder *holder.Widget
	var editBox *hzsqlcl.EditBox

	pages := []hzsqlcl.WizardPage{
		hzsqlcl.NewNameAndTypePage(),
		hzsqlcl.NewPageWidget2(),
	}

	palette := gowid.Palette{
		"hint":       gowid.MakePaletteEntry(gowid.ColorBlack, gowid.NewUrwidColor("light gray")),
		"error":      gowid.MakePaletteEntry(gowid.ColorRed, gowid.ColorDefault),
		"line":       gowid.MakeStyledPaletteEntry(gowid.NewUrwidColor("black"), gowid.NewUrwidColor("light gray"), gowid.StyleBold),
		"resultLine": gowid.MakePaletteEntry(gowid.ColorWhite, gowid.ColorDefault),
		"query":      gowid.MakePaletteEntry(gowid.ColorOrange, gowid.ColorDefault),
	}
	hline := styled.New(fill.New('-'), gowid.MakePaletteRef("line"))
	resultWidget := text.NewFromContentExt(hzsqlcl.CreateHintMessage(""),
		text.Options{
			Align: gowid.HAlignLeft{},
		},
	)
	createMappingWizard := hzsqlcl.NewWizard(pages, func(app gowid.IApp, state hzsqlcl.WizardState) {
		if generatedSQL, err := hzsqlcl.CreateSQLForCreateMapping(state); err != nil {
			panic(err)
		} else {
			editBox.SetText(app, generatedSQL)
		}
	})
	editBox = hzsqlcl.NewEditBox(resultWidget, func(app gowid.IApp, resultWidget gowid.IWidget, enteredText string) {
		if enteredText == "w;" {
			createMappingWizard.Open(viewHolder, gowid.RenderWithRatio{R: 0.5}, app)
			return
		}

		trimmedEnteredText := strings.TrimSuffix(enteredText, ";")
		//trimmedEnteredText := strings.TrimPrefix(strings.TrimSuffix(enteredText, ";\n"), "> ")
		//resultWidget.(*text.Widget).SetContent(app, hzsqlcl.CreateResultLineMessage(trimmedEnteredText))
		res, err := client.ExecuteSQL(trimmedEnteredText)

		currentContent := resultWidget.(*text.Widget).Content()

		currentContent.AddAt(currentContent.Length(), text.ContentSegment{Text: "> ", Style: nil})
		currentContent.AddAt(currentContent.Length(), text.ContentSegment{Text: enteredText, Style: gowid.MakePaletteRef("query")})
		currentContent.AddAt(currentContent.Length(), text.ContentSegment{Text: "\n", Style: nil})

		if err != nil {
			errorMessage := hzsqlcl.CreateErrorMessage(err.Error())
			currentContent.AddAt(currentContent.Length(), text.ContentSegment{Text: errorMessage.String(), Style: gowid.MakePaletteRef("error")})
		} else {
			currentContent.AddAt(currentContent.Length(), text.ContentSegment{Text: handleSqlResult(res), Style: gowid.MakePaletteRef("resultLine")})
		}
		resultWidget.(*text.Widget).SetContent(app, currentContent)
	})
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
	view := framed.New(pilew, framed.Options{
		Frame:       framed.UnicodeFrame,
		TitleWidget: holder.New(text.New(" localhost:5701 ")),
	})
	viewHolder = holder.New(view)
	return gowid.NewApp(gowid.AppArgs{
		View:    viewHolder,
		Palette: palette,
	})
}

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

	statusBar := hzsqlcl.NewStatusBar()
	app, err := createApp(statusBar)
	if err != nil {
		log.Fatal(err)
	}

	//go func() {
	//	for {
	//		time.Sleep(2 * time.Second)
	//		statusBar.SetHint(app, "create mapping MAPPING_NAME MAPPINT TYPE")
	//		time.Sleep(2 * time.Second)
	//		statusBar.SetError(app, "ERROR: connection to the server was lost")
	//		time.Sleep(2 * time.Second)
	//		statusBar.Clear(app)
	//	}
	//}()
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
	rows := result.Rows()

	var byteBuffer bytes.Buffer
	//var tempWriter io.StringWriter
	table := tablewriter.NewWriter(&byteBuffer)

	i := 0
	for rows.HasNext() {
		row := rows.Next()
		rowMetadata := row.Metadata()
		columnCount := rowMetadata.ColumnCount()

		tableRow := make([]string, columnCount)
		for columnIndex := 0; columnIndex < columnCount; columnIndex++ {
			if i == 0 {
				tableRow[columnIndex] = rowMetadata.Column(columnIndex).Name()
			} else {
				tableRow[columnIndex] = fmt.Sprintf("%v", row.ValueAtIndex(columnIndex))
			}
		}

		if i == 0 {
			table.SetHeader(tableRow)
		} else {
			table.Append(tableRow)
		}

		i++
	}

	table.Render() // Send output

	return byteBuffer.String()
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
