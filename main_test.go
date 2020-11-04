package main

import "testing"

func TestDown(t *testing.T) {
	sc := newScreenContent()
	sc.worktreeGroup.statusLines = []statusLine{
		{statusCode: 'm', file: "hoge"},
		{statusCode: 'm', file: "fuga"},
	}
	sc.untrackingGroup.statusLines = []statusLine{
		{statusCode: '?', file: "foo"},
		{statusCode: '?', file: "bar"},
	}
	sc.decideInitialCusor()
	if sc.currentGroup != 1 {
		t.Errorf("currentGroup != 1: %d %d", sc.currentGroup, sc.currentIdx)
	}
	sc.down()
	sc.down()

	if sc.currentGroup != 2 {
		t.Errorf("currentGroup != 2: %d %d", sc.currentGroup, sc.currentIdx)
	}
}
