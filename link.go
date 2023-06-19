package gb28181

import (
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"
)

// 对于录像查询，通过 queryKey (即 deviceId + channelId + sn) 唯一区分一次请求和响应
// 并将其关联起来，以实现异步响应的目的
// 提供单例实例供调用
var RecordQueryLink = NewRecordQueryLink(time.Second * 60)

type recordQueryLink struct {
	pendingResult map[string]recordQueryResult // queryKey 查询结果缓存
	pendingResp   map[string]recordQueryResp   // queryKey 待回复的查询请求
	timeout       time.Duration                // 查询结果的过期时间
	sync.RWMutex
}

type recordQueryResult struct {
	time     time.Time
	err      error
	sum      int
	finished bool
	list     []*Record
}
type recordQueryResp struct {
	respChan  chan<- recordQueryResult
	timeout   time.Duration
	startTime time.Time
}

func NewRecordQueryLink(resultTimeout time.Duration) *recordQueryLink {
	c := &recordQueryLink{
		timeout:       resultTimeout,
		pendingResult: make(map[string]recordQueryResult),
		pendingResp:   make(map[string]recordQueryResp),
	}
	return c
}

// 唯一区分一次录像查询
func recordQueryKey(deviceId, channelId string, sn int) string {
	return fmt.Sprintf("%s-%s-%d", deviceId, channelId, sn)
}

// 定期清理过期的查询结果和请求
func (c *recordQueryLink) cleanTimeout() {
	for k, s := range c.pendingResp {
		if time.Since(s.startTime) > s.timeout {
			if r, ok := c.pendingResult[k]; ok {
				c.notify(k, r)
			} else {
				c.notify(k, recordQueryResult{err: fmt.Errorf("query time out")})
			}
		}
	}
	for k, r := range c.pendingResult {
		if time.Since(r.time) > c.timeout {
			delete(c.pendingResult, k)
		}
	}
}

func (c *recordQueryLink) Put(deviceId, channelId string, sn int, sum int, record []*Record) {
	key, r := c.doPut(deviceId, channelId, sn, sum, record)
	if r.finished {
		c.notify(key, r)
	}
}

func (c *recordQueryLink) doPut(deviceId, channelId string, sn, sum int, record []*Record) (key string, r recordQueryResult) {
	c.Lock()
	defer c.Unlock()
	key = recordQueryKey(deviceId, channelId, sn)
	if v, ok := c.pendingResult[key]; ok {
		r = v
	} else {
		r = recordQueryResult{time: time.Now(), sum: sum, list: make([]*Record, 0)}
	}

	r.list = append(r.list, record...)
	if len(r.list) == sum {
		r.finished = true
	}
	c.pendingResult[key] = r
	GB28181Plugin.Logger.Debug("put record",
		zap.String("key", key),
		zap.Int("sum", sum),
		zap.Int("count", len(r.list)))
	return
}

func (c *recordQueryLink) WaitResult(
	deviceId, channelId string, sn int,
	timeout time.Duration) (resultCh <-chan recordQueryResult) {

	key := recordQueryKey(deviceId, channelId, sn)
	c.Lock()
	defer c.Unlock()
	respCh := make(chan recordQueryResult, 1)
	resultCh = respCh
	c.pendingResp[key] = recordQueryResp{startTime: time.Now(), timeout: timeout, respChan: respCh}
	return
}

func (c *recordQueryLink) notify(key string, r recordQueryResult) {
	if s, ok := c.pendingResp[key]; ok {
		s.respChan <- r
	}
	c.Lock()
	defer c.Unlock()
	delete(c.pendingResp, key)
	delete(c.pendingResult, key)
	GB28181Plugin.Logger.Debug("record notify", zap.String("key", key))
}
