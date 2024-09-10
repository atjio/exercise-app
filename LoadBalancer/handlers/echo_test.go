package handlers

import (
	"LoadBalancer/global"
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestPostEchoHandler_StandardCase(t *testing.T) {
	// Set State
	global.State = &global.AppState{HealthyNodes: []string{}}

	mockResponses := map[string]string{
		"https://localhost:1234": "Response from node 1",
		"https://localhost:5678": "Response from node 2",
		"https://localhost:9012": "Response from node 3",
	}

	global.State.AddNode("https://localhost:1234")
	global.State.AddNode("https://localhost:5678")
	global.State.AddNode("https://localhost:9012")

	passRequest = mockPassRequest(mockResponses)

	e := echo.New()

	req := httptest.NewRequest(
		"POST",
		"/echo",
		strings.NewReader("MOCK"),
	)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	err := PostEchoHandler(e.NewContext(req, rec))
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "Response from node 1", rec.Body.String())
}

func TestPostEchoHandler_FirstNodeFail_FailFirstSucceedLater(t *testing.T) {
	// Set State
	global.State = &global.AppState{HealthyNodes: []string{}}

	mockResponses := map[string]string{
		"https://localhost:1234": "",
		"https://localhost:5678": "Response from node 2",
		"https://localhost:9012": "Response from node 3",
	}

	global.State.AddNode("https://localhost:1234")
	global.State.AddNode("https://localhost:5678")
	global.State.AddNode("https://localhost:9012")

	passRequest = mockPassRequest(mockResponses)

	e := echo.New()

	req := httptest.NewRequest(
		"POST",
		"/echo",
		strings.NewReader("MOCK"),
	)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	err := PostEchoHandler(e.NewContext(req, rec))
	assert.Error(t, err)

	req = httptest.NewRequest(
		"POST",
		"/echo",
		strings.NewReader("MOCK"),
	)
	req.Header.Set("Content-Type", "application/json")

	rec = httptest.NewRecorder()

	err = PostEchoHandler(e.NewContext(req, rec))
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "Response from node 2", rec.Body.String())
}

func TestPostEchoHandler_NoHealthyNode_Fail(t *testing.T) {
	// Set State
	global.State = &global.AppState{HealthyNodes: []string{}}

	mockResponses := map[string]string{}

	passRequest = mockPassRequest(mockResponses)

	e := echo.New()

	req := httptest.NewRequest(
		"POST",
		"/echo",
		strings.NewReader("MOCK"),
	)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	err := PostEchoHandler(e.NewContext(req, rec))
	assert.Error(t, err)
}

func mockPassRequest(mockResponses map[string]string) func(httpClient *http.Client, url string, header string, body *bytes.Reader) (response []byte, err error) {
	return func(httpClient *http.Client, url string, header string, body *bytes.Reader) (response []byte, err error) {
		for k, v := range mockResponses {
			if strings.Contains(url, k) && (v != "") {
				return []byte(v), nil
			}
		}
		return nil, errors.New("not found")
	}
}
