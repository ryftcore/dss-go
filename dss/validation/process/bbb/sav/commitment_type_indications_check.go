// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/sav/checks/CommitmentTypeIndicationsCheck.java (DSS 6.5.RC1).
package sav

import (
	jaxb "github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	diagnosticjaxb "github.com/ryftcore/dss-go/dss/diagnostic/jaxb"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/utils"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// CommitmentTypeIndicationsCheck checks if the commitment type indications are
// acceptable.
type CommitmentTypeIndicationsCheck struct {
	*process.ChainItemBase[*jaxb.XmlSAV]

	// signature is the signature to check.
	signature *diagnostic.SignatureWrapper

	// constraint is the constraint.
	constraint policy.MultiValuesRule
}

// NewCommitmentTypeIndicationsCheck is the default constructor. Port of
// CommitmentTypeIndicationsCheck(Provider, XmlSAV, SignatureWrapper, MultiValuesRule).
func NewCommitmentTypeIndicationsCheck(i18nProvider *i18n.Provider, result *process.Result[*jaxb.XmlSAV],
	signature *diagnostic.SignatureWrapper, constraint policy.MultiValuesRule) *CommitmentTypeIndicationsCheck {
	c := &CommitmentTypeIndicationsCheck{
		ChainItemBase: process.NewChainItemBase(i18nProvider, result, constraint),
		signature:     signature,
		constraint:    constraint,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *CommitmentTypeIndicationsCheck) Process() bool {
	commitmentTypeIndications := c.signature.CommitmentTypeIndications()
	expectedValues := c.constraint.Values()

	if utils.IsCollectionEmpty(commitmentTypeIndications) {
		return false
	}

	if utils.IsCollectionNotEmpty(expectedValues) {
		presentIdentifiers := make([]string, len(commitmentTypeIndications))
		for i, cti := range commitmentTypeIndications {
			presentIdentifiers[i] = commitmentTypeIndicationIdentifier(cti)
		}
		return containsAll(expectedValues, presentIdentifiers)
	}

	return true
}

// commitmentTypeIndicationIdentifier reads XmlCommitmentTypeIndication#getIdentifier():
// the generated member is a pointer, whose nil is Java's null.
func commitmentTypeIndicationIdentifier(cti *diagnosticjaxb.XmlCommitmentTypeIndication) string {
	if cti.Identifier == nil {
		return ""
	}
	return *cti.Identifier
}

// containsAll ports java.util.List#containsAll(Collection): TRUE if every
// element of contained is present in container.
func containsAll(container []string, contained []string) bool {
	set := make(map[string]struct{}, len(container))
	for _, v := range container {
		set[v] = struct{}{}
	}
	for _, v := range contained {
		if _, ok := set[v]; !ok {
			return false
		}
	}
	return true
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *CommitmentTypeIndicationsCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTagBBBSAVISQPXTIP
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *CommitmentTypeIndicationsCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagBBBSAVISQPXTIPANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *CommitmentTypeIndicationsCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationIndeterminate
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure. Port
// of getFailedSubIndicationForConclusion().
func (c *CommitmentTypeIndicationsCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationSigConstraintsFailure
}
