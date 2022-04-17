package gb28181

import (
	"net/http"
	"strconv"
	"time"

	"m7s.live/engine/v4/util"
)

func (conf *GB28181Config) API_list(w http.ResponseWriter, r *http.Request) {
	sse := util.NewSSE(w, r.Context())

	for {
		var list []*Device
		Devices.Range(func(key, value interface{}) bool {
			device := value.(*Device)
			if time.Since(device.UpdateTime) > time.Duration(conf.RegisterValidity)*time.Second {
				Devices.Delete(key)
			} else {
				list = append(list, device)
			}
			return true
		})
		sse.WriteJSON(list)
		select {
		case <-time.After(time.Second * 5):
		case <-sse.Done():
			return
		}
	}
}

func (conf *GB28181Config) API_records(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	channel := r.URL.Query().Get("channel")
	startTime := r.URL.Query().Get("startTime")
	endTime := r.URL.Query().Get("endTime")
	if c := FindChannel(id, channel); c != nil {
		w.WriteHeader(c.QueryRecord(startTime, endTime))
	} else {
		w.WriteHeader(404)
	}
}

func (conf *GB28181Config) API_control(w http.ResponseWriter, r *http.Request) {
	// CORS(w, r)
	id := r.URL.Query().Get("id")
	channel := r.URL.Query().Get("channel")
	ptzcmd := r.URL.Query().Get("ptzcmd")
	if c := FindChannel(id, channel); c != nil {
		w.WriteHeader(c.Control(ptzcmd))
	} else {
		w.WriteHeader(404)
	}
}

func (conf *GB28181Config) API_invite(w http.ResponseWriter, r *http.Request) {
	// CORS(w, r)
	query := r.URL.Query()
	id := query.Get("id")
	channel := r.URL.Query().Get("channel")
	startTime := query.Get("startTime")
	endTime := query.Get("endTime")
	if c := FindChannel(id, channel); c != nil {
		if startTime == "" && c.LivePublisher != nil {
			w.WriteHeader(304) //直播流已存在
		} else {
			w.WriteHeader(c.Invite(startTime, endTime))
		}
	} else {
		w.WriteHeader(404)
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
		w.WriteHeader(404)
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

	expiresInt, _ := strconv.Atoi(expires)
	intervalInt, _ := strconv.Atoi(interval)

	if v, ok := Devices.Load(id); ok {
		d := v.(*Device)
		w.WriteHeader(d.MobilePositionSubscribe(id, expiresInt, intervalInt))
	} else {
		w.WriteHeader(404)
	}
}
