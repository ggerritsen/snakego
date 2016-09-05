package main

import (
	"fmt"
	"os"
	"time"
)

const framerate = 1 * time.Second
const height, width = 3, 3

func main() {
	b := newBoard(height, width)
	b.draw(os.Stdout)

	time.Sleep(framerate)

	fmt.Fprintf(os.Stdout, "stop")
}
