// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/trust/ServiceTypeIdentifier.java (DSS 6.5.RC1).
//
// Service type identifier (ETSI TS 119 612 V2.1.1). It specifies the identifier of the
// service type.
package qualification

// ServiceTypeIdentifier is a trust service type identifier (ETSI TS 119 612 V2.1.1).
type ServiceTypeIdentifier string

const (
	// ---- 5.5.1.1 Regulation (EU) No 910/2014 qualified trust service types

	// ServiceTypeIdentifier_CA_QC is a qualified certificate issuing trust service.
	ServiceTypeIdentifier_CA_QC ServiceTypeIdentifier = "CA_QC"
	// ServiceTypeIdentifier_OCSP_QC is a qualified OCSP certificate status service.
	ServiceTypeIdentifier_OCSP_QC ServiceTypeIdentifier = "OCSP_QC"
	// ServiceTypeIdentifier_CRL_QC is a qualified CRL certificate status service.
	ServiceTypeIdentifier_CRL_QC ServiceTypeIdentifier = "CRL_QC"
	// ServiceTypeIdentifier_TSA_QTST is a qualified electronic time stamp generation service.
	ServiceTypeIdentifier_TSA_QTST ServiceTypeIdentifier = "TSA_QTST"
	// ServiceTypeIdentifier_EDS_Q is a qualified electronic delivery service.
	ServiceTypeIdentifier_EDS_Q ServiceTypeIdentifier = "EDS_Q"
	// ServiceTypeIdentifier_EDS_REM_Q is a qualified electronic registered mail delivery
	// service.
	ServiceTypeIdentifier_EDS_REM_Q ServiceTypeIdentifier = "EDS_REM_Q"
	// ServiceTypeIdentifier_PSES_Q is a qualified preservation service.
	ServiceTypeIdentifier_PSES_Q ServiceTypeIdentifier = "PSES_Q"
	// ServiceTypeIdentifier_QESVALIDATION_Q is a qualified validation service.
	ServiceTypeIdentifier_QESVALIDATION_Q ServiceTypeIdentifier = "QESVALIDATION_Q"
	// ServiceTypeIdentifier_REMOTE_QSIG_CD_MANAGEMENT_Q is the management of remote qualified
	// electronic signature creation devices as a qualified trust service.
	ServiceTypeIdentifier_REMOTE_QSIG_CD_MANAGEMENT_Q ServiceTypeIdentifier = "REMOTE_QSIG_CD_MANAGEMENT_Q"
	// ServiceTypeIdentifier_REMOTE_QSEAL_CD_MANAGEMENT_Q is the management of remote
	// qualified electronic seal creation devices as a qualified trust service.
	ServiceTypeIdentifier_REMOTE_QSEAL_CD_MANAGEMENT_Q ServiceTypeIdentifier = "REMOTE_QSEAL_CD_MANAGEMENT_Q"
	// ServiceTypeIdentifier_EAA_Q is the issuance of qualified electronic attestations of
	// attributes.
	ServiceTypeIdentifier_EAA_Q ServiceTypeIdentifier = "EAA_Q"
	// ServiceTypeIdentifier_ELECTRONIC_ARCHIVING_Q is a qualified electronic archiving
	// service.
	ServiceTypeIdentifier_ELECTRONIC_ARCHIVING_Q ServiceTypeIdentifier = "ELECTRONIC_ARCHIVING_Q"
	// ServiceTypeIdentifier_LEDGERS_Q is a qualified trust service for the recording of
	// electronic data in qualified electronic ledgers.
	ServiceTypeIdentifier_LEDGERS_Q ServiceTypeIdentifier = "LEDGERS_Q"

	// ---- 5.5.1.2 Regulation (EU) No 910/2014 non qualified trust service types

	// ServiceTypeIdentifier_CA_PKC is a non-qualified certificate generation service.
	ServiceTypeIdentifier_CA_PKC ServiceTypeIdentifier = "CA_PKC"
	// ServiceTypeIdentifier_OCSP is a non-qualified OCSP certificate status service.
	ServiceTypeIdentifier_OCSP ServiceTypeIdentifier = "OCSP"
	// ServiceTypeIdentifier_CRL is a non-qualified CRL certificate status service.
	ServiceTypeIdentifier_CRL ServiceTypeIdentifier = "CRL"
	// ServiceTypeIdentifier_TSA is a non-qualified time-stamping generation service.
	ServiceTypeIdentifier_TSA ServiceTypeIdentifier = "TSA"
	// ServiceTypeIdentifier_TSA_TSS_QC is a non-qualified time-stamping service part of a
	// service issuing qualified certificates.
	ServiceTypeIdentifier_TSA_TSS_QC ServiceTypeIdentifier = "TSA_TSS_QC"
	// ServiceTypeIdentifier_TSA_TSS_ADESQC_AND_QES is a non-qualified time-stamping service.
	ServiceTypeIdentifier_TSA_TSS_ADESQC_AND_QES ServiceTypeIdentifier = "TSA_TSS_ADESQC_AND_QES"
	// ServiceTypeIdentifier_EDS is a non-qualified electronic delivery service.
	ServiceTypeIdentifier_EDS ServiceTypeIdentifier = "EDS"
	// ServiceTypeIdentifier_EDS_REM is a non-qualified Registered Electronic Mail delivery
	// service.
	ServiceTypeIdentifier_EDS_REM ServiceTypeIdentifier = "EDS_REM"
	// ServiceTypeIdentifier_PSES is a non-qualified preservation service.
	ServiceTypeIdentifier_PSES ServiceTypeIdentifier = "PSES"
	// ServiceTypeIdentifier_ADES_VALIDATION is a non-qualified validation service.
	ServiceTypeIdentifier_ADES_VALIDATION ServiceTypeIdentifier = "ADES_VALIDATION"
	// ServiceTypeIdentifier_ADES_GENERATION is a non-qualified generation service.
	ServiceTypeIdentifier_ADES_GENERATION ServiceTypeIdentifier = "ADES_GENERATION"
	// ServiceTypeIdentifier_REMOTE_SIG_CD_MANAGEMENT is a non-qualified trust service for the
	// management of remote electronic signature creation devices.
	ServiceTypeIdentifier_REMOTE_SIG_CD_MANAGEMENT ServiceTypeIdentifier = "REMOTE_SIG_CD_MANAGEMENT"
	// ServiceTypeIdentifier_REMOTE_SEAL_CD_MANAGEMENT is a non-qualified trust service for
	// the management of remote electronic seal creation devices.
	ServiceTypeIdentifier_REMOTE_SEAL_CD_MANAGEMENT ServiceTypeIdentifier = "REMOTE_SEAL_CD_MANAGEMENT"
	// ServiceTypeIdentifier_EAA is the issuance of non-qualified electronic attestations of
	// attributes.
	ServiceTypeIdentifier_EAA ServiceTypeIdentifier = "EAA"
	// ServiceTypeIdentifier_ELECTRONIC_ARCHIVING is a non-qualified electronic archiving
	// service.
	ServiceTypeIdentifier_ELECTRONIC_ARCHIVING ServiceTypeIdentifier = "ELECTRONIC_ARCHIVING"
	// ServiceTypeIdentifier_LEDGERS is a non-qualified trust service for the recording of
	// electronic data in non-qualified electronic ledgers.
	ServiceTypeIdentifier_LEDGERS ServiceTypeIdentifier = "LEDGERS"
	// ServiceTypeIdentifier_PKC_VALIDATION is a non-qualified trust service for the
	// validation of certificates for electronic signatures/seals/website authentication.
	ServiceTypeIdentifier_PKC_VALIDATION ServiceTypeIdentifier = "PKC_VALIDATION"
	// ServiceTypeIdentifier_PKC_PRESERVATION is a non-qualified trust service for the
	// preservation of certificates for electronic signatures/seals.
	ServiceTypeIdentifier_PKC_PRESERVATION ServiceTypeIdentifier = "PKC_PRESERVATION"
	// ServiceTypeIdentifier_EAA_VALIDATION is a non-qualified trust service for the
	// validation of electronic attestation of attributes.
	ServiceTypeIdentifier_EAA_VALIDATION ServiceTypeIdentifier = "EAA_VALIDATION"
	// ServiceTypeIdentifier_TST_VALIDATION is a non-qualified trust service for the
	// validation of electronic timestamps.
	ServiceTypeIdentifier_TST_VALIDATION ServiceTypeIdentifier = "TST_VALIDATION"
	// ServiceTypeIdentifier_EDS_VALIDATION is a non-qualified trust service for the
	// validation of data transmitted through electronic registered delivery services.
	ServiceTypeIdentifier_EDS_VALIDATION ServiceTypeIdentifier = "EDS_VALIDATION"
	// ServiceTypeIdentifier_EAA_PUBEAA is the issuance of non-qualified electronic
	// attestation of attributes by or on behalf of a public sector body.
	ServiceTypeIdentifier_EAA_PUBEAA ServiceTypeIdentifier = "EAA_PUBEAA"
	// ServiceTypeIdentifier_CA_PKC_CERTSOFOTHERTYPESOFTS is a non-qualified certificate
	// generation service for certificates of other trust service types.
	ServiceTypeIdentifier_CA_PKC_CERTSOFOTHERTYPESOFTS ServiceTypeIdentifier = "CA_PKC_CERTSOFOTHERTYPESOFTS"
	// ServiceTypeIdentifier_PKC_VALIDATION_CERTSOFOTHERTYPESOFTS is a non-qualified trust
	// service for the validation of certificates for other trust service types.
	ServiceTypeIdentifier_PKC_VALIDATION_CERTSOFOTHERTYPESOFTS ServiceTypeIdentifier = "PKC_VALIDATION_CERTSOFOTHERTYPESOFTS"

	// ---- 5.5.1.3 Trust service types not defined in Regulation (EU) No 910/2014 but
	// nationally defined

	// ServiceTypeIdentifier_RA is a registration service.
	ServiceTypeIdentifier_RA ServiceTypeIdentifier = "RA"
	// ServiceTypeIdentifier_RA_NOTHAVINGPKIID is a registration service not identified by a
	// PKI-based public key.
	ServiceTypeIdentifier_RA_NOTHAVINGPKIID ServiceTypeIdentifier = "RA_NOTHAVINGPKIID"
	// ServiceTypeIdentifier_ACA is an attribute certificate generation service.
	ServiceTypeIdentifier_ACA ServiceTypeIdentifier = "ACA"
	// ServiceTypeIdentifier_SIGNATUREPOLICYAUTHORITY is a signature policy authority
	// service.
	ServiceTypeIdentifier_SIGNATUREPOLICYAUTHORITY ServiceTypeIdentifier = "SIGNATUREPOLICYAUTHORITY"
	// ServiceTypeIdentifier_ARCHIV is an archival service.
	ServiceTypeIdentifier_ARCHIV ServiceTypeIdentifier = "ARCHIV"
	// ServiceTypeIdentifier_ARCHIV_NOTHAVINGPKIID is an archival service not identified by a
	// PKI-based public key.
	ServiceTypeIdentifier_ARCHIV_NOTHAVINGPKIID ServiceTypeIdentifier = "ARCHIV_NOTHAVINGPKIID"
	// ServiceTypeIdentifier_IDV is an identity verification service.
	ServiceTypeIdentifier_IDV ServiceTypeIdentifier = "IDV"
	// ServiceTypeIdentifier_IDV_NOTHAVINGPKIID is an identity verification service not
	// identified by a PKI-based public key.
	ServiceTypeIdentifier_IDV_NOTHAVINGPKIID ServiceTypeIdentifier = "IDV_NOTHAVINGPKIID"
	// ServiceTypeIdentifier_KESCROW is a key escrow service.
	ServiceTypeIdentifier_KESCROW ServiceTypeIdentifier = "KESCROW"
	// ServiceTypeIdentifier_KESCROW_NOTHAVINGPKIID is a key escrow service not identified by
	// a PKI-based public key.
	ServiceTypeIdentifier_KESCROW_NOTHAVINGPKIID ServiceTypeIdentifier = "KESCROW_NOTHAVINGPKIID"
	// ServiceTypeIdentifier_PPWD is an issuer of PIN- or password-based identity
	// credentials.
	ServiceTypeIdentifier_PPWD ServiceTypeIdentifier = "PPWD"
	// ServiceTypeIdentifier_PPWD_NOTHAVINGPKIID is an issuer of PIN- or password-based
	// identity credentials not identified by a PKI-based public key.
	ServiceTypeIdentifier_PPWD_NOTHAVINGPKIID ServiceTypeIdentifier = "PPWD_NOTHAVINGPKIID"
	// ServiceTypeIdentifier_TLISSUER is a service issuing trusted lists.
	ServiceTypeIdentifier_TLISSUER ServiceTypeIdentifier = "TLISSUER"
	// ServiceTypeIdentifier_NATIONALROOTCA_QC is a national root signing CA.
	ServiceTypeIdentifier_NATIONALROOTCA_QC ServiceTypeIdentifier = "NATIONALROOTCA_QC"
	// ServiceTypeIdentifier_UNSPECIFIED is a trust service of an unspecified type.
	ServiceTypeIdentifier_UNSPECIFIED ServiceTypeIdentifier = "UNSPECIFIED"
)

// serviceTypeIdentifierFields holds the (shortName, uri, qualified, national) tuple for
// each constant.
type serviceTypeIdentifierFields struct {
	shortName string
	uri       string
	qualified bool
	national  bool
}

// serviceTypeIdentifierData holds the fields for each constant.
//
// NOTE: Java's ServiceTypeIdentifier.EAA constant is declared with the shortName
// "RemoteSealCDManagement" (a copy-paste artifact in upstream); reproduced as-is.
var serviceTypeIdentifierData = map[ServiceTypeIdentifier]serviceTypeIdentifierFields{
	ServiceTypeIdentifier_CA_QC:     {"CA/QC", "http://uri.etsi.org/TrstSvc/Svctype/CA/QC", true, false},
	ServiceTypeIdentifier_OCSP_QC:   {"OCSP/QC", "http://uri.etsi.org/TrstSvc/Svctype/Certstatus/OCSP/QC", true, false},
	ServiceTypeIdentifier_CRL_QC:    {"CRL/QC", "http://uri.etsi.org/TrstSvc/Svctype/Certstatus/CRL/QC", true, false},
	ServiceTypeIdentifier_TSA_QTST:  {"TSA/QTST", "http://uri.etsi.org/TrstSvc/Svctype/TSA/QTST", true, false},
	ServiceTypeIdentifier_EDS_Q:     {"EDS/Q", "http://uri.etsi.org/TrstSvc/Svctype/EDS/Q", true, false},
	ServiceTypeIdentifier_EDS_REM_Q: {"EDS/REM/Q", "http://uri.etsi.org/TrstSvc/Svctype/EDS/REM/Q", true, false},
	ServiceTypeIdentifier_PSES_Q:    {"PSES/Q", "http://uri.etsi.org/TrstSvc/Svctype/PSES/Q", true, false},
	ServiceTypeIdentifier_QESVALIDATION_Q: {"QESValidation/Q",
		"http://uri.etsi.org/TrstSvc/Svctype/QESValidation/Q", true, false},
	ServiceTypeIdentifier_REMOTE_QSIG_CD_MANAGEMENT_Q: {"RemoteQSigCDManagement/Q",
		"http://uri.etsi.org/TrstSvc/Svctype/RemoteQSigCDManagement/Q", true, false},
	ServiceTypeIdentifier_REMOTE_QSEAL_CD_MANAGEMENT_Q: {"RemoteQSealCDManagement/Q",
		"http://uri.etsi.org/TrstSvc/Svctype/RemoteQSealCDManagement/Q", true, false},
	ServiceTypeIdentifier_EAA_Q: {"EAA/Q", "http://uri.etsi.org/TrstSvc/Svctype/EAA/Q", true, false},
	ServiceTypeIdentifier_ELECTRONIC_ARCHIVING_Q: {"ElectronicArchiving/Q",
		"http://uri.etsi.org/TrstSvc/Svctype/ElectronicArchiving/Q", true, false},
	ServiceTypeIdentifier_LEDGERS_Q: {"Ledgers/Q", "http://uri.etsi.org/TrstSvc/Svctype/Ledgers/Q", true, false},

	ServiceTypeIdentifier_CA_PKC: {"CA/PKC", "http://uri.etsi.org/TrstSvc/Svctype/CA/PKC", false, false},
	ServiceTypeIdentifier_OCSP:   {"OCSP", "http://uri.etsi.org/TrstSvc/Svctype/Certstatus/OCSP", false, false},
	ServiceTypeIdentifier_CRL:    {"CRL", "http://uri.etsi.org/TrstSvc/Svctype/Certstatus/CRL", false, false},
	ServiceTypeIdentifier_TSA:    {"TSA", "http://uri.etsi.org/TrstSvc/Svctype/TSA", false, false},
	ServiceTypeIdentifier_TSA_TSS_QC: {"TSA/TSS-QC", "http://uri.etsi.org/TrstSvc/Svctype/TSA/TSS-QC",
		false, false},
	ServiceTypeIdentifier_TSA_TSS_ADESQC_AND_QES: {"TSA/TSS-AdESQCandQES",
		"http://uri.etsi.org/TrstSvc/Svctype/TSA/TSS-AdESQCandQES", false, false},
	ServiceTypeIdentifier_EDS:             {"EDS", "http://uri.etsi.org/TrstSvc/Svctype/EDS", false, false},
	ServiceTypeIdentifier_EDS_REM:         {"EDS/REM", "http://uri.etsi.org/TrstSvc/Svctype/EDS/REM", false, false},
	ServiceTypeIdentifier_PSES:            {"PSES", "http://uri.etsi.org/TrstSvc/Svctype/PSES", false, false},
	ServiceTypeIdentifier_ADES_VALIDATION: {"AdESValidation", "http://uri.etsi.org/TrstSvc/Svctype/AdESValidation", false, false},
	ServiceTypeIdentifier_ADES_GENERATION: {"AdESGeneration", "http://uri.etsi.org/TrstSvc/Svctype/AdESGeneration", false, false},
	ServiceTypeIdentifier_REMOTE_SIG_CD_MANAGEMENT: {"RemoteSigCDManagement",
		"http://uri.etsi.org/TrstSvc/Svctype/RemoteSigCDManagement", false, false},
	ServiceTypeIdentifier_REMOTE_SEAL_CD_MANAGEMENT: {"RemoteSealCDManagement",
		"http://uri.etsi.org/TrstSvc/Svctype/RemoteSealCDManagement", false, false},
	ServiceTypeIdentifier_EAA: {"RemoteSealCDManagement", "http://uri.etsi.org/TrstSvc/Svctype/EAA", false, false},
	ServiceTypeIdentifier_ELECTRONIC_ARCHIVING: {"ElectronicArchiving",
		"http://uri.etsi.org/TrstSvc/Svctype/ElectronicArchiving", false, false},
	ServiceTypeIdentifier_LEDGERS: {"Ledgers", "http://uri.etsi.org/TrstSvc/Svctype/Ledgers", false, false},
	ServiceTypeIdentifier_PKC_VALIDATION: {"PKCValidation", "http://uri.etsi.org/TrstSvc/Svctype/PKCValidation",
		false, false},
	ServiceTypeIdentifier_PKC_PRESERVATION: {"PKCPreservation", "http://uri.etsi.org/TrstSvc/Svctype/PKCPreservation",
		false, false},
	ServiceTypeIdentifier_EAA_VALIDATION: {"EAAValidation", "http://uri.etsi.org/TrstSvc/Svctype/EAAValidation",
		false, false},
	ServiceTypeIdentifier_TST_VALIDATION: {"TSTValidation", "http://uri.etsi.org/TrstSvc/Svctype/TSTValidation",
		false, false},
	ServiceTypeIdentifier_EDS_VALIDATION: {"TSTValidatiEDSValidationon",
		"http://uri.etsi.org/TrstSvc/Svctype/EDSValidation", false, false},
	ServiceTypeIdentifier_EAA_PUBEAA: {"EAA/Pub-EAA", "http://uri.etsi.org/TrstSvc/Svctype/EAA/Pub-EAA", false, false},
	ServiceTypeIdentifier_CA_PKC_CERTSOFOTHERTYPESOFTS: {"CA/PKC/CertsforOtherTypesOfTS",
		"http://uri.etsi.org/TrstSvc/Svctype/CA/PKC/CertsforOtherTypesOfTS", false, false},
	ServiceTypeIdentifier_PKC_VALIDATION_CERTSOFOTHERTYPESOFTS: {"PKCValidation/CertsforOtherTypesOfTS",
		"http://uri.etsi.org/TrstSvc/Svctype/PKCValidation/CertsforOtherTypesOfTS", false, false},

	ServiceTypeIdentifier_RA: {"RA", "http://uri.etsi.org/TrstSvc/Svctype/RA", false, true},
	ServiceTypeIdentifier_RA_NOTHAVINGPKIID: {"RA/nothavingPKIid",
		"http://uri.etsi.org/TrstSvc/Svctype/RA/nothavingPKIid", false, true},
	ServiceTypeIdentifier_ACA: {"ACA", "http://uri.etsi.org/TrstSvc/Svctype/ACA", false, true},
	ServiceTypeIdentifier_SIGNATUREPOLICYAUTHORITY: {"SignaturePolicyAuthority",
		"http://uri.etsi.org/TrstSvc/Svctype/SignaturePolicyAuthority", false, true},
	ServiceTypeIdentifier_ARCHIV: {"Archiv", "http://uri.etsi.org/TrstSvc/Svctype/Archiv", false, true},
	ServiceTypeIdentifier_ARCHIV_NOTHAVINGPKIID: {"Archiv/nothavingPKIid",
		"http://uri.etsi.org/TrstSvc/Svctype/Archiv/nothavingPKIid", false, true},
	ServiceTypeIdentifier_IDV: {"IdV", "http://uri.etsi.org/TrstSvc/Svctype/IdV", false, true},
	ServiceTypeIdentifier_IDV_NOTHAVINGPKIID: {"IdV/nothavingPKIid",
		"http://uri.etsi.org/TrstSvc/Svctype/IdV/nothavingPKIid", false, true},
	ServiceTypeIdentifier_KESCROW: {"KEscrow", "http://uri.etsi.org/TrstSvc/Svctype/KEscrow", false, true},
	ServiceTypeIdentifier_KESCROW_NOTHAVINGPKIID: {"KEscrow/nothavingPKIid",
		"http://uri.etsi.org/TrstSvc/Svctype/KEscrow/nothavingPKIid", false, true},
	ServiceTypeIdentifier_PPWD: {"PPwd", "http://uri.etsi.org/TrstSvc/Svctype/PPwd", false, true},
	ServiceTypeIdentifier_PPWD_NOTHAVINGPKIID: {"PPwd/nothavingPKIid",
		"http://uri.etsi.org/TrstSvc/Svctype/PPwd/nothavingPKIid", false, true},
	ServiceTypeIdentifier_TLISSUER: {"TLIssuer", "http://uri.etsi.org/TrstSvd/Svctype/TLIssuer", false, true},
	ServiceTypeIdentifier_NATIONALROOTCA_QC: {"NationalRootCA-QC",
		"http://uri.etsi.org/TrstSvc/Svctype/NationalRootCA-QC", false, true},
	ServiceTypeIdentifier_UNSPECIFIED: {"unspecified", "http://uri.etsi.org/TrstSvc/Svctype/unspecified", false, true},
}

// serviceTypeIdentifierValues returns all constants in declaration order.
func serviceTypeIdentifierValues() []ServiceTypeIdentifier {
	return []ServiceTypeIdentifier{
		ServiceTypeIdentifier_CA_QC,
		ServiceTypeIdentifier_OCSP_QC,
		ServiceTypeIdentifier_CRL_QC,
		ServiceTypeIdentifier_TSA_QTST,
		ServiceTypeIdentifier_EDS_Q,
		ServiceTypeIdentifier_EDS_REM_Q,
		ServiceTypeIdentifier_PSES_Q,
		ServiceTypeIdentifier_QESVALIDATION_Q,
		ServiceTypeIdentifier_REMOTE_QSIG_CD_MANAGEMENT_Q,
		ServiceTypeIdentifier_REMOTE_QSEAL_CD_MANAGEMENT_Q,
		ServiceTypeIdentifier_EAA_Q,
		ServiceTypeIdentifier_ELECTRONIC_ARCHIVING_Q,
		ServiceTypeIdentifier_LEDGERS_Q,
		ServiceTypeIdentifier_CA_PKC,
		ServiceTypeIdentifier_OCSP,
		ServiceTypeIdentifier_CRL,
		ServiceTypeIdentifier_TSA,
		ServiceTypeIdentifier_TSA_TSS_QC,
		ServiceTypeIdentifier_TSA_TSS_ADESQC_AND_QES,
		ServiceTypeIdentifier_EDS,
		ServiceTypeIdentifier_EDS_REM,
		ServiceTypeIdentifier_PSES,
		ServiceTypeIdentifier_ADES_VALIDATION,
		ServiceTypeIdentifier_ADES_GENERATION,
		ServiceTypeIdentifier_REMOTE_SIG_CD_MANAGEMENT,
		ServiceTypeIdentifier_REMOTE_SEAL_CD_MANAGEMENT,
		ServiceTypeIdentifier_EAA,
		ServiceTypeIdentifier_ELECTRONIC_ARCHIVING,
		ServiceTypeIdentifier_LEDGERS,
		ServiceTypeIdentifier_PKC_VALIDATION,
		ServiceTypeIdentifier_PKC_PRESERVATION,
		ServiceTypeIdentifier_EAA_VALIDATION,
		ServiceTypeIdentifier_TST_VALIDATION,
		ServiceTypeIdentifier_EDS_VALIDATION,
		ServiceTypeIdentifier_EAA_PUBEAA,
		ServiceTypeIdentifier_CA_PKC_CERTSOFOTHERTYPESOFTS,
		ServiceTypeIdentifier_PKC_VALIDATION_CERTSOFOTHERTYPESOFTS,
		ServiceTypeIdentifier_RA,
		ServiceTypeIdentifier_RA_NOTHAVINGPKIID,
		ServiceTypeIdentifier_ACA,
		ServiceTypeIdentifier_SIGNATUREPOLICYAUTHORITY,
		ServiceTypeIdentifier_ARCHIV,
		ServiceTypeIdentifier_ARCHIV_NOTHAVINGPKIID,
		ServiceTypeIdentifier_IDV,
		ServiceTypeIdentifier_IDV_NOTHAVINGPKIID,
		ServiceTypeIdentifier_KESCROW,
		ServiceTypeIdentifier_KESCROW_NOTHAVINGPKIID,
		ServiceTypeIdentifier_PPWD,
		ServiceTypeIdentifier_PPWD_NOTHAVINGPKIID,
		ServiceTypeIdentifier_TLISSUER,
		ServiceTypeIdentifier_NATIONALROOTCA_QC,
		ServiceTypeIdentifier_UNSPECIFIED,
	}
}

// ShortName gets the identifier's label. Port of getShortName().
func (s ServiceTypeIdentifier) ShortName() string {
	return serviceTypeIdentifierData[s].shortName
}

// URI gets the identifier's URI. Port of getUri().
func (s ServiceTypeIdentifier) URI() string {
	return serviceTypeIdentifierData[s].uri
}

// IsQualified gets whether the identifier corresponds to a qualified status. Port of
// isQualified().
func (s ServiceTypeIdentifier) IsQualified() bool {
	return serviceTypeIdentifierData[s].qualified
}

// IsNational gets whether the identifier corresponds to a national status. Port of
// isNational().
func (s ServiceTypeIdentifier) IsNational() bool {
	return serviceTypeIdentifierData[s].national
}

// ServiceTypeIdentifierIsCaQc checks whether the serviceTypeIdentifier is CA/QC. Port of
// isCaQc(String).
func ServiceTypeIdentifierIsCaQc(serviceTypeIdentifier string) bool {
	return ServiceTypeIdentifier_CA_QC.URI() == serviceTypeIdentifier
}

// ServiceTypeIdentifierIsQTST checks whether the serviceTypeIdentifier is TSA/QTST. Port
// of isQTST(String).
func ServiceTypeIdentifierIsQTST(serviceTypeIdentifier string) bool {
	return ServiceTypeIdentifier_TSA_QTST.URI() == serviceTypeIdentifier
}

// ServiceTypeIdentifierIsQEAA checks whether the serviceTypeIdentifier is EAA/Q. Port of
// isQEAA(String).
func ServiceTypeIdentifierIsQEAA(serviceTypeIdentifier string) bool {
	return ServiceTypeIdentifier_EAA_Q.URI() == serviceTypeIdentifier
}

// ServiceTypeIdentifierFromUri returns a corresponding ServiceTypeIdentifier by the given
// uri, or "" (Java's null) if none matches. Port of fromUri(String).
func ServiceTypeIdentifierFromUri(uri string) ServiceTypeIdentifier {
	for _, sti := range serviceTypeIdentifierValues() {
		if serviceTypeIdentifierData[sti].uri == uri {
			return sti
		}
	}
	return ""
}
