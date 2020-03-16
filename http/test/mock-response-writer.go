package test

import "net/http"

// MockResponseWriter is a test friendly implementation of http.ResponseWriter
type MockResponseWriter struct {

	// Header()
	HeaderReturn http.Header

	// Write()
	WriteBytesArg    []byte
	WriteIntReturn   int
	WriteErrorReturn error

	// WriteHeader()
	WriteHeaderStatusCodeArg int
}

// Header mocks 'http.ResponseWriter.Header() http.Header'
func (writer *MockResponseWriter) Header() http.Header {
	return writer.HeaderReturn
}

// Write mocks 'http.ResponseWriter.Write([]byte) (int,error)'
func (writer *MockResponseWriter) Write(bytes []byte) (int, error) {
	writer.WriteBytesArg = bytes
	return writer.WriteIntReturn, writer.WriteErrorReturn
}

// WriteHeader mocks 'http.ResponseWriter.WriteHeader(int)'
func (writer *MockResponseWriter) WriteHeader(statusCode int) {
	writer.WriteHeaderStatusCodeArg = statusCode
}
