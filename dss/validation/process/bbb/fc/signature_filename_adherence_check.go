// Ported from dss-validation/.../validation/process/bbb/fc/checks/SignatureFilenameAdherenceCheck.java (DSS 6.5.RC1).
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
	signaturesFilename      = "signatures"
	cadesSignatureExtension = ".p7s"
	signaturesXML           = MetaInfFolder + signaturesFilename + XMLExtension
	signatureP7S            = MetaInfFolder + SignatureFilename + cadesSignatureExtension
)

// SignatureFilenameAdherenceCheck checks validity of the signature's filename against the
// ASiC specification.
type SignatureFilenameAdherenceCheck struct {
	FilenameAdherenceCheck[*diagnostic.SignatureWrapper]
}

// NewSignatureFilenameAdherenceCheck is the default constructor.
func NewSignatureFilenameAdherenceCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*drjaxb.XmlFC],
	diagnosticData *diagnostic.DiagnosticData, token *diagnostic.SignatureWrapper,
	constraint policy.LevelRule) *SignatureFilenameAdherenceCheck {
	c := &SignatureFilenameAdherenceCheck{}
	c.InitFilenameAdherenceCheck(i18nProvider, result, diagnosticData, token, constraint)
	c.InitChainItem(c)
	return c
}

// Process performs the check.
func (c *SignatureFilenameAdherenceCheck) Process() bool {
	filename := c.Token.Filename()
	if strings.TrimSpace(filename) == "" {
		return false
	}
	signatureForm, err := c.Token.SignatureFormat().SignatureForm()
	if err != nil {
		panic(err)
	}
	switch c.DiagnosticData.ContainerType() {
	case enumerations.ASiCContainerType_ASiC_S:
		switch signatureForm {
		case enumerations.SignatureForm_XAdES:
			return signaturesXML == filename
		case enumerations.SignatureForm_CAdES:
			return signatureP7S == filename
		default:
			panic(fmt.Sprintf("Only XAdES and CAdES ASiC container types are supported! Found : %s", signatureForm))
		}
	case enumerations.ASiCContainerType_ASiC_E:
		switch signatureForm {
		case enumerations.SignatureForm_XAdES:
			return strings.HasPrefix(filename, MetaInfFolder) && strings.Contains(filename, signaturesFilename) &&
				strings.HasSuffix(filename, XMLExtension)
		case enumerations.SignatureForm_CAdES:
			return strings.HasPrefix(filename, MetaInfFolder) && strings.Contains(filename, SignatureFilename) &&
				strings.HasSuffix(filename, cadesSignatureExtension)
		default:
			panic(fmt.Sprintf("Only XAdES and CAdES ASiC container types are supported! Found : %s", signatureForm))
		}
	default:
		panic(fmt.Sprintf("Container type '%s' is not supported!", c.DiagnosticData.ContainerType()))
	}
}

// MessageTag returns the constraint message i18n key.
func (c *SignatureFilenameAdherenceCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_FC_ISFCS
}

// ErrorMessageTag returns the error message i18n key.
func (c *SignatureFilenameAdherenceCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_FC_ISFCS_ANS
}
