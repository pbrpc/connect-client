//revive:disable:package-comments
package connectclient

import (
	"connectrpc.com/connect/v2"
	"connectrpc.com/connect/v2/connecthttp"
)

// BaseURL answers with the URL a Connect client is built against for a
// host:port address on the internal mesh: cleartext HTTP, the same posture as
// a gRPC client's insecure credentials.
func BaseURL(address string) string {
	return "http://" + address
}

// New creates a Connect client that dispatches over httpClient against
// baseURL, with the given interceptors and transport options. The generated
// service clients are built from it:
//
//	rpc := connectclient.New(httpClient, connectclient.BaseURL("service:50051"), nil)
//	svc := examplev1connect.NewExampleServiceClient(rpc)
//
// For a peer that speaks only gRPC, pass connecthttp.WithGRPC().
func New(
	httpClient connecthttp.HTTPClient,
	baseURL string,
	interceptors []connect.ClientInterceptor,
	opts ...connecthttp.Option,
) *connect.Client {
	return connect.NewClient(connecthttp.NewTransport(httpClient, baseURL, opts...), interceptors...)
}
