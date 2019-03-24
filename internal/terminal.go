package internal

import (
	"io"
	"log"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/gorilla/websocket"
	"github.com/kr/pty"
	"golang.org/x/crypto/ssh/terminal"
)

type webWriter struct{}

var (
	connUrl    = &url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/echo"}
	conn, _, _ = websocket.DefaultDialer.Dial(connUrl.String(), nil)
)

func PublishToServer(data []byte) error {
	lock.RLock()
	defer lock.RUnlock()

	err := conn.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		return nil
	}

	return nil
}

func (ww webWriter) Write(data []byte) (int, error) {
	PublishToServer(data)
	return len(data), nil
}

func RunPseudoTerminal() error {
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

	thenga2, _ := os.Create("/tmp/thenga2.txt")
	defer thenga2.Close()
	w := io.MultiWriter(ptmx)
	w2 := io.MultiWriter(os.Stdout, thenga2, new(webWriter))
	// Copy stdin to the pty and the pty to stdout.
	go func() {
		_, _ = io.Copy(w, os.Stdin)
	}()
	_, _ = io.Copy(w2, ptmx)
	return nil
}
