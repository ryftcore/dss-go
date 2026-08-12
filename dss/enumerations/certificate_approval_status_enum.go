// Ported from dss-enumerations/.../CertificateApprovalStatusEnum.java (DSS 6.5.RC1).
//
// NOTE: LoTETypeEnum, LoTEServiceTypeIdentifierEnum and LoTEServiceStatusEnum
// (with their exported constants) are defined outside this file's manifest
// and are assumed to exist per the porting brief.
package enumerations

// CertificateApprovalStatusEnum contains a list of known and supported
// certificate approval statuses, based on the ETSI TS 119 602 profile
// definitions. Implements CertificateApprovalStatus.
type CertificateApprovalStatusEnum string

const (
	// CertificateApprovalStatusEnum_PID_PROVIDER represents a PID provider
	// certificate, as defined in the ETSI TS 119 605, Annex C.1.
	CertificateApprovalStatusEnum_PID_PROVIDER CertificateApprovalStatusEnum = "PID_PROVIDER"
	// CertificateApprovalStatusEnum_CERT_FOR_PID_REVOCATION represents a
	// certificate for PID revocation, according to the profile defined in
	// ETSI TS 119 602, Annex C.1.
	CertificateApprovalStatusEnum_CERT_FOR_PID_REVOCATION CertificateApprovalStatusEnum = "CERT_FOR_PID_REVOCATION"
	// CertificateApprovalStatusEnum_CERT_FOR_WALLET_ISSUANCE represents a
	// certificate for Wallet solution issuance, according to the profile
	// defined in ETSI TS 119 602, Annex C.2.
	CertificateApprovalStatusEnum_CERT_FOR_WALLET_ISSUANCE CertificateApprovalStatusEnum = "CERT_FOR_WALLET_ISSUANCE"
	// CertificateApprovalStatusEnum_CERT_FOR_WALLET_REVOCATION represents a
	// certificate for Wallet solution revocation, according to the profile
	// defined in ETSI TS 119 602, Annex C.2.
	CertificateApprovalStatusEnum_CERT_FOR_WALLET_REVOCATION CertificateApprovalStatusEnum = "CERT_FOR_WALLET_REVOCATION"
	// CertificateApprovalStatusEnum_CERT_FOR_WRPAC_ISSUANCE represents a
	// certificate for WRPAC issuance, according to the profile defined in
	// ETSI TS 119 602, Annex F.
	CertificateApprovalStatusEnum_CERT_FOR_WRPAC_ISSUANCE CertificateApprovalStatusEnum = "CERT_FOR_WRPAC_ISSUANCE"
	// CertificateApprovalStatusEnum_CERT_FOR_WRPAC_REVOCATION represents a
	// certificate for WRPAC revocation, according to the profile defined in
	// ETSI TS 119 602, Annex F.
	CertificateApprovalStatusEnum_CERT_FOR_WRPAC_REVOCATION CertificateApprovalStatusEnum = "CERT_FOR_WRPAC_REVOCATION"
	// CertificateApprovalStatusEnum_CERT_FOR_WRPRC_ISSUANCE represents a
	// certificate for WRPRC issuance, according to the profile defined in
	// ETSI TS 119 602, Annex G.
	CertificateApprovalStatusEnum_CERT_FOR_WRPRC_ISSUANCE CertificateApprovalStatusEnum = "CERT_FOR_WRPRC_ISSUANCE"
	// CertificateApprovalStatusEnum_CERT_FOR_WRPRC_REVOCATION represents a
	// certificate for WRPRC revocation, according to the profile defined in
	// ETSI TS 119 602, Annex G.
	CertificateApprovalStatusEnum_CERT_FOR_WRPRC_REVOCATION CertificateApprovalStatusEnum = "CERT_FOR_WRPRC_REVOCATION"
	// CertificateApprovalStatusEnum_NOTIFIED_CERT_FOR_PUB_EAA_ISSUANCE
	// represents a notified certificate for Pub-EAA issuance, according to
	// the profile defined in ETSI TS 119 602, Annex H.
	CertificateApprovalStatusEnum_NOTIFIED_CERT_FOR_PUB_EAA_ISSUANCE CertificateApprovalStatusEnum = "NOTIFIED_CERT_FOR_PUB_EAA_ISSUANCE"
	// CertificateApprovalStatusEnum_NOTIFIED_CERT_FOR_PUB_EAA_REVOCATION
	// represents a notified certificate for Pub-EAA revocation, according to
	// the profile defined in ETSI TS 119 602, Annex H.
	CertificateApprovalStatusEnum_NOTIFIED_CERT_FOR_PUB_EAA_REVOCATION CertificateApprovalStatusEnum = "NOTIFIED_CERT_FOR_PUB_EAA_REVOCATION"
	// CertificateApprovalStatusEnum_WITHDRAWN_CERT_FOR_PUB_EAA_ISSUANCE
	// represents a withdrawn certificate for Pub-EAA issuance, according to
	// the profile defined in ETSI TS 119 602, Annex H.
	CertificateApprovalStatusEnum_WITHDRAWN_CERT_FOR_PUB_EAA_ISSUANCE CertificateApprovalStatusEnum = "WITHDRAWN_CERT_FOR_PUB_EAA_ISSUANCE"
	// CertificateApprovalStatusEnum_WITHDRAWN_CERT_FOR_PUB_EAA_REVOCATION
	// represents a withdrawn certificate for Pub-EAA revocation, according to
	// the profile defined in ETSI TS 119 602, Annex H.
	CertificateApprovalStatusEnum_WITHDRAWN_CERT_FOR_PUB_EAA_REVOCATION CertificateApprovalStatusEnum = "WITHDRAWN_CERT_FOR_PUB_EAA_REVOCATION"
	// CertificateApprovalStatusEnum_CERT_FOR_REGISTER represents a
	// certificate for Register, according to the profile defined in ETSI TS
	// 119 602, Annex G.
	CertificateApprovalStatusEnum_CERT_FOR_REGISTER CertificateApprovalStatusEnum = "CERT_FOR_REGISTER"
	// CertificateApprovalStatusEnum_CERT_FOR_UNKNOWN represents a
	// certificate of unknown type (e.g. an error or conflict on
	// validation).
	CertificateApprovalStatusEnum_CERT_FOR_UNKNOWN CertificateApprovalStatusEnum = "CERT_FOR_UNKNOWN"
	// CertificateApprovalStatusEnum_NA is not applicable (e.g. no LoTE
	// trust anchor reached).
	CertificateApprovalStatusEnum_NA CertificateApprovalStatusEnum = "NA"
)

// certificateApprovalStatusEnumFields holds all per-constant attributes,
// mirroring the Java enum's instance fields.
type certificateApprovalStatusEnumFields struct {
	label    string
	listType ListType
	sti      LoTEServiceTypeIdentifier
	status   LoTEServiceStatus
}

// certificateApprovalStatusEnumData holds the full field tuple for each
// constant, copied verbatim from the Java enum constructors.
var certificateApprovalStatusEnumData = map[CertificateApprovalStatusEnum]certificateApprovalStatusEnumFields{
	CertificateApprovalStatusEnum_PID_PROVIDER: {
		"PID Provider", LoTETypeEnum_EUPIDProvidersList, LoTEServiceTypeIdentifierEnum_PID_ISSUANCE, nil,
	},
	CertificateApprovalStatusEnum_CERT_FOR_PID_REVOCATION: {
		"Certificate for PID Revocation", LoTETypeEnum_EUPIDProvidersList, LoTEServiceTypeIdentifierEnum_PID_REVOCATION, nil,
	},
	CertificateApprovalStatusEnum_CERT_FOR_WALLET_ISSUANCE: {
		"Certificate for Wallet Solution Issuance", LoTETypeEnum_EUWalletProvidersList, LoTEServiceTypeIdentifierEnum_WALLET_ISSUANCE, nil,
	},
	CertificateApprovalStatusEnum_CERT_FOR_WALLET_REVOCATION: {
		"Certificate for Wallet Solution Revocation", LoTETypeEnum_EUWalletProvidersList, LoTEServiceTypeIdentifierEnum_WALLET_REVOCATION, nil,
	},
	CertificateApprovalStatusEnum_CERT_FOR_WRPAC_ISSUANCE: {
		"Certificate for WRPAC Issuance", LoTETypeEnum_EUWRPACProvidersList, LoTEServiceTypeIdentifierEnum_WRPAC_ISSUANCE, nil,
	},
	CertificateApprovalStatusEnum_CERT_FOR_WRPAC_REVOCATION: {
		"Certificate for WRPAC Revocation", LoTETypeEnum_EUWRPACProvidersList, LoTEServiceTypeIdentifierEnum_WRPAC_REVOCATION, nil,
	},
	CertificateApprovalStatusEnum_CERT_FOR_WRPRC_ISSUANCE: {
		"Certificate for WRPRC Issuance", LoTETypeEnum_EUWRPRCProvidersList, LoTEServiceTypeIdentifierEnum_WRPRC_ISSUANCE, nil,
	},
	CertificateApprovalStatusEnum_CERT_FOR_WRPRC_REVOCATION: {
		"Certificate for WRPRC Revocation", LoTETypeEnum_EUWRPRCProvidersList, LoTEServiceTypeIdentifierEnum_WRPRC_REVOCATION, nil,
	},
	CertificateApprovalStatusEnum_NOTIFIED_CERT_FOR_PUB_EAA_ISSUANCE: {
		"Notified Certificate for Pub-EAA Issuance", LoTETypeEnum_EUPubEAAProvidersList, LoTEServiceTypeIdentifierEnum_PUB_EAA_ISSUANCE, LoTEServiceStatusEnum_PUB_EAA_PROVIDER_NOTIFIED,
	},
	CertificateApprovalStatusEnum_NOTIFIED_CERT_FOR_PUB_EAA_REVOCATION: {
		"Notified Certificate for Pub-EAA Revocation", LoTETypeEnum_EUPubEAAProvidersList, LoTEServiceTypeIdentifierEnum_PUB_EAA_REVOCATION, LoTEServiceStatusEnum_PUB_EAA_PROVIDER_NOTIFIED,
	},
	CertificateApprovalStatusEnum_WITHDRAWN_CERT_FOR_PUB_EAA_ISSUANCE: {
		// NOTE: the Java source reuses the "Notified Certificate for Pub-EAA
		// Issuance" label verbatim for this WITHDRAWN constant; preserved
		// as-is (likely an upstream copy/paste bug, not corrected here).
		"Notified Certificate for Pub-EAA Issuance", LoTETypeEnum_EUPubEAAProvidersList, LoTEServiceTypeIdentifierEnum_PUB_EAA_ISSUANCE, LoTEServiceStatusEnum_PUB_EAA_PROVIDER_WITHDRAWN,
	},
	CertificateApprovalStatusEnum_WITHDRAWN_CERT_FOR_PUB_EAA_REVOCATION: {
		// NOTE: same upstream label reuse as above, preserved verbatim.
		"Notified Certificate for Pub-EAA Revocation", LoTETypeEnum_EUPubEAAProvidersList, LoTEServiceTypeIdentifierEnum_PUB_EAA_REVOCATION, LoTEServiceStatusEnum_PUB_EAA_PROVIDER_WITHDRAWN,
	},
	CertificateApprovalStatusEnum_CERT_FOR_REGISTER: {
		"Certificate for Register", LoTETypeEnum_EURegistrarsAndRegistersList, LoTEServiceTypeIdentifierEnum_REGISTER, nil,
	},
	CertificateApprovalStatusEnum_CERT_FOR_UNKNOWN: {
		"Certificate for Unknown usage", nil, nil, nil,
	},
	CertificateApprovalStatusEnum_NA: {
		"Not applicable", nil, nil, nil,
	},
}

// CertificateApprovalStatusEnumValues returns all constants in declaration order.
func CertificateApprovalStatusEnumValues() []CertificateApprovalStatusEnum {
	return []CertificateApprovalStatusEnum{
		CertificateApprovalStatusEnum_PID_PROVIDER,
		CertificateApprovalStatusEnum_CERT_FOR_PID_REVOCATION,
		CertificateApprovalStatusEnum_CERT_FOR_WALLET_ISSUANCE,
		CertificateApprovalStatusEnum_CERT_FOR_WALLET_REVOCATION,
		CertificateApprovalStatusEnum_CERT_FOR_WRPAC_ISSUANCE,
		CertificateApprovalStatusEnum_CERT_FOR_WRPAC_REVOCATION,
		CertificateApprovalStatusEnum_CERT_FOR_WRPRC_ISSUANCE,
		CertificateApprovalStatusEnum_CERT_FOR_WRPRC_REVOCATION,
		CertificateApprovalStatusEnum_NOTIFIED_CERT_FOR_PUB_EAA_ISSUANCE,
		CertificateApprovalStatusEnum_NOTIFIED_CERT_FOR_PUB_EAA_REVOCATION,
		CertificateApprovalStatusEnum_WITHDRAWN_CERT_FOR_PUB_EAA_ISSUANCE,
		CertificateApprovalStatusEnum_WITHDRAWN_CERT_FOR_PUB_EAA_REVOCATION,
		CertificateApprovalStatusEnum_CERT_FOR_REGISTER,
		CertificateApprovalStatusEnum_CERT_FOR_UNKNOWN,
		CertificateApprovalStatusEnum_NA,
	}
}

// ListType returns the applicable List type. Implements CertificateApprovalStatus.
func (c CertificateApprovalStatusEnum) ListType() ListType {
	return certificateApprovalStatusEnumData[c].listType
}

// ServiceTypeIdentifier returns the applicable service type identifier.
// Implements CertificateApprovalStatus.
func (c CertificateApprovalStatusEnum) ServiceTypeIdentifier() LoTEServiceTypeIdentifier {
	return certificateApprovalStatusEnumData[c].sti
}

// ServiceStatus returns the applicable service status. Implements
// CertificateApprovalStatus.
func (c CertificateApprovalStatusEnum) ServiceStatus() LoTEServiceStatus {
	return certificateApprovalStatusEnumData[c].status
}

// Label returns the user-friendly certificate label. Implements
// CertificateApprovalStatus.
func (c CertificateApprovalStatusEnum) Label() string {
	return certificateApprovalStatusEnumData[c].label
}

// compile-time interface assertion.
var _ CertificateApprovalStatus = CertificateApprovalStatusEnum("")
