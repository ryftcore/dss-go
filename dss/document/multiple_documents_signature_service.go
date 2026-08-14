// Ported from dss-document/src/main/java/eu/europa/esig/dss/signature/MultipleDocumentsSignatureService.java (DSS 6.5.RC1).
package document

import (
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/spi/validation"
)

// MultipleDocumentsSignatureService provides operations for the signature creation and for its
// extension over multiple documents, generic over the SP implementation of signature parameters
// and TP implementation of timestamp parameters corresponding to the supported document format.
type MultipleDocumentsSignatureService[SP model.SerializableSignatureParameters, TP model.SerializableTimestampParameters] interface {
	// GetContentTimestamp creates a content-timestamp attribute (to be included in the
	// signed-data). Port of #getContentTimestamp.
	GetContentTimestamp(toSignDocuments []model.DSSDocument, parameters SP) *validation.TimestampToken

	// GetDataToSign retrieves the bytes of the data that need to be signed based on the
	// toSignDocuments and parameters. When toSignDocuments contains an already existing
	// signature the returned bytes are related to a new parallel signature. Port of
	// #getDataToSign.
	GetDataToSign(toSignDocuments []model.DSSDocument, parameters SP) *model.ToBeSigned

	// IsValidSignatureValue verifies the signature value against a ToBeSigned and a
	// CertificateToken. Port of #isValidSignatureValue.
	IsValidSignatureValue(toBeSigned *model.ToBeSigned, signatureValue *model.SignatureValue, signingCertificate *model.CertificateToken) bool

	// SignDocument signs the toSignDocuments with the provided signatureValue. Port of
	// #signDocument.
	SignDocument(toSignDocuments []model.DSSDocument, parameters SP, signatureValue *model.SignatureValue) model.DSSDocument

	// ExtendDocument extends the level of the signatures in the toExtendDocument. Port of
	// #extendDocument.
	ExtendDocument(toExtendDocument model.DSSDocument, parameters SP) model.DSSDocument

	// Timestamp timestamps the toTimestampDocuments with the provided parameters. Port of
	// #timestamp.
	Timestamp(toTimestampDocuments []model.DSSDocument, parameters TP) model.DSSDocument
}
