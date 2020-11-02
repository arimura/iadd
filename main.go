package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/nsf/termbox-go"
)

const coldef = termbox.ColorDefault

func drawText(x, y int, text string, fg termbox.Attribute) {
	for i, c := range text {
		termbox.SetCell(x+i, y, c, fg, coldef)
	}
}

type statusLine struct {
	statusCode byte
	file       string
}

func (s *statusLine) string() string {
	return fmt.Sprintf("       %s %s", string(s.statusCode), s.file)
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

	stagingLines := make([]statusLine, 0)
	worktreeLines := make([]statusLine, 0)
	untrackingLines := make([]statusLine, 0)

	for file, fileStatus := range s {
		//just one line for Untracked file
		if fileStatus.Worktree == git.Untracked {
			untrackingLines = append(untrackingLines, statusLine{(byte)(fileStatus.Worktree), file})
			continue
		}

		if fileStatus.Staging == git.Untracked {
			panic("Abnormal status: unexpected untracking file " + file)
		}

		if fileStatus.Worktree != git.Unmodified {
			worktreeLines = append(worktreeLines, statusLine{(byte)(fileStatus.Worktree), file})
		}

		if fileStatus.Staging != git.Unmodified {
			stagingLines = append(stagingLines, statusLine{(byte)(fileStatus.Staging), file})
		}
		//TODO: handle copied, UpdatedButUnmarged
	}

	sort := func(l []statusLine) {
		sort.SliceStable(l, func(i, j int) bool {
			return strings.Compare((l)[i].file, (l)[j].file) == -1
		})
	}
	sort(stagingLines)
	sort(worktreeLines)
	sort(untrackingLines)

	statingLineStartPoint := len(stagingLines)
	worktreeLineStartPoint := len(worktreeLines)

	drawText(0, 0, "Changes to be committed:", coldef)
	for i, s := range stagingLines {
		drawText(0, 1+i, s.string(), termbox.ColorGreen)
	}

	drawText(0, 1+statingLineStartPoint, "Changes not staged for commit:", coldef)
	for i, s := range worktreeLines {
		drawText(0, 2+i+statingLineStartPoint, s.string(), termbox.ColorRed)
	}

	drawText(0, 2+statingLineStartPoint+worktreeLineStartPoint, "Untracked files:", coldef)
	for i, s := range untrackingLines {
		drawText(0, 3+i+statingLineStartPoint+worktreeLineStartPoint, s.string(), termbox.ColorRed)
	}
}

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()
	termbox.Clear(coldef, coldef)
	drawStatus()
	termbox.Flush()
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
		drawText(0, 0, "hoge", coldef)
		termbox.Flush()
	}
}
