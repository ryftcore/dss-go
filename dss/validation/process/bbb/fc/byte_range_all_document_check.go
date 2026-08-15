// Ported from dss-validation/.../validation/process/bbb/fc/checks/ByteRangeAllDocumentCheck.java (DSS 6.5.RC1).
package fc

import (
	drjaxb "github.com/utain/esig/dss/detailedreport/jaxb"
	"github.com/utain/esig/dss/diagnostic"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/i18n"
	policy "github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/validation/process"
)

// ByteRangeAllDocumentCheck verifies if all signatures and document timestamps present in a PDF are valid.
type ByteRangeAllDocumentCheck struct {
	*process.ChainItemBase[*drjaxb.XmlFC]

	diagnosticData *diagnostic.DiagnosticData
}

// NewByteRangeAllDocumentCheck is the default constructor.
func NewByteRangeAllDocumentCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*drjaxb.XmlFC],
	diagnosticData *diagnostic.DiagnosticData, constraint policy.LevelRule) *ByteRangeAllDocumentCheck {
	c := &ByteRangeAllDocumentCheck{diagnosticData: diagnosticData}
	c.ChainItemBase = process.NewChainItemBase(i18nProvider, result, constraint)
	c.InitChainItem(c)
	return c
}

// Process performs the check.
func (c *ByteRangeAllDocumentCheck) Process() bool {
	for _, signature := range c.diagnosticData.Signatures() {
		if signature.PDFRevision() != nil && !signature.IsSignatureByteRangeValid() {
			return false
		}
	}
	for _, timestamp := range c.diagnosticData.TimestampList() {
		if timestamp.PDFRevision() != nil && !timestamp.IsSignatureByteRangeValid() {
			return false
		}
	}
	return true
}

// MessageTag returns the constraint message i18n key.
func (c *ByteRangeAllDocumentCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_FC_DASTHVBR
}

// ErrorMessageTag returns the error message i18n key.
func (c *ByteRangeAllDocumentCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_FC_DASTHVBR_ANS
}

// FailedIndicationForConclusion returns the Indication on failure.
func (c *ByteRangeAllDocumentCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.Indication_FAILED
}

// FailedSubIndicationForConclusion returns the SubIndication on failure.
func (c *ByteRangeAllDocumentCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndication_FORMAT_FAILURE
}
