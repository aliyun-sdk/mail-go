package smtp

import (
	"flag"
	"testing"
)

var client *Client
var user, password string

func init() {
	flag.StringVar(&user, "user", "", "您的发信地址")
	flag.StringVar(&password, "password", "", "您的SMTP密码")
	flag.Parse()
	client = New("smtpdm.aliyun.com:80", user, password)
}

func TestClient_Send(t *testing.T) {
	err := client.Send(
		From("精彩不用等"),
		Subject("我的测试邮件"),
		SendTo("249008728@qq.com"),
		SendCC("zsj@99xs.com"),
		SendBCC("zsj19881218@qq.com"),
		Content(Plain, "我的测试邮件"),
		ReplyTo("249008728@qq.com"),
	)
	if err != nil {
		t.Error(err)
	}
}
