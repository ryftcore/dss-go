// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/certificate/CertQualificationAtTimeBlock.java (DSS 6.5.RC1).
//
// CROSS-CHUNK ASSUMPTION: TrustServiceFilter and the TrustServicesFilterFactory
// static factory methods (Java package qualification.trust.filter) are
// owned by a sibling porter of this shared package and were not present on
// disk while this file was written; their call sites here follow this
// package's established "Factory" static-method flattening convention
// (TrustServicesFilterFactoryCreateXxx package functions returning a
// TrustServiceFilter with a Filter([]*TrustServiceWrapper) []*TrustServiceWrapper
// method) but must be reconciled against the sibling porter's actual
// signatures once available.
package qualification

import (
	"fmt"
	"time"

	"github.com/utain/esig/dss/detailedreport/jaxb"
	"github.com/utain/esig/dss/diagnostic"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/i18n"
	"github.com/utain/esig/dss/utils"
	"github.com/utain/esig/dss/validation/process"
)

// CertQualificationAtTimeBlock verifies certificate's qualification at the
// given time.
type CertQualificationAtTimeBlock struct {
	*process.ChainBase[*jaxb.XmlValidationCertificateQualification]

	// validationTime is the time type to get the qualification at.
	validationTime enumerations.ValidationTime

	// date is the time to check against; nil is Java's null.
	date *time.Time

	// signingCertificate is the certificate to get qualification for.
	signingCertificate *diagnostic.CertificateWrapper

	// acceptableServices is the list of matching TrustServices.
	acceptableServices []*diagnostic.TrustServiceWrapper

	// certificateQualification is the internal cached variable representing
	// the qualification result.
	certificateQualification enumerations.CertificateQualification

	// filteredServices is the internal cached variable representing the
	// filtered value trust services allowed to issue qualified certificates.
	filteredServices []*diagnostic.TrustServiceWrapper

	overrides CertQualificationAtTimeBlockOverrides
}

// CertQualificationAtTimeBlockOverrides captures the member Java's
// CertQualificationAtTimeBlock treats virtually: executeQSCDCheck(),
// self-called from initChain(). CertQualificationAtTimeForQWACBlock overrides
// it and registers itself via InitCertQualificationAtTimeBlock.
type CertQualificationAtTimeBlockOverrides interface {
	// ExecuteQSCDCheck defines whether a QSCD check should be processed for
	// certificate qualification determination. NOTE: can be disabled in
	// some cases, e.g. for QWAC validation. Port of the overridable
	// protected boolean executeQSCDCheck().
	ExecuteQSCDCheck() bool
}

// NewCertQualificationAtTimeBlockAtIssuanceTime is the constructor to
// instantiate the validation at the certificate's issuance time. Port of
// CertQualificationAtTimeBlock(I18nProvider, ValidationTime, CertificateWrapper, List).
func NewCertQualificationAtTimeBlockAtIssuanceTime(i18nProvider *i18n.I18nProvider, validationTime enumerations.ValidationTime,
	signingCertificate *diagnostic.CertificateWrapper, acceptableServices []*diagnostic.TrustServiceWrapper) *CertQualificationAtTimeBlock {
	return NewCertQualificationAtTimeBlock(i18nProvider, validationTime, nil, signingCertificate, acceptableServices)
}

// NewCertQualificationAtTimeBlock is the constructor to instantiate the
// validation at the validation time. Port of
// CertQualificationAtTimeBlock(I18nProvider, ValidationTime, Date, CertificateWrapper, List).
func NewCertQualificationAtTimeBlock(i18nProvider *i18n.I18nProvider, validationTime enumerations.ValidationTime, date *time.Time,
	signingCertificate *diagnostic.CertificateWrapper, acceptableServices []*diagnostic.TrustServiceWrapper) *CertQualificationAtTimeBlock {
	xmlResult := &jaxb.XmlValidationCertificateQualification{}
	c := &CertQualificationAtTimeBlock{
		ChainBase: process.NewChainBase(i18nProvider, process.NewResult(xmlResult,
			&xmlResult.XmlConstraintsConclusionContent, &xmlResult.XmlConstraintsConclusionAttrs)),
		validationTime:           validationTime,
		signingCertificate:       signingCertificate,
		acceptableServices:       append([]*diagnostic.TrustServiceWrapper{}, acceptableServices...),
		certificateQualification: enumerations.CertificateQualification_NA,
	}
	c.Result.Value.Id = signingCertificate.Id()

	switch validationTime {
	case enumerations.ValidationTime_CERTIFICATE_ISSUANCE_TIME:
		c.date = signingCertificate.NotBefore()
	case enumerations.ValidationTime_VALIDATION_TIME, enumerations.ValidationTime_BEST_SIGNATURE_TIME:
		c.date = date
	default:
		panic(fmt.Sprintf("Unknown qualification time : %s", validationTime))
	}

	c.InitCertQualificationAtTimeBlock(c)
	return c
}

// InitCertQualificationAtTimeBlock registers the concrete chain with its
// base so that the base can dispatch to the overridden methods. It must be
// called exactly once, by the concrete chain's constructor, before Execute.
func (c *CertQualificationAtTimeBlock) InitCertQualificationAtTimeBlock(overrides CertQualificationAtTimeBlockOverrides) {
	c.overrides = overrides
	c.InitChainBase(c)
}

// BuildChainTitle builds the chain title. Port of buildChainTitle().
func (c *CertQualificationAtTimeBlock) BuildChainTitle() string {
	message := i18n.MessageTag_CERT_QUALIFICATION_AT_TIME
	param, err := process.GetValidationTimeMessageTag(c.validationTime)
	if err != nil {
		panic(err)
	}
	return c.I18nProvider.GetMessage(message, param)
}

// ExecuteQSCDCheck defines whether a QSCD check should be processed for
// certificate qualification determination. Port of executeQSCDCheck(),
// whose default is true.
func (c *CertQualificationAtTimeBlock) ExecuteQSCDCheck() bool {
	return true
}

// FilteredServices returns a list of filtered valid trust services allowed
// to issue qualified certificates. Port of getFilteredServices(); panics
// (Java's IllegalStateException) if Execute() has not run yet.
func (c *CertQualificationAtTimeBlock) FilteredServices() []*diagnostic.TrustServiceWrapper {
	if c.filteredServices == nil {
		panic("execute() method shall be called first!")
	}
	return c.filteredServices
}

// InitChain initializes the chain. Port of initChain().
func (c *CertQualificationAtTimeBlock) InitChain() {
	// Init internal variable to the provided list of extracted Trust Services
	c.filteredServices = append([]*diagnostic.TrustServiceWrapper{}, c.acceptableServices...)

	var item process.ChainItem[*jaxb.XmlValidationCertificateQualification]

	// 1a. Filter by date
	filterByDate := TrustServicesFilterFactoryCreateFilterByDate(c.date)
	c.filteredServices = filterByDate.Filter(c.filteredServices)

	// Execute only for Trusted Lists with defined MRA
	if c.isMRAEnactedForTrustedList(c.filteredServices) {
		filterByMRAEnacted := TrustServicesFilterFactoryCreateMRAEnactedFilter()
		c.filteredServices = filterByMRAEnacted.Filter(c.filteredServices)

		filterByMRAEnacted = TrustServicesFilterFactoryCreateFilterByMRAEquivalenceStartingDate(c.date)
		c.filteredServices = filterByMRAEnacted.Filter(c.filteredServices)

		item = c.hasMraEnactedTrustService(c.filteredServices)
		c.FirstItem = item

		item = item.SetNextItem(c.mraCertificateEquivalenceApplied())

	} else {
		item = c.hasTrustServiceAtTime(c.filteredServices)
		c.FirstItem = item
	}

	// 1b. Filter by service for CA/QC
	item = item.SetNextItem(c.hasCaQc(c.filteredServices))

	filterByCaQc := TrustServicesFilterFactoryCreateFilterByCaQc()
	caqcServices := filterByCaQc.Filter(c.filteredServices)

	// continue validation with available trust services if CA/QC not found
	if utils.IsCollectionNotEmpty(caqcServices) {
		c.filteredServices = caqcServices
	}

	// 2. Filter by cert type (current type or overruled)
	filterByCertificateType := TrustServicesFilterFactoryCreateFilterByCertificateType(c.signingCertificate)
	c.filteredServices = filterByCertificateType.Filter(c.filteredServices)

	item = item.SetNextItem(c.hasTrustServiceWithType(c.filteredServices))

	// 3. Run consistency checks to get warnings
	for _, trustService := range c.filteredServices {
		item = item.SetNextItem(c.serviceConsistency(trustService))
	}

	if len(c.filteredServices) > 1 {
		// 4. Simulate with all CA/QC (granted + withdrawn) to detect conflict
		var results []enumerations.CertificateQualification
		seen := map[enumerations.CertificateQualification]struct{}{}
		for _, trustService := range c.filteredServices {
			calculator := NewCertificateQualificationCalculator(c.signingCertificate, trustService)
			q := calculator.Qualification()
			if _, ok := seen[q]; !ok {
				seen[q] = struct{}{}
				results = append(results, q)
			}
		}
		item = item.SetNextItem(c.isNoConflictDetected(results))

		// interrupt in case of conflict
		if len(results) > 1 {
			c.certificateQualification = enumerations.CertificateQualification_NA
			return
		}
	}

	// 5a. Filter services with consistent status
	filterConsistentByStatus := TrustServicesFilterFactoryCreateConsistentServiceByStatusFilter()
	c.filteredServices = filterConsistentByStatus.Filter(c.filteredServices)

	// 5b. Filter by Granted
	item = item.SetNextItem(c.hasGrantedStatus(c.filteredServices))

	filterByGranted := TrustServicesFilterFactoryCreateFilterByGranted()
	grantedServices := filterByGranted.Filter(c.filteredServices)

	// continue validation with available trust services if granted not found
	if utils.IsCollectionNotEmpty(grantedServices) {
		c.filteredServices = grantedServices
	}

	// 6. Filter one trust service
	if utils.CollectionSize(c.filteredServices) > 1 {
		filterUnique := TrustServicesFilterFactoryCreateUniqueServiceFilter(c.signingCertificate)
		c.filteredServices = filterUnique.Filter(c.filteredServices)

		item = item.SetNextItem(c.isAbleToSelectOneTrustService(c.filteredServices))
	}

	var selectedTrustService *diagnostic.TrustServiceWrapper
	if len(c.filteredServices) > 0 {
		selectedTrustService = c.filteredServices[0]
	}

	// 7. Trusted certificate matches the trust service properties ?
	if selectedTrustService != nil {
		item = item.SetNextItem(c.isTrustedCertificateMatchTrustService(selectedTrustService))
	}

	// Keep only CA/QC and granted for further status determination
	if !containsTrustService(caqcServices, selectedTrustService) || !containsTrustService(grantedServices, selectedTrustService) {
		c.filteredServices = nil
		selectedTrustService = nil
	}

	item = item.SetNextItem(c.isValidCAQC(selectedTrustService))

	// 8. QC?
	filterConsistentByQC := TrustServicesFilterFactoryCreateConsistentServiceByQCFilter()
	trustServicesByQC := filterConsistentByQC.Filter(c.filteredServices)

	item = item.SetNextItem(c.hasConsistentByQCTrustService(trustServicesByQC))

	var qcSelectedTrustService *diagnostic.TrustServiceWrapper
	if len(trustServicesByQC) > 0 {
		qcSelectedTrustService = trustServicesByQC[0]
	}

	qcStrategy := CreateQualificationFromCertAndTL(c.signingCertificate, qcSelectedTrustService)
	qualifiedStatus := qcStrategy.QualifiedStatus()
	item = item.SetNextItem(c.isQualified(qualifiedStatus))

	// 9. Type?
	filterConsistentByType := TrustServicesFilterFactoryCreateConsistentServiceByCertificateTypeFilter()
	trustServicesByType := filterConsistentByType.Filter(c.filteredServices)

	item = item.SetNextItem(c.hasCertificateTypeCoverage(trustServicesByType))

	var typeSelectedTrustService *diagnostic.TrustServiceWrapper
	if len(trustServicesByType) > 0 {
		typeSelectedTrustService = trustServicesByType[0]
	}

	typeStrategy := CreateTypeFromCertAndTL(c.signingCertificate, typeSelectedTrustService, qualifiedStatus)
	certType := typeStrategy.Type()
	item = item.SetNextItem(c.certificateType(certType))

	// 11. QSCD ?
	filterConsistentByQSCD := TrustServicesFilterFactoryCreateConsistentServiceByQSCDFilter()
	trustServicesByQSCD := filterConsistentByQSCD.Filter(c.filteredServices)

	var qscdSelectedTrustService *diagnostic.TrustServiceWrapper
	if len(trustServicesByQSCD) > 0 {
		qscdSelectedTrustService = trustServicesByQSCD[0]
	}

	qscdStrategy := CreateQSCDFromCertAndTL(c.signingCertificate, qscdSelectedTrustService, qualifiedStatus)
	qscdStatus := qscdStrategy.QSCDStatus()

	if c.overrides.ExecuteQSCDCheck() {

		item = item.SetNextItem(c.hasConsistentByQSCDTrustService(trustServicesByQSCD))

		item = item.SetNextItem(c.isQscd(qscdStatus))

	}

	c.certificateQualification = GetCertQualification(qualifiedStatus, certType, qscdStatus)
}

// AddAdditionalInfo adds additional info to the chain. Port of the overridden
// protected void addAdditionalInfo().
func (c *CertQualificationAtTimeBlock) AddAdditionalInfo() {
	cq := jaxb.CertificateQualificationValue(c.certificateQualification)
	c.Result.Value.CertificateQualification = &cq
	vt := jaxb.ValidationTimeValue(c.validationTime)
	c.Result.Value.ValidationTime = &vt
	if c.date != nil {
		c.Result.Value.DateTime = jaxb.NewXSDateTime(*c.date)
	}
}

func (c *CertQualificationAtTimeBlock) hasMraEnactedTrustService(trustServices []*diagnostic.TrustServiceWrapper) process.ChainItem[*jaxb.XmlValidationCertificateQualification] {
	return NewRelatedToMraEnactedTrustServiceCheck(c.I18nProvider, c.Result, trustServices, c.FailLevelRule())
}

func (c *CertQualificationAtTimeBlock) mraCertificateEquivalenceApplied() process.ChainItem[*jaxb.XmlValidationCertificateQualification] {
	return NewMRACertificateEquivalenceApplied(c.I18nProvider, c.Result, c.signingCertificate, c.WarnLevelRule())
}

func (c *CertQualificationAtTimeBlock) hasCaQc(trustServices []*diagnostic.TrustServiceWrapper) process.ChainItem[*jaxb.XmlValidationCertificateQualification] {
	return NewCaQcCheck(c.I18nProvider, c.Result, trustServices, c.WarnLevelRule())
}

func (c *CertQualificationAtTimeBlock) hasTrustServiceAtTime(trustServices []*diagnostic.TrustServiceWrapper) process.ChainItem[*jaxb.XmlValidationCertificateQualification] {
	return NewTrustServiceAtTimeCheck(c.I18nProvider, c.Result, trustServices, c.validationTime, c.FailLevelRule())
}

func (c *CertQualificationAtTimeBlock) hasTrustServiceWithType(trustServices []*diagnostic.TrustServiceWrapper) process.ChainItem[*jaxb.XmlValidationCertificateQualification] {
	return NewTrustServicesByCertificateTypeCheck(c.I18nProvider, c.Result, trustServices, c.FailLevelRule())
}

func (c *CertQualificationAtTimeBlock) isNoConflictDetected(certificateQualificationsAtTime []enumerations.CertificateQualification) process.ChainItem[*jaxb.XmlValidationCertificateQualification] {
	return NewIsNoQualificationConflictDetectedCheck(c.I18nProvider, c.Result, certificateQualificationsAtTime, c.FailLevelRule())
}

func (c *CertQualificationAtTimeBlock) hasGrantedStatus(trustServices []*diagnostic.TrustServiceWrapper) process.ChainItem[*jaxb.XmlValidationCertificateQualification] {
	return NewGrantedStatusCheck(c.I18nProvider, c.Result, trustServices, c.WarnLevelRule())
}

func (c *CertQualificationAtTimeBlock) hasCertificateTypeCoverage(trustServices []*diagnostic.TrustServiceWrapper) process.ChainItem[*jaxb.XmlValidationCertificateQualification] {
	return NewCertificateTypeCoverageCheck(c.I18nProvider, c.Result, trustServices, c.FailLevelRule())
}

func (c *CertQualificationAtTimeBlock) hasConsistentByQCTrustService(trustServices []*diagnostic.TrustServiceWrapper) process.ChainItem[*jaxb.XmlValidationCertificateQualification] {
	return NewCertificateIssuedByConsistentByQCTrustServiceCheck(c.I18nProvider, c.Result, trustServices, c.FailLevelRule())
}

func (c *CertQualificationAtTimeBlock) hasConsistentByQSCDTrustService(trustServices []*diagnostic.TrustServiceWrapper) process.ChainItem[*jaxb.XmlValidationCertificateQualification] {
	return NewCertificateIssuedByConsistentByQSCDTrustServiceCheck(c.I18nProvider, c.Result, trustServices, c.FailLevelRule())
}

func (c *CertQualificationAtTimeBlock) isAbleToSelectOneTrustService(trustServices []*diagnostic.TrustServiceWrapper) process.ChainItem[*jaxb.XmlValidationCertificateQualification] {
	return NewIsAbleToSelectOneTrustService(c.I18nProvider, c.Result, trustServices, c.FailLevelRule())
}

func (c *CertQualificationAtTimeBlock) serviceConsistency(selectedTrustService *diagnostic.TrustServiceWrapper) process.ChainItem[*jaxb.XmlValidationCertificateQualification] {
	return NewServiceConsistencyCheck(c.I18nProvider, c.Result, selectedTrustService, c.WarnLevelRule())
}

func (c *CertQualificationAtTimeBlock) isTrustedCertificateMatchTrustService(selectedTrustService *diagnostic.TrustServiceWrapper) process.ChainItem[*jaxb.XmlValidationCertificateQualification] {
	return NewTrustedCertificateMatchTrustServiceCheck(c.I18nProvider, c.Result, selectedTrustService, c.WarnLevelRule())
}

func (c *CertQualificationAtTimeBlock) isValidCAQC(selectedTrustService *diagnostic.TrustServiceWrapper) process.ChainItem[*jaxb.XmlValidationCertificateQualification] {
	return NewValidCAQCCheck(c.I18nProvider, c.Result, selectedTrustService, c.FailLevelRule())
}

func (c *CertQualificationAtTimeBlock) isQualified(qualifiedStatus enumerations.CertificateQualifiedStatus) process.ChainItem[*jaxb.XmlValidationCertificateQualification] {
	return NewQualifiedCheck(c.I18nProvider, c.Result, qualifiedStatus, c.validationTime, c.WarnLevelRule())
}

func (c *CertQualificationAtTimeBlock) certificateType(certType enumerations.CertificateType) process.ChainItem[*jaxb.XmlValidationCertificateQualification] {
	return NewCertificateTypeCheck(c.I18nProvider, c.Result, certType, c.validationTime, c.WarnLevelRule())
}

func (c *CertQualificationAtTimeBlock) isQscd(qscdStatus enumerations.QSCDStatus) process.ChainItem[*jaxb.XmlValidationCertificateQualification] {
	return NewQSCDCheck(c.I18nProvider, c.Result, qscdStatus, c.validationTime, c.WarnLevelRule())
}

// isMRAEnactedForTrustedList ports the private isMRAEnactedForTrustedList(List).
func (c *CertQualificationAtTimeBlock) isMRAEnactedForTrustedList(trustServices []*diagnostic.TrustServiceWrapper) bool {
	for _, trustService := range trustServices {
		if trustService.TrustedList != nil && utils.IsTrue(trustService.TrustedList.Mra) {
			return true
		}
	}
	return false
}

// containsTrustService reports whether services contains target by pointer
// identity, mirroring Java's List#contains(Object) (identity-independent
// equals() is not overridden by TrustServiceWrapper, so Java also compares
// by reference here); a nil target is never contained, matching
// List#contains(null) returning false for a list without a null element.
func containsTrustService(services []*diagnostic.TrustServiceWrapper, target *diagnostic.TrustServiceWrapper) bool {
	if target == nil {
		return false
	}
	for _, s := range services {
		if s == target {
			return true
		}
	}
	return false
}
