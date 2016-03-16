package main

import (
	//"net/http"
)

func (t *testSuite) TestSystemStatus() {
	rr := t.Client.Get("http://example.com/api/v1/system/status")
	t.Equal(200, rr.Code)
	t.JSONEq(`{
		"result": {
			"request": {
				"protocol": "HTTP/1.1",
				"method": "GET",
				"uri": "",
				"remoteAddr": "",
				"headers": {},
				"contentType": "",
				"cookies": []
			}
		}}`,
		rr.Body.String())
}
