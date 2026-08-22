// Ported from dss-validation/.../validation/process/bbb/fc/checks/TimestampFilenameAdherenceCheck.java (DSS 6.5.RC1).
package fc

import (
	"fmt"
	"strings"

	drjaxb "github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	policy "github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

const (
	tstExtension = ".tst"
	timestampTST = MetaInfFolder + TimestampFilename + tstExtension
)

// TimestampFilenameAdherenceCheck verifies conformance of the timestamp's document filename to
// the ASiC specification.
type TimestampFilenameAdherenceCheck struct {
	FilenameAdherenceCheck[*diagnostic.TimestampWrapper]
}

// NewTimestampFilenameAdherenceCheck is the default constructor.
func NewTimestampFilenameAdherenceCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*drjaxb.XmlFC],
	diagnosticData *diagnostic.DiagnosticData, token *diagnostic.TimestampWrapper,
	constraint policy.LevelRule) *TimestampFilenameAdherenceCheck {
	c := &TimestampFilenameAdherenceCheck{}
	c.InitFilenameAdherenceCheck(i18nProvider, result, diagnosticData, token, constraint)
	c.InitChainItem(c)
	return c
}

func (c *TimestampFilenameAdherenceCheck) isTimestamp(filename string) bool {
	return strings.HasPrefix(filename, MetaInfFolder) && strings.Contains(filename, TimestampFilename) &&
		strings.HasSuffix(filename, tstExtension)
}

func (c *TimestampFilenameAdherenceCheck) isInitialTimestampToken(filename string) bool {
	return timestampTST == filename
}

func (c *TimestampFilenameAdherenceCheck) isArchiveTimestampToken(filename string) bool {
	manifestFile := c.DiagnosticData.ManifestFileForFilename(filename)
	if manifestFile != nil && manifestFile.Filename != nil && c.IsASiCArchiveManifest(*manifestFile.Filename) {
		return c.isTimestamp(filename)
	}
	return false
}

// Process performs the check.
func (c *TimestampFilenameAdherenceCheck) Process() bool {
	filename := c.Token.Filename()
	if strings.TrimSpace(filename) == "" {
		return false
	}
	switch c.DiagnosticData.ContainerType() {
	case enumerations.ASiCContainerType_ASiC_S:
		return c.isInitialTimestampToken(filename) || c.isArchiveTimestampToken(filename)
	case enumerations.ASiCContainerType_ASiC_E:
		return c.isTimestamp(filename)
	default:
		panic(fmt.Sprintf("Container type '%s' is not supported!", c.DiagnosticData.ContainerType()))
	}
}

// MessageTag returns the constraint message i18n key.
func (c *TimestampFilenameAdherenceCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_FC_ISFCS
}

// ErrorMessageTag returns the error message i18n key.
func (c *TimestampFilenameAdherenceCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_FC_ISFCS_ANS
}
