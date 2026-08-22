// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/validation/scope/SignatureScopeFinder.java (DSS 6.5.RC1).
package scope

import (
	"github.com/ryftcore/dss-go/dss/model/scope"
	"github.com/ryftcore/dss-go/dss/spi/validation"
)

// SignatureScopeFinder builds a list of SignatureScopes from an AdvancedSignature.
//
// T is the AdvancedSignature implementation, mirroring Java's
// SignatureScopeFinder<T extends AdvancedSignature>.
type SignatureScopeFinder[T validation.AdvancedSignature] interface {
	// FindSignatureScope returns a list of SignatureScopes from a signature. Port of
	// findSignatureScope(T).
	FindSignatureScope(advancedSignature T) []scope.SignatureScope
}
