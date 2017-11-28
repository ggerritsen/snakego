package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
)

const height, width = 11, 11

var frameRateMs = 500

func main() {
	// setup stty to forward keys directly
	cmd := exec.Command("stty", "-g")
	cmd.Stdin = os.Stdin
	b, err := cmd.Output()
	if err != nil {
		log.Fatalf("Could not backup terminal settings: %s", err)
	}

	// restoryStty to be called when game ends
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

	// trap exit signal
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, os.Kill, syscall.SIGTERM)

	board := newBoard(height, width)

	// catch user input to change snake's direction
	go func() {
		sc := bufio.NewScanner(os.Stdin)
		sc.Split(bufio.ScanBytes)

		for {
			// TODO: make it possible to use arrow keys
			if sc.Scan() {
				board.ChangeDirection(sc.Text())
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
			board.Draw(os.Stdout)
			time.Sleep(time.Duration(frameRateMs) * time.Millisecond)
			ateApple, err := board.PlayMove()
			if err != nil {
				// game over
				fmt.Printf("Game over: %s.\n", err)
				time.Sleep(time.Duration(frameRateMs) * time.Millisecond)
				stop <- os.Interrupt
				break
			}
			if ateApple {
				x := math.Pow(0.99, float64(len(board.snake))) 
				frameRateMs = int(float64(frameRateMs) * x)
				fmt.Printf("framerate %dms", frameRateMs)
				time.Sleep(1 * time.Second)
			}
		}
	}()

	<-stop
	restoreStty()
	println("The end.")
	os.Exit(0)
}
