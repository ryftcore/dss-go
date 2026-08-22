// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/signature/identifier/AbstractSignatureIdentifierBuilder.java (DSS 6.5.RC1).
package validation

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"
	"time"
	"unicode/utf16"

	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/utils"
)

// AbstractSignatureIdentifierBuilderMETAINFFolder is the META-INF folder (used to determine a
// signature file position in an ASiC container). Port of the public static final
// META_INF_FOLDER constant.
const AbstractSignatureIdentifierBuilderMETAINFFolder = "META-INF/"

// AbstractSignatureIdentifierBuilderOverrides declares the operations
// AbstractSignatureIdentifierBuilder calls back into virtually: every concrete builder (a later
// phase's CAdES/XAdES/JAdES/PAdES/CB-AdES SignatureIdentifierBuilder) registers itself with
// InitAbstractSignatureIdentifierBuilder so the base can dispatch them, the way model.TokenBase
// dispatches to model.TokenOverrides via InitToken.
//
// Java types these Object (widened for a heterogeneous StringBuilder#append(Object) call, i.e.
// String.valueOf semantics); concrete overrides return Integer in every subclass but PAdES's
// getSignaturePosition (String) and PAdES's getCounterSignaturePosition / XAdES's
// getSignatureFilePosition (Object). Go has no common numeric-or-string supertype, so these stay
// `any`, converted to text the same way Java's String.valueOf(Object) would (see
// abstractSignatureIdentifierBuilderStringValueOf).
type AbstractSignatureIdentifierBuilderOverrides interface {
	// CounterSignaturePosition returns the current counter signature position in its master
	// signature. Port of the protected abstract getCounterSignaturePosition(AdvancedSignature).
	CounterSignaturePosition(masterSignature AdvancedSignature) any

	// SignaturePosition returns a position of a signature in the provided file. Port of the
	// protected abstract getSignaturePosition().
	SignaturePosition() any

	// SignatureFilePosition returns a position of a signature file.
	// NOTE: this method returns a signature filename for ASiC containers, empty string for
	// others. Port of the protected getSignatureFilePosition() (default: empty, overridable).
	SignatureFilePosition() any
}

// AbstractSignatureIdentifierBuilder is the abstract SignatureIdentifier builder.
type AbstractSignatureIdentifierBuilder struct {
	// overrides points back at the concrete builder; see InitAbstractSignatureIdentifierBuilder.
	overrides AbstractSignatureIdentifierBuilderOverrides

	// signature is the signature to build the identifier for. Java declares the field protected
	// final; a Go subclass in another package reaches it through Signature().
	signature AdvancedSignature
}

// NewAbstractSignatureIdentifierBuilderBase builds the base state a subclass embeds, over the
// given signature. Port of the protected AbstractSignatureIdentifierBuilder(AdvancedSignature)
// constructor; the subclass constructor must follow it with
// InitAbstractSignatureIdentifierBuilder.
func NewAbstractSignatureIdentifierBuilderBase(signatureValue AdvancedSignature) AbstractSignatureIdentifierBuilder {
	return AbstractSignatureIdentifierBuilder{signature: signatureValue}
}

// InitAbstractSignatureIdentifierBuilder registers the concrete builder with its base so that
// the base can dispatch CounterSignaturePosition/SignaturePosition/SignatureFilePosition. Every
// concrete subclass constructor must call this once.
func (b *AbstractSignatureIdentifierBuilder) InitAbstractSignatureIdentifierBuilder(overrides AbstractSignatureIdentifierBuilderOverrides) {
	b.overrides = overrides
}

// abstractSignatureIdentifierBuilderOverrides returns the registered overrides, panicking when
// the concrete builder forgot to call InitAbstractSignatureIdentifierBuilder - Java's abstract
// class can never be instantiated bare.
func (b *AbstractSignatureIdentifierBuilder) abstractSignatureIdentifierBuilderOverrides() AbstractSignatureIdentifierBuilderOverrides {
	if b.overrides == nil {
		panic("AbstractSignatureIdentifierBuilder was not initialised: the concrete builder must call InitAbstractSignatureIdentifierBuilder in its constructor")
	}
	return b.overrides
}

// Signature returns the signature the identifier is built for. The Go counterpart of reading
// Java's protected final signature field from a subclass.
func (b *AbstractSignatureIdentifierBuilder) Signature() AdvancedSignature {
	return b.signature
}

// Build builds the SignatureIdentifier for the provided AdvancedSignature. Port of build(),
// satisfying model.IdentifierBuilder; BuildSignatureIdentifier is the narrowly typed counterpart
// Java expresses with a covariant return type.
func (b *AbstractSignatureIdentifierBuilder) Build() model.Identifier {
	return b.BuildSignatureIdentifier()
}

// BuildSignatureIdentifier builds the SignatureIdentifier for the provided AdvancedSignature.
// Port of build() with its Java return type.
func (b *AbstractSignatureIdentifierBuilder) BuildSignatureIdentifier() *SignatureIdentifier {
	return NewSignatureIdentifier(b.BuildBinaries())
}

// BuildBinaries builds unique binary data describing the signature object.
// Port of the protected buildBinaries().
//
// Java wraps the IOException a ByteArrayOutputStream can never raise in a DSSException; a
// bytes.Buffer cannot fail either, so the Go port has no error to return and the dead branch
// disappears.
func (b *AbstractSignatureIdentifierBuilder) BuildBinaries() []byte {
	buffer := &bytes.Buffer{}
	b.WriteSignedProperties(buffer)
	b.WriteSignaturePosition(buffer)
	return buffer.Bytes()
}

// WriteSignedProperties writes signed properties of a signature to the given buffer.
// Port of the protected writeSignedProperties(ByteArrayOutputStream).
func (b *AbstractSignatureIdentifierBuilder) WriteSignedProperties(buffer *bytes.Buffer) {
	b.writeSigningTime(buffer, b.signature.SigningTime())
	b.writeSigningCertificateRefs(buffer, b.signature.CertificateSource().SigningCertificateRefs())
	b.writeSignatureValue(buffer, b.signature.SignatureValue())
}

// writeSigningTime is the private writeSigningTime(ByteArrayOutputStream, Date): the epoch
// milliseconds of signingTime as an 8-byte big-endian value (DataOutputStream#writeLong), or
// nothing when signingTime is nil.
func (b *AbstractSignatureIdentifierBuilder) writeSigningTime(buffer *bytes.Buffer, signingTime *time.Time) {
	if signingTime != nil {
		_ = binary.Write(buffer, binary.BigEndian, signingTime.UnixMilli())
	}
}

// writeSigningCertificateRefs is the private writeSigningCertificateRefs(ByteArrayOutputStream,
// List<CertificateRef>).
func (b *AbstractSignatureIdentifierBuilder) writeSigningCertificateRefs(buffer *bytes.Buffer, signingCertificateRefs []*spi.CertificateRef) {
	if utils.IsCollectionNotEmpty(signingCertificateRefs) {
		for _, certificateRef := range signingCertificateRefs {
			b.WriteString(buffer, certificateRef.DSSIDAsString())
		}
	}
}

// writeSignatureValue is the private writeSignatureValue(ByteArrayOutputStream, byte[]).
func (b *AbstractSignatureIdentifierBuilder) writeSignatureValue(buffer *bytes.Buffer, signatureValue []byte) {
	if utils.IsArrayNotEmpty(signatureValue) {
		buffer.Write(signatureValue)
	}
}

// WriteString writes str into buffer. Port of the protected writeString(ByteArrayOutputStream,
// String).
//
// Java wraps the buffer in a DataOutputStream and calls writeChars, which emits the UTF-16BE
// code units of the string; that is what is reproduced here (see also
// TimestampIdentifierBuilder.WriteTimestampPosition, which does the same for the sibling
// TSP identifier builder). Go strings are never nil, so Java's `str != null` guard has no
// counterpart: an empty string simply contributes zero code units either way.
func (b *AbstractSignatureIdentifierBuilder) WriteString(buffer *bytes.Buffer, str string) {
	for _, unit := range utf16.Encode([]rune(str)) {
		_ = binary.Write(buffer, binary.BigEndian, unit)
	}
}

// WriteSignaturePosition writes the current signature position between other signature entries
// on the same level. Port of the protected writeSignaturePosition(ByteArrayOutputStream).
func (b *AbstractSignatureIdentifierBuilder) WriteSignaturePosition(buffer *bytes.Buffer) {
	b.WriteString(buffer, b.PositionId())
}

// PositionId returns the Id representing the current signature position in a file, considering
// its pre-siblings and master signatures when present.
// Port of the protected getPositionId().
func (b *AbstractSignatureIdentifierBuilder) PositionId() string {
	var result strings.Builder

	masterSignature := b.signature.MasterSignature()
	if masterSignature != nil {
		result.WriteString(masterSignature.ID())
		result.WriteString(abstractSignatureIdentifierBuilderStringValueOf(
			b.abstractSignatureIdentifierBuilderOverrides().CounterSignaturePosition(masterSignature)))
	} else {
		result.WriteString(abstractSignatureIdentifierBuilderStringValueOf(
			b.abstractSignatureIdentifierBuilderOverrides().SignaturePosition()))
		result.WriteString(abstractSignatureIdentifierBuilderStringValueOf(
			b.abstractSignatureIdentifierBuilderOverrides().SignatureFilePosition()))
	}

	return result.String()
}

// SignatureFilePosition returns a position of a signature file.
// NOTE: this method returns a signature filename for ASiC containers, empty string for others.
// Port of the protected getSignatureFilePosition(); empty by default, overridable (only XAdES
// overrides it upstream).
func (b *AbstractSignatureIdentifierBuilder) SignatureFilePosition() any {
	return utils.EmptyString
}

// abstractSignatureIdentifierBuilderStringValueOf mirrors Java's String.valueOf(Object): "null"
// for a nil value (StringBuilder#append(Object) prints the literal string "null" for a null
// argument), fmt.Sprint otherwise (matching Integer/String#toString for the concrete types every
// upstream override actually returns).
func abstractSignatureIdentifierBuilderStringValueOf(value any) string {
	if value == nil {
		return "null"
	}
	return fmt.Sprint(value)
}

// compile-time assertions: an AbstractSignatureIdentifierBuilder is an IdentifierBuilder and
// satisfies its own overrides contract (SignatureFilePosition has a base default; the two
// genuinely abstract methods are supplied by whatever concrete builder embeds this struct).
var (
	_ model.IdentifierBuilder = (*AbstractSignatureIdentifierBuilder)(nil)
)
