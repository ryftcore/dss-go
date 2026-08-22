// Ported from dss-validation/.../validation/process/bbb/fc/checks/FullScopeCheck.java (DSS 6.5.RC1).
package fc

import (
	drjaxb "github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	diagjaxb "github.com/ryftcore/dss-go/dss/diagnostic/jaxb"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	policy "github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// FullScopeCheck checks if the signature covers FULL scope documents.
type FullScopeCheck struct {
	*process.ChainItemBase[*drjaxb.XmlFC]

	signatureScopes []*diagjaxb.XmlSignatureScope
}

// NewFullScopeCheck is the default constructor.
func NewFullScopeCheck(i18nProvider *i18n.Provider, result *process.Result[*drjaxb.XmlFC],
	signatureScopes []*diagjaxb.XmlSignatureScope, constraint policy.LevelRule) *FullScopeCheck {
	c := &FullScopeCheck{signatureScopes: signatureScopes}
	c.ChainItemBase = process.NewChainItemBase(i18nProvider, result, constraint)
	c.InitChainItem(c)
	return c
}

// Process performs the check.
func (c *FullScopeCheck) Process() bool {
	for _, sigScope := range c.signatureScopes {
		if sigScope.Scope == nil || sigScope.Scope.SignatureScopeType() != enumerations.SignatureScopeTypeFull {
			return false
		}
	}
	return true
}

// MessageTag returns the constraint message i18n key.
func (c *FullScopeCheck) MessageTag() i18n.MessageTag { return i18n.MessageTagBBBFCICFD }

// ErrorMessageTag returns the error message i18n key.
func (c *FullScopeCheck) ErrorMessageTag() i18n.MessageTag { return i18n.MessageTagBBBFCICFDANS }

// FailedIndicationForConclusion returns the Indication on failure.
func (c *FullScopeCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationFailed
}

// FailedSubIndicationForConclusion returns the SubIndication on failure.
func (c *FullScopeCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationFormatFailure
}
