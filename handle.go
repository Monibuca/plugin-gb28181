package gb28181

import (
	"bytes"
	"crypto/md5"
	"encoding/xml"
	"fmt"

	"github.com/logrusorgru/aurora"
	"go.uber.org/zap"
	"m7s.live/plugin/gb28181/v4/utils"

	"github.com/ghettovoice/gosip/sip"

	"net/http"
	"time"

	"golang.org/x/net/html/charset"
)

type Authorization struct {
	*sip.Authorization
}

func (a *Authorization) Verify(username, passwd, realm, nonce string) bool {

	//1、将 username,realm,password 依次组合获取 1 个字符串，并用算法加密的到密文 r1
	s1 := fmt.Sprintf("%s:%s:%s", username, realm, passwd)
	r1 := a.getDigest(s1)
	//2、将 method，即REGISTER ,uri 依次组合获取 1 个字符串，并对这个字符串使用算法 加密得到密文 r2
	s2 := fmt.Sprintf("REGISTER:%s", a.Uri())
	r2 := a.getDigest(s2)

	if r1 == "" || r2 == "" {
		fmt.Println("Authorization algorithm wrong")
		return false
	}
	//3、将密文 1，nonce 和密文 2 依次组合获取 1 个字符串，并对这个字符串使用算法加密，获得密文 r3，即Response
	s3 := fmt.Sprintf("%s:%s:%s", r1, nonce, r2)
	r3 := a.getDigest(s3)

	//4、计算服务端和客户端上报的是否相等
	return r3 == a.Response()
}

func (a *Authorization) getDigest(raw string) string {
	switch a.Algorithm() {
	case "MD5":
		return fmt.Sprintf("%x", md5.Sum([]byte(raw)))
	default: //如果没有算法，默认使用MD5
		return fmt.Sprintf("%x", md5.Sum([]byte(raw)))
	}
}

func (config *GB28181Config) OnRegister(req sip.Request, tx sip.ServerTransaction) {
	from, _ := req.From()
	id := from.Address.User().String()
	plugin.Debug(id)
	via, _ := req.Via()
	fmt.Printf("%s", via.Value())
	passAuth := false
	// 不需要密码情况
	if config.Username == "" && config.Password == "" {
		passAuth = true
	} else {
		// 需要密码情况 设备第一次上报，返回401和加密算法
		if hdrs := req.GetHeaders("Authorization"); len(hdrs) > 0 {
			authenticateHeader := hdrs[0].(*sip.GenericHeader)
			auth := &Authorization{sip.AuthFromValue(authenticateHeader.Contents)}

			// 有些摄像头没有配置用户名的地方，用户名就是摄像头自己的国标id
			var username string
			if auth.Username() == id {
				username = id
			} else {
				username = config.Username
			}

			if dc, ok := DeviceRegisterCount.LoadOrStore(id, 1); ok && dc.(int) > MaxRegisterCount {
				response := sip.NewResponseFromRequest("", req, http.StatusForbidden, "Forbidden", "")
				tx.Respond(response)
				return
			} else {
				// 设备第二次上报，校验
				_nonce, loaded := DeviceNonce.Load(id)
				if loaded && auth.Verify(username, config.Password, config.Realm, _nonce.(string)) {
					passAuth = true
				} else {
					DeviceRegisterCount.Store(id, dc.(int)+1)
				}
			}
		}
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
	id := from.Address.User().String()
	if v, ok := Devices.Load(id); ok {
		d := v.(*Device)
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
				plugin.Error("decode catelog err", zap.Error(err))
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
			plugin.Sugar().Warnf("DeviceID:", aurora.Red(d.ID), " Not supported CmdType : "+temp.CmdType+" body:\n", req.Body)
			response := sip.NewResponseFromRequest("", req, http.StatusBadRequest, "", "")
			tx.Respond(response)
			return
		}

		tx.Respond(sip.NewResponseFromRequest("", req, http.StatusOK, "OK", body))
	}
}
func (config *GB28181Config) onBye(req sip.Request, tx sip.ServerTransaction) {
	tx.Respond(sip.NewResponseFromRequest("", req, http.StatusOK, "OK", ""))
}

// OnNotify 订阅通知处理
func (config *GB28181Config) OnNotify(req sip.Request, tx sip.ServerTransaction) {
	from, _ := req.From()
	id := from.Address.User().String()
	if v, ok := Devices.Load(id); ok {
		d := v.(*Device)
		d.UpdateTime = time.Now()
		temp := &struct {
			XMLName   xml.Name
			CmdType   string
			DeviceID  string
			Time      string //位置订阅-GPS时间
			Longitude string //位置订阅-经度
			Latitude  string //位置订阅-维度
			// Speed      string           //位置订阅-速度(km/h)(可选)
			// Direction  string           //位置订阅-方向(取值为当前摄像头方向与正北方的顺时针夹角,取值范围0°~360°,单位:°)(可选)
			// Altitude   string           //位置订阅-海拔高度,单位:m(可选)
			DeviceList []*notifyMessage `xml:"DeviceList>Item"` //目录订阅
		}{}
		decoder := xml.NewDecoder(bytes.NewReader([]byte(req.Body())))
		decoder.CharsetReader = charset.NewReaderLabel
		err := decoder.Decode(temp)
		if err != nil {
			err = utils.DecodeGbk(temp, []byte(req.Body()))
			if err != nil {
				plugin.Error("decode catelog err", zap.Error(err))
			}
		}
		var body string
		switch temp.CmdType {
		case "Catalog":
			//目录状态
			d.UpdateChannelStatus(temp.DeviceList)
		case "MobilePosition":
			//更新channel的坐标
			d.UpdateChannelPosition(temp.DeviceID, temp.Time, temp.Longitude, temp.Latitude)
		// case "Alarm":
		// 	//报警事件通知 TODO
		default:
			plugin.Sugar().Warnf("DeviceID:", aurora.Red(d.ID), " Not supported CmdType : "+temp.CmdType+" body:", req.Body)
			response := sip.NewResponseFromRequest("", req, http.StatusBadRequest, "", "")
			tx.Respond(response)
			return
		}

		tx.Respond(sip.NewResponseFromRequest("", req, http.StatusOK, "OK", body))
	}
}

type notifyMessage struct {
	DeviceID     string
	ParentID     string
	Name         string
	Manufacturer string
	Model        string
	Owner        string
	CivilCode    string
	Address      string
	Parental     int
	SafetyWay    int
	RegisterWay  int
	Secrecy      int
	Status       string
	//状态改变事件 ON:上线,OFF:离线,VLOST:视频丢失,DEFECT:故障,ADD:增加,DEL:删除,UPDATE:更新(必选)
	Event string
}
