package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

const framerate = 1 * time.Second

func main() {
	println("...")
	println("...")
	println("...")

	time.Sleep(framerate)

	clear()

	fmt.Fprintf(os.Stdout, "stop")
}

func clear() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}
