package main

import (
	"fmt"
	"encoding/json"
	"github.com/labstack/echo"
)

/*
// DefaultHTTPErrorHandler invokes the default HTTP error handler.
func (e *Echo) DefaultHTTPErrorHandler(err error, c *Context) {
	e.defaultHTTPErrorHandler(err, c)
}

// SetHTTPErrorHandler registers a custom Echo.HTTPErrorHandler.
func (e *Echo) SetHTTPErrorHandler(h HTTPErrorHandler) {
	e.httpErrorHandler = h
}
 */

func (t *testSuite) TestNotFound() {
	rr := t.Client.Get("http://example.com/not-existed")
	t.Equal(404, rr.Code)
	t.JSONEq(`{
		"error": {
			"id": "11111111111111111111111111111111",
			"message": "Not Found",
			"type": "notFound"
		}}`, rr.Body.String())
}

func (t *testSuite) TestMethodNotAllowed() {
	rr := t.Client.Delete("http://example.com/")
	t.Equal(405, rr.Code)
	t.JSONEq(`{
		"error": {
			"id": "11111111111111111111111111111111",
			"message": "Method Not Allowed",
			"type": "methodNotAllowed"
		}}`, rr.Body.String())
}

func getInternalError(c *echo.Context) error {
	zero := 0
	foo := 1 / zero
	fmt.Print(foo)
	return nil
}

func (t *testSuite) TestInternalError() {
	// TODO move server handler setup into SetupTest method
	t.Client.Server.Get("/error/internal", getInternalError)

	rr := t.Client.Get("http://example.com/error/internal")
	t.Equal(500, rr.Code)
	rv := Envelope{}
	json.Unmarshal(rr.Body.Bytes(), &rv)
	rv.Error.Stack = ""
	bytes, _ := json.Marshal(rv)
	t.JSONEq(`{
		"error": {
			"id": "11111111111111111111111111111111",
			"message": "runtime error: integer divide by zero",
			"type": "internal"
		}}`, string(bytes))
}

func (t *testSuite) TestSampleAppError() {
	// TODO implement
}