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

func TestHandlers_RefreshNodeHealthStatus_StandardCase(t *testing.T) {
	global.State = &global.AppState{HealthyNodes: []string{}}

	mockHealthcheckResponses := map[string]string{
		"http://localhost:1234": "HEALTHY",
		"http://localhost:5678": "UNHEALTHY",
		"http://localhost:9012": "HEALTHY",
		"http://localhost:1111": "UNHEALTHY",
	}

	global.State.AddNode("http://localhost:1234")
	global.State.AddNode("http://localhost:5678")
	global.State.AddNode("http://localhost:9012")
	global.State.UnhealthyNodes = append(global.State.UnhealthyNodes, "http://localhost:1111")

	healthcheck = mockHealthcheck(mockHealthcheckResponses)

	RefreshNodeHealthStatus()
	assert.Len(t, global.State.HealthyNodes, 2)
	assert.Len(t, global.State.UnhealthyNodes, 2)

	assert.Equal(t, global.State.HealthyNodes[0], "http://localhost:1234")
	assert.Equal(t, global.State.HealthyNodes[1], "http://localhost:9012")
}

func TestHandlers_RefreshNodeHealthStatus_AddHealthy_PutAtTheEnd(t *testing.T) {
	global.State = &global.AppState{HealthyNodes: []string{}}

	mockHealthcheckResponses := map[string]string{
		"http://localhost:1234": "HEALTHY",
		"http://localhost:5678": "HEALTHY",
		"http://localhost:9012": "HEALTHY",
		"http://localhost:1111": "UNHEALTHY",
	}

	global.State.AddNode("http://localhost:1234")
	global.State.AddNode("http://localhost:9012")
	global.State.UnhealthyNodes = append(global.State.UnhealthyNodes, "http://localhost:1111", "http://localhost:5678")

	healthcheck = mockHealthcheck(mockHealthcheckResponses)

	RefreshNodeHealthStatus()
	assert.Len(t, global.State.HealthyNodes, 3)
	assert.Len(t, global.State.UnhealthyNodes, 1)

	assert.Equal(t, global.State.HealthyNodes[0], "http://localhost:1234")
	assert.Equal(t, global.State.HealthyNodes[1], "http://localhost:9012")
	assert.Equal(t, global.State.HealthyNodes[2], "http://localhost:5678")
}

func TestHandlers_GetRegisterHandler_StandardCase(t *testing.T) {
	global.State = &global.AppState{HealthyNodes: []string{}}

	mockHealthcheckResponses := map[string]string{
		"http://localhost:1234": "HEALTHY",
		"http://localhost:5678": "HEALTHY",
		"http://localhost:9012": "HEALTHY",
	}

	global.State.AddNode("http://localhost:1234")
	global.State.AddNode("http://localhost:5678")

	healthcheck = mockHealthcheck(mockHealthcheckResponses)

	e := echo.New()

	req := httptest.NewRequest(
		"GET",
		"/register",
		strings.NewReader(""),
	)

	req.RemoteAddr = "localhost:9012"
	req.Header.Set("X-Client-Port", ":9012")

	rec := httptest.NewRecorder()

	err := GetRegisterHandler(e.NewContext(req, rec))
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)

	assert.Len(t, global.State.HealthyNodes, 3)
	assert.Equal(t, global.State.HealthyNodes[0], "http://localhost:1234")
	assert.Equal(t, global.State.HealthyNodes[1], "http://localhost:5678")
	assert.Equal(t, global.State.HealthyNodes[2], "http://localhost:9012")
}

func TestHandlers_GetRegisterHandler_MissingPort(t *testing.T) {
	global.State = &global.AppState{HealthyNodes: []string{}}

	mockHealthcheckResponses := map[string]string{
		"http://localhost:1234": "HEALTHY",
		"http://localhost:5678": "HEALTHY",
		"http://localhost:9012": "HEALTHY",
	}

	global.State.AddNode("http://localhost:1234")
	global.State.AddNode("http://localhost:5678")

	healthcheck = mockHealthcheck(mockHealthcheckResponses)

	e := echo.New()

	req := httptest.NewRequest(
		"GET",
		"/register",
		strings.NewReader(""),
	)

	req.RemoteAddr = "localhost:9012"

	rec := httptest.NewRecorder()

	err := GetRegisterHandler(e.NewContext(req, rec))
	assert.NoError(t, err)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	assert.Len(t, global.State.HealthyNodes, 2)
	assert.Equal(t, global.State.HealthyNodes[0], "http://localhost:1234")
	assert.Equal(t, global.State.HealthyNodes[1], "http://localhost:5678")
}

func TestHandlers_GetRegisterHandler_UnhealthyNewNode(t *testing.T) {
	global.State = &global.AppState{HealthyNodes: []string{}}

	mockHealthcheckResponses := map[string]string{
		"http://localhost:1234": "HEALTHY",
		"http://localhost:5678": "HEALTHY",
		"http://localhost:9012": "UNHEALTHY",
	}

	global.State.AddNode("http://localhost:1234")
	global.State.AddNode("http://localhost:5678")

	healthcheck = mockHealthcheck(mockHealthcheckResponses)

	e := echo.New()

	req := httptest.NewRequest(
		"GET",
		"/register",
		strings.NewReader(""),
	)

	req.RemoteAddr = "localhost:9012"
	req.Header.Set("X-Client-Port", ":9012")

	rec := httptest.NewRecorder()

	err := GetRegisterHandler(e.NewContext(req, rec))
	assert.NoError(t, err)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	assert.Len(t, global.State.HealthyNodes, 2)
	assert.Equal(t, global.State.HealthyNodes[0], "http://localhost:1234")
	assert.Equal(t, global.State.HealthyNodes[1], "http://localhost:5678")
}

func TestHandlers_GetRegisterHandler_NodeAlreadyRegistered(t *testing.T) {
	global.State = &global.AppState{HealthyNodes: []string{}}

	mockHealthcheckResponses := map[string]string{
		"http://localhost:1234": "HEALTHY",
		"http://localhost:5678": "HEALTHY",
		"http://localhost:9012": "HEALTHY",
	}

	global.State.AddNode("http://localhost:1234")
	global.State.AddNode("http://localhost:5678")
	global.State.AddNode("http://localhost:9012")

	healthcheck = mockHealthcheck(mockHealthcheckResponses)

	e := echo.New()

	req := httptest.NewRequest(
		"GET",
		"/register",
		strings.NewReader(""),
	)

	req.RemoteAddr = "localhost:9012"
	req.Header.Set("X-Client-Port", ":9012")

	rec := httptest.NewRecorder()

	err := GetRegisterHandler(e.NewContext(req, rec))
	assert.NoError(t, err)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	assert.Len(t, global.State.HealthyNodes, 3)
	assert.Equal(t, global.State.HealthyNodes[0], "http://localhost:1234")
	assert.Equal(t, global.State.HealthyNodes[1], "http://localhost:5678")
	assert.Equal(t, global.State.HealthyNodes[2], "http://localhost:9012")
}

// Let Healthcheck update the state
func TestHandlers_GetRegisterHandler_NodeAlreadyRegisteredInUnhealthy(t *testing.T) {
	global.State = &global.AppState{HealthyNodes: []string{}}

	mockHealthcheckResponses := map[string]string{
		"http://localhost:1234": "HEALTHY",
		"http://localhost:5678": "HEALTHY",
		"http://localhost:9012": "HEALTHY",
	}

	global.State.AddNode("http://localhost:1234")
	global.State.AddNode("http://localhost:5678")
	global.State.UnhealthyNodes = append(global.State.UnhealthyNodes, "http://localhost:9012")

	healthcheck = mockHealthcheck(mockHealthcheckResponses)

	e := echo.New()

	req := httptest.NewRequest(
		"GET",
		"/register",
		strings.NewReader(""),
	)

	req.RemoteAddr = "localhost:9012"
	req.Header.Set("X-Client-Port", ":9012")

	rec := httptest.NewRecorder()

	err := GetRegisterHandler(e.NewContext(req, rec))
	assert.NoError(t, err)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	assert.Len(t, global.State.HealthyNodes, 2)
	assert.Equal(t, global.State.HealthyNodes[0], "http://localhost:1234")
	assert.Equal(t, global.State.HealthyNodes[1], "http://localhost:5678")
}

func TestHandlers_PostEchoHandler_StandardCase(t *testing.T) {
	global.State = &global.AppState{HealthyNodes: []string{}}

	mockResponses := map[string]string{
		"http://localhost:1234": "Response from node 1",
		"http://localhost:5678": "Response from node 2",
		"http://localhost:9012": "Response from node 3",
	}

	global.State.AddNode("http://localhost:1234")
	global.State.AddNode("http://localhost:5678")
	global.State.AddNode("http://localhost:9012")

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

	assert.Len(t, global.State.HealthyNodes, 3)
	assert.Equal(t, global.State.HealthyNodes[0], "http://localhost:5678")
	assert.Equal(t, global.State.HealthyNodes[1], "http://localhost:9012")
	assert.Equal(t, global.State.HealthyNodes[2], "http://localhost:1234")
}

// Echo handler shouldn't involve with healthcheck, as failures can be expected to detect problems (e.g Blue-Green)
func TestHandlers_PostEchoHandler_FirstNodeFail_FailFirstSucceedLater(t *testing.T) {
	global.State = &global.AppState{HealthyNodes: []string{}}

	mockResponses := map[string]string{
		"http://localhost:1234": "",
		"http://localhost:5678": "Response from node 2",
		"http://localhost:9012": "Response from node 3",
	}

	global.State.AddNode("http://localhost:1234")
	global.State.AddNode("http://localhost:5678")
	global.State.AddNode("http://localhost:9012")

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

	assert.Len(t, global.State.HealthyNodes, 3)

	assert.Equal(t, global.State.HealthyNodes[0], "http://localhost:5678")
	assert.Equal(t, global.State.HealthyNodes[1], "http://localhost:9012")
	assert.Equal(t, global.State.HealthyNodes[2], "http://localhost:1234")

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

	assert.Len(t, global.State.HealthyNodes, 3)

	assert.Equal(t, global.State.HealthyNodes[0], "http://localhost:9012")
	assert.Equal(t, global.State.HealthyNodes[1], "http://localhost:1234")
	assert.Equal(t, global.State.HealthyNodes[2], "http://localhost:5678")
}

func TestHandlers_PostEchoHandler_NoHealthyNode_Fail(t *testing.T) {
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

func mockHealthcheck(mockResponses map[string]string) func(url string) bool {
	return func(url string) bool {
		v, ok := mockResponses[url]
		return (ok && v == "HEALTHY")
	}
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
