package iadd

import (
	"fmt"
	"sort"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/nsf/termbox-go"
)

const coldef = termbox.ColorDefault

type StatusLine struct {
	statusCode byte
	file       string
}

func (s *StatusLine) String() string {
	return fmt.Sprintf("       %s %s", string(s.statusCode), s.file)
}

type Group struct {
	header      string
	statusLines []StatusLine
	fg          termbox.Attribute
}

func (g *Group) Lines() []line {
	if len(g.statusLines) == 0 {
		return make([]line, 0)
	}
	lines := make([]line, 0)
	for _, s := range g.statusLines {
		lines = append(lines, line{
			String: s.String(),
			Fg:     g.fg,
			Bg:     coldef,
		})
	}
	return lines
}

func (g *Group) HasStatusLines() bool {
	return len(g.statusLines) > 0
}

type ScreenContent struct {
	stagingGroup    Group
	worktreeGroup   Group
	untrackingGroup Group
	currentIdx      int
	lineLen         int
	statusLines     []StatusLine
}

func NewScreenContent() *ScreenContent {
	return &ScreenContent{
		stagingGroup: Group{
			header:      "Changes to be committed:",
			statusLines: make([]StatusLine, 0),
			fg:          termbox.ColorGreen},
		worktreeGroup: Group{
			header:      "Changes not staged for commit:",
			statusLines: make([]StatusLine, 0),
			fg:          termbox.ColorRed},
		untrackingGroup: Group{
			header:      "Untracked files:",
			statusLines: make([]StatusLine, 0),
			fg:          termbox.ColorRed},
	}
}

func (s *ScreenContent) LoadCurrentStatus() {
	s.stagingGroup = Group{
		header:      "Changes to be committed:",
		statusLines: make([]StatusLine, 0),
		fg:          termbox.ColorGreen}
	s.worktreeGroup = Group{
		header:      "Changes not staged for commit:",
		statusLines: make([]StatusLine, 0),
		fg:          termbox.ColorRed}
	s.untrackingGroup = Group{
		header:      "Untracked files:",
		statusLines: make([]StatusLine, 0),
		fg:          termbox.ColorRed}

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
			s.untrackingGroup.statusLines = append(s.untrackingGroup.statusLines, StatusLine{(byte)(fileStatus.Worktree), file})
			continue
		}

		if fileStatus.Staging == git.Untracked {
			panic("Abnormal status: unexpected untracking file " + file)
		}

		if fileStatus.Worktree != git.Unmodified {
			s.worktreeGroup.statusLines = append(s.worktreeGroup.statusLines, StatusLine{(byte)(fileStatus.Worktree), file})
		}

		if fileStatus.Staging != git.Unmodified {
			s.stagingGroup.statusLines = append(s.stagingGroup.statusLines, StatusLine{(byte)(fileStatus.Staging), file})
		}
		//TODO: handle copied, UpdatedButUnmarged
	}
	sort := func(l []StatusLine) {
		sort.SliceStable(l, func(i, j int) bool {
			return strings.Compare((l)[i].file, (l)[j].file) == -1
		})
	}
	sort(s.stagingGroup.statusLines)
	sort(s.worktreeGroup.statusLines)
	sort(s.untrackingGroup.statusLines)

	s.lineLen = len(s.stagingGroup.statusLines) + len(s.worktreeGroup.statusLines) + len(s.untrackingGroup.statusLines)
	s.statusLines = s.stagingGroup.statusLines
	s.statusLines = append(s.statusLines, s.worktreeGroup.statusLines...)
	s.statusLines = append(s.statusLines, s.untrackingGroup.statusLines...)
}

func (s *ScreenContent) Lines() []line {
	a := make([]line, 0)
	a = append(a, s.stagingGroup.Lines()...)
	a = append(a, s.worktreeGroup.Lines()...)
	a = append(a, s.untrackingGroup.Lines()...)

	a[s.currentIdx].Bg = termbox.ColorYellow

	cntGroupHeader := 0
	if s.stagingGroup.HasStatusLines() {
		a = Insert(a, line{String: "hoge", Fg: coldef, Bg: coldef}, 0)
		cntGroupHeader = cntGroupHeader + 1 + len(s.stagingGroup.statusLines)
	}
	if s.worktreeGroup.HasStatusLines() {
		a = Insert(a, line{String: "hoge", Fg: coldef, Bg: coldef}, cntGroupHeader)
		cntGroupHeader = cntGroupHeader + 1 + len(s.worktreeGroup.statusLines)
	}
	if s.untrackingGroup.HasStatusLines() {
		a = Insert(a, line{String: "hoge", Fg: coldef, Bg: coldef}, cntGroupHeader)
	}

	return a
}

func Insert(lines []line, l line, at int) []line {
	latter := append([]line{l}, lines[at:]...)
	lines = append(lines[:at], latter...)
	return lines
}

func (s *ScreenContent) Down() {
	if s.currentIdx+1 >= s.lineLen {
		return
	}
	s.currentIdx++
}

func (s *ScreenContent) Up() {
	if s.currentIdx == 0 {
		return
	}
	s.currentIdx--
}

func (s *ScreenContent) Add() {
	f := s.statusLines[s.currentIdx].file
	r, e := git.PlainOpen(".")
	if e != nil {
		panic(e)
	}
	w, e := r.Worktree()
	if e != nil {
		panic(e)
	}
	w.Add(f)
}

func (s *ScreenContent) Revert() {
	f := s.statusLines[s.currentIdx].file
	r, e := git.PlainOpen(".")
	if e != nil {
		panic(e)
	}
	i, e := r.Storer.Index()
	if e != nil {
		panic(e)
	}
	_, err := i.Remove(f)
	if e != nil {
		panic(err)
	}
	r.Storer.SetIndex(i)
}

type line struct {
	String string
	Fg     termbox.Attribute
	Bg     termbox.Attribute
}
