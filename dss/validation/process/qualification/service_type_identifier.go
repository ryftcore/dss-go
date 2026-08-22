// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/trust/ServiceTypeIdentifier.java (DSS 6.5.RC1).
//
// Service type identifier (ETSI TS 119 612 V2.1.1). It specifies the identifier of the
// service type.
package qualification

// ServiceTypeIdentifier is a trust service type identifier (ETSI TS 119 612 V2.1.1).
type ServiceTypeIdentifier string

const (
	// ---- 5.5.1.1 Regulation (EU) No 910/2014 qualified trust service types

	// ServiceTypeIdentifierCAQC is a qualified certificate issuing trust service.
	ServiceTypeIdentifierCAQC ServiceTypeIdentifier = "CA_QC"
	// ServiceTypeIdentifierOCSPQC is a qualified OCSP certificate status service.
	ServiceTypeIdentifierOCSPQC ServiceTypeIdentifier = "OCSP_QC"
	// ServiceTypeIdentifierCRLQC is a qualified CRL certificate status service.
	ServiceTypeIdentifierCRLQC ServiceTypeIdentifier = "CRL_QC"
	// ServiceTypeIdentifierTSAQTST is a qualified electronic time stamp generation service.
	ServiceTypeIdentifierTSAQTST ServiceTypeIdentifier = "TSA_QTST"
	// ServiceTypeIdentifierEDSQ is a qualified electronic delivery service.
	ServiceTypeIdentifierEDSQ ServiceTypeIdentifier = "EDS_Q"
	// ServiceTypeIdentifierEDSREMQ is a qualified electronic registered mail delivery
	// service.
	ServiceTypeIdentifierEDSREMQ ServiceTypeIdentifier = "EDS_REM_Q"
	// ServiceTypeIdentifierPSESQ is a qualified preservation service.
	ServiceTypeIdentifierPSESQ ServiceTypeIdentifier = "PSES_Q"
	// ServiceTypeIdentifierQESValidationQ is a qualified validation service.
	ServiceTypeIdentifierQESValidationQ ServiceTypeIdentifier = "QESVALIDATION_Q"
	// ServiceTypeIdentifierRemoteQSigCDManagementQ is the management of remote qualified
	// electronic signature creation devices as a qualified trust service.
	ServiceTypeIdentifierRemoteQSigCDManagementQ ServiceTypeIdentifier = "REMOTE_QSIG_CD_MANAGEMENT_Q"
	// ServiceTypeIdentifierRemoteQSealCDManagementQ is the management of remote
	// qualified electronic seal creation devices as a qualified trust service.
	ServiceTypeIdentifierRemoteQSealCDManagementQ ServiceTypeIdentifier = "REMOTE_QSEAL_CD_MANAGEMENT_Q"
	// ServiceTypeIdentifierEAAQ is the issuance of qualified electronic attestations of
	// attributes.
	ServiceTypeIdentifierEAAQ ServiceTypeIdentifier = "EAA_Q"
	// ServiceTypeIdentifierElectronicArchivingQ is a qualified electronic archiving
	// service.
	ServiceTypeIdentifierElectronicArchivingQ ServiceTypeIdentifier = "ELECTRONIC_ARCHIVING_Q"
	// ServiceTypeIdentifierLedgersQ is a qualified trust service for the recording of
	// electronic data in qualified electronic ledgers.
	ServiceTypeIdentifierLedgersQ ServiceTypeIdentifier = "LEDGERS_Q"

	// ---- 5.5.1.2 Regulation (EU) No 910/2014 non qualified trust service types

	// ServiceTypeIdentifierCAPKC is a non-qualified certificate generation service.
	ServiceTypeIdentifierCAPKC ServiceTypeIdentifier = "CA_PKC"
	// ServiceTypeIdentifierOCSP is a non-qualified OCSP certificate status service.
	ServiceTypeIdentifierOCSP ServiceTypeIdentifier = "OCSP"
	// ServiceTypeIdentifierCRL is a non-qualified CRL certificate status service.
	ServiceTypeIdentifierCRL ServiceTypeIdentifier = "CRL"
	// ServiceTypeIdentifierTSA is a non-qualified time-stamping generation service.
	ServiceTypeIdentifierTSA ServiceTypeIdentifier = "TSA"
	// ServiceTypeIdentifierTSATSSQC is a non-qualified time-stamping service part of a
	// service issuing qualified certificates.
	ServiceTypeIdentifierTSATSSQC ServiceTypeIdentifier = "TSA_TSS_QC"
	// ServiceTypeIdentifierTSATSSAdESQCAndQES is a non-qualified time-stamping service.
	ServiceTypeIdentifierTSATSSAdESQCAndQES ServiceTypeIdentifier = "TSA_TSS_ADESQC_AND_QES"
	// ServiceTypeIdentifierEDS is a non-qualified electronic delivery service.
	ServiceTypeIdentifierEDS ServiceTypeIdentifier = "EDS"
	// ServiceTypeIdentifierEDSREM is a non-qualified Registered Electronic Mail delivery
	// service.
	ServiceTypeIdentifierEDSREM ServiceTypeIdentifier = "EDS_REM"
	// ServiceTypeIdentifierPSES is a non-qualified preservation service.
	ServiceTypeIdentifierPSES ServiceTypeIdentifier = "PSES"
	// ServiceTypeIdentifierAdESValidation is a non-qualified validation service.
	ServiceTypeIdentifierAdESValidation ServiceTypeIdentifier = "ADES_VALIDATION"
	// ServiceTypeIdentifierAdESGeneration is a non-qualified generation service.
	ServiceTypeIdentifierAdESGeneration ServiceTypeIdentifier = "ADES_GENERATION"
	// ServiceTypeIdentifierRemoteSigCDManagement is a non-qualified trust service for the
	// management of remote electronic signature creation devices.
	ServiceTypeIdentifierRemoteSigCDManagement ServiceTypeIdentifier = "REMOTE_SIG_CD_MANAGEMENT"
	// ServiceTypeIdentifierRemoteSealCDManagement is a non-qualified trust service for
	// the management of remote electronic seal creation devices.
	ServiceTypeIdentifierRemoteSealCDManagement ServiceTypeIdentifier = "REMOTE_SEAL_CD_MANAGEMENT"
	// ServiceTypeIdentifierEAA is the issuance of non-qualified electronic attestations of
	// attributes.
	ServiceTypeIdentifierEAA ServiceTypeIdentifier = "EAA"
	// ServiceTypeIdentifierElectronicArchiving is a non-qualified electronic archiving
	// service.
	ServiceTypeIdentifierElectronicArchiving ServiceTypeIdentifier = "ELECTRONIC_ARCHIVING"
	// ServiceTypeIdentifierLedgers is a non-qualified trust service for the recording of
	// electronic data in non-qualified electronic ledgers.
	ServiceTypeIdentifierLedgers ServiceTypeIdentifier = "LEDGERS"
	// ServiceTypeIdentifierPKCValidation is a non-qualified trust service for the
	// validation of certificates for electronic signatures/seals/website authentication.
	ServiceTypeIdentifierPKCValidation ServiceTypeIdentifier = "PKC_VALIDATION"
	// ServiceTypeIdentifierPKCPreservation is a non-qualified trust service for the
	// preservation of certificates for electronic signatures/seals.
	ServiceTypeIdentifierPKCPreservation ServiceTypeIdentifier = "PKC_PRESERVATION"
	// ServiceTypeIdentifierEAAValidation is a non-qualified trust service for the
	// validation of electronic attestation of attributes.
	ServiceTypeIdentifierEAAValidation ServiceTypeIdentifier = "EAA_VALIDATION"
	// ServiceTypeIdentifierTSTValidation is a non-qualified trust service for the
	// validation of electronic timestamps.
	ServiceTypeIdentifierTSTValidation ServiceTypeIdentifier = "TST_VALIDATION"
	// ServiceTypeIdentifierEDSValidation is a non-qualified trust service for the
	// validation of data transmitted through electronic registered delivery services.
	ServiceTypeIdentifierEDSValidation ServiceTypeIdentifier = "EDS_VALIDATION"
	// ServiceTypeIdentifierEAAPubEAA is the issuance of non-qualified electronic
	// attestation of attributes by or on behalf of a public sector body.
	ServiceTypeIdentifierEAAPubEAA ServiceTypeIdentifier = "EAA_PUBEAA"
	// ServiceTypeIdentifierCAPKCCertsOfOtherTypesOfTS is a non-qualified certificate
	// generation service for certificates of other trust service types.
	ServiceTypeIdentifierCAPKCCertsOfOtherTypesOfTS ServiceTypeIdentifier = "CA_PKC_CERTSOFOTHERTYPESOFTS"
	// ServiceTypeIdentifierPKCValidationCertsOfOtherTypesOfTS is a non-qualified trust
	// service for the validation of certificates for other trust service types.
	ServiceTypeIdentifierPKCValidationCertsOfOtherTypesOfTS ServiceTypeIdentifier = "PKC_VALIDATION_CERTSOFOTHERTYPESOFTS"

	// ---- 5.5.1.3 Trust service types not defined in Regulation (EU) No 910/2014 but
	// nationally defined

	// ServiceTypeIdentifierRA is a registration service.
	ServiceTypeIdentifierRA ServiceTypeIdentifier = "RA"
	// ServiceTypeIdentifierRANotHavingPKIID is a registration service not identified by a
	// PKI-based public key.
	ServiceTypeIdentifierRANotHavingPKIID ServiceTypeIdentifier = "RA_NOTHAVINGPKIID"
	// ServiceTypeIdentifierACA is an attribute certificate generation service.
	ServiceTypeIdentifierACA ServiceTypeIdentifier = "ACA"
	// ServiceTypeIdentifierSignaturePolicyAuthority is a signature policy authority
	// service.
	ServiceTypeIdentifierSignaturePolicyAuthority ServiceTypeIdentifier = "SIGNATUREPOLICYAUTHORITY"
	// ServiceTypeIdentifierArchiv is an archival service.
	ServiceTypeIdentifierArchiv ServiceTypeIdentifier = "ARCHIV"
	// ServiceTypeIdentifierArchivNotHavingPKIID is an archival service not identified by a
	// PKI-based public key.
	ServiceTypeIdentifierArchivNotHavingPKIID ServiceTypeIdentifier = "ARCHIV_NOTHAVINGPKIID"
	// ServiceTypeIdentifierIDV is an identity verification service.
	ServiceTypeIdentifierIDV ServiceTypeIdentifier = "IDV"
	// ServiceTypeIdentifierIDVNotHavingPKIID is an identity verification service not
	// identified by a PKI-based public key.
	ServiceTypeIdentifierIDVNotHavingPKIID ServiceTypeIdentifier = "IDV_NOTHAVINGPKIID"
	// ServiceTypeIdentifierKEscrow is a key escrow service.
	ServiceTypeIdentifierKEscrow ServiceTypeIdentifier = "KESCROW"
	// ServiceTypeIdentifierKEscrowNotHavingPKIID is a key escrow service not identified by
	// a PKI-based public key.
	ServiceTypeIdentifierKEscrowNotHavingPKIID ServiceTypeIdentifier = "KESCROW_NOTHAVINGPKIID"
	// ServiceTypeIdentifierPPWD is an issuer of PIN- or password-based identity
	// credentials.
	ServiceTypeIdentifierPPWD ServiceTypeIdentifier = "PPWD"
	// ServiceTypeIdentifierPPWDNotHavingPKIID is an issuer of PIN- or password-based
	// identity credentials not identified by a PKI-based public key.
	ServiceTypeIdentifierPPWDNotHavingPKIID ServiceTypeIdentifier = "PPWD_NOTHAVINGPKIID"
	// ServiceTypeIdentifierTLIssuer is a service issuing trusted lists.
	ServiceTypeIdentifierTLIssuer ServiceTypeIdentifier = "TLISSUER"
	// ServiceTypeIdentifierNationalRootCAQC is a national root signing CA.
	ServiceTypeIdentifierNationalRootCAQC ServiceTypeIdentifier = "NATIONALROOTCA_QC"
	// ServiceTypeIdentifierUnspecified is a trust service of an unspecified type.
	ServiceTypeIdentifierUnspecified ServiceTypeIdentifier = "UNSPECIFIED"
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
	ServiceTypeIdentifierCAQC:    {"CA/QC", "http://uri.etsi.org/TrstSvc/Svctype/CA/QC", true, false},
	ServiceTypeIdentifierOCSPQC:  {"OCSP/QC", "http://uri.etsi.org/TrstSvc/Svctype/Certstatus/OCSP/QC", true, false},
	ServiceTypeIdentifierCRLQC:   {"CRL/QC", "http://uri.etsi.org/TrstSvc/Svctype/Certstatus/CRL/QC", true, false},
	ServiceTypeIdentifierTSAQTST: {"TSA/QTST", "http://uri.etsi.org/TrstSvc/Svctype/TSA/QTST", true, false},
	ServiceTypeIdentifierEDSQ:    {"EDS/Q", "http://uri.etsi.org/TrstSvc/Svctype/EDS/Q", true, false},
	ServiceTypeIdentifierEDSREMQ: {"EDS/REM/Q", "http://uri.etsi.org/TrstSvc/Svctype/EDS/REM/Q", true, false},
	ServiceTypeIdentifierPSESQ:   {"PSES/Q", "http://uri.etsi.org/TrstSvc/Svctype/PSES/Q", true, false},
	ServiceTypeIdentifierQESValidationQ: {"QESValidation/Q",
		"http://uri.etsi.org/TrstSvc/Svctype/QESValidation/Q", true, false},
	ServiceTypeIdentifierRemoteQSigCDManagementQ: {"RemoteQSigCDManagement/Q",
		"http://uri.etsi.org/TrstSvc/Svctype/RemoteQSigCDManagement/Q", true, false},
	ServiceTypeIdentifierRemoteQSealCDManagementQ: {"RemoteQSealCDManagement/Q",
		"http://uri.etsi.org/TrstSvc/Svctype/RemoteQSealCDManagement/Q", true, false},
	ServiceTypeIdentifierEAAQ: {"EAA/Q", "http://uri.etsi.org/TrstSvc/Svctype/EAA/Q", true, false},
	ServiceTypeIdentifierElectronicArchivingQ: {"ElectronicArchiving/Q",
		"http://uri.etsi.org/TrstSvc/Svctype/ElectronicArchiving/Q", true, false},
	ServiceTypeIdentifierLedgersQ: {"Ledgers/Q", "http://uri.etsi.org/TrstSvc/Svctype/Ledgers/Q", true, false},

	ServiceTypeIdentifierCAPKC: {"CA/PKC", "http://uri.etsi.org/TrstSvc/Svctype/CA/PKC", false, false},
	ServiceTypeIdentifierOCSP:  {"OCSP", "http://uri.etsi.org/TrstSvc/Svctype/Certstatus/OCSP", false, false},
	ServiceTypeIdentifierCRL:   {"CRL", "http://uri.etsi.org/TrstSvc/Svctype/Certstatus/CRL", false, false},
	ServiceTypeIdentifierTSA:   {"TSA", "http://uri.etsi.org/TrstSvc/Svctype/TSA", false, false},
	ServiceTypeIdentifierTSATSSQC: {"TSA/TSS-QC", "http://uri.etsi.org/TrstSvc/Svctype/TSA/TSS-QC",
		false, false},
	ServiceTypeIdentifierTSATSSAdESQCAndQES: {"TSA/TSS-AdESQCandQES",
		"http://uri.etsi.org/TrstSvc/Svctype/TSA/TSS-AdESQCandQES", false, false},
	ServiceTypeIdentifierEDS:            {"EDS", "http://uri.etsi.org/TrstSvc/Svctype/EDS", false, false},
	ServiceTypeIdentifierEDSREM:         {"EDS/REM", "http://uri.etsi.org/TrstSvc/Svctype/EDS/REM", false, false},
	ServiceTypeIdentifierPSES:           {"PSES", "http://uri.etsi.org/TrstSvc/Svctype/PSES", false, false},
	ServiceTypeIdentifierAdESValidation: {"AdESValidation", "http://uri.etsi.org/TrstSvc/Svctype/AdESValidation", false, false},
	ServiceTypeIdentifierAdESGeneration: {"AdESGeneration", "http://uri.etsi.org/TrstSvc/Svctype/AdESGeneration", false, false},
	ServiceTypeIdentifierRemoteSigCDManagement: {"RemoteSigCDManagement",
		"http://uri.etsi.org/TrstSvc/Svctype/RemoteSigCDManagement", false, false},
	ServiceTypeIdentifierRemoteSealCDManagement: {"RemoteSealCDManagement",
		"http://uri.etsi.org/TrstSvc/Svctype/RemoteSealCDManagement", false, false},
	ServiceTypeIdentifierEAA: {"RemoteSealCDManagement", "http://uri.etsi.org/TrstSvc/Svctype/EAA", false, false},
	ServiceTypeIdentifierElectronicArchiving: {"ElectronicArchiving",
		"http://uri.etsi.org/TrstSvc/Svctype/ElectronicArchiving", false, false},
	ServiceTypeIdentifierLedgers: {"Ledgers", "http://uri.etsi.org/TrstSvc/Svctype/Ledgers", false, false},
	ServiceTypeIdentifierPKCValidation: {"PKCValidation", "http://uri.etsi.org/TrstSvc/Svctype/PKCValidation",
		false, false},
	ServiceTypeIdentifierPKCPreservation: {"PKCPreservation", "http://uri.etsi.org/TrstSvc/Svctype/PKCPreservation",
		false, false},
	ServiceTypeIdentifierEAAValidation: {"EAAValidation", "http://uri.etsi.org/TrstSvc/Svctype/EAAValidation",
		false, false},
	ServiceTypeIdentifierTSTValidation: {"TSTValidation", "http://uri.etsi.org/TrstSvc/Svctype/TSTValidation",
		false, false},
	ServiceTypeIdentifierEDSValidation: {"TSTValidatiEDSValidationon",
		"http://uri.etsi.org/TrstSvc/Svctype/EDSValidation", false, false},
	ServiceTypeIdentifierEAAPubEAA: {"EAA/Pub-EAA", "http://uri.etsi.org/TrstSvc/Svctype/EAA/Pub-EAA", false, false},
	ServiceTypeIdentifierCAPKCCertsOfOtherTypesOfTS: {"CA/PKC/CertsforOtherTypesOfTS",
		"http://uri.etsi.org/TrstSvc/Svctype/CA/PKC/CertsforOtherTypesOfTS", false, false},
	ServiceTypeIdentifierPKCValidationCertsOfOtherTypesOfTS: {"PKCValidation/CertsforOtherTypesOfTS",
		"http://uri.etsi.org/TrstSvc/Svctype/PKCValidation/CertsforOtherTypesOfTS", false, false},

	ServiceTypeIdentifierRA: {"RA", "http://uri.etsi.org/TrstSvc/Svctype/RA", false, true},
	ServiceTypeIdentifierRANotHavingPKIID: {"RA/nothavingPKIid",
		"http://uri.etsi.org/TrstSvc/Svctype/RA/nothavingPKIid", false, true},
	ServiceTypeIdentifierACA: {"ACA", "http://uri.etsi.org/TrstSvc/Svctype/ACA", false, true},
	ServiceTypeIdentifierSignaturePolicyAuthority: {"SignaturePolicyAuthority",
		"http://uri.etsi.org/TrstSvc/Svctype/SignaturePolicyAuthority", false, true},
	ServiceTypeIdentifierArchiv: {"Archiv", "http://uri.etsi.org/TrstSvc/Svctype/Archiv", false, true},
	ServiceTypeIdentifierArchivNotHavingPKIID: {"Archiv/nothavingPKIid",
		"http://uri.etsi.org/TrstSvc/Svctype/Archiv/nothavingPKIid", false, true},
	ServiceTypeIdentifierIDV: {"IdV", "http://uri.etsi.org/TrstSvc/Svctype/IdV", false, true},
	ServiceTypeIdentifierIDVNotHavingPKIID: {"IdV/nothavingPKIid",
		"http://uri.etsi.org/TrstSvc/Svctype/IdV/nothavingPKIid", false, true},
	ServiceTypeIdentifierKEscrow: {"KEscrow", "http://uri.etsi.org/TrstSvc/Svctype/KEscrow", false, true},
	ServiceTypeIdentifierKEscrowNotHavingPKIID: {"KEscrow/nothavingPKIid",
		"http://uri.etsi.org/TrstSvc/Svctype/KEscrow/nothavingPKIid", false, true},
	ServiceTypeIdentifierPPWD: {"PPwd", "http://uri.etsi.org/TrstSvc/Svctype/PPwd", false, true},
	ServiceTypeIdentifierPPWDNotHavingPKIID: {"PPwd/nothavingPKIid",
		"http://uri.etsi.org/TrstSvc/Svctype/PPwd/nothavingPKIid", false, true},
	ServiceTypeIdentifierTLIssuer: {"TLIssuer", "http://uri.etsi.org/TrstSvd/Svctype/TLIssuer", false, true},
	ServiceTypeIdentifierNationalRootCAQC: {"NationalRootCA-QC",
		"http://uri.etsi.org/TrstSvc/Svctype/NationalRootCA-QC", false, true},
	ServiceTypeIdentifierUnspecified: {"unspecified", "http://uri.etsi.org/TrstSvc/Svctype/unspecified", false, true},
}

// serviceTypeIdentifierValues returns all constants in declaration order.
func serviceTypeIdentifierValues() []ServiceTypeIdentifier {
	return []ServiceTypeIdentifier{
		ServiceTypeIdentifierCAQC,
		ServiceTypeIdentifierOCSPQC,
		ServiceTypeIdentifierCRLQC,
		ServiceTypeIdentifierTSAQTST,
		ServiceTypeIdentifierEDSQ,
		ServiceTypeIdentifierEDSREMQ,
		ServiceTypeIdentifierPSESQ,
		ServiceTypeIdentifierQESValidationQ,
		ServiceTypeIdentifierRemoteQSigCDManagementQ,
		ServiceTypeIdentifierRemoteQSealCDManagementQ,
		ServiceTypeIdentifierEAAQ,
		ServiceTypeIdentifierElectronicArchivingQ,
		ServiceTypeIdentifierLedgersQ,
		ServiceTypeIdentifierCAPKC,
		ServiceTypeIdentifierOCSP,
		ServiceTypeIdentifierCRL,
		ServiceTypeIdentifierTSA,
		ServiceTypeIdentifierTSATSSQC,
		ServiceTypeIdentifierTSATSSAdESQCAndQES,
		ServiceTypeIdentifierEDS,
		ServiceTypeIdentifierEDSREM,
		ServiceTypeIdentifierPSES,
		ServiceTypeIdentifierAdESValidation,
		ServiceTypeIdentifierAdESGeneration,
		ServiceTypeIdentifierRemoteSigCDManagement,
		ServiceTypeIdentifierRemoteSealCDManagement,
		ServiceTypeIdentifierEAA,
		ServiceTypeIdentifierElectronicArchiving,
		ServiceTypeIdentifierLedgers,
		ServiceTypeIdentifierPKCValidation,
		ServiceTypeIdentifierPKCPreservation,
		ServiceTypeIdentifierEAAValidation,
		ServiceTypeIdentifierTSTValidation,
		ServiceTypeIdentifierEDSValidation,
		ServiceTypeIdentifierEAAPubEAA,
		ServiceTypeIdentifierCAPKCCertsOfOtherTypesOfTS,
		ServiceTypeIdentifierPKCValidationCertsOfOtherTypesOfTS,
		ServiceTypeIdentifierRA,
		ServiceTypeIdentifierRANotHavingPKIID,
		ServiceTypeIdentifierACA,
		ServiceTypeIdentifierSignaturePolicyAuthority,
		ServiceTypeIdentifierArchiv,
		ServiceTypeIdentifierArchivNotHavingPKIID,
		ServiceTypeIdentifierIDV,
		ServiceTypeIdentifierIDVNotHavingPKIID,
		ServiceTypeIdentifierKEscrow,
		ServiceTypeIdentifierKEscrowNotHavingPKIID,
		ServiceTypeIdentifierPPWD,
		ServiceTypeIdentifierPPWDNotHavingPKIID,
		ServiceTypeIdentifierTLIssuer,
		ServiceTypeIdentifierNationalRootCAQC,
		ServiceTypeIdentifierUnspecified,
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
	return ServiceTypeIdentifierCAQC.URI() == serviceTypeIdentifier
}

// ServiceTypeIdentifierIsQTST checks whether the serviceTypeIdentifier is TSA/QTST. Port
// of isQTST(String).
func ServiceTypeIdentifierIsQTST(serviceTypeIdentifier string) bool {
	return ServiceTypeIdentifierTSAQTST.URI() == serviceTypeIdentifier
}

// ServiceTypeIdentifierIsQEAA checks whether the serviceTypeIdentifier is EAA/Q. Port of
// isQEAA(String).
func ServiceTypeIdentifierIsQEAA(serviceTypeIdentifier string) bool {
	return ServiceTypeIdentifierEAAQ.URI() == serviceTypeIdentifier
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
