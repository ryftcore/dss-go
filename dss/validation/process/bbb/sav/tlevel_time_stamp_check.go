// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/sav/checks/TLevelTimeStampCheck.java (DSS 6.5.RC1).
package sav

import (
	jaxb "github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// TLevelTimeStampCheck verifies if there is at least one valid T-level
// timestamp.
type TLevelTimeStampCheck[T any] struct {
	*AbstractTimeStampPresentCheck[T]

	// signature is the signature to check.
	signature *diagnostic.SignatureWrapper
}

// NewTLevelTimeStampCheck is the default constructor. Port of
// TLevelTimeStampCheck(I18nProvider, T, SignatureWrapper, Map, Collection, LevelRule).
func NewTLevelTimeStampCheck[T any](i18nProvider *i18n.I18nProvider, result *process.Result[T],
	signature *diagnostic.SignatureWrapper, bbbs map[string]*jaxb.XmlBasicBuildingBlocks,
	xmlTimestamps []*jaxb.XmlTimestamp, constraint policy.LevelRule) *TLevelTimeStampCheck[T] {
	c := &TLevelTimeStampCheck[T]{
		AbstractTimeStampPresentCheck: NewAbstractTimeStampPresentCheck(i18nProvider, result, bbbs, xmlTimestamps, constraint),
		signature:                     signature,
	}
	c.InitAbstractTimeStampPresentCheck(c)
	c.InitChainItem(c)
	return c
}

// Timestamps returns the collection of timestamps to be checked. Port of
// getTimestamps().
func (c *TLevelTimeStampCheck[T]) Timestamps() []*diagnostic.TimestampWrapper {
	return c.signature.TLevelTimestamps()
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *TLevelTimeStampCheck[T]) MessageTag() i18n.MessageTag {
	return i18n.MessageTagBBBSAVIVTTSTP
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *TLevelTimeStampCheck[T]) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagBBBSAVIVTTSTPANS
}
