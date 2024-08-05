package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type compressWriter struct {
	w            http.ResponseWriter
	zw           *gzip.Writer
	compressable bool
}

func newCompressWriter(w http.ResponseWriter) *compressWriter {
	return &compressWriter{
		w:            w,
		zw:           gzip.NewWriter(w),
		compressable: false,
	}
}

// Header returns the response header
func (c *compressWriter) Header() http.Header {
	return c.w.Header()
}

// Write write bytes to response body.
func (c *compressWriter) Write(p []byte) (int, error) {
	return c.writer().Write(p)
}

// WriteHeader write http status code.
func (c *compressWriter) WriteHeader(statusCode int) {

	contentType := c.w.Header().Get("Content-Type")

	if strings.Contains(contentType, "text/html") || strings.Contains(contentType, "application/json") {
		c.w.Header().Set("Content-Encoding", "gzip")
		c.compressable = true
	}

	c.w.WriteHeader(statusCode)
}

// Close закрывает gzip.Writer и досылает все данные из буфера.
func (c *compressWriter) Close() error {
	if c.compressable {
		return c.writer().(io.WriteCloser).Close()
	}
	return nil
}

// writer реализует интерфейс io.Writer
func (c *compressWriter) writer() io.Writer {
	if c.compressable {
		return c.zw
	} else {
		return c.w
	}
}

// compressReader realise interface of io.Reader
type compressReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

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

// Read implements io.Reader
func (c compressReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

// Close закрывает gzip.Reader и досылает все данные из gzip.Reader
func (c *compressReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zr.Close()
}

// GzipMiddleware middleware for decompressing client data
func GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ow := w

		acceptEncoding := r.Header.Get("Accept-Encoding")
		supportsGzip := strings.Contains(acceptEncoding, "gzip")
		if supportsGzip {
			cw := newCompressWriter(w)

			ow = cw

			defer cw.Close()
		}

		contentEncoding := r.Header.Get("Content-Encoding")
		sendsGzip := strings.Contains(contentEncoding, "gzip")
		if sendsGzip {
			cr, err := newCompressReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			r.Body = cr
			defer cr.Close()
		}

		next.ServeHTTP(ow, r)
	})
}
