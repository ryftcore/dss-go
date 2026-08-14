// Ported from dss-enumerations/.../Indication.java (DSS 6.5.RC1).

package enumerations

import "fmt"

// Indication lists the possible values for indications.
//
// # Source ETSI EN 319 102-1
//
// Implements UriBasedEnum.
type Indication string

const (
	// Indication_TOTAL_PASSED: when the cryptographic checks of the signature
	// (including checks of hashes of individual data objects that have been signed
	// indirectly) succeeded as well as all checks prescribed by the signature
	// validation policy have been passed.
	Indication_TOTAL_PASSED Indication = "TOTAL_PASSED"
	// Indication_TOTAL_FAILED: the cryptographic checks of the signature failed
	// (including checks of hashes of individual data objects that have been signed
	// indirectly), or it is proven that the signing certificate was invalid at the
	// time of generation of the signature, or because the signature is not conformant
	// to one of the base standards to the extent that the cryptographic verification
	// building block is unable to process it.
	Indication_TOTAL_FAILED Indication = "TOTAL_FAILED"
	// Indication_INDETERMINATE: the results of the performed checks do not allow to
	// ascertain the signature to be TOTAL-PASSED or TOTAL-FAILED.
	Indication_INDETERMINATE Indication = "INDETERMINATE"
	// Indication_PASSED: when an individual constrain validation succeeds.
	Indication_PASSED Indication = "PASSED"
	// Indication_FAILED: when an individual constrain validation fails.
	Indication_FAILED Indication = "FAILED"
	// Indication_NO_SIGNATURE_FOUND: when no signature is found within the document
	// (empty report is not permitted).
	Indication_NO_SIGNATURE_FOUND Indication = "NO_SIGNATURE_FOUND"
)

// indicationURI maps each Indication to its VR URI.
var indicationURI = map[Indication]string{
	Indication_TOTAL_PASSED:       "urn:etsi:019102:mainindication:total-passed",
	Indication_TOTAL_FAILED:       "urn:etsi:019102:mainindication:total-failed",
	Indication_INDETERMINATE:      "urn:etsi:019102:mainindication:indeterminate",
	Indication_PASSED:             "urn:etsi:019102:mainindication:passed",
	Indication_FAILED:             "urn:etsi:019102:mainindication:failed",
	Indication_NO_SIGNATURE_FOUND: "urn:cef:dss:mainindication:noSignatureFound",
}

// IndicationValues returns all Indication constants in declaration order.
func IndicationValues() []Indication {
	return []Indication{
		Indication_TOTAL_PASSED,
		Indication_TOTAL_FAILED,
		Indication_INDETERMINATE,
		Indication_PASSED,
		Indication_FAILED,
		Indication_NO_SIGNATURE_FOUND,
	}
}

// IndicationValueOf returns the Indication matching the given Java enum name.
func IndicationValueOf(name string) (Indication, error) {
	for _, v := range IndicationValues() {
		if string(v) == name {
			return v, nil
		}
	}
	return "", fmt.Errorf("no enum constant Indication.%s", name)
}

// URI returns the indication VR URI. Implements UriBasedEnum.
func (i Indication) URI() string {
	return indicationURI[i]
}
