package rapidapi_test

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/clambin/gotools/httpstub"
	"github.com/stretchr/testify/assert"

	"github.com/clambin/gotools/rapidapi"
)

func TestClient_Call(t *testing.T) {
	var testCases = []struct {
		name     string
		pass     bool
		hostname string
		apikey   string
		endpoint string
		response string
	}{
		{"happy", true, "example.com", "1234", "/", "OK"},
		{"no apikey", false, "example.com", "", "/", "Forbidden"},
		{"bad endpoint", false, "example.com", "1234", "/invalid", "Page not found"},
	}

	for _, testCase := range testCases {
		client := rapidapi.NewWithHTTPClient(httpstub.NewTestClient(loopback), testCase.hostname, testCase.apikey)

		response, err := client.Call(testCase.endpoint)

		if testCase.pass == true {
			assert.Nil(t, err, testCase.name)
			assert.Equal(t, testCase.response, string(response), testCase.name)
		} else {
			assert.NotNil(t, err, testCase.name)
			assert.Equal(t, testCase.response, err.Error(), testCase.name)
		}

	}
}

func TestClient_CallWithContext(t *testing.T) {
	client := rapidapi.NewWithHTTPClient(httpstub.NewTestClient(loopback), "example.com", "1234")

	ctx, cancel := context.WithCancel(context.Background())

	var err error
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		_, err = client.CallWithContext(ctx, "/slow")
		wg.Done()
	}()
	cancel()

	wg.Wait()
	assert.Error(t, err)
}

// loopback emulates a rapidAPI endpoint
func loopback(req *http.Request) *http.Response {
	if req.Header.Get("x-rapidapi-key") != "1234" {
		return &http.Response{
			StatusCode: 304,
			Status:     "Forbidden",
			Header:     make(http.Header),
			Body:       ioutil.NopCloser(bytes.NewBufferString("")),
		}
	}
	if req.URL.Host != "example.com" || !(req.URL.Path == "/" || req.URL.Path == "/slow") {
		return &http.Response{
			StatusCode: 404,
			Status:     "Page not found",
			Header:     make(http.Header),
			Body:       ioutil.NopCloser(bytes.NewBufferString("")),
		}
	}
	pass := true
	if req.URL.Path == "/slow" {
		select {
		case <-req.Context().Done():
			pass = false
		case <-time.After(60 * time.Second):
		}
	}
	if pass == false {
		return &http.Response{
			StatusCode: http.StatusInternalServerError,
		}
	}
	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     make(http.Header),
		Body:       ioutil.NopCloser(bytes.NewBufferString("OK")),
	}
}
