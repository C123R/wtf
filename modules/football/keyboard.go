package football

import (
	"github.com/gdamore/tcell"
)

func (widget *Widget) initializeKeyboardControls() {
	widget.InitializeCommonControls(widget.Refresh)

	widget.SetKeyboardKey(tcell.KeyRight, widget.next, "Select next item")
	widget.SetKeyboardKey(tcell.KeyLeft, widget.prev, "Select previous item")
}

func (widget *Widget) next() {
	offset++
	widget.Refresh()
}

func (widget *Widget) prev() {
	offset--
	widget.Refresh()
}
