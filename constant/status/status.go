package status

import (
	"fmt"
)

// Status http status code
type Status int

func (s Status) String() string {
	if msg, ok := status[s]; ok {
		return msg
	} else {
		return fmt.Sprintf("unknown status %d", s)
	}
}

func (s Status) Code() int {
	return int(s)
}

// HTTP status codes were stolen from net/http.
const (
	Continue           Status = 100 // RFC 7231, 6.2.1
	SwitchingProtocols Status = 101 // RFC 7231, 6.2.2
	Processing         Status = 102 // RFC 2518, 10.1

	OK                   Status = 200 // RFC 7231, 6.3.1
	Created              Status = 201 // RFC 7231, 6.3.2
	Accepted             Status = 202 // RFC 7231, 6.3.3
	NonAuthoritativeInfo Status = 203 // RFC 7231, 6.3.4
	NoContent            Status = 204 // RFC 7231, 6.3.5
	ResetContent         Status = 205 // RFC 7231, 6.3.6
	PartialContent       Status = 206 // RFC 7233, 4.1
	MultiStatus          Status = 207 // RFC 4918, 11.1
	AlreadyReported      Status = 208 // RFC 5842, 7.1
	IMUsed               Status = 226 // RFC 3229, 10.4.1

	MultipleChoices   Status = 300 // RFC 7231, 6.4.1
	MovedPermanently  Status = 301 // RFC 7231, 6.4.2
	Found             Status = 302 // RFC 7231, 6.4.3
	SeeOther          Status = 303 // RFC 7231, 6.4.4
	NotModified       Status = 304 // RFC 7232, 4.1
	UseProxy          Status = 305 // RFC 7231, 6.4.5
	_                 Status = 306 // RFC 7231, 6.4.6 (Unused)
	TemporaryRedirect Status = 307 // RFC 7231, 6.4.7
	PermanentRedirect Status = 308 // RFC 7538, 3

	BadRequest                   Status = 400 // RFC 7231, 6.5.1
	Unauthorized                 Status = 401 // RFC 7235, 3.1
	PaymentRequired              Status = 402 // RFC 7231, 6.5.2
	Forbidden                    Status = 403 // RFC 7231, 6.5.3
	NotFound                     Status = 404 // RFC 7231, 6.5.4
	MethodNotAllowed             Status = 405 // RFC 7231, 6.5.5
	NotAcceptable                Status = 406 // RFC 7231, 6.5.6
	ProxyAuthRequired            Status = 407 // RFC 7235, 3.2
	RequestTimeout               Status = 408 // RFC 7231, 6.5.7
	Conflict                     Status = 409 // RFC 7231, 6.5.8
	Gone                         Status = 410 // RFC 7231, 6.5.9
	LengthRequired               Status = 411 // RFC 7231, 6.5.10
	PreconditionFailed           Status = 412 // RFC 7232, 4.2
	RequestEntityTooLarge        Status = 413 // RFC 7231, 6.5.11
	RequestURITooLong            Status = 414 // RFC 7231, 6.5.12
	UnsupportedMediaType         Status = 415 // RFC 7231, 6.5.13
	RequestedRangeNotSatisfiable Status = 416 // RFC 7233, 4.4
	ExpectationFailed            Status = 417 // RFC 7231, 6.5.14
	Teapot                       Status = 418 // RFC 7168, 2.3.3
	UnprocessableEntity          Status = 422 // RFC 4918, 11.2
	Locked                       Status = 423 // RFC 4918, 11.3
	FailedDependency             Status = 424 // RFC 4918, 11.4
	UpgradeRequired              Status = 426 // RFC 7231, 6.5.15
	PreconditionRequired         Status = 428 // RFC 6585, 3
	TooManyRequests              Status = 429 // RFC 6585, 4
	RequestHeaderFieldsTooLarge  Status = 431 // RFC 6585, 5
	UnavailableForLegalReasons   Status = 451 // RFC 7725, 3

	InternalServerError           Status = 500 // RFC 7231, 6.6.1
	NotImplemented                Status = 501 // RFC 7231, 6.6.2
	BadGateway                    Status = 502 // RFC 7231, 6.6.3
	ServiceUnavailable            Status = 503 // RFC 7231, 6.6.4
	GatewayTimeout                Status = 504 // RFC 7231, 6.6.5
	HTTPVersionNotSupported       Status = 505 // RFC 7231, 6.6.6
	VariantAlsoNegotiates         Status = 506 // RFC 2295, 8.1
	InsufficientStorage           Status = 507 // RFC 4918, 11.5
	LoopDetected                  Status = 508 // RFC 5842, 7.2
	NotExtended                   Status = 510 // RFC 2774, 7
	NetworkAuthenticationRequired Status = 511 // RFC 6585, 6
)

var status = map[Status]string{
	Continue:           "Continue",
	SwitchingProtocols: "Switching Protocols",
	Processing:         "Processing",

	OK:                   "OK",
	Created:              "Created",
	Accepted:             "Accepted",
	NonAuthoritativeInfo: "Non-Authoritative Information",
	NoContent:            "No Content",
	ResetContent:         "Reset Content",
	PartialContent:       "Partial Content",
	MultiStatus:          "Multi-Status",
	AlreadyReported:      "Already Reported",
	IMUsed:               "IM Used",

	MultipleChoices:   "Multiple Choices",
	MovedPermanently:  "Moved Permanently",
	Found:             "Found",
	SeeOther:          "See Other",
	NotModified:       "Not Modified",
	UseProxy:          "Use Proxy",
	TemporaryRedirect: "Temporary Redirect",
	PermanentRedirect: "Permanent Redirect",

	BadRequest:                   "Bad Request",
	Unauthorized:                 "Unauthorized",
	PaymentRequired:              "Payment Required",
	Forbidden:                    "Forbidden",
	NotFound:                     "404 Page not found",
	MethodNotAllowed:             "Method Not Allowed",
	NotAcceptable:                "Not Acceptable",
	ProxyAuthRequired:            "Proxy Authentication Required",
	RequestTimeout:               "Request Timeout",
	Conflict:                     "Conflict",
	Gone:                         "Gone",
	LengthRequired:               "Length Required",
	PreconditionFailed:           "Precondition Failed",
	RequestEntityTooLarge:        "Request Entity Too Large",
	RequestURITooLong:            "Request URI Too Long",
	UnsupportedMediaType:         "Unsupported Media Type",
	RequestedRangeNotSatisfiable: "Requested Range Not Satisfiable",
	ExpectationFailed:            "Expectation Failed",
	Teapot:                       "I'm a teapot",
	UnprocessableEntity:          "Unprocessable Entity",
	Locked:                       "Locked",
	FailedDependency:             "Failed Dependency",
	UpgradeRequired:              "Upgrade Required",
	PreconditionRequired:         "Precondition Required",
	TooManyRequests:              "Too Many Requests",
	RequestHeaderFieldsTooLarge:  "Request Header Fields Too Large",
	UnavailableForLegalReasons:   "Unavailable For Legal Reasons",

	InternalServerError:           "Internal Server Error",
	NotImplemented:                "Not Implemented",
	BadGateway:                    "Bad Gateway",
	ServiceUnavailable:            "Service Unavailable",
	GatewayTimeout:                "Gateway Timeout",
	HTTPVersionNotSupported:       "HTTP Version Not Supported",
	VariantAlsoNegotiates:         "Variant Also Negotiates",
	InsufficientStorage:           "Insufficient Storage",
	LoopDetected:                  "Loop Detected",
	NotExtended:                   "Not Extended",
	NetworkAuthenticationRequired: "Network Authentication Required",
}
