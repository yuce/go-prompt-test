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

			if strings.HasSuffix(t, ";") {
				w.IWidget.(*edit.Widget).SetText("", app)
				if w.handler != nil {
					w.handler(app, w.resultWidget, t)
				}
			} else {
				inputWidget := w.IWidget.(*edit.Widget)
				inputWidget.SetText(t + "\n", app)
				inputWidget.SetCursorPos(inputWidget.CursorPos() + 1, app)
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
		default:
			res = w.IWidget.UserInput(ev, size, focus, app)
		}
	}
	return res
}
