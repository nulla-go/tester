package reciever

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"

	"github.com/strengine/Core/av/avutil"

	"github.com/strengine/Core/av/pubsub"
	"github.com/strengine/Core/format/rtmp"

	"github.com/strengine/Core/format"
)

type writeFlusher struct {
	httpflusher http.Flusher
	io.Writer
}

func (self writeFlusher) Flush() error {
	self.httpflusher.Flush()
	return nil
}

type Channel struct {
	que *pubsub.Queue
}

type Reciever struct {
	server   rtmp.Server
	l        sync.RWMutex
	channels map[string]*Channel
}

func (r *Reciever) Start(port string) {
	format.RegisterAll()
	r.server.HandlePublish = func(conn *rtmp.Conn) {
		streams, _ := conn.Streams()

		r.l.Lock()
		fmt.Println("Connected successefully")

		ch := r.channels[conn.URL.Path]
		if ch == nil {
			ch = &Channel{}
			ch.que = pubsub.NewQueue()
			ch.que.WriteHeader(streams)
			r.channels[conn.URL.Path] = ch
		} else {
			ch = nil
		}
		r.l.Unlock()
		if ch == nil {
			fmt.Println("Channel is nil")
			return
		}

		avutil.CopyPackets(ch.que, conn)

		r.l.Lock()
		delete(r.channels, conn.URL.Path)
		r.l.Unlock()
		ch.que.Close()
	}
	fmt.Println("Reciever created")
	fmt.Println("Listing port", port)
	r.server.Addr = port
	err := r.server.ListenAndServe()
	if err != nil {
		log.Println(err)
	}
}
