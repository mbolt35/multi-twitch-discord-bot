package httputil

import (
	"bytes"
	"encoding/json"
	"io"
)

const (
	// HttpContentTypeHeader Content Type Request Header Key
	HttpContentTypeHeader string = "Content-Type"

	// HttpClientIdHeader Client Id Request Header Key
	HttpClientIdHeader string = "Client-ID"

	// HttpAcceptHeader Accept Request Header Key
	HttpAcceptHeader string = "Accept"

	// JsonContentType JSON Content-Type
	JsonContentType string = "application/json"
)

// EncodeJson Encodes JSON from the provided interface and escapes html
func EncodeJson(obj interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(obj)

	return buffer.Bytes(), err
}

// DecodeJson Decodes JSON from the reader into the provided interface
func DecodeJson(r io.Reader, obj interface{}) error {
	decoder := json.NewDecoder(r)
	err := decoder.Decode(obj)
	return err
}
