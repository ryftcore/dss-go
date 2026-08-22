// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/OID.java (DSS 6.5.RC1).
//
// Java exposes the constants as org.bouncycastle.asn1.ASN1ObjectIdentifier; the Go port
// uses encoding/asn1.ObjectIdentifier. Because Go constants cannot hold slices these are
// package-level vars: they must be treated as immutable.
//
// The Java constants are class-scoped and lower-case, so - as PORTING.md prescribes for
// enum constants - each keeps its exact Java name behind an "OID_" prefix, which both
// exports it and keeps it collision-free inside the flattened spi package.
package spi

import "encoding/asn1"

// The constants Java derives with PKCSObjectIdentifiers.id_aa.branch(...) hang off
// id-aa OBJECT IDENTIFIER ::= { iso(1) member-body(2) us(840) rsadsi(113549) pkcs(1)
// pkcs-9(9) smime(16) 2 }, i.e. 1.2.840.113549.1.9.16.2, and are spelled out in full here.
var (
	// OIDIdAaEtsArchiveTimestampV2 is
	// id-aa-ets-archiveTimestampV2 OBJECT IDENTIFIER ::= { iso(1) member-body(2) us(840)
	// rsadsi(113549) pkcs(1) pkcs-9(9) smime(16) id-aa(2) 48}
	OIDIdAaEtsArchiveTimestampV2 = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 9, 16, 2, 48}

	// OIDIdEtsiElectronicSignatureStandardAttributes is
	// Attributes {itu-t(0) identified-organization(4) etsi(0)
	// electronic-signature-standard(1733) attributes(2)}
	OIDIdEtsiElectronicSignatureStandardAttributes = asn1.ObjectIdentifier{0, 4, 0, 1733, 2}

	// OIDIdEtsiSignerAttributes is
	// Signer attributes {itu-t(0) identified-organization(4) etsi(0) cades(19122)
	// attributes(1)} (see ETSI EN 319 122-1)
	OIDIdEtsiSignerAttributes = asn1.ObjectIdentifier{0, 4, 0, 19122, 1}

	// OIDIdEtsiSpq is
	// Signer attributes {itu-t(0) identified-organization(4) etsi(0) cades(19122)
	// id-spq (2)} (see ETSI EN 319 122-1)
	OIDIdEtsiSpq = asn1.ObjectIdentifier{0, 4, 0, 19122, 2}

	// OIDIdAaEtsMimeType is
	// id-aa-ets-mimeType OBJECT IDENTIFIER ::= { itu-t(0) identified-organization(4)
	// etsi(0) electronic-signature-standard(1733) attributes(2) 1 }
	OIDIdAaEtsMimeType = asn1.ObjectIdentifier{0, 4, 0, 1733, 2, 1}

	// OIDIdAaEtsArchiveTimestampV3 is
	// id-aa-ets-archiveTimestampV3 OBJECT IDENTIFIER ::= { itu-t(0)
	// identified-organization(4) etsi(0) electronic-signature-standard(1733)
	// attributes(2) 4 }
	OIDIdAaEtsArchiveTimestampV3 = asn1.ObjectIdentifier{0, 4, 0, 1733, 2, 4}

	// OIDIdAaATSHashIndex is
	// id-aa-ATSHashIndex OBJECT IDENTIFIER ::= { itu-t(0) identified-organization(4)
	// etsi(0) electronicsignature-standard(1733) attributes(2) 5 }
	OIDIdAaATSHashIndex = asn1.ObjectIdentifier{0, 4, 0, 1733, 2, 5}

	// OIDIdAaEtsSignerAttrV2 is
	// id-aa-ets-signerAttrV2 OBJECT IDENTIFIER ::= { itu-t(0) identified-organization(4)
	// etsi(0) cades(19122) attributes(1) 1 }
	OIDIdAaEtsSignerAttrV2 = asn1.ObjectIdentifier{0, 4, 0, 19122, 1, 1}

	// OIDIdAaEtsSigPolicyStore is
	// id-aa-ets-sigPolicyStore OBJECT IDENTIFIER ::= { itu-t(0) identified-organization(4)
	// etsi(0) cades(19122) attributes(1) 3 }
	OIDIdAaEtsSigPolicyStore = asn1.ObjectIdentifier{0, 4, 0, 19122, 1, 3}

	// OIDIdAaATSHashIndexV2 is
	// id-aa-ATSHashIndex-v2 OBJECT IDENTIFIER ::= { itu-t(0) identified-organization(4)
	// etsi(0) cades(19122) attributes(1) 4 }
	OIDIdAaATSHashIndexV2 = asn1.ObjectIdentifier{0, 4, 0, 19122, 1, 4}

	// OIDIdAaATSHashIndexV3 is
	// id-aa-ATSHashIndex-v3 OBJECT IDENTIFIER ::= { itu-t(0) identified-organization(4)
	// etsi(0) cades(19122) attributes(1) 5 }
	OIDIdAaATSHashIndexV3 = asn1.ObjectIdentifier{0, 4, 0, 19122, 1, 5}

	// OIDIdSpDocSpecification is
	// id-spq-ets-docspec OBJECT IDENTIFIER ::= { itu-t(0) identified-organization(4)
	// etsi(0) cades(19122) id-spq (2) 1 }
	OIDIdSpDocSpecification = asn1.ObjectIdentifier{0, 4, 0, 19122, 2, 1}

	// OIDAttributeCertificateRefsOid is
	// id-aa-ets-attrCertificateRefs OBJECT IDENTIFIER ::= { iso(1) member-body(2)
	// us(840) rsadsi(113549) pkcs(1) pkcs-9(9) smime(16) id-aa(2) 44 }
	OIDAttributeCertificateRefsOid = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 9, 16, 2, 44}

	// OIDAttributeRevocationRefsOid is
	// id-aa-ets-attrRevocationRefs OBJECT IDENTIFIER ::= { iso(1) member-body(2)
	// us(840) rsadsi(113549) pkcs(1) pkcs-9(9) smime(16) id-aa(2) 45}
	OIDAttributeRevocationRefsOid = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 9, 16, 2, 45}

	// OIDIdAtRole is
	// id-at OBJECT IDENTIFIER ::= { joint-iso-ccitt(2) ds(5) 4 }
	// id-at-role OBJECT IDENTIFIER ::= { id-at 72}
	OIDIdAtRole = asn1.ObjectIdentifier{2, 5, 4, 72}

	// OIDAdbeRevocationInfoArchival is defined as
	// adbe-revocationInfoArchival {adbe(1.2.840.113583) acrobat(1) security(1) 8} in
	// "PDF Reference, fifth edition: Adobe Portable Document Format, Version 1.6",
	// Adobe Systems Incorporated, 2004.
	OIDAdbeRevocationInfoArchival = asn1.ObjectIdentifier{1, 2, 840, 113583, 1, 1, 8}

	// --- ETSI TS 119 495 V1.4.1

	// OIDPsd2QcStatement is
	// etsi-psd2-qcStatement QC-STATEMENT ::= {SYNTAX PSD2QcType IDENTIFIED BY
	// id-etsi-psd2-qcStatement }
	// id-etsi-psd2-qcStatement OBJECT IDENTIFIER ::= {itu-t(0)
	// identified-organization(4) etsi(0) psd2(19495) qcstatement(2) }
	OIDPsd2QcStatement = asn1.ObjectIdentifier{0, 4, 0, 19495, 2}

	// OIDIdEtsiQcsQcCClegislation is
	// esi4-qcStatement-7 QC-STATEMENT ::= { SYNTAX QcCClegislation IDENTIFIED BY
	// id-etsi-qcsQcCClegislation }
	// id-etsi-qcs-QcCClegislation OBJECT IDENTIFIER ::= { id-etsi-qcs 7 }
	OIDIdEtsiQcsQcCClegislation = asn1.ObjectIdentifier{0, 4, 0, 1862, 1, 7}

	// OIDIdEtsiQcsQcIdentMethod is
	// esi4-qcStatement-8 QC-STATEMENT ::= { SYNTAX QcIdentMethod IDENTIFIED BY
	// id-etsi-qcs-QcIdentMethod }
	// id-etsi-qcs-QcIdentMethod OBJECT IDENTIFIER ::= { id-etsi-qcs 8 }
	OIDIdEtsiQcsQcIdentMethod = asn1.ObjectIdentifier{0, 4, 0, 1862, 1, 8}

	// OIDIdEtsiQcsQcQSCDlegislation is
	// esi4-qcStatement-9 QC-STATEMENT ::= { SYNTAX QcQSCDlegislation IDENTIFIED BY
	// id-etsi-qcs-QcQCSDlegislation }
	// id-etsi-qcs-QcQSCDlegislation OBJECT IDENTIFIER ::= { id-etsi-qcs 9 }
	OIDIdEtsiQcsQcQSCDlegislation = asn1.ObjectIdentifier{0, 4, 0, 1862, 1, 9}

	// --- ETSI TS 119 412-6 V1.4.1

	// OIDIdEtsiQctPid is
	// id-etsi-qct-pid OBJECT IDENTIFIER ::= { id-etsi-eidas2-qct-extensions 1 }
	// -- Certificate for PID provider sign/seal certificate
	OIDIdEtsiQctPid = asn1.ObjectIdentifier{0, 4, 0, 194126, 1, 1}

	// OIDIdEtsiQctWal is
	// id-etsi-qct-wal OBJECT IDENTIFIER ::= { id-etsi-eidas2-qct-extensions 2 }
	// -- Certificate for Wallet provider sign/seal certificate
	OIDIdEtsiQctWal = asn1.ObjectIdentifier{0, 4, 0, 194126, 1, 2}

	// OIDIdEtsiQcsQcPSB is
	// -- PSB certificate mandatory data
	// id-etsi-qcs-QcPSB   OBJECT IDENTIFIER ::= { id-etsi-eidas2-qct-extensions 3 }
	// esi4-qcStatement-10 QC-STATEMENT ::= { SYNTAX QcPSB IDENTIFIED BY id-etsi-qcs-QcPSB }
	OIDIdEtsiQcsQcPSB = asn1.ObjectIdentifier{0, 4, 0, 194126, 1, 3}

	// OIDIdEtsiExtValassuredSTCerts is
	// EN 319 412-1 "5.2.2 Validity Assured - Short Term"
	// id-etsi-ext OBJECT IDENTIFIER ::= { itu-t(0) identified-organization(4) etsi(0)
	// id-cert-profile(194121) 2 }
	// id-etsi-ext-valassured-ST-certs OBJECT IDENTIFIER ::= { id-etsi-ext 1 }
	OIDIdEtsiExtValassuredSTCerts = asn1.ObjectIdentifier{0, 4, 0, 194121, 2, 1}

	// OIDIdAaErInternal is the ASN.1 Internal EvidenceRecord Attribute
	// id-aa-er-internal  OBJECT IDENTIFIER ::= { iso(1) member-body(2)
	// us(840) rsadsi(113549) pkcs(1) pkcs9(9) smime(16) id-aa(2) 49 }
	OIDIdAaErInternal = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 9, 16, 2, 49}

	// OIDIdAaErExternal is the ASN.1 External EvidenceRecord Attribute
	// id-aa-er-external  OBJECT IDENTIFIER ::= { iso(1) member-body(2)
	// us(840) rsadsi(113549) pkcs(1) pkcs9(9) smime(16) id-aa(2) 50 }
	OIDIdAaErExternal = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 9, 16, 2, 50}
)
