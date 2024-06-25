package gzip

import (
	"compress/gzip"
	"io"
	"net/http"
)

// CompressResponseWriter

func newCompressResponseWriter(w http.ResponseWriter) *compressResponseWriter {
	return &compressResponseWriter{
		rw: w,
		zp: gzip.NewWriter(w),
	}
}

func (c *compressResponseWriter) Header() http.Header {
	return c.rw.Header()
}

func (c *compressResponseWriter) Write(p []byte) (int, error) {
	return c.zp.Write(p)
}

func (c *compressResponseWriter) WriteHeader(statusCode int) {
	if statusCode < 300 {
		c.rw.Header().Set("Content-Encoding", "gzip")
	}
	c.rw.WriteHeader(statusCode)
}

func (c *compressResponseWriter) Close() error {
	return c.zp.Close()
}

// CompressReader

func newCompressReader(r io.ReadCloser) (*compressReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &compressReader{
		r:  r,
		zr: zr,
	}, nil
}

func (c compressReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

func (c *compressReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zr.Close()
}
