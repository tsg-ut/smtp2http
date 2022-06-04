package main

import (
	"net"
	"net/mail"
	"net/netip"

	"github.com/alash3al/go-smtpsrv/v3"
)

type Context interface {
	From() *mail.Address
	RemoteAddr() net.Addr
	Parse() (*smtpsrv.Email, error)
	Read(p []byte) (int, error)
	To() *mail.Address
}

var _ Context = (*smtpsrv.Context)(nil)

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
func (tc *TestContext) Parse() (*smtpsrv.Email, error) {
	panic("unimplmeneted")
}
func (tc *TestContext) Read(p []byte) (int, error) {
	panic("unimplmeneted")
}
func (tc *TestContext) To() *mail.Address {
	panic("unimplmeneted")
}

var _ Context = (*TestContext)(nil)
