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
		snake:       []coord{{1, 1}}, // TODO start in the middle
		currentMove: NO,
		r:           rand.New(rand.NewSource(time.Now().UnixNano())),
	}

	board.addApple()
	return board
}

func (b *board) addApple() {
	// random
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

	// non-random
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

func (b *board) draw(w io.Writer) {
	clear(w)
	snakeHead := b.snake[0]

	bb := make([][]string, b.height)
	for i := range bb {
		bb[i] = make([]string, b.width)
	}
	for i := range bb {
		for j := range bb[i] {
			bb[i][j] = "."
		}
	}

	bb[snakeHead.x][snakeHead.y] = "$"
	for i := 1; i < len(b.snake); i++ {
		e := b.snake[i]
		bb[e.x][e.y] = "S"
	}
	bb[b.apple.x][b.apple.y] = "A"

	for i := range bb {
		for j := range bb[i] {
			fmt.Fprintf(w, bb[i][j])
		}
		fmt.Fprintf(w, "\n")
	}
}

func (b *board) changeDirection(s string) {
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

// playMove performs the next move of the snake.
// It returns false in case of game over.
func (b *board) playMove() bool {
	if b.currentMove == NO {
		return true
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

func (b *board) update(c coord) bool {
	for _, s := range b.snake {
		if s == c {
			// snake hit itself, game over
			return false
		}
	}

	oldSnake := b.snake
	newSnake := []coord{c}
	for i := 0; i < len(oldSnake)-1; i++ {
		newSnake = append(newSnake, oldSnake[i])
	}
	b.snake = newSnake

	if b.apple == c {
		newSnake = append(newSnake, oldSnake[len(oldSnake)-1])
		b.snake = newSnake
		b.addApple()
	}

	return true
}

func clear(w io.Writer) {
	cmd := exec.Command("clear")
	cmd.Stdout = w
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}
