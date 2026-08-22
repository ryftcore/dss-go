// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/certificate/checks/GrantedStatusCheck.java (DSS 6.5.RC1).
//
// Java's getStatusList() collects into a
// HashSet<String>, whose iteration order (bucket order, not insertion order)
// feeds the multi-value error message via Set#toString() when more than one
// distinct status is found. That bucket order is not reproducible statically
// in Go. This port sorts the identifiers lexicographically instead of using
// insertion order, for a deterministic (if not necessarily Java-bucket-
// identical) result.
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

// GrantedStatusCheck verifies if the certificate has TrustServices with a
// 'granted' status.
type GrantedStatusCheck[T any] struct {
	*process.ChainItemBase[T]

	// trustServicesAtTime is the list of TrustServiceWrappers at control time.
	trustServicesAtTime []*diagnostic.TrustServiceWrapper
}

// NewGrantedStatusCheck is the default constructor. Port of
// GrantedStatusCheck(I18nProvider, T, List, LevelRule).
func NewGrantedStatusCheck[T any](i18nProvider *i18n.I18nProvider, result *process.Result[T],
	trustServicesAtTime []*diagnostic.TrustServiceWrapper, constraint policy.LevelRule) *GrantedStatusCheck[T] {
	c := &GrantedStatusCheck[T]{
		ChainItemBase:       process.NewChainItemBase(i18nProvider, result, constraint),
		trustServicesAtTime: trustServicesAtTime,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *GrantedStatusCheck[T]) Process() bool {
	filterByGranted := TrustServicesFilterFactoryCreateFilterByGranted()
	return utils.IsCollectionNotEmpty(filterByGranted.Filter(c.trustServicesAtTime))
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *GrantedStatusCheck[T]) MessageTag() i18n.MessageTag {
	return i18n.MessageTagQualHasGranted
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *GrantedStatusCheck[T]) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagQualHasGrantedANS
}

// BuildErrorMessage builds an error message. Port of buildErrorMessage().
func (c *GrantedStatusCheck[T]) BuildErrorMessage() *jaxb.XmlMessage {
	statusList := c.getStatusList()
	errorTag := i18n.MessageTagQualHasGrantedANS
	var argument string
	if len(statusList) == 1 {
		argument = statusList[0]
	} else {
		errorTag = i18n.MessageTagQualHasGrantedANS2
		// Java's Set<String>#toString(). Kept file-local (no cross-file
		// helpers) rather than factored out, per the porting hard rules.
		argument = "[" + strings.Join(statusList, ", ") + "]"
	}
	return c.BuildXmlMessage(errorTag, argument)
}

// getStatusList ports the private getStatusList(). See the file header for
// the flagged HashSet iteration-order deviation.
func (c *GrantedStatusCheck[T]) getStatusList() []string {
	identifiers := map[string]struct{}{}
	var ordered []string
	for _, trustService := range c.trustServicesAtTime {
		status := trustService.Status
		tss := TrustServiceStatusFromUri(status)
		id := status
		if tss != "" {
			id = tss.ShortName()
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
func (c *GrantedStatusCheck[T]) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationFailed
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion(), whose default is null.
func (c *GrantedStatusCheck[T]) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return ""
}
