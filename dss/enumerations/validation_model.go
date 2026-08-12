// Ported from dss-enumerations/.../ValidationModel.java (DSS 6.5.RC1).
package enumerations

// ValidationModel represents a validation model of a certificate chain
// (e.g. SHELL, CHAIN, etc.).
type ValidationModel string

const (
	// ValidationModel_SHELL is the model for validation of X.509
	// certificate chains where all certificates have to be valid at a
	// given time.
	ValidationModel_SHELL ValidationModel = "SHELL"
	// ValidationModel_CHAIN is the model for validation of X.509
	// certificate chains where all CA certificates have to be valid at the
	// time they were used for issuing a certificate and the end-entity
	// certificate was valid when creating the signature.
	ValidationModel_CHAIN ValidationModel = "CHAIN"
	// ValidationModel_HYBRID is the hybrid validation model, evaluating
	// the signing-certificate at the validation time, while all other
	// intermediate CA certificates at the time of the signing-certificate's
	// issuance.
	ValidationModel_HYBRID ValidationModel = "HYBRID"
)

// ValidationModelValues returns all constants in declaration order.
func ValidationModelValues() []ValidationModel {
	return []ValidationModel{
		ValidationModel_SHELL,
		ValidationModel_CHAIN,
		ValidationModel_HYBRID,
	}
}
