// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/DSSUtils.java (DSS 6.5.RC1).
//
// This is DSS's general-purpose static utility class; per the established convention in this
// package (see DSSASN1Utils, DSSPKUtils, QcStatementUtils, CertificateExtensionsUtils), every
// exported function is prefixed with the Java class name "DSSUtils" since Go has no per-class
// method namespace and several *Utils classes flatten into this one package.
//
// NOT PORTED (deferred to the CMS phase, which owns the CMSSignedData/CMSSignedDataParser
// types they take): toCMSSignedData(InputStream), toCMSSignedData(DSSDocument),
// toCMSSignedData(byte[]), isTimestampToken(DSSDocument). Mirrors the same deferral already
// made in dss_asn1_utils.go for the CMS-shaped methods of DSSASN1Utils.
//
// DEVIATION: loadCertificate(File/byte[]/InputStream) and loadCertificateFromP7c(...) delegate
// upstream to DSSCertificateTokenSecurityFactory / DSSP7CCertificatesSecurityFactory
// (eu.europa.esig.dss.spi.security, dss-spi/.../spi/security/, not in this manifest), whose
// entire purpose is to retry certificate parsing across every configured JCA security
// provider until one succeeds. Go links crypto/x509 directly - as documented in
// DSSSecurityProviderInitSystemProviders, there is no provider registry to iterate - so these
// functions parse with crypto/x509 directly instead of calling out to that (not yet ported)
// factory type. Flagged for the integrator to reconcile against whichever manifest ports
// spi.security.
package spi

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/asn1"
	"encoding/binary"
	"encoding/pem"
	"errors"
	"fmt"
	"hash"
	"io"
	"math/big"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
	"unicode/utf16"

	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/internal/eccurve"
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/utils"
	"golang.org/x/crypto/ripemd160"
	"golang.org/x/crypto/sha3"
)

// DSSUtilsEmptyByteArray is the empty byte array. Port of EMPTY_BYTE_ARRAY.
var DSSUtilsEmptyByteArray = []byte{}

const (
	// DSSUtilsCarriageReturn represents a carriage return '\r' character. Port of
	// CARRIAGE_RETURN.
	DSSUtilsCarriageReturn byte = '\r'
	// DSSUtilsLineFeed represents a new line '\n' character. Port of LINE_FEED.
	DSSUtilsLineFeed byte = '\n'

	// DSSUtilsRFC3339TimeFormat is the RFC 3339 DateTime format used by default, spelled out
	// as a java.text.SimpleDateFormat pattern for fidelity with upstream call sites; it is
	// translated to a Go time layout by dssUtilsGoLayout. Port of RFC3339_TIME_FORMAT.
	DSSUtilsRFC3339TimeFormat = "yyyy-MM-dd'T'HH:mm:ss'Z'"
	// DSSUtilsISO8601DateFormat formats a date-time as specified in ISO 8601-1. Port of
	// ISO8601_DATE_FORMAT.
	DSSUtilsISO8601DateFormat = "yyyy-MM-dd"

	// DSSUtilsUTF8Encoding is the UTF-8 encoding name string. Port of UTF8_ENCODING. Retained
	// for documentation parity even though Go strings are UTF-8 natively and no call site
	// needs to name the encoding explicitly.
	DSSUtilsUTF8Encoding = "UTF-8"
)

// DSSUtilsUTCTimeZone is the UTC timezone (GMT+0), used by default. Port of UTC_TIMEZONE.
var DSSUtilsUTCTimeZone = time.UTC

// dssUtilsLineBreakChars contains character bytes representing a line break (new line,
// carriage return). Port of LINE_BREAK_CHARS.
var dssUtilsLineBreakChars = []byte{DSSUtilsCarriageReturn, DSSUtilsLineFeed}

// dssUtilsOidNamespacePrefix is the URN OID prefix (RFC 3061). Port of OID_NAMESPACE_PREFIX.
const dssUtilsOidNamespacePrefix = "urn:oid:"

// dssUtilsRFC3986URIPattern is the URI regex defined in RFC 3986 Appendix B. Port of
// RFC3986_URI_PATTERN.
var dssUtilsRFC3986URIPattern = regexp.MustCompile(`^(([^:/?#]+):)?(//([^/?#]*))?([^?#]*)(\?([^#]*))?(#(.*))?`)

var dssUtilsURNOidPattern = regexp.MustCompile(`(?i)^urn:oid:.*$`)
var dssUtilsOidCodePattern = regexp.MustCompile(`^([0-2])((\.0)|(\.[1-9][0-9]*))*$`)
var dssUtilsControlCharsPattern = regexp.MustCompile(`[\x00-\x1F\x7F]+`)
var dssUtilsNonAlphanumericPattern = regexp.MustCompile(`[^\p{L}\p{Nd}]+`)
var dssUtilsInvalidXmlCharsPattern = regexp.MustCompile("[^\x09\x0A\x0D\x20-\U0000D7FF\U0000E000-\U0000FFFD\U00010000-\U0010FFFF]")

// ---------------------------------------------------------------------------------------------
// Date formatting / parsing
// ---------------------------------------------------------------------------------------------

// DSSUtilsFormatDateToRFC formats a date according to RFC 3339, the date aligned to UTC.
// Example: "2019-11-19T17:28:15Z". Port of formatDateToRFC(Date); a zero Time (Java's null
// Date) formats as "N/A".
func DSSUtilsFormatDateToRFC(date time.Time) string {
	return DSSUtilsFormatDateWithCustomFormat(date, DSSUtilsRFC3339TimeFormat)
}

// DSSUtilsIsRFCDate checks silently whether dateTimeString is conformant to the RFC 3339
// date-time pattern "yyyy-MM-dd'T'HH:mm:ss'Z'". Port of isRFCDate(String).
func DSSUtilsIsRFCDate(dateTimeString string) bool {
	if utils.IsStringNotEmpty(dateTimeString) {
		_, err := dssUtilsParseWithJavaPattern(DSSUtilsRFC3339TimeFormat, dateTimeString)
		return err == nil
	}
	return false
}

// dssUtilsParseWithJavaPattern parses value against a java.text.SimpleDateFormat pattern in
// the UTC time zone.
//
// SimpleDateFormat#parse(String) starts at index 0 and stops as soon as the pattern is
// satisfied, IGNORING whatever follows - "2019-11-19T17:28:15Z" parses against "yyyy-MM-dd"
// and answers the date alone. time.Parse instead rejects the leftover as "extra text", so the
// leftover the ParseError reports is trimmed off and the prefix re-parsed, which reproduces
// the Java behaviour exactly. Text that does not satisfy the pattern still fails, as it does
// for the non-lenient SimpleDateFormat upstream configures.
func dssUtilsParseWithJavaPattern(pattern, value string) (time.Time, error) {
	layout := dssUtilsGoLayout(pattern)
	parsed, err := time.ParseInLocation(layout, value, time.UTC)
	if err != nil {
		var parseError *time.ParseError
		if !errors.As(err, &parseError) || !strings.Contains(parseError.Message, "extra text") ||
			parseError.ValueElem == "" || len(parseError.ValueElem) >= len(value) ||
			!strings.HasSuffix(value, parseError.ValueElem) {
			return time.Time{}, err
		}
		parsed, err = time.ParseInLocation(layout, value[:len(value)-len(parseError.ValueElem)], time.UTC)
		if err != nil {
			return time.Time{}, err
		}
	}
	// time.Parse accepts a fractional second right after the seconds field even when the
	// layout does not declare one; SimpleDateFormat does not, so a pattern with no
	// millisecond field must reject it.
	if parsed.Nanosecond() != 0 && !strings.ContainsRune(pattern, 'S') {
		return time.Time{}, fmt.Errorf("unparseable date: %q", value)
	}
	return parsed, nil
}

// DSSUtilsParseRFCDate parses a String date in RFC format, e.g. "2019-11-19T17:28:15Z". Port
// of parseRFCDate(String); returns the zero Time (Java's null) when dateTimeString is empty or
// fails to parse. The slf4j warning upstream logs on a parse failure is dropped (not
// load-bearing, per PORTING.md).
func DSSUtilsParseRFCDate(dateTimeString string) time.Time {
	if utils.IsStringNotEmpty(dateTimeString) {
		if t, err := dssUtilsParseWithJavaPattern(DSSUtilsRFC3339TimeFormat, dateTimeString); err == nil {
			return t
		}
	}
	return time.Time{}
}

// DSSUtilsFormatDateToISO8601 formats a date according to ISO/IEC 8601-1 date pattern
// "yyyy-MM-dd". Example: "2019-11-19". Port of formatDateToISO8601(Date).
func DSSUtilsFormatDateToISO8601(date time.Time) string {
	return DSSUtilsFormatDateWithCustomFormat(date, DSSUtilsISO8601DateFormat)
}

// DSSUtilsIsISO8601Date checks silently whether dateString is conformant to the ISO/IEC 8601-1
// date pattern "yyyy-MM-dd". Port of isISO8601Date(String).
func DSSUtilsIsISO8601Date(dateString string) bool {
	if utils.IsStringNotEmpty(dateString) {
		_, err := dssUtilsParseWithJavaPattern(DSSUtilsISO8601DateFormat, dateString)
		return err == nil
	}
	return false
}

// DSSUtilsParseISO8601Date parses an ISO 8601-1 date String, e.g. "2001-01-01". Port of
// parseISO8601Date(String); returns the zero Time when dateString is empty or fails to parse.
func DSSUtilsParseISO8601Date(dateString string) time.Time {
	if utils.IsStringNotEmpty(dateString) {
		if t, err := dssUtilsParseWithJavaPattern(DSSUtilsISO8601DateFormat, dateString); err == nil {
			return t
		}
	}
	return time.Time{}
}

// DSSUtilsFormatDateWithCustomFormat formats date (with UTC time zone) according to format, a
// java.text.SimpleDateFormat pattern. A zero date (Java's null) formats as "N/A". Port of
// formatDateWithCustomFormat(Date, String).
func DSSUtilsFormatDateWithCustomFormat(date time.Time, format string) string {
	return DSSUtilsFormatDateWithCustomFormatAndLocation(date, format, DSSUtilsUTCTimeZone)
}

// DSSUtilsFormatDateWithCustomFormatAndTimeZoneName formats date according to format (a
// java.text.SimpleDateFormat pattern) and timeZone (an IANA time zone name). Port of
// formatDateWithCustomFormat(Date, String, String).
//
// NOTE: when timeZone is empty, the system default time zone is used, matching upstream.
//
// DEVIATION: java.util.TimeZone#getTimeZone silently falls back to GMT for an unrecognised
// zone id; time.LoadLocation returns an error instead. This reproduces upstream's forgiving
// behaviour by falling back to UTC (GMT) in that case too.
func DSSUtilsFormatDateWithCustomFormatAndTimeZoneName(date time.Time, format, timeZone string) string {
	var location *time.Location
	if utils.IsStringNotEmpty(timeZone) {
		loaded, err := time.LoadLocation(timeZone)
		if err != nil {
			loaded = time.UTC
		}
		location = loaded
	}
	return DSSUtilsFormatDateWithCustomFormatAndLocation(date, format, location)
}

// DSSUtilsFormatDateWithCustomFormatAndLocation formats date according to format (a
// java.text.SimpleDateFormat pattern) and location.
//
// NOTE: when location is nil, the system default time zone is used, matching upstream's null
// TimeZone handling. Port of formatDateWithCustomFormat(Date, String, TimeZone).
func DSSUtilsFormatDateWithCustomFormatAndLocation(date time.Time, format string, location *time.Location) string {
	if date.IsZero() {
		return "N/A"
	}
	if location == nil {
		location = time.Local
	}
	return date.In(location).Format(dssUtilsGoLayout(format))
}

// dssUtilsGoLayout translates a java.text.SimpleDateFormat pattern into a Go reference-time
// layout. Covers the pattern letters DSS's own codebase uses for this method (y M d H h m s a
// E, and '...'-quoted literal text, ” being a literal quote); any other letter is passed
// through unchanged, which produces visibly (not silently) wrong output for a pattern outside
// this coverage.
func dssUtilsGoLayout(pattern string) string {
	runes := []rune(pattern)
	var out strings.Builder
	for i := 0; i < len(runes); {
		c := runes[i]
		switch c {
		case '\'':
			i++
			if i < len(runes) && runes[i] == '\'' {
				out.WriteByte('\'')
				i++
				continue
			}
			for i < len(runes) && runes[i] != '\'' {
				out.WriteRune(runes[i])
				i++
			}
			if i < len(runes) {
				i++ // consume closing quote
			}
		case 'y':
			n := dssUtilsRunLength(runes, i, 'y')
			if n >= 4 {
				out.WriteString("2006")
			} else {
				out.WriteString("06")
			}
			i += n
		case 'M':
			n := dssUtilsRunLength(runes, i, 'M')
			switch {
			case n >= 4:
				out.WriteString("January")
			case n == 3:
				out.WriteString("Jan")
			default:
				out.WriteString("01")
			}
			i += n
		case 'd':
			n := dssUtilsRunLength(runes, i, 'd')
			out.WriteString("02")
			i += n
		case 'H':
			n := dssUtilsRunLength(runes, i, 'H')
			out.WriteString("15")
			i += n
		case 'h':
			n := dssUtilsRunLength(runes, i, 'h')
			out.WriteString("03")
			i += n
		case 'm':
			n := dssUtilsRunLength(runes, i, 'm')
			out.WriteString("04")
			i += n
		case 's':
			n := dssUtilsRunLength(runes, i, 's')
			out.WriteString("05")
			i += n
		case 'S':
			// SimpleDateFormat's millisecond field. Go expresses fractional seconds as a run
			// of '0' attached to the '.' (or ',') that must precede it, so the decimal
			// separator the pattern already emitted carries the run.
			n := dssUtilsRunLength(runes, i, 'S')
			out.WriteString(strings.Repeat("0", n))
			i += n
		case 'a':
			out.WriteString("PM")
			i++
		case 'E':
			n := dssUtilsRunLength(runes, i, 'E')
			if n >= 4 {
				out.WriteString("Monday")
			} else {
				out.WriteString("Mon")
			}
			i += n
		default:
			out.WriteRune(c)
			i++
		}
	}
	return out.String()
}

// dssUtilsRunLength returns the number of consecutive occurrences of c starting at runes[i].
func dssUtilsRunLength(runes []rune, i int, c rune) int {
	n := 0
	for i+n < len(runes) && runes[i+n] == c {
		n++
	}
	return n
}

// ---------------------------------------------------------------------------------------------
// Hex / PEM / DER / certificate loading
// ---------------------------------------------------------------------------------------------

// DSSUtilsToHex converts a byte array into a String representing the hexadecimal values of
// each byte in order, lowercase. Returns "" for a nil value. Port of toHex(byte[]).
//
// NOTE: the upstream javadoc claims the result "is converted to uppercase", but the method
// body just delegates to Utils.toHex(byte[]), which - like utils.ToHex here - is lowercase
// (commons-codec Hex.encodeHexString semantics). The doc comment does not match upstream's own
// code, so this port follows the code, per PORTING.md's behavioural-fidelity mandate.
func DSSUtilsToHex(value []byte) string {
	if value != nil {
		return utils.ToHex(value)
	}
	return ""
}

// DSSUtilsConvertToPEM converts the given certificate into its PEM string. Port of
// convertToPEM(CertificateToken) (and the private convertToPEM(Object) it delegates to,
// inlined here since CertificateToken is the only caller).
//
// DEVIATION: upstream wraps BouncyCastle's PemWriter/JcaMiscPEMGenerator and can throw a
// DSSException; encoding/pem.EncodeToMemory has no failure mode for a well-formed
// pem.Block, so there is no error to propagate.
func DSSUtilsConvertToPEM(cert *model.CertificateToken) string {
	block := &pem.Block{Type: "CERTIFICATE", Bytes: cert.Certificate().Raw}
	return string(pem.EncodeToMemory(block))
}

// DSSUtilsIsStartWithASN1SequenceTag returns true if r starts with an ASN.1 Sequence. Port of
// isStartWithASN1SequenceTag(InputStream).
func DSSUtilsIsStartWithASN1SequenceTag(r io.Reader) (bool, error) {
	document, err := model.NewInMemoryDocumentFromStream(r)
	if err != nil {
		return false, err
	}
	b, err := DSSUtilsReadFirstByte(document)
	if err != nil {
		return false, err
	}
	return DSSASN1UtilsIsASN1SequenceTag(b), nil
}

// DSSUtilsConvertToDER converts a PEM encoded certificate/crl/... to DER encoded. Port of
// convertToDER(String).
func DSSUtilsConvertToDER(pemContent string) ([]byte, error) {
	block, _ := pem.Decode([]byte(pemContent))
	if block == nil {
		return nil, model.NewDSSErrorMessageCause("Unable to convert PEM to DER", errors.New("no PEM block found"))
	}
	return block.Bytes, nil
}

// dssUtilsParseCertificate parses a single certificate from DER or PEM encoded binaries.
func dssUtilsParseCertificate(input []byte) (*model.CertificateToken, error) {
	der := input
	if block, _ := pem.Decode(input); block != nil {
		der = block.Bytes
	}
	cert, err := eccurve.ParseCertificate(der)
	if err != nil {
		return nil, model.NewDSSErrorMessageCause(
			"Unable to load CertificateFactory for the given certificate. All security providers have failed.", err)
	}
	return model.NewCertificateToken(cert)
}

// DSSUtilsLoadCertificate loads a certificate from the given file path. The certificate must
// be DER-encoded and may be supplied in binary or printable (PEM/Base64) encoding. Port of
// loadCertificate(File); see the file-level DEVIATION note regarding the bypassed
// DSSCertificateTokenSecurityFactory.
func DSSUtilsLoadCertificate(filePath string) (*model.CertificateToken, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, model.NewDSSErrorMessageCause(fmt.Sprintf("Unable to find a file '%s'", filePath), err)
	}
	return DSSUtilsLoadCertificateFromBinary(data)
}

// DSSUtilsLoadCertificateFromBinary loads a certificate from the byte array. Port of
// loadCertificate(byte[]).
func DSSUtilsLoadCertificateFromBinary(input []byte) (*model.CertificateToken, error) {
	if input == nil {
		panic("Input binary cannot be null")
	}
	return dssUtilsParseCertificate(input)
}

// DSSUtilsLoadCertificateFromStream loads a certificate from the given InputStream. Port of
// loadCertificate(InputStream).
func DSSUtilsLoadCertificateFromStream(r io.Reader) (*model.CertificateToken, error) {
	if r == nil {
		panic("InputStream cannot be null")
	}
	data, err := utils.ToByteArray(r)
	if err != nil {
		return nil, model.NewDSSErrorMessageCause(fmt.Sprintf("Unable to read InputStream : %s", err.Error()), err)
	}
	return dssUtilsParseCertificate(data)
}

// DSSUtilsLoadCertificateFromP7c loads a collection of certificates from a p7c file. Port of
// loadCertificateFromP7c(File); see the file-level DEVIATION note.
func DSSUtilsLoadCertificateFromP7c(filePath string) ([]*model.CertificateToken, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, model.NewDSSErrorMessageCause(fmt.Sprintf("Unable to find a file '%s'", filePath), err)
	}
	return DSSUtilsLoadCertificateFromP7cBinary(data)
}

// DSSUtilsLoadCertificateFromP7cBinary loads a collection of certificates from a p7c byte
// array. Port of loadCertificateFromP7c(byte[]).
//
// A p7c file is a PKCS#7 SignedData structure carrying only a certificate bag (RFC 2315
// degenerate case, no actual CMS signature); crypto/x509 has no PKCS#7 parser, so the
// SignedData ContentInfo is unwrapped by hand to reach the [0] IMPLICIT SET OF Certificate
// field - see dssUtilsP7cCertificates.
func DSSUtilsLoadCertificateFromP7cBinary(input []byte) ([]*model.CertificateToken, error) {
	certsDER, err := dssUtilsP7cCertificates(input)
	if err != nil {
		return nil, err
	}
	if len(certsDER) == 0 {
		return nil, model.NewDSSError("No certificate found in the InputStream")
	}
	result := make([]*model.CertificateToken, 0, len(certsDER))
	for _, der := range certsDER {
		cert, err := eccurve.ParseCertificate(der)
		if err != nil {
			return nil, model.NewDSSErrorMessageCause(fmt.Sprintf("Failed to load certificate(s) : %s", err.Error()), err)
		}
		token, err := model.NewCertificateToken(cert)
		if err != nil {
			return nil, err
		}
		result = append(result, token)
	}
	return result, nil
}

// DSSUtilsLoadCertificateFromP7cStream loads a collection of certificates from a p7c
// InputStream. Port of loadCertificateFromP7c(InputStream).
func DSSUtilsLoadCertificateFromP7cStream(r io.Reader) ([]*model.CertificateToken, error) {
	data, err := utils.ToByteArray(r)
	if err != nil {
		return nil, model.NewDSSErrorMessageCause(fmt.Sprintf("Unable to read InputStream : %s", err.Error()), err)
	}
	return DSSUtilsLoadCertificateFromP7cBinary(data)
}

// DSSUtilsLoadCertificateFromBase64EncodedString loads a certificate from a base 64 encoded
// String. Port of loadCertificateFromBase64EncodedString(String).
func DSSUtilsLoadCertificateFromBase64EncodedString(base64Encoded string) (*model.CertificateToken, error) {
	return DSSUtilsLoadCertificateFromBinary(utils.FromBase64(base64Encoded))
}

// dssUtilsP7cCertificates extracts the DER-encoded Certificate entries from a (possibly PEM
// encoded) PKCS#7 SignedData ContentInfo, per RFC 2315:
//
//	ContentInfo ::= SEQUENCE {
//	    contentType OBJECT IDENTIFIER,             -- signedData (1.2.840.113549.1.7.2)
//	    content     [0] EXPLICIT ANY OPTIONAL }
//	SignedData ::= SEQUENCE {
//	    version          INTEGER,
//	    digestAlgorithms SET OF DigestAlgorithmIdentifier,
//	    contentInfo      ContentInfo,
//	    certificates     [0] IMPLICIT SET OF Certificate OPTIONAL,
//	    ... }
//
// Only the certificates field is of interest here; every other field is captured into an
// asn1.RawValue purely to advance past it in sequence order.
func dssUtilsP7cCertificates(input []byte) ([][]byte, error) {
	der := input
	if block, _ := pem.Decode(input); block != nil {
		der = block.Bytes
	}
	var outer struct {
		ContentType asn1.ObjectIdentifier
		Content     asn1.RawValue `asn1:"explicit,tag:0"`
	}
	if _, err := asn1.Unmarshal(der, &outer); err != nil {
		return nil, model.NewDSSErrorMessageCause(
			"Unable to load CertificateFactory for the given certificate. All security providers have failed.", err)
	}
	var signedData struct {
		Version          int
		DigestAlgorithms asn1.RawValue
		ContentInfo      asn1.RawValue
		Certificates     asn1.RawValue `asn1:"optional,tag:0"`
	}
	if _, err := asn1.Unmarshal(outer.Content.Bytes, &signedData); err != nil {
		return nil, model.NewDSSErrorMessageCause(
			"Unable to load CertificateFactory for the given certificate. All security providers have failed.", err)
	}
	if len(signedData.Certificates.Bytes) == 0 {
		return nil, nil
	}
	var certs [][]byte
	rest := signedData.Certificates.Bytes
	for len(rest) > 0 {
		var raw asn1.RawValue
		tail, err := asn1.Unmarshal(rest, &raw)
		if err != nil {
			return nil, model.NewDSSErrorMessageCause(
				"Unable to load CertificateFactory for the given certificate. All security providers have failed.", err)
		}
		certs = append(certs, raw.FullBytes)
		rest = tail
	}
	return certs, nil
}

// ---------------------------------------------------------------------------------------------
// Digests
// ---------------------------------------------------------------------------------------------

// DSSUtilsSHA1Digest digests the given string with the SHA1 algorithm and hex-encodes the
// resulting byte array. Port of getSHA1Digest(String).
func DSSUtilsSHA1Digest(stringToDigest string) (string, error) {
	d, err := DSSUtilsDigest(enumerations.DigestAlgorithm_SHA1, []byte(stringToDigest))
	if err != nil {
		return "", err
	}
	return utils.ToHex(d), nil
}

// DSSUtilsIsSHA1Digest checks if the provided str represents a SHA-1 digest. Port of
// isSHA1Digest(String).
func DSSUtilsIsSHA1Digest(str string) bool {
	return utils.IsStringNotBlank(str) && utils.IsHexEncoded(str) && len(str) == 40
}

// DSSUtilsDigest digests data with the given algorithm. Port of digest(DigestAlgorithm,
// byte[]).
//
// Panics with the Java message ("The data cannot be null") when data is nil, matching
// Objects.requireNonNull.
func DSSUtilsDigest(digestAlgorithm enumerations.DigestAlgorithm, data []byte) ([]byte, error) {
	if data == nil {
		panic("The data cannot be null")
	}
	switch digestAlgorithm {
	// The SHAKE output lengths are BouncyCastle's: SHAKEDigest#getDigestSize() answers
	// fixedOutputLength / 4, i.e. TWICE the security strength in bytes - 32 bytes for
	// SHAKE-128 and 64 for SHAKE-256, not 16 and 32.
	case enumerations.DigestAlgorithm_SHAKE128:
		out := make([]byte, 32)
		sha3.ShakeSum128(out, data)
		return out, nil
	case enumerations.DigestAlgorithm_SHAKE256:
		out := make([]byte, 64)
		sha3.ShakeSum256(out, data)
		return out, nil
	default:
		h, err := DSSUtilsMessageDigest(digestAlgorithm)
		if err != nil {
			return nil, err
		}
		h.Write(data)
		return h.Sum(nil), nil
	}
}

// DSSUtilsMessageDigest gets the message digest hash.Hash from the DigestAlgorithm. Port of
// getMessageDigest(DigestAlgorithm), i.e. of DigestAlgorithm#getMessageDigest() /
// MessageDigest.getInstance(javaName) against the BouncyCastle provider.
//
// DEVIATION: matches the availability dssMessageDigestCalculatorMessageDigest already
// establishes in dss_message_digest_calculator.go - MD2 and WHIRLPOOL have no Go
// implementation where BouncyCastle/the JDK provide one; SHAKE-128 and SHAKE-256 are
// unavailable in BOTH ports (BouncyCastle registers no JCA MessageDigest under those names
// either - see DSSUtilsDigest, which reaches them directly instead, exactly as upstream's
// digest(DigestAlgorithm, byte[]) does via SHAKEDigest). SHAKE256-512 IS available, since
// BouncyCastle registers "SHAKE256-512".
func DSSUtilsMessageDigest(digestAlgorithm enumerations.DigestAlgorithm) (hash.Hash, error) {
	switch digestAlgorithm {
	case enumerations.DigestAlgorithm_MD5:
		return md5.New(), nil
	case enumerations.DigestAlgorithm_SHA1:
		return sha1.New(), nil
	case enumerations.DigestAlgorithm_SHA224:
		return sha256.New224(), nil
	case enumerations.DigestAlgorithm_SHA256:
		return sha256.New(), nil
	case enumerations.DigestAlgorithm_SHA384:
		return sha512.New384(), nil
	case enumerations.DigestAlgorithm_SHA512:
		return sha512.New(), nil
	case enumerations.DigestAlgorithm_SHA3_224:
		return sha3.New224(), nil
	case enumerations.DigestAlgorithm_SHA3_256:
		return sha3.New256(), nil
	case enumerations.DigestAlgorithm_SHA3_384:
		return sha3.New384(), nil
	case enumerations.DigestAlgorithm_SHA3_512:
		return sha3.New512(), nil
	case enumerations.DigestAlgorithm_RIPEMD160:
		return ripemd160.New(), nil
	case enumerations.DigestAlgorithm_SHAKE256_512:
		return &dssUtilsShake{shake: sha3.NewShake256(), size: 64}, nil
	}
	return nil, model.NewDSSErrorMessageCause(
		fmt.Sprintf("Unable to create a MessageDigest for algorithm '%s'", string(digestAlgorithm)),
		fmt.Errorf("NoSuchAlgorithmException: %s", digestAlgorithm.JavaName()))
}

// dssUtilsShake adapts a SHAKE extendable-output function to hash.Hash by fixing the output
// length, the way BouncyCastle's SHAKEDigest exposes a fixed digest size.
type dssUtilsShake struct {
	shake sha3.ShakeHash
	size  int
}

// Write absorbs more input.
func (s *dssUtilsShake) Write(p []byte) (int, error) { return s.shake.Write(p) }

// Sum appends the fixed-length squeezed output to b, leaving this hash's state untouched.
func (s *dssUtilsShake) Sum(b []byte) []byte {
	clone := s.shake.Clone()
	out := make([]byte, s.size)
	_, _ = io.ReadFull(clone, out)
	return append(b, out...)
}

// Reset discards the absorbed input.
func (s *dssUtilsShake) Reset() { s.shake.Reset() }

// Size returns the fixed output length in bytes.
func (s *dssUtilsShake) Size() int { return s.size }

// BlockSize returns the sponge rate.
func (s *dssUtilsShake) BlockSize() int { return s.shake.BlockSize() }

// compile-time interface assertion.
var _ hash.Hash = (*dssUtilsShake)(nil)

// DSSUtilsToDigestDocument creates a model.DigestDocument with the provided model.Digest. Port
// of toDigestDocument(Digest).
func DSSUtilsToDigestDocument(digest model.Digest) *model.DigestDocument {
	return DSSUtilsToDigestDocumentFromValue(digest.Algorithm(), digest.Value())
}

// DSSUtilsToDigestDocumentFromValue creates a model.DigestDocument with the provided
// DigestAlgorithm and digestValue. Port of toDigestDocument(DigestAlgorithm, byte[]).
func DSSUtilsToDigestDocumentFromValue(digestAlgorithm enumerations.DigestAlgorithm, digestValue []byte) *model.DigestDocument {
	return model.NewDigestDocumentFromBase64(digestAlgorithm, utils.ToBase64(digestValue))
}

// DSSUtilsDigestReader digests the data in r with the given algorithm. Port of
// digest(DigestAlgorithm, InputStream).
func DSSUtilsDigestReader(digestAlgo enumerations.DigestAlgorithm, r io.Reader) ([]byte, error) {
	h, err := DSSUtilsMessageDigest(digestAlgo)
	if err != nil {
		return nil, err
	}
	buffer := make([]byte, 4096)
	for {
		n, readErr := r.Read(buffer)
		if n > 0 {
			h.Write(buffer[:n])
		}
		if readErr == io.EOF {
			break
		}
		if readErr != nil {
			return nil, model.NewDSSErrorMessageCause(fmt.Sprintf("Unable to compute digest : %s", readErr.Error()), readErr)
		}
	}
	return h.Sum(nil), nil
}

// DSSUtilsDigestOfDocument computes the digest for document. Port of digest(DigestAlgorithm,
// DSSDocument).
func DSSUtilsDigestOfDocument(digestAlgorithm enumerations.DigestAlgorithm, document model.DSSDocument) ([]byte, error) {
	return document.DigestValue(digestAlgorithm)
}

// DSSUtilsDigestConcat computes the digest over the concatenation of data. Port of
// digest(DigestAlgorithm, byte[]...).
func DSSUtilsDigestConcat(digestAlgorithm enumerations.DigestAlgorithm, data ...[]byte) ([]byte, error) {
	h, err := DSSUtilsMessageDigest(digestAlgorithm)
	if err != nil {
		return nil, err
	}
	for _, b := range data {
		h.Write(b)
	}
	return h.Sum(nil), nil
}

// DSSUtilsGetDigest returns the model.Digest of dssDocument. Port of getDigest(DigestAlgorithm,
// DSSDocument).
//
// NOTE: named with the "Get" prefix, against the usual drop-get convention, solely to avoid
// colliding with DSSUtilsDigest (the port of the unrelated overload digest(DigestAlgorithm,
// byte[])); Java disambiguates the two by parameter type the way Go cannot.
func DSSUtilsGetDigest(digestAlgo enumerations.DigestAlgorithm, dssDocument model.DSSDocument) (model.Digest, error) {
	value, err := DSSUtilsDigestOfDocument(digestAlgo, dssDocument)
	if err != nil {
		return model.Digest{}, err
	}
	return model.NewDigest(digestAlgo, value), nil
}

// DSSUtilsMD5Digest returns the hex encoding of the MD5 digest of bytes. Port of
// getMD5Digest(byte[]).
func DSSUtilsMD5Digest(data []byte) (string, error) {
	d, err := DSSUtilsDigest(enumerations.DigestAlgorithm_MD5, data)
	if err != nil {
		return "", err
	}
	return utils.ToHex(d), nil
}

// ---------------------------------------------------------------------------------------------
// Byte array / file / document plumbing
// ---------------------------------------------------------------------------------------------

// DSSUtilsToByteArray reads the contents of the file at filePath into a byte array. Port of
// toByteArray(File), with openInputStream(File) inlined (its only caller).
func DSSUtilsToByteArray(filePath string) ([]byte, error) {
	info, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, model.NewDSSErrorMessageCause(
				fmt.Sprintf("Unable to read content of file '%s'", filePath),
				fmt.Errorf("File '%s' does not exist", filePath))
		}
		return nil, model.NewDSSErrorMessageCause(fmt.Sprintf("Unable to read content of file '%s'", filePath), err)
	}
	if info.IsDir() {
		return nil, model.NewDSSErrorMessageCause(
			fmt.Sprintf("Unable to read content of file '%s'", filePath),
			fmt.Errorf("File '%s' exists but is a directory", filePath))
	}
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, model.NewDSSErrorMessageCause(fmt.Sprintf("Unable to read content of file '%s'", filePath), err)
	}
	return data, nil
}

// DSSUtilsToByteArrayOfDocument gets the contents of document as a byte array. Port of
// toByteArray(DSSDocument).
func DSSUtilsToByteArrayOfDocument(document model.DSSDocument) ([]byte, error) {
	stream, err := document.OpenStream()
	if err != nil {
		return nil, model.NewDSSErrorMessageCause(
			fmt.Sprintf("Unable to read content of document with name '%s'. Reason : %s", document.Name(), err.Error()), err)
	}
	defer stream.Close()
	data, err := utils.ToByteArray(stream)
	if err != nil {
		return nil, model.NewDSSErrorMessageCause(
			fmt.Sprintf("Unable to read content of document with name '%s'. Reason : %s", document.Name(), err.Error()), err)
	}
	return data, nil
}

// DSSUtilsToByteArrayFromReader gets the content of r as a byte array. Port of
// toByteArray(InputStream).
func DSSUtilsToByteArrayFromReader(r io.Reader) ([]byte, error) {
	if r == nil {
		panic("The InputStream cannot be null")
	}
	data, err := utils.ToByteArray(r)
	if err != nil {
		return nil, model.NewDSSErrorMessageCause(fmt.Sprintf("Unable to read InputStream : %s", err.Error()), err)
	}
	return data, nil
}

// DSSUtilsIsEmpty verifies whether document is empty (has no body). Port of
// isEmpty(DSSDocument).
func DSSUtilsIsEmpty(document model.DSSDocument) (bool, error) {
	if _, ok := document.(*model.DigestDocument); ok {
		return true, nil
	}
	stream, err := document.OpenStream()
	if err != nil {
		return false, model.NewDSSErrorMessageCause(fmt.Sprintf("Unable to check if document has a content: %s", err.Error()), err)
	}
	defer stream.Close()
	buf := make([]byte, 1)
	n, err := stream.Read(buf)
	if err != nil && err != io.EOF {
		return false, model.NewDSSErrorMessageCause(fmt.Sprintf("Unable to check if document has a content: %s", err.Error()), err)
	}
	return n == 0, nil
}

// DSSUtilsFileByteSize returns the byte size of dssDocument. Port of getFileByteSize(DSSDocument).
func DSSUtilsFileByteSize(dssDocument model.DSSDocument) (int64, error) {
	stream, err := dssDocument.OpenStream()
	if err != nil {
		return 0, model.NewDSSErrorMessageCause(fmt.Sprintf("Cannot read the document with name [%s]", dssDocument.Name()), err)
	}
	defer stream.Close()
	size, err := utils.GetInputStreamSize(stream)
	if err != nil {
		return 0, model.NewDSSErrorMessageCause(fmt.Sprintf("Cannot read the document with name [%s]", dssDocument.Name()), err)
	}
	return size, nil
}

// DSSUtilsSaveToFile saves data to the file at filePath, creating parent directories as
// needed. Port of saveToFile(byte[], File).
func DSSUtilsSaveToFile(data []byte, filePath string) error {
	if dir := filepath.Dir(filePath); dir != "" {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return model.NewDSSErrorMessageCause(fmt.Sprintf("Unable to save a file : %s", err.Error()), err)
		}
	}
	if err := os.WriteFile(filePath, data, 0o644); err != nil {
		return model.NewDSSErrorMessageCause(fmt.Sprintf("Unable to save a file : %s", err.Error()), err)
	}
	return nil
}

// ---------------------------------------------------------------------------------------------
// String / URI / identifier helpers
// ---------------------------------------------------------------------------------------------

// DSSUtilsNormalizedString replaces all special characters by an underscore. Port of
// getNormalizedString(String).
func DSSUtilsNormalizedString(str string) string {
	if str == "" {
		return str
	}
	normalized := DSSUtilsDecodeURI(str)
	return regexp.MustCompile(`\W`).ReplaceAllString(normalized, "_")
}

// DSSUtilsDeterministicID returns a unique id for signingTime and id's DSS Id. Port of
// getDeterministicId(Date, TokenIdentifier).
//
// id may be nil (Java's null TokenIdentifier is a valid argument).
func DSSUtilsDeterministicID(signingTime time.Time, id *model.TokenIdentifier) (string, error) {
	var buf bytes.Buffer
	if !signingTime.IsZero() {
		_ = binary.Write(&buf, binary.BigEndian, signingTime.UnixMilli())
	}
	if id != nil {
		dssUtilsWriteChars(&buf, id.AsXmlID())
	}
	digest, err := DSSUtilsMD5Digest(buf.Bytes())
	if err != nil {
		return "", model.NewDSSErrorMessageCause(fmt.Sprintf("Unable to compute a deterministic Id : %s", err.Error()), err)
	}
	return "id-" + digest, nil
}

// DSSUtilsCounterSignatureDeterministicID returns a unique id for a counter signature. Port of
// getCounterSignatureDeterministicId(Date, TokenIdentifier, String).
func DSSUtilsCounterSignatureDeterministicID(signingTime time.Time, id *model.TokenIdentifier, masterSignatureID string) (string, error) {
	var buf bytes.Buffer
	if !signingTime.IsZero() {
		_ = binary.Write(&buf, binary.BigEndian, signingTime.UnixMilli())
	}
	if id != nil {
		dssUtilsWriteChars(&buf, id.AsXmlID())
	}
	if masterSignatureID != "" {
		dssUtilsWriteChars(&buf, masterSignatureID)
	}
	digest, err := DSSUtilsMD5Digest(buf.Bytes())
	if err != nil {
		return "", model.NewDSSErrorMessageCause(
			fmt.Sprintf("Unable to compute a deterministic Id for a counter-signature : %s", err.Error()), err)
	}
	return "id-" + digest, nil
}

// dssUtilsWriteChars writes s the way java.io.DataOutputStream#writeChars does: each UTF-16
// code unit of s as two big-endian bytes.
func dssUtilsWriteChars(buf *bytes.Buffer, s string) {
	for _, unit := range utf16.Encode([]rune(s)) {
		buf.WriteByte(byte(unit >> 8))
		buf.WriteByte(byte(unit))
	}
}

// DSSUtilsUTCDate returns a UTC date based on year, month and day; month is 0-based (0 for
// January), matching java.util.Calendar. Port of getUtcDate(int, int, int).
func DSSUtilsUTCDate(year, month, day int) time.Time {
	return time.Date(year, time.Month(month+1), day, 0, 0, 0, 0, time.UTC)
}

// DSSUtilsPrintSecurityProviders lists all defined security providers. Port of
// printSecurityProviders().
//
// DEVIATION: documented no-op, like DSSSecurityProviderInitSystemProviders - Go has no
// java.security.Security-style provider registry to enumerate.
func DSSUtilsPrintSecurityProviders() {
	// empty: there is no provider registry to enumerate.
}

// DSSUtilsReadFirstByte reads the first byte from dssDocument. Port of
// readFirstByte(DSSDocument).
//
// Matches Java's behaviour of ignoring the read count: an empty document yields the zero byte
// rather than an error, since Java's InputStream#read(byte[],int,int) return value is
// similarly ignored upstream.
func DSSUtilsReadFirstByte(dssDocument model.DSSDocument) (byte, error) {
	result := make([]byte, 1)
	stream, err := dssDocument.OpenStream()
	if err != nil {
		return 0, model.NewDSSErrorMessageCause(fmt.Sprintf("Cannot read first byte of the document. Reason : %s", err.Error()), err)
	}
	defer stream.Close()
	if _, err := stream.Read(result); err != nil && err != io.EOF {
		return 0, model.NewDSSErrorMessageCause(fmt.Sprintf("Cannot read first byte of the document. Reason : %s", err.Error()), err)
	}
	return result[0], nil
}

// DSSUtilsDecodeURI decodes a URI to be compliant with RFC 3986 (see DSS-2411 for details).
// Port of decodeURI(String); the slf4j warning upstream logs on a malformed URI is dropped
// (not load-bearing).
func DSSUtilsDecodeURI(uri string) string {
	replaced := strings.ReplaceAll(uri, "+", "%2B") // preserve '+' characters
	decoded, err := url.QueryUnescape(replaced)
	if err != nil {
		return uri
	}
	return decoded
}

// DSSUtilsEncodeURI encodes a URI (e.g. to be used within a ds:Reference element). Port of
// encodeURI(String).
//
// DEVIATION: upstream reconstructs the URI with java.net.URI(scheme, authority, path, query,
// fragment), whose multi-argument constructor percent-encodes each component through JDK
// internal per-component legal/quote character tables (see sun.net.www.ParseUtil), while
// preserving non-ASCII Unicode letters unescaped - exactly as the upstream doc comment
// describes. Reproducing those JDK tables byte-for-byte is out of scope here: dss-spi has no
// caller of this function yet, and XAdES (the module that builds ds:Reference URIs with it) is
// a later phase. This port instead percent-encodes only ASCII control characters, space and
// RFC 2396's "unwise" characters (<>"{}|\^`[]) component-wise, leaving Unicode letters and
// everything else untouched - matching the documented intent without matching the JDK
// byte-for-byte. Revisit with real interop vectors once XAdES lands.
func DSSUtilsEncodeURI(fileURI string) string {
	if fileURI == "" {
		return fileURI
	}
	match := dssUtilsRFC3986URIPattern.FindStringSubmatch(fileURI)
	if match == nil {
		return fileURI
	}
	scheme, hasAuthority, authority, path, hasQuery, query, hasFragment, fragment :=
		match[2], match[3], match[4], match[5], match[6], match[7], match[8], match[9]

	var out strings.Builder
	if scheme != "" {
		out.WriteString(scheme)
		out.WriteByte(':')
	}
	if hasAuthority != "" {
		out.WriteString("//")
		out.WriteString(dssUtilsEncodeURIComponent(authority))
	}
	out.WriteString(dssUtilsEncodeURIComponent(path))
	if hasQuery != "" {
		out.WriteByte('?')
		out.WriteString(dssUtilsEncodeURIComponent(query))
	}
	if hasFragment != "" {
		out.WriteByte('#')
		out.WriteString(dssUtilsEncodeURIComponent(fragment))
	}
	return out.String()
}

// dssUtilsUnwiseURIChars are the RFC 2396 "unwise" characters, percent-encoded by
// DSSUtilsEncodeURI alongside space and ASCII control characters.
const dssUtilsUnwiseURIChars = " <>\"{}|\\^`[]"

// dssUtilsEncodeURIComponent percent-encodes ASCII control characters, space and RFC 2396's
// "unwise" characters, leaving everything else - including non-ASCII Unicode - untouched.
func dssUtilsEncodeURIComponent(s string) string {
	var out strings.Builder
	for i := 0; i < len(s); i++ {
		b := s[i]
		if b < 0x20 || b == 0x7f || strings.IndexByte(dssUtilsUnwiseURIChars, b) >= 0 {
			fmt.Fprintf(&out, "%%%02X", b)
		} else {
			out.WriteByte(b)
		}
	}
	return out.String()
}

// DSSUtilsExceptionMessage returns a message retrieved from err, its cause's message if err's
// own is empty, or err's Go type name if neither is available. Port of
// getExceptionMessage(Exception).
//
// Panics with the Java message when err is nil, matching the upstream
// Objects.requireNonNull-style guard.
func DSSUtilsExceptionMessage(err error) string {
	if err == nil {
		panic("Cannot retrieve a message. The exception is null!")
	}
	if msg := err.Error(); msg != "" {
		return msg
	}
	if cause := errors.Unwrap(err); cause != nil {
		if msg := cause.Error(); msg != "" {
			return msg
		}
	}
	return fmt.Sprintf("%T", err)
}

// DSSUtilsRemoveControlCharacters replaces ASCII control characters 00-31 and 127 with "".
// Port of removeControlCharacters(String); the slf4j warning upstream logs when characters
// were removed is dropped (not load-bearing).
func DSSUtilsRemoveControlCharacters(str string) string {
	return dssUtilsControlCharsPattern.ReplaceAllString(str, "")
}

// DSSUtilsReplaceAllNonAlphanumericCharacters replaces all non-alphanumeric (Unicode letter or
// decimal digit) characters in str with replacement. Port of
// replaceAllNonAlphanumericCharacters(String, String).
func DSSUtilsReplaceAllNonAlphanumericCharacters(str, replacement string) string {
	return dssUtilsNonAlphanumericPattern.ReplaceAllString(str, replacement)
}

// DSSUtilsReplaceInvalidXmlCharacters replaces all XML 1.0 invalid characters in str with
// replacement. Port of replaceInvalidXmlCharacters(String, String).
func DSSUtilsReplaceInvalidXmlCharacters(str, replacement string) string {
	return dssUtilsInvalidXmlCharsPattern.ReplaceAllString(str, replacement)
}

// DSSUtilsIsUrnOid checks if id is a URN representation of an OID according to IETF RFC 3061.
// Port of isUrnOid(String).
func DSSUtilsIsUrnOid(id string) bool {
	return id != "" && dssUtilsURNOidPattern.MatchString(id)
}

// DSSUtilsIsOidCode checks if oid is a valid OID, e.g. "1.3.6.1.4.1.343" is valid, "25.25" and
// "http://sample.com" are not. Port of isOidCode(String).
func DSSUtilsIsOidCode(oid string) bool {
	return oid != "" && dssUtilsOidCodePattern.MatchString(oid)
}

// DSSUtilsOidCode keeps only the code part of an OID string, e.g. "urn:oid:1.2.3" becomes
// "1.2.3". Port of getOidCode(String).
func DSSUtilsOidCode(urnOid string) string {
	if urnOid == "" {
		return urnOid
	}
	return urnOid[strings.LastIndex(urnOid, ":")+1:]
}

// DSSUtilsToUrnOid returns a URN URI generated from oid, e.g. "1.2.4.5.6.8" becomes
// "urn:oid:1.2.4.5.6.8" (RFC 3061). Port of toUrnOid(String).
func DSSUtilsToUrnOid(oid string) string {
	return dssUtilsOidNamespacePrefix + oid
}

// DSSUtilsURIOrUrnOID returns the URI if present, otherwise the URN-encoded OID (RFC 3061), or
// "" if neither is present. Port of getUriOrUrnOid(ObjectIdentifier).
func DSSUtilsURIOrUrnOID(objectIdentifier enumerations.ObjectIdentifier) string {
	// TS 119 182-1 5.4.1 "The oId data type": if both an OID and a URI exist identifying one
	// object, the URI value should be used in the id member.
	uri := objectIdentifier.URI()
	if uri == "" && objectIdentifier.OID() != "" {
		uri = DSSUtilsToUrnOid(objectIdentifier.OID())
	}
	return uri
}

// DSSUtilsObjectIdentifierValue normalizes and retrieves a String identifier (for non-XAdES
// processing). Examples: "http://website.com" -> "http://website.com"; "urn:oid:1.2.3" ->
// "1.2.3"; "1.2.3" -> "1.2.3". Port of getObjectIdentifierValue(String).
func DSSUtilsObjectIdentifierValue(oidOrUriString string) string {
	return dssUtilsObjectIdentifierValue(oidOrUriString)
}

// DSSUtilsObjectIdentifierValueWithQualifier returns a URI value of oidOrUriString taking the
// given ObjectIdentifierQualifier into account (for XAdES processing). Port of
// getObjectIdentifierValue(String, ObjectIdentifierQualifier).
//
// qualifier is accepted for documentation parity with the Java overload; upstream only uses it
// to pick between two slf4j.debug diagnostics (dropped, not load-bearing per PORTING.md), so it
// does not otherwise affect the result.
func DSSUtilsObjectIdentifierValueWithQualifier(oidOrUriString string, qualifier enumerations.ObjectIdentifierQualifier) string {
	_ = qualifier
	return dssUtilsObjectIdentifierValue(oidOrUriString)
}

// dssUtilsObjectIdentifierValue implements the shared logic of both
// getObjectIdentifierValue(String) overloads (the xades-only diagnostics are dropped).
func dssUtilsObjectIdentifierValue(oidOrUriString string) string {
	value := oidOrUriString
	if utils.IsStringNotEmpty(oidOrUriString) {
		value = DSSUtilsTrimWhitespacesAndNewlines(value)
		if DSSUtilsIsUrnOid(value) {
			value = DSSUtilsOidCode(value)
		}
		// else: value is already an OID code or a plain URI; kept as-is either way.
	}
	return value
}

// DSSUtilsTrimWhitespacesAndNewlines trims whitespace and new line characters. Port of
// trimWhitespacesAndNewlines(String).
func DSSUtilsTrimWhitespacesAndNewlines(str string) string {
	str = strings.ReplaceAll(str, "\n", "")
	str = strings.ReplaceAll(str, "\r", "")
	return utils.Trim(str)
}

// DSSUtilsStripFirstLeadingOccurrence trims leading if it is a leading part of text. Port of
// stripFirstLeadingOccurrence(String, String).
//
// DEVIATION-BY-FIDELITY: upstream implements this via text.replaceFirst("^"+leading, ""), i.e.
// leading is a regular expression fragment, not a literal prefix - a caller passing regex
// metacharacters gets regex semantics in Java, faithfully reproduced here with Go's regexp
// (RE2 syntax, which is not always identical to java.util.regex but agrees for the plain
// literal strings every known call site uses).
func DSSUtilsStripFirstLeadingOccurrence(text, leading string) (string, error) {
	if text == "" {
		return text, nil
	}
	if leading == "" {
		return text, nil
	}
	re, err := regexp.Compile("^" + leading)
	if err != nil {
		return "", err
	}
	return re.ReplaceAllString(text, ""), nil
}

// DSSUtilsDocumentNames returns a list of document names from dssDocuments. Port of
// getDocumentNames(List).
func DSSUtilsDocumentNames(dssDocuments []model.DSSDocument) []string {
	if len(dssDocuments) == 0 {
		return []string{}
	}
	names := make([]string, 0, len(dssDocuments))
	for _, d := range dssDocuments {
		names = append(names, d.Name())
	}
	return names
}

// DSSUtilsDocumentWithName returns the document with fileName from documents, or nil when not
// found. Port of getDocumentWithName(List, String).
func DSSUtilsDocumentWithName(documents []model.DSSDocument, fileName string) model.DSSDocument {
	for _, document := range documents {
		if fileName == document.Name() {
			return document
		}
	}
	return nil
}

// DSSUtilsDocumentWithLastName returns the last document in ascending alphabetical name order.
// Port of getDocumentWithLastName(List).
func DSSUtilsDocumentWithLastName(documents []model.DSSDocument) model.DSSDocument {
	if len(documents) == 0 {
		return nil
	}
	names := DSSUtilsDocumentNames(documents)
	sortedNames := append([]string(nil), names...)
	dssUtilsSortStrings(sortedNames)
	return DSSUtilsDocumentWithName(documents, sortedNames[len(sortedNames)-1])
}

// dssUtilsSortStrings sorts values ascending in place (java.util.Collections#sort's natural
// String ordering).
func dssUtilsSortStrings(values []string) {
	for i := 1; i < len(values); i++ {
		for j := i; j > 0 && values[j-1] > values[j]; j-- {
			values[j-1], values[j] = values[j], values[j-1]
		}
	}
}

// DSSUtilsEnrichCollection adds every element of toAddCollection into *currentCollection that
// is not already present. Port of enrichCollection(Collection, Collection).
//
// DEVIATION: Java's Collection#contains uses equals(); Go's comparable constraint uses ==,
// which for a pointer type T is identity comparison, not value equality. Callers instantiating
// this with a pointer type (e.g. *model.CertificateToken, whose Java counterpart has a
// content-based equals()) get different de-duplication semantics than upstream and must
// compare via an explicit loop instead if content equality is required.
func DSSUtilsEnrichCollection[T comparable](currentCollection *[]T, toAddCollection []T) {
	for _, object := range toAddCollection {
		found := false
		for _, existing := range *currentCollection {
			if existing == object {
				found = true
				break
			}
		}
		if !found {
			*currentCollection = append(*currentCollection, object)
		}
	}
}

// ---------------------------------------------------------------------------------------------
// Signature-value / algorithm helpers
// ---------------------------------------------------------------------------------------------

// DSSUtilsConvertECSignatureValue ensures signatureValue has the format expected for
// expectedAlgorithm, converting between the ASN.1 SEQUENCE{r,s} (ECDSA) and plain R||S
// (PLAIN_ECDSA) encodings when required. Port of convertECSignatureValue(SignatureAlgorithm,
// SignatureValue).
func DSSUtilsConvertECSignatureValue(expectedAlgorithm enumerations.SignatureAlgorithm, signatureValue *model.SignatureValue) (*model.SignatureValue, error) {
	newSignatureValue := model.NewSignatureValue()
	newSignatureValue.SetAlgorithm(expectedAlgorithm)

	expectedEncryptionAlgorithm := expectedAlgorithm.EncryptionAlgorithm()
	signatureEncryptionAlgorithm := signatureValue.Algorithm().EncryptionAlgorithm()

	var signatureValueBinaries []byte
	var err error
	switch {
	case expectedEncryptionAlgorithm == enumerations.EncryptionAlgorithm_ECDSA &&
		signatureEncryptionAlgorithm == enumerations.EncryptionAlgorithm_PLAIN_ECDSA:
		signatureValueBinaries, err = DSSASN1UtilsToStandardDSASignatureValue(signatureValue.Value())

	case expectedEncryptionAlgorithm == enumerations.EncryptionAlgorithm_PLAIN_ECDSA &&
		signatureEncryptionAlgorithm == enumerations.EncryptionAlgorithm_ECDSA:
		signatureValueBinaries, err = DSSASN1UtilsToPlainDSASignatureValue(signatureValue.Value())

	default:
		return nil, model.NewDSSError(fmt.Sprintf(
			"Not supported conversion from SignatureAlgorithm '%s' defined within SignatureValue to the target algorithm '%s'",
			signatureValue.Algorithm(), expectedAlgorithm))
	}
	if err != nil {
		return nil, err
	}
	newSignatureValue.SetValue(signatureValueBinaries)
	return newSignatureValue, nil
}

// DSSUtilsEdDSASignatureAlgorithm returns the EdDSA SignatureAlgorithm used to create
// signatureValue. Only EdDSA algorithms are recognised; for any other length "" (Java's null)
// is returned. Port of getEdDSASignatureAlgorithm(byte[]); the slf4j warning upstream logs on
// an unrecognised length is dropped (not load-bearing).
//
// See RFC 8032 "Edwards-Curve Digital Signature Algorithm (EdDSA)" 4: EdDSA uses small public
// keys (32 or 57 bytes) and signatures (64 or 114 bytes) for Ed25519 and Ed448 respectively.
func DSSUtilsEdDSASignatureAlgorithm(signatureValue []byte) enumerations.SignatureAlgorithm {
	switch len(signatureValue) {
	case 64:
		return enumerations.SignatureAlgorithm_ED25519
	case 114:
		return enumerations.SignatureAlgorithm_ED448
	default:
		return ""
	}
}

// DSSUtilsAssertSPUserNoticeConfigurationValid verifies the validity of userNotice. Port of
// assertSPUserNoticeConfigurationValid(UserNotice); Java's thrown IllegalArgumentException
// becomes the returned error.
func DSSUtilsAssertSPUserNoticeConfigurationValid(userNotice *model.UserNotice) error {
	organizationEmpty := utils.IsStringEmpty(userNotice.Organization())
	noticeNumbersEmpty := len(userNotice.NoticeNumbers()) == 0
	if organizationEmpty != noticeNumbersEmpty {
		return errors.New("Both Organization name and NoticeNumbers shall be defined within the UserNotice configuration!")
	}
	return nil
}

// DSSUtilsToBigIntegerList transforms integers into a list of *big.Int. Port of
// toBigIntegerList(int[]).
func DSSUtilsToBigIntegerList(integers []int) []*big.Int {
	result := make([]*big.Int, 0, len(integers))
	for _, i := range integers {
		result = append(result, big.NewInt(int64(i)))
	}
	return result
}

// DSSUtilsIsLineBreakByte verifies if b represents a line break character (new line or
// carriage return). Port of isLineBreakByte(byte).
func DSSUtilsIsLineBreakByte(b byte) bool {
	for _, m := range dssUtilsLineBreakChars {
		if b == m {
			return true
		}
	}
	return false
}

// DSSUtilsGenerateKid generates the 'kid' value as in IETF RFC 5035. Port of
// generateKid(CertificateToken).
func DSSUtilsGenerateKid(signingCertificate *model.CertificateToken) []byte {
	return DSSASN1UtilsIssuerSerialForCertificate(signingCertificate).DER()
}

// DSSUtilsTimeValueInSeconds strips the milliseconds component from timeInMillis. Port of
// getTimeValueInSeconds(long).
func DSSUtilsTimeValueInSeconds(timeInMillis int64) int64 {
	return timeInMillis / 1000
}

// DSSUtilsTimeValueInMilliseconds adds a zero milliseconds component to timeWithoutMillis
// (seconds). Port of getTimeValueInMilliseconds(long).
func DSSUtilsTimeValueInMilliseconds(timeWithoutMillis int64) int64 {
	return timeWithoutMillis * 1000
}

// DSSUtilsTimeValueInMillisecondsFromFractionalSeconds adds the milliseconds component to
// timeWithoutMillis (fractional seconds). Port of getTimeValueInMilliseconds(double).
func DSSUtilsTimeValueInMillisecondsFromFractionalSeconds(timeWithoutMillis float64) int64 {
	return int64(timeWithoutMillis * 1000)
}

// DSSUtilsDateFromMilliseconds creates a Date from milliseconds since 1970-01-01T00:00Z. Port
// of getDateFromMilliseconds(Number); dateTimeNumber nil returns the zero Time (Java's null).
func DSSUtilsDateFromMilliseconds(dateTimeNumber *float64) time.Time {
	if dateTimeNumber == nil {
		return time.Time{}
	}
	return time.UnixMilli(int64(*dateTimeNumber)).UTC()
}

// DSSUtilsHost gets the host name from urlString, e.g.
// "ldap://ldap.infonotary.com/dc=identity-ca,dc=infonotary,dc=com" -> "ldap.infonotary.com".
// Port of getHost(String).
func DSSUtilsHost(urlString string) string {
	if utils.IsStringEmpty(urlString) {
		return ""
	}

	doubleSlash := strings.Index(urlString, "//")
	if doubleSlash == -1 {
		doubleSlash = 0
	} else {
		doubleSlash += 2
	}

	end := strings.IndexByte(urlString[doubleSlash:], '/')
	if end >= 0 {
		end += doubleSlash
	} else {
		end = len(urlString)
	}

	port := strings.IndexByte(urlString[doubleSlash:], ':')
	if port >= 0 {
		port += doubleSlash
	}
	if port > 0 && port < end {
		end = port
	}

	return urlString[doubleSlash:end]
}
