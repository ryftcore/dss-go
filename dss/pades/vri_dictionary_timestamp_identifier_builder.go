// Ported from dss-pades/src/main/java/eu/europa/esig/dss/pades/validation/dss/VriDictionaryTimestampIdentifierBuilder.java (DSS 6.5.RC1).
package pades

import "github.com/utain/esig/dss/spi/validation"

// VriDictionaryTimestampIdentifierBuilder builds a unique identifier for a time-stamp
// encapsulated within a VRI dictionary.
type VriDictionaryTimestampIdentifierBuilder struct {
	validation.TimestampIdentifierBuilder

	// name is the identifier of the corresponding VRI dictionary in the PDF document.
	name string
}

// NewVriDictionaryTimestampIdentifierBuilder builds an identifier for a time-stamp extracted
// from a VRI dictionary, over the DER-encoded timestampTokenBinaries and the name identifier of
// the /VRI dictionary. Port of the
// VriDictionaryTimestampIdentifierBuilder(byte[], String) constructor.
//
// Panics with the Java message when the binaries are missing (Objects.requireNonNull in the
// super constructor).
func NewVriDictionaryTimestampIdentifierBuilder(timestampTokenBinaries []byte, name string) *VriDictionaryTimestampIdentifierBuilder {
	builder := &VriDictionaryTimestampIdentifierBuilder{
		TimestampIdentifierBuilder: validation.NewTimestampIdentifierBuilderBase(timestampTokenBinaries),
		name:                       name,
	}
	builder.InitTimestampIdentifierBuilder(builder)
	return builder
}

// TimestampPosition returns the name of the /VRI dictionary carrying the time-stamp.
// Port of the protected getTimestampPosition() override.
func (b *VriDictionaryTimestampIdentifierBuilder) TimestampPosition() string {
	return b.name
}

// compile-time assertions: a VriDictionaryTimestampIdentifierBuilder is a
// TimestampTokenIdentifierBuilder and overrides TimestampPosition.
var (
	_ validation.TimestampTokenIdentifierBuilder     = (*VriDictionaryTimestampIdentifierBuilder)(nil)
	_ validation.TimestampIdentifierBuilderOverrides = (*VriDictionaryTimestampIdentifierBuilder)(nil)
)
