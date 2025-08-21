package response

import (
	"io"

	"github.com/gabrielluizsf/tcp_to_http/internal/version"
	"github.com/i9si-sistemas/stringx"
)

func NewStatusCode(code int) StatusCode {
	return StatusCode(code).validate()
}

func WriteStatusLine(w io.Writer, statusCode StatusCode) (err error) {
	statusCodeLine := stringx.New(version.HTTP).Concat(stringx.Space, statusCode.validate(), stringx.New("\r\n"))
	_, err = w.Write([]byte(statusCodeLine))
	return
}

type StatusCode int

func (sc StatusCode) String() string {
	switch sc {
	// 1xx Informational
	case StatusContinue:
		return "100 Continue"
	case StatusSwitchingProtocols:
		return "101 Switching Protocols"
	case StatusProcessing:
		return "102 Processing"
	case StatusEarlyHints:
		return "103 Early Hints"

	// 2xx Success
	case StatusOK:
		return "200 OK"
	case StatusCreated:
		return "201 Created"
	case StatusAccepted:
		return "202 Accepted"
	case StatusNonAuthoritativeInfo:
		return "203 Non-Authoritative Information"
	case StatusNoContent:
		return "204 No Content"
	case StatusResetContent:
		return "205 Reset Content"
	case StatusPartialContent:
		return "206 Partial Content"
	case StatusMultiStatus:
		return "207 Multi-Status"
	case StatusAlreadyReported:
		return "208 Already Reported"
	case StatusIMUsed:
		return "226 IM Used"

	// 3xx Redirection
	case StatusMultipleChoices:
		return "300 Multiple Choices"
	case StatusMovedPermanently:
		return "301 Moved Permanently"
	case StatusFound:
		return "302 Found"
	case StatusSeeOther:
		return "303 See Other"
	case StatusNotModified:
		return "304 Not Modified"
	case StatusUseProxy:
		return "305 Use Proxy"
	case StatusTemporaryRedirect:
		return "307 Temporary Redirect"
	case StatusPermanentRedirect:
		return "308 Permanent Redirect"

	// 4xx Client Errors
	case StatusBadRequest:
		return "400 Bad Request"
	case StatusUnauthorized:
		return "401 Unauthorized"
	case StatusPaymentRequired:
		return "402 Payment Required"
	case StatusForbidden:
		return "403 Forbidden"
	case StatusNotFound:
		return "404 Not Found"
	case StatusMethodNotAllowed:
		return "405 Method Not Allowed"
	case StatusNotAcceptable:
		return "406 Not Acceptable"
	case StatusProxyAuthRequired:
		return "407 Proxy Authentication Required"
	case StatusRequestTimeout:
		return "408 Request Timeout"
	case StatusConflict:
		return "409 Conflict"
	case StatusGone:
		return "410 Gone"
	case StatusLengthRequired:
		return "411 Length Required"
	case StatusPreconditionFailed:
		return "412 Precondition Failed"
	case StatusRequestEntityTooLarge:
		return "413 Content Too Large"
	case StatusRequestURITooLong:
		return "414 URI Too Long"
	case StatusUnsupportedMediaType:
		return "415 Unsupported Media Type"
	case StatusRequestedRangeNotSatisfiable:
		return "416 Range Not Satisfiable"
	case StatusExpectationFailed:
		return "417 Expectation Failed"
	case StatusTeapot:
		return "418 I'm a Teapot"
	case StatusMisdirectedRequest:
		return "421 Misdirected Request"
	case StatusUnprocessableEntity:
		return "422 Unprocessable Entity"
	case StatusLocked:
		return "423 Locked"
	case StatusFailedDependency:
		return "424 Failed Dependency"
	case StatusTooEarly:
		return "425 Too Early"
	case StatusUpgradeRequired:
		return "426 Upgrade Required"
	case StatusPreconditionRequired:
		return "428 Precondition Required"
	case StatusTooManyRequests:
		return "429 Too Many Requests"
	case StatusRequestHeaderFieldsTooLarge:
		return "431 Request Header Fields Too Large"
	case StatusUnavailableForLegalReasons:
		return "451 Unavailable For Legal Reasons"

	// 5xx Server Errors
	case StatusInternalServerError:
		return "500 Internal Server Error"
	case StatusNotImplemented:
		return "501 Not Implemented"
	case StatusBadGateway:
		return "502 Bad Gateway"
	case StatusServiceUnavailable:
		return "503 Service Unavailable"
	case StatusGatewayTimeout:
		return "504 Gateway Timeout"
	case StatusHTTPVersionNotSupported:
		return "505 HTTP Version Not Supported"
	case StatusVariantAlsoNegotiates:
		return "506 Variant Also Negotiates"
	case StatusInsufficientStorage:
		return "507 Insufficient Storage"
	case StatusLoopDetected:
		return "508 Loop Detected"
	case StatusNotExtended:
		return "510 Not Extended"
	case StatusNetworkAuthenticationRequired:
		return "511 Network Authentication Required"
	default:
		return "Unknown Status"
	}
}

func (sc StatusCode) validate() StatusCode {
	if sc < 100 || sc > 599 {
		return StatusOK
	}
	return StatusCode(sc)
}

const (
	// 1xx
	StatusContinue           StatusCode = 100
	StatusSwitchingProtocols StatusCode = 101
	StatusProcessing         StatusCode = 102
	StatusEarlyHints         StatusCode = 103

	// 2xx
	StatusOK                   StatusCode = 200
	StatusCreated              StatusCode = 201
	StatusAccepted             StatusCode = 202
	StatusNonAuthoritativeInfo StatusCode = 203
	StatusNoContent            StatusCode = 204
	StatusResetContent         StatusCode = 205
	StatusPartialContent       StatusCode = 206
	StatusMultiStatus          StatusCode = 207
	StatusAlreadyReported      StatusCode = 208
	StatusIMUsed               StatusCode = 226

	// 3xx
	StatusMultipleChoices   StatusCode = 300
	StatusMovedPermanently  StatusCode = 301
	StatusFound             StatusCode = 302
	StatusSeeOther          StatusCode = 303
	StatusNotModified       StatusCode = 304
	StatusUseProxy          StatusCode = 305
	StatusTemporaryRedirect StatusCode = 307
	StatusPermanentRedirect StatusCode = 308

	// 4xx
	StatusBadRequest                   StatusCode = 400
	StatusUnauthorized                 StatusCode = 401
	StatusPaymentRequired              StatusCode = 402
	StatusForbidden                    StatusCode = 403
	StatusNotFound                     StatusCode = 404
	StatusMethodNotAllowed             StatusCode = 405
	StatusNotAcceptable                StatusCode = 406
	StatusProxyAuthRequired            StatusCode = 407
	StatusRequestTimeout               StatusCode = 408
	StatusConflict                     StatusCode = 409
	StatusGone                         StatusCode = 410
	StatusLengthRequired               StatusCode = 411
	StatusPreconditionFailed           StatusCode = 412
	StatusRequestEntityTooLarge        StatusCode = 413
	StatusRequestURITooLong            StatusCode = 414
	StatusUnsupportedMediaType         StatusCode = 415
	StatusRequestedRangeNotSatisfiable StatusCode = 416
	StatusExpectationFailed            StatusCode = 417
	StatusTeapot                       StatusCode = 418
	StatusMisdirectedRequest           StatusCode = 421
	StatusUnprocessableEntity          StatusCode = 422
	StatusLocked                       StatusCode = 423
	StatusFailedDependency             StatusCode = 424
	StatusTooEarly                     StatusCode = 425
	StatusUpgradeRequired              StatusCode = 426
	StatusPreconditionRequired         StatusCode = 428
	StatusTooManyRequests              StatusCode = 429
	StatusRequestHeaderFieldsTooLarge  StatusCode = 431
	StatusUnavailableForLegalReasons   StatusCode = 451

	// 5xx
	StatusInternalServerError           StatusCode = 500
	StatusNotImplemented                StatusCode = 501
	StatusBadGateway                    StatusCode = 502
	StatusServiceUnavailable            StatusCode = 503
	StatusGatewayTimeout                StatusCode = 504
	StatusHTTPVersionNotSupported       StatusCode = 505
	StatusVariantAlsoNegotiates         StatusCode = 506
	StatusInsufficientStorage           StatusCode = 507
	StatusLoopDetected                  StatusCode = 508
	StatusNotExtended                   StatusCode = 510
	StatusNetworkAuthenticationRequired StatusCode = 511
)
