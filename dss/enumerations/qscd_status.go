// Ported from dss-enumerations/.../QSCDStatus.java (DSS 6.5.RC1).
//
// Defines if the certificate is QSCD.
package enumerations

import "fmt"

// QSCDStatus defines if the certificate is QSCD.
type QSCDStatus string

const (
	// QSCDStatus_QSCD: certificate to be used with Qualified Electronic
	// Signature Creation Device.
	QSCDStatus_QSCD QSCDStatus = "QSCD"
	// QSCDStatus_NOT_QSCD: not a certificate to be used with Qualified
	// Electronic Signature Creation Device.
	QSCDStatus_NOT_QSCD QSCDStatus = "NOT_QSCD"
)

// QSCDStatusValues returns all constants in declaration order.
func QSCDStatusValues() []QSCDStatus {
	return []QSCDStatus{
		QSCDStatus_QSCD,
		QSCDStatus_NOT_QSCD,
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
	return QSCDStatus_QSCD == status
}
