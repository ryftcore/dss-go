// Ported from dss-model/.../TimestampBinary.java (DSS 6.5.RC1).
package model

// TimestampBinary contains only a binary representation of a timestamp.
type TimestampBinary struct {
	// bytes are the binaries of a timestamp.
	bytes []byte
}

// NewTimestampBinary creates a TimestampBinary. Ports the default
// constructor.
func NewTimestampBinary(data []byte) *TimestampBinary {
	return &TimestampBinary{bytes: data}
}

// Bytes gets the timestamp's binary.
func (t *TimestampBinary) Bytes() []byte { return t.bytes }
