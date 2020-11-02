package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/nsf/termbox-go"
)

const coldef = termbox.ColorDefault

var sortMap = map[byte]int{
	' ': 100,
	'?': 3,
	'M': 2,
	'A': 1,
	'D': 2,
	'R': 1,
	'C': 100,
	'U': 100,
}

func drawText(x, y int, text string) {
	for i, c := range text {
		termbox.SetCell(x+i, y, c, coldef, coldef)
	}
}

type group int8

const (
	placeStagingHeader group = iota
	placeStaging
	placeWorkintreeHeader
	placeWorkintree
	placeUntrackedHeader
	placeUntracked
)

type statusLine struct {
	rowGroup   group
	statusCode byte
	file       string
}

func (s *statusLine) string() string {
	return fmt.Sprintf("       %s %s", string(s.statusCode), s.file)
}

func (s *statusLine) group() group {
	return s.group()
}

type row interface {
	string() string
	group() group
}

type header struct {
	label string
}

func (h *header) string() string {
	return h.label
}
func (h *header) group() group {
	return h.group()
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

	statusLines := make([]row, 0)
	statusLines = append(statusLines, &header{label: "Changes to be committed:"})

	stagingLines := make([]statusLine, 0)
	worktreeLines := make([]statusLine, 0)
	untrackingLines := make([]statusLine, 0)

	for file, fileStatus := range s {
		//just one line for Untracked file
		if fileStatus.Worktree == git.Untracked {
			untrackingLines = append(untrackingLines, statusLine{placeUntracked, (byte)(fileStatus.Worktree), file})
			continue
		}

		if fileStatus.Staging == git.Untracked {
			panic("Abnormal status: unexpected untracking file " + file)
		}

		if fileStatus.Worktree != git.Unmodified {
			worktreeLines = append(worktreeLines, statusLine{placeWorkintree, (byte)(fileStatus.Worktree), file})
		}

		if fileStatus.Staging != git.Unmodified {
			stagingLines = append(stagingLines, statusLine{placeStaging, (byte)(fileStatus.Staging), file})
		}
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

	drawText(0, 0, "Changes to be committed:")
	for i, s := range stagingLines {
		drawText(0, 1+i, s.string())
	}

	drawText(0, 1+statingLineStartPoint, "Changes not staged for commit:")
	for i, s := range worktreeLines {
		drawText(0, 2+i+statingLineStartPoint, s.string())
	}

	drawText(0, 2+statingLineStartPoint+worktreeLineStartPoint, "Untracked files:")
	for i, s := range untrackingLines {
		drawText(0, 3+i+statingLineStartPoint+worktreeLineStartPoint, s.string())
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
		drawText(0, 0, "hoge")
		termbox.Flush()
	}
}
