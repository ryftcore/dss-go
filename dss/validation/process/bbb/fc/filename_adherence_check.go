// Ported from dss-validation/.../validation/process/bbb/fc/checks/FilenameAdherenceCheck.java (DSS 6.5.RC1).
package fc

import (
	"strings"

	drjaxb "github.com/utain/esig/dss/detailedreport/jaxb"
	"github.com/utain/esig/dss/diagnostic"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/i18n"
	policy "github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/validation/process"
)

// Filename/path constants shared by the ASiC filename-adherence checks.
const (
	MetaInfFolder           = "META-INF/"
	SignatureFilename       = "signature"
	TimestampFilename       = "timestamp"
	EvidenceRecordFilename  = "evidencerecord"
	ArchiveManifestFilename = "ASiCArchiveManifest"
	AsiceMetainfManifest    = MetaInfFolder + "manifest.xml"
	MetainfAsicManifest     = MetaInfFolder + "ASiCManifest"
	XMLExtension            = ".xml"
)

// FilenameAdherenceCheck verifies validity of the token's filename according to the ASiC
// specification. T is the AbstractSignatureWrapper implementation (signature or timestamp).
type FilenameAdherenceCheck[T diagnostic.AbstractSignatureWrapperOverrides] struct {
	*process.ChainItemBase[*drjaxb.XmlFC]

	Token          T
	DiagnosticData *diagnostic.DiagnosticData
}

// InitFilenameAdherenceCheck wires the shared state; called by the concrete constructor
// before InitChainItem.
func (c *FilenameAdherenceCheck[T]) InitFilenameAdherenceCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*drjaxb.XmlFC],
	diagnosticData *diagnostic.DiagnosticData, token T, constraint policy.LevelRule) {
	c.Token = token
	c.DiagnosticData = diagnosticData
	c.ChainItemBase = process.NewChainItemBase(i18nProvider, result, constraint)
}

// IsASiCManifest checks if the filename corresponds to the "META-INF/ASiCManifest*.xml" pattern.
func (c *FilenameAdherenceCheck[T]) IsASiCManifest(filename string) bool {
	return strings.HasPrefix(filename, MetainfAsicManifest) && strings.HasSuffix(filename, XMLExtension)
}

// IsASiCArchiveManifest checks if the filename corresponds to the
// "META-INF/*ASiCArchiveManifest*.xml" pattern.
func (c *FilenameAdherenceCheck[T]) IsASiCArchiveManifest(filename string) bool {
	return strings.HasPrefix(filename, MetaInfFolder) && strings.Contains(filename, ArchiveManifestFilename) &&
		strings.HasSuffix(filename, XMLExtension)
}

// FailedIndicationForConclusion returns the Indication on failure.
func (c *FilenameAdherenceCheck[T]) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.Indication_FAILED
}

// FailedSubIndicationForConclusion returns the SubIndication on failure.
func (c *FilenameAdherenceCheck[T]) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndication_FORMAT_FAILURE
}
