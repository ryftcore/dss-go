// Ported from
// dss-cms/src/main/java/eu/europa/esig/dss/cms/CMSSignedAttributeTableGenerator.java
// (DSS 6.5.RC1).
//
// Java implements org.bouncycastle.cms.CMSAttributeTableGenerator, whose getAttributes(Map)
// BouncyCastle's SignerInfoGenerator.generate(contentType) calls once it knows the real
// content-type and message-digest; the Go port folds that deferred call directly into
// SignerInfoGenerator.Generate (signer_info_generator.go), which is the only caller. This file
// keeps just the attribute-assembly logic, i.e. Java's createStandardAttributeTable, since the
// CMSAttributeTableGenerator interface itself and the Hashtable-vs-parameters-Map plumbing
// around it have no work left to do once "invoked with the map of five constants" and "return
// an AttributeTable" collapse into ordinary Go arguments and a cmscore.Attributes return.
//
// "This class replicates a org.bouncycastle.cms.DefaultAuthenticatedAttributeTableGenerator,
// but without the signing-time attribute, that should be provided externally. The class is
// used on both CMS for CAdES and CMS for PAdES generations." (Java doc, kept for context.)
package cms

import (
	"encoding/asn1"

	"github.com/utain/esig/dss/internal/asn1ber"
	"github.com/utain/esig/dss/internal/cmscore"
)

// OID_id_aa_cmsAlgorithmProtect is id-aa-cmsAlgorithmProtect OBJECT IDENTIFIER ::= {iso(1)
// member-body(2) us(840) rsadsi(113549) pkcs(1) pkcs-9(9) 52}, RFC 6211 - BouncyCastle's
// CMSAttributes.cmsAlgorithmProtect.
var OID_id_aa_cmsAlgorithmProtect = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 9, 52}

// cmsSignedAttributeTableGenerate creates a standard attribute table: normally including
// contentType, messageDigest and cms-algorithm-protection. Entries already present in
// initialAttributes for one of those three types override the generated one, exactly as
// providing an AttributeTable to CMSSignedAttributeTableGenerator's constructor does upstream.
//
// Port of #createStandardAttributeTable(Map), specialised to the five parameters DSS's own
// call sites ever pass (CMSAttributeTableGenerator.CONTENT_TYPE, DIGEST,
// DIGEST_ALGORITHM_IDENTIFIER, SIGNATURE_ALGORITHM_IDENTIFIER): a nil contentType is Java's
// "contentType will be null if we're trying to generate a counter signature", which omits the
// attribute instead of adding one holding a null OID.
func cmsSignedAttributeTableGenerate(initialAttributes cmscore.Attributes, contentType asn1.ObjectIdentifier,
	messageDigest []byte, digestAlgorithmIdentifier, signatureAlgorithmIdentifier *asn1ber.AlgorithmIdentifier) cmscore.Attributes {
	table := make(cmscore.Attributes, len(initialAttributes))
	copy(table, initialAttributes)

	if table.Get(cmscore.OIDContentType) == nil && contentType != nil {
		table = append(table, cmscore.NewAttribute(cmscore.OIDContentType, asn1ber.EncodeOID(contentType)))
	}
	if table.Get(cmscore.OIDMessageDigest) == nil {
		table = append(table, cmscore.NewAttribute(cmscore.OIDMessageDigest, asn1ber.WriteTLV(asn1ber.TagOctetString, messageDigest)))
	}
	if table.Get(OID_id_aa_cmsAlgorithmProtect) == nil {
		table = append(table, cmscore.NewAttribute(OID_id_aa_cmsAlgorithmProtect,
			cmsAlgorithmProtectionDER(digestAlgorithmIdentifier, signatureAlgorithmIdentifier)))
	}
	return table
}

// cmsAlgorithmProtectionDER returns the DER encoding of a CMSAlgorithmProtection value for the
// "signatureAlgorithm present" alternative (RFC 6211):
//
//	CMSAlgorithmProtection ::= SEQUENCE {
//	    digestAlgorithm    DigestAlgorithmIdentifier,
//	    signatureAlgorithm [1] SignatureAlgorithmIdentifier OPTIONAL,
//	    macAlgorithm       [2] MessageAuthenticationCodeAlgorithm OPTIONAL }
//
// replacing org.bouncycastle.asn1.cms.CMSAlgorithmProtection(AlgorithmIdentifier, int,
// AlgorithmIdentifier) with CMSAlgorithmProtection.SIGNATURE. The [1] tag is IMPLICIT (it
// replaces signatureAlgorithm's own SEQUENCE tag, the way SignerIdentifier's [0]
// subjectKeyIdentifier alternative does) - verified against a real BouncyCastle 1.84
// CMSAlgorithmProtection encoding (see testdata/gen), not derived from the ASN.1 module alone,
// since CMS modules mix implicit and explicit tagging by field.
func cmsAlgorithmProtectionDER(digestAlgorithmIdentifier, signatureAlgorithmIdentifier *asn1ber.AlgorithmIdentifier) []byte {
	body := digestAlgorithmIdentifier.DER()
	signatureAlgorithmBody := asn1ber.EncodeOID(signatureAlgorithmIdentifier.Algorithm)
	if signatureAlgorithmIdentifier.Parameters != nil {
		signatureAlgorithmBody = append(signatureAlgorithmBody, signatureAlgorithmIdentifier.Parameters...)
	}
	body = append(body, asn1ber.WriteTLV(asn1ber.ClassContextSpecific|asn1ber.Constructed|1, signatureAlgorithmBody)...)
	return asn1ber.WriteSequence(body)
}
