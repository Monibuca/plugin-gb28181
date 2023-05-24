package gb28181

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

func (c *GB28181Config) API_list(w http.ResponseWriter, r *http.Request) {
	list := make([]*Device, 0)
	Devices.Range(func(key, value interface{}) bool {
		device := value.(*Device)
		if time.Since(device.UpdateTime) > c.RegisterValidity {
			Devices.Delete(key)
		} else {
			list = append(list, device)
		}
		return true
	})
	WriteJSONOk(w, list)
}

func (c *GB28181Config) API_records(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	channel := r.URL.Query().Get("channel")
	startTime := r.URL.Query().Get("startTime")
	endTime := r.URL.Query().Get("endTime")
	if c := FindChannel(id, channel); c != nil {
		res, err := c.QueryRecord(startTime, endTime)
		if err == nil {
			WriteJSONOk(w, res)
		} else {
			WriteJSON(w, err.Error(), 500)
		}
	} else {
		WriteJSON(w, fmt.Sprintf("device %q channel %q not found", id, channel), 404)
	}
}

func (c *GB28181Config) API_control(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	channel := r.URL.Query().Get("channel")
	ptzcmd := r.URL.Query().Get("ptzcmd")
	if c := FindChannel(id, channel); c != nil {
		code := c.Control(ptzcmd)
		WriteJSON(w, "", code)
	} else {
		WriteJSON(w, fmt.Sprintf("device %q channel %q not found", id, channel), 404)
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
	opt.Validate(query.Get("startTime"), query.Get("endTime"))
	if c := FindChannel(id, channel); c == nil {
		WriteJSON(w, fmt.Sprintf("device %q channel %q not found", id, channel), 404)
	} else if opt.IsLive() && c.status.Load() > 0 {
		WriteJSON(w, "live stream already exists", 304) //直播流已存在
	} else if code, err := c.Invite(&opt); err == nil {
		WriteJSON(w, "", code)
	} else {
		WriteJSON(w, err.Error(), code)
	}
}

func (c *GB28181Config) API_bye(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	channel := r.URL.Query().Get("channel")
	streamPath := r.URL.Query().Get("streamPath")
	if c := FindChannel(id, channel); c != nil {
		code := c.Bye(streamPath)
		WriteJSON(w, "", code)
	} else {
		WriteJSON(w, "stream dose not exists", 404)
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
		code := d.MobilePositionSubscribe(id, expiresInt, intervalInt)
		WriteJSON(w, "", code)
	} else {
		WriteJSON(w, "device does not exist.", 404)
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

	var list []*DevicePosition
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
	WriteJSONOk(w, list)
}

func WriteJSONOk(w http.ResponseWriter, data interface{}) {
	WriteJSON(w, data, 200)
}

func WriteJSON(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
