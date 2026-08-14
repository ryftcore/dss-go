// SignedData of RFC 5652 clause 5.1, replacing org.bouncycastle.asn1.cms.SignedData and the
// read side of org.bouncycastle.cms.CMSSignedData.
package cmscore

import (
	"errors"
	"fmt"

	"github.com/utain/esig/dss/internal/asn1ber"
)

// SignedData is
//
//	SignedData ::= SEQUENCE {
//	    version          CMSVersion,
//	    digestAlgorithms DigestAlgorithmIdentifiers,
//	    encapContentInfo EncapsulatedContentInfo,
//	    certificates     [0] IMPLICIT CertificateSet OPTIONAL,
//	    crls             [1] IMPLICIT RevocationInfoChoices OPTIONAL,
//	    signerInfos      SignerInfos }
//
// A parsed SignedData keeps the encoding of each of the four fields DSS re-serialises when it
// computes an archive time-stamp message imprint (digestAlgorithms, encapContentInfo,
// certificates, crls and signerInfos - see CMSUtils#writeSignedData*Encoded), so that those
// bytes never have to be reconstructed.
type SignedData struct {
	// Version is the CMSVersion, see ComputeSignedDataVersion.
	Version int
	// DigestAlgorithms holds the members of the digestAlgorithms SET.
	DigestAlgorithms []*asn1ber.AlgorithmIdentifier
	// EncapContentInfo is the signed content and its type.
	EncapContentInfo *EncapsulatedContentInfo
	// Certificates is the certificates field, nil when it is absent.
	Certificates *CertificateSet
	// CRLs is the crls field, nil when it is absent.
	CRLs *RevocationInfoChoices
	// SignerInfos holds the members of the signerInfos SET.
	SignerInfos []*SignerInfo

	// element is the parsed SignedData, nil when it was built.
	element *asn1ber.Element
	// digestAlgorithmsElement and signerInfosElement are the parsed SET fields.
	digestAlgorithmsElement *asn1ber.Element
	signerInfosElement      *asn1ber.Element
}

// ParseSignedData decodes a bare SignedData, i.e. one not wrapped in a ContentInfo. Use
// ParseCMS for the usual case.
func ParseSignedData(input []byte) (*SignedData, error) {
	element, err := parseOne(input, "SignedData")
	if err != nil {
		return nil, err
	}
	return SignedDataFromElement(element)
}

// SignedDataFromElement decodes an already parsed SignedData.
func SignedDataFromElement(element *asn1ber.Element) (*SignedData, error) {
	children, err := expectSequence(element, "SignedData", 4)
	if err != nil {
		return nil, err
	}
	version, err := expectSmallInteger(children[0], "SignedData.version")
	if err != nil {
		return nil, err
	}
	digestAlgorithmElements, err := expectSet(children[1], "SignedData.digestAlgorithms")
	if err != nil {
		return nil, err
	}
	signedData := &SignedData{
		Version:                 version,
		element:                 element,
		digestAlgorithmsElement: children[1],
	}
	for _, child := range digestAlgorithmElements {
		algorithm, err := asn1ber.AlgorithmIdentifierFromElement(child)
		if err != nil {
			return nil, wrapField("cmscore: SignedData.digestAlgorithms", err)
		}
		signedData.DigestAlgorithms = append(signedData.DigestAlgorithms, algorithm)
	}

	encapContentInfo, err := EncapsulatedContentInfoFromElement(children[2])
	if err != nil {
		return nil, err
	}
	signedData.EncapContentInfo = encapContentInfo

	index := 3
	if children[index].IsContextSpecific(0) {
		if !children[index].IsConstructed() {
			return nil, errors.New("cmscore: SignedData.certificates is not constructed")
		}
		certificates, err := certificateSetFromElement(children[index])
		if err != nil {
			return nil, err
		}
		signedData.Certificates = certificates
		index++
	}
	if index < len(children) && children[index].IsContextSpecific(1) {
		if !children[index].IsConstructed() {
			return nil, errors.New("cmscore: SignedData.crls is not constructed")
		}
		crls, err := revocationInfoChoicesFromElement(children[index])
		if err != nil {
			return nil, err
		}
		signedData.CRLs = crls
		index++
	}
	if index >= len(children) {
		return nil, errors.New("cmscore: SignedData.signerInfos is missing")
	}

	signerInfoElements, err := expectSet(children[index], "SignedData.signerInfos")
	if err != nil {
		return nil, err
	}
	signedData.signerInfosElement = children[index]
	for _, child := range signerInfoElements {
		signerInfo, err := SignerInfoFromElement(child)
		if err != nil {
			return nil, err
		}
		signedData.SignerInfos = append(signedData.SignerInfos, signerInfo)
	}
	index++
	if index != len(children) {
		return nil, fmt.Errorf("cmscore: SignedData holds %d unexpected trailing components", len(children)-index)
	}
	return signedData, nil
}

// Element returns the parsed SignedData, nil when it was built.
func (s *SignedData) Element() *asn1ber.Element { return s.element }

// Encoded returns the SignedData's original encoding, and its DER encoding when it was built
// rather than parsed.
func (s *SignedData) Encoded() []byte {
	if s.element != nil {
		return s.element.Encoded()
	}
	return s.DER()
}

// DigestAlgorithmsElement returns the parsed digestAlgorithms SET, nil when the SignedData was
// built.
func (s *SignedData) DigestAlgorithmsElement() *asn1ber.Element { return s.digestAlgorithmsElement }

// SignerInfosElement returns the parsed signerInfos SET, nil when the SignedData was built.
func (s *SignedData) SignerInfosElement() *asn1ber.Element { return s.signerInfosElement }

// CertificatesElement returns the parsed certificates field, nil when it is absent or the
// SignedData was built.
func (s *SignedData) CertificatesElement() *asn1ber.Element { return s.Certificates.Element() }

// CRLsElement returns the parsed crls field, nil when it is absent or the SignedData was built.
func (s *SignedData) CRLsElement() *asn1ber.Element { return s.CRLs.Element() }

// IsDetached reports whether the signature is detached, i.e. encapContentInfo.eContent is
// absent. Port of CMS#isDetachedSignature.
func (s *SignedData) IsDetached() bool { return s.EncapContentInfo.IsDetached() }

// CertificateDERs returns the encoding of every plain X.509 certificate of the certificates
// field, nil when the field is absent.
func (s *SignedData) CertificateDERs() [][]byte { return s.Certificates.Certificates() }

// DER returns the DER encoding of the SignedData: definite lengths throughout, a segmented
// eContent collapsed into one primitive OCTET STRING, and every SET OF ordered.
func (s *SignedData) DER() []byte {
	body := encodeInt(s.Version)

	algorithms := make([][]byte, len(s.DigestAlgorithms))
	for index, algorithm := range s.DigestAlgorithms {
		algorithms[index] = algorithm.DER()
	}
	body = append(body, derSetOf([]byte{setIdentifier}, algorithms)...)
	body = append(body, s.EncapContentInfo.DER()...)
	if s.Certificates != nil {
		body = append(body, s.Certificates.DER()...)
	}
	if s.CRLs != nil {
		body = append(body, s.CRLs.DER()...)
	}

	signerInfos := make([][]byte, len(s.SignerInfos))
	for index, signerInfo := range s.SignerInfos {
		signerInfos[index] = signerInfo.DER()
	}
	body = append(body, derSetOf([]byte{setIdentifier}, signerInfos)...)
	return asn1ber.WriteSequence(body)
}

// ContentInfoDER returns the DER encoding of the SignedData wrapped in its id-signedData
// ContentInfo, i.e. a complete CMS document.
func (s *SignedData) ContentInfoDER() []byte {
	return NewContentInfo(OIDSignedData, s.DER()).DER()
}
