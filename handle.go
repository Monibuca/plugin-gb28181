package gb28181

import (
	"bytes"
	"encoding/xml"
	"github.com/Monibuca/plugin-gb28181/v3/sip"
	"github.com/Monibuca/plugin-gb28181/v3/transaction"
	"github.com/Monibuca/plugin-gb28181/v3/utils"
	"github.com/logrusorgru/aurora"

	. "github.com/Monibuca/utils/v3"
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

			if dc, ok := DeviceRegisterCount.LoadOrStore(id, 1); ok && dc.(int) > MaxRegisterCount {
				var response sip.Response
				response.Message = req.BuildResponse(http.StatusForbidden)
				_ = tx.Respond(&response)
				return
			} else {
				// 设备第二次上报，校验
				_nonce, loaded := DeviceNonce.Load(id)
				if loaded && req.Authorization.Verify(username, config.Password, config.Realm, _nonce.(string)) {
					passAuth = true
				} else {
					DeviceRegisterCount.Store(id, dc.(int)+1)
				}
			}
		}

	}
	if passAuth {
		storeDevice(id, tx.Core, req.Message)
		DeviceNonce.Delete(id)
		DeviceRegisterCount.Delete(id)
		m := req.BuildOK()
		resp := &sip.Response{Message: m}
		_ = tx.Respond(resp)

	} else {
		var response sip.Response
		response.Message = req.BuildResponseWithPhrase(401, "Unauthorized")
		_nonce, _ := DeviceNonce.LoadOrStore(id, utils.RandNumString(32))
		response.WwwAuthenticate = sip.NewWwwAuthenticate(config.Realm, _nonce.(string), sip.DIGEST_ALGO_MD5)
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
			go d.QueryDeviceInfo(req)
		}
		d.UpdateTime = time.Now()
		temp := &struct {
			XMLName      xml.Name
			CmdType      string
			DeviceID     string
			DeviceName   string
			Manufacturer string
			Model        string
			Channel      string
			DeviceList   []*Channel `xml:"DeviceList>Item"`
			RecordList   []*Record  `xml:"RecordList>Item"`
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
		var body string
		switch temp.CmdType {
		case "Keepalive":
			d.LastKeepaliveAt = time.Now()
			//callID !="" 说明是订阅的事件类型信息
			if d.Channels == nil || d.subscriber.CallID != "" && d.LastKeepaliveAt.After(d.subscriber.Timeout) {
				go d.Catalog()
			}
			d.CheckSubStream()
		case "Catalog":
			d.UpdateChannels(temp.DeviceList)
		case "RecordInfo":
			d.UpdateRecord(temp.DeviceID, temp.RecordList)
		case "DeviceInfo":
			// 主设备信息
			d.Name = temp.DeviceName
			d.Manufacturer = temp.Manufacturer
			d.Model = temp.Model
		case "Alarm":
			d.Status = "Alarmed"
			body = sip.BuildAlarmResponseXML(d.ID)
		default:
			Println("DeviceID:", aurora.Red(d.ID), " Not supported CmdType : "+temp.CmdType+" body:\n", req.Body)
			response := &sip.Response{req.BuildResponse(http.StatusBadRequest)}
			tx.Respond(response)
			return
		}

		buildOK := req.BuildOK()
		buildOK.Body = body
		response := &sip.Response{buildOK}
		tx.Respond(response)
	}
}
