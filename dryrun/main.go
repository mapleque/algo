package main

import (
	"fmt"
	"strconv"

	"github.com/mapleque/algo"
)

func main() {
	size := selectSize()
	root := algo.NewTree("model", size)
	algo.SetLogLevel(algo.Info)
	node := root
	steps := 0
	fmt.Println(node.GetState().GetBoard())
	for node != nil {
		node = userMove(node)
		steps++
		fmt.Println(
			fmt.Sprintf(
				"steps %d, %s",
				steps,
				node.GetAction(),
			),
			node.GetState().GetBoard())
	}
}

func selectSize() algo.BoardSize {
	var op string
	fmt.Println("Please enter board size:\n\t5(for dev)\t9\t13\t19")
	fmt.Scanln(&op)
	switch op {
	case "5":
		return algo.BoardSizeMini
	case "9":
		return algo.BoardSizeSmall
	case "13":
		return algo.BoardSizeMedium
	case "19":
		return algo.BoardSizeLarge
	default:
		fmt.Printf("invalid size value: %s, please enter 5 or 9 or 13 or 19\n", op)
		return selectSize()
	}
}

func userMove(root *algo.TreeNode) *algo.TreeNode {
	fmt.Println("Please enter the point you will move, like a1")
	var op string
	fmt.Scanln(&op)
	if len(op) < 2 {
		fmt.Printf("invalid move position: %s\n", op)
		return userMove(root)
	}
	y := int(op[0] - 'a')
	if y < 0 || y > root.GetState().GetBoard().Size() {
		fmt.Printf("invalid move position: %s\n", op)
		return userMove(root)
	}
	x, err := strconv.Atoi(op[1:])
	x = x - 1
	if err != nil || x < 0 || x > root.GetState().GetBoard().Size() {
		fmt.Printf("invalid move position: %s\n", op)
		return userMove(root)
	}
	node := root.FindChild(x, y)
	if node == nil {
		fmt.Printf("invalid move position: %s\n", op)
		return userMove(root)
	}
	return node
}
