package smtp

import (
	"bytes"
	"strings"
)

const CRLF = "\r\n"

// message SMTP消息体
type message struct {
	from    email
	replyTo email
	sendTo  []string
	sendCC  []string
	sendBCC []string
	subject string
	content content
}

func (m *message) allAddrs() []string {
	addrs := append(m.sendTo, m.sendCC...)
	return append(addrs, m.sendBCC...)
}

func (m *message) toBytes() []byte {
	buf := bytes.NewBuffer(nil)
	buf.WriteString("From: " + m.from.toString() + CRLF)
	buf.WriteString("To: " + strings.Join(m.sendTo, ";") + CRLF)
	buf.WriteString("Cc: " + strings.Join(m.sendCC, ";") + CRLF)
	buf.WriteString("Bcc: " + strings.Join(m.sendBCC, ";") + CRLF)
	buf.WriteString("Reply-To: " + m.replyTo.toString() + CRLF)
	buf.WriteString("Subject: " + m.subject + CRLF)
	buf.WriteString(m.content.toString())
	return buf.Bytes()
}

func newMessage(opts ...Option) *message {
	msg := new(message)
	for _, opt := range opts {
		opt(msg)
	}
	return msg
}

type Option func(m *message)

func From(name string, addr ...string) Option {
	return func(m *message) {
		if len(addr) == 0 {
			m.from = email{name: name}
		} else {
			m.from = email{name: name, addr: addr[0]}
		}
	}
}

func Subject(s string) Option {
	return func(m *message) {
		m.subject = s
	}
}

func SendTo(to ...string) Option {
	return func(m *message) {
		m.sendTo = to
	}
}

func SendCC(cc ...string) Option {
	return func(m *message) {
		m.sendCC = cc
	}
}

func SendBCC(bcc ...string) Option {
	return func(m *message) {
		m.sendBCC = bcc
	}
}

func ReplyTo(addr string, name ...string) Option {
	return func(m *message) {
		if len(name) == 0 {
			m.replyTo = email{addr: addr}
		} else {
			m.replyTo = email{addr: addr, name: name[0]}
		}
	}
}

func Content(typ ContentType, body string) Option {
	return func(m *message) {
		m.content = content{typ: typ, body: body}
	}
}
