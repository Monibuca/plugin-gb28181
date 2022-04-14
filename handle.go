package gb28181

import (
	"bytes"
	"encoding/xml"
	"fmt"

	"github.com/logrusorgru/aurora"
	"m7s.live/plugin/gb28181/v4/utils"

	"github.com/ghettovoice/gosip/sip"

	"net/http"
	"time"

	"golang.org/x/net/html/charset"
)

func (config *GB28181Config) OnRegister(req sip.Request, tx sip.ServerTransaction) {
	from, _ := req.From()

	id := from.Address.User().String()
	plugin.Debug(id)

	passAuth := false
	// 不需要密码情况
	if config.Username == "" && config.Password == "" {
		passAuth = true
	} else {
		// // 需要密码情况 设备第一次上报，返回401和加密算法
		// if req.Authorization != nil && req.Authorization.GetUsername() != "" {
		// 	// 有些摄像头没有配置用户名的地方，用户名就是摄像头自己的国标id
		// 	var username string
		// 	if req.Authorization.GetUsername() == id {
		// 		username = id
		// 	} else {
		// 		username = config.Username
		// 	}

		// 	if dc, ok := DeviceRegisterCount.LoadOrStore(id, 1); ok && dc.(int) > MaxRegisterCount {
		// 		var response sip.Response
		// 		response.Message = req.BuildResponse(http.StatusForbidden)
		// 		_ = tx.Respond(&response)
		// 		return
		// 	} else {
		// 		// 设备第二次上报，校验
		// 		_nonce, loaded := DeviceNonce.Load(id)
		// 		if loaded && req.Authorization.Verify(username, config.Password, config.Realm, _nonce.(string)) {
		// 			passAuth = true
		// 		} else {
		// 			DeviceRegisterCount.Store(id, dc.(int)+1)
		// 		}
		// 	}
		// }

	}
	if passAuth {
		config.StoreDevice(id, req, &tx)
		DeviceNonce.Delete(id)
		DeviceRegisterCount.Delete(id)
		_ = tx.Respond(sip.NewResponseFromRequest("", req, http.StatusOK, "OK", ""))
	} else {
		response := sip.NewResponseFromRequest("", req, http.StatusUnauthorized, "Unauthorized", "")
		_nonce, _ := DeviceNonce.LoadOrStore(id, utils.RandNumString(32))
		auth := fmt.Sprintf(
			`Digest realm="%s",algorithm=%s,nonce="%s"`,
			config.Realm,
			"MD5",
			_nonce.(string),
		)
		response.AppendHeader(&sip.GenericHeader{
			HeaderName: "WWW-Authenticate",
			Contents:   auth,
		})
		_ = tx.Respond(response)
	}
}
func (config *GB28181Config) OnMessage(req sip.Request, tx sip.ServerTransaction) {
	from, _ := req.From()
	if v, ok := Devices.Load(from.Address.User()); ok {
		d := v.(*Device)
		//d.SourceAddr = req.
		if d.Status == string(sip.REGISTER) {
			d.Status = "ONLINE"
			//go d.QueryDeviceInfo(req)
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
		decoder := xml.NewDecoder(bytes.NewReader([]byte(req.Body())))
		decoder.CharsetReader = charset.NewReaderLabel
		err := decoder.Decode(temp)
		if err != nil {
			err = utils.DecodeGbk(temp, []byte(req.Body()))
			if err != nil {
				fmt.Printf("decode catelog err: %s", err)
			}
		}
		var body string
		switch temp.CmdType {
		case "Keepalive":
			d.LastKeepaliveAt = time.Now()
			//callID !="" 说明是订阅的事件类型信息
			if d.Channels == nil {
				go d.Catalog()
			} else {
				if d.subscriber.CallID != "" && d.LastKeepaliveAt.After(d.subscriber.Timeout) {
					go d.Catalog()
				} else {
					for _, c := range d.Channels {
						if config.AutoInvite &&
							(c.LivePublisher == nil) {
							c.Invite("", "")
						}
					}
				}

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
			body = BuildAlarmResponseXML(d.ID)
		default:
			fmt.Println("DeviceID:", aurora.Red(d.ID), " Not supported CmdType : "+temp.CmdType+" body:\n", req.Body)
			response := sip.NewResponseFromRequest("", req, http.StatusBadRequest, "", "")
			tx.Respond(response)
			return
		}

		_ = tx.Respond(sip.NewResponseFromRequest("", req, http.StatusOK, "OK", body))
	}
}
func (config *GB28181Config) onBye(req sip.Request, tx sip.ServerTransaction) {
	_ = tx.Respond(sip.NewResponseFromRequest("", req, http.StatusOK, "OK", ""))
}
