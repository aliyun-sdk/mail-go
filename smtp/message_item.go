package smtp

import (
	"fmt"
)

type ContentType string

const (
	Html  ContentType = "html"
	Plain ContentType = "plain"
)

// content 邮件内容
type content struct {
	typ  ContentType
	body string
}

func (c *content) toString() string {
	switch c.typ {
	case Html:
		return "Content-Type: text/html; charset=UTF-8" + CRLF + CRLF + c.body
	default:
		return "Content-Type: text/plain; charset=UTF-8" + CRLF + CRLF + c.body
	}
}

// email 电子邮箱
type email struct {
	name string
	addr string
}

func (e *email) toString() string {
	if e.name == "" {
		return e.addr
	}
	return fmt.Sprintf("%s <%s>", e.name, e.addr)
}
