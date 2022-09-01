package gui

import (
	g "github.com/AllenDang/giu"
)

func loop() {
	mainWindow := g.SingleWindow()
	linkWindow := g.Window("Set Download Link").Flags(g.WindowFlagsNoResize)

	showMainWindow(mainWindow)
	showLinkWindow(linkWindow, mainWindow)
	showErrors()
}
