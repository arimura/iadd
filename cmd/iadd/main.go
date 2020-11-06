package main

import (
	"github.com/arimura/iadd"
	"github.com/nsf/termbox-go"
)

func drawText(x, y int, text string, fg, bg termbox.Attribute) {
	for i, c := range text {
		termbox.SetCell(x+i, y, c, fg, bg)
	}
}

func drawStatus(sc *iadd.ScreenContent) {
	for i, l := range sc.Lines() {
		drawText(0, i, l.String, l.Fg, l.Bg)
	}
}

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	//TODO: print final status
	defer termbox.Close()

	sc := iadd.NewScreenContent()
	sc.LoadCurrentStatus()
MAINLOOP:
	for {
		termbox.Clear(iadd.Coldef, iadd.Coldef)
		drawStatus(sc)
		termbox.Flush()
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Ch {
			case 'q':
				break MAINLOOP
			case 'a':
				sc.Add()
				sc.LoadCurrentStatus()
			case 'r':
				sc.Remove()
				sc.LoadCurrentStatus()
			}
			switch ev.Key {
			case termbox.KeyArrowDown:
				sc.Down()
			case termbox.KeyArrowUp:
				sc.Up()
			}
		}
	}
}
