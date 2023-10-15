package outbound

import (
	"context"
	"io"
	"net"
	"time"

	"github.com/Dreamacro/clash/component/dialer"
	C "github.com/Dreamacro/clash/constant"
)

type Reject struct {
	*Base
	drop bool
}

// DialContext implements C.ProxyAdapter
func (r *Reject) DialContext(ctx context.Context, metadata *C.Metadata, opts ...dialer.Option) (C.Conn, error) {
	return NewConn(&nopConn{r.drop}, r), nil
}

// ListenPacketContext implements C.ProxyAdapter
func (r *Reject) ListenPacketContext(ctx context.Context, metadata *C.Metadata, opts ...dialer.Option) (C.PacketConn, error) {
	return newPacketConn(&nopPacketConn{r.drop}, r), nil
}

func NewReject(drop bool) *Reject {
	var name string
	if drop {
		name = "REJECT-DROP"
	} else {
		name = "REJECT"
	}
	return &Reject{
		Base: &Base{
			name: name,
			tp:   C.Reject,
			udp:  true,
		},
		drop: drop,
	}
}

type nopConn struct {
	drop bool
}

func (rw *nopConn) Read(b []byte) (int, error) {
	if rw.drop {
		time.Sleep(C.DefaultRejectDropTimeout)
	}
	return 0, io.EOF
}

func (rw *nopConn) Write(b []byte) (int, error) {
	if rw.drop {
		time.Sleep(C.DefaultRejectDropTimeout)
	}
	return 0, io.EOF
}

func (rw *nopConn) Close() error                     { return nil }
func (rw *nopConn) LocalAddr() net.Addr              { return nil }
func (rw *nopConn) RemoteAddr() net.Addr             { return nil }
func (rw *nopConn) SetDeadline(time.Time) error      { return nil }
func (rw *nopConn) SetReadDeadline(time.Time) error  { return nil }
func (rw *nopConn) SetWriteDeadline(time.Time) error { return nil }

type nopPacketConn struct {
	drop bool
}

func (npc *nopPacketConn) WriteTo(b []byte, addr net.Addr) (n int, err error) {
	if npc.drop {
		time.Sleep(C.DefaultRejectDropTimeout)
	}
	return len(b), nil
}

func (npc *nopPacketConn) ReadFrom(b []byte) (int, net.Addr, error) {
	if npc.drop {
		time.Sleep(C.DefaultRejectDropTimeout)
	}
	return 0, nil, io.EOF
}

func (npc *nopPacketConn) Close() error                     { return nil }
func (npc *nopPacketConn) LocalAddr() net.Addr              { return &net.UDPAddr{IP: net.IPv4zero, Port: 0} }
func (npc *nopPacketConn) SetDeadline(time.Time) error      { return nil }
func (npc *nopPacketConn) SetReadDeadline(time.Time) error  { return nil }
func (npc *nopPacketConn) SetWriteDeadline(time.Time) error { return nil }
