package requester

import (
	"bufio"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/strengine/Core/av"
	"github.com/strengine/Core/av/avutil"
	"github.com/strengine/Core/av/pktque"
	"github.com/strengine/Core/format/rtmp"
)

type Requester struct {
	wg sync.WaitGroup
}

func (r *Requester) Push(addr []string) {
	file, _ := avutil.Open("./samples/1.flv")

	for _, v := range addr {
		r.wg.Add(1)
		go r.transfer(v, &file)
	}
	r.wg.Wait()
	file.Close()
}

func (r *Requester) transfer(addr string, file *av.DemuxCloser) {

	conn := r.connect(addr)

	demuxer := &pktque.FilterDemuxer{Demuxer: *file, Filter: &pktque.Walltime{}}
	r.copy(conn, demuxer)

	bufio.NewReader(os.Stdin)
	//file.Close()
	conn.Close()
}

func (r *Requester) connect(addr string) *rtmp.Conn {
	conn, err := rtmp.Dial(addr)
	if err != nil {
		fmt.Println("Connection failed to ", addr)
		time.Sleep(2 * time.Second)
		return r.connect(addr)
	}
	return conn
}

func (r *Requester) copy(connect *rtmp.Conn, demuxer *pktque.FilterDemuxer) {
	avutil.CopyFile(connect, demuxer)
}
