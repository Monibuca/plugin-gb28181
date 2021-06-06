package transaction

import (
	// "fmt"
	"time"

	"github.com/Monibuca/plugin-gb28181/v3/sip"
)

/*
                               |INVITE from TU
             Timer A fires     |INVITE sent
             Reset A,          V                      Timer B fires
             INVITE sent +-----------+                or Transport Err.
               +---------|           |---------------+inform TU
               |         |  Calling  |               |
               +-------->|           |-------------->|
                         +-----------+ 2xx           |
                            |  |       2xx to TU     |
                            |  |1xx                  |
    300-699 +---------------+  |1xx to TU            |
   ACK sent |                  |                     |
resp. to TU |  1xx             V                     |
            |  1xx to TU  -----------+               |
            |  +---------|           |               |
            |  |         |Proceeding |-------------->|
            |  +-------->|           | 2xx           |
            |            +-----------+ 2xx to TU     |
            |       300-699    |                     |
            |       ACK sent,  |                     |
            |       resp. to TU|                     |
            |                  |                     |      NOTE:
            |  300-699         V                     |
            |  ACK sent  +-----------+Transport Err. |  transitions
            |  +---------|           |Inform TU      |  labeled with
            |  |         | Completed |-------------->|  the event
            |  +-------->|           |               |  over the action
            |            +-----------+               |  to take
            |              ^   |                     |
            |              |   | Timer D fires       |
            +--------------+   | -                   |
                               |                     |
                               V                     |
                         +-----------+               |
                         |           |               |
                         | Terminated|<--------------+
                         |           |
                         +-----------+

                 Figure 5: INVITE client transaction
*/
func ict_snd_invite(t *Transaction, e *EventObj) error {
	msg := e.msg

	t.isReliable = msg.IsReliable()
	t.origRequest = msg
	t.state = ICT_CALLING

	//发送出去之后，开启 timer
	if msg.IsReliable() {
		//stop timer E in reliable transport
		//fmt.Println("Reliabel")
	} else {
		//fmt.Println("Not Reliable")
		//发送定时器，每次加倍，没有上限？
		t.timerA = NewSipTimer(T1, 0, func() {
			t.event <- &EventObj{
				evt: TIMEOUT_A,
				tid: t.id,
			}
		})
	}

	//事物定时器
	t.timerB = time.AfterFunc(TimeB, func() {
		t.event <- &EventObj{
			evt: TIMEOUT_B,
			tid: t.id,
		}
	})

	return nil
}

func osip_ict_timeout_a_event(t *Transaction, e *EventObj) error {
	err := t.SipSend(t.origRequest)
	if err != nil {
		//发送失败
		t.Terminate()
		return err
	}
	t.timerA.Reset(t.timerA.timeout * 2)

	return nil
}

func osip_ict_timeout_b_event(t *Transaction, e *EventObj) error {
	t.Terminate()
	return nil
}

func ict_rcv_1xx(t *Transaction, e *EventObj) error {
	t.lastResponse = e.msg
	t.state = ICT_PROCEEDING
	return nil
}
func ict_rcv_2xx(t *Transaction, e *EventObj) error {
	t.lastResponse = e.msg
	t.Terminate()

	return nil
}
func ict_rcv_3456xx(t *Transaction, e *EventObj) error {
	t.lastResponse = e.msg
	if t.state != ICT_COMPLETED {
		/* not a retransmission */
		/* automatic handling of ack! */
		ack := ict_create_ack(t, e.msg)
		t.ack = ack
		_ = t.SipSend(t.ack)
		t.Terminate()
	}

	/* start timer D (length is set to MAX (64*DEFAULT_T1 or 32000) */
	t.timerD = time.AfterFunc(TimeD, func() {
		t.event <- &EventObj{
			evt: TIMEOUT_D,
			tid: t.id,
		}
	})

	t.state = ICT_COMPLETED

	return nil
}

func ict_create_ack(t *Transaction, resp *sip.Message) *sip.Message {
	return &sip.Message{
		Mode: t.origRequest.Mode,
		Addr: t.origRequest.Addr,
		StartLine: &sip.StartLine{
			Method: sip.ACK,
			Uri:    t.origRequest.StartLine.Uri,
		},
		MaxForwards: t.origRequest.MaxForwards,
		CallID:      t.callID,
		Contact:     t.origRequest.Contact,
		UserAgent:   t.origRequest.UserAgent,
		Via:         t.via,
		From:        t.from,
		To:          t.to,
		CSeq: &sip.CSeq{
			ID:     1,
			Method: sip.ACK,
		},
	}
}

func ict_retransmit_ack(t *Transaction, e *EventObj) error {
	if t.ack == nil {
		/* ??? we should make a new ACK and send it!!! */
		return nil
	}

	err := t.SipSend(t.ack)
	if err != nil {
		return err
	}
	t.state = ICT_COMPLETED
	return nil
}

func osip_ict_timeout_d_event(t *Transaction, e *EventObj) error {
	t.Terminate()
	return nil
}
