package gb28181

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"m7s.live/engine/v4/util"
)

var (
	playScaleValues = map[float32]bool{0.25: true, 0.5: true, 1: true, 2: true, 4: true}
)

func (c *GB28181Config) API_list(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	if query.Get("interval") == "" {
		query.Set("interval", "5s")
	}
	util.ReturnFetchValue(func() (list []*Device) {
		list = make([]*Device, 0)
		Devices.Range(func(key, value interface{}) bool {
			list = append(list, value.(*Device))
			return true
		})
		return
	}, w, r)
}

func (c *GB28181Config) API_records(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	id := query.Get("id")
	channel := query.Get("channel")
	startTime := query.Get("startTime")
	endTime := query.Get("endTime")
	trange := strings.Split(query.Get("range"), "-")
	if len(trange) == 2 {
		startTime = trange[0]
		endTime = trange[1]
	}
	if c := FindChannel(id, channel); c != nil {
		res, err := c.QueryRecord(startTime, endTime)
		if err == nil {
			util.ReturnValue(res, w, r)
		} else {
			util.ReturnError(util.APIErrorInternal, err.Error(), w, r)
		}
	} else {
		util.ReturnError(util.APIErrorNotFound, fmt.Sprintf("device %q channel %q not found", id, channel), w, r)
	}
}

func (c *GB28181Config) API_control(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	channel := r.URL.Query().Get("channel")
	ptzcmd := r.URL.Query().Get("ptzcmd")
	if c := FindChannel(id, channel); c != nil {
		util.ReturnError(0, fmt.Sprintf("control code:%d", c.Control(ptzcmd)), w, r)
	} else {
		util.ReturnError(util.APIErrorNotFound, fmt.Sprintf("device %q channel %q not found", id, channel), w, r)
	}
}

func (c *GB28181Config) API_ptz(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	id := q.Get("id")
	channel := q.Get("channel")
	cmd := q.Get("cmd")   // 命令名称，见 ptz.go name2code 定义
	hs := q.Get("hSpeed") // 水平速度
	vs := q.Get("vSpeed") // 垂直速度
	zs := q.Get("zSpeed") // 缩放速度

	hsN, err := strconv.ParseUint(hs, 10, 8)
	if err != nil {
		util.ReturnError(util.APIErrorQueryParse, "hSpeed parameter is invalid", w, r)
		return
	}
	vsN, err := strconv.ParseUint(vs, 10, 8)
	if err != nil {
		util.ReturnError(util.APIErrorQueryParse, "vSpeed parameter is invalid", w, r)
		return
	}
	zsN, err := strconv.ParseUint(zs, 10, 8)
	if err != nil {
		util.ReturnError(util.APIErrorQueryParse, "zSpeed parameter is invalid", w, r)
		return
	}

	ptzcmd, err := toPtzStrByCmdName(cmd, uint8(hsN), uint8(vsN), uint8(zsN))
	if err != nil {
		util.ReturnError(util.APIErrorQueryParse, err.Error(), w, r)
		return
	}
	if c := FindChannel(id, channel); c != nil {
		code := c.Control(ptzcmd)
		util.ReturnError(code, "device received", w, r)
	} else {
		util.ReturnError(util.APIErrorNotFound, fmt.Sprintf("device %q channel %q not found", id, channel), w, r)
	}
}

func (c *GB28181Config) API_invite(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	id := query.Get("id")
	channel := query.Get("channel")
	streamPath := query.Get("streamPath")
	port, _ := strconv.Atoi(query.Get("mediaPort"))
	opt := InviteOptions{
		dump:       query.Get("dump"),
		MediaPort:  uint16(port),
		StreamPath: streamPath,
	}
	startTime := query.Get("startTime")
	endTime := query.Get("endTime")
	trange := strings.Split(query.Get("range"), "-")
	if len(trange) == 2 {
		startTime = trange[0]
		endTime = trange[1]
	}
	opt.Validate(startTime, endTime)
	if c := FindChannel(id, channel); c == nil {
		util.ReturnError(util.APIErrorNotFound, fmt.Sprintf("device %q channel %q not found", id, channel), w, r)
	} else if opt.IsLive() && c.State.Load() > 0 {
		util.ReturnError(util.APIErrorQueryParse, "live stream already exists", w, r)
	} else if code, err := c.Invite(&opt); err == nil {
		if code == 200 {
			util.ReturnOK(w, r)
		} else {
			util.ReturnError(util.APIErrorInternal, fmt.Sprintf("invite return code %d", code), w, r)
		}
	} else {
		util.ReturnError(util.APIErrorInternal, err.Error(), w, r)
	}
}

func (c *GB28181Config) API_bye(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	channel := r.URL.Query().Get("channel")
	streamPath := r.URL.Query().Get("streamPath")
	if c := FindChannel(id, channel); c != nil {
		util.ReturnError(0, fmt.Sprintf("bye code:%d", c.Bye(streamPath)), w, r)
	} else {
		util.ReturnError(util.APIErrorNotFound, fmt.Sprintf("device %q channel %q not found", id, channel), w, r)
	}
}

func (c *GB28181Config) API_play_pause(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	channel := r.URL.Query().Get("channel")
	streamPath := r.URL.Query().Get("streamPath")
	if c := FindChannel(id, channel); c != nil {
		util.ReturnError(0, fmt.Sprintf("pause code:%d", c.Pause(streamPath)), w, r)
	} else {
		util.ReturnError(util.APIErrorNotFound, fmt.Sprintf("device %q channel %q not found", id, channel), w, r)
	}
}

func (c *GB28181Config) API_play_resume(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	channel := r.URL.Query().Get("channel")
	streamPath := r.URL.Query().Get("streamPath")
	if c := FindChannel(id, channel); c != nil {
		util.ReturnError(0, fmt.Sprintf("resume code:%d", c.Resume(streamPath)), w, r)
	} else {
		util.ReturnError(util.APIErrorNotFound, fmt.Sprintf("device %q channel %q not found", id, channel), w, r)
	}
}

func (c *GB28181Config) API_play_seek(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	channel := r.URL.Query().Get("channel")
	streamPath := r.URL.Query().Get("streamPath")
	secStr := r.URL.Query().Get("second")
	sec, err := strconv.ParseUint(secStr, 10, 32)
	if err != nil {
		util.ReturnError(util.APIErrorQueryParse, "second parameter is invalid: "+err.Error(), w, r)
		return
	}
	if c := FindChannel(id, channel); c != nil {
		util.ReturnError(0, fmt.Sprintf("play code:%d", c.PlayAt(streamPath, uint(sec))), w, r)
	} else {
		util.ReturnError(util.APIErrorNotFound, fmt.Sprintf("device %q channel %q not found", id, channel), w, r)
	}
}

func (c *GB28181Config) API_play_forward(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	channel := r.URL.Query().Get("channel")
	streamPath := r.URL.Query().Get("streamPath")
	speedStr := r.URL.Query().Get("speed")
	speed, err := strconv.ParseFloat(speedStr, 32)
	secondErrMsg := "speed parameter is invalid, should be one of 0.25,0.5,1,2,4"
	if err != nil || !playScaleValues[float32(speed)] {
		util.ReturnError(util.APIErrorQueryParse, secondErrMsg, w, r)
		return
	}
	if c := FindChannel(id, channel); c != nil {
		util.ReturnError(0, fmt.Sprintf("playforward code:%d", c.PlayForward(streamPath, float32(speed))), w, r)
	} else {
		util.ReturnError(util.APIErrorNotFound, fmt.Sprintf("device %q channel %q not found", id, channel), w, r)
	}
}

func (c *GB28181Config) API_position(w http.ResponseWriter, r *http.Request) {
	//CORS(w, r)
	query := r.URL.Query()
	//设备id
	id := query.Get("id")
	//订阅周期(单位：秒)
	expires := query.Get("expires")
	//订阅间隔（单位：秒）
	interval := query.Get("interval")

	expiresInt, err := time.ParseDuration(expires)
	if expires == "" || err != nil {
		expiresInt = c.Position.Expires
	}
	intervalInt, err := time.ParseDuration(interval)
	if interval == "" || err != nil {
		intervalInt = c.Position.Interval
	}

	if v, ok := Devices.Load(id); ok {
		d := v.(*Device)
		util.ReturnError(0, fmt.Sprintf("mobileposition code:%d", d.MobilePositionSubscribe(id, expiresInt, intervalInt)), w, r)
	} else {
		util.ReturnError(util.APIErrorNotFound, fmt.Sprintf("device %q  not found", id), w, r)
	}
}

type DevicePosition struct {
	ID        string
	GpsTime   time.Time //gps时间
	Longitude string    //经度
	Latitude  string    //纬度
}

func (c *GB28181Config) API_get_position(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	//设备id
	id := query.Get("id")
	if query.Get("interval") == "" {
		query.Set("interval", fmt.Sprintf("%ds", c.Position.Interval.Seconds()))
	}
	util.ReturnFetchValue(func() (list []*DevicePosition) {
		if id == "" {
			Devices.Range(func(key, value interface{}) bool {
				d := value.(*Device)
				if time.Since(d.GpsTime) <= c.Position.Interval {
					list = append(list, &DevicePosition{ID: d.ID, GpsTime: d.GpsTime, Longitude: d.Longitude, Latitude: d.Latitude})
				}
				return true
			})
		} else if v, ok := Devices.Load(id); ok {
			d := v.(*Device)
			list = append(list, &DevicePosition{ID: d.ID, GpsTime: d.GpsTime, Longitude: d.Longitude, Latitude: d.Latitude})
		}
		return
	}, w, r)
}
