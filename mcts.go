package algo

import (
	"math"
	"math/rand"
)

// MCTS expend tree.
func (root *TreeNode) MCTS() {
	for !root.allRollout {
		root.rollout()
	}
}

func (root *TreeNode) Stop() {
	root.ctx.stop = true
}

// BestMove get best move
func (root *TreeNode) BestMove() *TreeNode {
	return root.bestMove(1.4)
}

func (root *TreeNode) rollout() Player {
	log.Trace("rollout:", root)
	if root.ctx.stop {
		log.Trace("rollout stop")
		return -1
	}
	root.lock()
	log.Trace("rollout get lock")
	node := root.rolloutPolicy()
	if root.state.hasResult() || node == nil {
		result := root.state.Result()
		root.backpropagate(result)
		root.allRollout = true
		log.Infof("rollout a result %v %s:", result, root)
		root.unlock()
		log.Trace("rollout unlock")
		return result
	}
	root.unlock()
	log.Trace("rollout unlock")
	log.Tracef("rollout a node: %s with board %s", node, node.state.board)
	return node.rollout()
}

func (root *TreeNode) rolloutPolicy() *TreeNode {
	// random select an move from all possible moves.
	root.expand()
	list := []*TreeNode{}
	for _, node := range root.children {
		if !node.allRollout {
			list = append(list, node)
		}
	}
	n := len(list)
	if n == 0 {
		return nil
	}
	i := rand.Intn(n)
	return list[i]
}

func (root *TreeNode) backpropagate(result Player) {
	root.visitTimes++
	root.result[result]++
	if root.parent != nil {
		root.parent.backpropagate(result)
	}
}

func (root *TreeNode) bestMove(c float64) *TreeNode {
	if len(root.children) == 0 {
		return root.rolloutPolicy()
	}
	var index int
	var max float64

	for i, node := range root.children {
		node.uct = float64(node.q()/node.n()) +
			c*math.Sqrt(
				(2*math.Log(float64(root.n()))/float64(node.n())),
			)
		if node.uct > max {
			max = node.uct
			index = i
		}
	}
	return root.children[index]
}

func (root *TreeNode) q() int {
	wins := root.result[root.state.nextMovePlayer]
	loss := root.result[root.state.nextMovePlayer.next()]
	return wins - loss
}

func (root *TreeNode) n() int {
	return root.visitTimes + 1
}
