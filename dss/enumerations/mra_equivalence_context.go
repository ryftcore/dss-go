// Ported from dss-enumerations/.../MRAEquivalenceContext.java (DSS 6.5.RC1).

package enumerations

import "fmt"

// MRAEquivalenceContext identifies the context of the machine processable
// declarative statement whose reference implementation(s) used by the pointing
// contracting party is (are) declared in the CertificateContentDeclarationPointingParty
// element and whose equivalent implementation(s) used by the pointed contracting
// party is (are) declared in the CertificateContentDeclarationPointedParty element.
//
// Implements UriBasedEnum.
type MRAEquivalenceContext string

const (
	// MRAEquivalenceContextQCCompliance indicates that the
	// CertificateContentReferenceEquivalence element applies to the context of mapping
	// the respective pointing party and pointed party reference machine processable
	// statement(s) included in a certificate to declare (as a statement made by the
	// issuing TSP) and to confirm (as a benchmark for establishing the content of the
	// corresponding TL trust service entry) that it has been issued as a qualified
	// certificate.
	MRAEquivalenceContextQCCompliance MRAEquivalenceContext = "QC_COMPLIANCE"
	// MRAEquivalenceContextQCType indicates that the
	// CertificateContentReferenceEquivalence element applies to the context of mapping
	// the respective pointing party and pointed party reference machine processable
	// statement(s) included in a certificate to declare (as a statement made by the
	// issuing TSP) and to confirm (as a benchmark for establishing the content of the
	// corresponding TL trust service entry) that it has been issued for a certain
	// usage type (i.e. for electronic signatures, for electronic seals, or for
	// website authentication).
	MRAEquivalenceContextQCType MRAEquivalenceContext = "QC_TYPE"
	// MRAEquivalenceContextQCQSCD indicates that the
	// CertificateContentReferenceEquivalence element applies to the context of mapping
	// the respective pointing party and pointed party reference machine processable
	// statement(s) included in a certificate to declare (as a statement made by the
	// issuing TSP) and to confirm (as a benchmark for establishing the content of the
	// corresponding TL trust service entry) that the private key, to which the
	// certified public key corresponds, resides in an EU qualified electronic
	// signature or seal creation device.
	MRAEquivalenceContextQCQSCD MRAEquivalenceContext = "QC_QSCD"
)

// mraEquivalenceContextURI maps each MRAEquivalenceContext to its defined URI.
var mraEquivalenceContextURI = map[MRAEquivalenceContext]string{
	MRAEquivalenceContextQCCompliance: "http://ec.europa.eu/tools/lotl/mra/QcCompliance",
	MRAEquivalenceContextQCType:       "http://ec.europa.eu/tools/lotl/mra/QcType",
	MRAEquivalenceContextQCQSCD:       "http://ec.europa.eu/tools/lotl/mra/QcQSCD",
}

// MRAEquivalenceContextValues returns all MRAEquivalenceContext constants in declaration order.
func MRAEquivalenceContextValues() []MRAEquivalenceContext {
	return []MRAEquivalenceContext{
		MRAEquivalenceContextQCCompliance,
		MRAEquivalenceContextQCType,
		MRAEquivalenceContextQCQSCD,
	}
}

// MRAEquivalenceContextValueOf returns the MRAEquivalenceContext matching the given
// Java enum name.
func MRAEquivalenceContextValueOf(name string) (MRAEquivalenceContext, error) {
	for _, v := range MRAEquivalenceContextValues() {
		if string(v) == name {
			return v, nil
		}
	}
	return "", fmt.Errorf("no enum constant MRAEquivalenceContext.%s", name)
}

// URI identifies the URI of the MRA equivalence context. Implements UriBasedEnum.
func (m MRAEquivalenceContext) URI() string {
	return mraEquivalenceContextURI[m]
}
