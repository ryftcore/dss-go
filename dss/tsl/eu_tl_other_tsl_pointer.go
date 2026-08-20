// Ported from dss-tsl-validation/src/main/java/eu/europa/esig/dss/tsl/function/EUTLOtherTSLPointer.java (DSS 6.5.RC1).
package tsl

// eutlOtherTSLPointerExpectedEUTLType is the private static EXPECTED_EU_TL_TYPE.
const eutlOtherTSLPointerExpectedEUTLType = "http://uri.etsi.org/TrstSvc/TrustedList/TSLType/EUgeneric"

// EUTLOtherTSLPointer selects OtherTSLPointerType(s) with a defined type equals to EUgeneric.
type EUTLOtherTSLPointer struct {
	TypeOtherTSLPointer
}

// NewEUTLOtherTSLPointer is the default constructor. Port of EUTLOtherTSLPointer().
func NewEUTLOtherTSLPointer() *EUTLOtherTSLPointer {
	return &EUTLOtherTSLPointer{TypeOtherTSLPointer: *NewTypeOtherTSLPointer(eutlOtherTSLPointerExpectedEUTLType)}
}
