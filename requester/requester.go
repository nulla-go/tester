package requester

import (
	"bufio"
	"os"
	"time"

	"github.com/strengine/core/av/avutil"
	"github.com/strengine/core/av/pktque"
	"github.com/strengine/core/format/rtmp"
)

type Requester struct {
}

func (r *Requester) Push(addr string) {
	file, _ := avutil.Open("./samples/1.flv")
	conn := r.connect(addr)

	demuxer := &pktque.FilterDemuxer{Demuxer: file, Filter: &pktque.Walltime{}}
	go r.copy(conn, demuxer)

	bufio.NewReader(os.Stdin)

	file.Close()
	conn.Close()
}

func (r *Requester) connect(addr string) *rtmp.Conn {
	conn, err := rtmp.Dial(addr)
	if err != nil {
		time.Sleep(2)
		return r.connect(addr)
	}
	return conn
}

func (r *Requester) copy(connect *rtmp.Conn, demuxer *pktque.FilterDemuxer) {
	avutil.CopyFile(connect, demuxer)
}
