// CMSVersion computation of RFC 5652 clauses 5.1 and 5.3, replacing the private
// CMSSignedHelper/CMSSignedGenerator logic BouncyCastle applies when it generates a SignedData.
package cmscore

import "encoding/asn1"

// CMS version numbers a SignedData or a SignerInfo can carry.
const (
	// CMSVersion1 is the base version.
	CMSVersion1 = 1
	// CMSVersion3 is required by a subjectKeyIdentifier SignerIdentifier, a version 1
	// attribute certificate, or an eContentType other than id-data.
	CMSVersion3 = 3
	// CMSVersion4 is required by a version 2 attribute certificate.
	CMSVersion4 = 4
	// CMSVersion5 is required by an OtherCertificateFormat or an OtherRevocationInfoFormat.
	CMSVersion5 = 5
)

// ComputeSignedDataVersion applies the rule of RFC 5652 clause 5.1 verbatim:
//
//	IF ((certificates is present) AND
//	    (any certificates with a type of other are present)) OR
//	   ((crls is present) AND
//	    (any crls with a type of other are present))
//	THEN version MUST be 5
//	ELSE
//	   IF (certificates is present) AND
//	      (any version 2 attribute certificates are present)
//	   THEN version MUST be 4
//	   ELSE
//	      IF ((certificates is present) AND
//	          (any version 1 attribute certificates are present)) OR
//	         (any SignerInfo structures are version 3) OR
//	         (encapContentInfo eContentType is other than id-data)
//	      THEN version MUST be 3
//	      ELSE version MUST be 1
//
// A nil certificates or crls argument means the field is absent.
func ComputeSignedDataVersion(eContentType asn1.ObjectIdentifier, certificates *CertificateSet,
	crls *RevocationInfoChoices, signerInfos []*SignerInfo) int {
	if certificates.hasChoice(CertificateChoiceOther) || crls.HasOtherFormat() {
		return CMSVersion5
	}
	if certificates.hasChoice(CertificateChoiceV2AttrCert) {
		return CMSVersion4
	}
	if certificates.hasChoice(CertificateChoiceV1AttrCert) {
		return CMSVersion3
	}
	for _, signerInfo := range signerInfos {
		if signerInfo.Version == CMSVersion3 {
			return CMSVersion3
		}
	}
	if eContentType != nil && !eContentType.Equal(OIDData) {
		return CMSVersion3
	}
	return CMSVersion1
}

// ComputeSignerInfoVersion applies the rule of RFC 5652 clause 5.3:
//
//	IF (sid is issuerAndSerialNumber) THEN version MUST be 1
//	ELSE version MUST be 3
func ComputeSignerInfoVersion(sid *SignerIdentifier) int {
	if sid != nil && sid.IsSubjectKeyIdentifier() {
		return CMSVersion3
	}
	return CMSVersion1
}

// hasChoice reports whether the set holds a member of the given alternative. A nil set - i.e.
// an absent certificates field - holds none.
func (c *CertificateSet) hasChoice(tagNo int) bool {
	if c == nil {
		return false
	}
	for _, choice := range c.Choices {
		if choice.TagNo == tagNo {
			return true
		}
	}
	return false
}
