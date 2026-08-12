package spi

import (
	"bufio"
	"encoding/asn1"
	"encoding/hex"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// oidTestConstants pairs every Go constant with the Java field it ports.
var oidTestConstants = map[string]asn1.ObjectIdentifier{
	"id_aa_ets_archiveTimestampV2":                   OID_id_aa_ets_archiveTimestampV2,
	"id_etsi_electronicSignatureStandard_attributes": OID_id_etsi_electronicSignatureStandard_attributes,
	"id_etsi_signer_attributes":                      OID_id_etsi_signer_attributes,
	"id_etsi_spq":                                    OID_id_etsi_spq,
	"id_aa_ets_mimeType":                             OID_id_aa_ets_mimeType,
	"id_aa_ets_archiveTimestampV3":                   OID_id_aa_ets_archiveTimestampV3,
	"id_aa_ATSHashIndex":                             OID_id_aa_ATSHashIndex,
	"id_aa_ets_signerAttrV2":                         OID_id_aa_ets_signerAttrV2,
	"id_aa_ets_sigPolicyStore":                       OID_id_aa_ets_sigPolicyStore,
	"id_aa_ATSHashIndexV2":                           OID_id_aa_ATSHashIndexV2,
	"id_aa_ATSHashIndexV3":                           OID_id_aa_ATSHashIndexV3,
	"id_sp_doc_specification":                        OID_id_sp_doc_specification,
	"attributeCertificateRefsOid":                    OID_attributeCertificateRefsOid,
	"attributeRevocationRefsOid":                     OID_attributeRevocationRefsOid,
	"id_at_role":                                     OID_id_at_role,
	"adbe_revocationInfoArchival":                    OID_adbe_revocationInfoArchival,
	"psd2_qcStatement":                               OID_psd2_qcStatement,
	"id_etsi_qcs_QcCClegislation":                    OID_id_etsi_qcs_QcCClegislation,
	"id_etsi_qcs_QcIdentMethod":                      OID_id_etsi_qcs_QcIdentMethod,
	"id_etsi_qcs_QcQSCDlegislation":                  OID_id_etsi_qcs_QcQSCDlegislation,
	"id_etsi_qct_pid":                                OID_id_etsi_qct_pid,
	"id_etsi_qct_wal":                                OID_id_etsi_qct_wal,
	"id_etsi_qcs_QcPSB":                              OID_id_etsi_qcs_QcPSB,
	"id_etsi_ext_valassured_ST_certs":                OID_id_etsi_ext_valassured_ST_certs,
	"id_aa_er_internal":                              OID_id_aa_er_internal,
	"id_aa_er_external":                              OID_id_aa_er_external,
}

// TestOIDKnownAnswers checks every constant against the dotted value and the DER encoding of
// the ASN1ObjectIdentifier the Java class declares - including the ones it derives with
// ASN1ObjectIdentifier#branch. Upstream OIDs are never to be "fixed": a change here is a
// compatibility break, so the answers come from a BouncyCastle 1.78.1 run over the same
// declarations (testdata/asn1/kat_oid.txt).
func TestOIDKnownAnswers(t *testing.T) {
	file, err := os.Open(filepath.Join("testdata", "asn1", "kat_oid.txt"))
	if err != nil {
		t.Fatalf("unable to open the known-answer file: %v", err)
	}
	defer func() { _ = file.Close() }()

	seen := make(map[string]bool)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fields := strings.Split(scanner.Text(), "|")
		if len(fields) != 3 {
			continue
		}
		name, value, der := fields[0], fields[1], fields[2]
		oid, ok := oidTestConstants[name]
		if !ok {
			t.Errorf("the port has no constant for the Java field %s", name)
			continue
		}
		seen[name] = true
		if got := oid.String(); got != value {
			t.Errorf("%s: got %s, want %s", name, got, value)
		}
		encoded, err := asn1.Marshal(oid)
		if err != nil {
			t.Errorf("%s: %v", name, err)
			continue
		}
		if got := hex.EncodeToString(encoded); got != der {
			t.Errorf("%s DER: got %s, want %s", name, got, der)
		}
	}
	if err := scanner.Err(); err != nil {
		t.Fatalf("unable to read the known-answer file: %v", err)
	}
	for name := range oidTestConstants {
		if !seen[name] {
			t.Errorf("the known-answer file has no entry for %s", name)
		}
	}
}

// TestOIDBranches checks the constants Java derives with branch() against their parents, so
// that a typo in a leaf is caught even when the known-answer file is regenerated.
func TestOIDBranches(t *testing.T) {
	idAA := asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 9, 16, 2}
	for _, entry := range []struct {
		name   string
		parent asn1.ObjectIdentifier
		child  asn1.ObjectIdentifier
		arc    int
	}{
		{"id_aa_ets_archiveTimestampV2", idAA, OID_id_aa_ets_archiveTimestampV2, 48},
		{"attributeCertificateRefsOid", idAA, OID_attributeCertificateRefsOid, 44},
		{"attributeRevocationRefsOid", idAA, OID_attributeRevocationRefsOid, 45},
		{"id_aa_er_internal", idAA, OID_id_aa_er_internal, 49},
		{"id_aa_er_external", idAA, OID_id_aa_er_external, 50},
		{"id_aa_ets_mimeType", OID_id_etsi_electronicSignatureStandard_attributes, OID_id_aa_ets_mimeType, 1},
		{"id_aa_ets_archiveTimestampV3", OID_id_etsi_electronicSignatureStandard_attributes, OID_id_aa_ets_archiveTimestampV3, 4},
		{"id_aa_ATSHashIndex", OID_id_etsi_electronicSignatureStandard_attributes, OID_id_aa_ATSHashIndex, 5},
		{"id_aa_ets_signerAttrV2", OID_id_etsi_signer_attributes, OID_id_aa_ets_signerAttrV2, 1},
		{"id_aa_ets_sigPolicyStore", OID_id_etsi_signer_attributes, OID_id_aa_ets_sigPolicyStore, 3},
		{"id_aa_ATSHashIndexV2", OID_id_etsi_signer_attributes, OID_id_aa_ATSHashIndexV2, 4},
		{"id_aa_ATSHashIndexV3", OID_id_etsi_signer_attributes, OID_id_aa_ATSHashIndexV3, 5},
		{"id_sp_doc_specification", OID_id_etsi_spq, OID_id_sp_doc_specification, 1},
	} {
		expected := append(append(asn1.ObjectIdentifier{}, entry.parent...), entry.arc)
		if !entry.child.Equal(expected) {
			t.Errorf("%s: got %s, want %s", entry.name, entry.child, expected)
		}
	}
}
