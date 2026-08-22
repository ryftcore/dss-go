// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/x509/revocation/crl/CRLRef.java (DSS 6.5.RC1).
//
// The CrlValidatedID constructor consumes the ESF structure of the same name. Upstream that
// structure comes from org.bouncycastle.asn1.esf; Go has no equivalent, so CrlValidatedID
// and CrlIdentifier are defined here, next to their only consumer. Their decoding keeps the
// original DER of the crlissuer Name, since the X500Principal it becomes must compare and
// serialise byte-identically to the reference it was read from.
package spi

import (
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ryftcore/dss-go/dss/internal/asn1ber"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/model/x509/revocation"
)

// crlRefDateFormat reproduces java.util.Date#toString(), which CRLRef's toString()
// interpolates for the CRL issue time.
const crlRefDateFormat = "Mon Jan 02 15:04:05 MST 2006"

// CrlIdentifier is the ESF structure
//
//	CrlIdentifier ::= SEQUENCE {
//	    crlissuer     Name,
//	    crlIssuedTime UTCTime,
//	    crlNumber     INTEGER OPTIONAL }
//
// replacing org.bouncycastle.asn1.esf.CrlIdentifier.
type CrlIdentifier struct {
	// CrlIssuer holds the original DER encoding of the crlissuer Name.
	CrlIssuer []byte
	// CrlIssuedTime is the decoded crlIssuedTime UTCTime.
	CrlIssuedTime time.Time
	// CrlNumber is the crlNumber, nil when absent.
	CrlNumber *big.Int
}

// CrlValidatedID is the ESF structure
//
//	CrlValidatedID ::= SEQUENCE {
//	    crlHash       OtherHash,
//	    crlIdentifier CrlIdentifier OPTIONAL }
//
// replacing org.bouncycastle.asn1.esf.CrlValidatedID.
type CrlValidatedID struct {
	// CrlHash is the digest of the referenced CRL.
	CrlHash *OtherHash
	// CrlIdentifier identifies the referenced CRL, nil when absent.
	CrlIdentifier *CrlIdentifier
}

// ParseCrlValidatedID decodes a CrlValidatedID from its DER encoding.
// Port of CrlValidatedID.getInstance(Object).
func ParseCrlValidatedID(der []byte) (*CrlValidatedID, error) {
	element, rest, err := asn1ber.Parse(der)
	if err != nil {
		return nil, err
	}
	if len(rest) != 0 {
		return nil, errors.New("extra data found after the CrlValidatedID")
	}
	if !element.IsConstructed() || len(element.Children()) == 0 || len(element.Children()) > 2 {
		return nil, errors.New("malformed CrlValidatedID")
	}
	crlHash, err := ParseOtherHash(element.Children()[0].Encoded())
	if err != nil {
		return nil, err
	}
	validatedID := &CrlValidatedID{CrlHash: crlHash}
	if len(element.Children()) == 2 {
		crlIdentifier, err := crlRefParseCrlIdentifier(element.Children()[1])
		if err != nil {
			return nil, err
		}
		validatedID.CrlIdentifier = crlIdentifier
	}
	return validatedID, nil
}

// crlRefParseCrlIdentifier decodes an already parsed CrlIdentifier.
func crlRefParseCrlIdentifier(element *asn1ber.Element) (*CrlIdentifier, error) {
	if !element.IsConstructed() || len(element.Children()) < 2 || len(element.Children()) > 3 {
		return nil, errors.New("malformed CrlIdentifier")
	}
	if !element.Children()[1].IsUniversal(asn1ber.TagUTCTime) {
		return nil, errors.New("malformed CrlIdentifier: crlIssuedTime is not a UTCTime")
	}
	crlIssuedTime, err := asn1ber.ParseUTCTime(string(element.Children()[1].Content()))
	if err != nil {
		return nil, err
	}
	identifier := &CrlIdentifier{
		CrlIssuer:     element.Children()[0].Encoded(),
		CrlIssuedTime: crlIssuedTime,
	}
	if len(element.Children()) == 3 {
		if !element.Children()[2].IsUniversal(asn1ber.TagInteger) {
			return nil, errors.New("malformed CrlIdentifier: crlNumber is not an INTEGER")
		}
		identifier.CrlNumber = element.Children()[2].Integer()
	}
	return identifier, nil
}

// CRLRef is a reference to an X509 CRL.
type CRLRef struct {
	RevocationRefBase[revocation.CRL]

	// crlIssuer is the name of the CRL issuer.
	crlIssuer *model.X500Principal
	// crlIssueTime is the time of CRL production; the zero time stands for Java's null.
	crlIssueTime time.Time
	// crlNumber is the CRL number, nil when unknown.
	crlNumber *big.Int
	// crlURI is a URI reference to the CRL (a hint).
	crlURI string
}

// NewCRLRef builds a reference carrying the digest only.
// Port of the CRLRef(Digest) constructor.
func NewCRLRef(digest model.Digest) *CRLRef {
	return crlRefNew(digest, nil, time.Time{}, nil)
}

// NewCRLRefWithIssuer builds a reference carrying the digest, the issuer and the issue time.
// Port of the CRLRef(Digest, X500Principal, Date) constructor.
func NewCRLRefWithIssuer(digest model.Digest, crlIssuer *model.X500Principal, crlIssueTime time.Time) *CRLRef {
	return crlRefNew(digest, crlIssuer, crlIssueTime, nil)
}

// NewCRLRefWithNumber builds a reference carrying the digest, the issuer, the issue time and
// the CRL number. Port of the CRLRef(Digest, X500Principal, Date, BigInteger) constructor.
func NewCRLRefWithNumber(digest model.Digest, crlIssuer *model.X500Principal, crlIssueTime time.Time,
	crlNumber *big.Int) *CRLRef {
	return crlRefNew(digest, crlIssuer, crlIssueTime, crlNumber)
}

// NewCRLRefFromCrlValidatedID builds a reference from the ESF structure a CAdES
// complete-revocation-references attribute carries.
// Port of the CRLRef(CrlValidatedID) constructor; every failure upstream wraps in
// "Unable to build CRLRef from CrlValidatedID" is returned as an error here.
func NewCRLRefFromCrlValidatedID(cmsRef *CrlValidatedID) (*CRLRef, error) {
	ref := crlRefNew(model.Digest{}, nil, time.Time{}, nil)
	if cmsRef.CrlIdentifier != nil {
		crlIssuer, err := DSSASN1UtilsToX500Principal(cmsRef.CrlIdentifier.CrlIssuer)
		if err != nil {
			return nil, model.NewDSSErrorMessageCause("Unable to build CRLRef from CrlValidatedID", err)
		}
		ref.crlIssuer = crlIssuer
		ref.crlIssueTime = cmsRef.CrlIdentifier.CrlIssuedTime
		ref.crlNumber = cmsRef.CrlIdentifier.CrlNumber
	}
	digest, err := DSSRevocationUtilsDigest(cmsRef.CrlHash)
	if err != nil {
		return nil, model.NewDSSErrorMessageCause("Unable to build CRLRef from CrlValidatedID", err)
	}
	ref.SetDigest(digest)
	return ref, nil
}

// crlRefNew assembles a CRLRef and registers it with its base so that the identifier
// dispatch of RevocationRefBase reaches this type.
func crlRefNew(digest model.Digest, crlIssuer *model.X500Principal, crlIssueTime time.Time,
	crlNumber *big.Int) *CRLRef {
	ref := &CRLRef{
		RevocationRefBase: NewRevocationRefBase[revocation.CRL](),
		crlIssuer:         crlIssuer,
		crlIssueTime:      crlIssueTime,
		crlNumber:         crlNumber,
	}
	ref.InitRevocationRef(ref)
	ref.SetDigest(digest)
	return ref
}

// CRLIssuer returns the name of the CRL issuer. Port of getCrlIssuer().
func (r *CRLRef) CRLIssuer() *model.X500Principal {
	return r.crlIssuer
}

// CRLIssueTime returns the CRL issue time, the zero time when unknown.
// Port of getCrlIssueTime().
func (r *CRLRef) CRLIssueTime() time.Time {
	return r.crlIssueTime
}

// CRLNumber returns the CRL number, nil when unknown. Port of getCrlNumber().
func (r *CRLRef) CRLNumber() *big.Int {
	return r.crlNumber
}

// CRLURI returns the URI reference to the CRL (a hint). Port of getCrlUri().
func (r *CRLRef) CRLURI() string {
	return r.crlURI
}

// SetCRLURI sets a URI reference to the CRL (a hint). Port of setCrlUri(String).
func (r *CRLRef) SetCRLURI(crlURI string) {
	r.crlURI = crlURI
}

// String renders the reference the way Java's toString() does; a missing issuer, issue time
// or CRL number is printed as "null", as Java's string concatenation does.
func (r *CRLRef) String() string {
	return "CRLRef [" +
		"crlIssuer=" + crlRefString(r.crlIssuer) +
		", crlIssueTime=" + crlRefDateString(r.crlIssueTime) +
		", crlNumber=" + crlRefNumberString(r.crlNumber) +
		"] " + r.RevocationRefBase.String()
}

// crlRefString renders an X500Principal the way Java's string concatenation does.
func crlRefString(principal *model.X500Principal) string {
	if principal == nil {
		return "null"
	}
	return principal.String()
}

// crlRefDateString renders a date the way java.util.Date#toString() does.
func crlRefDateString(date time.Time) string {
	if date.IsZero() {
		return "null"
	}
	return date.Format(crlRefDateFormat)
}

// crlRefNumberString renders a CRL number the way Java's string concatenation does.
func crlRefNumberString(number *big.Int) string {
	if number == nil {
		return "null"
	}
	return number.String()
}

// Equals reports whether both references carry the same digest, issuer, issue time and CRL
// number. Port of equals(Object).
func (r *CRLRef) Equals(other *CRLRef) bool {
	if r == other {
		return true
	}
	if other == nil {
		return false
	}
	// Port of super.equals(object), which compares the digests only; the getClass() check
	// it also performs is carried by the *CRLRef parameter type.
	if !r.Digest().Equals(other.Digest()) {
		return false
	}
	if !crlRefPrincipalsEqual(r.crlIssuer, other.crlIssuer) {
		return false
	}
	if !r.crlIssueTime.Equal(other.crlIssueTime) {
		return false
	}
	return crlRefNumbersEqual(r.crlNumber, other.crlNumber)
}

// crlRefPrincipalsEqual is Objects.equals for two X500Principals.
func crlRefPrincipalsEqual(first, second *model.X500Principal) bool {
	if first == nil || second == nil {
		return first == second
	}
	return first.Equals(second)
}

// crlRefNumbersEqual is Objects.equals for two CRL numbers.
func crlRefNumbersEqual(first, second *big.Int) bool {
	if first == nil || second == nil {
		return first == second
	}
	return first.Cmp(second) == 0
}

// compile-time assertions: a CRLRef is a revocation reference and prints itself.
var (
	_ RevocationRef[revocation.CRL] = (*CRLRef)(nil)
	_ fmt.Stringer                  = (*CRLRef)(nil)
)
