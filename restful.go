package gb28181

import (
	"net/http"
)

func restful() {
	http.HandleFunc("/api/gb28181/query/records", func(w http.ResponseWriter, r *http.Request) {
		// CORS(w, r)
		id := r.URL.Query().Get("id")
		channel := r.URL.Query().Get("channel")
		startTime := r.URL.Query().Get("startTime")
		endTime := r.URL.Query().Get("endTime")
		if c := FindChannel(id, channel); c != nil {
			w.WriteHeader(c.QueryRecord(startTime, endTime))
		} else {
			w.WriteHeader(404)
		}
	})
	// http.HandleFunc("/api/gb28181/list", func(w http.ResponseWriter, r *http.Request) {
	// 	// CORS(w, r)
	// 	sse := NewSSE(w, r.Context())
	// 	for {
	// 		var list []*Device
	// 		Devices.Range(func(key, value interface{}) bool {
	// 			device := value.(*Device)
	// 			if time.Since(device.UpdateTime) > time.Duration(serverConfig.RegisterValidity)*time.Second {
	// 				Devices.Delete(key)
	// 			} else {
	// 				list = append(list, device)
	// 			}
	// 			return true
	// 		})
	// 		sse.WriteJSON(list)
	// 		select {
	// 		case <-time.After(time.Second * 5):
	// 		case <-sse.Done():
	// 			return
	// 		}
	// 	}
	// })
	http.HandleFunc("/api/gb28181/control", func(w http.ResponseWriter, r *http.Request) {
		// CORS(w, r)
		id := r.URL.Query().Get("id")
		channel := r.URL.Query().Get("channel")
		ptzcmd := r.URL.Query().Get("ptzcmd")
		if c := FindChannel(id, channel); c != nil {
			w.WriteHeader(c.Control(ptzcmd))
		} else {
			w.WriteHeader(404)
		}
	})
	http.HandleFunc("/api/gb28181/invite", func(w http.ResponseWriter, r *http.Request) {
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
	})
	http.HandleFunc("/api/gb28181/bye", func(w http.ResponseWriter, r *http.Request) {
		// CORS(w, r)
		id := r.URL.Query().Get("id")
		channel := r.URL.Query().Get("channel")
		live := r.URL.Query().Get("live")
		if c := FindChannel(id, channel); c != nil {
			w.WriteHeader(c.Bye(live != "false"))
		} else {
			w.WriteHeader(404)
		}
	})
}
