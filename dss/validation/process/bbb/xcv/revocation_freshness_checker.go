// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/xcv/rfc/RevocationFreshnessChecker.java (DSS 6.5.RC1).
package xcv

import (
	"fmt"
	"time"

	jaxb "github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// RevocationFreshnessChecker is 5.2.5 Revocation freshness checker. This
// building block checks that a given revocation status information is
// "fresh" at a given validation time. The freshness of the revocation status
// information is the maximum accepted difference between the issuance time
// of the revocation status information and the current time. This process is
// used by other validation blocks when checking the revocation status of a
// certificate.
type RevocationFreshnessChecker struct {
	*process.ChainBase[*jaxb.XmlRFC]

	// revocationData is the revocation data to check.
	revocationData *diagnostic.RevocationWrapper

	// validationDate is the validation time.
	validationDate time.Time

	// policy is the validation policy.
	policy policy.ValidationPolicy

	// context is the validation context.
	context enumerations.Context

	// subContext is the validation subContext.
	subContext enumerations.SubContext
}

// NewRevocationFreshnessChecker is the default constructor. Port of
// RevocationFreshnessChecker(I18nProvider, RevocationWrapper, Date, Context, SubContext, ValidationPolicy).
func NewRevocationFreshnessChecker(i18nProvider *i18n.I18nProvider, revocationData *diagnostic.RevocationWrapper,
	validationDate time.Time, context enumerations.Context, subContext enumerations.SubContext,
	validationPolicy policy.ValidationPolicy) *RevocationFreshnessChecker {
	xmlRFC := &jaxb.XmlRFC{}
	c := &RevocationFreshnessChecker{
		ChainBase: process.NewChainBase(i18nProvider, process.NewResult(xmlRFC,
			&xmlRFC.XmlConstraintsConclusionContent, &xmlRFC.XmlConstraintsConclusionAttrs)),
		revocationData: revocationData,
		validationDate: validationDate,
		policy:         validationPolicy,
		context:        context,
		subContext:     subContext,
	}
	if revocationData != nil {
		id := revocationData.Id()
		c.Result.Value.Id = &id
	}
	c.InitChainBase(c)
	return c
}

// Title returns the title of the building block. Port of getTitle().
func (c *RevocationFreshnessChecker) Title() i18n.MessageTag {
	return i18n.MessageTagRevocationFreshnessChecker
}

// InitChain initializes the chain. Port of initChain().
func (c *RevocationFreshnessChecker) InitChain() {
	var item process.ChainItem[*jaxb.XmlRFC]

	if c.revocationData != nil {
		id := c.revocationData.Id()
		c.Result.Value.Id = &id

		/*
		 * 1) The building block shall get the maximum accepted revocation
		 * freshness from the X.509 validation constraints for the given
		 * certificate. If the constraints do not contain a value for the
		 * maximum accepted revocation freshness and the revocation
		 * information status is a CRL or an OCSP response IETF RFC 5280
		 * [1], IETF RFC 6960 [i.12] with a value in the nextUpdate field
		 * the time interval between the fields thisUpdate and nextUpdate
		 * shall be used as the value of maximum freshness. If nextUpdate is
		 * not set, the building block shall return with the indication
		 * FAILED.
		 *
		 * NOTE: This means that if the given validation time is after the
		 * nextUpdate time, the revocation status information will not be
		 * considered fresh.
		 */
		revocationFreshnessConstraint := c.policy.RevocationFreshnessConstraint(c.context, c.subContext)
		if revocationFreshnessConstraint == nil || enumerations.LevelIgnore == revocationFreshnessConstraint.Level() {
			switch c.revocationData.RevocationType() {
			case enumerations.RevocationTypeCRL:
				item = c.crlNextUpdateCheck(c.revocationData)
				c.FirstItem = item
			case enumerations.RevocationTypeOCSP:
				item = c.ocspNextUpdateCheck(c.revocationData)
				c.FirstItem = item
			default:
				panic(fmt.Errorf("the RevocationType '%s' is not supported!", c.revocationData.RevocationType()))
			}
		}
		/*
		 * 2) If the issuance time of the revocation information status is
		 * after the validation time minus the considered maximum freshness,
		 * the building block shall return the indication PASSED. Otherwise
		 * the building block shall return the indication FAILED.
		 */
		if item == nil {
			item = c.revocationDataFreshCheck(c.revocationData, revocationFreshnessConstraint)
			c.FirstItem = item
		} else {
			item = item.SetNextItem(c.revocationDataFreshCheck(c.revocationData, revocationFreshnessConstraint)) //nolint:staticcheck // mirrors upstream RevocationFreshnessChecker#initChain: Java's trailing `item = item.setNextItem(...)` is the same dead store - setNextItem links the item and returns it, and nothing reads the tail afterwards.
		}
	}
}

// crlNextUpdateCheck ports the private crlNextUpdateCheck(RevocationWrapper).
func (c *RevocationFreshnessChecker) crlNextUpdateCheck(revocationData *diagnostic.RevocationWrapper) process.ChainItem[*jaxb.XmlRFC] {
	constraint := c.policy.CRLNextUpdatePresentConstraint(c.context, c.subContext)
	return NewNextUpdateCheck(c.I18nProvider, c.Result, revocationData, constraint)
}

// ocspNextUpdateCheck ports the private ocspNextUpdateCheck(RevocationWrapper).
func (c *RevocationFreshnessChecker) ocspNextUpdateCheck(revocationData *diagnostic.RevocationWrapper) process.ChainItem[*jaxb.XmlRFC] {
	constraint := c.policy.OCSPNextUpdatePresentConstraint(c.context, c.subContext)
	return NewNextUpdateCheck(c.I18nProvider, c.Result, revocationData, constraint)
}

// revocationDataFreshCheck ports the private
// revocationDataFreshCheck(RevocationWrapper, DurationRule).
func (c *RevocationFreshnessChecker) revocationDataFreshCheck(revocationData *diagnostic.RevocationWrapper,
	revocationFreshnessConstraint policy.DurationRule) process.ChainItem[*jaxb.XmlRFC] {
	/*
	 * The building block shall get the maximum accepted revocation
	 * freshness from the X.509 validation constraints for the given
	 * certificate.
	 */
	if revocationFreshnessConstraint != nil {
		return NewRevocationDataFreshCheck(c.I18nProvider, c.Result, revocationData, c.validationDate, revocationFreshnessConstraint)
	}
	/*
	 * If the constraints do not contain a value for the maximum accepted
	 * revocation freshness and the revocation information status is a CRL
	 * or an OCSP response IETF RFC 5280 [1], IETF RFC 6960 [i.12] with a
	 * value in the nextUpdate field the time interval between the fields
	 * thisUpdate and nextUpdate shall be used as the value of maximum
	 * freshness.
	 */
	constraint := c.policy.RevocationFreshnessNextUpdateConstraint(c.context, c.subContext)
	return NewRevocationDataFreshCheckWithNullConstraint(c.I18nProvider, c.Result, revocationData, c.validationDate, constraint)
}
