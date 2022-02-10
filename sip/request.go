package sip

import (
	"fmt"
	"time"
)

// Request Request
type Request struct {
	*Message
}

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
<Type>time</Type>
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
)

// BuildDeviceInfoXML 获取设备详情指令
func BuildDeviceInfoXML(sn int,id string) string {
	return fmt.Sprintf(DeviceInfoXML,sn, id)
}

// BuildCatalogXML 获取NVR下设备列表指令
func BuildCatalogXML(sn int, id string) string {
	return fmt.Sprintf(CatalogXML,sn, id)
}

// BuildRecordInfoXML 获取录像文件列表指令
func BuildRecordInfoXML(sn int,id string, start, end int64) string {
	return fmt.Sprintf(RecordInfoXML, sn,id, time.Unix(start, 0).Format("2006-01-02T15:04:05"), time.Unix(end, 0).Format("2006-01-02T15:04:05"))
}
