package responses

const (
	OK                           = "SUCCESS"
	DOMAIN_EXCEPTION             = "DOMAIN_EXCEPTION"
	NOT_FOUND_EXCEPTION          = "NOT_FOUND"
	UNEXPECTED_EXCEPTION         = "UNHANDLED_EXCEPTION"
	REQUEST_VALIDATION_EXCEPTION = "REQUEST_VALIDATION_EXCEPTION"
	UNAUTHORIZED_REQUEST         = "UNAUTHORIZED"
	REQUIRES_CREDENTIALS         = "MISSING_CREDENTIALS"
	INVALID_CREDENTIALS          = "INVALID_CREDENTIALS"
	INVALID_ACCESS               = "INVALID_ACCESS"
	CONCURRENCY_EXCEPTION        = "CONCURRENCY_EXCEPTION"
	DOWNSTREAM_EXCEPTION         = "DOWNSTREAM_EXCEPTION"
	UNSUPPORTED_CONTENT          = "UNSUPPORTED_CONTENT"
)

const (
	StatusContinue           = 100 // RFC 9110, 15.2.1
	StatusSwitchingProtocols = 101 // RFC 9110, 15.2.2
	StatusProcessing         = 102 // RFC 2518, 10.1
	StatusEarlyHints         = 103 // RFC 8297

	StatusOK                   = 200 // RFC 9110, 15.3.1
	StatusCreated              = 201 // RFC 9110, 15.3.2
	StatusAccepted             = 202 // RFC 9110, 15.3.3
	StatusNonAuthoritativeInfo = 203 // RFC 9110, 15.3.4
	StatusNoContent            = 204 // RFC 9110, 15.3.5
	StatusResetContent         = 205 // RFC 9110, 15.3.6
	StatusPartialContent       = 206 // RFC 9110, 15.3.7
	StatusMultiStatus          = 207 // RFC 4918, 11.1
	StatusAlreadyReported      = 208 // RFC 5842, 7.1
	StatusIMUsed               = 226 // RFC 3229, 10.4.1

	StatusMultipleChoices   = 300 // RFC 9110, 15.4.1
	StatusMovedPermanently  = 301 // RFC 9110, 15.4.2
	StatusFound             = 302 // RFC 9110, 15.4.3
	StatusSeeOther          = 303 // RFC 9110, 15.4.4
	StatusNotModified       = 304 // RFC 9110, 15.4.5
	StatusUseProxy          = 305 // RFC 9110, 15.4.6
	_                       = 306 // RFC 9110, 15.4.7 (Unused)
	StatusTemporaryRedirect = 307 // RFC 9110, 15.4.8
	StatusPermanentRedirect = 308 // RFC 9110, 15.4.9

	StatusBadRequest                   = 400 // RFC 9110, 15.5.1
	StatusUnauthorized                 = 401 // RFC 9110, 15.5.2
	StatusPaymentRequired              = 402 // RFC 9110, 15.5.3
	StatusForbidden                    = 403 // RFC 9110, 15.5.4
	StatusNotFound                     = 404 // RFC 9110, 15.5.5
	StatusMethodNotAllowed             = 405 // RFC 9110, 15.5.6
	StatusNotAcceptable                = 406 // RFC 9110, 15.5.7
	StatusProxyAuthRequired            = 407 // RFC 9110, 15.5.8
	StatusRequestTimeout               = 408 // RFC 9110, 15.5.9
	StatusConflict                     = 409 // RFC 9110, 15.5.10
	StatusGone                         = 410 // RFC 9110, 15.5.11
	StatusLengthRequired               = 411 // RFC 9110, 15.5.12
	StatusPreconditionFailed           = 412 // RFC 9110, 15.5.13
	StatusRequestEntityTooLarge        = 413 // RFC 9110, 15.5.14
	StatusRequestURITooLong            = 414 // RFC 9110, 15.5.15
	StatusUnsupportedMediaType         = 415 // RFC 9110, 15.5.16
	StatusRequestedRangeNotSatisfiable = 416 // RFC 9110, 15.5.17
	StatusExpectationFailed            = 417 // RFC 9110, 15.5.18
	StatusTeapot                       = 418 // RFC 9110, 15.5.19 (Unused)
	StatusMisdirectedRequest           = 421 // RFC 9110, 15.5.20
	StatusUnprocessableEntity          = 422 // RFC 9110, 15.5.21
	StatusLocked                       = 423 // RFC 4918, 11.3
	StatusFailedDependency             = 424 // RFC 4918, 11.4
	StatusTooEarly                     = 425 // RFC 8470, 5.2.
	StatusUpgradeRequired              = 426 // RFC 9110, 15.5.22
	StatusPreconditionRequired         = 428 // RFC 6585, 3
	StatusTooManyRequests              = 429 // RFC 6585, 4
	StatusRequestHeaderFieldsTooLarge  = 431 // RFC 6585, 5
	StatusUnavailableForLegalReasons   = 451 // RFC 7725, 3

	StatusInternalServerError           = 500 // RFC 9110, 15.6.1
	StatusNotImplemented                = 501 // RFC 9110, 15.6.2
	StatusBadGateway                    = 502 // RFC 9110, 15.6.3
	StatusServiceUnavailable            = 503 // RFC 9110, 15.6.4
	StatusGatewayTimeout                = 504 // RFC 9110, 15.6.5
	StatusHTTPVersionNotSupported       = 505 // RFC 9110, 15.6.6
	StatusVariantAlsoNegotiates         = 506 // RFC 2295, 8.1
	StatusInsufficientStorage           = 507 // RFC 4918, 11.5
	StatusLoopDetected                  = 508 // RFC 5842, 7.2
	StatusNotExtended                   = 510 // RFC 2774, 7
	StatusNetworkAuthenticationRequired = 511 // RFC 6585, 6
)

// StatusText returns a text for the HTTP status code. It returns the empty
// string if the code is unknown.
func StatusBusinessCode(code int) string {
	switch code {
	case StatusContinue:
		return "HTTP:Continue"
	case StatusSwitchingProtocols:
		return "HTTP:Switching Protocols"
	case StatusProcessing:
		return "HTTP:Processing"
	case StatusEarlyHints:
		return "HTTP:Early Hints"
	case StatusOK:
		return "SUCCESS"
	case StatusCreated:
		return "SUCCESS"
	case StatusAccepted:
		return "SUCCESS"
	case StatusNonAuthoritativeInfo:
		return "HTTP:Non-Authoritative Information"
	case StatusNoContent:
		return "HTTP:No Content"
	case StatusResetContent:
		return "HTTP:Reset Content"
	case StatusPartialContent:
		return "HTTP:Partial Content"
	case StatusMultiStatus:
		return "HTTP:Multi-Status"
	case StatusAlreadyReported:
		return "HTTP:Already Reported"
	case StatusIMUsed:
		return "HTTP:IM Used"
	case StatusMultipleChoices:
		return "HTTP:Multiple Choices"
	case StatusMovedPermanently:
		return "HTTP:Moved Permanently"
	case StatusFound:
		return "HTTP:Found"
	case StatusSeeOther:
		return "HTTP:See Other"
	case StatusNotModified:
		return "HTTP:Not Modified"
	case StatusUseProxy:
		return "vUse Proxy"
	case StatusTemporaryRedirect:
		return "HTTP:Temporary Redirect"
	case StatusPermanentRedirect:
		return "HTTP:Permanent Redirect"
	case StatusBadRequest:
		return "REQUEST_VALIDATION_EXCEPTION"
	case StatusUnauthorized:
		return "UNAUTHORIZED"
	case StatusPaymentRequired:
		return "HTTP:Payment Required"
	case StatusForbidden:
		return "INVALID_ACCESS"
	case StatusNotFound:
		return "HTTP:Not Found"
	case StatusMethodNotAllowed:
		return "HTTP:Method Not Allowed"
	case StatusNotAcceptable:
		return "HTTP:Not Acceptable"
	case StatusProxyAuthRequired:
		return "HTTP:Proxy Authentication Required"
	case StatusRequestTimeout:
		return "HTTP:Request Timeout"
	case StatusConflict:
		return "CONCURRENCY_EXCEPTION"
	case StatusGone:
		return "HTTP:Gone"
	case StatusLengthRequired:
		return "HTTP:Length Required"
	case StatusPreconditionFailed:
		return "HTTP:Precondition Failed"
	case StatusRequestEntityTooLarge:
		return "HTTP:Request Entity Too Large"
	case StatusRequestURITooLong:
		return "HTTP:Request URI Too Long"
	case StatusUnsupportedMediaType:
		return "UNSUPPORTED_CONTENT"
	case StatusRequestedRangeNotSatisfiable:
		return "HTTP:Requested Range Not Satisfiable"
	case StatusExpectationFailed:
		return "HTTP:Expectation Failed"
	case StatusTeapot:
		return "HTTP:I'm a teapot"
	case StatusMisdirectedRequest:
		return "HTTP:Misdirected Request"
	case StatusUnprocessableEntity:
		return "DOMAIN_EXCEPTION"
	case StatusLocked:
		return "HTTP:Locked"
	case StatusFailedDependency:
		return "HTTP:Failed Dependency"
	case StatusTooEarly:
		return "HTTP:Too Early"
	case StatusUpgradeRequired:
		return "HTTP:Upgrade Required"
	case StatusPreconditionRequired:
		return "HTTP:Precondition Required"
	case StatusTooManyRequests:
		return "HTTP:Too Many Requests"
	case StatusRequestHeaderFieldsTooLarge:
		return "HTTP:Request Header Fields Too Large"
	case StatusUnavailableForLegalReasons:
		return "HTTP:Unavailable For Legal Reasons"
	case StatusInternalServerError:
		return "UNHANDLED_EXCEPTION"
	case StatusNotImplemented:
		return "HTTP:Not Implemented"
	case StatusBadGateway:
		return "HTTP:Bad Gateway"
	case StatusServiceUnavailable:
		return "DOWNSTREAM_EXCEPTION"
	case StatusGatewayTimeout:
		return "Gateway Timeout"
	case StatusHTTPVersionNotSupported:
		return "HTTP:HTTP Version Not Supported"
	case StatusVariantAlsoNegotiates:
		return "HTTP:Variant Also Negotiates"
	case StatusInsufficientStorage:
		return "HTTP:Insufficient Storage"
	case StatusLoopDetected:
		return "HTTP:Loop Detected"
	case StatusNotExtended:
		return "HTTP:Not Extended"
	case StatusNetworkAuthenticationRequired:
		return "HTTP:Network Authentication Required"
	default:
		return ""
	}
}
