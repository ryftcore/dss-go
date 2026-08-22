// Ported from dss-validation/.../validation/process/bbb/fc/checks/ByteRangeCollisionCheck.java (DSS 6.5.RC1).
package fc

import (
	"math/big"

	drjaxb "github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	policy "github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// ByteRangeCollisionCheck checks if the current signature /ByteRange does not collide with other
// signature byte ranges.
type ByteRangeCollisionCheck struct {
	*process.ChainItemBase[*drjaxb.XmlFC]

	currentSignature diagnostic.AbstractSignatureWrapperOverrides
	diagnosticData   *diagnostic.Data
}

// NewByteRangeCollisionCheck is the default constructor.
func NewByteRangeCollisionCheck(i18nProvider *i18n.Provider, result *process.Result[*drjaxb.XmlFC],
	signatureWrapper diagnostic.AbstractSignatureWrapperOverrides, diagnosticData *diagnostic.Data,
	constraint policy.LevelRule) *ByteRangeCollisionCheck {
	c := &ByteRangeCollisionCheck{currentSignature: signatureWrapper, diagnosticData: diagnosticData}
	c.ChainItemBase = process.NewChainItemBase(i18nProvider, result, constraint)
	c.InitChainItem(c)
	return c
}

// Process performs the check.
func (c *ByteRangeCollisionCheck) Process() bool {
	for _, signature := range c.diagnosticData.Signatures() {
		if c.currentSignature.Id() != signature.Id() && collideRevisions(c.currentSignature, signature) {
			return false
		}
	}
	for _, timestamp := range c.diagnosticData.TimestampList() {
		if c.currentSignature.Id() != timestamp.Id() && collideRevisions(c.currentSignature, timestamp) {
			return false
		}
	}
	return true
}

func collideRevisions(one, two diagnostic.AbstractSignatureWrapperOverrides) bool {
	oneRevision := one.PDFRevision()
	twoRevision := two.PDFRevision()
	if oneRevision == nil || twoRevision == nil {
		return false
	}
	return byteRangesCollide(oneRevision.SignatureByteRange(), twoRevision.SignatureByteRange()) ||
		byteRangesCollide(twoRevision.SignatureByteRange(), oneRevision.SignatureByteRange())
}

func byteRangesCollide(byteRangeOne, byteRangeTwo []*big.Int) bool {
	if len(byteRangeOne) != 4 || len(byteRangeTwo) != 4 {
		panic("Signature ByteRange shall have 4 integers!")
	}
	firstOne := firstByteRangePartLength(byteRangeOne)
	return (firstOne < firstByteRangePartLength(byteRangeTwo)) != (firstOne < secondByteRangePartLength(byteRangeTwo))
}

func firstByteRangePartLength(byteRange []*big.Int) int64 {
	return byteRange[0].Int64() + byteRange[1].Int64()
}

func secondByteRangePartLength(byteRange []*big.Int) int64 {
	return byteRange[2].Int64() + byteRange[3].Int64()
}

// MessageTag returns the constraint message i18n key.
func (c *ByteRangeCollisionCheck) MessageTag() i18n.MessageTag { return i18n.MessageTagBBBFCDBTOOST }

// ErrorMessageTag returns the error message i18n key.
func (c *ByteRangeCollisionCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagBBBFCDBTOOSTANS
}

// FailedIndicationForConclusion returns the Indication on failure.
func (c *ByteRangeCollisionCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationFailed
}

// FailedSubIndicationForConclusion returns the SubIndication on failure.
func (c *ByteRangeCollisionCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationFormatFailure
}
