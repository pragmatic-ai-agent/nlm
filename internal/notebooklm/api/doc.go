// Package api provides the high-level NotebookLM client.
//
// Client wraps the low-level batchexecute and gRPC-Web transports to expose
// notebook, source, note, and artifact operations as ordinary Go methods.
// It owns request construction, response decoding into the generated
// protobuf types, and the typed sentinel errors (ErrAuthExpired,
// ErrSourceCapReached, ErrArtifactGenerating, …) that callers match with
// errors.Is. Construct a client with New; the zero value is not usable.
package api
