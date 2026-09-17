//revive:disable:package-comments
package connectclient

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"

	"connectrpc.com/connect/v2"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/pbrpc/testing/mocks/roundtripper"
)

const procedure = "/example.ExampleService/Echo"

// echoSpec describes procedure as a unary RPC over StringValue messages.
var echoSpec = connect.Spec{
	StreamType: connect.StreamTypeUnary,
	Procedure:  procedure,
}

func TestBaseURL(t *testing.T) {
	if got := BaseURL("service:50051"); got != "http://service:50051" {
		t.Errorf("BaseURL = %q, want http://service:50051", got)
	}
}

func TestNew(t *testing.T) {
	wire := roundtripper.Record(func(request *http.Request) (*http.Response, error) {
		body, err := io.ReadAll(request.Body)
		if err != nil {
			return nil, err
		}

		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/proto"}},
			Body:       io.NopCloser(bytes.NewReader(body)),
			Request:    request,
		}, nil
	})
	httpClient := &http.Client{Transport: wire}

	intercepted := false
	observe := func(next connect.ClientFunc) connect.ClientFunc {
		return func(ctx context.Context, spec connect.Spec) (connect.ClientStream, error) {
			intercepted = true

			return next(ctx, spec)
		}
	}

	rpc := New(httpClient, BaseURL("service:50051"), []connect.ClientInterceptor{observe})

	var response wrapperspb.StringValue

	if err := rpc.CallUnary(t.Context(), echoSpec, wrapperspb.String("hello"), &response); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if response.GetValue() != "hello" {
		t.Errorf("response = %q, want hello", response.GetValue())
	}
	if !intercepted {
		t.Error("the interceptor did not run")
	}
	if got := wire.Sent()[0].Request.URL.String(); got != "http://service:50051"+procedure {
		t.Errorf("request URL = %q, want the procedure under the base URL", got)
	}
}
