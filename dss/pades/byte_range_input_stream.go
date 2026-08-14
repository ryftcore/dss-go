// Ported from dss-pades/src/main/java/eu/europa/esig/dss/pades/validation/ByteRangeInputStream.java
// (DSS 6.5.RC1).
//
// Java's FilterInputStream wraps an InputStream; the Go counterpart wraps an io.ReadCloser and
// implements io.ReadCloser itself (see pdf_byte_range_document.go, the sole caller, which needs
// the returned value to satisfy DSSDocument.OpenStream's io.ReadCloser contract). Close()
// delegates to the wrapped stream, matching FilterInputStream's default close() behavior.
package pades

import (
	"io"
)

// ByteRangeInputStream reads an io.ReadCloser according to the given ByteRange.
type ByteRangeInputStream struct {
	// in is the wrapped InputStream.
	in io.ReadCloser

	// byteRange is the ByteRange to be read.
	byteRange *ByteRange

	// position identifies the current position of the InputStream.
	position int
}

var _ io.ReadCloser = (*ByteRangeInputStream)(nil)

// NewByteRangeInputStream creates a ByteRangeInputStream. Port of the constructor
// ByteRangeInputStream(InputStream, ByteRange).
func NewByteRangeInputStream(is io.ReadCloser, byteRange *ByteRange) *ByteRangeInputStream {
	if is == nil {
		panic("InputStream cannot be null!")
	}
	if byteRange == nil {
		panic("ByteRange cannot be null!")
	}
	return &ByteRangeInputStream{in: is, byteRange: byteRange}
}

// Read ports read(byte[], int, int).
func (s *ByteRangeInputStream) Read(b []byte) (int, error) {
	if b == nil {
		panic("Byte array cannot be null!")
	}

	if _, err := s.ensureFirstByteRangePosition(); err != nil {
		return 0, err
	}

	totalRead := 0
	var lastErr error
	for totalRead < len(b) {
		remaining := len(b) - totalRead
		toRead := s.remainingBytesInCurrentPart()
		if toRead <= 0 {
			break
		}
		n := remaining
		if toRead < n {
			n = toRead
		}

		readBytes, err := s.in.Read(b[totalRead : totalRead+n])
		if readBytes < 1 {
			lastErr = err
			break
		}
		totalRead += readBytes
		s.position += readBytes
		if err != nil {
			lastErr = err
			break
		}
	}

	if totalRead < 1 {
		if lastErr != nil {
			return 0, lastErr
		}
		return 0, io.EOF
	}

	return totalRead, nil
}

// Close ports the inherited FilterInputStream#close() behavior by delegating to the wrapped
// stream.
func (s *ByteRangeInputStream) Close() error {
	return s.in.Close()
}

func (s *ByteRangeInputStream) ensureFirstByteRangePosition() (int64, error) {
	return s.skip(0)
}

// skip ports skip(long).
func (s *ByteRangeInputStream) skip(n int64) (int64, error) {
	finalPosition := int64(s.position) + n
	var offset int64
	if finalPosition >= int64(s.byteRange.FirstPartStart())+int64(s.byteRange.FirstPartEnd()) &&
		finalPosition < int64(s.byteRange.SecondPartStart()) {
		offset = int64(s.byteRange.SecondPartStart()) + n - int64(s.byteRange.FirstPartStart()) - int64(s.byteRange.FirstPartEnd())
	} else if s.position < s.byteRange.FirstPartStart() {
		offset = int64(s.byteRange.FirstPartStart()-s.position) + n
	} else {
		offset = n
	}

	skipped, err := io.CopyN(io.Discard, s.in, offset)
	if skipped > offset {
		skipped = offset
	}
	s.position += int(skipped)
	if err != nil && err != io.EOF {
		return skipped, err
	}
	return skipped, nil
}

func (s *ByteRangeInputStream) isPositionWithinRange(position int) bool {
	return s.isPositionWithinFirstPart(position) || s.isPositionWithinSecondPart(position)
}

func (s *ByteRangeInputStream) isPositionWithinFirstPart(position int) bool {
	return position >= s.byteRange.FirstPartStart() && position <= s.byteRange.FirstPartStart()+s.byteRange.FirstPartEnd()
}

func (s *ByteRangeInputStream) isPositionWithinSecondPart(position int) bool {
	return position >= s.byteRange.SecondPartStart() && position <= s.byteRange.SecondPartStart()+s.byteRange.SecondPartEnd()
}

func (s *ByteRangeInputStream) remainingBytesInCurrentPart() int {
	if s.isPositionWithinFirstPart(s.position) {
		return s.byteRange.FirstPartStart() + s.byteRange.FirstPartEnd() - s.position
	} else if s.isPositionWithinSecondPart(s.position) {
		return s.byteRange.SecondPartStart() + s.byteRange.SecondPartEnd() - s.position
	}
	return 0
}
