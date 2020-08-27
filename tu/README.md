
Transaction User(TU)事务用户：在transaction 层之上的协议层。TU包括了UAC core、UAS core,和proxy core。
tu处理业务逻辑，并对事物层进行操作。

#### 类型

SIP服务器典型有以下几类:

a. 注册服务器 -即只管Register消息,这里相当于location也在这里了

b. 重定向服务器 -给ua回一条302后,转给其它的服务器,这样保证全系统统一接入

c. 代理服务器 -只做proxy,即对SIP消息转发

d. 媒体服务器-只做rtp包相关处理,即media server

e. B2BUA - 这个里包实际一般是可以含以上几种服务器类型

暂时仅处理gb28181 相关

#### TU

tu负责根据应用层需求，发起操作。
比如注册到sip服务器、发起会话、取消会话等。