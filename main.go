package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
)

const frameRate = 1000 * time.Millisecond
const height, width = 3, 3

func main() {
	// setup stty to forward keys directly
	cmd := exec.Command("stty", "-g")
	cmd.Stdin = os.Stdin
	b, err := cmd.Output()
	if err != nil {
		log.Fatalf("Could not backup terminal settings: %s", err)
	}

	restoreStty := func() {
		cmd := exec.Command("/bin/stty", "-g", string(b))
		cmd.Stdin = os.Stdin
		if err := cmd.Run(); err != nil {
			log.Printf("Could not reset stty. Try it manually with '/bin/stty -g %s'", string(b))
			log.Fatal(err)
		}
	}

	cmd = exec.Command("/bin/stty", "cbreak", "-echo")
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		log.Fatalf("Could not setup terminal: %s", err)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, os.Kill, syscall.SIGTERM)

	board := newBoard(os.Stdout, height, width)

	// catch user input to change snake's direction
	go func() {
		sc := bufio.NewScanner(os.Stdin)
		sc.Split(bufio.ScanBytes)

		for {
			// TODO: make it possible to use arrow keys
			if sc.Scan() {
				board.changeDirection(sc.Text())
			} else {
				fmt.Printf("error: %s", sc.Err())
				break
			}
		}
	}()

	// update board every frameRate
	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Printf("[ERROR] Got panic: %s\n", err)
				restoreStty()
				os.Exit(1)
			}
		}()

		for {
			board.draw()
			time.Sleep(frameRate)
			if ok := board.playMove(); !ok {
				// game over
				fmt.Printf("Game over.\n")
				time.Sleep(frameRate)
				stop <- os.Interrupt
				break
			}
		}
	}()

	<-stop
	restoreStty()
	println("The end.")
	os.Exit(0)
}
