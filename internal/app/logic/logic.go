package logic

import (
	"bytes"
	"html/template"
	"quant/pkg/mail"
)

var Temp *template.Template

func SendEmail() {
	/**
	ret["time"] + " " + pairName + "("+ret["contractA"]+")" + " 流动性：+" + \
	                         ret["usdAmount"] + "U  当前价:" + ret["price"] + \
	                         " 总流动性"+str(ret["totalLiquidity"])
	*/
	buf := bytes.NewBufferString("")
	data := map[string]interface{}{
		"time":           0,
		"contractA":      0,
		"usdAmount":      0,
		"price":          0,
		"totalLiquidity": 0,
		"holders":        0,
	}
	err := Temp.ExecuteTemplate(buf, "email.html", data)
	if err != nil {
		panic(err)
	}
	content := buf.Bytes()
	// logger.Infof("email content => %s", content)
	mail.SendEmail(content)
}
