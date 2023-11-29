package codec

import (
	"bufio"
	"encoding/gob"
	"gii/glog"
	"io"
)

type GobCodec struct {
	conn io.ReadWriteCloser
	enc  *gob.Encoder
	dec  *gob.Decoder
	buf  *bufio.Writer
}

var _ Codec = (*GobCodec)(nil)

func NewGobCodec(conn io.ReadWriteCloser) Codec {
	buf := bufio.NewWriter(conn)
	return &GobCodec{
		conn: conn,
		enc:  gob.NewEncoder(buf),
		dec:  gob.NewDecoder(conn),
		buf:  buf,
	}
}

func (g *GobCodec) Close() error {
	return g.conn.Close()
}

func (g *GobCodec) ReadHeader(h *Header) error {
	return g.dec.Decode(h)
}

func (g *GobCodec) ReadBody(body interface{}) error {
	return g.dec.Decode(body)
}

func (g *GobCodec) Write(h *Header, body interface{}) (err error) {
	defer func() {
		_ = g.buf.Flush()
		if err != nil {
			glog.Error(err)
		}
	}()
	if err = g.enc.Encode(h); err != nil {
		return
	}
	if err = g.enc.Encode(body); err != nil {
		return
	}
	return nil
}
