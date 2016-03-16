package main

import (
	"testing"
	"github.com/stretchr/testify/suite"
	"net/http/httptest"
	"github.com/labstack/echo"
	"net/http"
)

type testSuite struct {
	suite.Suite
	HttpRecorder *httptest.ResponseRecorder
	Client *TestClient
}

func TestRunSuite(t *testing.T) {
	suite.Run(t, new(testSuite))
}

func (t *testSuite) SetupTest() {
	t.HttpRecorder = httptest.NewRecorder()
	t.Client = NewTestClient(t.Suite)
}

// TODO Move test client into separate module
type TestClient struct {
	Server *echo.Echo
	recorder *httptest.ResponseRecorder
	suite suite.Suite
}

func NewTestClient(suite suite.Suite) *TestClient {
	return &TestClient{
		Server: NewServer(),
		suite: suite,
	}
}

/* NOTE better to return ResponseWriter, so it would be possible to write
rr1 = GET
rr2 = POST
rr2 = PUT

assert rr1
assert rr2
assert rr3
*/
func (c *TestClient) Get(urlStr string) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	req, err := http.NewRequest("GET", urlStr, nil)
	c.suite.Nil(err)
	c.Server.ServeHTTP(recorder, req)
	return recorder
}

func (c *TestClient) Delete(urlStr string) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	req, err := http.NewRequest("DELETE", urlStr, nil)
	c.suite.Nil(err)
	c.Server.ServeHTTP(recorder, req)
	return recorder
}
