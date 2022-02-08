package sip

import (
	"fmt"
)

// Response Response
type Response struct {
	*Message
}
// AlarmResponseXML alarm response xml样式
var (AlarmResponseXML = `<?xml version="1.0"?>
<Response>
<CmdType>Alarm</CmdType>
<SN>17430</SN>
<DeviceID>%s</DeviceID>
</Response>
`
)

// BuildRecordInfoXML 获取录像文件列表指令
func BuildAlarmResponseXML(id string) string {
	return fmt.Sprintf(AlarmResponseXML, id)
}

