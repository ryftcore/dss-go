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
	// OID_id_aa_ets_archiveTimestampV2 is
	// id-aa-ets-archiveTimestampV2 OBJECT IDENTIFIER ::= { iso(1) member-body(2) us(840)
	// rsadsi(113549) pkcs(1) pkcs-9(9) smime(16) id-aa(2) 48}
	OID_id_aa_ets_archiveTimestampV2 = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 9, 16, 2, 48}

	// OID_id_etsi_electronicSignatureStandard_attributes is
	// Attributes {itu-t(0) identified-organization(4) etsi(0)
	// electronic-signature-standard(1733) attributes(2)}
	OID_id_etsi_electronicSignatureStandard_attributes = asn1.ObjectIdentifier{0, 4, 0, 1733, 2}

	// OID_id_etsi_signer_attributes is
	// Signer attributes {itu-t(0) identified-organization(4) etsi(0) cades(19122)
	// attributes(1)} (see ETSI EN 319 122-1)
	OID_id_etsi_signer_attributes = asn1.ObjectIdentifier{0, 4, 0, 19122, 1}

	// OID_id_etsi_spq is
	// Signer attributes {itu-t(0) identified-organization(4) etsi(0) cades(19122)
	// id-spq (2)} (see ETSI EN 319 122-1)
	OID_id_etsi_spq = asn1.ObjectIdentifier{0, 4, 0, 19122, 2}

	// OID_id_aa_ets_mimeType is
	// id-aa-ets-mimeType OBJECT IDENTIFIER ::= { itu-t(0) identified-organization(4)
	// etsi(0) electronic-signature-standard(1733) attributes(2) 1 }
	OID_id_aa_ets_mimeType = asn1.ObjectIdentifier{0, 4, 0, 1733, 2, 1}

	// OID_id_aa_ets_archiveTimestampV3 is
	// id-aa-ets-archiveTimestampV3 OBJECT IDENTIFIER ::= { itu-t(0)
	// identified-organization(4) etsi(0) electronic-signature-standard(1733)
	// attributes(2) 4 }
	OID_id_aa_ets_archiveTimestampV3 = asn1.ObjectIdentifier{0, 4, 0, 1733, 2, 4}

	// OID_id_aa_ATSHashIndex is
	// id-aa-ATSHashIndex OBJECT IDENTIFIER ::= { itu-t(0) identified-organization(4)
	// etsi(0) electronicsignature-standard(1733) attributes(2) 5 }
	OID_id_aa_ATSHashIndex = asn1.ObjectIdentifier{0, 4, 0, 1733, 2, 5}

	// OID_id_aa_ets_signerAttrV2 is
	// id-aa-ets-signerAttrV2 OBJECT IDENTIFIER ::= { itu-t(0) identified-organization(4)
	// etsi(0) cades(19122) attributes(1) 1 }
	OID_id_aa_ets_signerAttrV2 = asn1.ObjectIdentifier{0, 4, 0, 19122, 1, 1}

	// OID_id_aa_ets_sigPolicyStore is
	// id-aa-ets-sigPolicyStore OBJECT IDENTIFIER ::= { itu-t(0) identified-organization(4)
	// etsi(0) cades(19122) attributes(1) 3 }
	OID_id_aa_ets_sigPolicyStore = asn1.ObjectIdentifier{0, 4, 0, 19122, 1, 3}

	// OID_id_aa_ATSHashIndexV2 is
	// id-aa-ATSHashIndex-v2 OBJECT IDENTIFIER ::= { itu-t(0) identified-organization(4)
	// etsi(0) cades(19122) attributes(1) 4 }
	OID_id_aa_ATSHashIndexV2 = asn1.ObjectIdentifier{0, 4, 0, 19122, 1, 4}

	// OID_id_aa_ATSHashIndexV3 is
	// id-aa-ATSHashIndex-v3 OBJECT IDENTIFIER ::= { itu-t(0) identified-organization(4)
	// etsi(0) cades(19122) attributes(1) 5 }
	OID_id_aa_ATSHashIndexV3 = asn1.ObjectIdentifier{0, 4, 0, 19122, 1, 5}

	// OID_id_sp_doc_specification is
	// id-spq-ets-docspec OBJECT IDENTIFIER ::= { itu-t(0) identified-organization(4)
	// etsi(0) cades(19122) id-spq (2) 1 }
	OID_id_sp_doc_specification = asn1.ObjectIdentifier{0, 4, 0, 19122, 2, 1}

	// OID_attributeCertificateRefsOid is
	// id-aa-ets-attrCertificateRefs OBJECT IDENTIFIER ::= { iso(1) member-body(2)
	// us(840) rsadsi(113549) pkcs(1) pkcs-9(9) smime(16) id-aa(2) 44 }
	OID_attributeCertificateRefsOid = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 9, 16, 2, 44}

	// OID_attributeRevocationRefsOid is
	// id-aa-ets-attrRevocationRefs OBJECT IDENTIFIER ::= { iso(1) member-body(2)
	// us(840) rsadsi(113549) pkcs(1) pkcs-9(9) smime(16) id-aa(2) 45}
	OID_attributeRevocationRefsOid = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 9, 16, 2, 45}

	// OID_id_at_role is
	// id-at OBJECT IDENTIFIER ::= { joint-iso-ccitt(2) ds(5) 4 }
	// id-at-role OBJECT IDENTIFIER ::= { id-at 72}
	OID_id_at_role = asn1.ObjectIdentifier{2, 5, 4, 72}

	// OID_adbe_revocationInfoArchival is defined as
	// adbe-revocationInfoArchival {adbe(1.2.840.113583) acrobat(1) security(1) 8} in
	// "PDF Reference, fifth edition: Adobe Portable Document Format, Version 1.6",
	// Adobe Systems Incorporated, 2004.
	OID_adbe_revocationInfoArchival = asn1.ObjectIdentifier{1, 2, 840, 113583, 1, 1, 8}

	// --- ETSI TS 119 495 V1.4.1

	// OID_psd2_qcStatement is
	// etsi-psd2-qcStatement QC-STATEMENT ::= {SYNTAX PSD2QcType IDENTIFIED BY
	// id-etsi-psd2-qcStatement }
	// id-etsi-psd2-qcStatement OBJECT IDENTIFIER ::= {itu-t(0)
	// identified-organization(4) etsi(0) psd2(19495) qcstatement(2) }
	OID_psd2_qcStatement = asn1.ObjectIdentifier{0, 4, 0, 19495, 2}

	// OID_id_etsi_qcs_QcCClegislation is
	// esi4-qcStatement-7 QC-STATEMENT ::= { SYNTAX QcCClegislation IDENTIFIED BY
	// id-etsi-qcsQcCClegislation }
	// id-etsi-qcs-QcCClegislation OBJECT IDENTIFIER ::= { id-etsi-qcs 7 }
	OID_id_etsi_qcs_QcCClegislation = asn1.ObjectIdentifier{0, 4, 0, 1862, 1, 7}

	// OID_id_etsi_qcs_QcIdentMethod is
	// esi4-qcStatement-8 QC-STATEMENT ::= { SYNTAX QcIdentMethod IDENTIFIED BY
	// id-etsi-qcs-QcIdentMethod }
	// id-etsi-qcs-QcIdentMethod OBJECT IDENTIFIER ::= { id-etsi-qcs 8 }
	OID_id_etsi_qcs_QcIdentMethod = asn1.ObjectIdentifier{0, 4, 0, 1862, 1, 8}

	// OID_id_etsi_qcs_QcQSCDlegislation is
	// esi4-qcStatement-9 QC-STATEMENT ::= { SYNTAX QcQSCDlegislation IDENTIFIED BY
	// id-etsi-qcs-QcQCSDlegislation }
	// id-etsi-qcs-QcQSCDlegislation OBJECT IDENTIFIER ::= { id-etsi-qcs 9 }
	OID_id_etsi_qcs_QcQSCDlegislation = asn1.ObjectIdentifier{0, 4, 0, 1862, 1, 9}

	// --- ETSI TS 119 412-6 V1.4.1

	// OID_id_etsi_qct_pid is
	// id-etsi-qct-pid OBJECT IDENTIFIER ::= { id-etsi-eidas2-qct-extensions 1 }
	// -- Certificate for PID provider sign/seal certificate
	OID_id_etsi_qct_pid = asn1.ObjectIdentifier{0, 4, 0, 194126, 1, 1}

	// OID_id_etsi_qct_wal is
	// id-etsi-qct-wal OBJECT IDENTIFIER ::= { id-etsi-eidas2-qct-extensions 2 }
	// -- Certificate for Wallet provider sign/seal certificate
	OID_id_etsi_qct_wal = asn1.ObjectIdentifier{0, 4, 0, 194126, 1, 2}

	// OID_id_etsi_qcs_QcPSB is
	// -- PSB certificate mandatory data
	// id-etsi-qcs-QcPSB   OBJECT IDENTIFIER ::= { id-etsi-eidas2-qct-extensions 3 }
	// esi4-qcStatement-10 QC-STATEMENT ::= { SYNTAX QcPSB IDENTIFIED BY id-etsi-qcs-QcPSB }
	OID_id_etsi_qcs_QcPSB = asn1.ObjectIdentifier{0, 4, 0, 194126, 1, 3}

	// OID_id_etsi_ext_valassured_ST_certs is
	// EN 319 412-1 "5.2.2 Validity Assured - Short Term"
	// id-etsi-ext OBJECT IDENTIFIER ::= { itu-t(0) identified-organization(4) etsi(0)
	// id-cert-profile(194121) 2 }
	// id-etsi-ext-valassured-ST-certs OBJECT IDENTIFIER ::= { id-etsi-ext 1 }
	OID_id_etsi_ext_valassured_ST_certs = asn1.ObjectIdentifier{0, 4, 0, 194121, 2, 1}

	// OID_id_aa_er_internal is the ASN.1 Internal EvidenceRecord Attribute
	// id-aa-er-internal  OBJECT IDENTIFIER ::= { iso(1) member-body(2)
	// us(840) rsadsi(113549) pkcs(1) pkcs9(9) smime(16) id-aa(2) 49 }
	OID_id_aa_er_internal = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 9, 16, 2, 49}

	// OID_id_aa_er_external is the ASN.1 External EvidenceRecord Attribute
	// id-aa-er-external  OBJECT IDENTIFIER ::= { iso(1) member-body(2)
	// us(840) rsadsi(113549) pkcs(1) pkcs9(9) smime(16) id-aa(2) 50 }
	OID_id_aa_er_external = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 9, 16, 2, 50}
)
