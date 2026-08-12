// Ported from dss-utils/.../Utils.java + IUtils.java (DSS 6.5.RC1),
// stream-related methods, matching org.apache.commons.io.IOUtils
// semantics.
//
// Java's InputStream/OutputStream/Closeable map to io.Reader/io.Writer/
// io.Closer per PORTING.md.

package utils

import (
	"bytes"
	"io"
)

// ToByteArray reads the Reader and returns the resulting byte array.
func ToByteArray(r io.Reader) ([]byte, error) {
	return io.ReadAll(r)
}

// CloseQuietly closes the Closer, ignoring any error and tolerating nil.
func CloseQuietly(c io.Closer) {
	if c == nil {
		return
	}
	_ = c.Close()
}

// Copy copies all of r to w.
func Copy(r io.Reader, w io.Writer) error {
	_, err := io.Copy(w, r)
	return err
}

// Write writes content to w in full.
func Write(content []byte, w io.Writer) error {
	_, err := w.Write(content)
	return err
}

// NullWriter opens a Writer that discards all bytes written to it
// (Java: Utils.nullOutputStream()).
func NullWriter() io.Writer {
	return io.Discard
}

// GetInputStreamSize gets the size of the Reader's content, by consuming
// it.
func GetInputStreamSize(r io.Reader) (int64, error) {
	return io.Copy(io.Discard, r)
}

// CompareInputStreams compares the content of two Readers for equality.
func CompareInputStreams(r1, r2 io.Reader) (bool, error) {
	const chunkSize = 8192
	buf1 := make([]byte, chunkSize)
	buf2 := make([]byte, chunkSize)
	for {
		n1, err1 := io.ReadFull(r1, buf1)
		n2, err2 := io.ReadFull(r2, buf2)

		if err1 != nil && err1 != io.EOF && err1 != io.ErrUnexpectedEOF {
			return false, err1
		}
		if err2 != nil && err2 != io.EOF && err2 != io.ErrUnexpectedEOF {
			return false, err2
		}

		if n1 != n2 {
			return false, nil
		}
		if !bytes.Equal(buf1[:n1], buf2[:n2]) {
			return false, nil
		}

		eof1 := err1 == io.EOF || err1 == io.ErrUnexpectedEOF
		eof2 := err2 == io.EOF || err2 == io.ErrUnexpectedEOF
		if eof1 != eof2 {
			return false, nil
		}
		if eof1 && eof2 {
			return true, nil
		}
	}
}

// StartsWithStream checks if the Reader's content starts with prefixArray.
// (Java: Utils.startsWith(InputStream, byte[]))
//
// Mirrors IOUtils.read-then-compare behavior: exactly len(prefixArray)
// bytes are attempted; if the stream is shorter, the unread tail of the
// comparison buffer stays zero-valued, so a short stream compares equal
// only if prefixArray itself ends in that many zero bytes.
func StartsWithStream(r io.Reader, prefixArray []byte) (bool, error) {
	if r == nil || prefixArray == nil {
		return false, nil
	}
	temp := make([]byte, len(prefixArray))
	_, err := io.ReadFull(r, temp)
	if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
		return false, err
	}
	return bytes.Equal(prefixArray, temp), nil
}
