package main

import (
	"github.com/labstack/echo"
	"runtime"
	"fmt"
)

type ErrorInfo struct {
	Id string `json:"id"`
	Message string `json:"message"`
	Type string `json:"type"`
	Stack string `json:"stack,omitempty"`
}

// Recover returns a middleware which recovers from panics anywhere in the chain
// and handles the control to the centralized HTTPErrorHandler.
func Recover() echo.MiddlewareFunc {
	// TODO: Provide better stack trace `https://github.com/go-errors/errors` `https://github.com/docker/libcontainer/tree/master/stacktrace`
	return func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			defer func() {
				if err := recover(); err != nil {
					trace := make([]byte, 1<<16)
					n := runtime.Stack(trace, true)
					rv := Envelope{
						Error: &ErrorInfo{
							Id: "11111111111111111111111111111111",
							Message: fmt.Sprint(err),
							Type: "internal",
							Stack: string(trace[:n]),
						},
					}
					c.JSON(500, rv)
					// TODO 'func (c *Context) Error(err error)' global error handler must render JSON output
					//c.Error(errorValue)
				}
			}()
			return h(c)
		}
	}
}
