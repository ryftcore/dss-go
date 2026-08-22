// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/vpfswatsp/checks/psv/checks/POENotAfterCARevocationTimeCheck.java (DSS 6.5.RC1).
//
// See poe.go for the package-flattening note.
//
// Java's type parameter <R extends RevocationWrapper> only widens the accepted
// collection; the body reads nothing but RevocationWrapper#getId(). The Go
// class is therefore not generic and takes the CertificateRevocationWrapper
// list its only caller (PastSignatureValidation) passes.
package vpfswatsp

import (
	"time"

	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// POENotAfterCARevocationTimeCheck verifies if there is a POE for the
// revocation information of the signer certificate at (or before) the
// revocation time of the CA certificate.
type POENotAfterCARevocationTimeCheck struct {
	*process.ChainItemBase[*jaxb.XmlPSV]

	// revocationData is a collection of filtered acceptable revocation data.
	revocationData []*diagnostic.CertificateRevocationWrapper

	// caRevocationTime is the revocation time of the CA certificate; nil is
	// Java's null, which its Date#compareTo would dereference.
	caRevocationTime *time.Time

	// poeExtraction is a collection of POEs.
	poeExtraction *POEExtraction
}

// NewPOENotAfterCARevocationTimeCheck is the default constructor. Port of
// POENotAfterCARevocationTimeCheck(Provider, XmlPSV, Collection, Date, POEExtraction, LevelRule).
func NewPOENotAfterCARevocationTimeCheck(i18nProvider *i18n.Provider, result *process.Result[*jaxb.XmlPSV],
	revocationData []*diagnostic.CertificateRevocationWrapper, caRevocationTime *time.Time,
	poeExtraction *POEExtraction, constraint policy.LevelRule) *POENotAfterCARevocationTimeCheck {
	c := &POENotAfterCARevocationTimeCheck{
		ChainItemBase:    process.NewChainItemBase(i18nProvider, result, constraint),
		revocationData:   revocationData,
		caRevocationTime: caRevocationTime,
		poeExtraction:    poeExtraction,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *POENotAfterCARevocationTimeCheck) Process() bool {
	for _, revocationWrapper := range c.revocationData {
		if c.poeExtraction.IsPOEExists(revocationWrapper.Id(), *c.caRevocationTime) {
			return true
		}
	}
	return false
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *POENotAfterCARevocationTimeCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTagPSVITPRISCNARTCAC
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *POENotAfterCARevocationTimeCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagPSVITPRISCNARTCACANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *POENotAfterCARevocationTimeCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationIndeterminate
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *POENotAfterCARevocationTimeCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationRevokedCANoPOE
}
