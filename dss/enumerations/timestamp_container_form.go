// Ported from dss-enumerations/.../TimestampContainerForm.java (DSS 6.5.RC1).
package enumerations

import "strings"

// TimestampContainerForm is used for a detached timestamping service.
type TimestampContainerForm string

const (
	// TimestampContainerForm_PDF is used to timestamp a PDF document.
	TimestampContainerForm_PDF TimestampContainerForm = "PDF"
	// TimestampContainerForm_ASiC_E is used to timestamp provided
	// document(s) and creates an ASiC-E container.
	TimestampContainerForm_ASiC_E TimestampContainerForm = "ASiC_E"
	// TimestampContainerForm_ASiC_S is used to timestamp provided
	// document(s) and creates an ASiC-S container.
	TimestampContainerForm_ASiC_S TimestampContainerForm = "ASiC_S"
)

// TimestampContainerFormValues returns all constants in declaration order.
func TimestampContainerFormValues() []TimestampContainerForm {
	return []TimestampContainerForm{
		TimestampContainerForm_PDF,
		TimestampContainerForm_ASiC_E,
		TimestampContainerForm_ASiC_S,
	}
}

// Readable returns the user-friendly name of the form.
func (t TimestampContainerForm) Readable() string {
	return strings.ReplaceAll(string(t), "_", "-")
}

// TimestampContainerFormValueOf returns the constant matching the given
// Java enum name.
func TimestampContainerFormValueOf(name string) (TimestampContainerForm, error) {
	for _, v := range TimestampContainerFormValues() {
		if string(v) == name {
			return v, nil
		}
	}
	return "", &timestampContainerFormInvalidValueError{name}
}

type timestampContainerFormInvalidValueError struct {
	name string
}

func (e *timestampContainerFormInvalidValueError) Error() string {
	return "no enum constant TimestampContainerForm." + e.name
}
