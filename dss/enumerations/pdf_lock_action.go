// Ported from dss-enumerations/.../PdfLockAction.java (DSS 6.5.RC1).
package enumerations

import "fmt"

// PdfLockAction is a name which, in conjunction with Fields, indicates the
// set of fields that should be locked.
type PdfLockAction string

const (
	// PdfLockActionAll means all form fields do not permit changes.
	PdfLockActionAll PdfLockAction = "ALL"
	// PdfLockActionInclude means only those form fields specified in
	// fields do not permit changes.
	PdfLockActionInclude PdfLockAction = "INCLUDE"
	// PdfLockActionExclude means only those form fields not specified in
	// fields do not permit changes.
	PdfLockActionExclude PdfLockAction = "EXCLUDE"
)

var pdfLockActionName = map[PdfLockAction]string{
	PdfLockActionAll:     "All",
	PdfLockActionInclude: "Include",
	PdfLockActionExclude: "Exclude",
}

// PdfLockActionValues returns all constants in declaration order.
func PdfLockActionValues() []PdfLockAction {
	return []PdfLockAction{
		PdfLockActionAll,
		PdfLockActionInclude,
		PdfLockActionExclude,
	}
}

// Name returns the name value of the field parameter.
func (p PdfLockAction) Name() string {
	return pdfLockActionName[p]
}

// PdfLockActionForName returns a PdfLockAction corresponding to the given
// name. Returns an error if no constant matches, mirroring Java's
// IllegalArgumentException.
func PdfLockActionForName(name string) (PdfLockAction, error) {
	for _, v := range PdfLockActionValues() {
		if pdfLockActionName[v] == name {
			return v, nil
		}
	}
	return "", fmt.Errorf("unsupported /Action field value : %s", name)
}

// PdfLockActionValueOf returns the constant matching the given Java enum
// name.
func PdfLockActionValueOf(name string) (PdfLockAction, error) {
	for _, v := range PdfLockActionValues() {
		if string(v) == name {
			return v, nil
		}
	}
	return "", &pdfLockActionInvalidValueError{name}
}

type pdfLockActionInvalidValueError struct {
	name string
}

func (e *pdfLockActionInvalidValueError) Error() string {
	return "no enum constant PdfLockAction." + e.name
}
