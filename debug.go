package shopify

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"
)

// Debug is a boolean variable that, if set to true, will cause HTTP requests
// and responses to be printed.
var Debug = getDebug()

func getDebug() bool {
	return os.Getenv("GO_SHOPIFY_DEBUG") == "1"
}

// A DebugTransport logs to the error output all the requests and responses
// that go through it.
type DebugTransport struct {
	Transport http.RoundTripper
}

// RoundTrip logs the request and its response.
func (t *DebugTransport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	if data, err := httputil.DumpRequest(req, true); err == nil {
		fmt.Fprintf(os.Stderr, "%s", string(data))
	}

	resp, err = t.Transport.RoundTrip(req)

	if resp != nil {
		if data, err := httputil.DumpResponse(resp, true); err == nil {
			fmt.Fprintf(os.Stderr, "%s", string(data))
		}
	}

	return resp, err
}
