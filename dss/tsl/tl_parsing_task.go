// Ported from dss-tsl-validation/src/main/java/eu/europa/esig/dss/tsl/parsing/TLParsingTask.java (DSS 6.5.RC1).
//
// CROSS-CHUNK DEPENDENCY (see this batch's porter notes): the three structural predicates
// (NonEmptyTSPInformation, NonEmptyServiceInformation, NonEmptyTrustService) and
// TrustServiceProviderConverter live in dss-tsl-validation's "function" package, which the TSLJOB
// chunk ports into this same Go package (tsl). They are referenced here by their Java names, with
// this codebase's constructor (New<Name>) and functional-interface (Test / Apply) spellings.
package tsl

import (
	"github.com/ryftcore/dss-go/dss/model"
	tslmodel "github.com/ryftcore/dss-go/dss/model/tsl"
	"github.com/ryftcore/dss-go/dss/trustedlist/jaxb"
	"github.com/ryftcore/dss-go/dss/utils"
)

// TLParsingTask parses a TL and returns a TLParsingResult.
type TLParsingTask struct {
	AbstractParsingTaskBase

	// tlSource is the TLSource to parse.
	tlSource *TLSource
}

var _ AbstractParsingTaskOverrides = (*TLParsingTask)(nil)

// NewTLParsingTask is the default constructor, taking the TL document to parse and its TLSource.
// Port of TLParsingTask(DSSDocument, TLSource).
//
// Panics with the Java messages when either argument is nil (Objects.requireNonNull; the document
// check lives in the superclass constructor).
func NewTLParsingTask(document model.DSSDocument, tlSource *TLSource) *TLParsingTask {
	task := &TLParsingTask{AbstractParsingTaskBase: NewAbstractParsingTaskBase(document)}
	if tlSource == nil {
		panic("The TLSource is null")
	}
	task.tlSource = tlSource
	task.InitAbstractParsingTask(task)
	return task
}

// Get parses the TL. Port of the get() override (Supplier#get); the exceptions Java's
// getJAXBObject and verifyTLVersionConformity raise become returned errors, per PORTING.md.
//
// NOTE: Java's covariant return type (TLParsingResult, narrowing ParsingTask's ParsingResult) is
// kept here; Go has no covariance, so a caller holding this task through the ParsingTask interface
// needs the trivial adapter the porter notes describe.
func (t *TLParsingTask) Get() (*TLParsingResult, error) {
	result := NewTLParsingResult()
	jaxbObject, err := t.JAXBObject()
	if err != nil {
		return nil, err
	}

	t.parseSchemeInformation(result, jaxbObject.SchemeInformation)
	t.parseTrustServiceProviderList(result, jaxbObject.TrustServiceProviderList)
	if err := t.VerifyTLVersionConformity(result, result.Version(), t.tlSource.TLVersions()); err != nil {
		return nil, err
	}

	return result, nil
}

// parseSchemeInformation ports the private parseSchemeInformation(TLParsingResult,
// TSLSchemeInformationType).
func (t *TLParsingTask) parseSchemeInformation(result *TLParsingResult,
	schemeInformation *jaxb.TSLSchemeInformationType) {
	t.CommonParseSchemeInformation(&result.AbstractTLParsingResult, schemeInformation)
}

// parseTrustServiceProviderList ports the private
// parseTrustServiceProviderList(TLParsingResult, TrustServiceProviderListType).
func (t *TLParsingTask) parseTrustServiceProviderList(result *TLParsingResult,
	trustServiceProviderList *jaxb.TrustServiceProviderListType) {
	if trustServiceProviderList != nil && utils.IsCollectionNotEmpty(trustServiceProviderList.TrustServiceProvider) {
		filteredTrustServiceProviders := t.filter(trustServiceProviderList.TrustServiceProvider)
		converter := NewTrustServiceProviderConverter().SetTerritory(result.Territory())
		trustServiceProviders := make([]*tslmodel.TrustServiceProvider, 0, len(filteredTrustServiceProviders))
		for _, tspType := range filteredTrustServiceProviders {
			trustServiceProviders = append(trustServiceProviders, converter.Apply(tspType))
		}
		result.SetTrustServiceProviders(trustServiceProviders)
	} else {
		result.SetTrustServiceProviders([]*tslmodel.TrustServiceProvider{})
	}
}

// filter ports the private filter(List<TSPType>), including its in-place mutation of each
// surviving TSPType's TSPServices list (step 3/4).
func (t *TLParsingTask) filter(trustServiceProviders []*jaxb.TSPType) []*jaxb.TSPType {
	filteredTSP := trustServiceProviders

	// 1. Remove TSPs with invalid structure
	nonEmptyTSPInformation := NewNonEmptyTSPInformation()
	kept := make([]*jaxb.TSPType, 0, len(filteredTSP))
	for _, tspType := range filteredTSP {
		if nonEmptyTSPInformation.Test(tspType) {
			kept = append(kept, tspType)
		}
	}
	filteredTSP = kept

	// 2. Filter the TSP with the predicate
	if t.tlSource.TrustServiceProviderPredicate() != nil {
		kept = make([]*jaxb.TSPType, 0, len(filteredTSP))
		for _, tspType := range filteredTSP {
			if t.tlSource.TrustServiceProviderPredicate().Test(tspType) {
				kept = append(kept, tspType)
			}
		}
		filteredTSP = kept
	}

	// 3. Foreach TSP, remove invalid trust services
	nonEmptyServiceInformation := NewNonEmptyServiceInformation()
	for _, tspType := range filteredTSP {
		tspServices := tspType.TSPServices
		if tspServices != nil && utils.IsCollectionNotEmpty(tspServices.TSPService) {
			filteredTrustServices := make([]*jaxb.TSPServiceType, 0, len(tspServices.TSPService))
			for _, tspService := range tspServices.TSPService {
				if nonEmptyServiceInformation.Test(tspService) {
					filteredTrustServices = append(filteredTrustServices, tspService)
				}
			}

			// 4. Filter the trust services with the predicate
			if t.tlSource.TrustServicePredicate() != nil {
				keptServices := make([]*jaxb.TSPServiceType, 0, len(filteredTrustServices))
				for _, tspService := range filteredTrustServices {
					if t.tlSource.TrustServicePredicate().Test(tspService) {
						keptServices = append(keptServices, tspService)
					}
				}
				filteredTrustServices = keptServices
			}

			newTspServices := &jaxb.TSPServicesListType{}
			if len(filteredTrustServices) != 0 {
				newTspServices.TSPService = append(newTspServices.TSPService, filteredTrustServices...)
			}
			tspType.TSPServices = newTspServices
		}
	}

	// 5. Remove TSPs with empty trust services
	nonEmptyTrustService := NewNonEmptyTrustService()
	kept = make([]*jaxb.TSPType, 0, len(filteredTSP))
	for _, tspType := range filteredTSP {
		if nonEmptyTrustService.Test(tspType) {
			kept = append(kept, tspType)
		}
	}
	return kept
}

// CreateTrustedListFacade keeps the base's TrustedListFacade: TLParsingTask does not override
// createTrustedListFacade(). The method is restated so that *TLParsingTask - and not the embedded
// base value - is what InitAbstractParsingTask registers, which matters only for LOTLParsingTask
// but is spelled out here for symmetry.
func (t *TLParsingTask) CreateTrustedListFacade() trustedListFacade {
	return t.AbstractParsingTaskBase.CreateTrustedListFacade()
}
