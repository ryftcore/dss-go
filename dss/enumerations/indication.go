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
	// IndicationTotalPassed: when the cryptographic checks of the signature
	// (including checks of hashes of individual data objects that have been signed
	// indirectly) succeeded as well as all checks prescribed by the signature
	// validation policy have been passed.
	IndicationTotalPassed Indication = "TOTAL_PASSED"
	// IndicationTotalFailed: the cryptographic checks of the signature failed
	// (including checks of hashes of individual data objects that have been signed
	// indirectly), or it is proven that the signing certificate was invalid at the
	// time of generation of the signature, or because the signature is not conformant
	// to one of the base standards to the extent that the cryptographic verification
	// building block is unable to process it.
	IndicationTotalFailed Indication = "TOTAL_FAILED"
	// IndicationIndeterminate: the results of the performed checks do not allow to
	// ascertain the signature to be TOTAL-PASSED or TOTAL-FAILED.
	IndicationIndeterminate Indication = "INDETERMINATE"
	// IndicationPassed: when an individual constrain validation succeeds.
	IndicationPassed Indication = "PASSED"
	// IndicationFailed: when an individual constrain validation fails.
	IndicationFailed Indication = "FAILED"
	// IndicationNoSignatureFound: when no signature is found within the document
	// (empty report is not permitted).
	IndicationNoSignatureFound Indication = "NO_SIGNATURE_FOUND"
)

// indicationURI maps each Indication to its VR URI.
var indicationURI = map[Indication]string{
	IndicationTotalPassed:      "urn:etsi:019102:mainindication:total-passed",
	IndicationTotalFailed:      "urn:etsi:019102:mainindication:total-failed",
	IndicationIndeterminate:    "urn:etsi:019102:mainindication:indeterminate",
	IndicationPassed:           "urn:etsi:019102:mainindication:passed",
	IndicationFailed:           "urn:etsi:019102:mainindication:failed",
	IndicationNoSignatureFound: "urn:cef:dss:mainindication:noSignatureFound",
}

// IndicationValues returns all Indication constants in declaration order.
func IndicationValues() []Indication {
	return []Indication{
		IndicationTotalPassed,
		IndicationTotalFailed,
		IndicationIndeterminate,
		IndicationPassed,
		IndicationFailed,
		IndicationNoSignatureFound,
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
