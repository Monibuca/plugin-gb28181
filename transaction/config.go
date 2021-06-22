package transaction

//SIP服务器静态配置信息

/*
# sip监听udp端口
listen              5060;

# SIP server ID(SIP服务器ID).
# 设备端配置编号需要与该值一致，否则无法注册
serial              34020000002000000001;

# SIP server domain(SIP服务器域)
realm               3402000000;

# 服务端发送ack后，接收回应的超时时间，单位为秒
# 如果指定时间没有回应，认为失败
ack_timeout         30;

# 设备心跳维持时间，如果指定时间内(秒）没有接收一个心跳
# 认为设备离线
keepalive_timeout   120;

# 注册之后是否自动给设备端发送invite
# on: 是  off 不是，需要通过api控制
auto_play           on;
# 设备将流发送的端口，是否固定
# on 发送流到多路复用端口 如9000
# off 自动从rtp_mix_port - rtp_max_port 之间的值中
# 选一个可以用的端口
invite_port_fixed     on;

# 向设备或下级域查询设备列表的间隔，单位(秒)
# 默认60秒
query_catalog_interval  60;
*/

type Config struct {
	//sip服务器的配置
	SipNetwork string //传输协议，默认UDP，可选TCP
	SipIP      string //sip 服务器公网IP
	SipPort    uint16 //sip 服务器端口，默认 5060
	Serial     string //sip 服务器 id, 默认 34020000002000000001
	Realm      string //sip 服务器域，默认 3402000000

	AckTimeout        uint16 //sip 服务应答超时，单位秒
	RegisterValidity  int    //注册有效期，单位秒，默认 3600
	RegisterInterval  int    //注册间隔，单位秒，默认 60
	HeartbeatInterval int    //心跳间隔，单位秒，默认 60
	HeartbeatRetry    int    //心跳超时次数，默认 3

	//媒体服务器配置
	MediaIP          string //媒体服务器地址
	MediaPort        uint16 //媒体服务器端口
	MediaPortMin     uint16
	MediaPortMax     uint16
	MediaIdleTimeout uint16 //推流超时时间，超过则断开链接，让设备重连
	AudioEnable      bool   //是否开启音频
	WaitKeyFrame     bool   //是否等待关键帧，如果等待，则在收到第一个关键帧之前，忽略所有媒体流
	Debug            bool   //是否打印调试信息
	CatalogInterval      int    //目录查询间隔
}
