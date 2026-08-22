// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/timestamp/TimestampQualificationAtTimeBlock.java (DSS 6.5.RC1).
package qualification

import (
	"fmt"
	"time"

	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/utils"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// TimestampQualificationAtTimeBlock verifies timestamp's qualification at
// the given time.
type TimestampQualificationAtTimeBlock struct {
	*process.ChainBase[*jaxb.XmlValidationTimestampQualificationAtTime]

	// validationTime is the time type to get the qualification at.
	validationTime enumerations.ValidationTime

	// date is the time to check against; nil is Java's null.
	date *time.Time

	// acceptableServices is the list of matching TrustServices.
	acceptableServices []*diagnostic.TrustServiceWrapper

	// tstQualif is the determined timestamp qualification.
	tstQualif enumerations.TimestampQualification
}

// NewTimestampQualificationAtTimeBlockAtGenerationTime is the constructor
// for validation at the timestamp generation time. Port of
// TimestampQualificationAtTimeBlock(I18nProvider, ValidationTime, TimestampWrapper, List).
func NewTimestampQualificationAtTimeBlockAtGenerationTime(i18nProvider *i18n.I18nProvider, validationTime enumerations.ValidationTime,
	timestamp *diagnostic.TimestampWrapper, acceptableServices []*diagnostic.TrustServiceWrapper) *TimestampQualificationAtTimeBlock {
	return NewTimestampQualificationAtTimeBlock(i18nProvider, validationTime, nil, timestamp, acceptableServices)
}

// NewTimestampQualificationAtTimeBlock is the constructor with a custom
// validation date. Port of
// TimestampQualificationAtTimeBlock(I18nProvider, ValidationTime, Date, TimestampWrapper, List).
func NewTimestampQualificationAtTimeBlock(i18nProvider *i18n.I18nProvider, validationTime enumerations.ValidationTime, date *time.Time,
	timestamp *diagnostic.TimestampWrapper, acceptableServices []*diagnostic.TrustServiceWrapper) *TimestampQualificationAtTimeBlock {
	xmlResult := &jaxb.XmlValidationTimestampQualificationAtTime{}
	c := &TimestampQualificationAtTimeBlock{
		ChainBase: process.NewChainBase(i18nProvider, process.NewResult(xmlResult,
			&xmlResult.XmlConstraintsConclusionContent, &xmlResult.XmlConstraintsConclusionAttrs)),
		validationTime:     validationTime,
		acceptableServices: acceptableServices,
		tstQualif:          enumerations.TimestampQualificationNA,
	}

	switch validationTime {
	case enumerations.ValidationTimeTimestampGenerationTime:
		c.date = timestamp.ProductionTime()
	case enumerations.ValidationTimeTimestampPOETime:
		c.date = date
	default:
		panic(fmt.Sprintf("Unsupported time-stamp qualification time : %s", validationTime))
	}

	c.InitChainBase(c)
	return c
}

// BuildChainTitle builds the chain title. Port of buildChainTitle().
func (c *TimestampQualificationAtTimeBlock) BuildChainTitle() string {
	message := i18n.MessageTagTSTQualificationAtTime
	param, err := process.GetValidationTimeMessageTag(c.validationTime)
	if err != nil {
		panic(err)
	}
	return c.I18nProvider.GetMessage(message, param)
}

// InitChain initializes the chain. Port of initChain().
func (c *TimestampQualificationAtTimeBlock) InitChain() {
	filteredServices := append([]*diagnostic.TrustServiceWrapper{}, c.acceptableServices...)

	var item process.ChainItem[*jaxb.XmlValidationTimestampQualificationAtTime]

	// Execute only for Trusted Lists with defined MRA
	if c.isMRAEnactedForTrustedList(filteredServices) {
		filter := TrustServicesFilterFactoryCreateMRAEnactedFilter()
		filteredServices = filter.Filter(filteredServices)

		filter = TrustServicesFilterFactoryCreateFilterByMRAEquivalenceStartingDate(c.date)
		filteredServices = filter.Filter(filteredServices)

		item = c.hasMraEnactedTrustService(filteredServices)
		c.FirstItem = item
	}

	// 1. filter by service for QTST
	filter := TrustServicesFilterFactoryCreateFilterByQTST()
	qtstServices := filter.Filter(filteredServices)

	if item == nil {
		item = c.hasQTST(qtstServices)
		c.FirstItem = item
	} else {
		item = item.SetNextItem(c.hasQTST(qtstServices))
	}

	// 2. filter by granted
	filter = TrustServicesFilterFactoryCreateFilterByGranted()
	grantedServices := filter.Filter(qtstServices)

	item = item.SetNextItem(c.hasGrantedStatus(grantedServices))

	// 3. filter by date (generation time)
	filter = TrustServicesFilterFactoryCreateFilterByDate(c.date)
	grantedAtDateServices := filter.Filter(grantedServices)

	item = item.SetNextItem(c.hasGrantedStatusAtDate(grantedAtDateServices)) //nolint:staticcheck // mirrors upstream TimestampQualificationAtTimeBlock#initChain: Java's trailing `item = item.setNextItem(...)` is the same dead store - setNextItem links the item and returns it, and nothing reads the tail afterwards.

	// Determine qualification status
	if utils.IsCollectionNotEmpty(grantedAtDateServices) {
		c.tstQualif = enumerations.TimestampQualificationQTSA
	} else {
		c.tstQualif = enumerations.TimestampQualificationTSA
	}
}

// AddAdditionalInfo adds additional info to the chain. Port of the overridden
// protected void addAdditionalInfo().
func (c *TimestampQualificationAtTimeBlock) AddAdditionalInfo() {
	tq := jaxb.TimestampQualificationValue(c.tstQualif)
	c.Result.Value.TimestampQualification = tq
	vt := jaxb.ValidationTimeValue(c.validationTime)
	c.Result.Value.ValidationTime = &vt
	if c.date != nil {
		c.Result.Value.DateTime = jaxb.NewXSDateTime(*c.date)
	}
}

func (c *TimestampQualificationAtTimeBlock) hasMraEnactedTrustService(services []*diagnostic.TrustServiceWrapper) process.ChainItem[*jaxb.XmlValidationTimestampQualificationAtTime] {
	return NewRelatedToMraEnactedTrustServiceCheck(c.I18nProvider, c.Result, services, c.FailLevelRule())
}

func (c *TimestampQualificationAtTimeBlock) hasQTST(services []*diagnostic.TrustServiceWrapper) process.ChainItem[*jaxb.XmlValidationTimestampQualificationAtTime] {
	return NewQTSTCheck(c.I18nProvider, c.Result, services, c.FailLevelRule())
}

func (c *TimestampQualificationAtTimeBlock) hasGrantedStatus(services []*diagnostic.TrustServiceWrapper) process.ChainItem[*jaxb.XmlValidationTimestampQualificationAtTime] {
	return NewGrantedStatusCheck(c.I18nProvider, c.Result, services, c.FailLevelRule())
}

func (c *TimestampQualificationAtTimeBlock) hasGrantedStatusAtDate(services []*diagnostic.TrustServiceWrapper) process.ChainItem[*jaxb.XmlValidationTimestampQualificationAtTime] {
	return NewGrantedStatusAtTimeCheck(c.I18nProvider, c.Result, services, c.validationTime, c.FailLevelRule())
}

// isMRAEnactedForTrustedList ports the private isMRAEnactedForTrustedList(List).
func (c *TimestampQualificationAtTimeBlock) isMRAEnactedForTrustedList(trustServices []*diagnostic.TrustServiceWrapper) bool {
	for _, trustService := range trustServices {
		if trustService.TrustedList != nil && utils.IsTrue(trustService.TrustedList.Mra) {
			return true
		}
	}
	return false
}
