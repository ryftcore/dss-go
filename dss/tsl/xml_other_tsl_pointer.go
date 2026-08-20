// Ported from dss-tsl-validation/src/main/java/eu/europa/esig/dss/tsl/function/XMLOtherTSLPointer.java (DSS 6.5.RC1).
package tsl

// xmlOtherTSLPointerExpectedMimetype is the private static EXPECTED_MIMETYPE.
const xmlOtherTSLPointerExpectedMimetype = "application/vnd.etsi.tsl+xml"

// XMLOtherTSLPointer allows filtering of TSL pointers by an XML MimeType.
type XMLOtherTSLPointer struct {
	MimetypeOtherTSLPointer
}

// NewXMLOtherTSLPointer is the default constructor. Port of XMLOtherTSLPointer().
func NewXMLOtherTSLPointer() *XMLOtherTSLPointer {
	return &XMLOtherTSLPointer{MimetypeOtherTSLPointer: *NewMimetypeOtherTSLPointer(xmlOtherTSLPointerExpectedMimetype)}
}
