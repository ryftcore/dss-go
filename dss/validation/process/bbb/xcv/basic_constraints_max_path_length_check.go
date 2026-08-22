// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/xcv/sub/checks/BasicConstraintsMaxPathLengthCheck.java (DSS 6.5.RC1).
package xcv

import (
	jaxb "github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// BasicConstraintsMaxPathLengthCheck verifies whether the certificate path
// depth of the current certificate is conformant with
// BasicConstraints.pathLenConstraint value defined within intermediate CA
// certificates precessing in the chain.
type BasicConstraintsMaxPathLengthCheck struct {
	*process.ChainItemBase[*jaxb.XmlSubXCV]

	// certificate is the certificate to check.
	certificate *diagnostic.CertificateWrapper
}

// NewBasicConstraintsMaxPathLengthCheck is the default constructor. Port of
// BasicConstraintsMaxPathLengthCheck(Provider, XmlSubXCV, CertificateWrapper, LevelRule).
func NewBasicConstraintsMaxPathLengthCheck(i18nProvider *i18n.Provider, result *process.Result[*jaxb.XmlSubXCV],
	certificate *diagnostic.CertificateWrapper, constraint policy.LevelRule) *BasicConstraintsMaxPathLengthCheck {
	c := &BasicConstraintsMaxPathLengthCheck{
		ChainItemBase: process.NewChainItemBase(i18nProvider, result, constraint),
		certificate:   certificate,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *BasicConstraintsMaxPathLengthCheck) Process() bool {
	certificateChain := c.certificate.CertificateChain()
	/*
	 * (k) max_path_length: this integer is initialized to n, is
	 * decremented for each non-self-issued certificate in the path,
	 * and may be reduced to the value in the path length constraint
	 * field within the basic constraints extension of a CA
	 * certificate.
	 */
	maxPathLength := len(certificateChain) + 1 // certificate chain does not return current certificate
	for i := len(certificateChain) - 1; i > -1; i-- {
		/*
		 * (l) If the certificate was not self-issued, verify that
		 * max_path_length is greater than zero and decrement
		 * max_path_length by 1.
		 */
		cert := certificateChain[i]
		if !cert.IsSelfSigned() {
			maxPathLength--
		}
		pathLenConstraint := cert.PathLenConstraint()
		/*
		 * (m) If pathLenConstraint is present in the certificate and is
		 * less than max_path_length, set max_path_length to the value
		 * of pathLenConstraint.
		 */
		if pathLenConstraint != -1 && pathLenConstraint < maxPathLength {
			maxPathLength = pathLenConstraint
		}
	}
	return maxPathLength > 0
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *BasicConstraintsMaxPathLengthCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTagBBBXCVICPDV
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *BasicConstraintsMaxPathLengthCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagBBBXCVICPDVANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *BasicConstraintsMaxPathLengthCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationIndeterminate
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *BasicConstraintsMaxPathLengthCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationCertificateChainGeneralFailure
}
