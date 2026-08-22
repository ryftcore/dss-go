// Ported from
// dss-jades/src/main/java/eu/europa/esig/dss/jades/JWSCompactSerializationParser.java
// (DSS 6.5.RC1).
package jades

import (
	"bufio"
	"errors"
	"fmt"
	"io"

	"github.com/ryftcore/dss-go/dss/internal/jose"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/utils"
)

// jwsCompactSerializationParserNumberDots defines the maximum number of '.' characters inside a
// JWS signature. Port of NUMBER_DOTS.
const jwsCompactSerializationParserNumberDots = 2

// jwsCompactSerializationParserDotCharacter is the dot character, used as a separator of parts
// within a JWS Compact signature. Port of DOT_CHARACTER.
const jwsCompactSerializationParserDotCharacter = byte('.')

// JWSCompactSerializationParser is used to parse a Compact JWS. Port of the class
// JWSCompactSerializationParser.
type JWSCompactSerializationParser struct {
	// document is the document to be parsed. Port of the private final field of the same name.
	document model.DSSDocument
}

// NewJWSCompactSerializationParser is the constructor to parse a DSSDocument. Port of
// JWSCompactSerializationParser(DSSDocument).
func NewJWSCompactSerializationParser(document model.DSSDocument) *JWSCompactSerializationParser {
	return &JWSCompactSerializationParser{document: document}
}

// NewJWSCompactSerializationParserFromBinaries is the constructor to parse a byte array. Port of
// JWSCompactSerializationParser(byte[]).
func NewJWSCompactSerializationParserFromBinaries(binaries []byte) *JWSCompactSerializationParser {
	return &JWSCompactSerializationParser{document: model.NewInMemoryDocument(binaries)}
}

// Parse parses the provided document and returns the JWS Compact signature it holds. Port of
// parse().
//
// Upstream reads with `new Scanner(stream, UTF-8).nextLine()`, which takes everything up to the
// first line terminator and throws NoSuchElementException on an empty document; bufio.Scanner's
// default split does the same thing, and the empty case becomes an error rather than a panic.
func (p *JWSCompactSerializationParser) Parse() (*JWS, error) {
	stream, err := p.document.OpenStream()
	if err != nil {
		return nil, model.NewDSSErrorMessageCause(
			fmt.Sprintf("Cannot read the document. Reason : %s", err.Error()), err)
	}
	defer utils.CloseQuietly(stream)

	scanner := bufio.NewScanner(stream)
	scanner.Buffer(make([]byte, 0, 64*1024), 1<<24)
	if !scanner.Scan() {
		if scanErr := scanner.Err(); scanErr != nil {
			return nil, model.NewDSSErrorMessageCause(
				fmt.Sprintf("Cannot read the document. Reason : %s", scanErr.Error()), scanErr)
		}
		return nil, model.NewDSSError("No line found")
	}
	compactSerialization := scanner.Text()
	parts := jose.CompactDeserialize(compactSerialization)
	return NewJWSFromCompactSerializationParts(parts)
}

// IsSupported verifies whether the provided file is a Compact JWS supported by the parser. Port
// of isSupported().
//
// The byte-by-byte walk is not a shortcut for "does it parse": it decides, before any decoding,
// that the document is a compact JWS rather than one of the other formats DSS accepts. The three
// accepting branches, in order, are the ones upstream applies - base64url alphabet anywhere, a
// period up to twice, and any url-safe character once the first period has been seen, since an
// RFC 7797 unencoded payload sits between the first and second period and is not base64url.
func (p *JWSCompactSerializationParser) IsSupported() (bool, error) {
	if !DSSJsonUtilsIsAllowedSignatureDocumentType(p.document) {
		return false, nil
	}

	separatorCounter := 0
	ending := false // used to detect and "trim" line breaks in the end of JWS string

	stream, err := p.document.OpenStream()
	if err != nil {
		return false, model.NewDSSErrorMessageCause(
			fmt.Sprintf("Cannot read the document. Reason : %s", err.Error()), err)
	}
	defer utils.CloseQuietly(stream)

	reader := bufio.NewReader(stream)
	for {
		currentByte, readErr := reader.ReadByte()
		if readErr != nil {
			if errors.Is(readErr, io.EOF) {
				break
			}
			return false, model.NewDSSErrorMessageCause(
				fmt.Sprintf("Cannot read the document. Reason : %s", readErr.Error()), readErr)
		}

		switch {
		case spi.DSSUtilsIsLineBreakByte(currentByte):
			ending = true
		case ending:
			return false, nil
		case currentByte == jwsCompactSerializationParserDotCharacter:
			separatorCounter++
			if separatorCounter > jwsCompactSerializationParserNumberDots {
				return false, nil
			}
		case DSSJsonUtilsIsBase64UrlEncodedByte(currentByte):
			// continue
		case separatorCounter == 1 && DSSJsonUtilsIsUrlSafe(currentByte):
			// continue (payload can be not Base64Url encoded)
		default:
			return false, nil
		}
	}

	if separatorCounter != jwsCompactSerializationParserNumberDots {
		return false, nil
	}
	// if ending: "Line break characters found within the JWS Compact Serialization signature
	// document!" - a warning upstream, not a rejection.
	return true, nil
}
