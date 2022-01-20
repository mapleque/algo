package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/mapleque/algo"
)

func main() {
	fmt.Println("Welcome to our algo game!")
	fmt.Println("First of all, we have to config some options.")
	conf := config()
	root := algo.NewTree("model", conf.Size)
	algo.SetLogLevel(algo.Info)
	fmt.Println("\nLoading AI database...")
	root.LoadCheckpoint()
	fmt.Println("Now, let's start, good luck!")
	fmt.Println()

	node := root
	steps := 0
	fmt.Sprintln(node.GetState().GetBoard())
	for node != nil {
		switch node.NextPlayer() {
		case conf.UserPlayer:
			fmt.Println("you turn:")
			node = userMove(node)
		default:
			fmt.Println("AI turn:")
			go node.MCTS()
			<-time.After(conf.EachStepDuration)
			node.Stop()
			node = node.BestMove()
		}
		if node != nil {
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
	fmt.Println("Game over, be happy!")
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

type Config struct {
	Size             algo.BoardSize
	UserPlayer       algo.Player
	EachStepDuration time.Duration
}

func config() *Config {
	conf := &Config{}
	conf.Size = selectSize()
	conf.UserPlayer = selectPlayer()
	conf.EachStepDuration = selectDuration()

	fmt.Printf(
		"Now, we will start the game with following configure:\n"+
			"\tboard size: %d*%d\n"+
			"\tyou are using: %s\n"+
			"\tAI each step will take: %s\n\n",
		conf.Size, conf.Size,
		conf.UserPlayer,
		conf.EachStepDuration,
	)
	fmt.Printf("Press <Enter> for start, or input 'no' for reconfigure:")
	var op string
	fmt.Scanln(&op)
	if op == "no" {
		return config()
	}
	return conf
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

func selectPlayer() algo.Player {
	var op string
	fmt.Println("Please enter your role:\n\t1/b:Black\t2/w:White")
	fmt.Scanln(&op)
	switch op {
	case "1", "b":
		return algo.PlayerBlack
	case "2", "w":
		return algo.PlayerWhite
	default:
		fmt.Printf("invalid role value: %s, please enter b or w or 1 or 2\n", op)
		return selectPlayer()
	}
}

func selectDuration() time.Duration {
	var op string
	fmt.Println("Please enter the time of each step spend for ai: (Second, 2-300)")
	fmt.Scanln(&op)
	sec, err := strconv.Atoi(op)
	if err != nil || sec < 2 || sec > 300 {
		fmt.Printf("invalid second value: %s, please enter: 2-300\n", op)
		return selectDuration()
	}
	return time.Duration(sec) * time.Second
}
