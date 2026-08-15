// Ported from dss-validation/.../validation/process/bbb/fc/AbstractFormatChecking.java (DSS 6.5.RC1).
package fc

import (
	drjaxb "github.com/utain/esig/dss/detailedreport/jaxb"
	"github.com/utain/esig/dss/diagnostic"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/i18n"
	policy "github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/validation/process"
)

// AbstractFormatChecking contains common code to be processed as a part of a "5.2.2 Format
// Checking" building block for validation of signatures and timestamps. S is the token wrapper
// type (SignatureWrapper, TimestampWrapper, EAAWrapper, EAARevocationTokenWrapper, ...).
type AbstractFormatChecking[S diagnostic.AbstractTokenProxyOverrides] struct {
	*process.ChainBase[*drjaxb.XmlFC]

	DiagnosticData *diagnostic.DiagnosticData
	Token          S
	Context        enumerations.Context
	Policy         policy.ValidationPolicy
}

// InitAbstractFormatChecking wires the shared state; called by the concrete constructor
// before InitChainBase.
func (c *AbstractFormatChecking[S]) InitAbstractFormatChecking(i18nProvider *i18n.I18nProvider,
	diagnosticData *diagnostic.DiagnosticData, token S, context enumerations.Context, pol policy.ValidationPolicy) {
	c.DiagnosticData = diagnosticData
	c.Token = token
	c.Context = context
	c.Policy = pol

	xmlFC := &drjaxb.XmlFC{}
	result := process.NewResult(xmlFC, &xmlFC.XmlConstraintsConclusionContent, &xmlFC.XmlConstraintsConclusionAttrs)
	c.ChainBase = process.NewChainBase(i18nProvider, result)
}

// Title returns the title of the chain (i.e. the BasicBuildingBlock title).
func (c *AbstractFormatChecking[S]) Title() i18n.MessageTag { return i18n.MessageTag_FORMAT_CHECKING }
