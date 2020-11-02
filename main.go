package main

import (
	"log"

	"github.com/go-git/go-git/v5"
	"github.com/nsf/termbox-go"
)

const coldef = termbox.ColorDefault

func drawText(x, y int, text string) {
	for i, c := range text {
		termbox.SetCell(x+i, y, c, coldef, coldef)
	}
}

func gitHello() {
	r, e := git.PlainOpen(".")
	if e != nil {
		panic(e)
	}
	w, e := r.Worktree()
	if e != nil {
		panic(e)
	}
	s, e := w.Status()
	if e != nil {
		panic(e)
	}
	log.Printf("w: %v", s)
}

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()
	gitHello()
MAINLOOP:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Ch {
			case 'q':
				break MAINLOOP
			}
		}
		termbox.Clear(coldef, coldef)
		drawText(0, 0, "hoge")
		termbox.Flush()
	}
}
