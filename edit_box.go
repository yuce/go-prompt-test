package hzsqlcl

import (
	"container/ring"
	"fmt"
	"strings"

	"github.com/gcla/gowid"
	"github.com/gcla/gowid/widgets/edit"
	"github.com/gcla/gowid/widgets/text"
	"github.com/gdamore/tcell"
)

type EditBox struct {
	gowid.IWidget
	resultWidget *text.Widget
	handler      EditBoxHandler
	history      *ring.Ring
	lastHistory  *ring.Ring
}

type EditBoxHandler func(app gowid.IApp, widget gowid.IWidget, enteredText string)

func NewEditBox(resultWidget *text.Widget, handler EditBoxHandler) *EditBox {
	editWidget := edit.New(edit.Options{Caption: "SQL> "})
	commandHistory := ring.New(30)
	commandHistory.Value = "<end>"
	return &EditBox{
		IWidget:      editWidget,
		resultWidget: resultWidget,
		handler:      handler,
		history:      commandHistory,
		lastHistory:  commandHistory,
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
					w.lastHistory = w.lastHistory.Next()
					w.lastHistory.Value = t
					w.history = w.lastHistory

					w.handler(app, w.resultWidget, t)
				}
			} else {
				inputWidget := w.IWidget.(*edit.Widget)
				inputWidget.SetText(t+"\n", app)
				inputWidget.SetCursorPos(inputWidget.CursorPos()+1, app)
			}

		case tcell.KeyUp:
			command := fmt.Sprint(w.history.Value)
			inputWidget := w.IWidget.(*edit.Widget)
			inputWidget.SetText(command, app)
			inputWidget.SetCursorPos(len(command), app)
			if fmt.Sprint(w.history.Value) != "<end>" {
				w.history = w.history.Prev()
			}

		case tcell.KeyDown:
			w.history = w.history.Next()
			command := fmt.Sprint(w.history.Value)
			if command == "<nil>" {
				command = ""
				w.history = w.history.Prev()
			}

			inputWidget := w.IWidget.(*edit.Widget)
			inputWidget.SetText(command, app)
			inputWidget.SetCursorPos(len(command), app)

		case tcell.KeyCtrlL:
			inputWidget := w.IWidget.(*edit.Widget)
			inputWidget.SetText("", app)
			inputWidget.SetCursorPos(0, app)

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

func (w *EditBox) SetText(app gowid.IApp, text string) {
	w.IWidget.(*edit.Widget).SetText(text, app)
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
	}
}
