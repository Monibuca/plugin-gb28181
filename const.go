package gb28181

var reasons = map[int]string{
	100: "Trying",
	180: "Ringing",
	181: "Call Is Being Forwarded",
	182: "Queued",
	183: "Session Progress",
	200: "OK",
	202: "Accepted",
	300: "Multiple Choices",
	301: "Moved Permanently",
	302: "Moved Temporarily",
	305: "Use Proxy",
	380: "Alternative Service",
	400: "Bad Request",
	401: "Unauthorized",
	402: "Payment Required",
	403: "Forbidden",
	404: "Not Found",
	405: "Method Not Allowed",
	406: "Not Acceptable",
	407: "Proxy Authentication Required",
	408: "Request Timeout",
	410: "Gone",
	413: "Request Entity Too Large",
	414: "Request-URI Too Long",
	415: "Unsupported Media Type",
	416: "Unsupported URI Scheme",
	420: "Bad Extension",
	421: "Extension Required",
	423: "Interval Too Brief",
	480: "Temporarily Unavailable",
	481: "Call transaction Does Not Exist",
	482: "Loop Detected",
	483: "Too Many Hops",
	484: "Address Incomplete",
	485: "Ambiguous",
	486: "Busy Here",
	487: "Request Terminated",
	488: "Not Acceptable Here",
	489: "Bad Event",
	491: "Request Pending",
	493: "Undecipherable",
	500: "Server Internal Error",
	501: "Not Implemented",
	502: "Bad Gateway",
	503: "Service Unavailable",
	504: "Server Tim",
	505: "Version Not Supported",
	513: "message Too Large",
	600: "Busy Everywhere",
	603: "Decline",
	604: "Does Not Exist Anywhere",
	606: "SESSION NOT ACCEPTABLE",
}

func Explain(statusCode int) string {
	return reasons[statusCode]
}

const (
	INVIDE_MODE_MANUAL = iota
	INVIDE_MODE_AUTO
	INVIDE_MODE_ONSUBSCRIBE
)
