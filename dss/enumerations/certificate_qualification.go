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
	// CertificateQualificationQCERTForESigQSCD is a Qualified Certificate
	// for Electronic Signatures with private key on QSCD.
	CertificateQualificationQCERTForESigQSCD CertificateQualification = "QCERT_FOR_ESIG_QSCD"
	// CertificateQualificationQCERTForESealQSCD is a Qualified
	// Certificate for Electronic Seals with private key on QSCD.
	CertificateQualificationQCERTForESealQSCD CertificateQualification = "QCERT_FOR_ESEAL_QSCD"
	// CertificateQualificationQCERTForUnknownQSCD is a Qualified
	// Certificate for unidentified type with private key in a QSCD.
	CertificateQualificationQCERTForUnknownQSCD CertificateQualification = "QCERT_FOR_UNKNOWN_QSCD"
	// CertificateQualificationQCERTForESig is a Qualified Certificate
	// for Electronic Signatures.
	CertificateQualificationQCERTForESig CertificateQualification = "QCERT_FOR_ESIG"
	// CertificateQualificationQCERTForESeal is a Qualified Certificate
	// for Electronic Seals.
	CertificateQualificationQCERTForESeal CertificateQualification = "QCERT_FOR_ESEAL"
	// CertificateQualificationQCERTForWSA is a Qualified Certificate
	// for Web Site Authentications.
	CertificateQualificationQCERTForWSA CertificateQualification = "QCERT_FOR_WSA"
	// CertificateQualificationQCERTForUnknown is a Qualified Certificate
	// for unidentified type.
	CertificateQualificationQCERTForUnknown CertificateQualification = "QCERT_FOR_UNKNOWN"
	// CertificateQualificationCertForESig is a Certificate for
	// Electronic Signatures.
	CertificateQualificationCertForESig CertificateQualification = "CERT_FOR_ESIG"
	// CertificateQualificationCertForESeal is a Certificate for
	// Electronic Seals.
	CertificateQualificationCertForESeal CertificateQualification = "CERT_FOR_ESEAL"
	// CertificateQualificationCertForWSA is a Certificate for Web Site
	// Authentications.
	CertificateQualificationCertForWSA CertificateQualification = "CERT_FOR_WSA"
	// CertificateQualificationCertForUnknown is a Certificate for
	// unidentified type.
	CertificateQualificationCertForUnknown CertificateQualification = "CERT_FOR_UNKNOWN"
	// CertificateQualificationNA is Not Applicable.
	CertificateQualificationNA CertificateQualification = "NA"
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
	CertificateQualificationQCERTForESigQSCD: {
		"QC for eSig with QSCD", "Qualified Certificate for Electronic Signatures with private key on QSCD",
		CertificateQualifiedStatusQC, CertificateTypeESign, QSCDStatusQSCD,
	},
	CertificateQualificationQCERTForESealQSCD: {
		"QC for eSeal with QSCD", "Qualified Certificate for Electronic Seals with private key on QSCD",
		CertificateQualifiedStatusQC, CertificateTypeESeal, QSCDStatusQSCD,
	},
	CertificateQualificationQCERTForUnknownQSCD: {
		"QC for unknown type with QSCD", "Qualified Certificate for unknown type with its private key residing in a QSCD",
		CertificateQualifiedStatusQC, CertificateTypeUnknown, QSCDStatusQSCD,
	},
	CertificateQualificationQCERTForESig: {
		"QC for eSig", "Qualified Certificate for Electronic Signatures",
		CertificateQualifiedStatusQC, CertificateTypeESign, QSCDStatusNotQSCD,
	},
	CertificateQualificationQCERTForESeal: {
		"QC for eSeal", "Qualified Certificate for Electronic Seals",
		CertificateQualifiedStatusQC, CertificateTypeESeal, QSCDStatusNotQSCD,
	},
	CertificateQualificationQCERTForWSA: {
		"QC for WSA", "Qualified Certificate for Web Site Authentications",
		CertificateQualifiedStatusQC, CertificateTypeWSA, QSCDStatusNotQSCD,
	},
	CertificateQualificationQCERTForUnknown: {
		"QC for unknown type", "Qualified Certificate for unknown type",
		CertificateQualifiedStatusQC, CertificateTypeUnknown, QSCDStatusNotQSCD,
	},
	CertificateQualificationCertForESig: {
		"Cert for eSig", "Certificate for Electronic Signatures",
		CertificateQualifiedStatusNotQC, CertificateTypeESign, QSCDStatusNotQSCD,
	},
	CertificateQualificationCertForESeal: {
		"Cert for eSeal", "Certificate for Electronic Seals",
		CertificateQualifiedStatusNotQC, CertificateTypeESeal, QSCDStatusNotQSCD,
	},
	CertificateQualificationCertForWSA: {
		"Cert for WSA", "Certificate for Web Site Authentications",
		CertificateQualifiedStatusNotQC, CertificateTypeWSA, QSCDStatusNotQSCD,
	},
	CertificateQualificationCertForUnknown: {
		"Cert for unknown type", "Certificate for unknown type",
		CertificateQualifiedStatusNotQC, CertificateTypeUnknown, QSCDStatusNotQSCD,
	},
	CertificateQualificationNA: {
		"N/A", "Not applicable",
		CertificateQualifiedStatusNotQC, CertificateTypeUnknown, QSCDStatusNotQSCD,
	},
}

// CertificateQualificationValues returns all constants in declaration order.
func CertificateQualificationValues() []CertificateQualification {
	return []CertificateQualification{
		CertificateQualificationQCERTForESigQSCD,
		CertificateQualificationQCERTForESealQSCD,
		CertificateQualificationQCERTForUnknownQSCD,
		CertificateQualificationQCERTForESig,
		CertificateQualificationQCERTForESeal,
		CertificateQualificationQCERTForWSA,
		CertificateQualificationQCERTForUnknown,
		CertificateQualificationCertForESig,
		CertificateQualificationCertForESeal,
		CertificateQualificationCertForWSA,
		CertificateQualificationCertForUnknown,
		CertificateQualificationNA,
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
	return CertificateTypeESign == certificateQualificationData[c].certType
}

// IsForEseal returns if the certificate can be used for an electronic seal.
func (c CertificateQualification) IsForEseal() bool {
	return CertificateTypeESeal == certificateQualificationData[c].certType
}

// IsQscd returns if the certificate is used on a Qualified Signature
// Creation Device.
func (c CertificateQualification) IsQscd() bool {
	return QSCDStatusQSCD == certificateQualificationData[c].qscdStatus
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
