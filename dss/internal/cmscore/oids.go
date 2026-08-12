// Object identifiers of RFC 5652 and RFC 3161, replacing
// org.bouncycastle.asn1.cms.CMSObjectIdentifiers, org.bouncycastle.asn1.pkcs.PKCSObjectIdentifiers
// (the pkcs-9 attribute branch) and org.bouncycastle.asn1.ocsp.OCSPObjectIdentifiers.
//
// Go constants cannot hold slices, so these are package-level vars: treat them as immutable.
package cmscore

import "encoding/asn1"

// CMS content types, i.e. the { iso(1) member-body(2) us(840) rsadsi(113549) pkcs(1) pkcs7(7) }
// branch of RFC 5652 clause 4 and the id-ct branch of RFC 5911.
var (
	// OIDData is id-data, the content type of an arbitrary octet string.
	OIDData = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 7, 1}
	// OIDSignedData is id-signedData, the content type this package parses and builds.
	OIDSignedData = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 7, 2}
	// OIDEnvelopedData is id-envelopedData.
	OIDEnvelopedData = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 7, 3}
	// OIDSignedAndEnvelopedData is id-signedAndEnvelopedData (PKCS#7 only).
	OIDSignedAndEnvelopedData = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 7, 4}
	// OIDDigestedData is id-digestedData.
	OIDDigestedData = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 7, 5}
	// OIDEncryptedData is id-encryptedData.
	OIDEncryptedData = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 7, 6}
	// OIDCTTSTInfo is id-ct-TSTInfo, the eContentType of an RFC 3161 TimeStampToken.
	OIDCTTSTInfo = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 9, 16, 1, 4}
)

// Attribute types of the pkcs-9 branch that RFC 5652 clause 11 defines for SignerInfo.
var (
	// OIDContentType is id-contentType, a mandatory signed attribute.
	OIDContentType = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 9, 3}
	// OIDMessageDigest is id-messageDigest, a mandatory signed attribute.
	OIDMessageDigest = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 9, 4}
	// OIDSigningTime is id-signingTime.
	OIDSigningTime = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 9, 5}
	// OIDCounterSignature is id-countersignature, an unsigned attribute.
	OIDCounterSignature = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 9, 6}
)

// Revocation information formats carried by an OtherRevocationInfoFormat.
var (
	// OIDRIOCSPResponse is id-ri-ocsp-response OBJECT IDENTIFIER ::= { id-ri 2 }, where
	// id-ri ::= { iso(1) identified-organization(3) dod(6) internet(1) security(5)
	// mechanisms(5) pkix(7) ri(16) }. Its otherRevInfo is an RFC 6960 OCSPResponse.
	OIDRIOCSPResponse = asn1.ObjectIdentifier{1, 3, 6, 1, 5, 5, 7, 16, 2}
	// OIDPKIXOCSPBasic is id-pkix-ocsp-basic. When it appears as an otherRevInfoFormat the
	// otherRevInfo is a bare BasicOCSPResponse rather than a complete OCSPResponse; DSS
	// reads both (see CMSOCSPSource).
	OIDPKIXOCSPBasic = asn1.ObjectIdentifier{1, 3, 6, 1, 5, 5, 7, 48, 1, 1}
)
