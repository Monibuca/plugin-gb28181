package transaction

import (
	"fmt"
	"time"
)

/*
                         |Request received
                         |pass to TU
                         V
                   +-----------+
                   |           |
                   | Trying    |-------------+
                   |           |             |
                   +-----------+             |200-699 from TU
                         |                   |send response
                         |1xx from TU        |
                         |send response      |
                         |                   |
      Request            V      1xx from TU  |
      send response+-----------+send response|
          +--------|           |--------+    |
          |        | Proceeding|        |    |
          +------->|           |<-------+    |
   +<--------------|           |             |
   |Trnsprt Err    +-----------+             |
   |Inform TU            |                   |
   |                     |                   |
   |                     |200-699 from TU    |
   |                     |send response      |
   |  Request            V                   |
   |  send response+-----------+             |
   |      +--------|           |             |
   |      |        | Completed |<------------+
   |      +------->|           |
   +<--------------|           |
   |Trnsprt Err    +-----------+
   |Inform TU            |
   |                     |Timer J fires
   |                     |-
   |                     |
   |                     V
   |               +-----------+
   |               |           |
   +-------------->| Terminated|
                   |           |
                   +-----------+

       Figure 8: non-INVITE server transaction

*/

func nist_rcv_request(t *Transaction, e *EventObj) error {
	fmt.Println("rcv request: ", e.msg.GetMethod())
	fmt.Println("transaction state: ", t.state.String())
	if t.state != NIST_PRE_TRYING {
		fmt.Println("rcv request retransmission,do response")
		if t.lastResponse != nil {
			err := t.SipSend(t.lastResponse)
			if err != nil {
				//transport error
				return err
			}
		}
		return nil
	} else {
		t.origRequest = e.msg
		t.state = NIST_TRYING
		t.isReliable = e.msg.IsReliable()
	}

	return nil
}

func nist_snd_1xx(t *Transaction, e *EventObj) error {
	t.lastResponse = e.msg
	err := t.SipSend(t.lastResponse)
	if err != nil {
		return err
	}

	t.state = NIST_PROCEEDING
	return nil
}

func nist_snd_23456xx(t *Transaction, e *EventObj) error {
	t.lastResponse = e.msg
	err := t.SipSend(t.lastResponse)
	if err != nil {
		return err
	}
	if t.state != NIST_COMPLETED {
		if !t.isReliable {
			t.timerJ = time.AfterFunc(T1*64, func() {
				t.event <- &EventObj{
					evt: TIMEOUT_J,
					tid: t.id,
				}
			})
		}
	}

	t.state = NIST_COMPLETED
	return nil
}
func osip_nist_timeout_j_event(t *Transaction, e *EventObj) error {
	t.Terminate()
	return nil
}
