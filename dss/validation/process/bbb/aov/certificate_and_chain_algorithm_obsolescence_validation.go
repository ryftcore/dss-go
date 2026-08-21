// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/aov/CertificateAndChainAlgorithmObsolescenceValidation.java (DSS 6.5.RC1).
//
// Performs cryptographic validation of the certificate and its chain.
package aov

import (
	"time"

	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/utils"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// CertificateAndChainAlgorithmObsolescenceValidation performs cryptographic
// validation of the certificate and its chain.
type CertificateAndChainAlgorithmObsolescenceValidation struct {
	TokenAlgorithmObsolescenceValidation[*diagnostic.CertificateWrapper]
}

// NewCertificateAndChainAlgorithmObsolescenceValidation is the default
// constructor.
func NewCertificateAndChainAlgorithmObsolescenceValidation(i18nProvider *i18n.I18nProvider, token *diagnostic.CertificateWrapper,
	context enumerations.Context, validationDate time.Time,
	validationPolicy policy.ValidationPolicy) *CertificateAndChainAlgorithmObsolescenceValidation {
	c := &CertificateAndChainAlgorithmObsolescenceValidation{}
	c.InitAlgorithmObsolescenceValidation(i18nProvider, token, context, validationDate, validationPolicy, c)
	c.InitChainBase(c)
	return c
}

// BuildChain builds a chain of checks to be executed during the process. Port
// of the overridden protected ChainItem<XmlAOV> buildChain().
func (c *CertificateAndChainAlgorithmObsolescenceValidation) BuildChain() process.ChainItem[*jaxb.XmlAOV] {
	return c.buildCertificateChainValidationChain(c.FirstItem, c.token, c.getFullCertificateChain())
}

// getFullCertificateChain ports the private getFullCertificateChain().
func (c *CertificateAndChainAlgorithmObsolescenceValidation) getFullCertificateChain() []*diagnostic.CertificateWrapper {
	var certificateChain []*diagnostic.CertificateWrapper
	if c.token != nil {
		certificateChain = append(certificateChain, c.token)
		if utils.IsCollectionNotEmpty(c.token.CertificateChain()) {
			certificateChain = append(certificateChain, c.token.CertificateChain()...)
		}
	}
	return certificateChain
}
