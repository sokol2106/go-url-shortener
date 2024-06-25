package gzip

import (
	"compress/gzip"
	"io"
	"net/http"
)

type compressResponseWriter struct {
	rw http.ResponseWriter
	zp *gzip.Writer
}

type compressReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}
