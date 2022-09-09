package gb28181

import (
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/pion/rtp/v2"
	"m7s.live/engine/v4/util"
)

func (conf *GB28181Config) API_list(w http.ResponseWriter, r *http.Request) {
	util.ReturnJson(func() (list []*Device) {
		Devices.Range(func(key, value interface{}) bool {
			device := value.(*Device)
			if time.Since(device.UpdateTime) > time.Duration(conf.RegisterValidity)*time.Second {
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
		w.WriteHeader(404)
	}
}

func (conf *GB28181Config) API_control(w http.ResponseWriter, r *http.Request) {
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
	query := r.URL.Query()
	id := query.Get("id")
	channel := query.Get("channel")
	port, _ := strconv.Atoi(query.Get("mediaPort"))
	opt := InviteOptions{
		query.Get("startTime"),
		query.Get("endTime"),
		query.Get("dump"),
		"", 0, uint16(port),
	}
	if c := FindChannel(id, channel); c != nil {
		if opt.IsLive() && c.LivePublisher != nil {
			w.WriteHeader(304) //直播流已存在
		} else {
			w.WriteHeader(c.Invite(opt))
		}
	} else {
		w.WriteHeader(404)
	}
}

func (conf *GB28181Config) API_replay(w http.ResponseWriter, r *http.Request) {
	dump := r.URL.Query().Get("dump")
	if dump == "" {
		dump = conf.DumpPath
	}
	f, err := os.OpenFile(dump, os.O_RDONLY, 0644)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		go func() {
			defer f.Close()
			streamPath := dump
			if strings.HasPrefix(dump, "/") {
				streamPath = "replay" + dump
			} else {
				streamPath = "replay/" + dump
			}
			var pub GBPublisher
			var rtpPacket rtp.Packet
			if err = plugin.Publish(streamPath, &pub); err == nil {
				for l := make([]byte, 6); !pub.IsClosed(); time.Sleep(time.Millisecond * time.Duration(util.ReadBE[uint16](l[4:]))) {
					_, err = f.Read(l)
					if err != nil {
						return
					}
					payload := make([]byte, util.ReadBE[int](l[:4]))
					_, err = f.Read(payload)
					if err != nil {
						return
					}
					rtpPacket.Unmarshal(payload)
					pub.PushPS(&rtpPacket)
				}
			}
		}()
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
