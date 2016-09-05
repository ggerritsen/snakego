package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
)

const framerate = 1 * time.Second
const height, width = 3, 3

func main() {
	// setup stty to forward keys directly
	cmd := exec.Command("stty", "-g")
	cmd.Stdin = os.Stdin
	b, err := cmd.Output()
	if err != nil {
		log.Fatalf("Could not backup terminal settings: %s", err)
	}
	restoreStty := exec.Command("/bin/stty", "-g", string(b))
	restoreStty.Stdin = os.Stdin
	defer func() {
		if err := restoreStty.Run(); err != nil {
			log.Printf("Could not reset stty. Try it manually with '/bin/stty -g %s'", string(b))
			log.Fatal(err)
		}
	}()

	cmd = exec.Command("/bin/stty", "cbreak", "-echo")
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		log.Fatalf("Could not setup terminal: %s", err)
	}

	board := newBoard(os.Stdout, height, width)
	board.draw()

	println("What's your move?")

	sc := bufio.NewScanner(os.Stdin)
	sc.Split(bufio.ScanRunes)
	if sc.Scan() {
		input := sc.Bytes()
		fmt.Printf("you inputted %q\n", input)
	}

	println("Done.")
}
