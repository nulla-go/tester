package requester

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/strengine/Core/av"
	"github.com/strengine/Core/av/avutil"
	"github.com/strengine/Core/av/pktque"
	"github.com/strengine/Core/format"
	"github.com/strengine/Core/format/rtmp"
)

func init() {
	format.RegisterAll()
}

type Requester struct {
	wg sync.WaitGroup
}

func (r *Requester) Push(addr []string) {
	file, err := avutil.Open("./samples/1.flv")
	if err != nil {
		log.Println("Failed open file:", err)
		return
	}

	for _, v := range addr {
		r.wg.Add(1)
		go r.transfer(v, &file)
	}
	r.wg.Wait()
	file.Close()
}

func (r *Requester) transfer(addr string, file *av.DemuxCloser) {
	// messages := make(chan struct{})
	conn := r.connect(addr)

	demuxer := &pktque.FilterDemuxer{Demuxer: *file, Filter: &pktque.Walltime{}}
	r.copy(conn, demuxer)
	conn.Close()
	log.Println("Connection closed: ", addr)
}

func (r *Requester) connect(addr string) *rtmp.Conn {
	conn, err := rtmp.Dial(addr)
	if err != nil {
		fmt.Println("Connection failed to ", addr)
		time.Sleep(5 * time.Second)
		return r.connect(addr)
	}
	log.Println("Connection success to ", addr)
	return conn
}

func (r *Requester) copy(connect *rtmp.Conn, demuxer *pktque.FilterDemuxer) {
	avutil.CopyFile(connect, demuxer)
}
