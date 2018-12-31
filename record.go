package main

import (
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/kr/pty"
	"golang.org/x/crypto/ssh/terminal"
)

func test() error {
	// Get the SHELL which is getting currently used
	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "sh"
	}

	// Create arbitrary command.
	c := exec.Command(shell)

	// Start the command with a pty.
	ptmx, err := pty.Start(c)
	if err != nil {
		return err
	}
	// Make sure to close the pty at the end.
	defer ptmx.Close() // Best effort.

	// Handle pty size.
	ch := make(chan os.Signal, 1)
	defer close(ch)
	signal.Notify(ch, syscall.SIGWINCH)
	go func() {
		for range ch {
			if err := pty.InheritSize(os.Stdin, ptmx); err != nil {
				log.Printf("error resizing pty: %s", err)
			}
		}
	}()
	ch <- syscall.SIGWINCH // Initial resize.
	// Set stdin in raw mode.
	oldState, err := terminal.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer terminal.Restore(int(os.Stdin.Fd()), oldState) // Best effort.

	thenga, _ := os.Create("/tmp/thenga.txt")
	defer thenga.Close()
	w := io.MultiWriter(ptmx, thenga)
	w2 := io.MultiWriter(os.Stdout, thenga)
	// Copy stdin to the pty and the pty to stdout.
	go func() {
		_, _ = io.Copy(w, os.Stdin)
	}()
	_, _ = io.Copy(w2, ptmx)

	return nil
}

func main() {
	if err := test(); err != nil {
		log.Fatal(err)
	}
}
