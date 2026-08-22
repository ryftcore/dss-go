// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/lote/LoTEValidationJobSummary.java (DSS 6.5.RC1).
package lote

import "github.com/ryftcore/dss-go/dss/model"

// ValidationJobSummary contains summary of the validation result of a List of Trusted
// Entities validation job.
//
// Implements the assumed job.ValidationJobSummary[Info, LoLoTEInfo] interface
// (getDocumentListInfos/getOtherDocumentInfos, ported as DocumentListInfos()/
// OtherDocumentInfos()); java.io.Serializable has no Go counterpart and is dropped.
type ValidationJobSummary struct {
	// loloteInfos is a list of infos for processed LoLoTESource's.
	loloteInfos []*LoLoTEInfo
	// otherLoTEInfos is a list of infos for processed other LoTESource's.
	otherLoTEInfos []*Info
}

// NewLoTEValidationJobSummary is the default constructor.
//
// Java's IllegalArgumentException("LoTE Info shall be provided!") when both loloteInfos and
// otherLoTEInfos are empty is data-dependent, so it becomes a returned error rather than a
// panic (mirrors tsl.NewTLValidationJobSummary).
func NewLoTEValidationJobSummary(loloteInfos []*LoLoTEInfo, otherLoTEInfos []*Info) (*ValidationJobSummary, error) {
	if len(loloteInfos) == 0 && len(otherLoTEInfos) == 0 {
		return nil, model.NewDSSError("LoTE Info shall be provided!")
	}
	return &ValidationJobSummary{loloteInfos: loloteInfos, otherLoTEInfos: otherLoTEInfos}, nil
}

// LoLoTEInfos gets a list of LoLoTE infos.
func (s *ValidationJobSummary) LoLoTEInfos() []*LoLoTEInfo {
	return s.loloteInfos
}

// OtherLoTEInfos gets a list of other LoTE infos.
func (s *ValidationJobSummary) OtherLoTEInfos() []*Info {
	return s.otherLoTEInfos
}

// NumberOfProcessedLoTEs gets a number of processed LoTEs.
func (s *ValidationJobSummary) NumberOfProcessedLoTEs() int {
	amount := 0
	amount += len(s.otherLoTEInfos)
	for _, loloteInfo := range s.loloteInfos {
		amount += len(loloteInfo.ChildrenInfos())
	}
	return amount
}

// NumberOfProcessedLoLoTEs returns an amount of processed LoLoTEs.
func (s *ValidationJobSummary) NumberOfProcessedLoLoTEs() int {
	return len(s.loloteInfos)
}

// LoTEInfoByID gets a LoTE by a unique identifier, or nil.
func (s *ValidationJobSummary) LoTEInfoByID(identifier model.Identifier) *Info {
	for _, listInfo := range s.otherLoTEInfos {
		if identifier.Equals(listInfo.DSSID()) {
			return listInfo
		}
	}
	for _, loloteInfo := range s.loloteInfos {
		for _, loteInfo := range loloteInfo.ChildrenInfos() {
			if identifier.Equals(loteInfo.DSSID()) {
				return loteInfo
			}
		}
	}
	return nil
}

// LoLoTEInfoByID returns a LoLoTEInfo object by Identifier, or nil.
func (s *ValidationJobSummary) LoLoTEInfoByID(identifier model.Identifier) *LoLoTEInfo {
	for _, loloteInfo := range s.loloteInfos {
		if identifier.Equals(loloteInfo.DSSID()) {
			return loloteInfo
		}
	}
	return nil
}

// DocumentListInfos returns the list of LoLoTEInfos.
func (s *ValidationJobSummary) DocumentListInfos() []*LoLoTEInfo {
	return s.LoLoTEInfos()
}

// OtherDocumentInfos returns the list of other LoTEInfos.
func (s *ValidationJobSummary) OtherDocumentInfos() []*Info {
	return s.OtherLoTEInfos()
}
