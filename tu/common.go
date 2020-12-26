package tu

import (
	"fmt"
	"github.com/Monibuca/plugin-gb28181/sip"
	"github.com/Monibuca/plugin-gb28181/utils"
)

//根据参数构建各种消息
//参数来自于session/transaction等会话管理器
/*
method:请求方法
transport：UDP/TCP
sipSerial: sip server ID
sipRealm: sip domain
username: 用户名/设备序列号
srcIP： 源IP
srcPort：源端口
expires: 过期时间
cseq：消息序列号，当前对话递增
*/
//构建消息：以客户端（可能是IPC，也可能是SIP Server）的角度
func BuildMessageRequest(method sip.Method, transport, sipSerial, sipRealm, username, srcIP string, srcPort uint16, expires, cseq int, body string) *sip.Message {
	server := fmt.Sprintf("%s@%s", sipSerial, sipRealm)
	client := fmt.Sprintf("%s@%s", username, sipRealm)

	msg := &sip.Message{
		Mode:          sip.SIP_MESSAGE_REQUEST,
		MaxForwards:   70,
		UserAgent:     "IPC",
		Expires:       expires,
		ContentLength: 0,
	}
	msg.StartLine = &sip.StartLine{
		Method: method,
		Uri:    sip.NewURI(server),
	}
	msg.Via = &sip.Via{
		Transport: transport,
		Host:      client,
	}
	msg.Via.Params = map[string]string{
		"branch": randBranch(),
		"rport":  "-1", //only key,no-value
	}
	msg.From = &sip.Contact{
		Uri:    sip.NewURI(client),
		Params: nil,
	}
	msg.From.Params = map[string]string{
		"tag": utils.RandNumString(10),
	}
	msg.To = &sip.Contact{
		Uri: sip.NewURI(client),
	}
	msg.CallID = utils.RandNumString(8)
	msg.CSeq = &sip.CSeq{
		ID:     uint32(cseq),
		Method: method,
	}

	msg.Contact = &sip.Contact{
		Uri: sip.NewURI(fmt.Sprintf("%s@%s:%d", username, srcIP, srcPort)),
	}
	if len(body) > 0 {
		msg.ContentLength = len(body)
		msg.Body = body
	}
	return msg
}

//z9hG4bK + 10个随机数字
func randBranch() string {
	return fmt.Sprintf("z9hG4bK%s", utils.RandNumString(8))
}
