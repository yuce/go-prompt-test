package main

import (
	"bytes"
	"fmt"
	"hzsqlcl/components"
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

var (
	DefaultBackground = gowid.NewUrwidColor("white")
	DefaultButton     = gowid.NewUrwidColor("dark blue")
	DefaultButtonText = gowid.NewUrwidColor("yellow")
	DefaultText       = gowid.NewUrwidColor("black")
)

func createApp(statusBar *hzsqlcl.StatusBar) (*gowid.App, error) {
	var viewHolder *holder.Widget
	var editBox *hzsqlcl.EditBox

	pages := []hzsqlcl.WizardPage{
		hzsqlcl.NewNameAndTypePage(),
		hzsqlcl.NewFieldsPage(),
		hzsqlcl.NewOptionsPage(),
	}

	palette := gowid.Palette{
		"hint":       gowid.MakePaletteEntry(gowid.ColorBlack, gowid.NewUrwidColor("light gray")),
		"error":      gowid.MakePaletteEntry(gowid.ColorRed, gowid.ColorDefault),
		"line":       gowid.MakeStyledPaletteEntry(gowid.NewUrwidColor("black"), gowid.NewUrwidColor("light gray"), gowid.StyleBold),
		"resultLine": gowid.MakePaletteEntry(gowid.ColorWhite, gowid.ColorDefault),
		"query":      gowid.MakePaletteEntry(gowid.ColorLightBlue, gowid.ColorDefault),
		"keyword":    gowid.MakePaletteEntry(gowid.ColorBlue, gowid.ColorDefault),
		"form":       gowid.MakePaletteEntry(gowid.ColorWhite, gowid.ColorBlack),
		"background": gowid.MakePaletteEntry(DefaultText, DefaultBackground),
		"border":     gowid.MakePaletteEntry(DefaultButton, DefaultBackground),
		"button":     gowid.MakePaletteEntry(DefaultButtonText, DefaultButton),
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
		} else if strings.Trim(strings.TrimRight(enteredText, ";"), " \n\t") == "help" {
			currentContent := resultWidget.(*text.Widget).Content()
			currentContent.AddAt(currentContent.Length(), text.ContentSegment{Text: "> ", Style: nil})
			addQuery(currentContent, enteredText)
			currentContent.AddAt(currentContent.Length(), text.ContentSegment{Text: "\n", Style: nil})

			currentContent.AddAt(currentContent.Length(), text.ContentSegment{Text: "Welcome! Some available commands are: \n", Style: gowid.MakePaletteRef("resultLine")})
			currentContent.AddAt(currentContent.Length(), text.ContentSegment{Text: "SELECT: ", Style: gowid.MakePaletteRef("query")})
			currentContent.AddAt(currentContent.Length(), text.ContentSegment{Text: "You can select from a map or a mapping \n", Style: gowid.MakePaletteRef("resultLine")})
			currentContent.AddAt(currentContent.Length(), text.ContentSegment{Text: "INSERT INTO: ", Style: gowid.MakePaletteRef("query")})
			currentContent.AddAt(currentContent.Length(), text.ContentSegment{Text: "You can insert a data into a map or a mapping \n", Style: gowid.MakePaletteRef("resultLine")})
			currentContent.AddAt(currentContent.Length(), text.ContentSegment{Text: "CREATE MAPPING: ", Style: gowid.MakePaletteRef("query")})
			currentContent.AddAt(currentContent.Length(), text.ContentSegment{Text: "You can create a mapping \n", Style: gowid.MakePaletteRef("resultLine")})
			currentContent.AddAt(currentContent.Length(), text.ContentSegment{Text: "CREATE JOB: ", Style: gowid.MakePaletteRef("query")})
			currentContent.AddAt(currentContent.Length(), text.ContentSegment{Text: "You can create a job \n", Style: gowid.MakePaletteRef("resultLine")})

			resultWidget.(*text.Widget).SetContent(app, currentContent)
			return
		}

		trimmedEnteredText := strings.TrimSuffix(enteredText, ";")
		//trimmedEnteredText := strings.TrimPrefix(strings.TrimSuffix(enteredText, ";\n"), "> ")
		//resultWidget.(*text.Widget).SetContent(app, hzsqlcl.CreateResultLineMessage(trimmedEnteredText))
		res, err := client.ExecuteSQL(trimmedEnteredText)

		currentContent := resultWidget.(*text.Widget).Content()

		currentContent.AddAt(currentContent.Length(), text.ContentSegment{Text: "> ", Style: nil})
		addQuery(currentContent, enteredText)
		currentContent.AddAt(currentContent.Length(), text.ContentSegment{Text: "\n", Style: nil})

		if err != nil {
			errorMessage := hzsqlcl.CreateErrorMessage(err.Error())
			currentContent.AddAt(currentContent.Length(), text.ContentSegment{Text: errorMessage.String() + "\n", Style: gowid.MakePaletteRef("error")})
		} else {
			currentContent.AddAt(currentContent.Length(), text.ContentSegment{Text: handleSqlResult(res, enteredText), Style: gowid.MakePaletteRef("resultLine")})
		}
		resultWidget.(*text.Widget).SetContent(app, currentContent)
	})
	flow := gowid.RenderFlow{}
	pilew := components.NewResizeablePile([]gowid.IContainerWidget{
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

func addQuery(currentContent text.IContent, enteredText string) {
	keywords := []string{"select", "insert", "create", "mapping", "job", "type", "options", "int", "varchar", "as", "sink", "into", "from", "options"}

	accumulator := ""

	for i := 0; i < len(enteredText); i++ {
		currentChar := enteredText[i]
		if currentChar == ' ' || currentChar == '\n' || currentChar == '=' || currentChar == ';' || currentChar == ',' || currentChar == '(' || currentChar == ')' {
			var found bool
			for _, k := range keywords {

				if strings.ToLower(accumulator) == k {
					found = true
				}
			}
			if found {
				currentContent.AddAt(currentContent.Length(), text.ContentSegment{Text: accumulator, Style: gowid.MakePaletteRef("keyword")})
			} else {
				currentContent.AddAt(currentContent.Length(), text.ContentSegment{Text: accumulator, Style: gowid.MakePaletteRef("query")})
			}
			if currentChar == ' ' || currentChar == '\n' {
				currentContent.AddAt(currentContent.Length(), text.ContentSegment{Text: string(currentChar), Style: gowid.MakePaletteRef("query")})
			} else {
				currentContent.AddAt(currentContent.Length(), text.ContentSegment{Text: string(currentChar), Style: gowid.MakePaletteRef("resultLine")})
			}
			accumulator = ""
		} else {
			accumulator = accumulator + string(currentChar)
		}

	}
	if accumulator != "" {
		var found bool
		for _, k := range keywords {

			if strings.ToLower(accumulator) == k {
				found = true
			}
		}
		if found {
			currentContent.AddAt(currentContent.Length(), text.ContentSegment{Text: accumulator, Style: gowid.MakePaletteRef("keyword")})
		} else {
			currentContent.AddAt(currentContent.Length(), text.ContentSegment{Text: accumulator, Style: gowid.MakePaletteRef("query")})
		}

	}

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

func handleSqlResult(result sql.Result, query string) string {
	rows := result.Rows()
	if !rows.HasNext() && strings.HasPrefix(strings.ToLower(strings.TrimLeft(query, " \n\t")), "select") {
		return "NO ROWS\n"
	} else if !rows.HasNext() {
		return "OK\n"
	}

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
