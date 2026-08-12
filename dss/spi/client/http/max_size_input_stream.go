// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/client/http/MaxSizeInputStream.java (DSS 6.5.RC1).
package http

import (
	"fmt"
	"io"
)

// MaxSizeInputStream is used to limit the size of fetched data. Read returns
// an error once the data limit has been reached. Inspired by
// org.apache.commons.fileupload.util.LimitedInputStream.
type MaxSizeInputStream struct {
	// wrapped is the wrapped stream.
	wrapped io.Reader

	// maxSize is the maximum InputStream size. Zero means no limit.
	maxSize int

	// url is the requested URL.
	url string

	// count is the running byte counter.
	count int
}

// NewMaxSizeInputStream is the default constructor.
//
// Panics if wrappedStream is nil, mirroring Java's
// Objects.requireNonNull(wrappedStream, ...) programmer-error contract.
func NewMaxSizeInputStream(wrappedStream io.Reader, maxSize int, url string) *MaxSizeInputStream {
	if wrappedStream == nil {
		panic("InputStream shall be provided!")
	}
	return &MaxSizeInputStream{wrapped: wrappedStream, maxSize: maxSize, url: url}
}

// Read reads from the wrapped stream, returning an error once the total
// bytes read exceeds maxSize.
func (s *MaxSizeInputStream) Read(p []byte) (int, error) {
	n, err := s.wrapped.Read(p)
	if n > 0 {
		s.count += n
		if sizeErr := s.checkSize(); sizeErr != nil {
			return n, sizeErr
		}
	}
	return n, err
}

// checkSize returns an error once maxSize (when non-zero) has been exceeded.
func (s *MaxSizeInputStream) checkSize() error {
	if s.maxSize != 0 && s.count > s.maxSize {
		return fmt.Errorf("cannot fetch data limit=%d, url=%s", s.maxSize, s.url)
	}
	return nil
}
