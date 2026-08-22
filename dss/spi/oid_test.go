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
	"id_aa_ets_archiveTimestampV2":                   OIDIdAaEtsArchiveTimestampV2,
	"id_etsi_electronicSignatureStandard_attributes": OIDIdEtsiElectronicSignatureStandardAttributes,
	"id_etsi_signer_attributes":                      OIDIdEtsiSignerAttributes,
	"id_etsi_spq":                                    OIDIdEtsiSpq,
	"id_aa_ets_mimeType":                             OIDIdAaEtsMimeType,
	"id_aa_ets_archiveTimestampV3":                   OIDIdAaEtsArchiveTimestampV3,
	"id_aa_ATSHashIndex":                             OIDIdAaATSHashIndex,
	"id_aa_ets_signerAttrV2":                         OIDIdAaEtsSignerAttrV2,
	"id_aa_ets_sigPolicyStore":                       OIDIdAaEtsSigPolicyStore,
	"id_aa_ATSHashIndexV2":                           OIDIdAaATSHashIndexV2,
	"id_aa_ATSHashIndexV3":                           OIDIdAaATSHashIndexV3,
	"id_sp_doc_specification":                        OIDIdSpDocSpecification,
	"attributeCertificateRefsOid":                    OIDAttributeCertificateRefsOid,
	"attributeRevocationRefsOid":                     OIDAttributeRevocationRefsOid,
	"id_at_role":                                     OIDIdAtRole,
	"adbe_revocationInfoArchival":                    OIDAdbeRevocationInfoArchival,
	"psd2_qcStatement":                               OIDPsd2QcStatement,
	"id_etsi_qcs_QcCClegislation":                    OIDIdEtsiQcsQcCClegislation,
	"id_etsi_qcs_QcIdentMethod":                      OIDIdEtsiQcsQcIdentMethod,
	"id_etsi_qcs_QcQSCDlegislation":                  OIDIdEtsiQcsQcQSCDlegislation,
	"id_etsi_qct_pid":                                OIDIdEtsiQctPid,
	"id_etsi_qct_wal":                                OIDIdEtsiQctWal,
	"id_etsi_qcs_QcPSB":                              OIDIdEtsiQcsQcPSB,
	"id_etsi_ext_valassured_ST_certs":                OIDIdEtsiExtValassuredSTCerts,
	"id_aa_er_internal":                              OIDIdAaErInternal,
	"id_aa_er_external":                              OIDIdAaErExternal,
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
		{"id_aa_ets_archiveTimestampV2", idAA, OIDIdAaEtsArchiveTimestampV2, 48},
		{"attributeCertificateRefsOid", idAA, OIDAttributeCertificateRefsOid, 44},
		{"attributeRevocationRefsOid", idAA, OIDAttributeRevocationRefsOid, 45},
		{"id_aa_er_internal", idAA, OIDIdAaErInternal, 49},
		{"id_aa_er_external", idAA, OIDIdAaErExternal, 50},
		{"id_aa_ets_mimeType", OIDIdEtsiElectronicSignatureStandardAttributes, OIDIdAaEtsMimeType, 1},
		{"id_aa_ets_archiveTimestampV3", OIDIdEtsiElectronicSignatureStandardAttributes, OIDIdAaEtsArchiveTimestampV3, 4},
		{"id_aa_ATSHashIndex", OIDIdEtsiElectronicSignatureStandardAttributes, OIDIdAaATSHashIndex, 5},
		{"id_aa_ets_signerAttrV2", OIDIdEtsiSignerAttributes, OIDIdAaEtsSignerAttrV2, 1},
		{"id_aa_ets_sigPolicyStore", OIDIdEtsiSignerAttributes, OIDIdAaEtsSigPolicyStore, 3},
		{"id_aa_ATSHashIndexV2", OIDIdEtsiSignerAttributes, OIDIdAaATSHashIndexV2, 4},
		{"id_aa_ATSHashIndexV3", OIDIdEtsiSignerAttributes, OIDIdAaATSHashIndexV3, 5},
		{"id_sp_doc_specification", OIDIdEtsiSpq, OIDIdSpDocSpecification, 1},
	} {
		expected := append(append(asn1.ObjectIdentifier{}, entry.parent...), entry.arc)
		if !entry.child.Equal(expected) {
			t.Errorf("%s: got %s, want %s", entry.name, entry.child, expected)
		}
	}
}
