// Ported from dss-enumerations/.../SignatureQualification.java (DSS 6.5.RC1).
package enumerations

// SignatureQualification defines available signature qualification types.
// Implements UriBasedEnum.
type SignatureQualification string

const (
	// SignatureQualification_QESIG is a Qualified Electronic Signature.
	SignatureQualification_QESIG SignatureQualification = "QESIG"
	// SignatureQualification_QESEAL is a Qualified Electronic Seal.
	SignatureQualification_QESEAL SignatureQualification = "QESEAL"
	// SignatureQualification_UNKNOWN_QC_QSCD is a signature supported by a
	// Qualified Certificate with the private key in a QSCD.
	SignatureQualification_UNKNOWN_QC_QSCD SignatureQualification = "UNKNOWN_QC_QSCD"
	// SignatureQualification_ADESIG_QC is an Advanced Electronic Signature
	// supported by a Qualified Certificate.
	SignatureQualification_ADESIG_QC SignatureQualification = "ADESIG_QC"
	// SignatureQualification_ADESEAL_QC is an Advanced Electronic Seal
	// supported by a Qualified Certificate.
	SignatureQualification_ADESEAL_QC SignatureQualification = "ADESEAL_QC"
	// SignatureQualification_UNKNOWN_QC is a signature supported by a
	// Qualified Certificate.
	SignatureQualification_UNKNOWN_QC SignatureQualification = "UNKNOWN_QC"
	// SignatureQualification_ADESIG is an Advanced Electronic Signature.
	SignatureQualification_ADESIG SignatureQualification = "ADESIG"
	// SignatureQualification_ADESEAL is an Advanced Electronic Seal.
	SignatureQualification_ADESEAL SignatureQualification = "ADESEAL"
	// SignatureQualification_UNKNOWN is a signature of unknown type.
	SignatureQualification_UNKNOWN SignatureQualification = "UNKNOWN"
	// SignatureQualification_INDETERMINATE_QESIG is an Indeterminate
	// Qualified Electronic Signature.
	SignatureQualification_INDETERMINATE_QESIG SignatureQualification = "INDETERMINATE_QESIG"
	// SignatureQualification_INDETERMINATE_QESEAL is an Indeterminate
	// Qualified Electronic Seal.
	SignatureQualification_INDETERMINATE_QESEAL SignatureQualification = "INDETERMINATE_QESEAL"
	// SignatureQualification_INDETERMINATE_UNKNOWN_QC_QSCD is an
	// Indeterminate signature supported by a Qualified Certificate with
	// the private key in a QSCD.
	SignatureQualification_INDETERMINATE_UNKNOWN_QC_QSCD SignatureQualification = "INDETERMINATE_UNKNOWN_QC_QSCD"
	// SignatureQualification_INDETERMINATE_ADESIG_QC is an Indeterminate
	// Advanced Electronic Signature supported by a Qualified Certificate.
	SignatureQualification_INDETERMINATE_ADESIG_QC SignatureQualification = "INDETERMINATE_ADESIG_QC"
	// SignatureQualification_INDETERMINATE_ADESEAL_QC is an Indeterminate
	// Advanced Electronic Seal supported by a Qualified Certificate.
	SignatureQualification_INDETERMINATE_ADESEAL_QC SignatureQualification = "INDETERMINATE_ADESEAL_QC"
	// SignatureQualification_INDETERMINATE_UNKNOWN_QC is an Indeterminate
	// signature supported by a Qualified Certificate.
	SignatureQualification_INDETERMINATE_UNKNOWN_QC SignatureQualification = "INDETERMINATE_UNKNOWN_QC"
	// SignatureQualification_INDETERMINATE_ADESIG is an Indeterminate
	// Advanced Electronic Signature.
	SignatureQualification_INDETERMINATE_ADESIG SignatureQualification = "INDETERMINATE_ADESIG"
	// SignatureQualification_INDETERMINATE_ADESEAL is an Indeterminate
	// Advanced Electronic Seal.
	SignatureQualification_INDETERMINATE_ADESEAL SignatureQualification = "INDETERMINATE_ADESEAL"
	// SignatureQualification_INDETERMINATE_UNKNOWN is a signature of
	// unknown type.
	SignatureQualification_INDETERMINATE_UNKNOWN SignatureQualification = "INDETERMINATE_UNKNOWN"
	// SignatureQualification_NOT_ADES_QC_QSCD is Not Advanced Electronic
	// Signature but supported by a Qualified Certificate.
	SignatureQualification_NOT_ADES_QC_QSCD SignatureQualification = "NOT_ADES_QC_QSCD"
	// SignatureQualification_NOT_ADES_QC is Not Advanced Electronic
	// Signature but supported by a Qualified Certificate.
	SignatureQualification_NOT_ADES_QC SignatureQualification = "NOT_ADES_QC"
	// SignatureQualification_NOT_ADES is Not Advanced Electronic Signature.
	SignatureQualification_NOT_ADES SignatureQualification = "NOT_ADES"
	// SignatureQualification_NA is Not Applicable.
	SignatureQualification_NA SignatureQualification = "NA"
)

type signatureQualificationFields struct {
	readable string
	label    string
	uri      string
}

// signatureQualificationData holds the (readable, label, uri) tuple for
// each constant.
var signatureQualificationData = map[SignatureQualification]signatureQualificationFields{
	SignatureQualification_QESIG:  {"QESig", "Qualified Electronic Signature", "urn:cef:dss:signatureQualification:QESig"},
	SignatureQualification_QESEAL: {"QESeal", "Qualified Electronic Seal", "urn:cef:dss:signatureQualification:QESeal"},
	SignatureQualification_UNKNOWN_QC_QSCD: {
		"Unknown-QC-QSCD",
		"Signature produced by a Qualified Certificate of unknown type with its private key residing in a QSCD",
		"urn:cef:dss:signatureQualification:UnknownQCwithQSCD",
	},
	SignatureQualification_ADESIG_QC:            {"AdESig-QC", "Advanced Electronic Signature supported by a Qualified Certificate", "urn:cef:dss:signatureQualification:AdESigQC"},
	SignatureQualification_ADESEAL_QC:           {"AdESeal-QC", "Advanced Electronic Seal supported by a Qualified Certificate", "urn:cef:dss:signatureQualification:AdESealQC"},
	SignatureQualification_UNKNOWN_QC:           {"Unknown-QC", "Signature produced by a Qualified Certificate of unknown type", "urn:cef:dss:signatureQualification:UnknownQC"},
	SignatureQualification_ADESIG:               {"AdESig", "Advanced Electronic Signature", "urn:cef:dss:signatureQualification:AdESig"},
	SignatureQualification_ADESEAL:              {"AdESeal", "Advanced Electronic Seal", "urn:cef:dss:signatureQualification:AdESeal"},
	SignatureQualification_UNKNOWN:              {"Unknown", "Signature produced by a Certificate of unknown type", "urn:cef:dss:signatureQualification:Unknown"},
	SignatureQualification_INDETERMINATE_QESIG:  {"Indeterminate QESig", "Indeterminate Qualified Electronic Signature", "urn:cef:dss:signatureQualification:indeterminateQESig"},
	SignatureQualification_INDETERMINATE_QESEAL: {"Indeterminate QESeal", "Indeterminate Qualified Electronic Seal", "urn:cef:dss:signatureQualification:indeterminateQESeal"},
	SignatureQualification_INDETERMINATE_UNKNOWN_QC_QSCD: {
		"Indeterminate Unknown-QC-QSCD",
		"Indeterminate Signature produced by a Qualified Certificate of unknown type with its private key residing in a QSCD",
		"urn:cef:dss:signatureQualification:indeterminateUnknownQCwithQSCD",
	},
	SignatureQualification_INDETERMINATE_ADESIG_QC: {
		"Indeterminate AdESig-QC",
		"Indeterminate Advanced Electronic Signature supported by a Qualified Certificate",
		"urn:cef:dss:signatureQualification:indeterminateAdESigQC",
	},
	SignatureQualification_INDETERMINATE_ADESEAL_QC: {
		"Indeterminate AdESeal-QC",
		"Indeterminate Advanced Electronic Seal supported by a Qualified Certificate",
		"urn:cef:dss:signatureQualification:indeterminateAdESealQC",
	},
	SignatureQualification_INDETERMINATE_UNKNOWN_QC: {
		"Indeterminate Unknown-QC",
		"Indeterminate Signature produced by a Qualified Certificate of unknown type",
		"urn:cef:dss:signatureQualification:indeterminateUnknownQC",
	},
	SignatureQualification_INDETERMINATE_ADESIG: {
		"Indeterminate AdESig", "Indeterminate Advanced Electronic Signature", "urn:cef:dss:signatureQualification:indeterminateAdESig",
	},
	SignatureQualification_INDETERMINATE_ADESEAL: {
		"Indeterminate AdESeal", "Indeterminate Advanced Electronic Seal", "urn:cef:dss:signatureQualification:indeterminateAdESeal",
	},
	SignatureQualification_INDETERMINATE_UNKNOWN: {
		"Indeterminate Unknown", "Indeterminate Signature produced by a Certificate of unknown type", "urn:cef:dss:signatureQualification:indeterminateUnknown",
	},
	SignatureQualification_NOT_ADES_QC_QSCD: {
		"Not AdES but QC with QSCD", "Not Advanced Electronic Signature but supported by a Qualified Certificate", "urn:cef:dss:signatureQualification:notAdESbutQCwithQSCD",
	},
	SignatureQualification_NOT_ADES_QC: {
		"Not AdES but QC", "Not Advanced Electronic Signature but supported by a Qualified Certificate", "urn:cef:dss:signatureQualification:notAdESbutQC",
	},
	SignatureQualification_NOT_ADES: {"Not AdES", "Not Advanced Electronic Signature", "urn:cef:dss:signatureQualification:notAdES"},
	SignatureQualification_NA:       {"N/A", "Not applicable", "urn:cef:dss:signatureQualification:notApplicable"},
}

// SignatureQualificationValues returns all constants in declaration order.
func SignatureQualificationValues() []SignatureQualification {
	return []SignatureQualification{
		SignatureQualification_QESIG,
		SignatureQualification_QESEAL,
		SignatureQualification_UNKNOWN_QC_QSCD,
		SignatureQualification_ADESIG_QC,
		SignatureQualification_ADESEAL_QC,
		SignatureQualification_UNKNOWN_QC,
		SignatureQualification_ADESIG,
		SignatureQualification_ADESEAL,
		SignatureQualification_UNKNOWN,
		SignatureQualification_INDETERMINATE_QESIG,
		SignatureQualification_INDETERMINATE_QESEAL,
		SignatureQualification_INDETERMINATE_UNKNOWN_QC_QSCD,
		SignatureQualification_INDETERMINATE_ADESIG_QC,
		SignatureQualification_INDETERMINATE_ADESEAL_QC,
		SignatureQualification_INDETERMINATE_UNKNOWN_QC,
		SignatureQualification_INDETERMINATE_ADESIG,
		SignatureQualification_INDETERMINATE_ADESEAL,
		SignatureQualification_INDETERMINATE_UNKNOWN,
		SignatureQualification_NOT_ADES_QC_QSCD,
		SignatureQualification_NOT_ADES_QC,
		SignatureQualification_NOT_ADES,
		SignatureQualification_NA,
	}
}

// signatureQualificationByReadable is a lookup map from the readable
// short-string to the SignatureQualification constant.
var signatureQualificationByReadable = func() map[string]SignatureQualification {
	m := make(map[string]SignatureQualification, len(signatureQualificationData))
	for _, v := range SignatureQualificationValues() {
		m[signatureQualificationData[v].readable] = v
	}
	return m
}()

// Readable gets user-friendly name of the enumeration.
func (s SignatureQualification) Readable() string {
	return signatureQualificationData[s].readable
}

// Label gets description of the enumeration.
func (s SignatureQualification) Label() string {
	return signatureQualificationData[s].label
}

// URI returns the URI. Implements UriBasedEnum.
func (s SignatureQualification) URI() string {
	return signatureQualificationData[s].uri
}

// SignatureQualificationForName converts the given qualification name into
// the enum. Returns "" (zero value) and no error if value is empty,
// mirroring Java's null return; an unknown name is an error, mirroring
// Java's valueOf IllegalArgumentException.
func SignatureQualificationForName(value string) (SignatureQualification, error) {
	if value == "" {
		return "", nil
	}
	return SignatureQualificationValueOf(value)
}

// SignatureQualificationValueOf returns the constant matching the given
// Java enum name.
func SignatureQualificationValueOf(name string) (SignatureQualification, error) {
	for _, v := range SignatureQualificationValues() {
		if string(v) == name {
			return v, nil
		}
	}
	return "", &signatureQualificationInvalidValueError{name}
}

type signatureQualificationInvalidValueError struct {
	name string
}

func (e *signatureQualificationInvalidValueError) Error() string {
	return "no enum constant SignatureQualification." + e.name
}

// SignatureQualificationFromReadable converts the given readable
// description into the enum. Returns "" (zero value) if readable is empty
// or unknown, mirroring Java's null return.
func SignatureQualificationFromReadable(readable string) SignatureQualification {
	if readable == "" {
		return ""
	}
	return signatureQualificationByReadable[readable]
}

// SignatureQualificationForURI converts the given uri into the enum.
// Returns "" (zero value) if uri is empty or unknown, mirroring Java's null
// return.
func SignatureQualificationForURI(uri string) SignatureQualification {
	if uri == "" {
		return ""
	}
	for _, v := range SignatureQualificationValues() {
		if signatureQualificationData[v].uri == uri {
			return v
		}
	}
	return ""
}

// compile-time interface assertion.
var _ UriBasedEnum = SignatureQualification("")
