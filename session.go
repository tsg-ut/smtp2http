package main

import (
	"io"
	"net"
	"net/mail"

	"github.com/emersion/go-smtp"
)

type HandlerFunc func(Context) error

type backend struct {
	handler HandlerFunc
}

func (be *backend) Login(_ *smtp.ConnectionState, _, _ string) (smtp.Session, error) {
	return nil, smtp.ErrAuthUnsupported
}
func (be *backend) AnonymousLogin(con *smtp.ConnectionState) (smtp.Session, error) {
	return &session{
		handler: be.handler,
		addr:    con.RemoteAddr,
	}, nil
}

type session struct {
	from *mail.Address
	to   *mail.Address
	addr net.Addr

	handler HandlerFunc
}

func (s *session) Reset() {}
func (s *session) Logout() error {
	return nil
}
func (s *session) Mail(from string, opts smtp.MailOptions) (err error) {
	s.from, err = mail.ParseAddress(from)
	return
}
func (s *session) Rcpt(to string) (err error) {
	s.to, err = mail.ParseAddress(to)
	return
}
func (s *session) Data(r io.Reader) error {
	return s.handler(&ctxt{
		from:       s.from,
		to:         s.to,
		remoteAddr: s.addr,
		data:       r,
	})
}
