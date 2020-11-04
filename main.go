package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/nsf/termbox-go"
)

const coldef = termbox.ColorDefault

func drawText(x, y int, text string, fg, bg termbox.Attribute) {
	for i, c := range text {
		termbox.SetCell(x+i, y, c, fg, bg)
	}
}

type statusLine struct {
	statusCode byte
	file       string
}

func (s *statusLine) string() string {
	return fmt.Sprintf("       %s %s", string(s.statusCode), s.file)
}

type group struct {
	header      string
	statusLines []statusLine
	fg          termbox.Attribute
}

func (g *group) lines() []line {
	if len(g.statusLines) == 0 {
		return make([]line, 0)
	}
	lines := make([]line, 0)
	for _, s := range g.statusLines {
		lines = append(lines, line{
			string: s.string(),
			fg:     g.fg,
			bg:     coldef,
		})
	}
	return lines
}

func (g *group) hasStatusLines() bool {
	return len(g.statusLines) > 0
}

type screenContent struct {
	stagingGroup    group
	worktreeGroup   group
	untrackingGroup group
	currentIdx      int
	lineLen         int
}

func newScreenContent() *screenContent {
	return &screenContent{
		stagingGroup: group{
			header:      "Changes to be committed:",
			statusLines: make([]statusLine, 0),
			fg:          termbox.ColorGreen},
		worktreeGroup: group{
			header:      "Changes not staged for commit:",
			statusLines: make([]statusLine, 0),
			fg:          termbox.ColorRed},
		untrackingGroup: group{
			header:      "Untracked files:",
			statusLines: make([]statusLine, 0),
			fg:          termbox.ColorRed},
	}
}
func (s *screenContent) loadCurrentStatus() {
	r, e := git.PlainOpen(".")
	if e != nil {
		panic(e)
	}
	w, e := r.Worktree()
	if e != nil {
		panic(e)
	}
	st, e := w.Status()
	if e != nil {
		panic(e)
	}

	for file, fileStatus := range st {
		//just one line for Untracked file
		if fileStatus.Worktree == git.Untracked {
			s.untrackingGroup.statusLines = append(s.untrackingGroup.statusLines, statusLine{(byte)(fileStatus.Worktree), file})
			continue
		}

		if fileStatus.Staging == git.Untracked {
			panic("Abnormal status: unexpected untracking file " + file)
		}

		if fileStatus.Worktree != git.Unmodified {
			s.worktreeGroup.statusLines = append(s.worktreeGroup.statusLines, statusLine{(byte)(fileStatus.Worktree), file})
		}

		if fileStatus.Staging != git.Unmodified {
			s.stagingGroup.statusLines = append(s.stagingGroup.statusLines, statusLine{(byte)(fileStatus.Staging), file})
		}
		//TODO: handle copied, UpdatedButUnmarged
	}
	sort := func(l []statusLine) {
		sort.SliceStable(l, func(i, j int) bool {
			return strings.Compare((l)[i].file, (l)[j].file) == -1
		})
	}
	sort(s.stagingGroup.statusLines)
	sort(s.worktreeGroup.statusLines)
	sort(s.untrackingGroup.statusLines)

	s.lineLen = len(s.stagingGroup.statusLines) + len(s.worktreeGroup.statusLines) + len(s.untrackingGroup.statusLines)
}

func (s *screenContent) lines() []line {
	a := make([]line, 0)
	a = append(a, s.stagingGroup.lines()...)
	a = append(a, s.worktreeGroup.lines()...)
	a = append(a, s.untrackingGroup.lines()...)

	a[s.currentIdx].bg = termbox.ColorYellow

	cntGroupHeader := 0
	if s.stagingGroup.hasStatusLines() {
		a = insert(a, line{string: "hoge", fg: coldef, bg: coldef}, 0)
		cntGroupHeader = cntGroupHeader + 1 + len(s.stagingGroup.statusLines)
	}
	if s.worktreeGroup.hasStatusLines() {
		a = insert(a, line{string: "hoge", fg: coldef, bg: coldef}, cntGroupHeader)
		cntGroupHeader = cntGroupHeader + 1 + len(s.worktreeGroup.statusLines)
	}
	if s.untrackingGroup.hasStatusLines() {
		a = insert(a, line{string: "hoge", fg: coldef, bg: coldef}, cntGroupHeader)
	}

	return a
}

func insert(lines []line, l line, at int) []line {
	latter := append([]line{l}, lines[at:]...)
	lines = append(lines[:at], latter...)
	return lines
}

func (s *screenContent) down() {
	if s.currentIdx+1 >= s.lineLen {
		return
	}
	s.currentIdx++
}

func (s *screenContent) up() {
	if s.currentIdx == 0 {
		return
	}
	s.currentIdx--
}

type line struct {
	string string
	fg     termbox.Attribute
	bg     termbox.Attribute
}

func drawStatus(sc *screenContent) {
	for i, l := range sc.lines() {
		drawText(0, i, l.string, l.fg, l.bg)
	}
}

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	sc := newScreenContent()
	sc.loadCurrentStatus()
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
			}
			switch ev.Key {
			case termbox.KeyArrowDown:
				sc.down()
			case termbox.KeyArrowUp:
				sc.up()
			}
		}
	}
}
