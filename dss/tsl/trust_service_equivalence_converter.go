// Ported from dss-tsl-validation/src/main/java/eu/europa/esig/dss/tsl/function/converter/TrustServiceEquivalenceConverter.java (DSS 6.5.RC1).
package tsl

import (
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model/timedependent"
	tslmodel "github.com/ryftcore/dss-go/dss/model/tsl"
	"github.com/ryftcore/dss-go/dss/trustedlist/jaxb"
	"github.com/ryftcore/dss-go/dss/utils"
)

// TrustServiceEquivalenceConverter extracts the MRA equivalence scheme for a Trusted List.
type TrustServiceEquivalenceConverter struct {
	// criteriaConverter is the used CriteriaListConverter.
	criteriaConverter *CriteriaListConverter
}

// NewTrustServiceEquivalenceConverter is the default constructor, instantiating a
// CriteriaListConverter. Port of TrustServiceEquivalenceConverter().
func NewTrustServiceEquivalenceConverter() *TrustServiceEquivalenceConverter {
	return &TrustServiceEquivalenceConverter{criteriaConverter: NewCriteriaListConverter()}
}

// Apply ports apply(TrustServiceEquivalenceInformationType).
func (c *TrustServiceEquivalenceConverter) Apply(t *jaxb.TrustServiceEquivalenceInformationType) *timedependent.MutableTimeDependentValues[*tslmodel.ServiceEquivalence] {
	result := timedependent.NewMutableTimeDependentValues[*tslmodel.ServiceEquivalence]()

	var status enumerations.MRAStatus
	if t.TrustServiceEquivalenceStatus != nil {
		status = enumerations.MRAStatus(*t.TrustServiceEquivalenceStatus)
	}
	serviceEquivalence := tslmodel.NewServiceEquivalenceBuilder().
		SetLegalInfoIdentifier(t.TrustServiceLegalIdentifier).
		SetStartDate(abstractParsingTaskConvertToDate(t.TrustServiceEquivalenceStatusStartingTime)).
		SetStatus(status).
		SetTypeAsiEquivalence(c.getTypeASiEquivalence(t.TrustServiceTSLTypeEquivalenceList)).
		SetStatusEquivalence(c.getStatusEquivalence(t.TrustServiceTSLStatusEquivalenceList)).
		SetCertificateContentEquivalences(c.getCertificateEquivalence(t.CertificateContentReferencesEquivalenceList)).
		SetQualifierEquivalence(c.getQualifierEquivalence(t.TrustServiceTSLQualificationExtensionEquivalenceList)).
		Build()
	result.AddOldest(serviceEquivalence)

	oldestStartDate := serviceEquivalence.StartDate()

	// TODO : review lists nesting
	for _, trustServiceEquivalenceHistory := range t.TrustServiceEquivalenceHistory {
		for _, h := range trustServiceEquivalenceHistory.TrustServiceEquivalenceHistoryInstance {
			var historyStatus enumerations.MRAStatus
			if h.TrustServiceEquivalenceStatus != nil {
				historyStatus = enumerations.MRAStatus(*h.TrustServiceEquivalenceStatus)
			}
			historyServiceEquivalence := tslmodel.NewServiceEquivalenceBuilder().
				SetLegalInfoIdentifier(serviceEquivalence.LegalInfoIdentifier()).
				SetStartDate(abstractParsingTaskConvertToDate(h.TrustServiceEquivalenceStatusStartingTime)).
				SetEndDate(oldestStartDate).
				SetStatus(historyStatus).
				SetTypeAsiEquivalence(c.getTypeASiEquivalence(h.TrustServiceTSLTypeEquivalenceList)).
				SetStatusEquivalence(c.getStatusEquivalence(h.TrustServiceTSLStatusEquivalenceList)).
				SetCertificateContentEquivalences(c.getCertificateEquivalence(h.CertificateContentReferencesEquivalenceList)).
				SetQualifierEquivalence(c.getQualifierEquivalence(h.TrustServiceTSLQualificationExtensionEquivalenceList)).
				Build()
			result.AddOldest(historyServiceEquivalence)

			oldestStartDate = historyServiceEquivalence.StartDate()
		}
	}

	return result
}

func (c *TrustServiceEquivalenceConverter) getTypeASiEquivalence(
	serviceTSLTypeEquivalenceList *jaxb.TrustServiceTSLTypeEquivalenceListType) map[tslmodel.ServiceTypeASi]tslmodel.ServiceTypeASi {
	typeAsiEquivalence := make(map[tslmodel.ServiceTypeASi]tslmodel.ServiceTypeASi)
	if serviceTSLTypeEquivalenceList != nil {
		expected := serviceTSLTypeEquivalenceList.TrustServiceTSLTypeListPointedParty
		substitute := serviceTSLTypeEquivalenceList.TrustServiceTSLTypeListPointingParty
		if expected == nil || substitute == nil {
			return typeAsiEquivalence
		}
		for _, expectedTypeASI := range expected.TrustServiceTSLType {
			staExpected := c.getServiceTypeASi(expectedTypeASI)
			for _, substituteTypeASI := range substitute.TrustServiceTSLType {
				staSubstitute := c.getServiceTypeASi(substituteTypeASI)
				typeAsiEquivalence[staExpected] = staSubstitute
			}
		}
	}
	return typeAsiEquivalence
}

func (c *TrustServiceEquivalenceConverter) getServiceTypeASi(expectedTypeASI *jaxb.TrustServiceTSLTypeType) tslmodel.ServiceTypeASi {
	sta := tslmodel.NewServiceTypeASi()
	sta.SetType(expectedTypeASI.ServiceTypeIdentifier)
	additionalServiceInformation := expectedTypeASI.AdditionalServiceInformation
	if additionalServiceInformation != nil && additionalServiceInformation.URI != nil {
		sta.SetAsi(additionalServiceInformation.URI.Value)
	}
	return *sta
}

func (c *TrustServiceEquivalenceConverter) getStatusEquivalence(
	serviceTSLStatusEquivalenceList *jaxb.TrustServiceTSLStatusEquivalenceListType) []tslmodel.StatusEquivalenceMapping {
	var statusEquivalence []tslmodel.StatusEquivalenceMapping
	if serviceTSLStatusEquivalenceList == nil {
		return statusEquivalence
	}
	statusEquivalence = c.extractEquivalences(serviceTSLStatusEquivalenceList.TrustServiceTSLStatusValidEquivalence, statusEquivalence)
	statusEquivalence = c.extractEquivalences(serviceTSLStatusEquivalenceList.TrustServiceTSLStatusInvalidEquivalence, statusEquivalence)
	return statusEquivalence
}

func (c *TrustServiceEquivalenceConverter) extractEquivalences(statusEquivalence *jaxb.TrustServiceTSLStatusEquivalenceType,
	statusEquivalenceMapping []tslmodel.StatusEquivalenceMapping) []tslmodel.StatusEquivalenceMapping {
	if statusEquivalence == nil {
		return statusEquivalenceMapping
	}
	serviceTSLStatusListExpected := statusEquivalence.TrustServiceTSLStatusListPointedParty
	serviceTSLStatusListSubstitute := statusEquivalence.TrustServiceTSLStatusListPointingParty
	if serviceTSLStatusListExpected == nil || serviceTSLStatusListSubstitute == nil {
		return statusEquivalenceMapping
	}
	return append(statusEquivalenceMapping, tslmodel.StatusEquivalenceMapping{
		PointedStatuses:  serviceTSLStatusListExpected.ServiceStatus,
		PointingStatuses: serviceTSLStatusListSubstitute.ServiceStatus,
	})
}

func (c *TrustServiceEquivalenceConverter) getCertificateEquivalence(
	certificateContentEquivalenceList *jaxb.CertificateContentReferencesEquivalenceListType) []*tslmodel.CertificateContentEquivalence {
	var certificateContentEquivalences []*tslmodel.CertificateContentEquivalence
	if certificateContentEquivalenceList != nil && utils.IsCollectionNotEmpty(certificateContentEquivalenceList.CertificateContentReferenceEquivalence) {
		for _, certEquiv := range certificateContentEquivalenceList.CertificateContentReferenceEquivalence {
			expected := certEquiv.CertificateContentDeclarationPointedParty
			substitute := certEquiv.CertificateContentDeclarationPointingParty
			condition := c.criteriaConverter.Apply(substitute)

			equiv := tslmodel.NewCertificateContentEquivalence()
			equiv.SetContext(enumerations.MRAEquivalenceContext(certEquiv.CertificateContentReferenceEquivalenceContext))
			equiv.SetCondition(c.criteriaConverter.Apply(expected))
			equiv.SetContentReplacement(c.getQCStatementOids(condition))

			certificateContentEquivalences = append(certificateContentEquivalences, equiv)
		}
	}
	return certificateContentEquivalences
}

// getQCStatementOids ports the private getQCStatementOids(Condition).
func (c *TrustServiceEquivalenceConverter) getQCStatementOids(condition tslmodel.Condition) *tslmodel.QCStatementOids {
	result := tslmodel.NewQCStatementOids()

	var qcStatementIds, qcTypeIds, qcCClegislations []string
	var qcStatementIdsToRemove, qcTypeIdsToRemove, qcCClegislationsToRemove []string

	if composite, ok := condition.(*CompositeCondition); ok {
		switch composite.MatchingCriteriaIndicator() {
		case enumerations.AssertAll:
			for _, childCondition := range composite.Children() {
				qcStatementIds, qcTypeIds, qcCClegislations, qcStatementIdsToRemove, qcTypeIdsToRemove, qcCClegislationsToRemove =
					c.populateFromChild(childCondition, qcStatementIds, qcTypeIds, qcCClegislations,
						qcStatementIdsToRemove, qcTypeIdsToRemove, qcCClegislationsToRemove)
			}
		case enumerations.AssertAtLeastOne:
			children := composite.Children()
			if len(children) > 0 {
				qcStatementIds, qcTypeIds, qcCClegislations, qcStatementIdsToRemove, qcTypeIdsToRemove, qcCClegislationsToRemove =
					c.populateFromChild(children[0], qcStatementIds, qcTypeIds, qcCClegislations,
						qcStatementIdsToRemove, qcTypeIdsToRemove, qcCClegislationsToRemove)
			}
		case enumerations.AssertNone:
			for _, childCondition := range composite.Children() {
				// Reversed lists for NONE.
				qcStatementIdsToRemove, qcTypeIdsToRemove, qcCClegislationsToRemove, qcStatementIds, qcTypeIds, qcCClegislations =
					c.populateFromChild(childCondition, qcStatementIdsToRemove, qcTypeIdsToRemove, qcCClegislationsToRemove,
						qcStatementIds, qcTypeIds, qcCClegislations)
			}
		}
	}

	if qcCondition, ok := condition.(*QCStatementCondition); ok {
		if oid := qcCondition.Oid(); utils.IsStringNotEmpty(oid) {
			qcStatementIds = append(qcStatementIds, oid)
		}
		if typ := qcCondition.Type(); utils.IsStringNotEmpty(typ) {
			qcTypeIds = append(qcTypeIds, typ)
		}
		if legislation := qcCondition.Legislation(); utils.IsStringNotEmpty(legislation) {
			qcCClegislations = append(qcCClegislations, legislation)
		}
	}

	result.SetQcStatementIds(qcStatementIds)
	result.SetQcTypeIds(qcTypeIds)
	result.SetQcCClegislations(qcCClegislations)
	result.SetQcStatementIdsToRemove(qcStatementIdsToRemove)
	result.SetQcTypeIdsToRemove(qcTypeIdsToRemove)
	result.SetQcCClegislationsToRemove(qcCClegislationsToRemove)
	return result
}

// populateFromChild ports the private populateFromChild(...), merging a child condition's
// extracted OIDs into the accumulator slices without duplicates.
func (c *TrustServiceEquivalenceConverter) populateFromChild(condition tslmodel.Condition,
	qcStatementIds, qcTypeIds, qcCClegislations, qcStatementIdsToRemove, qcTypeIdsToRemove, qcCClegislationsToRemove []string) (
	[]string, []string, []string, []string, []string, []string) {
	conditionResult := c.getQCStatementOids(condition)
	qcStatementIds = appendUnique(qcStatementIds, conditionResult.QcStatementIds()...)
	qcTypeIds = appendUnique(qcTypeIds, conditionResult.QcTypeIds()...)
	qcCClegislations = appendUnique(qcCClegislations, conditionResult.QcCClegislations()...)
	qcStatementIdsToRemove = appendUnique(qcStatementIdsToRemove, conditionResult.QcStatementIdsToRemove()...)
	qcTypeIdsToRemove = appendUnique(qcTypeIdsToRemove, conditionResult.QcTypeIdsToRemove()...)
	qcCClegislationsToRemove = appendUnique(qcCClegislationsToRemove, conditionResult.QcCClegislationsToRemove()...)
	return qcStatementIds, qcTypeIds, qcCClegislations, qcStatementIdsToRemove, qcTypeIdsToRemove, qcCClegislationsToRemove
}

// appendUnique appends each of values to list unless already present, mirroring the repeated
// `if (!list.contains(x)) list.add(x)` pattern in populateFromChild.
func appendUnique(list []string, values ...string) []string {
	for _, v := range values {
		found := false
		for _, existing := range list {
			if existing == v {
				found = true
				break
			}
		}
		if !found {
			list = append(list, v)
		}
	}
	return list
}

func (c *TrustServiceEquivalenceConverter) getQualifierEquivalence(
	qualificationExtensionEquivalenceListType *jaxb.TrustServiceTSLQualificationExtensionEquivalenceListType) map[string]string {
	qualifierEquivalenceMap := make(map[string]string)
	if qualificationExtensionEquivalenceListType != nil && utils.IsCollectionNotEmpty(qualificationExtensionEquivalenceListType.QualifierEquivalenceList) {
		for _, qualifierEquivalenceList := range qualificationExtensionEquivalenceListType.QualifierEquivalenceList {
			for _, qualifierEquivalenceType := range qualifierEquivalenceList.QualifierEquivalence {
				qualifierExpected := qualifierEquivalenceType.QualifierPointedParty
				qualifierSubstitute := qualifierEquivalenceType.QualifierPointingParty
				if qualifierExpected == nil || qualifierSubstitute == nil {
					continue
				}
				var expectedURI, substituteURI string
				if qualifierExpected.URI != nil {
					expectedURI = *qualifierExpected.URI
				}
				if qualifierSubstitute.URI != nil {
					substituteURI = *qualifierSubstitute.URI
				}
				qualifierEquivalenceMap[expectedURI] = substituteURI
			}
		}
	}
	return qualifierEquivalenceMap
}
