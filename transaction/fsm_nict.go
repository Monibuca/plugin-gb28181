package transaction

import (
	// "fmt"
	"time"
)

/*
非invite事物的状态机：

                                   |Request from TU
                                   |send request
               Timer E             V
               send request  +-----------+
                   +---------|           |-------------------+
                   |         |  Trying   |  Timer F          |
                   +-------->|           |  or Transport Err.|
                             +-----------+  inform TU        |
                200-699         |  |                         |
                resp. to TU     |  |1xx                      |
                +---------------+  |resp. to TU              |
                |                  |                         |
                |   Timer E        V       Timer F           |
                |   send req +-----------+ or Transport Err. |
                |  +---------|           | inform TU         |
                |  |         |Proceeding |------------------>|
                |  +-------->|           |-----+             |
                |            +-----------+     |1xx          |
                |              |      ^        |resp to TU   |
                | 200-699      |      +--------+             |
                | resp. to TU  |                             |
                |              |                             |
                |              V                             |
                |            +-----------+                   |
                |            |           |                   |
                |            | Completed |                   |
                |            |           |                   |
                |            +-----------+                   |
                |              ^   |                         |
                |              |   | Timer K                 |
                +--------------+   | -                       |
                                   |                         |
                                   V                         |
             NOTE:           +-----------+                   |
                             |           |                   |
         transitions         | Terminated|<------------------+
         labeled with        |           |
         the event           +-----------+
         over the action
         to take

                 Figure 6: non-INVITE client transaction
*/
func nict_snd_request(t *Transaction, e *EventObj) error {
	msg := e.msg
	//fmt.Println("nict request:", msg.GetMethod())

	t.origRequest = msg
	t.state = NICT_TRYING

	err := t.SipSend(msg)
	if err != nil {
		t.Terminate()
		return err
	}

	//发送出去之后，开启 timer
	if msg.IsReliable() {
		//stop timer E in reliable transport
		//fmt.Println("Reliabel")
	} else {
		//fmt.Println("Not Reliable")
		//发送定时器
		t.timerE = NewSipTimer(T1, T2, func() {
			t.event <- &EventObj{
				evt: TIMEOUT_E,
				tid: t.id,
			}
		})
	}

	//事物定时器
	t.timerF = time.AfterFunc(TimeF, func() {
		t.event <- &EventObj{
			evt: TIMEOUT_F,
			tid: t.id,
		}
	})

	return nil
}

//事物超时
func osip_nict_timeout_f_event(t *Transaction, e *EventObj) error {
	t.Terminate()
	return nil
}

func osip_nict_timeout_e_event(t *Transaction, e *EventObj) error {
	if t.state == NICT_TRYING {
		//reset timer
		t.timerE.Reset(t.timerE.timeout * 2)
	} else {
		//in PROCEEDING STATE, TIMER is always T2
		t.timerE.Reset(T2)
	}

	//resend origin request
	err := t.SipSend(t.origRequest)
	if err != nil {
		t.Terminate()
		return err
	}

	return nil
}

func nict_rcv_1xx(t *Transaction, e *EventObj) error {
	t.lastResponse = e.msg
	t.state = NICT_PROCEEDING

	//重置发送定时器
	t.timerE.Reset(T2)

	return nil
}

func nict_rcv_23456xx(t *Transaction, e *EventObj) error {
	t.lastResponse = e.msg
	t.state = NICT_COMPLETED

	if e.msg.IsReliable() {
		//不设置timerK
	} else {
		t.timerK = time.AfterFunc(T4, func() {
			t.event <- &EventObj{
				evt: TIMEOUT_K,
				tid: t.id,
			}
		})
	}

	return nil
}

func osip_nict_timeout_k_event(t *Transaction, e *EventObj) error {
	t.Terminate()
	return nil
}
