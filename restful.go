package gb28181

import (
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"m7s.live/engine/v4/util"
)

func (conf *GB28181Config) API_list(w http.ResponseWriter, r *http.Request) {
	util.ReturnJson(func() (list []*Device) {
		Devices.Range(func(key, value interface{}) bool {
			device := value.(*Device)
			if time.Since(device.UpdateTime) > conf.RegisterValidity {
				Devices.Delete(key)
			} else {
				list = append(list, device)
			}
			return true
		})
		return
	}, time.Second*5, w, r)
}

func (conf *GB28181Config) API_records(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	channel := r.URL.Query().Get("channel")
	startTime := r.URL.Query().Get("startTime")
	endTime := r.URL.Query().Get("endTime")
	if c := FindChannel(id, channel); c != nil {
		w.WriteHeader(c.QueryRecord(startTime, endTime))
	} else {
		http.NotFound(w, r)
	}
}

func (conf *GB28181Config) API_control(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	channel := r.URL.Query().Get("channel")
	ptzcmd := r.URL.Query().Get("ptzcmd")
	if c := FindChannel(id, channel); c != nil {
		w.WriteHeader(c.Control(ptzcmd))
	} else {
		http.NotFound(w, r)
	}
}

func (conf *GB28181Config) API_invite(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	id := query.Get("id")
	channel := query.Get("channel")
	port, _ := strconv.Atoi(query.Get("mediaPort"))
	opt := InviteOptions{
		dump:      query.Get("dump"),
		MediaPort: uint16(port),
	}
	opt.Validate(query.Get("startTime"), query.Get("endTime"))
	if c := FindChannel(id, channel); c == nil {
		http.NotFound(w, r)
	} else if opt.IsLive() && c.LivePublisher != nil {
		w.WriteHeader(304) //直播流已存在
	} else if code, err := c.Invite(opt); err == nil {
		w.WriteHeader(code)
	} else {
		http.Error(w, err.Error(), code)
	}
}

func (conf *GB28181Config) API_replay(w http.ResponseWriter, r *http.Request) {
	dump := r.URL.Query().Get("dump")
	printOut := r.URL.Query().Get("print")
	if dump == "" {
		dump = conf.DumpPath
	}
	f, err := os.OpenFile(dump, os.O_RDONLY, 0644)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		streamPath := dump
		if strings.HasPrefix(dump, "/") {
			streamPath = "replay" + dump
		} else {
			streamPath = "replay/" + dump
		}
		var pub GBPublisher
		pub.SetIO(f)
		if err = plugin.Publish(streamPath, &pub); err == nil {
			if printOut != "" {
				pub.dumpPrint = w
				pub.SetParentCtx(r.Context())
				err = pub.Replay(f)
			} else {
				go pub.Replay(f)
				w.Write([]byte("ok"))
			}
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func (conf *GB28181Config) API_bye(w http.ResponseWriter, r *http.Request) {
	// CORS(w, r)
	id := r.URL.Query().Get("id")
	channel := r.URL.Query().Get("channel")
	live := r.URL.Query().Get("live")
	if c := FindChannel(id, channel); c != nil {
		w.WriteHeader(c.Bye(live != "false"))
	} else {
		http.NotFound(w, r)
	}
}

func (conf *GB28181Config) API_position(w http.ResponseWriter, r *http.Request) {
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
		expiresInt = conf.Position.Expires
	}
	intervalInt, err := time.ParseDuration(interval)
	if interval == "" || err != nil {
		intervalInt = conf.Position.Interval
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

func (conf *GB28181Config) API_get_position(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	//设备id
	id := query.Get("id")

	util.ReturnJson(func() (list []*DevicePosition) {
		if id == "" {
			Devices.Range(func(key, value interface{}) bool {
				d := value.(*Device)
				if time.Since(d.GpsTime) <= time.Duration(conf.Position.Interval)*time.Second {
					list = append(list, &DevicePosition{ID: d.ID, GpsTime: d.GpsTime, Longitude: d.Longitude, Latitude: d.Latitude})
				}
				return true
			})
		} else if v, ok := Devices.Load(id); ok {
			d := v.(*Device)
			list = append(list, &DevicePosition{ID: d.ID, GpsTime: d.GpsTime, Longitude: d.Longitude, Latitude: d.Latitude})
		}
		return
	}, time.Second*time.Duration(conf.Position.Interval), w, r)
}
