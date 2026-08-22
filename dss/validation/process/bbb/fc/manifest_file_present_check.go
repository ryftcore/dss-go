// Ported from dss-validation/.../validation/process/bbb/fc/checks/ManifestFilePresentCheck.java (DSS 6.5.RC1).
package fc

import (
	drjaxb "github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	diagjaxb "github.com/ryftcore/dss-go/dss/diagnostic/jaxb"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	policy "github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// ManifestFilePresentCheck checks if the manifest file is present inside an ASiC container.
type ManifestFilePresentCheck struct {
	*process.ChainItemBase[*drjaxb.XmlFC]

	containerInfo *diagjaxb.XmlContainerInfo
}

// NewManifestFilePresentCheck is the default constructor.
func NewManifestFilePresentCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*drjaxb.XmlFC],
	containerInfo *diagjaxb.XmlContainerInfo, constraint policy.LevelRule) *ManifestFilePresentCheck {
	c := &ManifestFilePresentCheck{containerInfo: containerInfo}
	c.ChainItemBase = process.NewChainItemBase(i18nProvider, result, constraint)
	c.InitChainItem(c)
	return c
}

// Process performs the check.
func (c *ManifestFilePresentCheck) Process() bool {
	if c.containerInfo.ContainerType != nil &&
		c.containerInfo.ContainerType.ASiCContainerType() == enumerations.ASiCContainerTypeASiCE {
		return len(c.containerInfo.ManifestFiles.All()) > 0
	}
	// ASiC-S container may contain a manifest file
	return true
}

// MessageTag returns the constraint message i18n key.
func (c *ManifestFilePresentCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_FC_IMFP_ASICE
}

// ErrorMessageTag returns the error message i18n key.
func (c *ManifestFilePresentCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_FC_IMFP_ASICE_ANS
}

// FailedIndicationForConclusion returns the Indication on failure.
func (c *ManifestFilePresentCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationFailed
}

// FailedSubIndicationForConclusion returns the SubIndication on failure.
func (c *ManifestFilePresentCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationFormatFailure
}
