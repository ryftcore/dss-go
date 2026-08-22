// Ported from dss-enumerations/.../COSEStructureType.java (DSS 6.5.RC1).

package enumerations

import "fmt"

// COSEStructureType represents COSE signature structure types defined in RFC 9052,
// 4. Signing Objects.
type COSEStructureType string

const (
	// COSEStructureTypeCoseSign is the COSE_Sign structure: signing with one or more
	// signers (RFC 9052 4.1).
	COSEStructureTypeCoseSign COSEStructureType = "COSE_SIGN"
	// COSEStructureTypeCoseSign1 is the COSE_Sign1 structure: signing with one
	// signer (RFC 9052 4.2).
	COSEStructureTypeCoseSign1 COSEStructureType = "COSE_SIGN1"
)

// COSEStructureTypeValues returns all COSEStructureType constants in declaration order.
func COSEStructureTypeValues() []COSEStructureType {
	return []COSEStructureType{COSEStructureTypeCoseSign, COSEStructureTypeCoseSign1}
}

// COSEStructureTypeValueOf returns the COSEStructureType matching the given Java enum name.
func COSEStructureTypeValueOf(name string) (COSEStructureType, error) {
	for _, v := range COSEStructureTypeValues() {
		if string(v) == name {
			return v, nil
		}
	}
	return "", fmt.Errorf("no enum constant COSEStructureType.%s", name)
}
