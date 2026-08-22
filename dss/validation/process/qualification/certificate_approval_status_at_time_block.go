// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/certificate/usage/CertificateApprovalStatusAtTimeBlock.java (DSS 6.5.RC1).
//
// Java's slf4j logging (the service-type-identifier / service-status
// conflict warnings) has no Go equivalent and is not ported.
package qualification

import (
	"fmt"
	"time"

	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// CertificateApprovalStatusAtTimeBlock verifies certificate's approval
// status at the given time.
type CertificateApprovalStatusAtTimeBlock struct {
	*process.ChainBase[*jaxb.XmlValidationCertificateApprovalStatus]

	validationTime     enumerations.ValidationTime
	date               *time.Time
	listTypeUri        string
	stiUri             string
	acceptableServices []*diagnostic.TrustedEntityServiceWrapper
	filteredServices   []*diagnostic.TrustedEntityServiceWrapper
}

// NewCertificateApprovalStatusAtTimeBlockAtIssuanceTime is the constructor
// to instantiate the validation at the certificate's issuance time. Port of
// CertificateApprovalStatusAtTimeBlock(Provider, ValidationTime, CertificateWrapper, String, String, List).
func NewCertificateApprovalStatusAtTimeBlockAtIssuanceTime(i18nProvider *i18n.Provider, validationTime enumerations.ValidationTime,
	signingCertificate *diagnostic.CertificateWrapper, listTypeUri, stiUri string,
	acceptableServices []*diagnostic.TrustedEntityServiceWrapper) *CertificateApprovalStatusAtTimeBlock {
	return NewCertificateApprovalStatusAtTimeBlock(i18nProvider, validationTime, nil, signingCertificate, listTypeUri, stiUri, acceptableServices)
}

// NewCertificateApprovalStatusAtTimeBlock is the constructor to instantiate
// the validation at the validation time. Port of
// CertificateApprovalStatusAtTimeBlock(Provider, ValidationTime, Date, CertificateWrapper, String, String, List).
func NewCertificateApprovalStatusAtTimeBlock(i18nProvider *i18n.Provider, validationTime enumerations.ValidationTime, date *time.Time,
	signingCertificate *diagnostic.CertificateWrapper, listTypeUri, stiUri string,
	acceptableServices []*diagnostic.TrustedEntityServiceWrapper) *CertificateApprovalStatusAtTimeBlock {
	xmlResult := &jaxb.XmlValidationCertificateApprovalStatus{}
	c := &CertificateApprovalStatusAtTimeBlock{
		ChainBase: process.NewChainBase(i18nProvider, process.NewResult(xmlResult,
			&xmlResult.XmlConstraintsConclusionContent, &xmlResult.XmlConstraintsConclusionAttrs)),
		validationTime:     validationTime,
		listTypeUri:        listTypeUri,
		stiUri:             stiUri,
		acceptableServices: append([]*diagnostic.TrustedEntityServiceWrapper{}, acceptableServices...),
	}
	c.Result.Value.Id = signingCertificate.Id()

	switch validationTime {
	case enumerations.ValidationTimeCertificateIssuanceTime:
		c.date = signingCertificate.NotBefore()
	case enumerations.ValidationTimeBESTSignatureTime, enumerations.ValidationTimeValidationTime:
		c.date = date
	default:
		panic(fmt.Sprintf("Unknown qualification time : %s", validationTime))
	}

	c.InitChainBase(c)
	return c
}

// BuildChainTitle builds the chain title. Port of buildChainTitle().
func (c *CertificateApprovalStatusAtTimeBlock) BuildChainTitle() string {
	message := i18n.MessageTagCertUsageAtTime
	param, err := process.GetValidationTimeMessageTag(c.validationTime)
	if err != nil {
		panic(err)
	}
	return c.I18nProvider.GetMessage(message, c.getStiUserFriendlyLabel(), param)
}

// getStiUserFriendlyLabel ports the private getStiUserFriendlyLabel().
func (c *CertificateApprovalStatusAtTimeBlock) getStiUserFriendlyLabel() string {
	sti := enumerations.LoTEServiceTypeIdentifierFromURI(c.stiUri)
	if sti != nil && sti.Label() != "" {
		return sti.Label()
	}
	return c.stiUri
}

// InitChain initializes the chain. Port of initChain().
func (c *CertificateApprovalStatusAtTimeBlock) InitChain() {

	// Init internal variable to the provided list of extracted Trust Services
	c.filteredServices = append([]*diagnostic.TrustedEntityServiceWrapper{}, c.acceptableServices...)

	// 1b. Filter by date
	filter := TrustedEntitiesFilterFactoryCreateFilterByDate(c.date)
	c.filteredServices = filter.Filter(c.filteredServices)

	item := c.hasTrustedServiceAtTime(c.filteredServices)
	c.FirstItem = item

	// 2a. Check sti consistency

	item = item.SetNextItem(c.trustedServiceTypeIdentifierKnown(c.stiUri))

	filter = TrustedEntitiesFilterFactoryCreateFilterByServiceTypeIdentifierUri(c.stiUri)
	c.filteredServices = filter.Filter(c.filteredServices)

	item = item.SetNextItem(c.trustedServicesWithSti(c.filteredServices))

	// 2b. Check status consistency
	item = item.SetNextItem(c.trustedServicesStatusConsistent(c.filteredServices))

	serviceStatusUri := c.getServiceStatusUri(c.filteredServices)
	if serviceStatusUri != nil {
		item = item.SetNextItem(c.trustedServiceStatusKnown(*serviceStatusUri)) //nolint:staticcheck // mirrors upstream CertificateApprovalStatusAtTimeBlock#initChain: Java's trailing `item = item.setNextItem(...)` is the same dead store - setNextItem links the item and returns it, and nothing reads the tail afterwards.
	}
	// NOTE: status can be null, validate successfully in this case
}

// AddAdditionalInfo adds additional info to the chain. Port of the overridden
// protected void addAdditionalInfo().
func (c *CertificateApprovalStatusAtTimeBlock) AddAdditionalInfo() {
	certificateApprovalStatus := &jaxb.XmlCertificateApprovalStatus{}
	listType := enumerations.ListTypeFromURI(c.listTypeUri)
	certificateApprovalStatus.ListType = listType

	serviceStiUriPtr := c.getServiceStiUri(c.filteredServices)
	var serviceStiUri string
	if serviceStiUriPtr != nil {
		serviceStiUri = *serviceStiUriPtr
	}
	sti := enumerations.LoTEServiceTypeIdentifierFromURI(serviceStiUri)
	certificateApprovalStatus.ServiceTypeIdentifier = sti

	serviceStatusUriPtr := c.getServiceStatusUri(c.filteredServices)
	var serviceStatusUri string
	if serviceStatusUriPtr != nil {
		serviceStatusUri = *serviceStatusUriPtr
	}
	status := enumerations.LoTEServiceStatusFromURI(serviceStatusUri)
	certificateApprovalStatus.ServiceStatus = status

	certApprovalStatus := enumerations.CertificateApprovalStatusFromDefinition(listType, sti, status)
	if certApprovalStatus == nil {
		certApprovalStatus = enumerations.CertificateApprovalStatusEnumCertForUnknown
	}
	label := certApprovalStatus.Label()
	certificateApprovalStatus.Label = &label
	c.Result.Value.CertificateApprovalStatus = certificateApprovalStatus

	vt := jaxb.ValidationTimeValue(c.validationTime)
	c.Result.Value.ValidationTime = &vt
	if c.date != nil {
		c.Result.Value.DateTime = jaxb.NewXSDateTime(*c.date)
	}
}

func (c *CertificateApprovalStatusAtTimeBlock) hasTrustedServiceAtTime(trustedServices []*diagnostic.TrustedEntityServiceWrapper) process.ChainItem[*jaxb.XmlValidationCertificateApprovalStatus] {
	return NewTrustedEntityServiceAtTimeCheck(c.I18nProvider, c.Result, trustedServices, c.validationTime, c.FailLevelRule())
}

func (c *CertificateApprovalStatusAtTimeBlock) trustedServicesWithSti(trustedServices []*diagnostic.TrustedEntityServiceWrapper) process.ChainItem[*jaxb.XmlValidationCertificateApprovalStatus] {
	return NewTrustedEntityServiceWithStiCheck(c.I18nProvider, c.Result, trustedServices, c.stiUri, c.FailLevelRule())
}

func (c *CertificateApprovalStatusAtTimeBlock) trustedServiceTypeIdentifierKnown(serviceStatusUri string) process.ChainItem[*jaxb.XmlValidationCertificateApprovalStatus] {
	return NewTrustedEntityServiceTypeIdentifierKnownCheck(c.I18nProvider, c.Result, serviceStatusUri, c.WarnLevelRule())
}

func (c *CertificateApprovalStatusAtTimeBlock) trustedServicesStatusConsistent(trustedServices []*diagnostic.TrustedEntityServiceWrapper) process.ChainItem[*jaxb.XmlValidationCertificateApprovalStatus] {
	return NewTrustedEntityServiceStatusConsistencyCheck(c.I18nProvider, c.Result, trustedServices, c.FailLevelRule())
}

func (c *CertificateApprovalStatusAtTimeBlock) trustedServiceStatusKnown(serviceStatusUri string) process.ChainItem[*jaxb.XmlValidationCertificateApprovalStatus] {
	return NewTrustedEntityServiceStatusKnownCheck(c.I18nProvider, c.Result, serviceStatusUri, c.WarnLevelRule())
}

// getServiceStiUri ports the private getServiceStiUri(List); the LOG.warn on
// a detected conflict is dropped per PORTING.md.
func (c *CertificateApprovalStatusAtTimeBlock) getServiceStiUri(filteredServices []*diagnostic.TrustedEntityServiceWrapper) *string {
	return distinctSingleValue(filteredServices, func(s *diagnostic.TrustedEntityServiceWrapper) string { return s.Type })
}

// getServiceStatusUri ports the private getServiceStatusUri(List); the
// LOG.warn on a detected conflict is dropped per PORTING.md.
func (c *CertificateApprovalStatusAtTimeBlock) getServiceStatusUri(filteredServices []*diagnostic.TrustedEntityServiceWrapper) *string {
	return distinctSingleValue(filteredServices, func(s *diagnostic.TrustedEntityServiceWrapper) string { return s.Status })
}

// distinctSingleValue ports the shared shape of getServiceStiUri/getServiceStatusUri:
// nil if the set of distinct field values is empty, the single value if
// there is exactly one, and a literal "?" (matching Java's fallback on a
// detected conflict) if there is more than one.
func distinctSingleValue(services []*diagnostic.TrustedEntityServiceWrapper, field func(*diagnostic.TrustedEntityServiceWrapper) string) *string {
	seen := map[string]struct{}{}
	for _, s := range services {
		seen[field(s)] = struct{}{}
	}
	switch len(seen) {
	case 0:
		return nil
	case 1:
		for v := range seen {
			if v == "" {
				// Java's set here holds the single element null (an absent
				// <ServiceType>/<Status> leaves the wrapper field null), and the
				// method hands that null straight back; this port carries an
				// absent value as the empty string, so it has to answer nil
				// rather than a pointer to "". The caller's
				// `if serviceStatusUri != nil` guard around
				// TrustedEntityServiceStatusKnownCheck depends on the
				// difference: with a non-nil "" it appended a
				// CERT_USAGE_STATUS_KNOWN WARNING constraint upstream never
				// emits.
				return nil
			}
			return &v
		}
	}
	conflict := "?"
	return &conflict
}
