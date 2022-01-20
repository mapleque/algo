package algo

import (
	"testing"
)

var (
	B = BoardStatusBlack
	W = BoardStatusWhite
)

func Test_clean1(t *testing.T) {
	//SetLogLevel(Debug)
	b := [][]int{{0, 1}, {1, 0}, {1, 2}, {2, 1}}
	w := [][]int{{1, 1}}
	e := [][]int{}
	f := w
	testClean(t, b, w, e, f)
}

func Test_clean2(t *testing.T) {
	b := [][]int{{0, 1}, {0, 2}, {1, 0}, {1, 3}, {2, 1}, {2, 2}}
	w := [][]int{{1, 1}, {1, 2}}
	e := w
	f := [][]int{}
	testClean(t, b, w, e, f)
}

func Test_clean3(t *testing.T) {
	b := [][]int{{0, 1}, {1, 0}}
	w := [][]int{{0, 0}}
	e := [][]int{}
	f := w
	testClean(t, b, w, e, f)
}

func Test_clean4(t *testing.T) {
	st := NewState(19)
	st.board[18][18] = BoardStatusBlack
	st.clean(18, 18)
	if st.board[18][18] != BoardStatusBlack {
		t.Error("18-18 should be keep, but:", st.board)
	}
}

func Test_tijie(t *testing.T) {
	st := NewState(19)
	for _, p := range [][]int{{0, 1}, {1, 0}, {1, 2}, {2, 1}} {
		st.board[p[0]][p[1]] = BoardStatusBlack
	}
	for _, p := range [][]int{{1, 1}, {2, 0}, {2, 2}, {3, 1}} {
		st.board[p[0]][p[1]] = BoardStatusWhite
	}
	st.clean(1, 1)
	if st.board[2][1] != BoardStatusForbidden {
		t.Error("2-1 should be clean, but:", st.board)
	}
	if st.board[1][1] != BoardStatusWhite {
		t.Error("1-1 should be keep, but:", st.board)
	}
}

func testClean(t *testing.T, b, w, e, f [][]int) {
	st := NewState(19)
	for _, p := range b {
		st.board[p[0]][p[1]] = BoardStatusBlack
	}
	for _, p := range w {
		st.board[p[0]][p[1]] = BoardStatusWhite
	}
	st.clean(9, 9)
	for _, p := range e {
		if st.board[p[0]][p[1]] != BoardStatusEmpty {
			t.Errorf("%d-%d should be clean to empty, but:%s\n", p[0], p[1], st.GetBoard())
		}
	}
	for _, p := range f {
		if st.board[p[0]][p[1]] != BoardStatusForbidden {
			t.Errorf("%d-%d should clean to forbidden, but:%s\n", p[0], p[1], st.GetBoard())
		}
	}
}

func Test_pRange(t *testing.T) {
	assertPRange(t, 0, 0, [][]int{{1, 0}, {0, 1}})
	assertPRange(t, 0, 1, [][]int{{0, 0}, {1, 1}, {0, 2}})
	assertPRange(t, 1, 0, [][]int{{2, 0}, {1, 1}, {0, 0}})
	assertPRange(t, 0, 18, [][]int{{0, 17}, {1, 18}})
	assertPRange(t, 18, 0, [][]int{{18, 1}, {17, 0}})
	assertPRange(t, 18, 18, [][]int{{18, 17}, {17, 18}})
	assertPRange(t, 1, 1, [][]int{{1, 0}, {2, 1}, {1, 2}, {0, 1}})
}

func assertPRange(t *testing.T, x, y int, e [][]int) {
	r := pRange(19, x, y)
	if len(r) != len(e) {
		t.Error("test", x, y, "should:", e, "but:", r)
	}
	for i, p := range r {
		if p[0] != e[i][0] || p[1] != e[i][1] {
			t.Error("test", x, y, "should:", e, "but:", r)
		}
	}
}
