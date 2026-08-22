// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/cv/CryptographicVerification.java (DSS 6.5.RC1).
package cv

import (
	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	diagnosticjaxb "github.com/ryftcore/dss-go/dss/diagnostic/jaxb"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/utils"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// digestMatcherType reads XmlDigestMatcher#getType(). The generated model
// carries the schema's optional members as pointers where the Java getters
// return null, so the checks of this package read them through these three
// nil-tolerant accessors rather than repeating the nil test at every use.
func digestMatcherType(digestMatcher *diagnosticjaxb.XmlDigestMatcher) enumerations.DigestMatcherType {
	if digestMatcher.Type == nil {
		return ""
	}
	return digestMatcher.Type.DigestMatcherType()
}

// digestMatcherId reads XmlDigestMatcher#getId().
func digestMatcherId(digestMatcher *diagnosticjaxb.XmlDigestMatcher) string {
	if digestMatcher.Id == nil {
		return ""
	}
	return *digestMatcher.Id
}

// digestMatcherUri reads XmlDigestMatcher#getUri().
func digestMatcherUri(digestMatcher *diagnosticjaxb.XmlDigestMatcher) string {
	if digestMatcher.Uri == nil {
		return ""
	}
	return *digestMatcher.Uri
}

// digestMatcherDocumentName reads XmlDigestMatcher#getDocumentName().
func digestMatcherDocumentName(digestMatcher *diagnosticjaxb.XmlDigestMatcher) string {
	if digestMatcher.DocumentName == nil {
		return ""
	}
	return *digestMatcher.DocumentName
}

// CryptographicVerification is 5.2.7 Cryptographic verification. This building
// block checks the integrity of the signed data by performing the cryptographic
// verifications.
type CryptographicVerification struct {
	*process.ChainBase[*jaxb.XmlCV]

	// diagnosticData is the Diagnostic data.
	diagnosticData *diagnostic.Data

	// token is the token to verify.
	token diagnostic.TokenProxy

	// validationPolicy is the validation policy.
	validationPolicy policy.ValidationPolicy

	// context is the validation context.
	context enumerations.Context
}

// NewCryptographicVerification is the default constructor. Port of
// CryptographicVerification(Provider, Data, TokenProxy, Context, ValidationPolicy).
func NewCryptographicVerification(i18nProvider *i18n.Provider, diagnosticData *diagnostic.Data,
	token diagnostic.TokenProxy, context enumerations.Context,
	validationPolicy policy.ValidationPolicy) *CryptographicVerification {
	xmlCV := &jaxb.XmlCV{}
	c := &CryptographicVerification{
		ChainBase: process.NewChainBase(i18nProvider, process.NewResult(xmlCV,
			&xmlCV.XmlConstraintsConclusionContent, &xmlCV.XmlConstraintsConclusionAttrs)),
		diagnosticData:   diagnosticData,
		token:            token,
		context:          context,
		validationPolicy: validationPolicy,
	}
	c.InitChainBase(c)
	return c
}

// Title returns the title of the building block. Port of getTitle().
func (c *CryptographicVerification) Title() i18n.MessageTag {
	return i18n.MessageTagCryptographicVerification
}

// InitChain initializes the chain. Port of initChain().
func (c *CryptographicVerification) InitChain() {

	var item process.ChainItem[*jaxb.XmlCV]

	/*
	 * 5.2.7.4 Processing The first and second steps as well as
	 * the Data To Be Signed depend on the signature type. The technical details
	 * on how to do this correctly are out of scope for the present document.
	 * See ETSI EN 319 122-1 [i.2], ETSI EN 319 122-2 [i.3], ETSI EN 319 132-1 [i.4],
	 * ETSI EN 319 132-2 [i.5], ETSI EN 319 142-1 [i.6], ETSI EN 319 142-2 [i.7]
	 * and IETF RFC 3852 [i.8] for details.
	 */

	digestMatchers := c.token.DigestMatchers()
	containsManifest := c.containsManifest(digestMatchers)

	if utils.IsCollectionNotEmpty(digestMatchers) {
		for _, digestMatcher := range digestMatchers {
			if enumerations.DigestMatcherTypeEvidenceRecordOrphanReference == digestMatcherType(digestMatcher) ||
				enumerations.DigestMatcherTypeEAAOrphanSelectivelyDisclosableClaim == digestMatcherType(digestMatcher) {
				// Evidence Records optionally allow additional digests to be present within first data group
				// EAAs allow non-disclosed hashes
				continue
			}
			if containsManifest && enumerations.DigestMatcherTypeManifestEntry == digestMatcherType(digestMatcher) {
				// move XML Manifest entries validation to a separate validation block
				continue
			}
			/*
			 * 1) The building block shall obtain the signed data object(s) if not provided
			 * in the inputs (e.g. by dereferencing an URI present in the signature). If the
			 * signed data object(s) cannot be obtained, the building block shall return the
			 * indication INDETERMINATE with the sub-indication SIGNED_DATA_NOT_FOUND.
			 */
			referenceDataFound := c.referenceDataFound(digestMatcher)
			if item == nil {
				item = referenceDataFound
				c.FirstItem = item
			} else {
				item = item.SetNextItem(referenceDataFound)
			}
			/*
			 * 2) The SVA shall check the integrity of the signed data objects. In case of
			 * failure, the building block shall return the indication FAILED with the
			 * sub-indication HASH_FAILURE.
			 */
			// to allow customizable validation of only identified entries
			if digestMatcher.DataFound {
				item = item.SetNextItem(c.referenceDataIntact(digestMatcher))
			}

			// perform validation when only both URI and document name are present
			if utils.IsStringNotEmpty(digestMatcherUri(digestMatcher)) &&
				utils.IsStringNotEmpty(digestMatcherDocumentName(digestMatcher)) {
				item = item.SetNextItem(c.referenceDataNameCheck(digestMatcher))
			}
		}

		if c.isEvidenceRecordHashTreeRenewalTimestamp() {

			evidenceRecordHashTreeRenewalTimestamp := c.evidenceRecordHashTreeRenewalTimestamp()
			if item == nil {
				item = evidenceRecordHashTreeRenewalTimestamp
				c.FirstItem = item
			} else {
				item = item.SetNextItem(evidenceRecordHashTreeRenewalTimestamp)
			}

		}
	}

	if containsManifest {

		manifestEntryFound := c.manifestEntryExistence(digestMatchers)
		if item == nil {
			item = manifestEntryFound
			c.FirstItem = item
		} else {
			item = item.SetNextItem(manifestEntryFound)
		}

		// manifest entries may be omitted when no detached data is provided to the validation
		if c.containsManifestEntries(digestMatchers) {

			item = item.SetNextItem(c.manifestEntryGroup(digestMatchers))

			for _, digestMatcher := range digestMatchers {
				if enumerations.DigestMatcherTypeManifestEntry == digestMatcherType(digestMatcher) {

					if digestMatcher.DataFound {
						item = item.SetNextItem(c.manifestEntryIntact(digestMatcher))
					}

					if utils.IsStringNotEmpty(digestMatcherUri(digestMatcher)) &&
						utils.IsStringNotEmpty(digestMatcherDocumentName(digestMatcher)) {
						item = item.SetNextItem(c.manifestEntryNameCheck(digestMatcher))
					}

				}
			}

		}

	}

	/*
	 * 3) The building block shall verify the cryptographic signature using the public key extracted from the
	 * signing certificate in the chain, the signature value and the signature algorithm extracted from the
	 * signature. If this cryptographic verification outputs a success indication, the building block shall return
	 * the indication PASSED.
	 *
	 * 4) Otherwise, the building block shall return the indication FAILED and the sub-indication
	 * SIG_CRYPTO_FAILURE.
	 */
	signatureIntact := c.signatureIntact()
	if item == nil {
		item = signatureIntact
		c.FirstItem = item
	} else {
		item = item.SetNextItem(signatureIntact) //nolint:staticcheck // mirrors upstream CryptographicVerification#initChain: Java's trailing `item = item.setNextItem(...)` is the same dead store - setNextItem links the item and returns it, and nothing reads the tail afterwards.
	}

}

// containsManifest ports the private containsManifest(List).
func (c *CryptographicVerification) containsManifest(digestMatchers []*diagnosticjaxb.XmlDigestMatcher) bool {
	for _, d := range digestMatchers {
		if enumerations.DigestMatcherTypeManifest == digestMatcherType(d) {
			return true
		}
	}
	return false
}

// containsManifestEntries ports the private containsManifestEntries(List).
func (c *CryptographicVerification) containsManifestEntries(digestMatchers []*diagnosticjaxb.XmlDigestMatcher) bool {
	for _, d := range digestMatchers {
		if enumerations.DigestMatcherTypeManifestEntry == digestMatcherType(d) && d.DataFound {
			return true
		}
	}
	return false
}

// referenceDataFound ports the private referenceDataFound(XmlDigestMatcher).
func (c *CryptographicVerification) referenceDataFound(digestMatcher *diagnosticjaxb.XmlDigestMatcher) process.ChainItem[*jaxb.XmlCV] {
	var constraint policy.LevelRule
	if enumerations.ContextEAA == c.context {
		constraint = c.validationPolicy.EAADisclosureFoundConstraint()
	} else {
		constraint = c.validationPolicy.ReferenceDataExistenceConstraint(c.context)
	}
	return NewReferenceDataExistenceCheck(c.I18nProvider, c.Result, digestMatcher, constraint)
}

// referenceDataIntact ports the private referenceDataIntact(XmlDigestMatcher).
func (c *CryptographicVerification) referenceDataIntact(digestMatcher *diagnosticjaxb.XmlDigestMatcher) process.ChainItem[*jaxb.XmlCV] {
	var constraint policy.LevelRule
	if enumerations.ContextEAA == c.context {
		constraint = c.validationPolicy.EAADisclosureIntactConstraint()
	} else {
		constraint = c.validationPolicy.ReferenceDataIntactConstraint(c.context)
	}
	return NewReferenceDataIntactCheck(c.I18nProvider, c.Result, digestMatcher, constraint)
}

// referenceDataNameCheck ports the private referenceDataNameCheck(XmlDigestMatcher).
func (c *CryptographicVerification) referenceDataNameCheck(digestMatcher *diagnosticjaxb.XmlDigestMatcher) process.ChainItem[*jaxb.XmlCV] {
	constraint := c.validationPolicy.ReferenceDataNameMatchConstraint(c.context)
	return NewReferenceDataNameMatchCheck(c.I18nProvider, c.Result, digestMatcher, constraint)
}

// manifestEntryExistence ports the private manifestEntryExistence(List).
func (c *CryptographicVerification) manifestEntryExistence(digestMatchers []*diagnosticjaxb.XmlDigestMatcher) process.ChainItem[*jaxb.XmlCV] {
	constraint := c.validationPolicy.ManifestEntryObjectExistenceConstraint(c.context)
	return NewManifestEntryExistenceCheck(c.I18nProvider, c.Result, digestMatchers, constraint)
}

// manifestEntryGroup ports the private manifestEntryGroup(List).
func (c *CryptographicVerification) manifestEntryGroup(digestMatchers []*diagnosticjaxb.XmlDigestMatcher) process.ChainItem[*jaxb.XmlCV] {
	constraint := c.validationPolicy.ManifestEntryObjectGroupConstraint(c.context)
	return NewManifestEntryGroupCheck(c.I18nProvider, c.Result, digestMatchers, constraint)
}

// manifestEntryIntact ports the private manifestEntryIntact(XmlDigestMatcher).
func (c *CryptographicVerification) manifestEntryIntact(digestMatcher *diagnosticjaxb.XmlDigestMatcher) process.ChainItem[*jaxb.XmlCV] {
	constraint := c.validationPolicy.ManifestEntryObjectIntactConstraint(c.context)
	return NewReferenceDataIntactCheck(c.I18nProvider, c.Result, digestMatcher, constraint)
}

// manifestEntryNameCheck ports the private manifestEntryNameCheck(XmlDigestMatcher).
func (c *CryptographicVerification) manifestEntryNameCheck(digestMatcher *diagnosticjaxb.XmlDigestMatcher) process.ChainItem[*jaxb.XmlCV] {
	constraint := c.validationPolicy.ManifestEntryNameMatchConstraint(c.context)
	return NewReferenceDataNameMatchCheck(c.I18nProvider, c.Result, digestMatcher, constraint)
}

// signatureIntact ports the private signatureIntact().
func (c *CryptographicVerification) signatureIntact() process.ChainItem[*jaxb.XmlCV] {
	constraint := c.validationPolicy.SignatureIntactConstraint(c.context)
	return NewSignatureIntactCheck(c.I18nProvider, c.Result, c.token, c.context, constraint)
}

// isEvidenceRecordHashTreeRenewalTimestamp ports the private
// isEvidenceRecordHashTreeRenewalTimestamp().
func (c *CryptographicVerification) isEvidenceRecordHashTreeRenewalTimestamp() bool {
	if timestampWrapper, ok := c.token.(*diagnostic.TimestampWrapper); ok {
		return timestampWrapper.Type().IsEvidenceRecordTimestamp() &&
			enumerations.EvidenceRecordTimestampTypeHashTreeRenewalArchiveTimestamp ==
				timestampWrapper.EvidenceRecordTimestampType()
	}
	return false
}

// evidenceRecordHashTreeRenewalTimestamp ports the private
// evidenceRecordHashTreeRenewalTimestamp().
func (c *CryptographicVerification) evidenceRecordHashTreeRenewalTimestamp() process.ChainItem[*jaxb.XmlCV] {
	constraint := c.validationPolicy.EvidenceRecordHashTreeRenewalConstraint()
	return NewEvidenceRecordHashTreeRenewalTimestampCheck(c.I18nProvider, c.Result, c.diagnosticData,
		c.token.(*diagnostic.TimestampWrapper), constraint)
}
