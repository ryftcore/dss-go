// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/certificate/checks/CaQcCheck.java (DSS 6.5.RC1).
//
// FLAGGED HASH-ORDER SITE: getStis() collects into a HashSet<String> in
// Java, whose bucket order is not reproducible here; identifiers are instead
// sorted lexicographically for a deterministic (if not necessarily
// Java-bucket-identical) BuildErrorMessage argument on the multi-value
// branch. See the porter brief's hard rule on hash-order leaks.
package qualification

import (
	"sort"
	"strings"

	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/utils"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// CaQcCheck checks whether there are CA/QC TrustServices.
type CaQcCheck struct {
	*process.ChainItemBase[*jaxb.XmlValidationCertificateQualification]

	// trustServices is the list of TrustServiceWrappers at control time.
	trustServices []*diagnostic.TrustServiceWrapper
}

// NewCaQcCheck is the default constructor. Port of
// CaQcCheck(I18nProvider, XmlValidationCertificateQualification, List, LevelRule).
func NewCaQcCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlValidationCertificateQualification],
	trustServices []*diagnostic.TrustServiceWrapper, constraint policy.LevelRule) *CaQcCheck {
	c := &CaQcCheck{
		ChainItemBase: process.NewChainItemBase(i18nProvider, result, constraint),
		trustServices: trustServices,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *CaQcCheck) Process() bool {
	filterByCaQc := TrustServicesFilterFactoryCreateFilterByCaQc()
	return utils.IsCollectionNotEmpty(filterByCaQc.Filter(c.trustServices))
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *CaQcCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTagQualHasCAQC
}

// BuildErrorMessage builds an error message. Port of buildErrorMessage().
func (c *CaQcCheck) BuildErrorMessage() *jaxb.XmlMessage {
	stiList := c.getStis()
	errorTag := i18n.MessageTagQualHasCAQCANS
	var argument string
	if len(stiList) == 1 {
		argument = stiList[0]
	} else {
		errorTag = i18n.MessageTagQualHasCAQCANS2
		argument = "[" + strings.Join(stiList, ", ") + "]"
	}
	return c.BuildXmlMessage(errorTag, argument)
}

// getStis ports the private getStis(). See the file header for the flagged
// HashSet iteration-order deviation.
func (c *CaQcCheck) getStis() []string {
	identifiers := map[string]struct{}{}
	var ordered []string
	for _, trustService := range c.trustServices {
		t := trustService.Type
		sti := ServiceTypeIdentifierFromUri(t)
		id := t
		if sti != "" {
			id = sti.ShortName()
		}
		if _, ok := identifiers[id]; !ok {
			identifiers[id] = struct{}{}
			ordered = append(ordered, id)
		}
	}
	sort.Strings(ordered)
	return ordered
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *CaQcCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationFailed
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion(), whose default is null.
func (c *CaQcCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return ""
}
