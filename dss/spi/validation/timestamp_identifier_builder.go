// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/x509/tsp/TimestampIdentifierBuilder.java (DSS 6.5.RC1).
package validation

import (
	"bytes"
	"encoding/binary"
	"unicode/utf16"

	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/utils"
)

// TimestampIdentifierBuilderOverrides is the contract a subclass of TimestampIdentifierBuilder
// implements. Port of the single protected method every upstream subclass
// (SignatureTimestampIdentifierBuilder, EvidenceRecordTimestampIdentifierBuilder,
// PdfTimestampTokenIdentifierBuilder, VriDictionaryTimestampIdentifierBuilder) overrides;
// dispatched the way model.TokenBase dispatches to model.TokenOverrides via InitToken.
type TimestampIdentifierBuilderOverrides interface {
	// TimestampPosition returns a position of a time-stamp token within a document among other
	// time-stamps. Port of the protected getTimestampPosition().
	//
	// Java declares the return type Object and interpolates it through
	// StringBuilder#append(Object), i.e. String.valueOf; every override already returns a
	// String, and the base returns Utils.EMPTY_STRING, so the Go port narrows it to string.
	TimestampPosition() string
}

// TimestampTokenIdentifierBuilder is the contract TimestampToken needs from the
// TimestampIdentifierBuilder it is handed - the name the javadoc of
// TimestampToken#getTimestampIdentifierBuilder uses for it. Go cannot pass a subclass where the
// base struct type is declared, so the Java parameter type becomes this interface, which
// *TimestampIdentifierBuilder and every type embedding it satisfy.
type TimestampTokenIdentifierBuilder interface {
	// BuildTimestampTokenIdentifier builds the identifier of the time-stamp token.
	BuildTimestampTokenIdentifier() *TimestampTokenIdentifier
}

// TimestampIdentifierBuilder builds a TimestampTokenIdentifier for the provided TimestampToken.
// This class provides a format independent implementation; use the inherited classes for
// format-specific implementations.
type TimestampIdentifierBuilder struct {
	// overrides points back at the concrete builder; see InitTimestampIdentifierBuilder. It is
	// nil for a bare base instance, which then behaves exactly as the (non-abstract) Java base
	// class does.
	overrides TimestampIdentifierBuilderOverrides

	// timestampTokenBinaries is the time-stamp token to build an identifier for. Java declares
	// the field protected; a Go subclass in another package reaches it through
	// TimestampTokenBinaries().
	timestampTokenBinaries []byte

	// filename is the name of the document containing the time-stamp token.
	filename string
}

// NewTimestampIdentifierBuilder builds an implementation-independent identifier builder over
// the given DER-encoded time-stamp binaries. Port of the
// TimestampIdentifierBuilder(byte[]) constructor used directly, i.e. not through a subclass.
//
// Panics with the Java message when the binaries are missing (Objects.requireNonNull).
func NewTimestampIdentifierBuilder(timestampTokenBinaries []byte) *TimestampIdentifierBuilder {
	builder := NewTimestampIdentifierBuilderBase(timestampTokenBinaries)
	builder.InitTimestampIdentifierBuilder(&builder)
	return &builder
}

// NewTimestampIdentifierBuilderBase builds the base state a subclass embeds. Port of the same
// TimestampIdentifierBuilder(byte[]) constructor reached through super(...); the subclass
// constructor must follow it with InitTimestampIdentifierBuilder.
//
// Panics with the Java message when the binaries are missing (Objects.requireNonNull).
func NewTimestampIdentifierBuilderBase(timestampTokenBinaries []byte) TimestampIdentifierBuilder {
	if timestampTokenBinaries == nil {
		panic("Timestamp token binaries cannot be null!")
	}
	return TimestampIdentifierBuilder{timestampTokenBinaries: timestampTokenBinaries}
}

// InitTimestampIdentifierBuilder registers the concrete builder with its base so that the base
// can dispatch TimestampPosition. Every subclass constructor must call this once.
func (b *TimestampIdentifierBuilder) InitTimestampIdentifierBuilder(overrides TimestampIdentifierBuilderOverrides) {
	b.overrides = overrides
}

// TimestampTokenBinaries returns the time-stamp binaries the identifier is built over; the Go
// counterpart of reading Java's protected timestampTokenBinaries field from a subclass.
func (b *TimestampIdentifierBuilder) TimestampTokenBinaries() []byte {
	return b.timestampTokenBinaries
}

// SetFilename sets the time-stamp document filename. Port of setFilename(String), which returns
// this for chaining.
func (b *TimestampIdentifierBuilder) SetFilename(filename string) *TimestampIdentifierBuilder {
	b.filename = filename
	return b
}

// Build builds the TimestampTokenIdentifier for the current time-stamp token. Port of the
// build() override, satisfying model.IdentifierBuilder; BuildTimestampTokenIdentifier is the
// narrowly typed counterpart Java expresses with a covariant return type.
func (b *TimestampIdentifierBuilder) Build() model.Identifier {
	return b.BuildTimestampTokenIdentifier()
}

// BuildTimestampTokenIdentifier builds the TimestampTokenIdentifier for the current time-stamp
// token. Port of build() with its Java return type.
func (b *TimestampIdentifierBuilder) BuildTimestampTokenIdentifier() *TimestampTokenIdentifier {
	return NewTimestampTokenIdentifier(b.BuildBinaries())
}

// BuildBinaries builds the unique binary data describing the time-stamp token.
// Port of the protected buildBinaries().
//
// Java wraps the IOException a ByteArrayOutputStream can never raise in a DSSException; a
// bytes.Buffer cannot fail either, so the Go port has no error to return and the dead branch
// disappears.
func (b *TimestampIdentifierBuilder) BuildBinaries() []byte {
	buffer := &bytes.Buffer{}
	b.WriteTimestampBinaries(buffer)
	b.WriteTimestampPosition(buffer)
	return buffer.Bytes()
}

// WriteTimestampBinaries writes the DER-encoded binaries of the current time-stamp token to the
// given buffer. Port of the protected writeTimestampBinaries(ByteArrayOutputStream).
func (b *TimestampIdentifierBuilder) WriteTimestampBinaries(buffer *bytes.Buffer) {
	buffer.Write(b.timestampTokenBinaries)
}

// WriteTimestampPosition writes the current time-stamp position within a document.
// Port of the protected writeTimestampPosition(ByteArrayOutputStream).
//
// Java wraps the buffer in a DataOutputStream and calls writeChars, which emits the UTF-16BE
// code units of the string; that is what is reproduced here. The Java positionId is never null
// (getUniquePositionId always returns a StringBuilder's toString), so its null check has no
// Go counterpart.
func (b *TimestampIdentifierBuilder) WriteTimestampPosition(buffer *bytes.Buffer) {
	for _, unit := range utf16.Encode([]rune(b.UniquePositionId())) {
		_ = binary.Write(buffer, binary.BigEndian, unit)
	}
}

// UniquePositionId returns the id representing the current time-stamp position in a file,
// considering its pre-siblings and master signatures when present.
// Port of the protected getUniquePositionId().
//
// Java's `filename != null` test becomes an empty-string test: Go has no nullable string, and
// appending the empty string is the same operation as skipping it.
func (b *TimestampIdentifierBuilder) UniquePositionId() string {
	return b.timestampPosition() + b.filename
}

// TimestampPosition returns the empty string, the position of a time-stamp token within a
// document among other time-stamps for the format-independent builder.
// Port of the protected getTimestampPosition() base implementation.
func (b *TimestampIdentifierBuilder) TimestampPosition() string {
	return utils.EmptyString
}

// timestampPosition dispatches to the registered subclass, falling back to this class's own
// implementation for a bare base instance (Java's TimestampIdentifierBuilder is not abstract).
func (b *TimestampIdentifierBuilder) timestampPosition() string {
	if b.overrides == nil {
		return b.TimestampPosition()
	}
	return b.overrides.TimestampPosition()
}

// compile-time assertions: a TimestampIdentifierBuilder is an IdentifierBuilder and satisfies
// its own overrides contract.
var (
	_ model.IdentifierBuilder             = (*TimestampIdentifierBuilder)(nil)
	_ TimestampIdentifierBuilderOverrides = (*TimestampIdentifierBuilder)(nil)
	_ TimestampTokenIdentifierBuilder     = (*TimestampIdentifierBuilder)(nil)
)
