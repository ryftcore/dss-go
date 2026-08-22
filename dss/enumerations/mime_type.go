// Ported from dss-enumerations/.../MimeType.java (DSS 6.5.RC1).
package enumerations

import "strings"

// MimeType identifies a file MimeType and its attributes.
type MimeType interface {
	// MimeTypeString returns the String identifying the MimeType.
	MimeTypeString() string

	// Extension returns file extension corresponding to the MimeType.
	Extension() string
}

// MimeTypeFromMimeTypeString returns the first representation of the
// MimeType corresponding to the given mime-type string, composed of two
// parts: a "type" and a "subtype". Falls back to MimeTypeEnumBinary if no
// registered MimeTypeLoader recognizes the string.
func MimeTypeFromMimeTypeString(mimeTypeString string) MimeType {
	for _, loader := range mimeTypeLoaders() {
		if mimeType := loader.FromMimeTypeString(mimeTypeString); mimeType != nil {
			return mimeType
		}
	}
	return MimeTypeEnumBinary
}

// MimeTypeFromFileExtension returns the MimeType matching the provided
// fileExtension, or MimeTypeEnumBinary if none is found.
func MimeTypeFromFileExtension(fileExtension string) MimeType {
	for _, loader := range mimeTypeLoaders() {
		if mimeType := loader.FromFileExtension(fileExtension); mimeType != nil {
			return mimeType
		}
	}
	return MimeTypeEnumBinary
}

// MimeTypeFromFileName returns the mime-type extrapolated from the file
// name, or MimeTypeEnumBinary if none is found.
//
// Upstream distinguishes two "no extension" cases that both collapse to "" in Go:
// getFileExtension returns null only for a null/blank file name, but returns an empty
// string for a name with no '.' (or with a leading '.'). Only the null case short-circuits
// to BINARY; a name like "README" still reaches fromFileExtension(""), consulting every
// registered loader with an empty extension. MimeTypeGetFileExtensionOK reports which case
// applies so that dispatch is preserved — it is observable to any custom MimeTypeLoader,
// even though the built-in MimeTypeEnumLoader matches no empty extension and so ends up at
// BINARY either way.
func MimeTypeFromFileName(fileName string) MimeType {
	fileExtension, ok := MimeTypeGetFileExtensionOK(fileName)
	if ok {
		lowerCaseExtension := strings.ToLower(fileExtension)
		return MimeTypeFromFileExtension(lowerCaseExtension)
	}
	return MimeTypeEnumBinary
}

// MimeTypeFromFilePath returns the mime-type extrapolated from the base
// name of filePath. Corresponds to Java's MimeType.fromFile(File), adapted
// to take a path string rather than a java.io.File.
func MimeTypeFromFilePath(filePath string) MimeType {
	fileName := filePath
	if idx := strings.LastIndexAny(filePath, `/\`); idx >= 0 {
		fileName = filePath[idx+1:]
	}
	return MimeTypeFromFileName(fileName)
}

// MimeTypeGetFileExtension returns the file extension based on the position
// of the '.' in fileName. File paths such as "xxx.y/toto" are not handled.
// Returns "" if fileName is blank or has no extension; use
// MimeTypeGetFileExtensionOK to tell those two cases apart.
func MimeTypeGetFileExtension(fileName string) string {
	ext, _ := MimeTypeGetFileExtensionOK(fileName)
	return ext
}

// MimeTypeGetFileExtensionOK is MimeTypeGetFileExtension with upstream's null/empty
// distinction preserved: ok is false only where Java returns null (a blank file name),
// and true — with a possibly empty extension — where Java returns a string.
func MimeTypeGetFileExtensionOK(fileName string) (string, bool) {
	if strings.TrimSpace(fileName) == "" {
		return "", false
	}

	extension := ""
	lastIndexOf := strings.LastIndex(fileName, ".")
	if lastIndexOf > 0 {
		extension = fileName[lastIndexOf+1:]
	}
	return extension, true
}
