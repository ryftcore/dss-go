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
	// CertificateApprovalStatusEnumPIDProvider represents a PID provider
	// certificate, as defined in the ETSI TS 119 605, Annex C.1.
	CertificateApprovalStatusEnumPIDProvider CertificateApprovalStatusEnum = "PID_PROVIDER"
	// CertificateApprovalStatusEnumCertForPIDRevocation represents a
	// certificate for PID revocation, according to the profile defined in
	// ETSI TS 119 602, Annex C.1.
	CertificateApprovalStatusEnumCertForPIDRevocation CertificateApprovalStatusEnum = "CERT_FOR_PID_REVOCATION"
	// CertificateApprovalStatusEnumCertForWalletIssuance represents a
	// certificate for Wallet solution issuance, according to the profile
	// defined in ETSI TS 119 602, Annex C.2.
	CertificateApprovalStatusEnumCertForWalletIssuance CertificateApprovalStatusEnum = "CERT_FOR_WALLET_ISSUANCE"
	// CertificateApprovalStatusEnumCertForWalletRevocation represents a
	// certificate for Wallet solution revocation, according to the profile
	// defined in ETSI TS 119 602, Annex C.2.
	CertificateApprovalStatusEnumCertForWalletRevocation CertificateApprovalStatusEnum = "CERT_FOR_WALLET_REVOCATION"
	// CertificateApprovalStatusEnumCertForWRPACIssuance represents a
	// certificate for WRPAC issuance, according to the profile defined in
	// ETSI TS 119 602, Annex F.
	CertificateApprovalStatusEnumCertForWRPACIssuance CertificateApprovalStatusEnum = "CERT_FOR_WRPAC_ISSUANCE"
	// CertificateApprovalStatusEnumCertForWRPACRevocation represents a
	// certificate for WRPAC revocation, according to the profile defined in
	// ETSI TS 119 602, Annex F.
	CertificateApprovalStatusEnumCertForWRPACRevocation CertificateApprovalStatusEnum = "CERT_FOR_WRPAC_REVOCATION"
	// CertificateApprovalStatusEnumCertForWRPRCIssuance represents a
	// certificate for WRPRC issuance, according to the profile defined in
	// ETSI TS 119 602, Annex G.
	CertificateApprovalStatusEnumCertForWRPRCIssuance CertificateApprovalStatusEnum = "CERT_FOR_WRPRC_ISSUANCE"
	// CertificateApprovalStatusEnumCertForWRPRCRevocation represents a
	// certificate for WRPRC revocation, according to the profile defined in
	// ETSI TS 119 602, Annex G.
	CertificateApprovalStatusEnumCertForWRPRCRevocation CertificateApprovalStatusEnum = "CERT_FOR_WRPRC_REVOCATION"
	// CertificateApprovalStatusEnumNotifiedCertForPubEAAIssuance
	// represents a notified certificate for Pub-EAA issuance, according to
	// the profile defined in ETSI TS 119 602, Annex H.
	CertificateApprovalStatusEnumNotifiedCertForPubEAAIssuance CertificateApprovalStatusEnum = "NOTIFIED_CERT_FOR_PUB_EAA_ISSUANCE"
	// CertificateApprovalStatusEnumNotifiedCertForPubEAARevocation
	// represents a notified certificate for Pub-EAA revocation, according to
	// the profile defined in ETSI TS 119 602, Annex H.
	CertificateApprovalStatusEnumNotifiedCertForPubEAARevocation CertificateApprovalStatusEnum = "NOTIFIED_CERT_FOR_PUB_EAA_REVOCATION"
	// CertificateApprovalStatusEnumWithdrawnCertForPubEAAIssuance
	// represents a withdrawn certificate for Pub-EAA issuance, according to
	// the profile defined in ETSI TS 119 602, Annex H.
	CertificateApprovalStatusEnumWithdrawnCertForPubEAAIssuance CertificateApprovalStatusEnum = "WITHDRAWN_CERT_FOR_PUB_EAA_ISSUANCE"
	// CertificateApprovalStatusEnumWithdrawnCertForPubEAARevocation
	// represents a withdrawn certificate for Pub-EAA revocation, according to
	// the profile defined in ETSI TS 119 602, Annex H.
	CertificateApprovalStatusEnumWithdrawnCertForPubEAARevocation CertificateApprovalStatusEnum = "WITHDRAWN_CERT_FOR_PUB_EAA_REVOCATION"
	// CertificateApprovalStatusEnumCertForRegister represents a
	// certificate for Register, according to the profile defined in ETSI TS
	// 119 602, Annex G.
	CertificateApprovalStatusEnumCertForRegister CertificateApprovalStatusEnum = "CERT_FOR_REGISTER"
	// CertificateApprovalStatusEnumCertForUnknown represents a
	// certificate of unknown type (e.g. an error or conflict on
	// validation).
	CertificateApprovalStatusEnumCertForUnknown CertificateApprovalStatusEnum = "CERT_FOR_UNKNOWN"
	// CertificateApprovalStatusEnumNA is not applicable (e.g. no LoTE
	// trust anchor reached).
	CertificateApprovalStatusEnumNA CertificateApprovalStatusEnum = "NA"
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
	CertificateApprovalStatusEnumPIDProvider: {
		"PID Provider", LoTETypeEnumEUPIDProvidersList, LoTEServiceTypeIdentifierEnumPIDIssuance, nil,
	},
	CertificateApprovalStatusEnumCertForPIDRevocation: {
		"Certificate for PID Revocation", LoTETypeEnumEUPIDProvidersList, LoTEServiceTypeIdentifierEnumPIDRevocation, nil,
	},
	CertificateApprovalStatusEnumCertForWalletIssuance: {
		"Certificate for Wallet Solution Issuance", LoTETypeEnumEUWalletProvidersList, LoTEServiceTypeIdentifierEnumWalletIssuance, nil,
	},
	CertificateApprovalStatusEnumCertForWalletRevocation: {
		"Certificate for Wallet Solution Revocation", LoTETypeEnumEUWalletProvidersList, LoTEServiceTypeIdentifierEnumWalletRevocation, nil,
	},
	CertificateApprovalStatusEnumCertForWRPACIssuance: {
		"Certificate for WRPAC Issuance", LoTETypeEnumEUWRPACProvidersList, LoTEServiceTypeIdentifierEnumWRPACIssuance, nil,
	},
	CertificateApprovalStatusEnumCertForWRPACRevocation: {
		"Certificate for WRPAC Revocation", LoTETypeEnumEUWRPACProvidersList, LoTEServiceTypeIdentifierEnumWRPACRevocation, nil,
	},
	CertificateApprovalStatusEnumCertForWRPRCIssuance: {
		"Certificate for WRPRC Issuance", LoTETypeEnumEUWRPRCProvidersList, LoTEServiceTypeIdentifierEnumWRPRCIssuance, nil,
	},
	CertificateApprovalStatusEnumCertForWRPRCRevocation: {
		"Certificate for WRPRC Revocation", LoTETypeEnumEUWRPRCProvidersList, LoTEServiceTypeIdentifierEnumWRPRCRevocation, nil,
	},
	CertificateApprovalStatusEnumNotifiedCertForPubEAAIssuance: {
		"Notified Certificate for Pub-EAA Issuance", LoTETypeEnumEUPubEAAProvidersList, LoTEServiceTypeIdentifierEnumPubEAAIssuance, LoTEServiceStatusEnumPubEAAProviderNotified,
	},
	CertificateApprovalStatusEnumNotifiedCertForPubEAARevocation: {
		"Notified Certificate for Pub-EAA Revocation", LoTETypeEnumEUPubEAAProvidersList, LoTEServiceTypeIdentifierEnumPubEAARevocation, LoTEServiceStatusEnumPubEAAProviderNotified,
	},
	CertificateApprovalStatusEnumWithdrawnCertForPubEAAIssuance: {
		// NOTE: the Java source reuses the "Notified Certificate for Pub-EAA
		// Issuance" label verbatim for this WITHDRAWN constant; preserved
		// as-is (likely an upstream copy/paste bug, not corrected here).
		"Notified Certificate for Pub-EAA Issuance", LoTETypeEnumEUPubEAAProvidersList, LoTEServiceTypeIdentifierEnumPubEAAIssuance, LoTEServiceStatusEnumPubEAAProviderWithdrawn,
	},
	CertificateApprovalStatusEnumWithdrawnCertForPubEAARevocation: {
		// NOTE: same upstream label reuse as above, preserved verbatim.
		"Notified Certificate for Pub-EAA Revocation", LoTETypeEnumEUPubEAAProvidersList, LoTEServiceTypeIdentifierEnumPubEAARevocation, LoTEServiceStatusEnumPubEAAProviderWithdrawn,
	},
	CertificateApprovalStatusEnumCertForRegister: {
		"Certificate for Register", LoTETypeEnumEURegistrarsAndRegistersList, LoTEServiceTypeIdentifierEnumRegister, nil,
	},
	CertificateApprovalStatusEnumCertForUnknown: {
		"Certificate for Unknown usage", nil, nil, nil,
	},
	CertificateApprovalStatusEnumNA: {
		"Not applicable", nil, nil, nil,
	},
}

// CertificateApprovalStatusEnumValues returns all constants in declaration order.
func CertificateApprovalStatusEnumValues() []CertificateApprovalStatusEnum {
	return []CertificateApprovalStatusEnum{
		CertificateApprovalStatusEnumPIDProvider,
		CertificateApprovalStatusEnumCertForPIDRevocation,
		CertificateApprovalStatusEnumCertForWalletIssuance,
		CertificateApprovalStatusEnumCertForWalletRevocation,
		CertificateApprovalStatusEnumCertForWRPACIssuance,
		CertificateApprovalStatusEnumCertForWRPACRevocation,
		CertificateApprovalStatusEnumCertForWRPRCIssuance,
		CertificateApprovalStatusEnumCertForWRPRCRevocation,
		CertificateApprovalStatusEnumNotifiedCertForPubEAAIssuance,
		CertificateApprovalStatusEnumNotifiedCertForPubEAARevocation,
		CertificateApprovalStatusEnumWithdrawnCertForPubEAAIssuance,
		CertificateApprovalStatusEnumWithdrawnCertForPubEAARevocation,
		CertificateApprovalStatusEnumCertForRegister,
		CertificateApprovalStatusEnumCertForUnknown,
		CertificateApprovalStatusEnumNA,
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
