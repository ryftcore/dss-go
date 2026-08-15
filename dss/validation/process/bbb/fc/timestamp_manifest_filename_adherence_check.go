// Ported from dss-validation/.../validation/process/bbb/fc/checks/TimestampManifestFilenameAdherenceCheck.java (DSS 6.5.RC1).
package fc

import (
	"strings"

	drjaxb "github.com/utain/esig/dss/detailedreport/jaxb"
	"github.com/utain/esig/dss/diagnostic"
	diagjaxb "github.com/utain/esig/dss/diagnostic/jaxb"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/i18n"
	policy "github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/validation/process"
)

// lastArchiveManifestFilename is the ASiC Archive Manifest name (META-INF/ASiCArchiveManifest.xml).
const lastArchiveManifestFilename = MetaInfFolder + ArchiveManifestFilename + XMLExtension

// TimestampManifestFilenameAdherenceCheck verifies conformance of the manifest filename related
// to a timestamp.
type TimestampManifestFilenameAdherenceCheck struct {
	FilenameAdherenceCheck[*diagnostic.TimestampWrapper]
}

// NewTimestampManifestFilenameAdherenceCheck is the default constructor.
func NewTimestampManifestFilenameAdherenceCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*drjaxb.XmlFC],
	diagnosticData *diagnostic.DiagnosticData, token *diagnostic.TimestampWrapper,
	constraint policy.LevelRule) *TimestampManifestFilenameAdherenceCheck {
	c := &TimestampManifestFilenameAdherenceCheck{}
	c.InitFilenameAdherenceCheck(i18nProvider, result, diagnosticData, token, constraint)
	c.InitChainItem(c)
	return c
}

func coversArchivalContent(manifestFile *diagjaxb.XmlManifestFile) bool {
	for _, entry := range manifestFile.Entries.All() {
		if strings.HasPrefix(entry, MetaInfFolder) &&
			(strings.Contains(entry, SignatureFilename) || strings.Contains(entry, TimestampFilename) ||
				strings.Contains(entry, EvidenceRecordFilename)) {
			return true
		}
	}
	return false
}

func coversFilename(manifestFile *diagjaxb.XmlManifestFile, filename string) bool {
	if manifestFile == nil {
		return false
	}
	for _, entry := range manifestFile.Entries.All() {
		if entry == filename {
			return true
		}
	}
	return false
}

func (c *TimestampManifestFilenameAdherenceCheck) isLastArchivalTimestamp() bool {
	for _, timestampWrapper := range c.DiagnosticData.TimestampList() {
		if timestampWrapper.Filename() != "" && timestampWrapper.Filename() != c.Token.Filename() {
			tstManifest := c.DiagnosticData.ManifestFileForFilename(timestampWrapper.Filename())
			if coversFilename(tstManifest, c.Token.Filename()) {
				return false
			}
		}
	}
	return true
}

// Process performs the check.
func (c *TimestampManifestFilenameAdherenceCheck) Process() bool {
	if c.DiagnosticData.ContainerType() == enumerations.ASiCContainerType_ASiC_S {
		// 4.3.3.2 Contents of the container: the META-INF folder may contain other
		// application specific information - can be of any format.
		return true
	}

	manifestFile := c.DiagnosticData.ManifestFileForFilename(c.Token.Filename())
	if manifestFile == nil {
		return false // required for a timestamp
	}

	manifestFilename := ""
	if manifestFile.Filename != nil {
		manifestFilename = *manifestFile.Filename
	}
	if strings.TrimSpace(manifestFilename) == "" {
		return false
	}

	if coversArchivalContent(manifestFile) {
		if c.isLastArchivalTimestamp() {
			return lastArchiveManifestFilename == manifestFilename
		}
		return c.IsASiCArchiveManifest(manifestFilename)
	}
	return c.IsASiCManifest(manifestFilename)
}

// MessageTag returns the constraint message i18n key.
func (c *TimestampManifestFilenameAdherenceCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_FC_IMFCS
}

// ErrorMessageTag returns the error message i18n key.
func (c *TimestampManifestFilenameAdherenceCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_FC_IMFCS_ANS
}
