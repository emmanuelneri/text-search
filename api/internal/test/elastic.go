package test

import "net/http"

type MockTransport struct {
	MockResponse *http.Response
}

func (t *MockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return t.MockResponse, nil
}

func NewMockTransport(r *http.Response) *MockTransport {
	return &MockTransport{MockResponse: r}
}
