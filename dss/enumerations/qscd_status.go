// Ported from dss-enumerations/.../QSCDStatus.java (DSS 6.5.RC1).
//
// Defines if the certificate is QSCD.
package enumerations

import "fmt"

// QSCDStatus defines if the certificate is QSCD.
type QSCDStatus string

const (
	// QSCDStatusQSCD: certificate to be used with Qualified Electronic
	// Signature Creation Device.
	QSCDStatusQSCD QSCDStatus = "QSCD"
	// QSCDStatusNotQSCD: not a certificate to be used with Qualified
	// Electronic Signature Creation Device.
	QSCDStatusNotQSCD QSCDStatus = "NOT_QSCD"
)

// QSCDStatusValues returns all constants in declaration order.
func QSCDStatusValues() []QSCDStatus {
	return []QSCDStatus{
		QSCDStatusQSCD,
		QSCDStatusNotQSCD,
	}
}

// QSCDStatusValueOf returns the QSCDStatus matching the given Java enum name.
func QSCDStatusValueOf(name string) (QSCDStatus, error) {
	for _, v := range QSCDStatusValues() {
		if string(v) == name {
			return v, nil
		}
	}
	return "", fmt.Errorf("no enum constant QSCDStatus.%s", name)
}

// QSCDStatusIsQSCD checks if the certificate is a QSCD.
func QSCDStatusIsQSCD(status QSCDStatus) bool {
	return QSCDStatusQSCD == status
}
