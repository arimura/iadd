package iadd

import "testing"

func TestDown(t *testing.T) {
	sc := NewScreenContent()
	sc.worktreeGroup.statusLines = []StatusLine{
		{statusCode: 'm', file: "hoge"},
		{statusCode: 'm', file: "fuga"},
	}
	sc.untrackingGroup.statusLines = []StatusLine{
		{statusCode: '?', file: "foo"},
		{statusCode: '?', file: "bar"},
	}
	if sc.currentGroup != 1 {
		t.Errorf("currentGroup != 1: %d %d", sc.currentGroup, sc.currentIdx)
	}
	sc.Down()
	sc.Down()

	if sc.currentGroup != 2 {
		t.Errorf("currentGroup != 2: %d %d", sc.currentGroup, sc.currentIdx)
	}
}
