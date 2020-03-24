package smtp

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"strings"
)

const CRLF = "\r\n"

// message SMTP消息体
type message struct {
	from        email
	replyTo     email
	sendTo      []string
	sendCC      []string
	sendBCC     []string
	subject     string
	content     content
	attachments []attachment
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
	buf.WriteString("MIME-Version: 1.0" + CRLF)
	m.writeMixed(buf)
	return buf.Bytes()
}

func (m *message) writeMixed(buf *bytes.Buffer) {
	buf.WriteString("Content-Type: multipart/mixed; boundary=\"MixedBoundaryString\"" + CRLF + CRLF)
	buf.WriteString("--MixedBoundaryString" + CRLF)
	m.writeRelated(buf)
	for _, att := range m.attachments {
		writeAttachment(buf, att.filename, att.contentType, att.data)
	}
	buf.WriteString("--MixedBoundaryString--")
}

func (m *message) writeRelated(buf *bytes.Buffer) {
	buf.WriteString("Content-Type: multipart/related; boundary=\"RelatedBoundaryString\"" + CRLF + CRLF)
	buf.WriteString("--RelatedBoundaryString" + CRLF)
	m.writeAlternative(buf)
	buf.WriteString("--RelatedBoundaryString--" + CRLF + CRLF)
}

func (m *message) writeAlternative(buf *bytes.Buffer) {
	buf.WriteString("Content-Type: multipart/alternative; boundary=\"AlternativeBoundaryString\"" + CRLF + CRLF)
	buf.WriteString("--AlternativeBoundaryString" + CRLF)

	buf.WriteString(m.content.toString())

	buf.WriteString("--AlternativeBoundaryString--" + CRLF + CRLF)
}

func writeAttachment(buf *bytes.Buffer, filename string, contentType string, data []byte) {
	buf.WriteString("--MixedBoundaryString" + CRLF)
	buf.WriteString(fmt.Sprintf("Content-Type: %s;name=\"%s\"\r\n", contentType, filename))
	buf.WriteString("Content-Transfer-Encoding: base64\r\n")
	buf.WriteString(fmt.Sprintf("Content-Disposition: attachment;filename=\"%s\"\r\n\r\n", filename))

	encodedData := make([]byte, base64.StdEncoding.EncodedLen(len(data)))
	base64.StdEncoding.Encode(encodedData, data)
	buf.Write(encodedData)
	buf.WriteString(CRLF)
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

func Attachment(filename, contentType string, data []byte) Option {
	return func(m *message) {
		m.attachments = append(m.attachments, attachment{
			filename:    filename,
			data:        data,
			contentType: contentType,
		})
	}
}
