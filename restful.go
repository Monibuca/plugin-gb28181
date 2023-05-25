package gb28181

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"m7s.live/engine/v4/util"
)

func (c *GB28181Config) API_list(w http.ResponseWriter, r *http.Request) {
	util.ReturnJson(func() (list []*Device) {
		list = make([]*Device, 0)
		Devices.Range(func(key, value interface{}) bool {
			device := value.(*Device)
			if time.Since(device.UpdateTime) > c.RegisterValidity {
				Devices.Delete(key)
			} else {
				list = append(list, device)
			}
			return true
		})
		return
	}, time.Second*5, w, r)
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
			WriteJSONOk(w, res)
		} else {
			WriteJSON(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		http.NotFound(w, r)
	}
}

func (c *GB28181Config) API_control(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	channel := r.URL.Query().Get("channel")
	ptzcmd := r.URL.Query().Get("ptzcmd")
	if c := FindChannel(id, channel); c != nil {
		w.WriteHeader(c.Control(ptzcmd))
	} else {
		http.NotFound(w, r)
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
		http.NotFound(w, r)
	} else if opt.IsLive() && c.status.Load() > 0 {
		http.Error(w, "live stream already exists", http.StatusNotModified)
	} else if code, err := c.Invite(&opt); err == nil {
		w.WriteHeader(code)
	} else {
		http.Error(w, err.Error(), code)
	}
}

func (c *GB28181Config) API_bye(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	channel := r.URL.Query().Get("channel")
	streamPath := r.URL.Query().Get("streamPath")
	if c := FindChannel(id, channel); c != nil {
		w.WriteHeader(c.Bye(streamPath))
	} else {
		http.NotFound(w, r)
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
		w.WriteHeader(d.MobilePositionSubscribe(id, expiresInt, intervalInt))
	} else {
		http.NotFound(w, r)
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

	util.ReturnJson(func() (list []*DevicePosition) {
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
	}, c.Position.Interval, w, r)
}

func WriteJSONOk(w http.ResponseWriter, data interface{}) {
	WriteJSON(w, data, 200)
}

func WriteJSON(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
