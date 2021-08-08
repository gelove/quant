package mail

import (
	"fmt"
	"log"
	"net/smtp"
	"quant/pkg/utils/logging"
	"time"

	"github.com/jordan-wright/email"
)

var ch = make(chan *email.Email)

const title = "allen"
const user = "geloves@163.com"
const token = "XLTSZQTCSRKPZXJJ"
const host = "smtp.163.com"
const port = 25 // ssl 465/994

func init() {
	p, err := email.NewPool(
		fmt.Sprintf("%s:%d", host, port),
		4,
		smtp.PlainAuth(title, user, token, host),
	)
	if err != nil {
		logging.ErrorF("email init failed => %+v", err)
		panic(err)
	}
	for i := 0; i < 4; i++ {
		go func() {
			for e := range ch {
				log.Printf("Send Email => %+v\n", e)
				p.Send(e, 10*time.Second)
			}
		}()
	}
}

func SendEmail(html []byte) {
	e := email.NewEmail()
	e.From = "Allen <geloves@163.com>"
	e.To = []string{"61114099@qq.com"}
	// e.Cc = []string{"test_cc@example.com"} // 抄送
	// e.Bcc = []string{"test_bcc@example.com"}
	e.Subject = "流动性监控"
	e.HTML = html
	ch <- e
}
