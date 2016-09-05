package main

import (
	"fmt"
	"io"
	"os/exec"
)

type el int

const (
	NONE el = iota
	SNAKE
	APPLE
	WALL
)

type board struct {
	w io.Writer
	b [][]el
}

func newBoard(w io.Writer, height, width int) *board {
	b := make([][]el, height)
	for i := range b {
		b[i] = make([]el, width)
	}
	b[1][1] = SNAKE
	b[2][2] = APPLE
	return &board{
		w: w,
		b: b,
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

func clear(w io.Writer) {
	cmd := exec.Command("clear")
	cmd.Stdout = w
	cmd.Run()
}
