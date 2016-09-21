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

type coord struct {
	x, y int
}

type board struct {
	b             [][]el
	height, width int
	snake         []coord
	curDirection  move
	r             *rand.Rand
}

func newBoard(height, width int) *board {
	b := make([][]el, height)
	for i := range b {
		b[i] = make([]el, width)
	}
	board := &board{
		b:            b,
		height:       height,
		width:        width,
		curDirection: NO,
		r:            rand.New(rand.NewSource(time.Now().UnixNano())),
	}

	board.snake = []coord{{1, 1}}
	board.addApple()

	return board
}

func (b *board) addApple() {
	b.refresh()

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

func (b *board) draw(w io.Writer) {
	clear(w)
	b.refresh()

	snakeHead := b.snake[0]

	for i := range b.b {
		s := ""
		for j := range b.b[i] {
			switch b.b[i][j] {
			case SNAKE:
				if i == snakeHead.x && j == snakeHead.y {
					s += "$"
				} else {
					s += "S"
				}
			case APPLE:
				s += "A"
			default:
				s += "."
			}
		}
		fmt.Fprintf(w, s)
		fmt.Fprintf(w, "\n")
	}
}

func (b *board) refresh() {
	for _, c := range b.snake {
		b.b[c.x][c.y] = SNAKE
	}
}

func (b *board) changeDirection(s string) {
	switch s {
	case "w":
		if b.curDirection != DOWN {
			b.curDirection = UP
		}
	case "s":
		if b.curDirection != UP {
			b.curDirection = DOWN
		}
	case "a":
		if b.curDirection != RIGHT {
			b.curDirection = LEFT
		}
	case "d":
		if b.curDirection != LEFT {
			b.curDirection = RIGHT
		}
	default:
	}
}

func (b *board) playMove() bool {
	if b.curDirection == NO {
		return true
	}
	b.refresh()

	x, y := b.snake[0].x, b.snake[0].y
	moves := map[move]func() (int, int){
		UP:    func() (int, int) { return x - 1, y },
		DOWN:  func() (int, int) { return x + 1, y },
		LEFT:  func() (int, int) { return x, y - 1 },
		RIGHT: func() (int, int) { return x, y + 1 },
	}

	c1, c2 := moves[b.curDirection]()
	if c1 < 0 {
		c1 = c1 + width
	}
	if c2 < 0 {
		c2 = c2 + height
	}
	c1 = c1 % b.width
	c2 = c2 % b.height

	return b.update(c1, c2)
}

func (b *board) update(x, y int) bool {
	for _, c := range b.snake {
		if c.x == x && c.y == y {
			return false
		}
	}

	newSnake := []coord{{x, y}}

	if b.b[x][y] == NONE {
		for i := 0; i < len(b.snake)-1; i++ {
			newSnake = append(newSnake, b.snake[i])
		}
	}

	if b.b[x][y] == APPLE {
		for i := 0; i < len(b.snake); i++ {
			newSnake = append(newSnake, b.snake[i])
		}
		b.addApple()
	}

	for _, c := range b.snake {
		b.b[c.x][c.y] = NONE
	}

	b.snake = newSnake
	for _, c := range b.snake {
		b.b[c.x][c.y] = SNAKE
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
