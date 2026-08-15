// Ported from dss-jades/src/main/java/eu/europa/esig/dss/jades/validation/JAdESAttributeIdentifier.java
// (DSS 6.5.RC1).
//
// Java's DataOutputStream#writeChars(String) writes each UTF-16 code unit of the string as a
// 2-byte big-endian value; DataOutputStream#writeInt(int) writes a 4-byte big-endian value. Both
// are reproduced directly rather than through AbstractSignatureIdentifierBuilder.WriteString
// (spi/validation, a different package with no exported helper for this), since the byte layout
// is small enough to inline and this file has no other reason to import that package.
//
// org.jose4j.json.internal.json_simple.JSONValue.toJSONString(Object) is internal/jose's JSON
// function (see internal/jose/writer.go's doc.go mapping and DSSJsonUtils.toBase64Url's own use
// of it for the same non-Map-argument case).
package jades

import (
	"bytes"
	"encoding/binary"
	"unicode/utf16"

	"github.com/utain/esig/dss/internal/jose"
	"github.com/utain/esig/dss/spi/validation/identifier"
)

// JAdESAttributeIdentifier represents an identifier of a JAdES Attribute (or 'etsiU' component).
// Port of the class JAdESAttributeIdentifier, extending identifier.SignatureAttributeIdentifier.
type JAdESAttributeIdentifier struct {
	identifier.SignatureAttributeIdentifier
}

// newJAdESAttributeIdentifier is the port of the package-private JAdESAttributeIdentifier(byte[])
// constructor.
func newJAdESAttributeIdentifier(data []byte) *JAdESAttributeIdentifier {
	return &JAdESAttributeIdentifier{
		SignatureAttributeIdentifier: identifier.NewSignatureAttributeIdentifierBase("JAdESAttributeIdentifier", data),
	}
}

// JAdESAttributeIdentifierBuild builds a JAdES Attribute identifier. Port of the static
// build(String, Object), which delegates to build(headerName, value, null).
func JAdESAttributeIdentifierBuild(headerName string, value any) *JAdESAttributeIdentifier {
	return JAdESAttributeIdentifierBuildWithOrder(headerName, value, nil)
}

// JAdESAttributeIdentifierBuildWithOrder builds the identifier for an 'etsiU' component. Port of
// the static build(String, Object, Integer).
//
// Java wraps the ByteArrayOutputStream/DataOutputStream pair in a try-with-resources and converts
// any IOException into a DSSException; a bytes.Buffer cannot fail to write, so that branch is
// unreachable here and there is accordingly no error return.
func JAdESAttributeIdentifierBuildWithOrder(headerName string, value any, order *int) *JAdESAttributeIdentifier {
	var buf bytes.Buffer
	if headerName != "" {
		jadesAttributeIdentifierWriteChars(&buf, headerName)
	}
	if value != nil {
		jadesAttributeIdentifierWriteChars(&buf, jose.JSON(value))
	}
	if order != nil {
		_ = binary.Write(&buf, binary.BigEndian, int32(*order))
	}
	return newJAdESAttributeIdentifier(buf.Bytes())
}

// jadesAttributeIdentifierWriteChars ports DataOutputStream#writeChars(String): every UTF-16 code
// unit of str, written as a 2-byte big-endian value.
func jadesAttributeIdentifierWriteChars(buf *bytes.Buffer, str string) {
	for _, unit := range utf16.Encode([]rune(str)) {
		_ = binary.Write(buf, binary.BigEndian, unit)
	}
}
