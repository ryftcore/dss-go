// Ported from dss-validation/.../validation/process/bbb/fc/checks/SignedFilesPresentCheck.java (DSS 6.5.RC1).
package fc

import (
	"strings"

	drjaxb "github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	diagjaxb "github.com/ryftcore/dss-go/dss/diagnostic/jaxb"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	policy "github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// SignedFilesPresentCheck checks if signed files are present in an ASiC container.
type SignedFilesPresentCheck struct {
	*process.ChainItemBase[*drjaxb.XmlFC]

	containerInfo *diagjaxb.XmlContainerInfo
}

// NewSignedFilesPresentCheck is the default constructor.
func NewSignedFilesPresentCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*drjaxb.XmlFC],
	containerInfo *diagjaxb.XmlContainerInfo, constraint policy.LevelRule) *SignedFilesPresentCheck {
	c := &SignedFilesPresentCheck{containerInfo: containerInfo}
	c.ChainItemBase = process.NewChainItemBase(i18nProvider, result, constraint)
	c.InitChainItem(c)
	return c
}

func (c *SignedFilesPresentCheck) isASiCS() bool {
	return c.containerInfo.ContainerType != nil &&
		c.containerInfo.ContainerType.ASiCContainerType() == enumerations.ASiCContainerTypeASiCS
}

func isRootDirectoryFile(fileName string) bool {
	return !strings.Contains(fileName, "/") && !strings.Contains(fileName, "\\")
}

func rootLevelFiles(fileNames []string) []string {
	var result []string
	for _, fileName := range fileNames {
		if isRootDirectoryFile(fileName) {
			result = append(result, fileName)
		}
	}
	return result
}

// Process performs the check.
func (c *SignedFilesPresentCheck) Process() bool {
	contentFiles := c.containerInfo.ContentFiles.All()
	if c.isASiCS() {
		// ASiC-S one signed file in the root directory
		root := rootLevelFiles(contentFiles)
		return len(root) == 1
	}
	// ASiC-E one or more signed files outside META-INF
	return len(contentFiles) > 0
}

// MessageTag returns the constraint message i18n key.
func (c *SignedFilesPresentCheck) MessageTag() i18n.MessageTag {
	if c.isASiCS() {
		return i18n.MessageTag_BBB_FC_ISFP_ASICS
	}
	return i18n.MessageTag_BBB_FC_ISFP_ASICE
}

// ErrorMessageTag returns the error message i18n key.
func (c *SignedFilesPresentCheck) ErrorMessageTag() i18n.MessageTag {
	if c.isASiCS() {
		return i18n.MessageTag_BBB_FC_ISFP_ASICS_ANS
	}
	return i18n.MessageTag_BBB_FC_ISFP_ASICE_ANS
}

// FailedIndicationForConclusion returns the Indication on failure.
func (c *SignedFilesPresentCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationFailed
}

// FailedSubIndicationForConclusion returns the SubIndication on failure.
func (c *SignedFilesPresentCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationFormatFailure
}
