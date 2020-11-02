package main

import (
	"fmt"
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

type statusLine struct {
	statusCode byte
	file       string
}

func (s *statusLine) string() string {
	return fmt.Sprintf("%s %s", string(s.statusCode), s.file)
}

func drawStatus() {
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

	statusLines := make([]statusLine, 0)
	for file, fileStatus := range s {
		//just one line for Untracked file
		if fileStatus.Worktree == git.Untracked {
			statusLines = append(statusLines, statusLine{(byte)(fileStatus.Worktree), file})
			continue
		}

		if fileStatus.Worktree != git.Unmodified {
			statusLines = append(statusLines, statusLine{(byte)(fileStatus.Worktree), file})
		}

		if fileStatus.Staging != git.Unmodified {
			statusLines = append(statusLines, statusLine{(byte)(fileStatus.Staging), file})
		}
	}
	for _, s := range statusLines {
		log.Println(s.string())
	}
}

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()
	drawStatus()
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
