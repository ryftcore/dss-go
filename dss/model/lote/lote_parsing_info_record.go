// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/lote/record/LoTEParsingInfoRecord.java (DSS 6.5.RC1).
//
// Java package eu.europa.esig.dss.model.lote.record is flattened into this lote package.
package lote

import (
	"time"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model/job"
)

// LoTEParsingInfoRecord represents a List of Trusted Entities parsing result record.
type LoTEParsingInfoRecord interface {
	job.ParsingInfoRecord

	// Type gets the List Type.
	Type() enumerations.ListType
	// SequenceNumber gets the List sequence number. Nil ports Java's null Integer.
	SequenceNumber() *int
	// Version gets List version. Nil ports Java's null Integer.
	Version() *int
	// Territory gets the List territory (country).
	Territory() string
	// IssueDate gets issuing date.
	IssueDate() time.Time
	// NextUpdateDate gets next update date.
	NextUpdateDate() time.Time
	// DistributionPoints gets distribution points.
	DistributionPoints() []string
	// TrustedEntities gets trusted entities.
	TrustedEntities() []*TrustedEntity
	// CurrentListPointers gets List of Lists other TSL pointers.
	CurrentListPointers() []*OtherListPointer
	// OtherListPointers gets Lists other TSL pointers.
	OtherListPointers() []*OtherListPointer
	// PivotUrls gets pivot URLs.
	PivotUrls() []string
	// SigningCertificateAnnouncementUrl gets signing certificate announcement URL.
	SigningCertificateAnnouncementUrl() string
	// TrustedEntitiesNumber returns a number of all trusted entities present in the List.
	TrustedEntitiesNumber() int
	// TrustedServicesNumber returns a number of all trusted services present in the List.
	TrustedServicesNumber() int
	// CertNumber returns a number of all CertificateTokens present in the List.
	CertNumber() int
	// StructureValidationMessages gets a list of error messages when occurred during the
	// structure validation.
	StructureValidationMessages() []string
}
