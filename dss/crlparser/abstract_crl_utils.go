// Ported from dss-crl-parser/src/main/java/eu/europa/esig/dss/crl/AbstractCRLUtils.java (DSS 6.5.RC1).
//
// DEVIATION: upstream is an abstract base class CRLUtilsX509CRLImpl (crl_utils_x509crl_impl.go)
// extends to share this code with the (unported) dss-crl-parser-stream implementation. With a
// single native implementation there is no subclass to share it with, so these become
// unexported package-level helpers that crl_utils_x509crl_impl.go calls directly; the
// "crlUtils" prefix keeps them grouped, standing in for the class they came from.
//
// A second deviation runs through every extraction helper below: upstream receives its input
// from java.security.cert.X509Extension#getExtensionValue(String), which returns the DER
// encoding of the extnValue OCTET STRING itself (tag, length AND content), so Java peels one
// OCTET STRING layer off before reading the real value. crypto/x509's pkix.Extension.Value is
// already that peeled content (see crypto/x509/parser.go's parseExtension, which reads the
// OCTET STRING and stores only its content octets), so this port reads it directly and the
// extra unwrap has no counterpart here.
package crlparser

import (
	encoding_asn1 "encoding/asn1"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ryftcore/dss-go/dss/model"
	"golang.org/x/crypto/cryptobyte"
	cryptobyte_asn1 "golang.org/x/crypto/cryptobyte/asn1"
)

// crlUtilsDERSequenceTag is the leading byte of a DER encoded SEQUENCE, i.e.
// BERTags.SEQUENCE | BERTags.CONSTRUCTED in BouncyCastle's terms.
const crlUtilsDERSequenceTag = 0x30

// crlUtilsGetDERContent returns the DER encoded content of binaries, converting it from PEM
// first when needed. Port of the private buildCRLBinary()/getDERContent()/isDerEncoded()/
// isPemEncoded() helpers of AbstractCRLUtils.
func crlUtilsGetDERContent(binaries []byte) ([]byte, error) {
	if len(binaries) == 0 {
		return nil, model.NewDSSError("Unsupported CRL. The obtained CRL content is empty!")
	}
	switch first := binaries[0]; {
	case first == crlUtilsDERSequenceTag:
		return binaries, nil
	case first == '-':
		return PemToDerConverterConvert(binaries)
	default:
		return nil, model.NewDSSError("Unable to load CRL. Not possible to convert to DER!")
	}
}

// crlUtilsExtractExpiredCertsOnCRL parses and sets the 'expiredCertsOnCRL' value.
// Port of the protected extractExpiredCertsOnCRL(CRLValidity, byte[]).
//
// A malformed or wrongly-typed value is silently ignored, matching upstream's catch-and-log
// behaviour (LOG.warn, no propagation).
func crlUtilsExtractExpiredCertsOnCRL(validity *CRLValidity, expiredCertsOnCRLValue []byte) {
	if expiredCertsOnCRLValue == nil {
		return
	}
	t, ok := crlUtilsParseGeneralizedTime(expiredCertsOnCRLValue)
	if !ok {
		return
	}
	validity.SetExpiredCertsOnCRL(&t)
}

// crlUtilsParseGeneralizedTime decodes value as a DER GeneralizedTime, reporting false (and no
// time) for anything else - including a syntactically valid but differently-tagged Time choice
// (UTCTime), which upstream also refuses via its
// "time.toASN1Primitive() instanceof ASN1GeneralizedTime" check.
func crlUtilsParseGeneralizedTime(value []byte) (time.Time, bool) {
	if len(value) == 0 || value[0] != byte(cryptobyte_asn1.GeneralizedTime) {
		return time.Time{}, false
	}
	var t time.Time
	rest, err := encoding_asn1.UnmarshalWithParams(value, &t, "generalized")
	if err != nil || len(rest) != 0 {
		return time.Time{}, false
	}
	return t, true
}

// crlUtilsExtractCrlNumber sets the CRL Number extension value, when present.
// Port of the protected extractCrlNumber(CRLValidity, byte[]).
//
// DEVIATION: upstream re-parses the raw extension bytes itself, reading the encoded INTEGER
// with ASN1Integer#getPositiveValue() (an unsigned magnitude read, defensive against CRL
// Number encodings that omit the leading padding zero byte DER requires for a non-negative
// value with its high bit set). crypto/x509.ParseRevocationList already exposes this
// extension, correctly DER-decoded, as RevocationList.Number; every CRL Number this port has
// been exercised against is a small non-negative integer, so its ordinary (signed) INTEGER
// decoding and upstream's defensive unsigned one agree, and crl_utils_x509crl_impl.go passes
// that field straight through here instead of re-parsing.
func crlUtilsExtractCrlNumber(validity *CRLValidity, crlNumber *big.Int) {
	if crlNumber == nil {
		return
	}
	validity.SetCRLNumber(crlNumber)
}

// crlUtilsIssuingDistributionPoint carries the fields decoded out of an IssuingDistributionPoint.
type crlUtilsIssuingDistributionPoint struct {
	url                 string
	onlyUserCerts       bool
	onlyCaCerts         bool
	onlySomeReasonFlags *encoding_asn1.BitString
	indirectCrl         bool
	onlyAttributeCerts  bool
}

// crlUtilsExtractIssuingDistributionPointBinary parses the issuing distribution point binaries
// and sets the corresponding CRLValidity fields.
// Port of the protected extractIssuingDistributionPointBinary(CRLValidity, byte[]).
//
// Unlike crlUtilsExtractExpiredCertsOnCRL/crlUtilsExtractCrlNumber, upstream's method has no
// try/catch: a malformed IssuingDistributionPoint propagates the BouncyCastle
// IllegalArgumentException out of buildCRLValidity as an unchecked exception. This port
// preserves that by returning the error instead of swallowing it, which
// crl_utils_x509crl_impl.go's caller in turn returns from CRLUtilsBuildCRLValidity.
func crlUtilsExtractIssuingDistributionPointBinary(validity *CRLValidity, issuingDistributionPointValue []byte) error {
	if issuingDistributionPointValue == nil {
		// Upstream logs a debug message here; the IDP fields are simply left unset.
		return nil
	}
	idp, err := crlUtilsParseIssuingDistributionPoint(issuingDistributionPointValue)
	if err != nil {
		return err
	}
	validity.SetOnlyAttributeCerts(idp.onlyAttributeCerts)
	validity.SetOnlyCaCerts(idp.onlyCaCerts)
	validity.SetOnlyUserCerts(idp.onlyUserCerts)
	validity.SetIndirectCrl(idp.indirectCrl)
	validity.SetReasonFlags(idp.onlySomeReasonFlags)
	validity.SetURL(idp.url)
	return nil
}

// crlUtilsParseIssuingDistributionPoint decodes the DER content of the issuingDistributionPoint
// extension:
//
//	IssuingDistributionPoint ::= SEQUENCE {
//	     distributionPoint          [0] DistributionPointName OPTIONAL,
//	     onlyContainsUserCerts      [1] BOOLEAN DEFAULT FALSE,
//	     onlyContainsCACerts        [2] BOOLEAN DEFAULT FALSE,
//	     onlySomeReasons            [3] ReasonFlags OPTIONAL,
//	     indirectCRL                [4] BOOLEAN DEFAULT FALSE,
//	     onlyContainsAttributeCerts [5] BOOLEAN DEFAULT FALSE }
func crlUtilsParseIssuingDistributionPoint(der []byte) (crlUtilsIssuingDistributionPoint, error) {
	var result crlUtilsIssuingDistributionPoint

	input := cryptobyte.String(der)
	var seq cryptobyte.String
	if !input.ReadASN1(&seq, cryptobyte_asn1.SEQUENCE) {
		return result, errors.New("invalid IssuingDistributionPoint: not a SEQUENCE")
	}

	// distributionPoint [0] DistributionPointName OPTIONAL. The tag is EXPLICIT because
	// DistributionPointName is a CHOICE.
	var distributionPoint cryptobyte.String
	var hasDistributionPoint bool
	if !seq.ReadOptionalASN1(&distributionPoint, &hasDistributionPoint, cryptobyte_asn1.Tag(0).Constructed().ContextSpecific()) {
		return result, errors.New("invalid IssuingDistributionPoint: malformed distributionPoint")
	}
	if hasDistributionPoint {
		url, err := crlUtilsParseDistributionPointNameURL(distributionPoint)
		if err != nil {
			return result, err
		}
		result.url = url
	}

	var err error
	if result.onlyUserCerts, err = crlUtilsReadOptionalImplicitBoolean(&seq, 1); err != nil {
		return result, fmt.Errorf("invalid IssuingDistributionPoint: malformed onlyContainsUserCerts: %w", err)
	}
	if result.onlyCaCerts, err = crlUtilsReadOptionalImplicitBoolean(&seq, 2); err != nil {
		return result, fmt.Errorf("invalid IssuingDistributionPoint: malformed onlyContainsCACerts: %w", err)
	}

	// onlySomeReasons [3] ReasonFlags OPTIONAL, IMPLICIT (ReasonFlags is a BIT STRING, not a
	// CHOICE, so implicit tagging applies and the content is the BIT STRING content directly:
	// one byte of unused-bit count followed by the data bytes).
	var reasonsContent cryptobyte.String
	var hasReasons bool
	if !seq.ReadOptionalASN1(&reasonsContent, &hasReasons, cryptobyte_asn1.Tag(3).ContextSpecific()) {
		return result, errors.New("invalid IssuingDistributionPoint: malformed onlySomeReasons")
	}
	if hasReasons {
		if len(reasonsContent) < 1 {
			return result, errors.New("invalid IssuingDistributionPoint: empty onlySomeReasons")
		}
		unusedBits := int(reasonsContent[0])
		dataBytes := []byte(reasonsContent[1:])
		// ASN1BitString#createPrimitive refuses these ("invalid pad bits detected"); taken
		// as they come they made BitLength negative.
		if unusedBits > 7 || (unusedBits > 0 && len(dataBytes) == 0) {
			return result, errors.New("invalid IssuingDistributionPoint: invalid onlySomeReasons pad bits")
		}
		result.onlySomeReasonFlags = &encoding_asn1.BitString{
			Bytes:     dataBytes,
			BitLength: len(dataBytes)*8 - unusedBits,
		}
	}

	if result.indirectCrl, err = crlUtilsReadOptionalImplicitBoolean(&seq, 4); err != nil {
		return result, fmt.Errorf("invalid IssuingDistributionPoint: malformed indirectCRL: %w", err)
	}
	if result.onlyAttributeCerts, err = crlUtilsReadOptionalImplicitBoolean(&seq, 5); err != nil {
		return result, fmt.Errorf("invalid IssuingDistributionPoint: malformed onlyContainsAttributeCerts: %w", err)
	}

	return result, nil
}

// crlUtilsReadOptionalImplicitBoolean reads an OPTIONAL/DEFAULT FALSE BOOLEAN carrying an
// IMPLICIT context-specific primitive tag out of seq, returning false when absent.
//
// cryptobyte's own String.ReadOptionalASN1Boolean cannot be used here: it re-reads its content
// through ReadASN1Boolean, which expects to see the universal BOOLEAN tag (0x01) again - correct
// for an EXPLICIT or untagged optional BOOLEAN, but IMPLICIT tagging (as used throughout
// IssuingDistributionPoint) replaces that tag rather than wrapping it, so by the time
// ReadOptionalASN1 has stripped the context tag, only the bare content octet is left.
func crlUtilsReadOptionalImplicitBoolean(seq *cryptobyte.String, tagNumber byte) (bool, error) {
	var content cryptobyte.String
	var present bool
	if !seq.ReadOptionalASN1(&content, &present, cryptobyte_asn1.Tag(tagNumber).ContextSpecific()) {
		return false, errors.New("malformed BOOLEAN")
	}
	if !present {
		return false, nil
	}
	if len(content) != 1 {
		return false, errors.New("BOOLEAN content is not a single byte")
	}
	switch content[0] {
	case 0x00:
		return false, nil
	case 0xFF:
		return true, nil
	default:
		return false, errors.New("BOOLEAN content is not DER canonical (0x00/0xFF)")
	}
}

// crlUtilsParseDistributionPointNameURL extracts the uniformResourceIdentifier GeneralName from
// a DistributionPointName's fullName choice, mirroring AbstractCRLUtils's private getUrl(...):
// only the FULL_NAME alternative is inspected, and the first URI GeneralName found wins.
//
//	DistributionPointName ::= CHOICE {
//	     fullName                [0] GeneralNames,
//	     nameRelativeToCRLIssuer [1] RelativeDistinguishedName }
//
//	GeneralNames ::= SEQUENCE SIZE (1..MAX) OF GeneralName
//
// GeneralName's uniformResourceIdentifier [6] IA5String alternative is IMPLICIT (IA5String is
// not a CHOICE), so its content is the URI's ASCII bytes directly.
func crlUtilsParseDistributionPointNameURL(distributionPointName cryptobyte.String) (string, error) {
	var fullName cryptobyte.String
	var isFullName bool
	if !distributionPointName.ReadOptionalASN1(&fullName, &isFullName, cryptobyte_asn1.Tag(0).Constructed().ContextSpecific()) {
		return "", errors.New("invalid DistributionPointName")
	}
	if !isFullName {
		return "", nil
	}
	for !fullName.Empty() {
		var element cryptobyte.String
		var tag cryptobyte_asn1.Tag
		if !fullName.ReadAnyASN1(&element, &tag) {
			return "", errors.New("invalid GeneralName")
		}
		if tag == cryptobyte_asn1.Tag(6).ContextSpecific() {
			return string(element), nil
		}
	}
	return "", nil
}
