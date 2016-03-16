package main

import (
	"net/http"
	"github.com/labstack/echo"
)

// Handler
func getIndex(c *echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!\n")
}

func init() {
	// TODO Is this a good place to init global vars?
}

type RequestInfo struct {
	Headers http.Header `json:"headers"`
	Method string `json:"method"`
	ContentType string `json:"contentType"`
	URI string `json:"uri"`
	Cookies []*http.Cookie `json:"cookies"`
	RemoteAddr string `json:"remoteAddr"`
	Proto string `json:"protocol"`
}

type Request struct {
	// TODO collapse headers from map[string][]string into map[string]string
	Headers http.Header `json:"headers"`
	Method string `json:"method"`
	ContentType string `json:"contentType"`
	URI string `json:"uri"`
	// TODO consider use map[string]string instead of *http.Cookie
	Cookies []*http.Cookie `json:"cookies"`
	RemoteAddr string `json:"remoteAddr"`
	Proto string `json:"protocol"`
	// TODO consider also return settings, HTTP context variables like requestId
}

type Status struct {
	Request Request `json:"request"`
}

type Envelope struct {
	Result interface{} `json:"result,omitempty"`
	Error *ErrorInfo `json:"error,omitempty"`
}

func getStatus(c *echo.Context) error {
	req := c.Request()
	rv := Envelope{
		Result: Status{
			Request{
				Method: req.Method,
				Headers: req.Header,
				ContentType: req.Header.Get("Content-Type"),
				// Build full uri like http://localhost:10081/api/v1/system/status
				// Consider build full url from req.URL
				URI: req.RequestURI,
				Cookies: req.Cookies(),
				// TODO drop port 127.0.0.1:33939 -> 127.0.0.1
				RemoteAddr: req.RemoteAddr,
				Proto: req.Proto,
			},
		},
		Error: nil,
	}
	return c.JSON(http.StatusOK, rv)
}

/*
// Use high-order function for DI - func EchoHandler(db Database, redis Redis, ...) http.Handler
func EchoHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello")
	})
}
*/

var errorType = map[int]string{
	http.StatusBadRequest:                   "badRequest",
	http.StatusUnauthorized:                 "unauthorized",
	http.StatusPaymentRequired:              "paymentRequired",
	http.StatusForbidden:                    "forbidden",
	http.StatusNotFound:                     "notFound",
	http.StatusMethodNotAllowed:             "methodNotAllowed",
	http.StatusNotAcceptable:                "notAcceptable",
	http.StatusProxyAuthRequired:            "proxyAuthenticationRequired",
	http.StatusRequestTimeout:               "requestTimeout",
	http.StatusConflict:                     "conflict",
	http.StatusGone:                         "gone",
	http.StatusLengthRequired:               "lengthRequired",
	http.StatusPreconditionFailed:           "preconditionFailed",
	http.StatusRequestEntityTooLarge:        "requestEntityTooLarge",
	http.StatusRequestURITooLong:            "requestUriTooLong",
	http.StatusUnsupportedMediaType:         "unsupportedMediaType",
	http.StatusRequestedRangeNotSatisfiable: "requestedRangeNotSatisfiable",
	http.StatusExpectationFailed:            "rxpectationFailed",
	http.StatusTeapot:                       "teapot",

	http.StatusInternalServerError:     "internal",
	http.StatusNotImplemented:          "notImplemented",
	http.StatusBadGateway:              "badGateway",
	http.StatusServiceUnavailable:      "serviceUnavailable",
	http.StatusGatewayTimeout:          "gatewayTimeout",
	http.StatusHTTPVersionNotSupported: "httpVersionNotSupported",

	428: "rreconditionRequired",
	429: "tooManyRequests",
	431: "requestHeaderFieldsTooLarge",
	511: "networkAuthenticationRequired",
}

// TODO map HTTP code to error.type - invalidOperation, notFound,
//   methodNotAllowed, etc
func apiErrorHandler(err error, c *echo.Context) {
	code := http.StatusInternalServerError
	msg := http.StatusText(code)
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code()
		msg = he.Error()
	}
	// TODO display original error only if debug is enabled
	msg = err.Error()

	rv := Envelope{
		Error: &ErrorInfo{
			Id: "11111111111111111111111111111111",
			Message: msg,
			Type: errorType[code],
		},
	}
	// TODO echo's default error handler checks !c.response.committed before writing to the response
	c.JSON(code, rv)
	// TODO add logging. echo's default error handler calls `e.logger.Error(err)`
}

func NewServer() *echo.Echo {
	e := echo.New()
	e.SetHTTPErrorHandler(apiErrorHandler)

	// Middleware
	//e.Use(mw.Logger())
	e.Use(Recover())

	// Routes
	e.Get("/", getIndex)
	e.Get("/api/v1/system/status", getStatus)
	return e
}

func main() {
	e := NewServer()
	e.Run(":4000")
}

/*
Response to curl
----------------
{
    "result": {
        "request": {
            "contentType": "",
            "cookies": [],
            "headers": {
                "Accept": [
                    "*//*"
                ],
                "User-Agent": [
                    "curl/7.38.0"
                ]
            },
            "method": "GET",
            "protocol": "HTTP/1.1",
            "remoteAddr": "127.0.0.1:34589",
            "uri": "/api/v1/system/status"
        }
    }
}

Response to Firefox
-------------------
{
    "result": {
        "request": {
            "contentType": "",
            "cookies": [],
            "headers": {
                "Accept": [
                    "text/html,application/xhtml+xml,application/xml;q=0.9,*//*;q=0.8"
                ],
                "Accept-Encoding": [
                    "deflate, gzip"
                ],
                "Accept-Language": [
                    "en-US,en;q=0.5"
                ],
                "Cache-Control": [
                    "max-age=0"
                ],
                "Connection": [
                    "keep-alive"
                ],
                "User-Agent": [
                    "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:43.0) Gecko/20100101 Firefox/43.0"
                ]
            },
            "method": "GET",
            "protocol": "HTTP/1.1",
            "remoteAddr": "127.0.0.1:34514",
            "uri": "/api/v1/system/status"
        }
    }
}
*/