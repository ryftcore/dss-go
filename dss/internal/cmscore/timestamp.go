// RFC 3161 time-stamp structures, replacing org.bouncycastle.asn1.tsp.* (TSTInfo,
// MessageImprint, Accuracy, TimeStampResp), org.bouncycastle.asn1.cmp.PKIStatusInfo and
// org.bouncycastle.tsp.TimeStampToken / TimeStampTokenInfo.
//
// Parse only: a validator reads time-stamps, a TSA writes them.
package cmscore

import (
	"encoding/asn1"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ryftcore/dss-go/dss/internal/asn1ber"
)

// PKIStatus values of RFC 3161 clause 2.4.2 (imported from RFC 2510).
const (
	// PKIStatusGranted means the token was issued as requested.
	PKIStatusGranted = 0
	// PKIStatusGrantedWithMods means a token was issued, with modifications.
	PKIStatusGrantedWithMods = 1
	// PKIStatusRejection means the request was refused; no token is present.
	PKIStatusRejection = 2
	// PKIStatusWaiting means the request has not been processed yet.
	PKIStatusWaiting = 3
	// PKIStatusRevocationWarning warns of an imminent revocation.
	PKIStatusRevocationWarning = 4
	// PKIStatusRevocationNotification notifies of a revocation.
	PKIStatusRevocationNotification = 5
)

// PKIFailureInfo bit numbers of RFC 3161 clause 2.4.2.
const (
	// PKIFailureBadAlg means the requested algorithm is not recognised or supported.
	PKIFailureBadAlg = 0
	// PKIFailureBadRequest means the transaction is not permitted or supported.
	PKIFailureBadRequest = 2
	// PKIFailureBadDataFormat means the data submitted has the wrong format.
	PKIFailureBadDataFormat = 5
	// PKIFailureTimeNotAvailable means the TSA's time source is not available.
	PKIFailureTimeNotAvailable = 14
	// PKIFailureUnacceptedPolicy means the requested policy is not supported.
	PKIFailureUnacceptedPolicy = 15
	// PKIFailureUnacceptedExtension means a requested extension is not supported.
	PKIFailureUnacceptedExtension = 16
	// PKIFailureAddInfoNotAvailable means the additional information asked for is unavailable.
	PKIFailureAddInfoNotAvailable = 17
	// PKIFailureSystemFailure means the request cannot be handled owing to a system failure.
	PKIFailureSystemFailure = 25
)

// PKIStatusInfo is
//
//	PKIStatusInfo ::= SEQUENCE {
//	    status       PKIStatus,
//	    statusString PKIFreeText OPTIONAL,
//	    failInfo     PKIFailureInfo OPTIONAL }
type PKIStatusInfo struct {
	// Status is one of the PKIStatus values.
	Status int
	// StatusString holds the UTF8String components of statusString, nil when absent.
	StatusString []string
	// FailInfo holds the value bits of the failInfo BIT STRING, nil when absent.
	FailInfo []byte
	// FailInfoUnusedBits is the number of unused bits of the failInfo BIT STRING.
	FailInfoUnusedBits int
}

// IsGranted reports whether a token was issued, i.e. the status is granted or
// grantedWithMods. Port of TimeStampResponse#getStatus() checks in DSS.
func (p *PKIStatusInfo) IsGranted() bool {
	return p.Status == PKIStatusGranted || p.Status == PKIStatusGrantedWithMods
}

// HasFailure reports whether the given PKIFailureInfo bit is set.
func (p *PKIStatusInfo) HasFailure(bit int) bool {
	index := bit / 8
	if index >= len(p.FailInfo) {
		return false
	}
	// A BIT STRING numbers its bits from the most significant one of the first octet.
	return p.FailInfo[index]&(0x80>>(bit%8)) != 0
}

// pkiStatusInfoFromElement decodes an already parsed PKIStatusInfo.
func pkiStatusInfoFromElement(element *asn1ber.Element) (*PKIStatusInfo, error) {
	children, err := expectSequence(element, "PKIStatusInfo", 1)
	if err != nil {
		return nil, err
	}
	status, err := expectSmallInteger(children[0], "PKIStatusInfo.status")
	if err != nil {
		return nil, err
	}
	info := &PKIStatusInfo{Status: status}
	for _, child := range children[1:] {
		switch {
		case child.IsUniversal(asn1ber.TagSequence):
			for _, text := range child.Children() {
				info.StatusString = append(info.StatusString, text.AsString())
			}
		case child.IsUniversal(asn1ber.TagBitString):
			content := child.DERContent()
			if len(content) == 0 {
				return nil, errors.New("cmscore: PKIStatusInfo.failInfo is an empty BIT STRING")
			}
			info.FailInfoUnusedBits = int(content[0])
			info.FailInfo = content[1:]
		default:
			return nil, errors.New("cmscore: unexpected PKIStatusInfo component")
		}
	}
	return info, nil
}

// TimeStampResp is
//
//	TimeStampResp ::= SEQUENCE {
//	    status         PKIStatusInfo,
//	    timeStampToken TimeStampToken OPTIONAL }
//
// The token is absent when the request was refused.
type TimeStampResp struct {
	// Status is the outcome of the request.
	Status *PKIStatusInfo
	// Token is the issued token, nil when the request was refused.
	Token *TimeStampToken

	// element is the parsed TimeStampResp.
	element *asn1ber.Element
}

// ParseTimeStampResp decodes a TimeStampResp, i.e. the body of an RFC 3161 HTTP reply.
func ParseTimeStampResp(input []byte) (*TimeStampResp, error) {
	element, err := parseOne(input, "TimeStampResp")
	if err != nil {
		return nil, err
	}
	children, err := expectSequence(element, "TimeStampResp", 1)
	if err != nil {
		return nil, err
	}
	status, err := pkiStatusInfoFromElement(children[0])
	if err != nil {
		return nil, err
	}
	response := &TimeStampResp{Status: status, element: element}
	if len(children) > 1 {
		token, err := TimeStampTokenFromElement(children[1])
		if err != nil {
			return nil, err
		}
		response.Token = token
	}
	return response, nil
}

// Element returns the parsed TimeStampResp.
func (t *TimeStampResp) Element() *asn1ber.Element { return t.element }

// Encoded returns the TimeStampResp's original encoding.
func (t *TimeStampResp) Encoded() []byte { return t.element.Encoded() }

// TimeStampToken is the RFC 3161 token: a CMS SignedData whose eContentType is id-ct-TSTInfo
// and whose eContent is the DER encoding of a TSTInfo.
//
// The token's own bytes are preserved, since DSS identifies a time-stamp by digesting them and
// re-embeds them unchanged as an unsigned attribute.
type TimeStampToken struct {
	cms     *CMS
	tstInfo *TSTInfo
}

// ParseTimeStampToken decodes a bare TimeStampToken, i.e. a ContentInfo.
func ParseTimeStampToken(input []byte) (*TimeStampToken, error) {
	cms, err := ParseCMS(input)
	if err != nil {
		return nil, err
	}
	return TimeStampTokenFromCMS(cms)
}

// TimeStampTokenFromElement decodes an already parsed ContentInfo as a TimeStampToken.
func TimeStampTokenFromElement(element *asn1ber.Element) (*TimeStampToken, error) {
	contentInfo, err := ContentInfoFromElement(element)
	if err != nil {
		return nil, err
	}
	cms, err := CMSFromContentInfo(contentInfo)
	if err != nil {
		return nil, err
	}
	return TimeStampTokenFromCMS(cms)
}

// TimeStampTokenFromCMS reads the TSTInfo out of an already parsed CMS. Port of
// TimeStampToken(ContentInfo), including its check of the encapsulated content type.
func TimeStampTokenFromCMS(cms *CMS) (*TimeStampToken, error) {
	if !cms.SignedContentType().Equal(OIDCTTSTInfo) {
		return nil, fmt.Errorf("cmscore: the token encapsulates %s, not id-ct-TSTInfo", cms.SignedContentType())
	}
	content := cms.SignedContent()
	if len(content) == 0 {
		return nil, errors.New("cmscore: the token carries no TSTInfo")
	}
	tstInfo, err := ParseTSTInfo(content)
	if err != nil {
		return nil, err
	}
	return &TimeStampToken{cms: cms, tstInfo: tstInfo}, nil
}

// CMS returns the token's CMS structure.
func (t *TimeStampToken) CMS() *CMS { return t.cms }

// SignedData returns the token's SignedData.
func (t *TimeStampToken) SignedData() *SignedData { return t.cms.SignedData() }

// TSTInfo returns the encapsulated TSTInfo.
func (t *TimeStampToken) TSTInfo() *TSTInfo { return t.tstInfo }

// Encoded returns the token's original encoding, the bytes DSS digests to identify it.
func (t *TimeStampToken) Encoded() []byte { return t.cms.Encoded() }

// MessageImprint is
//
//	MessageImprint ::= SEQUENCE {
//	    hashAlgorithm AlgorithmIdentifier,
//	    hashedMessage OCTET STRING }
type MessageImprint struct {
	// HashAlgorithm is the digest algorithm the TSA was asked to time-stamp under.
	HashAlgorithm *asn1ber.AlgorithmIdentifier
	// HashedMessage is the digest of the time-stamped data.
	HashedMessage []byte

	// element is the parsed MessageImprint.
	element *asn1ber.Element
}

// Element returns the parsed MessageImprint.
func (m *MessageImprint) Element() *asn1ber.Element { return m.element }

// Encoded returns the MessageImprint's original encoding.
func (m *MessageImprint) Encoded() []byte { return m.element.Encoded() }

// messageImprintFromElement decodes an already parsed MessageImprint.
func messageImprintFromElement(element *asn1ber.Element) (*MessageImprint, error) {
	children, err := expectSequence(element, "MessageImprint", 2)
	if err != nil {
		return nil, err
	}
	hashAlgorithm, err := asn1ber.AlgorithmIdentifierFromElement(children[0])
	if err != nil {
		return nil, wrapField("cmscore: MessageImprint.hashAlgorithm", err)
	}
	hashedMessage, err := expectOctetString(children[1], "MessageImprint.hashedMessage")
	if err != nil {
		return nil, err
	}
	return &MessageImprint{HashAlgorithm: hashAlgorithm, HashedMessage: hashedMessage, element: element}, nil
}

// Accuracy is
//
//	Accuracy ::= SEQUENCE {
//	    seconds INTEGER           OPTIONAL,
//	    millis  [0] INTEGER (1..999) OPTIONAL,
//	    micros  [1] INTEGER (1..999) OPTIONAL }
//
// A nil field means the component was absent, which RFC 3161 distinguishes from zero.
type Accuracy struct {
	// Seconds is the seconds component, nil when absent.
	Seconds *int
	// Millis is the millis component, nil when absent.
	Millis *int
	// Micros is the micros component, nil when absent.
	Micros *int
}

// Duration returns the accuracy as a duration, absent components counting as zero.
func (a *Accuracy) Duration() time.Duration {
	if a == nil {
		return 0
	}
	total := time.Duration(0)
	if a.Seconds != nil {
		total += time.Duration(*a.Seconds) * time.Second
	}
	if a.Millis != nil {
		total += time.Duration(*a.Millis) * time.Millisecond
	}
	if a.Micros != nil {
		total += time.Duration(*a.Micros) * time.Microsecond
	}
	return total
}

// accuracyFromElement decodes an already parsed Accuracy.
func accuracyFromElement(element *asn1ber.Element) (*Accuracy, error) {
	if !element.IsUniversal(asn1ber.TagSequence) || !element.IsConstructed() {
		return nil, errors.New("cmscore: TSTInfo.accuracy is not a SEQUENCE")
	}
	accuracy := &Accuracy{}
	for _, child := range element.Children() {
		switch {
		case child.IsUniversal(asn1ber.TagInteger):
			// DIVERGENCE, deliberate: BouncyCastle's Accuracy keeps seconds as an unrestricted
			// ASN1Integer, so a value beyond 2^63 parses there. Accuracy.Seconds is an int
			// here, and clamping such a value to 0 (what int(Int64()) does) would report a
			// wrong accuracy without a word, so it is refused instead. No DSS code reads the
			// accuracy, and no TSA writes seconds anywhere near that size.
			value, err := expectSmallInteger(child, "Accuracy.seconds")
			if err != nil {
				return nil, err
			}
			accuracy.Seconds = &value
		case child.IsContextSpecific(0):
			value, err := accuracySubSecond(child, "Accuracy.millis")
			if err != nil {
				return nil, err
			}
			accuracy.Millis = &value
		case child.IsContextSpecific(1):
			value, err := accuracySubSecond(child, "Accuracy.micros")
			if err != nil {
				return nil, err
			}
			accuracy.Micros = &value
		default:
			return nil, errors.New("cmscore: unexpected Accuracy component")
		}
	}
	return accuracy, nil
}

// accuracySubSecond decodes the [0] millis or [1] micros component of an Accuracy, an IMPLICIT
// INTEGER (1..999). BouncyCastle's Accuracy(ASN1Sequence) refuses a value outside that range
// (IllegalArgumentException "Invalid millis field : not in (1..999)", ArithmeticException for
// one that does not even fit an int), which fails the parse of the whole TimeStampToken there.
func accuracySubSecond(element *asn1ber.Element, name string) (int, error) {
	if element.IsConstructed() {
		return 0, fmt.Errorf("cmscore: %s is not an INTEGER", name)
	}
	value := element.Integer()
	if !value.IsInt64() || value.Int64() < 1 || value.Int64() > 999 {
		return 0, fmt.Errorf("cmscore: %s is not in (1..999)", name)
	}
	return int(value.Int64()), nil
}

// TSTInfo is
//
//	TSTInfo ::= SEQUENCE {
//	    version        INTEGER { v1(1) },
//	    policy         TSAPolicyId,
//	    messageImprint MessageImprint,
//	    serialNumber   INTEGER,
//	    genTime        GeneralizedTime,
//	    accuracy       Accuracy                OPTIONAL,
//	    ordering       BOOLEAN DEFAULT FALSE,
//	    nonce          INTEGER                 OPTIONAL,
//	    tsa            [0] EXPLICIT GeneralName OPTIONAL,
//	    extensions     [1] IMPLICIT Extensions  OPTIONAL }
type TSTInfo struct {
	// Version is the TSTInfo version, 1 for RFC 3161.
	Version int
	// Policy is the TSA policy the token was issued under.
	Policy asn1.ObjectIdentifier
	// MessageImprint is the digest the token covers.
	MessageImprint *MessageImprint
	// SerialNumber is the TSA's serial number for the token.
	SerialNumber *big.Int
	// GenTime is the instant the token was created.
	GenTime time.Time
	// GenTimeString is the genTime as it was written, kept because its number of fractional
	// digits states the precision the TSA claimed.
	GenTimeString string
	// Accuracy is the accuracy field, nil when absent.
	Accuracy *Accuracy
	// Ordering is the ordering field, false by default.
	Ordering bool
	// Nonce is the nonce echoed back from the request, nil when absent.
	Nonce *big.Int
	// TSA is the tsa field, nil when absent.
	TSA *asn1ber.GeneralName
	// Extensions holds the members of the extensions field as parsed Extension SEQUENCEs,
	// nil when the field is absent. They are left uninterpreted here.
	Extensions []*asn1ber.Element

	// element is the parsed TSTInfo.
	element *asn1ber.Element
}

// ParseTSTInfo decodes a TSTInfo, i.e. the content of a token's eContent OCTET STRING.
func ParseTSTInfo(input []byte) (*TSTInfo, error) {
	element, err := parseOne(input, "TSTInfo")
	if err != nil {
		return nil, err
	}
	return TSTInfoFromElement(element)
}

// TSTInfoFromElement decodes an already parsed TSTInfo.
func TSTInfoFromElement(element *asn1ber.Element) (*TSTInfo, error) {
	children, err := expectSequence(element, "TSTInfo", 5)
	if err != nil {
		return nil, err
	}
	version, err := expectSmallInteger(children[0], "TSTInfo.version")
	if err != nil {
		return nil, err
	}
	policy, err := expectOID(children[1], "TSTInfo.policy")
	if err != nil {
		return nil, err
	}
	messageImprint, err := messageImprintFromElement(children[2])
	if err != nil {
		return nil, err
	}
	serialNumber, err := expectInteger(children[3], "TSTInfo.serialNumber")
	if err != nil {
		return nil, err
	}
	if !children[4].IsUniversal(asn1ber.TagGeneralizedTime) {
		return nil, errors.New("cmscore: TSTInfo.genTime is not a GeneralizedTime")
	}
	genTimeString := string(children[4].Content())
	// The fractional seconds a TSA writes are significant - "the accuracy the TSA claims" -
	// and ParseGeneralizedTime keeps them to the millisecond, as ASN1GeneralizedTime#getDate
	// does and as every DSS date comparison expects.
	genTime, err := asn1ber.ParseGeneralizedTime(genTimeString)
	if err != nil {
		return nil, wrapField("cmscore: TSTInfo.genTime", err)
	}

	info := &TSTInfo{
		Version:        version,
		Policy:         policy,
		MessageImprint: messageImprint,
		SerialNumber:   serialNumber,
		GenTime:        genTime,
		GenTimeString:  genTimeString,
		element:        element,
	}

	for _, child := range children[5:] {
		switch {
		case child.IsUniversal(asn1ber.TagSequence):
			accuracy, err := accuracyFromElement(child)
			if err != nil {
				return nil, err
			}
			info.Accuracy = accuracy
		case child.IsUniversal(asn1ber.TagBoolean):
			content := child.Content()
			info.Ordering = len(content) == 1 && content[0] != 0x00
		case child.IsUniversal(asn1ber.TagInteger):
			info.Nonce = child.Integer()
		case child.IsContextSpecific(0):
			tagged, err := explicitContent(child, "TSTInfo.tsa")
			if err != nil {
				return nil, err
			}
			generalName, err := ParseGeneralName(tagged)
			if err != nil {
				return nil, err
			}
			info.TSA = generalName
		case child.IsContextSpecific(1):
			if !child.IsConstructed() {
				return nil, errors.New("cmscore: TSTInfo.extensions is not constructed")
			}
			info.Extensions = child.Children()
		default:
			return nil, errors.New("cmscore: unexpected TSTInfo component")
		}
	}
	return info, nil
}

// Element returns the parsed TSTInfo.
func (t *TSTInfo) Element() *asn1ber.Element { return t.element }

// Encoded returns the TSTInfo's original encoding.
func (t *TSTInfo) Encoded() []byte { return t.element.Encoded() }
