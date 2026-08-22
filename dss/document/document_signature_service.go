// Ported from dss-document/src/main/java/eu/europa/esig/dss/signature/DocumentSignatureService.java (DSS 6.5.RC1).
package document

import (
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi/validation"
)

// DocumentSignatureService provides operations for the signature creation and for its extension,
// generic over the SP implementation of signature parameters and TP implementation of timestamp
// parameters corresponding to the supported signature/document format.
type DocumentSignatureService[SP model.SerializableSignatureParameters, TP model.SerializableTimestampParameters] interface {
	// GetDataToSign retrieves the bytes of the data that need to be signed based on the
	// toSignDocument and parameters. When toSignDocument contains an already existing signature
	// the returned bytes are related to a new parallel signature.
	//
	//   - Enveloped signature (XML): a new signature is added and the signed data corresponds to
	//     that pointed by the first signature;
	//   - Enveloping signature:
	//   - XML: The parallel signature is not possible
	//   - CMS: A new parallel signature is added
	//   - Detached signature:
	//   - XML: The parallel signature is added
	//   - CMS: A new parallel signature is added
	//
	// Port of #getDataToSign.
	GetDataToSign(toSignDocument model.DSSDocument, parameters SP) *model.ToBeSigned

	// IsValidSignatureValue verifies the signature value against a ToBeSigned and a
	// CertificateToken. Port of #isValidSignatureValue.
	IsValidSignatureValue(toBeSigned *model.ToBeSigned, signatureValue *model.SignatureValue, signingCertificate *model.CertificateToken) bool

	// SignDocument signs the toSignDocument with the provided signatureValue. Port of
	// #signDocument.
	SignDocument(toSignDocument model.DSSDocument, parameters SP, signatureValue *model.SignatureValue) model.DSSDocument

	// ExtendDocument extends the level of the signatures in the toExtendDocument. Port of
	// #extendDocument.
	ExtendDocument(toExtendDocument model.DSSDocument, parameters SP) model.DSSDocument

	// SetTspSource defines the TSP (timestamp provider) source, used when timestamping the
	// signature. Port of #setTspSource.
	SetTspSource(tspSource validation.TSPSource)

	// GetContentTimestamp allows computing a content-timestamp (which is added in the signed
	// properties). Port of #getContentTimestamp.
	GetContentTimestamp(toSignDocument model.DSSDocument, parameters SP) *validation.TimestampToken

	// Timestamp allows adding a timestamp to an unsigned document. Port of #timestamp.
	Timestamp(toTimestampDocument model.DSSDocument, parameters TP) model.DSSDocument
}
