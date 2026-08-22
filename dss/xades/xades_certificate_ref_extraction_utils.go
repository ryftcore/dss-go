// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/validation/XAdESCertificateRefExtractionUtils.java
// (DSS 6.5.RC1).
//
// xadesASN1UtilsGetX500PrincipalOrNull ports eu.europa.esig.dss.spi.DSSASN1Utils#getX500PrincipalOrNull(String):
// DSSASN1Utils (Java eu.europa.esig.dss.spi.DSSASN1Utils, flattened into the frozen spi package
// per PORTING.md) is landed everywhere else this port needs it (spi/dss_asn1_utils.go), but that
// one static overload - parsing an RFC 2253/4514 distinguished-name STRING back into an
// X500Principal, the mirror image of model.X500Principal.RFC2253Name()/RFC2253NameWithOIDMap -
// has no Go counterpart there. PORTING.md forbids editing a frozen package, so it is reproduced
// here, scoped to this package; it would be a good candidate to hoist into spi/dss_asn1_utils.go
// later (every future caller of the same Java static method would otherwise have to duplicate it
// again).
//
// DEVIATION: the parser supports the RFC 2253 keyword set model.X500Principal's own RFC2253Name()
// encoder emits verbatim (CN, C, L, ST, O, OU, STREET, DC, UID - see model/x500_principal.go's
// x500PrincipalRFC2253Keywords table, which this table is the deliberate inverse of) plus a
// handful of the other common keywords javax.security.auth.x500.X500Principal recognises by
// default (T, SURNAME/SN, GIVENNAME/GN, SERIALNUMBER, EMAILADDRESS) and bare dotted-decimal
// OIDs; it does not additionally recognise the uppercase long-form X520Attributes descriptions
// Java's caller merges in as a supplementary keyword map (e.g. "ORGANIZATIONNAME=" instead of
// "O="), since no DSS-generated XML ever emits those. A hex "#..." AVA value (a raw pre-encoded
// DER value, per RFC 2253 section 2.4) is passed through as the AttributeTypeAndValue's raw
// bytes rather than decoded into a Go string, preserving the exact encoding for that rare case.
package xades

import (
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/ryftcore/dss-go/dss/internal/xmldom"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/utils"
	"github.com/ryftcore/dss-go/dss/xades/definition"
	xmlutils "github.com/ryftcore/dss-go/dss/xml/utils"
)

// CertificateRefExtractionUtilsCreateCertificateRefFromV1 extracts a CertificateRef from a
// V1 certRefElement. Port of the public static createCertificateRefFromV1(Element, XAdESPath).
func CertificateRefExtractionUtilsCreateCertificateRefFromV1(certRefElement *xmldom.Node, xadesPaths definition.XAdESPath) *spi.CertificateRef {
	if certRefElement == nil {
		return nil
	}
	digestElement, _ := xmlutils.XPathUtilsGetElement(certRefElement, xadesPaths.CurrentCertDigest())
	certDigest := DSSXMLUtilsGetDigestAndValue(digestElement)
	if certDigest.IsEmpty() {
		return nil
	}
	certRef := spi.NewCertificateRef()
	certRef.SetCertDigest(certDigest)
	certRef.SetCertificateIdentifier(xadesCertificateRefExtractionUtilsCertificateIdentifierV1(certRefElement, xadesPaths))
	return certRef
}

// CertificateRefExtractionUtilsCreateCertificateRefFromV2 extracts a CertificateRef from a
// V2 certRefElement. Port of the public static createCertificateRefFromV2(Element, XAdESPath).
func CertificateRefExtractionUtilsCreateCertificateRefFromV2(certRefElement *xmldom.Node, xadesPaths definition.XAdESPath) *spi.CertificateRef {
	if certRefElement == nil {
		return nil
	}
	digestElement, _ := xmlutils.XPathUtilsGetElement(certRefElement, xadesPaths.CurrentCertDigest())
	certDigest := DSSXMLUtilsGetDigestAndValue(digestElement)
	if certDigest.IsEmpty() {
		return nil
	}
	certRef := spi.NewCertificateRef()
	certRef.SetCertDigest(certDigest)
	certRef.SetCertificateIdentifier(xadesCertificateRefExtractionUtilsCertificateIdentifierV2(certRefElement, xadesPaths))
	return certRef
}

// xadesCertificateRefExtractionUtilsCertificateIdentifierV1 ports the private static
// getCertificateIdentifierV1(Element, XAdESPath).
//
// LOG.warn("Unable to build a SignerIdentifier from CertIDTypeV2!") - copied verbatim from
// upstream, including its V1/V2 message mismatch - is dropped per PORTING.md.
func xadesCertificateRefExtractionUtilsCertificateIdentifierV1(certRefElement *xmldom.Node, xadesPaths definition.XAdESPath) *spi.SignerIdentifier {
	var issuerName *model.X500Principal
	issuerNameEl, _ := xmlutils.XPathUtilsGetElement(certRefElement, xadesPaths.CurrentIssuerSerialIssuerNamePath())
	if issuerNameEl != nil {
		issuerName = xadesASN1UtilsGetX500PrincipalOrNull(issuerNameEl.TextContent())
	}

	var serialNumber *big.Int
	serialNumberEl, _ := xmlutils.XPathUtilsGetElement(certRefElement, xadesPaths.CurrentIssuerSerialSerialNumberPath())
	if serialNumberEl != nil {
		serialNumberText := strings.TrimSpace(serialNumberEl.TextContent())
		if utils.IsStringDigits(serialNumberText) {
			serialNumber, _ = new(big.Int).SetString(serialNumberText, 10)
		}
		// Upstream logs "Unable to parse SerialNumber from 'CertIDTypeV1' element. Not a numeric!".
	}

	if issuerName == nil || serialNumber == nil {
		return nil
	}

	signerIdentifier := spi.NewSignerIdentifier()
	signerIdentifier.SetIssuerName(issuerName)
	signerIdentifier.SetSerialNumber(serialNumber)
	return signerIdentifier
}

// xadesCertificateRefExtractionUtilsCertificateIdentifierV2 ports the private static
// getCertificateIdentifierV2(Element, XAdESPath).
func xadesCertificateRefExtractionUtilsCertificateIdentifierV2(certRefElement *xmldom.Node, xadesPaths definition.XAdESPath) *spi.SignerIdentifier {
	issuerSerialV2Element, _ := xmlutils.XPathUtilsGetElement(certRefElement, xadesPaths.CurrentIssuerSerialV2Path())
	if issuerSerialV2Element == nil {
		// Tag issuerSerialV2 is optional.
		return nil
	}

	textContent := issuerSerialV2Element.TextContent()
	if !utils.IsBase64Encoded(textContent) {
		// Upstream logs "The IssuerSerialV2 value is not base64-encoded!".
		return nil
	}
	binaries := utils.FromBase64(textContent)
	issuerSerial := spi.DSSASN1UtilsIssuerSerial(binaries)
	if issuerSerial == nil {
		return nil
	}
	return spi.DSSASN1UtilsToSignerIdentifierFromIssuerSerial(issuerSerial)
}

// xadesX500PrincipalKeywords maps an RFC 2253/4514 keyword to its attribute type OID, the
// deliberate inverse of model.X500Principal's own RFC2253Name() keyword table plus a small set
// of other common javax.security.auth.x500.X500Principal default keywords. See this file's
// header for the fidelity note.
var xadesX500PrincipalKeywords = map[string]string{
	"CN":           "2.5.4.3",
	"C":            "2.5.4.6",
	"L":            "2.5.4.7",
	"ST":           "2.5.4.8",
	"O":            "2.5.4.10",
	"OU":           "2.5.4.11",
	"STREET":       "2.5.4.9",
	"DC":           "0.9.2342.19200300.100.1.25",
	"UID":          "0.9.2342.19200300.100.1.1",
	"T":            "2.5.4.12",
	"TITLE":        "2.5.4.12",
	"SN":           "2.5.4.4",
	"SURNAME":      "2.5.4.4",
	"GN":           "2.5.4.42",
	"GIVENNAME":    "2.5.4.42",
	"SERIALNUMBER": "2.5.4.5",
	"EMAILADDRESS": "1.2.840.113549.1.9.1",
	"DNQ":          "2.5.4.46",
	"DNQUALIFIER":  "2.5.4.46",
	"GENERATION":   "2.5.4.44",
	"INITIALS":     "2.5.4.43",
}

// xadesASN1UtilsGetX500PrincipalOrNull ports DSSASN1Utils#getX500PrincipalOrNull(String); see
// this file's header for the parser's scope.
func xadesASN1UtilsGetX500PrincipalOrNull(x500PrincipalString string) *model.X500Principal {
	rdnSequence, err := xadesParseRFC2253DN(x500PrincipalString)
	if err != nil {
		// Upstream logs "Unable to create an instance of X500Principal".
		return nil
	}
	der, err := asn1.Marshal(rdnSequence)
	if err != nil {
		return nil
	}
	principal, err := model.NewX500Principal(der)
	if err != nil {
		return nil
	}
	return principal
}

// xadesParseRFC2253DN parses an RFC 2253/4514 distinguished name string into an ASN.1 Name
// (SEQUENCE OF RelativeDistinguishedName), root-first order (the reverse of the string's
// leaf-first order, mirroring model.X500Principal.generateRFC2253DN's own reversal).
func xadesParseRFC2253DN(s string) (pkix.RDNSequence, error) {
	leafFirst, err := xadesSplitDNIntoRDNs(s)
	if err != nil {
		return nil, err
	}
	rdnSequence := make(pkix.RDNSequence, 0, len(leafFirst))
	for i := len(leafFirst) - 1; i >= 0; i-- {
		rdnText := leafFirst[i]
		avaTexts, err := xadesSplitRDNIntoAVAs(rdnText)
		if err != nil {
			return nil, err
		}
		set := make(pkix.RelativeDistinguishedNameSET, 0, len(avaTexts))
		for _, avaText := range avaTexts {
			ava, err := xadesParseAVA(avaText)
			if err != nil {
				return nil, err
			}
			set = append(set, ava)
		}
		rdnSequence = append(rdnSequence, set)
	}
	return rdnSequence, nil
}

// xadesSplitDNIntoRDNs splits a DN string into its RDNs (leaf-first, as written), on unescaped
// ',' or ';' separators.
func xadesSplitDNIntoRDNs(s string) ([]string, error) {
	return xadesSplitUnescaped(s, ",;")
}

// xadesSplitRDNIntoAVAs splits a (possibly multi-valued) RDN into its AttributeTypeAndValues, on
// unescaped '+' separators.
func xadesSplitRDNIntoAVAs(s string) ([]string, error) {
	return xadesSplitUnescaped(s, "+")
}

// xadesSplitUnescaped splits s on any of seps, honouring backslash escapes and double-quoted
// substrings (in which none of seps, nor '\', triggers a split).
func xadesSplitUnescaped(s string, seps string) ([]string, error) {
	var parts []string
	var current strings.Builder
	inQuotes := false
	runes := []rune(s)
	for i := 0; i < len(runes); i++ {
		c := runes[i]
		switch {
		case c == '"' && !inQuotes:
			inQuotes = true
			current.WriteRune(c)
		case c == '"' && inQuotes:
			inQuotes = false
			current.WriteRune(c)
		case c == '\\' && i+1 < len(runes):
			current.WriteRune(c)
			i++
			current.WriteRune(runes[i])
		case !inQuotes && strings.ContainsRune(seps, c):
			parts = append(parts, current.String())
			current.Reset()
		default:
			current.WriteRune(c)
		}
	}
	if inQuotes {
		return nil, fmt.Errorf("xades: unterminated quoted value in distinguished name %q", s)
	}
	parts = append(parts, current.String())
	return parts, nil
}

// xadesParseAVA parses a single "type=value" AttributeTypeAndValue.
func xadesParseAVA(s string) (pkix.AttributeTypeAndValue, error) {
	eq := xadesFindUnescapedEquals(s)
	if eq < 0 {
		return pkix.AttributeTypeAndValue{}, fmt.Errorf("xades: missing '=' in AVA %q", s)
	}
	keyword := strings.TrimSpace(s[:eq])
	rawValue := strings.TrimSpace(s[eq+1:])

	oid, err := xadesResolveAttributeOID(keyword)
	if err != nil {
		return pkix.AttributeTypeAndValue{}, err
	}

	if strings.HasPrefix(rawValue, "#") {
		raw, err := hex.DecodeString(rawValue[1:])
		if err != nil {
			return pkix.AttributeTypeAndValue{}, fmt.Errorf("xades: invalid hex AVA value %q: %w", rawValue, err)
		}
		return pkix.AttributeTypeAndValue{Type: oid, Value: asn1.RawValue{FullBytes: raw}}, nil
	}

	value, err := xadesUnescapeAVAValue(rawValue)
	if err != nil {
		return pkix.AttributeTypeAndValue{}, err
	}
	return pkix.AttributeTypeAndValue{Type: oid, Value: value}, nil
}

// xadesFindUnescapedEquals returns the index of the first unescaped, unquoted '=' in s, -1 if
// none is found.
func xadesFindUnescapedEquals(s string) int {
	inQuotes := false
	runes := []rune(s)
	for i := 0; i < len(runes); i++ {
		switch {
		case runes[i] == '"':
			inQuotes = !inQuotes
		case runes[i] == '\\' && i+1 < len(runes):
			i++
		case !inQuotes && runes[i] == '=':
			return len(string(runes[:i]))
		}
	}
	return -1
}

// xadesResolveAttributeOID resolves an RFC 2253/4514 attribute type keyword or dotted-decimal
// OID string to an asn1.ObjectIdentifier.
func xadesResolveAttributeOID(keyword string) (asn1.ObjectIdentifier, error) {
	upper := strings.ToUpper(keyword)
	if oidText, ok := xadesX500PrincipalKeywords[upper]; ok {
		return xadesParseOID(oidText)
	}
	if strings.HasPrefix(upper, "OID.") {
		return xadesParseOID(keyword[len("OID."):])
	}
	if len(keyword) > 0 && (keyword[0] >= '0' && keyword[0] <= '9') {
		return xadesParseOID(keyword)
	}
	return nil, fmt.Errorf("xades: unrecognized attribute type keyword %q", keyword)
}

// xadesParseOID parses a dotted-decimal OID string.
func xadesParseOID(text string) (asn1.ObjectIdentifier, error) {
	parts := strings.Split(text, ".")
	oid := make(asn1.ObjectIdentifier, 0, len(parts))
	for _, p := range parts {
		n, err := strconv.Atoi(p)
		if err != nil {
			return nil, fmt.Errorf("xades: invalid OID component %q in %q: %w", p, text, err)
		}
		oid = append(oid, n)
	}
	return oid, nil
}

// xadesUnescapeAVAValue removes RFC 2253/4514 quoting and backslash-escaping from an AVA value,
// and trims unescaped leading/trailing spaces per the RFC.
func xadesUnescapeAVAValue(s string) (string, error) {
	runes := []rune(s)
	if len(runes) >= 2 && runes[0] == '"' && runes[len(runes)-1] == '"' {
		runes = runes[1 : len(runes)-1]
	}

	var buf strings.Builder
	for i := 0; i < len(runes); i++ {
		c := runes[i]
		if c == '\\' && i+1 < len(runes) {
			next := runes[i+1]
			if xadesIsHexDigit(next) && i+2 < len(runes) && xadesIsHexDigit(runes[i+2]) {
				b, err := hex.DecodeString(string([]rune{next, runes[i+2]}))
				if err != nil {
					return "", fmt.Errorf("xades: invalid hex escape in AVA value %q: %w", s, err)
				}
				buf.WriteByte(b[0])
				i += 2
				continue
			}
			buf.WriteRune(next)
			i++
			continue
		}
		buf.WriteRune(c)
	}
	return strings.TrimSpace(buf.String()), nil
}

// xadesIsHexDigit reports whether r is an ASCII hexadecimal digit.
func xadesIsHexDigit(r rune) bool {
	return (r >= '0' && r <= '9') || (r >= 'a' && r <= 'f') || (r >= 'A' && r <= 'F')
}
