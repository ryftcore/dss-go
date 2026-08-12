// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/x509/revocation/ocsp/OCSPResponseBinary.java (DSS 6.5.RC1).
package spi

import (
	"encoding/asn1"

	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/model/x509/revocation"
)

// OCSPResponseBinary is the binary of an OCSP response token. The identifier digests the
// encoding of the OCSPResponse wrapping the basic response, so the wrapped BasicOCSPResp
// keeps the bytes it was decoded from rather than being re-encoded.
type OCSPResponseBinary struct {
	model.EncapsulatedRevocationTokenIdentifier[revocation.OCSP]

	// basicOCSPResp is the basic OCSP response the binaries encapsulate.
	basicOCSPResp *BasicOCSPResp
	// asn1ObjectIdentifier specifies the origin of the OCSP response from the
	// SignedData.crls element. NOTE: used in CAdES only.
	asn1ObjectIdentifier asn1.ObjectIdentifier
}

// OCSPResponseBinaryBuild builds the binary of the given basic OCSP response, wrapping it in
// a successful OCSPResponse first. Port of the static build(BasicOCSPResp).
func OCSPResponseBinaryBuild(basicOCSPResp *BasicOCSPResp) (*OCSPResponseBinary, error) {
	ocspRespBinary, err := DSSRevocationUtilsEncodedFromBasicResp(basicOCSPResp)
	if err != nil {
		return nil, err
	}
	return NewOCSPResponseBinary(basicOCSPResp, ocspRespBinary), nil
}

// NewOCSPResponseBinary pairs a basic OCSP response with the binaries the identifier is
// computed over. Port of the package-private OCSPResponseBinary(BasicOCSPResp, byte[])
// constructor, which the Go port exports because the CMS phase builds binaries itself.
func NewOCSPResponseBinary(basicOCSPResp *BasicOCSPResp, encoded []byte) *OCSPResponseBinary {
	return &OCSPResponseBinary{
		EncapsulatedRevocationTokenIdentifier: *model.NewEncapsulatedRevocationTokenIdentifierWithClassName[revocation.OCSP](
			"OCSPResponseBinary", encoded),
		basicOCSPResp: basicOCSPResp,
	}
}

// BasicOCSPResp returns the encapsulated basic OCSP response.
// Port of getBasicOCSPResp().
func (b *OCSPResponseBinary) BasicOCSPResp() *BasicOCSPResp {
	return b.basicOCSPResp
}

// BasicOCSPRespContent returns the encoding of the basic OCSP response, i.e. the binaries of
// the BasicOCSPResponse without the OCSPResponse wrapper. Port of getBasicOCSPRespContent().
//
// Java catches the IOException BasicOCSPResp#getEncoded() may raise and returns an empty
// array; the Go response hands back the bytes it was decoded from and cannot fail.
func (b *OCSPResponseBinary) BasicOCSPRespContent() []byte {
	return b.basicOCSPResp.Encoded()
}

// ASN1ObjectIdentifier returns the origin of the response within SignedData.crls, nil when
// unset. Port of getAsn1ObjectIdentifier().
func (b *OCSPResponseBinary) ASN1ObjectIdentifier() asn1.ObjectIdentifier {
	return b.asn1ObjectIdentifier
}

// SetASN1ObjectIdentifier sets the origin of the response within SignedData.crls.
// Port of setAsn1ObjectIdentifier(ASN1ObjectIdentifier).
func (b *OCSPResponseBinary) SetASN1ObjectIdentifier(asn1ObjectIdentifier asn1.ObjectIdentifier) {
	b.asn1ObjectIdentifier = asn1ObjectIdentifier
}

// compile-time assertion: an OCSPResponseBinary is an encapsulated revocation identifier.
var _ EncapsulatedRevocationTokenIdentifier[revocation.OCSP] = (*OCSPResponseBinary)(nil)
