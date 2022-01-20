package algo

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
)

type Context struct {
	stop bool
}

// TreeNode is a search tree
type TreeNode struct {
	mux    sync.Mutex
	ckfile string

	ctx    *Context
	parent *TreeNode

	action *Action

	allRollout bool
	uct        float64
	total      int64
	visitTimes int
	result     [2]int

	state    *State
	children []*TreeNode
}

// NewTree ...
func NewTree(ckfile string, size BoardSize) *TreeNode {
	return &TreeNode{
		ckfile: fmt.Sprintf("%s.%d.ck", ckfile, size),
		ctx:    &Context{},
		state:  NewState(size),
	}
}

func (root *TreeNode) newChildFromAction(action *Action) *TreeNode {
	return &TreeNode{
		ctx:    root.ctx,
		parent: root,
		action: action,
		state:  root.state.MoveTo(action),
	}
}

func (root *TreeNode) expand() {
	log.Tracef("expand the node: %s", root)
	if root.children == nil {
		root.children = []*TreeNode{}
	}
	total := 0
	for _, action := range root.state.GetLegalActions() {
		if root.findChild(int(action.x), int(action.y)) == nil {
			root.children = append(
				root.children,
				root.newChildFromAction(action),
			)
			total++
		}
	}
	log.Tracef("expand found %d new ations", total)
	root.updateTotal(int64(total))
}

func (root *TreeNode) updateTotal(total int64) {
	if total == 0 {
		return
	}
	root.total += total
	if root.parent != nil {
		root.parent.updateTotal(total)
	}
}

func (root *TreeNode) updateState() {
	for _, node := range root.children {
		node.state = root.state.MoveTo(node.action)
		node.updateState()
	}
}

func (root *TreeNode) FindChild(x, y int) *TreeNode {
	root.expand()
	return root.findChild(x, y)
}

func (root *TreeNode) findChild(x, y int) *TreeNode {
	for _, node := range root.children {
		if int(node.action.x) == x && int(node.action.y) == y {
			return node
		}
	}
	return nil
}

func (root *TreeNode) NextPlayer() Player {
	return root.state.nextMovePlayer
}

func (root *TreeNode) GetAction() *Action {
	return root.action
}

func (root *TreeNode) GetState() *State {
	return root.state
}

func (root *TreeNode) GetChildren() []*TreeNode {
	return root.children
}

func (root *TreeNode) GetUCT() float64 {
	return root.uct
}

func (root *TreeNode) GetWins() int {
	return root.result[root.state.nextMovePlayer]
}

func (root *TreeNode) GetLoss() int {
	return root.result[root.state.nextMovePlayer.next()]
}

func (root *TreeNode) GetN() int {
	return root.visitTimes
}

func (root *TreeNode) String() string {
	if root == nil {
		return "nil"
	}
	var allRollout int
	if root.allRollout {
		allRollout = 1
	}
	return fmt.Sprintf(
		"id:%s,p:%s,a:%s,n:%d,t:%d,r:%s,u:%d",
		fmt.Sprintf("%p", root),
		fmt.Sprintf("%p", root.parent),
		root.action,
		root.visitTimes,
		root.total,
		fmt.Sprintf("%d-%d", root.result[0], root.result[1]),
		allRollout,
	)
}

func (root *TreeNode) SaveCheckpoint() {
	head := root.head()
	if head.ckfile == "" {
		panic("need ckfile")
	}
	log.Infof("save checkpoint to file: %s", root.ckfile)
	head.mux.Lock()
	defer head.mux.Unlock()
	log.Trace("save checkpoint get lock")
	f, err := os.OpenFile(head.ckfile, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	f.Write([]byte(fmt.Sprintf("%d\n", head.state.size)))
	write(head, f)
	log.Infof("save checkpoint finished with root: %s", head)
}

func write(root *TreeNode, f *os.File) {
	f.Write([]byte(root.String() + "\n"))
	for _, node := range root.children {
		write(node, f)
	}
}

func (root *TreeNode) LoadCheckpoint() {
	log.Infof("load checkpoint from file: %s", root.ckfile)
	if root.parent != nil {
		panic("only empty tree can load checkpoint")
	}
	root.mux.Lock()
	defer root.mux.Unlock()
	log.Trace("load checkpoint get lock")

	f, err := os.Open(root.ckfile)
	if os.IsNotExist(err) {
		log.Warnf("no checkpoint to load, file will be created: %s", root.ckfile)
		return
	}
	if err != nil {
		panic(err)
	}
	rd := bufio.NewReader(f)
	// first line is size
	line, err := readline(rd)
	if err != nil {
		panic(err)
	}

	size, err := strconv.Atoi(line)
	if err != nil {
		panic(fmt.Sprintf("invalid ck first line:\n\t\"%s\"", line))
	}
	if int(root.state.size) != size {
		panic(fmt.Sprintf(
			"different size checkpoint file is loading, need: %d, but %d",
			root.state.size,
			size,
		))
	}
	log.Tracef("found size %d checkpoint, start read lines", size)
	// build tree
	cknodes := map[string]*CkNode{}
	total := 0
	for line, err = readline(rd); err == nil; line, err = readline(rd) {
		total++
		ckn, err := newCkNode(line)
		if err != nil {
			panic(err)
		}
		cknodes[ckn.id] = ckn
		if total%500000 == 0 {
			log.Tracef("read %d lines", total)
		}
	}

	log.Tracef("read lines finished, total: %d, start build nodes", total)
	nodes := map[string]*TreeNode{}
	for id, ckn := range cknodes {
		if _, exist := cknodes[ckn.p]; !exist {
			// this is root
			if root.visitTimes > 0 {
				panic("multiple root exist")
			}
			root.visitTimes = 299
			root.result = ckn.r
			nodes[id] = root
		} else {
			node := ckn.newTreeNode(size)
			nodes[id] = node
		}
	}
	log.Trace("nodes build finished, start build tree")
	for id, n := range nodes {
		if pn, exist := nodes[cknodes[id].p]; !exist && n != root {
			panic(fmt.Sprintf("parent node should be exist, id: %s", id))
		} else {
			if pn != nil {
				n.parent = pn
				pn.children = append(pn.children, n)
			}
		}
	}
	log.Trace("tree build finished, start update states")
	// dfs build state
	root.updateState()
	log.Infof("load checkpoint finished with root: %s", root)
}

type CkNode struct {
	id string
	p  string
	a  string
	n  int
	t  int64
	r  [2]int
	u  int
}

func newCkNode(line string) (*CkNode, error) {
	arr := strings.Split(line, ",")
	if len(arr) != 7 {
		return nil, fmt.Errorf("invalid line format: %s", line)
	}
	var col [7]string
	for i, pre := range []string{"id:", "p:", "a:", "n:", "t:", "r:", "u:"} {
		if !strings.HasPrefix(arr[i], pre) {
			return nil, fmt.Errorf("invalid field 0: %s", arr[0])
		} else {
			col[i] = strings.TrimPrefix(arr[i], pre)
		}
	}
	var err error
	ckn := &CkNode{}
	ckn.id = col[0]
	ckn.p = col[1]
	ckn.a = col[2]
	ckn.n, err = strconv.Atoi(col[3])
	if err != nil {
		return nil, fmt.Errorf("invalid field 3: %s %v", arr[3], err)
	}
	ckn.t, err = strconv.ParseInt(col[4], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid field 4: %s %v", arr[4], err)
	}
	ckn.r, err = func() (rt [2]int, err error) {
		arr := strings.Split(col[5], "-")
		rt[0], err = strconv.Atoi(arr[0])
		if err != nil {
			return
		}
		rt[1], err = strconv.Atoi(arr[1])
		if err != nil {
			return
		}
		return
	}()
	if err != nil {
		return nil, fmt.Errorf("invalid field 5: %s %v", arr[5], err)
	}
	ckn.u, err = strconv.Atoi(col[6])
	if err != nil {
		return nil, fmt.Errorf("invalid field 6: %s %v", arr[6], err)
	}
	return ckn, nil
}

func (ckn *CkNode) newTreeNode(size int) *TreeNode {
	node := NewTree("", BoardSize(size))
	node.action = &Action{}
	_ = node.action.FromString(ckn.a)
	if ckn.u == 1 {
		node.allRollout = true
	}
	node.total = ckn.t
	node.visitTimes = ckn.n
	node.result = ckn.r
	return node
}

func readline(rd *bufio.Reader) (string, error) {
	var (
		isPrefix bool  = true
		err      error = nil
		line, l  []byte
	)
	for isPrefix && err == nil {
		l, isPrefix, err = rd.ReadLine()
		line = append(line, l...)
	}
	log.Debugf("read line: [%s]", string(line))
	return string(line), err
}

func (root *TreeNode) lock() {
	root.head().mux.Lock()
}
func (root *TreeNode) unlock() {
	root.head().mux.Unlock()
}
func (root *TreeNode) head() *TreeNode {
	head := root
	for head.parent != nil {
		head = head.parent
	}
	return head
}
