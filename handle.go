package gb28181

import (
	"bytes"
	"encoding/xml"
	"github.com/Monibuca/plugin-gb28181/v3/sip"
	"github.com/Monibuca/plugin-gb28181/v3/transaction"
	"github.com/Monibuca/plugin-gb28181/v3/utils"

	. "github.com/Monibuca/utils/v3"
	. "github.com/logrusorgru/aurora"
	"golang.org/x/net/html/charset"
	"net/http"
	"time"
)

func OnRegister(req *sip.Request, tx *transaction.GBTx) {
	id := req.From.Uri.UserInfo()

	passAuth := false
	// 不需要密码情况
	if config.Username == "" && config.Password == "" {
		passAuth = true
	} else {
		// 需要密码情况 设备第一次上报，返回401和加密算法
		if req.Authorization != nil && req.Authorization.GetUsername() != "" {
			// 有些摄像头没有配置用户名的地方，用户名就是摄像头自己的国标id
			var username string
			if req.Authorization.GetUsername() == id {
				username = id
			} else {
				username = config.Username
			}

			if DeviceRegisterCount[id] >= MaxRegisterCount {
				var response sip.Response
				response.Message = req.BuildResponse(http.StatusForbidden)
				_ = tx.Respond(&response)
				return
			} else {
				// 设备第二次上报，校验
				if req.Authorization.Verify(username, config.Password, config.Realm, DeviceNonce[id]) {
					passAuth = true
				} else {
					DeviceRegisterCount[id]++
				}
			}
		}

	}
	if passAuth {
		storeDevice(id, tx.Core, req.Message)
		delete(DeviceNonce, id)
		delete(DeviceRegisterCount, id)
		m := req.BuildOK()
		resp := &sip.Response{Message: m}
		_ = tx.Respond(resp)
	} else {
		var response sip.Response
		response.Message = req.BuildResponseWithPhrase(401, "Unauthorized")
		if DeviceNonce[id] == "" {
			nonce := utils.RandNumString(32)
			DeviceNonce[id] = nonce
		}
		response.WwwAuthenticate = sip.NewWwwAuthenticate(config.Realm, DeviceNonce[id], sip.DIGEST_ALGO_MD5)
		response.SourceAdd = req.DestAdd
		response.DestAdd = req.SourceAdd
		_ = tx.Respond(&response)
	}
}
func OnMessage(req *sip.Request, tx *transaction.GBTx) {

	if v, ok := Devices.Load(req.From.Uri.UserInfo()); ok {
		d := v.(*Device)
		if d.Status == string(sip.REGISTER) {
			d.Status = "ONLINE"
			go d.Query(req)
		}
		d.UpdateTime = time.Now()
		temp := &struct {
			XMLName    xml.Name
			CmdType    string
			DeviceID   string
			DeviceList []*Channel `xml:"DeviceList>Item"`
			RecordList []*Record  `xml:"RecordList>Item"`
		}{}
		decoder := xml.NewDecoder(bytes.NewReader([]byte(req.Body)))
		decoder.CharsetReader = charset.NewReaderLabel
		err := decoder.Decode(temp)
		if err != nil {
			err = utils.DecodeGbk(temp, []byte(req.Body))
			if err != nil {
				Printf("decode catelog err: %s", err)
			}
		}
		switch temp.CmdType {
		case "Keepalive":
			if d.subscriber.CallID != "" && time.Now().After(d.subscriber.Timeout) {
				go d.Subscribe(req)
			}
			d.CheckSubStream()
			break
		case "Catalog":
			d.UpdateChannels(temp.DeviceList)
			break
		case "RecordInfo":
			d.UpdateRecord(temp.DeviceID, temp.RecordList)
			break
		case "DeviceInfo":
			// 主设备信息
			break
		default:
			Println(Red("Not supported CmdType"), temp.CmdType)
			response := &sip.Response{req.BuildResponse(http.StatusBadRequest)}
			tx.Respond(response)
			return
		}
		response := &sip.Response{req.BuildOK()}
		tx.Respond(response)
	}
}
