// Ported from dss-validation/.../validation/process/bbb/fc/TimestampFormatChecking.java (DSS 6.5.RC1).
package fc

import (
	drjaxb "github.com/utain/esig/dss/detailedreport/jaxb"
	"github.com/utain/esig/dss/diagnostic"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/i18n"
	policy "github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/validation/process"
)

// TimestampFormatChecking performs "5.2.2 Format Checking" building block execution for a
// document or container timestamp.
type TimestampFormatChecking struct {
	AbstractSignatureFormatChecking[*diagnostic.TimestampWrapper]
}

// NewTimestampFormatChecking is the default constructor.
func NewTimestampFormatChecking(i18nProvider *i18n.I18nProvider, diagnosticData *diagnostic.DiagnosticData,
	timestamp *diagnostic.TimestampWrapper, context enumerations.Context, pol policy.ValidationPolicy) *TimestampFormatChecking {
	c := &TimestampFormatChecking{}
	c.InitAbstractSignatureFormatChecking(i18nProvider, diagnosticData, timestamp, context, pol, c)
	c.InitChainBase(c)
	return c
}

// InitChain builds the constraint chain. Port of the overridden protected void initChain().
func (c *TimestampFormatChecking) InitChain() {
	var item process.ChainItem[*drjaxb.XmlFC] = c.FirstItem

	// CAdES-V3 timestamp
	if c.Token.Type() == enumerations.TimestampType_ARCHIVE_TIMESTAMP &&
		c.Token.ArchiveTimestampType() == enumerations.ArchiveTimestampType_CAdES_V3 {
		item = c.cadesAtsV3HashIndex()
		c.FirstItem = item
	}

	// PAdES
	if c.Token.PDFRevision() != nil {
		item = c.GetPDFRevisionValidationChain(item)
	}

	// PDF/A (only for a detached document timestamp)
	if c.DiagnosticData.IsPDFAValidationPerformed() && c.Token.Type().IsDocumentTimestamp() &&
		len(c.Token.TimestampedSignatures()) == 0 {
		item = c.GetPdfaValidationChain(item)
	}

	// ASiC timestamps
	if c.DiagnosticData.IsContainerInfoPresent() && c.Token.Type().IsContainerTimestamp() {

		// only for a detached container timestamp
		if len(c.Token.TimestampedSignatures()) == 0 {
			item = c.GetASiCContainerValidationChain(item)
		}

		// when signature, timestamp or evidence record is covered
		if c.coversSignatureOrTimestampOrEvidenceRecord() {
			if item == nil {
				item = c.signedAndTimestampedFilesCovered()
				c.FirstItem = item
			} else {
				item = item.SetNextItem(c.signedAndTimestampedFilesCovered())
			}
		}
	}
}

func (c *TimestampFormatChecking) cadesAtsV3HashIndex() process.ChainItem[*drjaxb.XmlFC] {
	constraint := c.Policy.AtsHashIndexConstraint()
	return NewCAdESV3HashIndexCheck(c.I18nProvider, c.Result, c.Token, constraint)
}

func (c *TimestampFormatChecking) signedAndTimestampedFilesCovered() process.ChainItem[*drjaxb.XmlFC] {
	constraint := c.Policy.TimestampContainerSignedAndTimestampedFilesCoveredConstraint()
	return NewSignedAndTimestampedFilesCoveredCheck(c.I18nProvider, c.Result, c.DiagnosticData, c.Token, constraint)
}

// FilenameAdherenceCheck creates the timestamp filename adherence check. Port of the overridden
// protected ChainItem<XmlFC> filenameAdherenceCheck().
func (c *TimestampFormatChecking) FilenameAdherenceCheck() process.ChainItem[*drjaxb.XmlFC] {
	constraint := c.Policy.FilenameAdherenceConstraint()
	return NewTimestampFilenameAdherenceCheck(c.I18nProvider, c.Result, c.DiagnosticData, c.Token, constraint)
}

// ManifestFilenameAdherenceCheck creates the timestamp manifest filename adherence check. Port
// of the overridden protected ChainItem<XmlFC> manifestFilenameAdherenceCheck().
func (c *TimestampFormatChecking) ManifestFilenameAdherenceCheck() process.ChainItem[*drjaxb.XmlFC] {
	constraint := c.Policy.FilenameAdherenceConstraint()
	return NewTimestampManifestFilenameAdherenceCheck(c.I18nProvider, c.Result, c.DiagnosticData, c.Token, constraint)
}

func (c *TimestampFormatChecking) coversSignatureOrTimestampOrEvidenceRecord() bool {
	return len(c.Token.TimestampedSignatures()) > 0 || len(c.Token.TimestampedTimestamps()) > 0 ||
		len(c.Token.TimestampedEvidenceRecords()) > 0
}
