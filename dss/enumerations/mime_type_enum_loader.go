// Ported from dss-enumerations/.../MimeTypeEnumLoader.java (DSS 6.5.RC1).
//
// NOTE: MimeTypeEnum (with its MimeTypeEnumValues(), MimeTypeString() and
// Extensions() []string accessors) is defined outside this file's manifest
// and is assumed to exist per the porting brief. Extensions() is assumed to
// expose the full extensions list (Java's package-private `extensions`
// array field, looped over in fromFileExtension), as opposed to the single
// first extension returned by MimeType's Extension().
package enumerations

import "strings"

// MimeTypeEnumLoader is a MimeTypeLoader implementation backed by the
// built-in MimeTypeEnum enumeration.
type MimeTypeEnumLoader struct{}

// NewMimeTypeEnumLoader creates a new MimeTypeEnumLoader.
func NewMimeTypeEnumLoader() *MimeTypeEnumLoader {
	return &MimeTypeEnumLoader{}
}

// FromMimeTypeString returns the MimeTypeEnum matching mimeTypeString
// (case-insensitively), or nil if none matches.
func (l *MimeTypeEnumLoader) FromMimeTypeString(mimeTypeString string) MimeType {
	for _, mimeTypeEnum := range MimeTypeEnumValues() {
		if strings.EqualFold(mimeTypeString, mimeTypeEnum.MimeTypeString()) {
			return mimeTypeEnum
		}
	}
	return nil
}

// FromFileExtension returns the MimeTypeEnum matching fileExtension
// (case-insensitively), or nil if none matches.
func (l *MimeTypeEnumLoader) FromFileExtension(fileExtension string) MimeType {
	for _, mimeTypeEnum := range MimeTypeEnumValues() {
		for _, extension := range mimeTypeEnum.Extensions() {
			if strings.EqualFold(fileExtension, extension) {
				return mimeTypeEnum
			}
		}
	}
	return nil
}

// compile-time interface assertion.
var _ MimeTypeLoader = (*MimeTypeEnumLoader)(nil)
