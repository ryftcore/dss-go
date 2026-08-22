// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/validation/timestamp/SignatureTimestampIdentifierBuilder.java (DSS 6.5.RC1).
//
// org.bouncycastle.tsp.TimeStampToken becomes internal/cmscore.TimeStampToken, per PORTING.md's
// BouncyCastle-replacement rule; TimeStampToken.Encoded() (unlike Java's getEncoded(), which can
// raise IOException) cannot fail; it hands back the bytes the token was parsed from, so the
// private getEncoded(TimeStampToken) helper's try/catch has no Go counterpart.
package timestamp

import (
	"strconv"

	"github.com/ryftcore/dss-go/dss/internal/cmscore"
	"github.com/ryftcore/dss-go/dss/spi/validation"
)

// signatureTimestampIdentifierBuilderOrderOfAttributePrefix is the prefix string for the order
// of attribute value. Port of the private ORDER_OF_ATTRIBUTE_PREFIX constant.
const signatureTimestampIdentifierBuilderOrderOfAttributePrefix = "-OOA-"

// signatureTimestampIdentifierBuilderOrderWithinAttributePrefix is the prefix string for order
// within attribute value. Port of the private ORDER_WITHIN_ATTRIBUTE_PREFIX constant.
const signatureTimestampIdentifierBuilderOrderWithinAttributePrefix = "-OWA-"

// SignatureTimestampIdentifierBuilder builds an identifier for a time-stamp token encapsulated
// within a signature.
type SignatureTimestampIdentifierBuilder struct {
	validation.TimestampIdentifierBuilder

	// signature is the corresponding signature of the time-stamp.
	signature validation.AdvancedSignature

	// attribute is the attribute encapsulating the time-stamp token.
	attribute validation.SignatureAttribute

	// orderOfAttribute is the position of the attribute within the signature; nil stands for
	// Java's null Integer, i.e. "not set".
	orderOfAttribute *int

	// orderWithinAttribute is the position of the current time-stamp within the time-stamp
	// attribute; nil stands for Java's null Integer.
	orderWithinAttribute *int
}

// NewSignatureTimestampIdentifierBuilder builds an identifier for time-stamp token binaries from
// a signature. Port of the SignatureTimestampIdentifierBuilder(byte[]) constructor.
func NewSignatureTimestampIdentifierBuilder(timestampTokenBinaries []byte) *SignatureTimestampIdentifierBuilder {
	builder := &SignatureTimestampIdentifierBuilder{
		TimestampIdentifierBuilder: validation.NewTimestampIdentifierBuilderBase(timestampTokenBinaries),
	}
	builder.InitTimestampIdentifierBuilder(builder)
	return builder
}

// NewSignatureTimestampIdentifierBuilderFromToken builds an identifier for a time-stamp token
// from a signature. Port of the SignatureTimestampIdentifierBuilder(TimeStampToken) constructor.
//
// Panics with the Java message when timeStampToken is missing (Objects.requireNonNull).
func NewSignatureTimestampIdentifierBuilderFromToken(timeStampToken *cmscore.TimeStampToken) *SignatureTimestampIdentifierBuilder {
	if timeStampToken == nil {
		panic("TimeStampToken cannot be null!")
	}
	return NewSignatureTimestampIdentifierBuilder(timeStampToken.Encoded())
}

// SetSignature sets the signature corresponding to the time-stamp token.
// Port of setSignature(AdvancedSignature), which returns this for chaining.
func (b *SignatureTimestampIdentifierBuilder) SetSignature(signature validation.AdvancedSignature) *SignatureTimestampIdentifierBuilder {
	b.signature = signature
	return b
}

// SetAttribute sets a signature attribute encapsulating the time-stamp token.
// Port of setAttribute(SignatureAttribute), which returns this for chaining.
func (b *SignatureTimestampIdentifierBuilder) SetAttribute(attribute validation.SignatureAttribute) *SignatureTimestampIdentifierBuilder {
	b.attribute = attribute
	return b
}

// SetOrderOfAttribute sets the position of the time-stamp carrying attribute within the
// signature. Port of setOrderOfAttribute(Integer), which returns this for chaining.
func (b *SignatureTimestampIdentifierBuilder) SetOrderOfAttribute(orderOfAttribute *int) *SignatureTimestampIdentifierBuilder {
	b.orderOfAttribute = orderOfAttribute
	return b
}

// SetOrderWithinAttribute sets the position of the time-stamp within its carrying attribute.
// Port of setOrderWithinAttribute(Integer), which returns this for chaining.
func (b *SignatureTimestampIdentifierBuilder) SetOrderWithinAttribute(orderWithinAttribute *int) *SignatureTimestampIdentifierBuilder {
	b.orderWithinAttribute = orderWithinAttribute
	return b
}

// TimestampPosition overrides the base TimestampIdentifierBuilder.TimestampPosition().
// Port of the protected getTimestampPosition() override.
func (b *SignatureTimestampIdentifierBuilder) TimestampPosition() string {
	sb := ""
	if b.signature != nil {
		sb += b.signature.ID()
	}
	if b.attribute != nil {
		// SignatureAttribute.Identifier() returns identifier.SignatureAttributeIdentifier by
		// value (not addressable directly off the call result), but AsXmlID is defined on
		// *IdentifierBase; assign to a local first so it can be addressed.
		attributeIdentifier := b.attribute.Identifier()
		sb += attributeIdentifier.AsXmlID()
	}
	if b.orderOfAttribute != nil {
		sb += signatureTimestampIdentifierBuilderOrderOfAttributePrefix
		sb += strconv.Itoa(*b.orderOfAttribute)
	}
	if b.orderWithinAttribute != nil {
		sb += signatureTimestampIdentifierBuilderOrderWithinAttributePrefix
		sb += strconv.Itoa(*b.orderWithinAttribute)
	}
	return sb
}

// compile-time assertions: a SignatureTimestampIdentifierBuilder is a
// TimestampTokenIdentifierBuilder and overrides TimestampPosition.
var (
	_ validation.TimestampTokenIdentifierBuilder     = (*SignatureTimestampIdentifierBuilder)(nil)
	_ validation.TimestampIdentifierBuilderOverrides = (*SignatureTimestampIdentifierBuilder)(nil)
)
