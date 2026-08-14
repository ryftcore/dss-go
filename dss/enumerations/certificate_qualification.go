// Ported from dss-enumerations/.../CertificateQualification.java (DSS 6.5.RC1).
//
// CertificateType and QSCDStatus are defined outside this file's manifest
// and are assumed to exist per the porting brief (type X string with
// X_<NAME> constants and an IsQSCD-style helper).
package enumerations

// CertificateQualification represents the available certificate
// qualification types.
type CertificateQualification string

const (
	// CertificateQualification_QCERT_FOR_ESIG_QSCD is a Qualified Certificate
	// for Electronic Signatures with private key on QSCD.
	CertificateQualification_QCERT_FOR_ESIG_QSCD CertificateQualification = "QCERT_FOR_ESIG_QSCD"
	// CertificateQualification_QCERT_FOR_ESEAL_QSCD is a Qualified
	// Certificate for Electronic Seals with private key on QSCD.
	CertificateQualification_QCERT_FOR_ESEAL_QSCD CertificateQualification = "QCERT_FOR_ESEAL_QSCD"
	// CertificateQualification_QCERT_FOR_UNKNOWN_QSCD is a Qualified
	// Certificate for unidentified type with private key in a QSCD.
	CertificateQualification_QCERT_FOR_UNKNOWN_QSCD CertificateQualification = "QCERT_FOR_UNKNOWN_QSCD"
	// CertificateQualification_QCERT_FOR_ESIG is a Qualified Certificate
	// for Electronic Signatures.
	CertificateQualification_QCERT_FOR_ESIG CertificateQualification = "QCERT_FOR_ESIG"
	// CertificateQualification_QCERT_FOR_ESEAL is a Qualified Certificate
	// for Electronic Seals.
	CertificateQualification_QCERT_FOR_ESEAL CertificateQualification = "QCERT_FOR_ESEAL"
	// CertificateQualification_QCERT_FOR_WSA is a Qualified Certificate
	// for Web Site Authentications.
	CertificateQualification_QCERT_FOR_WSA CertificateQualification = "QCERT_FOR_WSA"
	// CertificateQualification_QCERT_FOR_UNKNOWN is a Qualified Certificate
	// for unidentified type.
	CertificateQualification_QCERT_FOR_UNKNOWN CertificateQualification = "QCERT_FOR_UNKNOWN"
	// CertificateQualification_CERT_FOR_ESIG is a Certificate for
	// Electronic Signatures.
	CertificateQualification_CERT_FOR_ESIG CertificateQualification = "CERT_FOR_ESIG"
	// CertificateQualification_CERT_FOR_ESEAL is a Certificate for
	// Electronic Seals.
	CertificateQualification_CERT_FOR_ESEAL CertificateQualification = "CERT_FOR_ESEAL"
	// CertificateQualification_CERT_FOR_WSA is a Certificate for Web Site
	// Authentications.
	CertificateQualification_CERT_FOR_WSA CertificateQualification = "CERT_FOR_WSA"
	// CertificateQualification_CERT_FOR_UNKNOWN is a Certificate for
	// unidentified type.
	CertificateQualification_CERT_FOR_UNKNOWN CertificateQualification = "CERT_FOR_UNKNOWN"
	// CertificateQualification_NA is Not Applicable.
	CertificateQualification_NA CertificateQualification = "NA"
)

type certificateQualificationFields struct {
	readable        string
	label           string
	qualifiedStatus CertificateQualifiedStatus
	certType        CertificateType
	qscdStatus      QSCDStatus
}

// certificateQualificationData holds the (readable, label, qualifiedStatus,
// type, qscdStatus) tuple for each constant.
var certificateQualificationData = map[CertificateQualification]certificateQualificationFields{
	CertificateQualification_QCERT_FOR_ESIG_QSCD: {
		"QC for eSig with QSCD", "Qualified Certificate for Electronic Signatures with private key on QSCD",
		CertificateQualifiedStatus_QC, CertificateType_ESIGN, QSCDStatus_QSCD,
	},
	CertificateQualification_QCERT_FOR_ESEAL_QSCD: {
		"QC for eSeal with QSCD", "Qualified Certificate for Electronic Seals with private key on QSCD",
		CertificateQualifiedStatus_QC, CertificateType_ESEAL, QSCDStatus_QSCD,
	},
	CertificateQualification_QCERT_FOR_UNKNOWN_QSCD: {
		"QC for unknown type with QSCD", "Qualified Certificate for unknown type with its private key residing in a QSCD",
		CertificateQualifiedStatus_QC, CertificateType_UNKNOWN, QSCDStatus_QSCD,
	},
	CertificateQualification_QCERT_FOR_ESIG: {
		"QC for eSig", "Qualified Certificate for Electronic Signatures",
		CertificateQualifiedStatus_QC, CertificateType_ESIGN, QSCDStatus_NOT_QSCD,
	},
	CertificateQualification_QCERT_FOR_ESEAL: {
		"QC for eSeal", "Qualified Certificate for Electronic Seals",
		CertificateQualifiedStatus_QC, CertificateType_ESEAL, QSCDStatus_NOT_QSCD,
	},
	CertificateQualification_QCERT_FOR_WSA: {
		"QC for WSA", "Qualified Certificate for Web Site Authentications",
		CertificateQualifiedStatus_QC, CertificateType_WSA, QSCDStatus_NOT_QSCD,
	},
	CertificateQualification_QCERT_FOR_UNKNOWN: {
		"QC for unknown type", "Qualified Certificate for unknown type",
		CertificateQualifiedStatus_QC, CertificateType_UNKNOWN, QSCDStatus_NOT_QSCD,
	},
	CertificateQualification_CERT_FOR_ESIG: {
		"Cert for eSig", "Certificate for Electronic Signatures",
		CertificateQualifiedStatus_NOT_QC, CertificateType_ESIGN, QSCDStatus_NOT_QSCD,
	},
	CertificateQualification_CERT_FOR_ESEAL: {
		"Cert for eSeal", "Certificate for Electronic Seals",
		CertificateQualifiedStatus_NOT_QC, CertificateType_ESEAL, QSCDStatus_NOT_QSCD,
	},
	CertificateQualification_CERT_FOR_WSA: {
		"Cert for WSA", "Certificate for Web Site Authentications",
		CertificateQualifiedStatus_NOT_QC, CertificateType_WSA, QSCDStatus_NOT_QSCD,
	},
	CertificateQualification_CERT_FOR_UNKNOWN: {
		"Cert for unknown type", "Certificate for unknown type",
		CertificateQualifiedStatus_NOT_QC, CertificateType_UNKNOWN, QSCDStatus_NOT_QSCD,
	},
	CertificateQualification_NA: {
		"N/A", "Not applicable",
		CertificateQualifiedStatus_NOT_QC, CertificateType_UNKNOWN, QSCDStatus_NOT_QSCD,
	},
}

// CertificateQualificationValues returns all constants in declaration order.
func CertificateQualificationValues() []CertificateQualification {
	return []CertificateQualification{
		CertificateQualification_QCERT_FOR_ESIG_QSCD,
		CertificateQualification_QCERT_FOR_ESEAL_QSCD,
		CertificateQualification_QCERT_FOR_UNKNOWN_QSCD,
		CertificateQualification_QCERT_FOR_ESIG,
		CertificateQualification_QCERT_FOR_ESEAL,
		CertificateQualification_QCERT_FOR_WSA,
		CertificateQualification_QCERT_FOR_UNKNOWN,
		CertificateQualification_CERT_FOR_ESIG,
		CertificateQualification_CERT_FOR_ESEAL,
		CertificateQualification_CERT_FOR_WSA,
		CertificateQualification_CERT_FOR_UNKNOWN,
		CertificateQualification_NA,
	}
}

// certificateQualificationByReadable is a lookup map from the readable
// short-string to the CertificateQualification constant.
var certificateQualificationByReadable = func() map[string]CertificateQualification {
	m := make(map[string]CertificateQualification, len(certificateQualificationData))
	for _, v := range CertificateQualificationValues() {
		m[certificateQualificationData[v].readable] = v
	}
	return m
}()

// Readable returns a short string defining the qualification type.
func (c CertificateQualification) Readable() string {
	return certificateQualificationData[c].readable
}

// Label returns a complete name of the qualification type.
func (c CertificateQualification) Label() string {
	return certificateQualificationData[c].label
}

// Type returns the type of an electronic signature the certificate can be
// used for.
func (c CertificateQualification) Type() CertificateType {
	return certificateQualificationData[c].certType
}

// IsQc returns if the certificate is qualified.
func (c CertificateQualification) IsQc() bool {
	return CertificateQualifiedStatusIsQC(certificateQualificationData[c].qualifiedStatus)
}

// IsForEsig returns if the certificate can be used for an electronic signature.
func (c CertificateQualification) IsForEsig() bool {
	return CertificateType_ESIGN == certificateQualificationData[c].certType
}

// IsForEseal returns if the certificate can be used for an electronic seal.
func (c CertificateQualification) IsForEseal() bool {
	return CertificateType_ESEAL == certificateQualificationData[c].certType
}

// IsQscd returns if the certificate is used on a Qualified Signature
// Creation Device.
func (c CertificateQualification) IsQscd() bool {
	return QSCDStatus_QSCD == certificateQualificationData[c].qscdStatus
}

// CertificateQualificationForName converts the given qualification name to
// the enum. Returns "" (zero value) and no error if value is empty,
// mirroring Java's null return; an unknown name is an error, mirroring
// Java's valueOf IllegalArgumentException.
func CertificateQualificationForName(value string) (CertificateQualification, error) {
	if value == "" {
		return "", nil
	}
	return CertificateQualificationValueOf(value)
}

// CertificateQualificationValueOf returns the constant matching the given
// Java enum name.
func CertificateQualificationValueOf(name string) (CertificateQualification, error) {
	for _, v := range CertificateQualificationValues() {
		if string(v) == name {
			return v, nil
		}
	}
	return "", &certificateQualificationInvalidValueError{name}
}

type certificateQualificationInvalidValueError struct {
	name string
}

func (e *certificateQualificationInvalidValueError) Error() string {
	return "no enum constant CertificateQualification." + e.name
}

// CertificateQualificationFromReadable converts the given readable
// description of the qualification to the enum. Returns "" (zero value) if
// readable is empty or unknown, mirroring Java's null return.
func CertificateQualificationFromReadable(readable string) CertificateQualification {
	if readable == "" {
		return ""
	}
	return certificateQualificationByReadable[readable]
}
