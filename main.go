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

const frameRate = 3 * time.Second
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
		log.Printf("Restoring stty to %q", string(b))
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
			board.nextMove()
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, os.Kill, syscall.SIGTERM)

	<-stop
	restoreStty()
	println("Done.")
	os.Exit(0)
}
