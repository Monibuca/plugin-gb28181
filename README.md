# GB28181插件

该插件提供SIP server的服务，以及流媒体服务器能力，可以将NVR和摄像头的流抓到m7s中，可获取的设备的录像数据以及访问录像视频。也可以控制摄像头的旋转、缩放等。

## 插件地址

github.com/Monibuca/plugin-gb28181

## 插件引入

```go
import (
_ "m7s.live/plugin/gb28181/v4"
)
```

## 默认插件配置

```yaml
gb28181:
  autoinvite:     true
  prefetchrecord: false
  udpcachesize:   0
  sipnetwork:     udp
  sipip:          127.0.0.1
  sipport:        5060
  serial:         "34020000002000000001"
  realm:          "3402000000"
  username:       ""
  password:       ""
  
  acktimeout:        10
  registervalidity:  60
  registerinterval:  60
  heartbeatinterval: 60
  heartbeatretry:    3
  
  mediaip:          127.0.0.1
  mediaport:        58200
  mediaidletimeout: 30
  medianetwork:     udp
  
  removebaninterval: 600
  logverbose:        false
  audioenable:       true
```

- `AutoInvite`     bool 表示自动发起invite，当Server（SIP）接收到设备信息时，立即向设备发送invite命令获取流
- `PreFetchRecord` bool

* sip服务器的配置
- `SipNetwork` string 传输协议，默认UDP，可选TCP
- `SipIP`      string sip 服务器公网IP
- `SipPort`    uint16 sip 服务器端口，默认 5060
- `Serial`     string sip 服务器 id, 默认 34020000002000000001
- `Realm`      string sip 服务器域，默认 3402000000
- `Username`   string sip 服务器账号
- `Password`   string sip 服务器密码

- `AckTimeout`        uint16 sip 服务应答超时，单位秒
- `RegisterValidity`  int    注册有效期，单位秒，默认 3600
- `RegisterInterval`  int    注册间隔，单位秒，默认 60
- `HeartbeatInterval` int    心跳间隔，单位秒，默认 60
- `HeartbeatRetry`    int    心跳超时次数，默认 3

* 媒体服务器配置
- `MediaIP`          string 媒体服务器地址
- `MediaPort`        uint16 媒体服务器端口
- `MediaNetwork`     string 媒体传输协议，默认UDP，可选TCP
- `MediaIdleTimeout` uint16 推流超时时间，超过则断开链接，让设备重连

- `AudioEnable`       bool 是否开启音频
- `LogVerbose`        bool
- `RemoveBanInterval` int  定时移除注册失败的设备黑名单，单位秒，默认10分钟（600秒）
- `UdpCacheSize`      int  表示UDP缓存大小，默认为0，不开启。仅当TCP关闭，切缓存大于0时才开启，会最多缓存最多N个包，并排序，修复乱序造成的无法播放问题，注意开启后，会有一定的性能损耗，并丢失部分包。

**注意某些摄像机没有设置用户名的地方，摄像机会以自身的国标id作为用户名，这个时候m7s会忽略使用摄像机的用户名，忽略配置的用户名**
如果设备配置了错误的用户名和密码，连续三次上报错误后，m7s会记录设备id，并在10分钟内禁止设备注册

## 插件功能

### 使用SIP协议接受NVR或其他GB28181设备的注册

- 服务器启动时自动监听SIP协议端口，当有设备注册时，会记录该设备信息，可以从UI的列表中看到设备
- 定时发送Catalog命令查询设备的目录信息，可获得通道数据或者子设备
- 发送RecordInfo命令查询设备对录像数据
- 发送Invite命令获取设备的实时视频或者录像视频
- 发送PTZ命令来控制摄像头云台
- 自动同步设备位置

### 作为GB28281的流媒体服务器接受设备的媒体流

- 当invite设备的**实时**视频流时，会在m7s中创建对应的流，StreamPath由设备编号和通道编号组成，即[设备编号]/[通道编号],如果有多个层级，通道编号是最后一个层级的编号
- 当invite设备的**录像**视频流时，StreamPath由设备编号和通道编号以及录像的起止时间拼接而成即[设备编号]/[通道编号]/[开始时间]-[结束时间]

### 如何设置UDP缓存大小

通过wireshark抓包，分析rtp，然后看一下大概多少个包可以有序

## 接口API

### 罗列所有的gb28181协议的设备

`/gb28181/api/list`
设备的结构体如下

```go
type Device struct {
	ID              string
	Name            string
	Manufacturer    string
	Model           string
	Owner           string
	RegisterTime    time.Time
	UpdateTime      time.Time
	LastKeepaliveAt time.Time
	Status          string
	Channels        []*Channel
	NetAddr         string
}
```

> 根据golang的规则，小写字母开头的变量不会被序列化

### 从设备拉取视频流

`/gb28181/api/invite`

参数名 | 必传 | 含义
|----|---|---
id|是 | 设备ID
channel|是|通道编号
startTime|否|开始时间（纯数字Unix时间戳）
endTime|否|结束时间（纯数字Unix时间戳）

返回200代表成功

### 停止从设备拉流

`/gb28181/api/bye`

参数名 | 必传 | 含义
|----|---|---
id|是 | 设备ID
channel|是|通道编号

### 发送控制命令

`/gb28181/api/control`

参数名 | 必传 | 含义
|----|---|---
id|是 | 设备ID
channel|是|通道编号
ptzcmd|是|PTZ控制指令

### 查询录像

`/gb28181/api/query/records`

参数名 | 必传 | 含义
|----|---|---
id|是 | 设备ID
channel|是|通道编号
startTime|否|开始时间（字符串，格式：2021-7-23T12:00:00）
endTime|否|结束时间（字符串格式同上）
