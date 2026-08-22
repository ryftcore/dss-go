// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/signature/XAdESLevelBaselineB.java (DSS 6.5.RC1).
//
// Java's static initializer calls SantuarioInitializer.init(); the Go port keeps that call in
// the package's init-time entry point (Service), where the equivalent block also lives,
// and repeats it in the constructor here so an LevelBaselineB built without going through
// the service still runs against an initialized stack - exactly the reason upstream duplicates
// the static block in both classes.
//
// Java's two getDataToSign and two signDocument overloads cannot share one Go name each:
//
//	getDataToSign(DSSDocument, SignatureParameters)         -> GetDataToSign
//	getDataToSign(List<DSSDocument>, SignatureParameters)   -> GetDataToSignForDocuments
//	signDocument(DSSDocument, ..., byte[])                       -> SignDocument
//	signDocument(List<DSSDocument>, ..., byte[])                 -> SignDocuments
//
// The last two are the SignatureProfile interface methods.
package xades

import (
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi/validation"
	xmlutils "github.com/ryftcore/dss-go/dss/xml/utils"
)

// XAdESSignatureBuilderRef is the polymorphic slice of Java's abstract XAdESSignatureBuilder
// that LevelBaselineB's local variables are typed with: getDataToSign holds an
// AbstractSignatureBuilder and calls build(), signDocument holds a SignatureBuilder and calls
// signDocument(byte[]). The concrete builders (Enveloped/Enveloping/Detached/InternallyDetached)
// are distinct Go types embedding AbstractSignatureBuilder, so the factory needs an interface
// value; every one of them satisfies this through that embedded base. Same technique, and same
// reason, as cades.SignatureExtender (see cades/cades_service.go).
type SignatureBuilderRef interface {
	SignatureBuilder

	// Build returns the canonicalized ds:SignedInfo octets to be signed.
	// Port of XAdESSignatureBuilder#build().
	Build() ([]byte, error)
}

// LevelBaselineB contains the B level baseline profile for a XAdES signature.
type LevelBaselineB struct {
	// certificateVerifier provides information on the sources to be used in the validation
	// process in the context of a signature.
	certificateVerifier validation.CertificateVerifier
}

// NewLevelBaselineB is the default constructor for LevelBaselineB.
// Port of XAdESLevelBaselineB(CertificateVerifier).
func NewLevelBaselineB(certificateVerifier validation.CertificateVerifier) *LevelBaselineB {
	xmlutils.SantuarioInitializerInit()
	return &LevelBaselineB{certificateVerifier: certificateVerifier}
}

// GetDataToSign returns the canonicalized ds:SignedInfo XML segment. This method is used for
// signing a document. Port of the #getDataToSign(DSSDocument, XAdESSignatureParameters)
// overload; Java's DSSException becomes the returned error.
func (b *LevelBaselineB) GetDataToSign(dssDocument model.DSSDocument,
	parameters *SignatureParameters) ([]byte, error) {
	return b.GetDataToSignForDocuments([]model.DSSDocument{dssDocument}, parameters)
}

// GetDataToSignForDocuments returns the canonicalized ds:SignedInfo XML segment. This method is
// used for signing multiple documents.
// Port of the #getDataToSign(List<DSSDocument>, XAdESSignatureParameters) overload.
func (b *LevelBaselineB) GetDataToSignForDocuments(documents []model.DSSDocument,
	parameters *SignatureParameters) ([]byte, error) {
	signatureBuilder, err := SignatureBuilderGetSignatureBuilderForDocuments(parameters, documents,
		b.certificateVerifier)
	if err != nil {
		return nil, err
	}
	parameters.GetContext().SetBuilder(signatureBuilder)
	return signatureBuilder.Build()
}

// SignDocument adds the signature value to the signature.
// Port of the overridden #signDocument(DSSDocument, XAdESSignatureParameters, byte[]).
func (b *LevelBaselineB) SignDocument(document model.DSSDocument,
	parameters *SignatureParameters, signatureValue []byte) (model.DSSDocument, error) {
	return b.SignDocuments([]model.DSSDocument{document}, parameters, signatureValue)
}

// SignDocuments adds the signature value to the signature over multiple documents.
// Port of the overridden #signDocument(List<DSSDocument>, XAdESSignatureParameters, byte[]).
func (b *LevelBaselineB) SignDocuments(toSignDocuments []model.DSSDocument,
	parameters *SignatureParameters, signatureValue []byte) (model.DSSDocument, error) {
	builder := parameters.GetContext().Builder()
	if builder != nil {
		// Upstream re-reads the very same value here; kept verbatim as a no-op.
		builder = parameters.GetContext().Builder()
	} else {
		newBuilder, err := SignatureBuilderGetSignatureBuilderForDocuments(parameters, toSignDocuments,
			b.certificateVerifier)
		if err != nil {
			return nil, err
		}
		builder = newBuilder
	}
	dssDocument, err := builder.SignDocument(signatureValue)
	if err != nil {
		return nil, err
	}
	parameters.GetContext().SetBuilder(builder)
	return dssDocument, nil
}
