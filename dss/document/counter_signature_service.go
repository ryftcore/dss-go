// Ported from dss-document/src/main/java/eu/europa/esig/dss/signature/CounterSignatureService.java (DSS 6.5.RC1).
package document

import (
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi/validation"
)

// CounterSignatureService provides operations for a counter-signature creation, generic over the
// CSP implementation of certain format signature parameters.
type CounterSignatureService[CSP model.SerializableCounterSignatureParameters] interface {
	// GetDataToBeCounterSigned retrieves the bytes of the data that need to be counter-signed
	// from signatureDocument. signatureDocument shall be a valid signature of the same type.
	// Port of #getDataToBeCounterSigned.
	GetDataToBeCounterSigned(signatureDocument model.DSSDocument, parameters CSP) *model.ToBeSigned

	// CounterSignSignature counter-signs the signatureDocument with the provided signatureValue.
	// Port of #counterSignSignature.
	CounterSignSignature(signatureDocument model.DSSDocument, parameters CSP, signatureValue *model.SignatureValue) model.DSSDocument

	// SetTspSource defines the TSP (timestamp provider) source, used when timestamping the
	// signature. Port of #setTspSource.
	SetTspSource(tspSource validation.TSPSource)
}
