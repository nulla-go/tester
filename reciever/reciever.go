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

type channel struct {
	que *pubsub.Queue
}

type reciever struct {
	server   *rtmp.Server
	l        sync.RWMutex
	channels map[string]*channel
}

func NewReciever() (reciverServer *reciever) {
	reciverServer = &reciever{}
	reciverServer.server = &rtmp.Server{}
	reciverServer.channels = map[string]*channel{}
	reciverServer.l = sync.RWMutex{}
	return
}

func (r *reciever) Start(port string) {
	format.RegisterAll()
	r.server.HandlePublish = func(conn *rtmp.Conn) {
		streams, _ := conn.Streams()
		key := conn.URL.Query().Get("key")
		if key == "" {
			return
		}
		r.l.Lock()
		fmt.Println("Connected successefully '", key, "'")

		ch := r.channels[key]
		if ch == nil {
			ch = &channel{}
			ch.que = pubsub.NewQueue()
			ch.que.WriteHeader(streams)
			r.channels[key] = ch
		} else {
			ch = nil
		}
		r.l.Unlock()
		if ch == nil {
			fmt.Println("Channel is nil")
			return
		}

		avutil.CopyPackets(ch.que, conn)
		fmt.Println("Connection '", key, "' is done")
		r.l.Lock()
		delete(r.channels, key)
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
