package twitter

import (
	"context"
	"fmt"
	"github.com/tj/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

// RewriteTransport rewrites https requests to http to avoid TLS cert issues
// during testing.
type RewriteTransport struct {
	Transport http.RoundTripper
}

// RoundTrip rewrites the request scheme to http and calls through to the
// composed RoundTripper or if it is nil, to the http.DefaultTransport.
func (t *RewriteTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.URL.Scheme = "http"
	if t.Transport == nil {
		return http.DefaultTransport.RoundTrip(req)
	}
	return t.Transport.RoundTrip(req)
}


func TestTwitter_Post(t *testing.T) {

	mux := http.NewServeMux()
	mux.HandleFunc("/1.1/statuses/update.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprintf(w, `{"id": 12345, "text": "hello world"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	transport := &RewriteTransport{&http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			return url.Parse(server.URL)
		},
	}}
	client := &http.Client{Transport: transport}

	twitter := New(client)

	ctx := context.Background()
	err := twitter.Post(ctx, []byte("hello world"))
	assert.Nil(t, err)

}
