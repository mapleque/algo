package main

import (
	"time"

	"github.com/mapleque/algo"
)

func main() {
	root := algo.NewTree("model", algo.BoardSizeSmall)
	algo.SetLogLevel(algo.Info)
	root.LoadCheckpoint()
	go root.MCTS()
	saveCheckpoint(root)
}

func saveCheckpoint(root *algo.TreeNode) {
	<-time.After(10 * time.Second)
	root.SaveCheckpoint()
	saveCheckpoint(root)
}
