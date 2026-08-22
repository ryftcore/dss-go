// Ported from dss-enumerations/.../SignatureForm.java (DSS 6.5.RC1).

package enumerations

import "fmt"

// SignatureForm lists the different signature forms.
type SignatureForm string

const (
	// SignatureFormXAdES is an XML-based signature according to EN 319 132.
	SignatureFormXAdES SignatureForm = "XAdES"
	// SignatureFormCAdES is a CMS-based signature according to EN 319 122.
	SignatureFormCAdES SignatureForm = "CAdES"
	// SignatureFormJAdES is a JSON-based signature according to TS 119 182.
	SignatureFormJAdES SignatureForm = "JAdES"
	// SignatureFormCBAdES is a CBOR-based signature according to TS 119 152.
	SignatureFormCBAdES SignatureForm = "CBAdES"
	// SignatureFormPAdES is a PDF-based signature according to EN 319 142.
	SignatureFormPAdES SignatureForm = "PAdES"
	// SignatureFormPKCS7 is a PDF-based signature according to ISO 32000.
	SignatureFormPKCS7 SignatureForm = "PKCS7"
)

// SignatureFormValues returns all SignatureForm constants in declaration order.
func SignatureFormValues() []SignatureForm {
	return []SignatureForm{
		SignatureFormXAdES,
		SignatureFormCAdES,
		SignatureFormJAdES,
		SignatureFormCBAdES,
		SignatureFormPAdES,
		SignatureFormPKCS7,
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
