// Ported from dss-tsl-validation/src/main/java/eu/europa/esig/dss/tsl/function/converter/MRAConverter.java (DSS 6.5.RC1).
package tsl

import (
	"github.com/ryftcore/dss-go/dss/model/timedependent"
	tslmodel "github.com/ryftcore/dss-go/dss/model/tsl"
	"github.com/ryftcore/dss-go/dss/trustedlist/jaxb"
)

// MRAConverter converts a JAXB MutualRecognitionAgreementInformationType to a Go MRA.
type MRAConverter struct {
	// converter is the TrustServiceEquivalence converter.
	converter *TrustServiceEquivalenceConverter
}

// NewMRAConverter is the default constructor. Port of MRAConverter().
func NewMRAConverter() *MRAConverter {
	return &MRAConverter{converter: NewTrustServiceEquivalenceConverter()}
}

// Apply ports apply(MutualRecognitionAgreementInformationType).
func (c *MRAConverter) Apply(t *jaxb.MutualRecognitionAgreementInformationType) *tslmodel.MRA {
	result := tslmodel.NewMRA()
	if t.TechnicalType != nil {
		result.SetTechnicalType(t.TechnicalType.String())
	}
	if t.Version != nil {
		result.SetVersion(t.Version.String())
	}
	result.SetPointingContractingPartyLegislation(t.PointingContractingPartyLegislation)
	result.SetPointedContractingPartyLegislation(t.PointedContractingPartyLegislation)

	var serviceEquivalences []*timedependent.MutableTimeDependentValues[*tslmodel.ServiceEquivalence]
	for _, trustServiceEquivalenceInformationType := range t.TrustServiceEquivalenceInformation {
		serviceEquivalences = append(serviceEquivalences, c.converter.Apply(trustServiceEquivalenceInformationType))
	}
	result.SetServiceEquivalence(serviceEquivalences)

	return result
}
