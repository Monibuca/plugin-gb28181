module github.com/Monibuca/plugin-gb28181/v3

go 1.13

require (
	github.com/Monibuca/engine/v3 v3.4.5
	github.com/Monibuca/utils/v3 v3.0.5
	github.com/agiledragon/gomonkey/v2 v2.2.0
	github.com/gofrs/uuid v4.2.0+incompatible
	github.com/google/uuid v1.3.0
	github.com/logrusorgru/aurora v2.0.3+incompatible
	github.com/panjjo/gosip v0.0.0-20210428012454-2314e85e3fd8
	github.com/pion/rtp v1.7.4
	github.com/sirupsen/logrus v1.8.1
	golang.org/x/net v0.0.0-20211112202133-69e39bad7dc2
	golang.org/x/text v0.3.7
)

replace github.com/Monibuca/engine/v3 => ../engine

replace github.com/Monibuca/utils/v3 => ../utils
