package tui

import (
	"github.com/rivo/tview"
)

var (
	right      = false
	frameRight *tview.Frame
	app        *tview.Application
)

func Run() {
	app = tview.NewApplication().EnableMouse(true)
}
