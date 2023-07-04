package gb28181

import (
	"fmt"
	"strconv"
	"time"
)

var (
	// CatalogXML 获取设备列表xml样式
	CatalogXML = `<?xml version="1.0"?><Query>
<CmdType>Catalog</CmdType>
<SN>%d</SN>
<DeviceID>%s</DeviceID>
</Query>
`
	// RecordInfoXML 获取录像文件列表xml样式
	RecordInfoXML = `<?xml version="1.0"?>
<Query>
<CmdType>RecordInfo</CmdType>
<SN>%d</SN>
<DeviceID>%s</DeviceID>
<StartTime>%s</StartTime>
<EndTime>%s</EndTime>
<Secrecy>0</Secrecy>
<Type>all</Type>
</Query>
`
	// DeviceInfoXML 查询设备详情xml样式
	DeviceInfoXML = `<?xml version="1.0"?>
<Query>
<CmdType>DeviceInfo</CmdType>
<SN>%d</SN>
<DeviceID>%s</DeviceID>
</Query>
`
	// DevicePositionXML 订阅设备位置
	DevicePositionXML = `<?xml version="1.0"?>
<Query>
<CmdType>MobilePosition</CmdType>
<SN>%d</SN>
<DeviceID>%s</DeviceID>
<Interval>%d</Interval>
</Query>`
)

func intTotime(t int64) time.Time {
	tstr := strconv.FormatInt(t, 10)
	if len(tstr) == 10 {
		return time.Unix(t, 0)
	}
	if len(tstr) == 13 {
		return time.UnixMilli(t)
	}
	return time.Now()
}

// BuildDeviceInfoXML 获取设备详情指令
func BuildDeviceInfoXML(sn int, id string) string {
	return fmt.Sprintf(DeviceInfoXML, sn, id)
}

// BuildCatalogXML 获取NVR下设备列表指令
func BuildCatalogXML(sn int, id string) string {
	return fmt.Sprintf(CatalogXML, sn, id)
}

// BuildRecordInfoXML 获取录像文件列表指令
func BuildRecordInfoXML(sn int, id string, start, end int64) string {
	return fmt.Sprintf(RecordInfoXML, sn, id, intTotime(start).Format("2006-01-02T15:04:05"), intTotime(end).Format("2006-01-02T15:04:05"))
}

// BuildDevicePositionXML 订阅设备位置
func BuildDevicePositionXML(sn int, id string, interval int) string {
	return fmt.Sprintf(DevicePositionXML, sn, id, interval)
}

// AlarmResponseXML alarm response xml样式
var (
	AlarmResponseXML = `<?xml version="1.0"?>
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
