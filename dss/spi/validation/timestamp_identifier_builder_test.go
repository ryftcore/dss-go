package validation

import (
	"bytes"
	"encoding/binary"
	"strings"
	"testing"
	"unicode/utf16"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
)

// timestampIdentifierBuilderTestSubclass stands in for the format-specific builders of the later
// phases (SignatureTimestampIdentifierBuilder and friends): it overrides nothing but the
// time-stamp position, which is what every upstream subclass does.
type timestampIdentifierBuilderTestSubclass struct {
	TimestampIdentifierBuilder
	position string
}

// TimestampPosition returns the pinned position.
func (b *timestampIdentifierBuilderTestSubclass) TimestampPosition() string {
	return b.position
}

// newTimestampIdentifierBuilderTestSubclass wires the subclass the way a ported one has to.
func newTimestampIdentifierBuilderTestSubclass(binaries []byte, position string) *timestampIdentifierBuilderTestSubclass {
	builder := &timestampIdentifierBuilderTestSubclass{
		TimestampIdentifierBuilder: NewTimestampIdentifierBuilderBase(binaries),
		position:                   position,
	}
	builder.InitTimestampIdentifierBuilder(builder)
	return builder
}

// TestTimestampIdentifierBuilderBinaries checks what the identifier is computed over: the
// time-stamp binaries followed by the UTF-16BE code units of the position id, which is the
// position concatenated with the filename.
func TestTimestampIdentifierBuilderBinaries(t *testing.T) {
	binaries := []byte{0x01, 0x02, 0x03}

	if got := NewTimestampIdentifierBuilder(binaries).BuildBinaries(); !bytes.Equal(got, binaries) {
		t.Errorf("BuildBinaries() = %x, want the bare binaries %x", got, binaries)
	}
	if got := NewTimestampIdentifierBuilder(binaries).UniquePositionId(); got != "" {
		t.Errorf("UniquePositionId() = %q, want the empty string", got)
	}

	withFilename := NewTimestampIdentifierBuilder(binaries).SetFilename("ts.tst")
	if got, want := withFilename.BuildBinaries(), append(append([]byte{}, binaries...),
		timestampIdentifierBuilderTestChars("ts.tst")...); !bytes.Equal(got, want) {
		t.Errorf("BuildBinaries() = %x, want %x", got, want)
	}

	subclass := newTimestampIdentifierBuilderTestSubclass(binaries, "-OOA-1")
	subclass.SetFilename("ts.tst")
	if got := subclass.UniquePositionId(); got != "-OOA-1ts.tst" {
		t.Errorf("UniquePositionId() = %q, want -OOA-1ts.tst", got)
	}
	if got, want := subclass.BuildBinaries(), append(append([]byte{}, binaries...),
		timestampIdentifierBuilderTestChars("-OOA-1ts.tst")...); !bytes.Equal(got, want) {
		t.Errorf("BuildBinaries() = %x, want %x", got, want)
	}
}

// TestTimestampIdentifierBuilderIdentifiers checks that the position and the filename actually
// separate two identifiers built over the same time-stamp.
func TestTimestampIdentifierBuilderIdentifiers(t *testing.T) {
	binaries := timestampTokenKATFile(t, "timestamp-token.tst")

	bare := NewTimestampIdentifierBuilder(binaries).BuildTimestampTokenIdentifier()
	named := NewTimestampIdentifierBuilder(binaries).SetFilename("ts.tst").BuildTimestampTokenIdentifier()
	positioned := newTimestampIdentifierBuilderTestSubclass(binaries, "-OOA-1").BuildTimestampTokenIdentifier()

	for _, identifier := range []*TimestampTokenIdentifier{bare, named, positioned} {
		if !strings.HasPrefix(identifier.AsXmlID(), "T-") {
			t.Errorf("AsXmlID() = %s, want the T- prefix", identifier.AsXmlID())
		}
	}
	if bare.AsXmlID() == named.AsXmlID() {
		t.Error("the filename did not change the identifier")
	}
	if bare.AsXmlID() == positioned.AsXmlID() {
		t.Error("the position did not change the identifier")
	}
	// Build() is the model.IdentifierBuilder view of the same value.
	var builder model.IdentifierBuilder = NewTimestampIdentifierBuilder(binaries)
	if builder.Build().AsXmlID() != bare.AsXmlID() {
		t.Error("Build() and BuildTimestampTokenIdentifier() disagree")
	}

	// A token built with an explicit builder takes its identifier from it.
	token, err := NewTimestampTokenWithIdentifierBuilder(binaries,
		enumerations.TimestampTypeSignatureTimestamp, nil,
		newTimestampIdentifierBuilderTestSubclass(binaries, "-OOA-1"))
	if err != nil {
		t.Fatalf("NewTimestampTokenWithIdentifierBuilder() failed: %v", err)
	}
	if token.DSSIDAsString() != positioned.AsXmlID() {
		t.Errorf("DSSIDAsString() = %s, want the builder's %s", token.DSSIDAsString(), positioned.AsXmlID())
	}

	// Without one, the token builds the format-independent builder over its own filename.
	detached, err := NewTimestampToken(binaries, enumerations.TimestampTypeContainerTimestamp)
	if err != nil {
		t.Fatalf("NewTimestampToken() failed: %v", err)
	}
	detached.SetFilename("ts.tst")
	if detached.DSSIDAsString() != named.AsXmlID() {
		t.Errorf("DSSIDAsString() = %s, want %s", detached.DSSIDAsString(), named.AsXmlID())
	}
}

// TestTimestampIdentifierBuilderRequiresBinaries reproduces the Objects.requireNonNull guard.
func TestTimestampIdentifierBuilderRequiresBinaries(t *testing.T) {
	defer func() {
		if recovered := recover(); recovered != "Timestamp token binaries cannot be null!" {
			t.Errorf("recover() = %v, want the requireNonNull message", recovered)
		}
	}()
	NewTimestampIdentifierBuilder(nil)
	t.Error("NewTimestampIdentifierBuilder(nil) returned, want a panic")
}

// timestampIdentifierBuilderTestChars reproduces DataOutputStream#writeChars.
func timestampIdentifierBuilderTestChars(value string) []byte {
	var encoded []byte
	for _, unit := range utf16.Encode([]rune(value)) {
		encoded = binary.BigEndian.AppendUint16(encoded, unit)
	}
	return encoded
}
