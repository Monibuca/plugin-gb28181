package transaction

import (
	"net/http"
	"sync"
	"time"

	"m7s.live/plugin-gb28181/v4/sip"
	. "m7s.live/plugin-gb28181/v4/transport"
	"m7s.live/plugin-gb28181/v4/utils"
)

var ActiveTX *GBTxs

// GBTxs a GBTxs stands for a Gb28181 Transaction collection
type GBTxs struct {
	Txs map[string]*GBTx
	RWM *sync.RWMutex
}

func (txs *GBTxs) NewTX(key string, conn Connection) *GBTx {
	tx := NewTransaction(key, conn)
	txs.RWM.Lock()
	txs.Txs[key] = tx
	txs.RWM.Unlock()
	return tx
}

func (txs *GBTxs) GetTX(key string) *GBTx {
	txs.RWM.RLock()
	tx, ok := txs.Txs[key]
	if !ok {
		tx = nil
	}
	txs.RWM.RUnlock()
	return tx
}

func (txs *GBTxs) rmTX(tx *GBTx) {
	txs.RWM.Lock()
	delete(txs.Txs, tx.key)
	txs.RWM.Unlock()
}

// GBTx Gb28181 Transaction
type GBTx struct {
	conn   Connection
	key    string
	resp   chan *sip.Response
	active chan int
	*Core
}

// NewTransaction create a new GBtx
func NewTransaction(key string, conn Connection) *GBTx {
	tx := &GBTx{conn: conn, key: key, resp: make(chan *sip.Response, 10), active: make(chan int, 1)}
	go tx.watch()
	return tx
}

// Key returns the GBTx Key
func (tx *GBTx) Key() string {
	return tx.key
}

func (tx *GBTx) watch() {
	for {
		select {
		case <-tx.active:
			//Println("active tx", tx.Key(), time.Now().Format("2006-01-02 15:04:05"))
		case <-time.After(20 * time.Second):
			tx.Close()
			//Println("watch closed tx", tx.key, time.Now().Format("2006-01-02 15:04:05"))
			return
		}
	}
}

// GetResponse GetResponse
func (tx *GBTx) GetResponse() *sip.Response {
	for {
		res := <-tx.resp
		if res == nil {
			return res
		}
		tx.active <- 2
		//Println("response tx", tx.key, time.Now().Format("2006-01-02 15:04:05"))
		if res.GetStatusCode() == http.StatusContinue || res.GetStatusCode() == http.StatusSwitchingProtocols {
			// Trying and Dialog Establishement 等待下一个返回
			continue
		}
		return res
	}
}

// Close the Close function closes the GBTx
func (tx *GBTx) Close() {
	//Printf("closed tx: %s   %s     TXs: %d", tx.key, time.Now().Format("2006-01-02 15:04:05"), len(ActiveTX.Txs))
	ActiveTX.rmTX(tx)
	close(tx.resp)
	close(tx.active)
}

// ReceiveResponse receive a Response
func (tx *GBTx) ReceiveResponse(msg *sip.Response) {
	defer func() {
		if r := recover(); r != nil {
			//Println("send to closed channel, txkey:", tx.key, "message: \n", msg)
		}
	}()
	//Println("receiveResponse tx", tx.Key(), time.Now().Format("2006-01-02 15:04:05"))
	tx.resp <- msg
	tx.active <- 1
}

// Respond Respond
func (tx *GBTx) Respond(res *sip.Response) error {
	str, _ := sip.Encode(res.Message)
	//Println("send response,to:", (res.DestAdd).String(), "txkey:", tx.key, "message: \n", string(str))
	_, err := tx.conn.WriteTo(str, res.DestAdd)
	return err
}

// Request Request
func (tx *GBTx) Request(req *sip.Request) error {
	str, _ := sip.Encode(req.Message)
	//Println("send Request,to:", (req.DestAdd).String(), "txkey:", tx.key, "message: \n", string(str))
	_, err := tx.conn.WriteTo(str, req.DestAdd)
	return err
}

func GetTXKey(msg *sip.Message) (key string) {
	if len(msg.CallID) > 0 {
		key = msg.CallID
	} else {
		key = utils.RandString(10)
	}
	return
}

func (tx *GBTx) SipResponse() (*sip.Response, error) {
	response := tx.GetResponse()
	if response == nil {
		return nil, utils.NewError(nil, "response timeout", "tx key:", tx.Key())
	}
	if response.GetStatusCode() != http.StatusOK {
		return response, utils.NewError(nil, "response fail", response.GetStatusCode(), response.GetReason(), "tx key:", tx.Key())
	}
	return response, nil
}

func (tx *GBTx) SipRequestForResponse(req *sip.Request) (response *sip.Response, err error) {
	err = tx.Request(req)
	if err == nil {
		return tx.SipResponse()
	}
	return
}
