package main

import (
	"fmt"
	"os"
	"time"
)

const framerate = 1 * time.Second
const height, width = 3, 3

func main() {
	b := newBoard(os.Stdout, height, width)
	b.draw()

	time.Sleep(framerate)

	fmt.Fprintf(os.Stdout, "stop")
}
