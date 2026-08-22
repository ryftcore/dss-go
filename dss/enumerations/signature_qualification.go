// Ported from dss-enumerations/.../SignatureQualification.java (DSS 6.5.RC1).
package enumerations

// SignatureQualification defines available signature qualification types.
// Implements UriBasedEnum.
type SignatureQualification string

const (
	// SignatureQualificationQESig is a Qualified Electronic Signature.
	SignatureQualificationQESig SignatureQualification = "QESIG"
	// SignatureQualificationQESeal is a Qualified Electronic Seal.
	SignatureQualificationQESeal SignatureQualification = "QESEAL"
	// SignatureQualificationUnknownQCQSCD is a signature supported by a
	// Qualified Certificate with the private key in a QSCD.
	SignatureQualificationUnknownQCQSCD SignatureQualification = "UNKNOWN_QC_QSCD"
	// SignatureQualificationAdESigQC is an Advanced Electronic Signature
	// supported by a Qualified Certificate.
	SignatureQualificationAdESigQC SignatureQualification = "ADESIG_QC"
	// SignatureQualificationAdESealQC is an Advanced Electronic Seal
	// supported by a Qualified Certificate.
	SignatureQualificationAdESealQC SignatureQualification = "ADESEAL_QC"
	// SignatureQualificationUnknownQC is a signature supported by a
	// Qualified Certificate.
	SignatureQualificationUnknownQC SignatureQualification = "UNKNOWN_QC"
	// SignatureQualificationAdESig is an Advanced Electronic Signature.
	SignatureQualificationAdESig SignatureQualification = "ADESIG"
	// SignatureQualificationAdESeal is an Advanced Electronic Seal.
	SignatureQualificationAdESeal SignatureQualification = "ADESEAL"
	// SignatureQualificationUnknown is a signature of unknown type.
	SignatureQualificationUnknown SignatureQualification = "UNKNOWN"
	// SignatureQualificationIndeterminateQESig is an Indeterminate
	// Qualified Electronic Signature.
	SignatureQualificationIndeterminateQESig SignatureQualification = "INDETERMINATE_QESIG"
	// SignatureQualificationIndeterminateQESeal is an Indeterminate
	// Qualified Electronic Seal.
	SignatureQualificationIndeterminateQESeal SignatureQualification = "INDETERMINATE_QESEAL"
	// SignatureQualificationIndeterminateUnknownQCQSCD is an
	// Indeterminate signature supported by a Qualified Certificate with
	// the private key in a QSCD.
	SignatureQualificationIndeterminateUnknownQCQSCD SignatureQualification = "INDETERMINATE_UNKNOWN_QC_QSCD"
	// SignatureQualificationIndeterminateAdESigQC is an Indeterminate
	// Advanced Electronic Signature supported by a Qualified Certificate.
	SignatureQualificationIndeterminateAdESigQC SignatureQualification = "INDETERMINATE_ADESIG_QC"
	// SignatureQualificationIndeterminateAdESealQC is an Indeterminate
	// Advanced Electronic Seal supported by a Qualified Certificate.
	SignatureQualificationIndeterminateAdESealQC SignatureQualification = "INDETERMINATE_ADESEAL_QC"
	// SignatureQualificationIndeterminateUnknownQC is an Indeterminate
	// signature supported by a Qualified Certificate.
	SignatureQualificationIndeterminateUnknownQC SignatureQualification = "INDETERMINATE_UNKNOWN_QC"
	// SignatureQualificationIndeterminateAdESig is an Indeterminate
	// Advanced Electronic Signature.
	SignatureQualificationIndeterminateAdESig SignatureQualification = "INDETERMINATE_ADESIG"
	// SignatureQualificationIndeterminateAdESeal is an Indeterminate
	// Advanced Electronic Seal.
	SignatureQualificationIndeterminateAdESeal SignatureQualification = "INDETERMINATE_ADESEAL"
	// SignatureQualificationIndeterminateUnknown is a signature of
	// unknown type.
	SignatureQualificationIndeterminateUnknown SignatureQualification = "INDETERMINATE_UNKNOWN"
	// SignatureQualificationNotAdESQCQSCD is Not Advanced Electronic
	// Signature but supported by a Qualified Certificate.
	SignatureQualificationNotAdESQCQSCD SignatureQualification = "NOT_ADES_QC_QSCD"
	// SignatureQualificationNotAdESQC is Not Advanced Electronic
	// Signature but supported by a Qualified Certificate.
	SignatureQualificationNotAdESQC SignatureQualification = "NOT_ADES_QC"
	// SignatureQualificationNotAdES is Not Advanced Electronic Signature.
	SignatureQualificationNotAdES SignatureQualification = "NOT_ADES"
	// SignatureQualificationNA is Not Applicable.
	SignatureQualificationNA SignatureQualification = "NA"
)

type signatureQualificationFields struct {
	readable string
	label    string
	uri      string
}

// signatureQualificationData holds the (readable, label, uri) tuple for
// each constant.
var signatureQualificationData = map[SignatureQualification]signatureQualificationFields{
	SignatureQualificationQESig:  {"QESig", "Qualified Electronic Signature", "urn:cef:dss:signatureQualification:QESig"},
	SignatureQualificationQESeal: {"QESeal", "Qualified Electronic Seal", "urn:cef:dss:signatureQualification:QESeal"},
	SignatureQualificationUnknownQCQSCD: {
		"Unknown-QC-QSCD",
		"Signature produced by a Qualified Certificate of unknown type with its private key residing in a QSCD",
		"urn:cef:dss:signatureQualification:UnknownQCwithQSCD",
	},
	SignatureQualificationAdESigQC:            {"AdESig-QC", "Advanced Electronic Signature supported by a Qualified Certificate", "urn:cef:dss:signatureQualification:AdESigQC"},
	SignatureQualificationAdESealQC:           {"AdESeal-QC", "Advanced Electronic Seal supported by a Qualified Certificate", "urn:cef:dss:signatureQualification:AdESealQC"},
	SignatureQualificationUnknownQC:           {"Unknown-QC", "Signature produced by a Qualified Certificate of unknown type", "urn:cef:dss:signatureQualification:UnknownQC"},
	SignatureQualificationAdESig:              {"AdESig", "Advanced Electronic Signature", "urn:cef:dss:signatureQualification:AdESig"},
	SignatureQualificationAdESeal:             {"AdESeal", "Advanced Electronic Seal", "urn:cef:dss:signatureQualification:AdESeal"},
	SignatureQualificationUnknown:             {"Unknown", "Signature produced by a Certificate of unknown type", "urn:cef:dss:signatureQualification:Unknown"},
	SignatureQualificationIndeterminateQESig:  {"Indeterminate QESig", "Indeterminate Qualified Electronic Signature", "urn:cef:dss:signatureQualification:indeterminateQESig"},
	SignatureQualificationIndeterminateQESeal: {"Indeterminate QESeal", "Indeterminate Qualified Electronic Seal", "urn:cef:dss:signatureQualification:indeterminateQESeal"},
	SignatureQualificationIndeterminateUnknownQCQSCD: {
		"Indeterminate Unknown-QC-QSCD",
		"Indeterminate Signature produced by a Qualified Certificate of unknown type with its private key residing in a QSCD",
		"urn:cef:dss:signatureQualification:indeterminateUnknownQCwithQSCD",
	},
	SignatureQualificationIndeterminateAdESigQC: {
		"Indeterminate AdESig-QC",
		"Indeterminate Advanced Electronic Signature supported by a Qualified Certificate",
		"urn:cef:dss:signatureQualification:indeterminateAdESigQC",
	},
	SignatureQualificationIndeterminateAdESealQC: {
		"Indeterminate AdESeal-QC",
		"Indeterminate Advanced Electronic Seal supported by a Qualified Certificate",
		"urn:cef:dss:signatureQualification:indeterminateAdESealQC",
	},
	SignatureQualificationIndeterminateUnknownQC: {
		"Indeterminate Unknown-QC",
		"Indeterminate Signature produced by a Qualified Certificate of unknown type",
		"urn:cef:dss:signatureQualification:indeterminateUnknownQC",
	},
	SignatureQualificationIndeterminateAdESig: {
		"Indeterminate AdESig", "Indeterminate Advanced Electronic Signature", "urn:cef:dss:signatureQualification:indeterminateAdESig",
	},
	SignatureQualificationIndeterminateAdESeal: {
		"Indeterminate AdESeal", "Indeterminate Advanced Electronic Seal", "urn:cef:dss:signatureQualification:indeterminateAdESeal",
	},
	SignatureQualificationIndeterminateUnknown: {
		"Indeterminate Unknown", "Indeterminate Signature produced by a Certificate of unknown type", "urn:cef:dss:signatureQualification:indeterminateUnknown",
	},
	SignatureQualificationNotAdESQCQSCD: {
		"Not AdES but QC with QSCD", "Not Advanced Electronic Signature but supported by a Qualified Certificate", "urn:cef:dss:signatureQualification:notAdESbutQCwithQSCD",
	},
	SignatureQualificationNotAdESQC: {
		"Not AdES but QC", "Not Advanced Electronic Signature but supported by a Qualified Certificate", "urn:cef:dss:signatureQualification:notAdESbutQC",
	},
	SignatureQualificationNotAdES: {"Not AdES", "Not Advanced Electronic Signature", "urn:cef:dss:signatureQualification:notAdES"},
	SignatureQualificationNA:      {"N/A", "Not applicable", "urn:cef:dss:signatureQualification:notApplicable"},
}

// SignatureQualificationValues returns all constants in declaration order.
func SignatureQualificationValues() []SignatureQualification {
	return []SignatureQualification{
		SignatureQualificationQESig,
		SignatureQualificationQESeal,
		SignatureQualificationUnknownQCQSCD,
		SignatureQualificationAdESigQC,
		SignatureQualificationAdESealQC,
		SignatureQualificationUnknownQC,
		SignatureQualificationAdESig,
		SignatureQualificationAdESeal,
		SignatureQualificationUnknown,
		SignatureQualificationIndeterminateQESig,
		SignatureQualificationIndeterminateQESeal,
		SignatureQualificationIndeterminateUnknownQCQSCD,
		SignatureQualificationIndeterminateAdESigQC,
		SignatureQualificationIndeterminateAdESealQC,
		SignatureQualificationIndeterminateUnknownQC,
		SignatureQualificationIndeterminateAdESig,
		SignatureQualificationIndeterminateAdESeal,
		SignatureQualificationIndeterminateUnknown,
		SignatureQualificationNotAdESQCQSCD,
		SignatureQualificationNotAdESQC,
		SignatureQualificationNotAdES,
		SignatureQualificationNA,
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
