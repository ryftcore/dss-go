// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/evidencerecord/EvidenceRecordValidator.java (DSS 6.5.RC1).
package evidencerecord

import (
	"github.com/ryftcore/dss-go/dss/enumerations"
	spivalidation "github.com/ryftcore/dss-go/dss/spi/validation"
	dssvalidation "github.com/ryftcore/dss-go/dss/validation"
)

// EvidenceRecordValidator is the interface to be used for evidence record validation.
type EvidenceRecordValidator interface {
	dssvalidation.DocumentValidator

	// EvidenceRecord returns a single EvidenceRecord to be validated. Port of getEvidenceRecord().
	EvidenceRecord() spivalidation.EvidenceRecord

	// EvidenceRecordType returns a type of the evidence record supported by the current
	// validator. Port of getEvidenceRecordType().
	EvidenceRecordType() enumerations.EvidenceRecordTypeEnum
}
