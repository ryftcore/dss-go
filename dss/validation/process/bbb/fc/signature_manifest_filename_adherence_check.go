// Ported from dss-validation/.../validation/process/bbb/fc/checks/SignatureManifestFilenameAdherenceCheck.java (DSS 6.5.RC1).
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

// SignatureManifestFilenameAdherenceCheck verifies conformance of the manifest filename related
// to a signature.
type SignatureManifestFilenameAdherenceCheck struct {
	FilenameAdherenceCheck[*diagnostic.SignatureWrapper]
}

// NewSignatureManifestFilenameAdherenceCheck is the default constructor.
func NewSignatureManifestFilenameAdherenceCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*drjaxb.XmlFC],
	diagnosticData *diagnostic.DiagnosticData, token *diagnostic.SignatureWrapper,
	constraint policy.LevelRule) *SignatureManifestFilenameAdherenceCheck {
	c := &SignatureManifestFilenameAdherenceCheck{}
	c.InitFilenameAdherenceCheck(i18nProvider, result, diagnosticData, token, constraint)
	c.InitChainItem(c)
	return c
}

// Process performs the check.
func (c *SignatureManifestFilenameAdherenceCheck) Process() bool {
	if c.DiagnosticData.ContainerType() == enumerations.ASiCContainerTypeASiCS {
		// 4.3.3.2 Contents of the container: the META-INF folder may contain other
		// application specific information - can be of any format.
		return true
	}
	manifestFile := c.DiagnosticData.ManifestFileForFilename(c.Token.Filename())
	signatureForm, err := c.Token.SignatureFormat().SignatureForm()
	if err != nil {
		panic(err)
	}
	if manifestFile == nil {
		// optional for XAdES, required for CAdES
		return signatureForm == enumerations.SignatureFormXAdES
	}

	manifestFilename := ""
	if manifestFile.Filename != nil {
		manifestFilename = *manifestFile.Filename
	}
	if strings.TrimSpace(manifestFilename) == "" {
		return false
	}
	switch signatureForm {
	case enumerations.SignatureFormXAdES:
		return AsiceMetainfManifest == manifestFilename
	case enumerations.SignatureFormCAdES:
		return c.IsASiCManifest(manifestFilename)
	default:
		panic(fmt.Sprintf("Only XAdES and CAdES ASiC container types are supported! Found : %s", signatureForm))
	}
}

// MessageTag returns the constraint message i18n key.
func (c *SignatureManifestFilenameAdherenceCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTagBBBFCIMFCS
}

// ErrorMessageTag returns the error message i18n key.
func (c *SignatureManifestFilenameAdherenceCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagBBBFCIMFCSANS
}
