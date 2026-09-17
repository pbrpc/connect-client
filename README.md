# connect-client

`connect-client` constructs Connect protocol clients from an HTTP client and a
service address.

## Installation

```bash
go get github.com/pbrpc/connect-client
```

## Client

`BaseURL` converts an internal mesh address to its cleartext HTTP URL. `New`
builds the protocol client used by generated Connect service clients:

```go
rpc := connectclient.New(
	httpClient,
	connectclient.BaseURL("service:50051"),
	nil,
)
service := examplev1connect.NewExampleServiceClient(rpc)
```

Client interceptors are supplied in invocation order. Connect transport options
are passed after the interceptors; `connecthttp.WithGRPC()` selects the gRPC
protocol for a peer that serves it exclusively.
