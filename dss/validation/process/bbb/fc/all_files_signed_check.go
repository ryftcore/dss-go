// Ported from dss-validation/.../validation/process/bbb/fc/checks/AllFilesSignedCheck.java (DSS 6.5.RC1).
package fc

import (
	drjaxb "github.com/utain/esig/dss/detailedreport/jaxb"
	"github.com/utain/esig/dss/diagnostic"
	diagjaxb "github.com/utain/esig/dss/diagnostic/jaxb"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/i18n"
	policy "github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/validation/process"
)

// AllFilesSignedCheck checks if all files are signed inside an ASiC container.
type AllFilesSignedCheck struct {
	*process.ChainItemBase[*drjaxb.XmlFC]

	signature     *diagnostic.SignatureWrapper
	containerInfo *diagjaxb.XmlContainerInfo
}

// NewAllFilesSignedCheck is the default constructor.
func NewAllFilesSignedCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*drjaxb.XmlFC],
	signature *diagnostic.SignatureWrapper, containerInfo *diagjaxb.XmlContainerInfo,
	constraint policy.LevelRule) *AllFilesSignedCheck {
	c := &AllFilesSignedCheck{signature: signature, containerInfo: containerInfo}
	c.ChainItemBase = process.NewChainItemBase(i18nProvider, result, constraint)
	c.InitChainItem(c)
	return c
}

func coversAllOriginalFiles(coveredFiles, originalFiles []string) bool {
	for _, file := range originalFiles {
		if !containsString(coveredFiles, file) {
			return false
		}
	}
	return true
}

func (c *AllFilesSignedCheck) relatedManifestFile(signatureFilename string) *diagjaxb.XmlManifestFile {
	for _, manifestFile := range c.containerInfo.ManifestFiles.All() {
		sigFilename := ""
		if manifestFile.SignatureFilename != nil {
			sigFilename = *manifestFile.SignatureFilename
		}
		if sigFilename == signatureFilename {
			return manifestFile
		}
	}
	return nil
}

func (c *AllFilesSignedCheck) coveredFilesFromScope() []string {
	var result []string
	for _, sigScope := range c.signature.SignatureScopes() {
		if sigScope.Scope != nil && sigScope.Scope.SignatureScopeType() == enumerations.SignatureScopeType_FULL {
			name := ""
			if sigScope.Name != nil {
				name = *sigScope.Name
			}
			result = append(result, name)
		}
	}
	return result
}

// Process performs the check.
func (c *AllFilesSignedCheck) Process() bool {
	containerType := enumerations.ASiCContainerType("")
	if c.containerInfo.ContainerType != nil {
		containerType = c.containerInfo.ContainerType.ASiCContainerType()
	}

	// ASiC-S -> nb files = 1
	if containerType == enumerations.ASiCContainerType_ASiC_S {
		return len(c.containerInfo.ContentFiles.All()) == 1
	} else if containerType == enumerations.ASiCContainerType_ASiC_E {
		signatureFilename := c.signature.Filename()
		contentFiles := c.containerInfo.ContentFiles.All()

		signatureForm, _ := c.signature.SignatureFormat().SignatureForm()

		manifestFile := c.relatedManifestFile(signatureFilename)
		if manifestFile != nil {
			coveredFiles := manifestFile.Entries.All()
			// check manifest <> content
			if !coversAllOriginalFiles(coveredFiles, contentFiles) {
				return false
			}
		} else if signatureForm == enumerations.SignatureForm_CAdES {
			// CAdES -> manifest file shall be present and signed
			return false
		}

		// XAdES -> check signature scope
		if signatureForm == enumerations.SignatureForm_XAdES {
			return coversAllOriginalFiles(c.coveredFilesFromScope(), contentFiles)
		}

		return true
	}

	return false
}

// MessageTag returns the constraint message i18n key.
func (c *AllFilesSignedCheck) MessageTag() i18n.MessageTag { return i18n.MessageTag_BBB_CV_IAFS }

// ErrorMessageTag returns the error message i18n key.
func (c *AllFilesSignedCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_CV_IAFS_ANS
}

// FailedIndicationForConclusion returns the Indication on failure.
func (c *AllFilesSignedCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.Indication_FAILED
}

// FailedSubIndicationForConclusion returns the SubIndication on failure.
func (c *AllFilesSignedCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndication_FORMAT_FAILURE
}
