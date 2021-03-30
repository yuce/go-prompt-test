package hzsqlcl

import (
	"github.com/gcla/gowid"
	"github.com/gcla/gowid/widgets/edit"
	"github.com/gcla/gowid/widgets/text"
	"github.com/gdamore/tcell"
	"strings"
)

type EditBox struct {
	gowid.IWidget
	resultWidget *text.Widget
	handler      EditBoxHandler
}

type EditBoxHandler func(app gowid.IApp, widget gowid.IWidget, enteredText string)

func NewEditBox(resultWidget *text.Widget, handler EditBoxHandler) *EditBox {
	editWidget := edit.New(edit.Options{Caption: "SQL> "})
	return &EditBox{
		IWidget:      editWidget,
		resultWidget: resultWidget,
		handler:      handler,
	}
}

func (w *EditBox) UserInput(ev interface{}, size gowid.IRenderSize, focus gowid.Selector, app gowid.IApp) bool {
	res := true
	if evk, ok := ev.(*tcell.EventKey); ok {
		switch evk.Key() {
		case tcell.KeyEnter:
			t := w.IWidget.(*edit.Widget).Text()
			w.IWidget.(*edit.Widget).SetText("", app)
			if w.handler != nil {
				w.handler(app, w.resultWidget, t)
			}
			//w.resultWidget.SetContent(app, CreateHintMessage(t))

			/*
				result, err := client.ExecuteSQL(t)
				var responseMessage string
				if err != nil {
					responseMessage = fmt.Sprintf("could not execute sql %s", err)
					w.resultWidget.SetContent(app, CreateHintMessage(responseMessage))
				} else {
					responseMessage = handleSqlResult(result)
					w.resultWidget.SetContent(app, CreateErrorMessage(responseMessage))
				}

			*/
		case tcell.KeyTAB:
			t := w.IWidget.(*edit.Widget).Text()

			splitted := strings.Split(t, " ")

			if len(splitted) == 1 {
				w.autoComplete(t, app)
			}

		default:
			res = w.IWidget.UserInput(ev, size, focus, app)
		}
	}
	return res
}

func (w *EditBox) autoComplete(t string, app gowid.IApp) {
	if strings.HasPrefix(t, "s") {
		w.IWidget.(*edit.Widget).SetText("select", app)
		w.IWidget.(*edit.Widget).SetCursorPos(len(w.IWidget.(*edit.Widget).Text()), app)
	} else if strings.HasPrefix(t, "S") {
		w.IWidget.(*edit.Widget).SetText("SELECT", app)
		w.IWidget.(*edit.Widget).SetCursorPos(len(w.IWidget.(*edit.Widget).Text()), app)
	} else if strings.HasPrefix(t, "i") {
		w.IWidget.(*edit.Widget).SetText("insert into", app)
		w.IWidget.(*edit.Widget).SetCursorPos(len(w.IWidget.(*edit.Widget).Text()), app)
	} else if strings.HasPrefix(t, "I") {
		w.IWidget.(*edit.Widget).SetText("INSERT INTO", app)
		w.IWidget.(*edit.Widget).SetCursorPos(len(w.IWidget.(*edit.Widget).Text()), app)
	} else if strings.HasPrefix(t, "c") {
		w.IWidget.(*edit.Widget).SetText("create mapping", app)
		w.IWidget.(*edit.Widget).SetCursorPos(len(w.IWidget.(*edit.Widget).Text()), app)
	} else if strings.HasPrefix(t, "C") {
		w.IWidget.(*edit.Widget).SetText("CREATE MAPPING", app)
		w.IWidget.(*edit.Widget).SetCursorPos(len(w.IWidget.(*edit.Widget).Text()), app)
	} else if t == "" {
		helpText := "Welcome! Some available commands are: \n"
		helpText += "SELECT: You can select from a map or a mapping\n"
		helpText += "INSERT INTO: You can insert a data into a map or a mapping\n"
		helpText += "CREATE MAPPING: You can create a mapping\n"
		w.resultWidget.SetText(helpText, app)
	}
}
