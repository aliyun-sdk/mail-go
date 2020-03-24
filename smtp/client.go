package smtp

import (
	"net"
	"net/smtp"
)

type Client struct {
	addr string
	user string
	auth smtp.Auth
}

func (c *Client) Send(opts ...Option) error {
	msg := newMessage(opts...)
	if msg.from.addr == "" {
		msg.from.addr = c.user
	}
	if msg.replyTo.name == "" {
		msg.replyTo.name = msg.from.name
	}
	if msg.replyTo.addr == "" {
		msg.replyTo.addr = msg.from.addr
	}
	return smtp.SendMail(c.addr, c.auth, msg.from.addr, msg.allAddrs(), msg.toBytes())
}

func New(addr, user, pwd string) *Client {
	host, _, _ := net.SplitHostPort(addr)
	return &Client{addr: addr, user: user, auth: smtp.PlainAuth("", user, pwd, host)}
}
