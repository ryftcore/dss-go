// Ported from dss-enumerations/.../TimestampContainerForm.java (DSS 6.5.RC1).
package enumerations

import "strings"

// TimestampContainerForm is used for a detached timestamping service.
type TimestampContainerForm string

const (
	// TimestampContainerFormPDF is used to timestamp a PDF document.
	TimestampContainerFormPDF TimestampContainerForm = "PDF"
	// TimestampContainerFormASiCE is used to timestamp provided
	// document(s) and creates an ASiC-E container.
	TimestampContainerFormASiCE TimestampContainerForm = "ASiC_E"
	// TimestampContainerFormASiCS is used to timestamp provided
	// document(s) and creates an ASiC-S container.
	TimestampContainerFormASiCS TimestampContainerForm = "ASiC_S"
)

// TimestampContainerFormValues returns all constants in declaration order.
func TimestampContainerFormValues() []TimestampContainerForm {
	return []TimestampContainerForm{
		TimestampContainerFormPDF,
		TimestampContainerFormASiCE,
		TimestampContainerFormASiCS,
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
