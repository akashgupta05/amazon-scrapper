package clients

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

var client *http.Client
var once sync.Once

// HTTPClient holds client and host
type HTTPClient struct {
	client *http.Client
	host   string
}

// HTTPClientRequest holds request parameters
type HTTPClientRequest struct {
	Path   string
	Method string
}

// HTTPClientResponse holds response parameters
type HTTPClientResponse struct {
	Response   []byte
	StatusCode int
}

// NewHTTPClient returns instance of HTTPClient
func NewHTTPClient(host string) *HTTPClient {
	return &HTTPClient{
		host:   host,
		client: getHTTPClient(),
	}
}

// Send requests and provide response for given http call
func (hc *HTTPClient) Send(req *HTTPClientRequest) (*HTTPClientResponse, error) {
	request, err := http.NewRequest(req.Method, fmt.Sprintf("%s%s", hc.host, req.Path), nil)
	if err != nil {
		return nil, err
	}

	resp, err := hc.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &HTTPClientResponse{Response: respBytes, StatusCode: resp.StatusCode}, nil
}

func getHTTPClient() *http.Client {
	once.Do(func() {
		tr := &http.Transport{
			MaxIdleConns:       10,
			IdleConnTimeout:    time.Second,
			DisableCompression: true,
		}
		client = &http.Client{Transport: tr, Timeout: time.Second}
	})

	return client
}
