package gui

import (
	g "github.com/AllenDang/giu"
)

func Start() {
	wnd := g.NewMasterWindow("Shyn Download Manager", 800, 600, g.MasterWindowFlagsNotResizable)
	wnd.Run(loop)
}
