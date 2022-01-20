package main

import (
	"fmt"
	"time"

	"github.com/mapleque/algo"
)

func main() {
	root := algo.NewTree("model", algo.BoardSizeSmall)
	algo.SetLogLevel(algo.Trace)
	root.LoadCheckpoint()

	node := root
	steps := 0
	for node != nil {
		fmt.Println(
			fmt.Sprintf(
				"steps %d, %s",
				steps,
				node.GetAction(),
			),
			node.GetState().GetBoard())
		node = node.BestMove()
		steps++
		<-time.After(2 * time.Second)
	}
}
