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

type coord struct {
	x, y int
}

type board struct {
	height, width int
	apple         coord
	snake         []coord
	wall          []coord // TODO use
	currentMove   move
	r             *rand.Rand
}

func newBoard(height, width int) *board {
	board := &board{
		height:      height,
		width:       width,
		snake:       []coord{{height / 2, width / 2}},
		currentMove: NO,
		r:           rand.New(rand.NewSource(time.Now().UnixNano())),
	}

	board.addApple()
	return board
}

func (b *board) addApple() {
	for i := 0; i < 10; i++ {
		c := coord{
			b.r.Intn(b.height),
			b.r.Intn(b.width),
		}
		for _, s := range b.snake {
			if s != c {
				b.apple = c
				return
			}
		}
	}

	// fallback option: place the apple non-random
	for i := 0; i < b.height; i++ {
		for j := 0; j < b.width; j++ {
			c := coord{i, j}
			for _, s := range b.snake {
				if s != c {
					b.apple = c
					return
				}
			}
		}
	}
}

func (b *board) Draw(w io.Writer) {
	clear(w)

	bb := make([][]string, b.height)
	for i := range bb {
		bb[i] = make([]string, b.width)
	}
	for i := range bb {
		for j := range bb[i] {
			bb[i][j] = "."
		}
	}

	snakeHead := b.snake[0]
	bb[snakeHead.x][snakeHead.y] = "$"
	for i := 1; i < len(b.snake); i++ {
		e := b.snake[i]
		bb[e.x][e.y] = "S"
	}
	bb[b.apple.x][b.apple.y] = "A"

	// print vertical buffer
	for i := 0; i < 17; i++ {
		fmt.Fprintf(w, "\n")
	}
	var horbuf string
	for i := 0; i < 90; i++ {
		horbuf = horbuf + " "
	}

	// print board
	for i := range bb {
		fmt.Fprintf(w, horbuf)
		for j := range bb[i] {
			fmt.Fprintf(w, bb[i][j])
		}
		fmt.Fprintf(w, "\n")
	}
}

func (b *board) ChangeDirection(s string) {
	switch s {
	case "w":
		if b.currentMove != DOWN {
			b.currentMove = UP
		}
	case "s":
		if b.currentMove != UP {
			b.currentMove = DOWN
		}
	case "a":
		if b.currentMove != RIGHT {
			b.currentMove = LEFT
		}
	case "d":
		if b.currentMove != LEFT {
			b.currentMove = RIGHT
		}
	default:
	}
}

// PlayMove performs the next move of the snake.
// It returns error in case of game over.
// It returns true if the snake ate the apple, false otherwise.
func (b *board) PlayMove() (bool, error) {
	if b.currentMove == NO {
		return false, nil
	}

	x, y := b.snake[0].x, b.snake[0].y
	moves := map[move]func() coord{
		UP:    func() coord { return coord{x - 1, y} },
		DOWN:  func() coord { return coord{x + 1, y} },
		LEFT:  func() coord { return coord{x, y - 1} },
		RIGHT: func() coord { return coord{x, y + 1} },
	}

	c := moves[b.currentMove]()
	if c.x < 0 {
		c.x = c.x + width
	}
	if c.y < 0 {
		c.y = c.y + height
	}
	c.x = c.x % b.width
	c.y = c.y % b.height

	return b.update(c)
}

func (b *board) update(c coord) (a bool, err error) {
	for _, s := range b.snake {
		if s == c {
			// snake hit itself, game over
			return false, fmt.Errorf("snake hit itself")
		}
	}

	oldSnake := b.snake
	newSnake := []coord{c}
	for i := 0; i < len(oldSnake)-1; i++ {
		newSnake = append(newSnake, oldSnake[i])
	}

	if b.apple == c {
		newSnake = append(newSnake, oldSnake[len(oldSnake)-1])
		b.addApple()
		a = true
	}

	b.snake = newSnake
	return a, nil
}

func clear(w io.Writer) {
	cmd := exec.Command("clear")
	cmd.Stdout = w
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}
