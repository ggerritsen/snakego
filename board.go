package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"os/exec"
	"time"
)

type move int

const (
	UP move = iota
	DOWN
	LEFT
	RIGHT
)

type el int

const (
	NONE el = iota
	SNAKE
	APPLE
	WALL
)

type board struct {
	w             io.Writer
	b             [][]el
	height, width int
	snake         []int
	r             *rand.Rand
}

func newBoard(w io.Writer, height, width int) *board {
	b := make([][]el, height)
	for i := range b {
		b[i] = make([]el, width)
	}
	board := &board{
		w:      w,
		b:      b,
		height: height,
		width:  width,
		r:      rand.New(rand.NewSource(time.Now().UnixNano())),
	}

	board.b[1][1] = SNAKE
	board.snake = []int{1, 1}
	board.addApple()

	return board
}

func (b *board) addApple() {
	// random
	for i := 0; i < 10; i++ {
		x, y := b.r.Intn(b.height), b.r.Intn(b.width)
		if b.b[x][y] == NONE {
			b.b[x][y] = APPLE
			return
		}
	}

	// non-random
	for i := 0; i < b.height; i++ {
		for j := 0; j < b.width; j++ {
			if b.b[i][j] == NONE {
				b.b[i][j] = APPLE
				return
			}
		}
	}
}

func (b *board) draw() {
	clear(b.w)
	for i := range b.b {
		s := ""
		for j := range b.b[i] {
			switch b.b[i][j] {
			case SNAKE:
				s += "S"
			case APPLE:
				s += "A"
			default:
				s += "."
			}
		}
		fmt.Fprintf(b.w, s)
		fmt.Fprintf(b.w, "\n")
	}
}

func (b *board) move(s string) {
	switch s {
	case "w":
		b.playMove(UP)
	case "s":
		b.playMove(DOWN)
	case "a":
		b.playMove(LEFT)
	case "d":
		b.playMove(RIGHT)
	default:
	}
}

func (b *board) playMove(m move) {
	x, y := b.snake[0], b.snake[1]

	moves := map[move]func() (int, int){
		UP:    func() (int, int) { return x - 1, y },
		DOWN:  func() (int, int) { return x + 1, y },
		LEFT:  func() (int, int) { return x, y - 1 },
		RIGHT: func() (int, int) { return x, y + 1 },
	}

	c1, c2 := moves[m]()
	b.b[c1][c2] = SNAKE
	b.b[x][y] = NONE
	b.snake = []int{c1, c2}
}

func clear(w io.Writer) {
	cmd := exec.Command("clear")
	cmd.Stdout = w
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}
