// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/job/OtherDocumentPointer.java (DSS 6.5.RC1).
package job

import "github.com/ryftcore/dss-go/dss/model"

// OtherDocumentPointer contains a pointer to another document to be processed.
type OtherDocumentPointer interface {
	// Location gets the location url. Port of getLocation().
	Location() string
	// SdiCertificates gets a list of ServiceDigitalIdentity X509 certificates. Port of
	// getSdiCertificates().
	SdiCertificates() []*model.CertificateToken
}
