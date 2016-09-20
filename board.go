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
	NO move = iota
	UP
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
	curDirection  move
	r             *rand.Rand
}

func newBoard(w io.Writer, height, width int) *board {
	b := make([][]el, height)
	for i := range b {
		b[i] = make([]el, width)
	}
	board := &board{
		w:            w,
		b:            b,
		height:       height,
		width:        width,
		curDirection: NO,
		r:            rand.New(rand.NewSource(time.Now().UnixNano())),
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

// TODO only go sideways, not the opposite direction
func (b *board) changeDirection(s string) {
	switch s {
	case "w":
		b.curDirection = UP
	case "s":
		b.curDirection = DOWN
	case "a":
		b.curDirection = LEFT
	case "d":
		b.curDirection = RIGHT
	default:
	}
}

func (b *board) nextMove() {
	if b.curDirection == NO {
		return
	}
	b.playMove(b.curDirection)
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
	if c1 < 0 {
		c1 = c1 + width
	}
	if c2 < 0 {
		c2 = c2 + height
	}
	c1 = c1 % b.width
	c2 = c2 % b.height

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
