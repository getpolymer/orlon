package internal

import (
	"io"
	"sync"
)

var subscribers []io.Writer

func Subscribe(w io.Writer) error {
	lock.Lock()
	defer lock.Unlock()
	subscribers = append(subscribers, w)
	return nil
}

var lock sync.RWMutex

// Publish publishes the message to all the listening writers
// not using io.Multiwriter because it stops if one if them fails and it is sequential
// preferably pool the goroutines instead of creating one per writer
func Publish(data []byte) error {
	lock.RLock()
	defer lock.RUnlock()
	for _, subscriber := range subscribers {
		go subscriber.Write(data)
	}
	return nil
}
