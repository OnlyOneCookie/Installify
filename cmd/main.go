package main

import (
	"Installify/internal/gui"
	"github.com/AllenDang/giu"
)

func main() {
	gui.Setup()
	wnd := giu.NewMasterWindow("Installify", 800, 600, 0)
	wnd.Run(gui.Loop)
}
