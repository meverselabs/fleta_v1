package p2p

import (
	"log"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/fletaio/fleta_v1/common/binutil"
	"github.com/fletaio/fleta_v1/core/types"
)

// TCPPeer manages send and recv of the connection
type TCPPeer struct {
	sync.Mutex
	conn          net.Conn
	id            string
	name          string
	isClose       bool
	connectedTime int64
	pingCount     uint64
	pingType      uint16
}

// NewTCPPeer returns a TCPPeer
func NewTCPPeer(conn net.Conn, ID string, Name string, connectedTime int64) *TCPPeer {
	if len(Name) == 0 {
		Name = ID
	}
	p := &TCPPeer{
		conn:          conn,
		id:            ID,
		name:          Name,
		connectedTime: connectedTime,
		pingType:      types.DefineHashedType("p2p.PingMessage"),
	}

	go func() {
		defer p.Close()

		pingCountLimit := uint64(3)
		for !p.isClose {
			p.Lock()
			if err := p.conn.SetWriteDeadline(time.Now().Add(5 * time.Second)); err != nil {
				p.Unlock()
				return
			}
			_, err := p.conn.Write(binutil.LittleEndian.Uint16ToBytes(p.pingType))
			if err != nil {
				p.Unlock()
				return
			}
			p.Unlock()
			if atomic.AddUint64(&p.pingCount, 1) > pingCountLimit {
				return
			}
			time.Sleep(3 * time.Second)
		}
	}()
	return p
}

// ID returns the id of the peer
func (p *TCPPeer) ID() string {
	return p.id
}

// Name returns the name of the peer
func (p *TCPPeer) Name() string {
	return p.name
}

// Close closes TCPPeer
func (p *TCPPeer) Close() {
	p.isClose = true
	p.conn.Close()
}

// IsClosed returns it is closed or not
func (p *TCPPeer) IsClosed() bool {
	return p.isClose
}

// ReadPacket returns a packet data
func (p *TCPPeer) ReadPacket() ([]byte, error) {
	for {
		if t, _, err := ReadUint16(p.conn); err != nil {
			return nil, err
		} else {
			atomic.StoreUint64(&p.pingCount, 0)
			if t == p.pingType {
				continue
			} else {
				if Len, _, err := ReadUint32(p.conn); err != nil {
					return nil, err
				} else {
					bs := make([]byte, 6+Len)
					binutil.LittleEndian.PutUint16(bs, t)
					binutil.LittleEndian.PutUint32(bs[2:], Len)
					if _, err := FillBytes(p.conn, bs[6:]); err != nil {
						return nil, err
					}
					return bs, nil
				}
			}
		}
	}
}

// SendPacket sends packet to the WebsocketPeer
func (p *TCPPeer) SendPacket(bs []byte) {
	if p.isClose {
		return
	}

	p.Lock()
	defer p.Unlock()

	if err := p.conn.SetWriteDeadline(time.Now().Add(5 * time.Second)); err != nil {
		log.Println(p.name, "SendPacket.SetWriteDeadline", err)
		p.Close()
		return
	}
	if _, err := p.conn.Write(bs); err != nil {
		log.Println(p.name, "SendPacket.Write", err)
		p.Close()
		return
	}
}

// ConnectedTime returns peer connected time
func (p *TCPPeer) ConnectedTime() int64 {
	return p.connectedTime
}
