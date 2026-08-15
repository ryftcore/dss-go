// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/xcv/rac/checks/RevocationAcceptanceCheckerResultCheck.java (DSS 6.5.RC1).
package xcv

import (
	"time"

	"github.com/utain/esig/dss/detailedreport/jaxb"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/i18n"
	"github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/validation/process"
)

// RevocationAcceptanceCheckerResultCheck verifies if the RAC result is valid. T
// is the XmlConstraintsConclusion result type.
type RevocationAcceptanceCheckerResultCheck[T any] struct {
	*process.ChainItemBase[T]

	// racResult is the Revocation Acceptance Checker result.
	racResult *jaxb.XmlRAC
}

// NewRevocationAcceptanceCheckerResultCheck is the default constructor. Port of
// RevocationAcceptanceCheckerResultCheck(I18nProvider, T, XmlRAC, LevelRule).
//
// Java always calls the ChainItem constructor that takes a bbbId, with the RAC
// id, which is a nullable String; the Go port picks the id-less constructor for
// a null one, the two being the same call there.
func NewRevocationAcceptanceCheckerResultCheck[T any](i18nProvider *i18n.I18nProvider, result *process.Result[T],
	racResult *jaxb.XmlRAC, constraint policy.LevelRule) *RevocationAcceptanceCheckerResultCheck[T] {
	var base *process.ChainItemBase[T]
	if racResult.Id != nil {
		base = process.NewChainItemBaseWithId(i18nProvider, result, constraint, *racResult.Id)
	} else {
		base = process.NewChainItemBase(i18nProvider, result, constraint)
	}
	c := &RevocationAcceptanceCheckerResultCheck[T]{
		ChainItemBase: base,
		racResult:     racResult,
	}
	c.InitChainItem(c)
	return c
}

// BlockType returns the validating block type. Port of getBlockType().
func (c *RevocationAcceptanceCheckerResultCheck[T]) BlockType() jaxb.XmlBlockType {
	return jaxb.XmlBlockType_RAC
}

// Process performs the check. Port of process().
func (c *RevocationAcceptanceCheckerResultCheck[T]) Process() bool {
	return c.IsValid(&c.racResult.XmlConstraintsConclusionContent)
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *RevocationAcceptanceCheckerResultCheck[T]) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_XCV_RAC
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *RevocationAcceptanceCheckerResultCheck[T]) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_XCV_RAC_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *RevocationAcceptanceCheckerResultCheck[T]) FailedIndicationForConclusion() enumerations.Indication {
	return c.racResult.Conclusion.Indication.Indication()
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure. Port
// of getFailedSubIndicationForConclusion(): the generated Go SubIndication
// member is a pointer, whose nil is Java's null.
func (c *RevocationAcceptanceCheckerResultCheck[T]) FailedSubIndicationForConclusion() enumerations.SubIndication {
	if c.racResult.Conclusion.SubIndication == nil {
		return ""
	}
	return c.racResult.Conclusion.SubIndication.SubIndication()
}

// BuildAdditionalInfo builds an additional information. Port of
// buildAdditionalInfo().
//
// Java tests getRevocationProductionDate() for null. The generated Go model
// carries RevocationThisUpdate and RevocationProductionDate as plain XSDateTime
// members (DetailedReport.xsd declares both required), so a Java null shows up
// as the zero time and the null tests are on that.
//
// The dates then go through ValidationProcessUtils#getFormattedDate, which
// answers null in Java and the empty string in Go (see its header): a
// revocation with a production date but no thisUpdate - which
// ThisUpdatePresenceCheck exists to catch - therefore renders "" here where
// Java renders "null". A null RAC id, which RevocationAcceptanceChecker never
// produces, renders the same way.
func (c *RevocationAcceptanceCheckerResultCheck[T]) BuildAdditionalInfo() *string {
	if !time.Time(c.racResult.RevocationProductionDate).IsZero() {
		var thisUpdateDate *time.Time
		if t := time.Time(c.racResult.RevocationThisUpdate); !t.IsZero() {
			thisUpdateDate = &t
		}
		productionDateValue := time.Time(c.racResult.RevocationProductionDate)
		thisUpdate := process.GetFormattedDate(thisUpdateDate)
		productionDate := process.GetFormattedDate(&productionDateValue)
		var id string
		if c.racResult.Id != nil {
			id = *c.racResult.Id
		}
		message := c.I18nProvider.GetMessage(i18n.MessageTag_REVOCATION_ACCEPTANCE_CHECK, id,
			thisUpdate, productionDate)
		return &message
	}
	return nil
}

// PreviousErrors returns a list of previous errors occurred in the chain. Port
// of getPreviousErrors().
func (c *RevocationAcceptanceCheckerResultCheck[T]) PreviousErrors() []*jaxb.XmlMessage {
	return c.racResult.Conclusion.Errors
}
