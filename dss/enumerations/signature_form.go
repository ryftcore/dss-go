// Ported from dss-enumerations/.../SignatureForm.java (DSS 6.5.RC1).

package enumerations

import "fmt"

// SignatureForm lists the different signature forms.
type SignatureForm string

const (
	// SignatureForm_XAdES is an XML-based signature according to EN 319 132.
	SignatureForm_XAdES SignatureForm = "XAdES"
	// SignatureForm_CAdES is a CMS-based signature according to EN 319 122.
	SignatureForm_CAdES SignatureForm = "CAdES"
	// SignatureForm_JAdES is a JSON-based signature according to TS 119 182.
	SignatureForm_JAdES SignatureForm = "JAdES"
	// SignatureForm_CBAdES is a CBOR-based signature according to TS 119 152.
	SignatureForm_CBAdES SignatureForm = "CBAdES"
	// SignatureForm_PAdES is a PDF-based signature according to EN 319 142.
	SignatureForm_PAdES SignatureForm = "PAdES"
	// SignatureForm_PKCS7 is a PDF-based signature according to ISO 32000.
	SignatureForm_PKCS7 SignatureForm = "PKCS7"
)

// SignatureFormValues returns all SignatureForm constants in declaration order.
func SignatureFormValues() []SignatureForm {
	return []SignatureForm{
		SignatureForm_XAdES,
		SignatureForm_CAdES,
		SignatureForm_JAdES,
		SignatureForm_CBAdES,
		SignatureForm_PAdES,
		SignatureForm_PKCS7,
	}
}

// SignatureFormValueOf returns the SignatureForm matching the given Java enum name.
func SignatureFormValueOf(name string) (SignatureForm, error) {
	for _, v := range SignatureFormValues() {
		if string(v) == name {
			return v, nil
		}
	}
	return "", fmt.Errorf("no enum constant SignatureForm.%s", name)
}
