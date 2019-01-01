package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/websocket"
	"github.com/kr/pty"
	"golang.org/x/crypto/ssh/terminal"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var data chan []byte

func writeToSocket(ws *websocket.Conn) {
	ws.SetWriteDeadline(time.Time{})
	for p := range data {
		if err := ws.WriteMessage(websocket.TextMessage, p); err != nil {
			log.Println(err)
			return
		}
	}

}

type webWriter struct{}

func (wr *webWriter) Write(p []byte) (n int, err error) {
	data <- p
	return len(p), nil
}

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

	thenga1, _ := os.Create("/tmp/thenga1.txt")
	thenga2, _ := os.Create("/tmp/thenga2.txt")
	defer thenga1.Close()
	defer thenga2.Close()
	w := io.MultiWriter(ptmx, thenga1)
	w2 := io.MultiWriter(os.Stdout, thenga2, new(webWriter))
	// Copy stdin to the pty and the pty to stdout.
	go func() {
		_, _ = io.Copy(w, os.Stdin)
	}()
	_, _ = io.Copy(w2, ptmx)
	return nil
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		if _, ok := err.(websocket.HandshakeError); !ok {
			log.Println(err)
		}
		return
	}

	go writeToSocket(ws)
}

func main() {
	data = make(chan []byte, 1024)
	http.HandleFunc("/ws", serveWs)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, "./index.html") })
	go test()
	if err := http.ListenAndServe("0.0.0.0:8080", nil); err != nil {
		log.Fatal(err)
	}
}
