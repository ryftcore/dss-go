// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/tsl/TLValidationJobSummary.java (DSS 6.5.RC1).
package tsl

import (
	"github.com/ryftcore/dss-go/dss/model"
)

// TLValidationJobSummary computes a summary for TLValidationJob.
//
// Implements the assumed job.ValidationJobSummary[TLInfo, LOTLInfo] interface
// (getDocumentListInfos/getOtherDocumentInfos, ported as DocumentListInfos()/
// OtherDocumentInfos()); java.io.Serializable has no Go counterpart and is dropped.
type TLValidationJobSummary struct {
	// lotlInfos is a list of LOTLs with a relationship between their TLs and pivots.
	lotlInfos []*LOTLInfo
	// otherTLInfos is the list of TL infos for otherTLSources.
	otherTLInfos []*TLInfo
}

// NewTLValidationJobSummary is the default constructor.
//
// Java's IllegalArgumentException("LOTL or TL Info shall be provided!") when both lotlInfos
// and otherTLInfos are empty is data-dependent, so it becomes a returned error rather than a
// panic.
func NewTLValidationJobSummary(lotlInfos []*LOTLInfo, otherTLInfos []*TLInfo) (*TLValidationJobSummary, error) {
	if len(lotlInfos) == 0 && len(otherTLInfos) == 0 {
		return nil, model.NewDSSError("LOTL or TL Info shall be provided!")
	}
	return &TLValidationJobSummary{lotlInfos: lotlInfos, otherTLInfos: otherTLInfos}, nil
}

// LOTLInfos returns a list of LOTLInfos for all processed LOTLs.
func (s *TLValidationJobSummary) LOTLInfos() []*LOTLInfo {
	return s.lotlInfos
}

// OtherTLInfos returns a list of TLInfos for other TLs.
func (s *TLValidationJobSummary) OtherTLInfos() []*TLInfo {
	return s.otherTLInfos
}

// NumberOfProcessedTLs returns an amount of processed TLs during the TL Validation job.
func (s *TLValidationJobSummary) NumberOfProcessedTLs() int {
	amount := 0
	amount += len(s.otherTLInfos)
	for _, lotlInfo := range s.lotlInfos {
		amount += len(lotlInfo.TLInfos())
	}
	return amount
}

// NumberOfProcessedLOTLs returns an amount of processed LOTLs during the TL Validation job.
func (s *TLValidationJobSummary) NumberOfProcessedLOTLs() int {
	return len(s.lotlInfos)
}

// TLInfoByID returns a TLInfo object by Identifier, or nil.
func (s *TLValidationJobSummary) TLInfoByID(identifier model.Identifier) *TLInfo {
	for _, tlInfo := range s.otherTLInfos {
		if identifier.Equals(tlInfo.DSSID()) {
			return tlInfo
		}
	}
	for _, lotlInfo := range s.lotlInfos {
		for _, tlInfo := range lotlInfo.TLInfos() {
			if identifier.Equals(tlInfo.DSSID()) {
				return tlInfo
			}
		}
	}
	return nil
}

// LOTLInfoByID returns a LOTLInfo object by Identifier, or nil.
func (s *TLValidationJobSummary) LOTLInfoByID(identifier model.Identifier) *LOTLInfo {
	for _, lotlInfo := range s.lotlInfos {
		if identifier.Equals(lotlInfo.DSSID()) {
			return lotlInfo
		}
	}
	return nil
}

// DocumentListInfos returns the list of LOTLInfos.
func (s *TLValidationJobSummary) DocumentListInfos() []*LOTLInfo {
	return s.LOTLInfos()
}

// OtherDocumentInfos returns the list of other TLInfos.
func (s *TLValidationJobSummary) OtherDocumentInfos() []*TLInfo {
	return s.OtherTLInfos()
}
