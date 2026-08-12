// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/tsl/TLParsingInfoRecord.java (DSS 6.5.RC1).
package tsl

import (
	"time"

	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/model/job"
)

// TLParsingInfoRecord represents a Trusted List parsing result record.
type TLParsingInfoRecord interface {
	job.ParsingInfoRecord

	// TSLType gets the LOTL/TL TSLType.
	TSLType() enumerations.TSLType
	// SequenceNumber gets the LOTL/TL sequence number. Nil ports Java's null Integer.
	SequenceNumber() *int
	// Version gets the LOTL/TL version. Nil ports Java's null Integer.
	Version() *int
	// Territory gets the LOTL/TL territory (country).
	Territory() string
	// IssueDate gets the issuing date.
	IssueDate() time.Time
	// NextUpdateDate gets the next update date.
	NextUpdateDate() time.Time
	// DistributionPoints gets distribution points.
	DistributionPoints() []string
	// TrustServiceProviders gets trust service providers.
	TrustServiceProviders() []*TrustServiceProvider
	// LotlOtherPointers gets LOTL other TSL pointers.
	LotlOtherPointers() []*OtherTSLPointer
	// TlOtherPointers gets TL other TSL pointers.
	TlOtherPointers() []*OtherTSLPointer
	// PivotUrls gets pivot URLs.
	PivotUrls() []string
	// SigningCertificateAnnouncementUrl gets the signing certificate announcement URL.
	SigningCertificateAnnouncementUrl() string
	// TSPNumber returns a number of all TrustServiceProviders present in the TL.
	TSPNumber() int
	// TSNumber returns a number of all TrustServices present in the TL.
	TSNumber() int
	// CertNumber returns a number of all CertificateTokens present in the TL.
	CertNumber() int
}
