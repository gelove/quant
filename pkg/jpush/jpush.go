package jpush

import (
	jpushclient "github.com/ylywyn/jpush-api-go-client"
)

const (
	appKey = "d87ab1334c0e95510451b601"
	secret = "5346211581a9395945024784"
)

var pf jpushclient.Platform
var ad jpushclient.Audience

func init() {
	pf.Add(jpushclient.ANDROID)
	pf.Add(jpushclient.IOS)
	pf.Add(jpushclient.WINPHONE)

	s := []string{"t1", "t2", "t3"}
	ad.SetTag(s)
	id := []string{"1", "2", "3"}
	ad.SetID(id)
	//ad.All()
}

// NewNotice 创建通知
func NewNotice(alter string) *jpushclient.Notice {
	notice := &jpushclient.Notice{}
	notice.SetAlert(alter)
	notice.SetAndroidNotice(&jpushclient.AndroidNotice{Alert: "AndroidNotice"})
	notice.SetIOSNotice(&jpushclient.IOSNotice{Alert: "IOSNotice"})
	notice.SetWinPhoneNotice(&jpushclient.WinPhoneNotice{Alert: "WinPhoneNotice"})
	return notice
}

// NewMessage 创建消息
func NewMessage(title, cotnet string) *jpushclient.Message {
	msg := &jpushclient.Message{}
	msg.Title = title
	msg.Content = cotnet
	return msg
}
