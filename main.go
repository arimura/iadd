package main

import (
	"github.com/arimura/iadd"
	"github.com/nsf/termbox-go"
)

const coldef = termbox.ColorDefault

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
	defer termbox.Close()

	sc := iadd.NewScreenContent()
	sc.LoadCurrentStatus()
MAINLOOP:
	for {
		termbox.Clear(coldef, coldef)
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
				sc.Revert()
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
