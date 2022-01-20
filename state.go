package algo

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

type Board [][]BoardStatus

func (b Board) Size() int {
	return len(b)
}

func isStarPos(x, y, size int) bool {
	b := size / 2
	if x == b && y == b {
		return true
	}
	sq := int(math.Sqrt(float64(b)))
	if sq <= 1 {
		return false
	}
	if (x == b || x == sq || size-x-1 == sq) &&
		(y == b || y == sq || size-y-1 == sq) {
		return true
	}
	return false
}

func (board Board) String() string {
	str := "\n   "
	for x := range board {
		str += string('a'+x) + " "
	}
	for x := range board {
		str += fmt.Sprintf("\n%2d ", x+1)
		for y := range board[x] {
			var c string
			switch board[x][y] {
			case BoardStatusBlack:
				c = "\u26AB"
			case BoardStatusWhite:
				c = "\u26AA"
			case BoardStatusEmpty:
				if isStarPos(x, y, board.Size()) {
					c = "\u205C"
				} else {
					c = "\u253C"
				}
				if y < len(board[x])-1 {
					c += "\u2500"
				}
			case BoardStatusForbidden:
				c = "*\u2500"
			default:
				c = "?\u2500"
			}
			str += c
		}
	}
	return str
}

type BoardSize int

const (
	BoardSizeLarge  BoardSize = 19
	BoardSizeMedium BoardSize = 13
	BoardSizeSmall  BoardSize = 9
	BoardSizeMini   BoardSize = 5
)

// BoardStatus fill all board position
type BoardStatus uint8

const (
	BoardStatusEmpty BoardStatus = iota
	BoardStatusForbidden
	BoardStatusBlack
	BoardStatusWhite

	BoardStatusFalse BoardStatus = 0
	BoardStatusTrue              = 1
)

// Space is a search space
type State struct {
	size           BoardSize
	board          Board
	nextMovePlayer Player
}

// NewSpace ...
func NewState(size BoardSize) *State {
	return &State{
		size:           size,
		board:          NewBoard(size),
		nextMovePlayer: PlayerBlack,
	}
}

func NewBoard(size BoardSize) Board {
	board := make([][]BoardStatus, size)
	for i := range board {
		board[i] = make([]BoardStatus, size)
	}
	return board
}

func (state *State) GetBoard() Board {
	return state.board
}

func (state *State) GetLegalActions() (actions []*Action) {
	for x := range state.board {
		for y := range state.board[x] {
			action := NewAction(x, y, state.nextMovePlayer)
			if !state.isForbidden(action) {
				actions = append(actions, action)
			}
		}
	}
	return
}

func (state *State) MoveTo(action *Action) *State {
	ns := NewState(state.size)
	ns.nextMovePlayer = state.nextMovePlayer.next()
	for x := range state.board {
		for y := range state.board[x] {
			ns.board[x][y] = state.board[x][y]
		}
	}
	ns.board[action.x][action.y] = state.nextMovePlayer.BoardStatus()
	ns.clean(int(action.x), int(action.y))
	return ns
}

func (state *State) isForbidden(action *Action) bool {
	if state.board[action.x][action.y] != BoardStatusEmpty {
		return true
	}
	state.board[action.x][action.y] = state.nextMovePlayer.BoardStatus()
	defer func() {
		state.board[action.x][action.y] = BoardStatusEmpty
	}()
	nv := NewBoard(state.size)
	x, y := int(action.x), int(action.y)
	nv[x][y] = BoardStatusTrue
	judge := &Judge{
		sets:   [][2]int{{x, y}},
		status: state.board[action.x][action.y],
		visit:  nv,
		state:  state,
	}
	judge.isAlive(x, y)
	if judge.alive {
		return false
	}
	// deal da jie
	if len(judge.sets) == 1 {
		for _, np := range pRange(judge.state.board.Size(), x, y) {
			judge.sets = [][2]int{{x, y}}
			judge.status = state.nextMovePlayer.next().BoardStatus()
			judge.isAlive(np[0], np[1])
			if !judge.alive {
				return false
			}
		}
	}
	return true
}

func (state *State) hasResult() bool {
	var b, w, t int
	for x := range state.board {
		for y := range state.board[x] {
			t++
			switch state.board[x][y] {
			case BoardStatusBlack:
				b++
			case BoardStatusWhite:
				w++
			}
		}
	}
	return b > t/2 || w > t/2
}

func (state *State) Result() Player {
	var base int
	for x := range state.board {
		for y := range state.board[x] {
			if state.guess(x, y) == PlayerBlack {
				base++
				if base > 185 {
					return PlayerBlack
				}
			}
		}
	}
	return PlayerWhite
}

func (state *State) guess(x, y int) Player {
	// find all edge node, calc which is more
	var base int
	visit := NewBoard(state.size)
	for _, action := range state.findEdges(x, y, visit) {
		if action.player == PlayerBlack {
			base++
		} else {
			base--
		}
	}
	if base >= 0 {
		return PlayerBlack
	}
	return PlayerWhite
}

func (state *State) findEdges(x, y int, visit Board) []*Action {
	visit[x][y] = BoardStatusTrue
	var res []*Action
	for _, np := range pRange(state.board.Size(), x, y) {
		nx, ny := np[0], np[1]
		if visit[nx][ny] == BoardStatusTrue {
			continue
		}
		switch state.board[nx][ny] {
		case BoardStatusBlack:
			visit[nx][ny] = BoardStatusTrue
			res = append(res, NewAction(nx, ny, PlayerBlack))
		case BoardStatusWhite:
			visit[nx][ny] = BoardStatusTrue
			res = append(res, NewAction(nx, ny, PlayerWhite))
		default:
			res = append(res, state.findEdges(nx, ny, visit)...)
		}
	}
	return res
}

func (state *State) clean(cx, cy int) {
	visit := NewBoard(state.size)
	for x := range visit {
		for y := range visit[x] {
			if visit[x][y] == BoardStatusFalse {
				visit[x][y] = BoardStatusTrue
				switch state.board[x][y] {
				case BoardStatusEmpty:
				case BoardStatusForbidden:
					state.board[x][y] = BoardStatusEmpty
				case BoardStatusBlack, BoardStatusWhite:
					nv := NewBoard(state.size)
					nv[x][y] = BoardStatusTrue
					judge := &Judge{
						sets:   [][2]int{{x, y}},
						status: state.board[x][y],
						visit:  nv,
						state:  state,
					}
					judge.deal(cx, cy)
				}
			}
		}
	}
}

func (judge *Judge) isAlive(cx, cy int) {
	next := false
	for _, p := range judge.sets {
		for _, np := range pRange(judge.state.board.Size(), p[0], p[1]) {
			x, y := np[0], np[1]
			if judge.visit[x][y] == BoardStatusTrue {
				continue
			}
			switch judge.state.board[x][y] {
			case judge.status:
				judge.visit[x][y] = BoardStatusTrue
				judge.sets = append(judge.sets, [2]int{x, y})
				next = true
			case BoardStatusEmpty, BoardStatusForbidden:
				judge.visit[x][y] = BoardStatusTrue
				judge.alive = true
			}
		}
	}
	if next {
		judge.isAlive(cx, cy)
	}
}

type Judge struct {
	sets   [][2]int // (x, y )
	status BoardStatus

	alive bool
	visit Board
	state *State
}

func (judge *Judge) deal(cx, cy int) {
	judge.isAlive(cx, cy)
	if !judge.alive {
		log.Debug("remove not alive:", judge.sets)
		if len(judge.sets) == 1 {
			p := judge.sets[0]
			x, y := p[0], p[1]
			if cx != x || cy != y {
				judge.state.board[x][y] = BoardStatusForbidden
			} else {
				log.Info("jie is here, keep this, deal the other one")
			}
		} else {
			for _, p := range judge.sets {
				x, y := p[0], p[1]
				judge.state.board[x][y] = BoardStatusEmpty
			}
		}
	}

	log.Debug("judge:", judge.sets)
	log.Debug("visit:", judge.visit)
	log.Debug("board:", judge.state.board)
}

func pRange(size, x, y int) (p [][2]int) {
	if y-1 >= 0 {
		p = append(p, [2]int{x, y - 1})
	}
	if x+1 < size {
		p = append(p, [2]int{x + 1, y})
	}
	if y+1 < size {
		p = append(p, [2]int{x, y + 1})
	}
	if x-1 >= 0 {
		p = append(p, [2]int{x - 1, y})
	}
	return
}

type Player int

const (
	PlayerBlack Player = iota
	PlayerWhite
)

func (p Player) next() Player {
	if p == PlayerBlack {
		return PlayerWhite
	}
	return PlayerBlack
}

func (p Player) BoardStatus() BoardStatus {
	switch p {
	case PlayerBlack:
		return BoardStatusBlack
	case PlayerWhite:
		return BoardStatusWhite
	}
	return BoardStatusEmpty
}

func (p Player) String() string {
	switch p {
	case PlayerBlack:
		return "Black"
	case PlayerWhite:
		return "White"
	}
	return "unknown"
}

type Action struct {
	x uint8
	y uint8

	player Player
}

func NewAction(x, y int, player Player) *Action {
	return &Action{
		x:      uint8(x),
		y:      uint8(y),
		player: player,
	}
}

func (action *Action) Detail() (x, y uint8, player Player) {
	return action.x, action.y, action.player
}

func (action *Action) String() string {
	if action == nil {
		return "nil"
	}
	return fmt.Sprintf("x%dy%dp%d", action.x, action.y, action.player)
}

func (action *Action) FromString(str string) error {
	xi := strings.Index(str, "x")
	yi := strings.Index(str, "y")
	pi := strings.Index(str, "p")
	if xi < 0 || yi < 0 || pi < 0 {
		return fmt.Errorf("invalid string: %s", str)
	}
	x, err := strconv.Atoi(str[xi:yi])
	if err != nil {
		return fmt.Errorf("invalid string: %s %v", str, err)
	}
	y, err := strconv.Atoi(str[yi:pi])
	if err != nil {
		return fmt.Errorf("invalid string: %s %v", str, err)
	}
	action.x = uint8(x)
	action.y = uint8(y)
	switch str[pi:] {
	case "0":
		action.player = PlayerBlack
	case "1":
		action.player = PlayerWhite
	default:
		return fmt.Errorf("invalid p value: %s", str)
	}
	return nil
}
