// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/x509/revocation/ocsp/OCSPRef.java (DSS 6.5.RC1).
//
// The OcspResponsesID constructor consumes the ESF structure of the same name. Upstream it
// comes from org.bouncycastle.asn1.esf; Go has no equivalent, so OcspResponsesID and
// OcspIdentifier are defined here, next to their only consumer. Their decoding reuses the
// RFC 6960 ResponderID decoder of DSSRevocationUtils, which keeps the original DER of the
// responder's Name.
package spi

import (
	"errors"
	"fmt"
	"time"

	"github.com/ryftcore/dss-go/dss/internal/asn1ber"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/model/x509/revocation"
	"github.com/ryftcore/dss-go/dss/utils"
)

// OcspIdentifier is the ESF structure
//
//	OcspIdentifier ::= SEQUENCE {
//	    ocspResponderID ResponderID,
//	    producedAt      GeneralizedTime }
//
// replacing org.bouncycastle.asn1.esf.OcspIdentifier.
type OcspIdentifier struct {
	// OcspResponderID identifies the responder of the referenced response.
	OcspResponderID *ResponderID
	// ProducedAt holds the original DER encoding of the producedAt GeneralizedTime.
	ProducedAt []byte
}

// OcspResponsesID is the ESF structure
//
//	OcspResponsesID ::= SEQUENCE {
//	    ocspIdentifier OcspIdentifier,
//	    ocspRepHash    OtherHash OPTIONAL }
//
// replacing org.bouncycastle.asn1.esf.OcspResponsesID.
type OcspResponsesID struct {
	// OcspIdentifier identifies the referenced response.
	OcspIdentifier *OcspIdentifier
	// OcspRepHash is the digest of the referenced response, nil when absent.
	OcspRepHash *OtherHash
}

// ParseOcspResponsesID decodes an OcspResponsesID from its DER encoding.
// Port of OcspResponsesID.getInstance(Object).
func ParseOcspResponsesID(der []byte) (*OcspResponsesID, error) {
	element, rest, err := asn1ber.Parse(der)
	if err != nil {
		return nil, err
	}
	if len(rest) != 0 {
		return nil, errors.New("extra data found after the OcspResponsesID")
	}
	if !element.IsConstructed() || len(element.Children()) == 0 || len(element.Children()) > 2 {
		return nil, errors.New("malformed OcspResponsesID")
	}
	identifier, err := ocspRefParseOcspIdentifier(element.Children()[0])
	if err != nil {
		return nil, err
	}
	responsesID := &OcspResponsesID{OcspIdentifier: identifier}
	if len(element.Children()) == 2 {
		hash, err := ParseOtherHash(element.Children()[1].Encoded())
		if err != nil {
			return nil, err
		}
		responsesID.OcspRepHash = hash
	}
	return responsesID, nil
}

// ocspRefParseOcspIdentifier decodes an already parsed OcspIdentifier.
func ocspRefParseOcspIdentifier(element *asn1ber.Element) (*OcspIdentifier, error) {
	if !element.IsConstructed() || len(element.Children()) != 2 {
		return nil, errors.New("malformed OcspIdentifier")
	}
	responderID, err := dssRevocationUtilsResponderIDFromElement(element.Children()[0])
	if err != nil {
		return nil, err
	}
	if !element.Children()[1].IsUniversal(asn1ber.TagGeneralizedTime) {
		return nil, errors.New("malformed OcspIdentifier: producedAt is not a GeneralizedTime")
	}
	return &OcspIdentifier{OcspResponderID: responderID, ProducedAt: element.Children()[1].Encoded()}, nil
}

// OCSPRef references an OCSP response.
type OCSPRef struct {
	RevocationRefBase[revocation.OCSP]

	// producedAt is the OCSP production time; the zero time stands for Java's null.
	producedAt time.Time
	// responderId identifies the OCSP responder.
	responderId *ResponderId
	// uri is a URI reference to the OCSP response (a hint).
	uri string
}

// NewOCSPRef builds a reference carrying the digest, the production time and the responder.
// Port of the OCSPRef(Digest, Date, ResponderId) constructor.
func NewOCSPRef(digest model.Digest, producedAt time.Time, responderId *ResponderId) *OCSPRef {
	ref := &OCSPRef{
		RevocationRefBase: NewRevocationRefBase[revocation.OCSP](),
		producedAt:        producedAt,
		responderId:       responderId,
	}
	ref.InitRevocationRef(ref)
	ref.SetDigest(digest)
	return ref
}

// NewOCSPRefFromOcspResponsesID builds a reference from the ESF structure a CAdES
// complete-revocation-references attribute carries.
// Port of the OCSPRef(OcspResponsesID) constructor; upstream lets the decoding failures
// propagate as runtime exceptions, they are returned as errors here.
func NewOCSPRefFromOcspResponsesID(ocspResponsesID *OcspResponsesID) (*OCSPRef, error) {
	digest, err := DSSRevocationUtilsDigest(ocspResponsesID.OcspRepHash)
	if err != nil {
		return nil, err
	}
	responderId, err := DSSRevocationUtilsDSSResponderID(ocspResponsesID.OcspIdentifier.OcspResponderID)
	if err != nil {
		return nil, err
	}
	producedAt := DSSASN1UtilsDate(ocspResponsesID.OcspIdentifier.ProducedAt)
	return NewOCSPRef(digest, producedAt, responderId), nil
}

// ProducedAt returns the OCSP production time, the zero time when unknown.
// Port of getProducedAt().
func (r *OCSPRef) ProducedAt() time.Time {
	return r.producedAt
}

// ResponderId returns the identity of the OCSP responder. Port of getResponderId().
func (r *OCSPRef) ResponderId() *ResponderId {
	return r.responderId
}

// URI returns the URI reference to the OCSP response (a hint). Port of getUri().
func (r *OCSPRef) URI() string {
	return r.uri
}

// SetURI sets the URI reference to the OCSP response (a hint). Port of setUri(String).
func (r *OCSPRef) SetURI(uri string) {
	r.uri = uri
}

// CreateIdentifier builds the reference's unique identifier.
// Port of the protected createIdentifier() override.
func (r *OCSPRef) CreateIdentifier() model.Identifier {
	return NewOCSPRefIdentifier(r)
}

// String renders the reference the way Java's toString() does: by responder name when the
// responder identified itself by name, by base64-encoded key hash otherwise.
func (r *OCSPRef) String() string {
	if r.responderId.X500Principal() != nil {
		return "OCSP Reference produced at [" + DSSUtilsFormatDateToRFC(r.producedAt) + "] " +
			"with Responder Name: [" + r.responderId.X500Principal().String() + "]"
	}
	return "OCSP Reference produced at [" + DSSUtilsFormatDateToRFC(r.producedAt) + "] " +
		"with Responder key 64base: [" + utils.ToBase64(r.responderId.Ski()) + "]"
}

// Equals reports whether both references carry the same digest, production time and
// responder. Port of equals(Object).
func (r *OCSPRef) Equals(other *OCSPRef) bool {
	if r == other {
		return true
	}
	if other == nil {
		return false
	}
	// Port of super.equals(obj), which compares the digests only; the getClass() check it
	// also performs is carried by the *OCSPRef parameter type.
	if !r.Digest().Equals(other.Digest()) {
		return false
	}
	if !r.producedAt.Equal(other.producedAt) {
		return false
	}
	return ocspRefResponderIdsEqual(r.responderId, other.responderId)
}

// ocspRefResponderIdsEqual is Objects.equals for two responder identities.
func ocspRefResponderIdsEqual(first, second *ResponderId) bool {
	if first == nil || second == nil {
		return first == second
	}
	return first.Equals(second)
}

// compile-time assertions: an OCSPRef is a revocation reference and prints itself.
var (
	_ RevocationRef[revocation.OCSP] = (*OCSPRef)(nil)
	_ fmt.Stringer                   = (*OCSPRef)(nil)
)
