package iadd

import (
	"fmt"
	"sort"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/nsf/termbox-go"
)

const Coldef = termbox.ColorDefault

type status struct {
	statusCode byte
	file       string
}

func (s *status) String() string {
	return fmt.Sprintf("       %s %s", string(s.statusCode), s.file)
}

type group struct {
	statuses []status
	fg       termbox.Attribute
}

func (g *group) Lines() []line {
	if len(g.statuses) == 0 {
		return make([]line, 0)
	}
	lines := make([]line, 0)
	for _, s := range g.statuses {
		lines = append(lines, line{
			String: s.String(),
			Fg:     g.fg,
			Bg:     Coldef,
		})
	}
	return lines
}

func (g *group) sortStatuses(){
	sort.SliceStable(g.statuses, func(i, j int) bool {
		return strings.Compare((g.statuses)[i].file, (g.statuses)[j].file) == -1
	})
}

func (g *group) HasStatuses() bool {
	return len(g.statuses) > 0
}

type ScreenContent struct {
	indexGroup      group
	worktreeGroup   group
	untrackingGroup group
	currentIdx      int
	lineLen         int
	statuses        []status
}

func NewScreenContent() *ScreenContent {
	return &ScreenContent{
		indexGroup:      *newIndexGroup(),
		worktreeGroup:   *newWorktreeGroup(),
		untrackingGroup: *newUntrackingGroup(),
	}
}

func (s *ScreenContent) LoadCurrentStatus() {
	s.indexGroup = *newIndexGroup()
	s.worktreeGroup = *newWorktreeGroup()
	s.untrackingGroup = *newUntrackingGroup()

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
			s.untrackingGroup.statuses = append(s.untrackingGroup.statuses, status{(byte)(fileStatus.Worktree), file})
			continue
		}

		if fileStatus.Staging == git.Untracked {
			panic("Abnormal status: unexpected untracking file " + file)
		}

		if fileStatus.Worktree != git.Unmodified {
			s.worktreeGroup.statuses = append(s.worktreeGroup.statuses, status{(byte)(fileStatus.Worktree), file})
		}

		if fileStatus.Staging != git.Unmodified {
			s.indexGroup.statuses = append(s.indexGroup.statuses, status{(byte)(fileStatus.Staging), file})
		}
		//TODO: handle copied, UpdatedButUnmarged
	}

	s.indexGroup.sortStatuses()
	s.worktreeGroup.sortStatuses()
	s.untrackingGroup.sortStatuses()

	s.lineLen = len(s.indexGroup.statuses) + len(s.worktreeGroup.statuses) + len(s.untrackingGroup.statuses)
	s.statuses = s.indexGroup.statuses
	s.statuses = append(s.statuses, s.worktreeGroup.statuses...)
	s.statuses = append(s.statuses, s.untrackingGroup.statuses...)
}

func newUntrackingGroup() *group {
	return &group{
		statuses: make([]status, 0),
		fg:       termbox.ColorRed}
}

func newWorktreeGroup() *group {
	return &group{
		statuses: make([]status, 0),
		fg:       termbox.ColorRed}
}

func newIndexGroup() *group {
	return &group{
		statuses: make([]status, 0),
		fg:       termbox.ColorGreen}
}

func (s *ScreenContent) Lines() []line {
	a := make([]line, 0)
	a = append(a, s.indexGroup.Lines()...)
	a = append(a, s.worktreeGroup.Lines()...)
	a = append(a, s.untrackingGroup.Lines()...)

	a[s.currentIdx].Bg = termbox.ColorYellow

	cntGroupHeader := 1
	a = insert(a, line{String: "a: add, r: remove, q: quit", Fg: Coldef, Bg: Coldef}, 0)
	if s.indexGroup.HasStatuses() {
		a = insert(a, line{String: "Changes to be committed:", Fg: Coldef, Bg: Coldef}, 1)
		cntGroupHeader = cntGroupHeader + 1 + len(s.indexGroup.statuses)
	}
	if s.worktreeGroup.HasStatuses() {
		a = insert(a, line{String: "Changes not staged for commit:", Fg: Coldef, Bg: Coldef}, cntGroupHeader)
		cntGroupHeader = cntGroupHeader + 1 + len(s.worktreeGroup.statuses)
	}
	if s.untrackingGroup.HasStatuses() {
		a = insert(a, line{String: "Untracked files:", Fg: Coldef, Bg: Coldef}, cntGroupHeader)
	}

	return a
}

func insert(lines []line, l line, at int) []line {
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
	f := s.statuses[s.currentIdx].file
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

func (s *ScreenContent) Remove() {
	f := s.statuses[s.currentIdx].file
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
