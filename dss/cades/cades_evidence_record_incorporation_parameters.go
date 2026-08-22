// Ported from dss-cades/src/main/java/eu/europa/esig/dss/cades/evidencerecord/CAdESEvidenceRecordIncorporationParameters.java (DSS 6.5.RC1).
package cades

import "github.com/ryftcore/dss-go/dss/document"

// EvidenceRecordIncorporationParameters holds parameters for an existing evidence record
// embedding into an existing CAdES signature.
type EvidenceRecordIncorporationParameters struct {
	document.AbstractEvidenceRecordIncorporationParameters
}

// NewCAdESEvidenceRecordIncorporationParameters is the default constructor.
func NewCAdESEvidenceRecordIncorporationParameters() *EvidenceRecordIncorporationParameters {
	return &EvidenceRecordIncorporationParameters{
		AbstractEvidenceRecordIncorporationParameters: document.NewAbstractEvidenceRecordIncorporationParameters(),
	}
}
