package main

import (
	"io"
	"net"
	"net/mail"
	"net/netip"

	"github.com/jhillyerd/enmime"
)

type Context interface {
	RemoteAddr() net.Addr
	From() *mail.Address
	To() *mail.Address
	Parse() (*enmime.Envelope, error)
}

type ctxt struct {
	from       *mail.Address
	to         *mail.Address
	remoteAddr net.Addr
	data       io.Reader
}

func (c *ctxt) RemoteAddr() net.Addr {
	return c.remoteAddr
}
func (c *ctxt) From() *mail.Address {
	return c.from
}
func (c *ctxt) To() *mail.Address {
	return c.to
}
func (c *ctxt) Parse() (*enmime.Envelope, error) {
	return enmime.NewParser().ReadEnvelope(c.data)
}

type TestContext struct {
	FromEmail        string
	RemoteAddrString string
}

func (tc *TestContext) From() *mail.Address {
	return &mail.Address{
		Address: tc.FromEmail,
	}
}
func (tc *TestContext) RemoteAddr() net.Addr {
	return net.TCPAddrFromAddrPort(netip.AddrPort(netip.MustParseAddrPort(tc.RemoteAddrString)))
}
func (tc *TestContext) Parse() (*enmime.Envelope, error) {
	panic("unimplmeneted")
}
func (tc *TestContext) To() *mail.Address {
	panic("unimplmeneted")
}

var _ Context = (*TestContext)(nil)
