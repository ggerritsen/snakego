package main

import (
	"fmt"
	"io"
	"os/exec"
)

type board struct {
	h, w int
}

func newBoard(height, width int) *board {
	return &board{
		h: height,
		w: width,
	}
}

func (b *board) draw(w io.Writer) {
	clear(w)
	for i := 0; i < b.h; i++ {
		s := ""
		for j := 0; j < b.w; j++ {
			s += "."
		}
		fmt.Fprintf(w, s)
		fmt.Fprintf(w, "\n")
	}
}

func clear(w io.Writer) {
	cmd := exec.Command("clear")
	cmd.Stdout = w
	cmd.Run()
}
