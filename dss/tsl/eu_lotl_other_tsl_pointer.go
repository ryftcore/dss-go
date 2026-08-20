// Ported from dss-tsl-validation/src/main/java/eu/europa/esig/dss/tsl/function/EULOTLOtherTSLPointer.java (DSS 6.5.RC1).
package tsl

// eulotlOtherTSLPointerExpectedEULOTLType is the private static EXPECTED_EU_LOTL_TYPE.
const eulotlOtherTSLPointerExpectedEULOTLType = "http://uri.etsi.org/TrstSvc/TrustedList/TSLType/EUlistofthelists"

// EULOTLOtherTSLPointer selects OtherTSLPointerType(s) with a defined type equals to
// EUlistofthelists.
type EULOTLOtherTSLPointer struct {
	TypeOtherTSLPointer
}

// NewEULOTLOtherTSLPointer is the default constructor. Port of EULOTLOtherTSLPointer().
func NewEULOTLOtherTSLPointer() *EULOTLOtherTSLPointer {
	return &EULOTLOtherTSLPointer{TypeOtherTSLPointer: *NewTypeOtherTSLPointer(eulotlOtherTSLPointerExpectedEULOTLType)}
}
