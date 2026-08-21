// Ported from dss-cades/src/main/java/eu/europa/esig/dss/cades/evidencerecord/CAdESEvidenceRecordIncorporationParameters.java (DSS 6.5.RC1).
package cades

import "github.com/ryftcore/dss-go/dss/document"

// CAdESEvidenceRecordIncorporationParameters holds parameters for an existing evidence record
// embedding into an existing CAdES signature.
type CAdESEvidenceRecordIncorporationParameters struct {
	document.AbstractEvidenceRecordIncorporationParameters
}

// NewCAdESEvidenceRecordIncorporationParameters is the default constructor.
func NewCAdESEvidenceRecordIncorporationParameters() *CAdESEvidenceRecordIncorporationParameters {
	return &CAdESEvidenceRecordIncorporationParameters{
		AbstractEvidenceRecordIncorporationParameters: document.NewAbstractEvidenceRecordIncorporationParameters(),
	}
}
