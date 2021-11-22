# GB28181插件

该插件提供SIP server的服务，以及流媒体服务器能力，可以将NVR和摄像头的流抓到m7s中，可获取的设备的录像数据以及访问录像视频。也可以控制摄像头的旋转、缩放等。

## 插件地址

github.com/Monibuca/plugin-gb28181

## 插件引入

```go
import (
_ "github.com/Monibuca/plugin-gb28181"
)
```

## 默认插件配置

```toml
[GB28181]
Serial = "34020000002000000001"
Realm = "3402000000"
Expires = 3600
ListenAddr = "127.0.0.1:5060"
AutoCloseAfter = -1
AutoInvite = false
MediaPort = 58200
CatalogInterval = 30
RemoveBanInterval = 600
Username = ""
Password = ""
UdpCacheSize = 0 
TCP = false
```

- `ListenAddr`是监听的地址，这里需要注意的是必须要带上Server的IP地址，这个IP地址是向设备发送信息的时候需要带上的。
- `Serial` Server（SIP）的编号
- `Realm` Server（SIP）的域
- `AutoCloseAfter` 如果设置大于等于0，则当某个流最后一个订阅者取消订阅时会延迟N秒，会自动发送bye，节省流量。如果为了响应及时，可以设置成-1，保持流的连接
- `AutoInvite` 表示自动发起invite，当Server（SIP）接收到设备信息时，立即向设备发送invite命令获取流
- `MediaPort` 表示用于接收设备流的端口号
- `CatalogInterval` 定时获取设备目录的间隔，单位秒
- `RemoveBanInterval` 定时移除注册失败的设备黑名单，单位秒，默认10分钟（600秒）
- `Username` 国标用户名
- `Password` 国标密码
- `TCP` 是否开启TCP接收国标流，默认false
- `UdpCacheSize` 表示UDP缓存大小，默认为0，不开启。仅当TCP关闭，切缓存大于0时才开启，会最多缓存最多N个包，并排序，修复乱序造成的无法播放问题，注意开启后，会有一定的性能损耗，并丢失部分包。

**注意某些摄像机没有设置用户名的地方，摄像机会以自身的国标id作为用户名，这个时候m7s会忽略使用摄像机的用户名，忽略配置的用户名**
如果设备配置了错误的用户名和密码，连续三次上报错误后，m7s会记录设备id，并在10分钟内禁止设备注册

## 插件功能

### 使用SIP协议接受NVR或其他GB28181设备的注册

- 服务器启动时自动监听SIP协议端口，当有设备注册时，会记录该设备信息，可以从UI的列表中看到设备
- 定时发送Catalog命令查询设备的目录信息，可获得通道数据或者子设备
- 发送RecordInfo命令查询设备对录像数据
- 发送Invite命令获取设备的实时视频或者录像视频
- 发送PTZ命令来控制摄像头云台

### 作为GB28281的流媒体服务器接受设备的媒体流

- 当invite设备的**实时**视频流时，会在m7s中创建对应的流，StreamPath由设备编号和通道编号组成，即[设备编号]/[通道编号],如果有多个层级，通道编号是最后一个层级的编号
- 当invite设备的**录像**视频流时，StreamPath由设备编号和通道编号以及录像的起止时间拼接而成即[设备编号]/[通道编号]/[开始时间]-[结束时间]

### 如何设置UDP缓存大小

通过wireshark抓包，分析rtp，然后看一下大概多少个包可以有序

## 接口API

### 罗列所有的gb28181协议的设备

`/api/gb28181/list`
设备的结构体如下

```go
type Device struct {
*transaction.Core `json:"-"`
ID                string
RegisterTime      time.Time
UpdateTime        time.Time
Status            string
Channels          []*Channel
queryChannel      bool
sn                int
from              *sip.Contact
to                *sip.Contact
Addr              string
SipIP             string //暴露的IP
channelMap        map[string]*Channel
channelMutex      sync.RWMutex
}
```

> 根据golang的规则，小写字母开头的变量不会被序列化

### 从设备拉取视频流

`/api/gb28181/invite`

参数名 | 必传 | 含义
|----|---|---
id|是 | 设备ID
channel|是|通道编号
startTime|否|开始时间（纯数字Unix时间戳）
endTime|否|结束时间（纯数字Unix时间戳）

返回200代表成功

### 停止从设备拉流

`/api/gb28181/bye`

参数名 | 必传 | 含义
|----|---|---
id|是 | 设备ID
channel|是|通道编号

### 发送控制命令

`/api/gb28181/control`

参数名 | 必传 | 含义
|----|---|---
id|是 | 设备ID
channel|是|通道编号
ptzcmd|是|PTZ控制指令

### 查询录像

`/api/gb28181/query/records`

参数名 | 必传 | 含义
|----|---|---
id|是 | 设备ID
channel|是|通道编号
startTime|否|开始时间（字符串，格式：2021-7-23T12:00:00）
endTime|否|结束时间（字符串格式同上）
