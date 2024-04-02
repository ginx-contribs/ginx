package headers

const (
	Date = "Date"

	IfModifiedSince = "If-Modified-Since"
	LastModified    = "Last-Modified"

	// Redirects
	Location = "Location"

	// Transfer coding
	TE               = "TE"
	Trailer          = "Trailer"
	TrailerLower     = "trailer"
	TransferEncoding = "Transfer-Encoding"

	// Controls
	Cookie         = "Cookie"
	Expect         = "Expect"
	MaxForwards    = "Max-Forwards"
	SetCookie      = "Set-Cookie"
	SetCookieLower = "set-cookie"

	// Connection management
	Connection      = "Connection"
	KeepAlive       = "Keep-Alive"
	ProxyConnection = "Proxy-Connection"

	// Authentication
	Authorization      = "Authorization"
	ProxyAuthenticate  = "Proxy-Authenticate"
	ProxyAuthorization = "Proxy-Authorization"
	WWWAuthenticate    = "WWW-Authenticate"

	// Range requests
	AcceptRanges = "Accept-Ranges"
	ContentRange = "Content-Range"
	IfRange      = "If-Range"
	Range        = "Range"

	// Response context
	Allow       = "Allow"
	Server      = "Server"
	ServerLower = "server"

	// Request context
	From           = "From"
	Host           = "Host"
	Referer        = "Referer"
	ReferrerPolicy = "Referrer-Policy"
	UserAgent      = "User-Agent"

	// Message body information
	ContentEncoding = "Content-Encoding"
	ContentLanguage = "Content-Language"
	ContentLength   = "Content-Length"
	ContentLocation = "Content-Location"
	ContentType     = "Content-Type"

	// Content negotiation
	Accept         = "Accept"
	AcceptCharset  = "Accept-Charset"
	AcceptEncoding = "Accept-Encoding"
	AcceptLanguage = "Accept-Language"
	AltSvc         = "Alt-Svc"

	// Cors
	Origin                           = "origin"
	AccessControlRequestMethod       = "Access-Control-Request-Method"
	AccessControlAllowOrigin         = "Access-Control-Allow-Origin"
	AccessControlAllowHeaders        = "Access-Control-Allow-AccessAllowHeaders"
	AccessControlAllowMethods        = "Access-Control-Allow-AccessAllowsMethods"
	AccessControlExposeHeaders       = "Access-Control-Expose-AccessAllowHeaders"
	AccessControlMaxAge              = "Access-Control-Max-Age"
	AccessControlAllowCredentials    = "Access-Control-Allow-Credentials"
	AccessControlAllowPrivateNetwork = "Access-Control-Allow-Private-Network"

	// User custom
	XRequestId = "X-Request-ID"
)
