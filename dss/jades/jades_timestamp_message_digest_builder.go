// Ported from dss-jades/src/main/java/eu/europa/esig/dss/jades/validation/timestamp/JAdESTimestampMessageDigestBuilder.java (DSS 6.5.RC1).
//
// SCC flattening: Java's eu.europa.esig.dss.jades.validation.timestamp package folds into this
// one Go package per the phase-6 package layout (S6_BRIEF.md).
//
// FORWARD DEPENDENCY: *JAdESSignature - see abstract_jws_document_analyzer.go's file header. This
// file additionally needs:
//
//	func (s *JAdESSignature) Jws() *JWS                                             // getJws()
//	func (s *JAdESSignature) SigDMechanism() *enumerations.SigDMechanism            // getSigDMechanism(), nil when absent
//	func (s *JAdESSignature) SignedDocumentsByHTTPHeaderName() []model.DSSDocument  // getSignedDocumentsByHTTPHeaderName()
//	func (s *JAdESSignature) SignedDocumentsForObjectIdByUriMechanism() []model.DSSDocument // getSignedDocumentsForObjectIdByUriMechanism()
//	func (s *JAdESSignature) EtsiUHeader() *JAdESEtsiUHeader                        // getEtsiUHeader() (already used by other landed jades files, e.g. jades_level_baseline_lt.go)
//
// # Canonicalization
//
// getCanonicalizedValue (Java) leaves JSON canonicalization unimplemented ("TODO: canonicalization
// is not supported yet"), falling back to plain JSON serialization; this port reproduces that
// exact deviation (jadesTMDBCanonicalizedValue below), which is a documented upstream limitation,
// not a porting shortcut.
package jades

import (
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/internal/jose"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/spi/validation"
)

// JAdESTimestampMessageDigestBuilder builds the message-imprint digest for JAdES timestamps.
// Port of the class JAdESTimestampMessageDigestBuilder, implementing
// timestamp.TimestampMessageDigestBuilder.
type JAdESTimestampMessageDigestBuilder struct {
	// signature is the signature.
	signature *JAdESSignature

	// digestAlgorithm is the digest algorithm to be used for message-imprint digest computation.
	digestAlgorithm enumerations.DigestAlgorithm

	// timestampToken is the timestamp token to compute message-digest for.
	timestampToken *validation.TimestampToken

	// canonicalizationAlgorithm is the canonicalization algorithm to be used for message-imprint
	// computation.
	canonicalizationAlgorithm string

	// timestampAttribute is the signature element containing the time-stamp token.
	timestampAttribute *JAdESAttribute
}

// newJAdESTimestampMessageDigestBuilder is the port of the private (JAdESSignature) constructor.
//
// Panics with the Java message when signature is nil (Objects.requireNonNull).
func newJAdESTimestampMessageDigestBuilder(signature *JAdESSignature) *JAdESTimestampMessageDigestBuilder {
	if signature == nil {
		panic("Signature cannot be null!")
	}
	return &JAdESTimestampMessageDigestBuilder{signature: signature}
}

// NewJAdESTimestampMessageDigestBuilder is the constructor to compute message-imprint for
// timestamps related to signature, to be used on timestamp creation. Port of the (JAdESSignature,
// DigestAlgorithm) constructor.
//
// Panics with the Java message when digestAlgorithm is empty (Objects.requireNonNull).
func NewJAdESTimestampMessageDigestBuilder(signature *JAdESSignature, digestAlgorithm enumerations.DigestAlgorithm) *JAdESTimestampMessageDigestBuilder {
	b := newJAdESTimestampMessageDigestBuilder(signature)
	if digestAlgorithm == "" {
		panic("DigestAlgorithm cannot be null!")
	}
	b.digestAlgorithm = digestAlgorithm
	return b
}

// NewJAdESTimestampMessageDigestBuilderForToken is the constructor to compute message-imprint for
// timestamps related to signature, containing timestamps. Port of the (JAdESSignature,
// TimestampToken) constructor.
//
// Panics with the Java message when timestampToken is nil (Objects.requireNonNull).
func NewJAdESTimestampMessageDigestBuilderForToken(signature *JAdESSignature, timestampToken *validation.TimestampToken) *JAdESTimestampMessageDigestBuilder {
	b := newJAdESTimestampMessageDigestBuilder(signature)
	if timestampToken == nil {
		panic("TimestampToken cannot be null!")
	}
	b.timestampToken = timestampToken
	b.digestAlgorithm = timestampToken.DigestAlgorithm()
	b.canonicalizationAlgorithm = timestampToken.CanonicalizationMethod()
	return b
}

// SetCanonicalizationAlgorithm sets the canonicalization algorithm to be used for message-digest
// computation. Port of setCanonicalizationAlgorithm(String).
func (b *JAdESTimestampMessageDigestBuilder) SetCanonicalizationAlgorithm(canonicalizationAlgorithm string) *JAdESTimestampMessageDigestBuilder {
	b.canonicalizationAlgorithm = canonicalizationAlgorithm
	return b
}

// SetTimestampAttribute sets a signature attribute identifying the time-stamp token. Port of
// setTimestampAttribute(JAdESAttribute).
func (b *JAdESTimestampMessageDigestBuilder) SetTimestampAttribute(timestampAttribute *JAdESAttribute) *JAdESTimestampMessageDigestBuilder {
	b.timestampAttribute = timestampAttribute
	return b
}

// ContentTimestampMessageDigest implements timestamp.TimestampMessageDigestBuilder. Port of
// getContentTimestampMessageDigest().
func (b *JAdESTimestampMessageDigestBuilder) ContentTimestampMessageDigest() model.DSSMessageDigest {
	messageDigest, err := b.contentTimestampMessageDigest()
	if err != nil {
		// Upstream logs MESSAGE_IMPRINT_ERROR[_WITH_ID].
		return model.CreateEmptyDSSMessageDigest()
	}
	return messageDigest
}

// contentTimestampMessageDigest is the fallible core of ContentTimestampMessageDigest. Port of
// the try block of getContentTimestampMessageDigest().
func (b *JAdESTimestampMessageDigestBuilder) contentTimestampMessageDigest() (model.DSSMessageDigest, error) {
	digestCalculator, err := spi.NewDSSMessageDigestCalculator(b.digestAlgorithm)
	if err != nil {
		return model.DSSMessageDigest{}, err
	}
	if err := b.writeSignedDataBinaries(digestCalculator); err != nil {
		return model.DSSMessageDigest{}, err
	}
	return digestCalculator.MessageDigest(b.digestAlgorithm), nil
}

// writeSignedDataBinaries ports the private writeSignedDataBinaries(DSSMessageDigestCalculator).
func (b *JAdESTimestampMessageDigestBuilder) writeSignedDataBinaries(digestCalculator *spi.DSSMessageDigestCalculator) error {
	if sigDMechanism := b.signature.SigDMechanism(); sigDMechanism != nil {
		return b.writeSigDReferencedOctets(digestCalculator, *sigDMechanism)
	}
	return b.writeJWSPayloadValue(digestCalculator)
}

// writeJWSPayloadValue ports the private writeJWSPayloadValue(DSSMessageDigestCalculator).
func (b *JAdESTimestampMessageDigestBuilder) writeJWSPayloadValue(digestCalculator *spi.DSSMessageDigestCalculator) error {
	var payload []byte
	if b.signature.Jws().IsRfc7797UnencodedPayload() {
		payload = b.signature.Jws().UnverifiedPayloadBytes()
	} else {
		payload = []byte(b.signature.Jws().EncodedPayload())
	}
	if len(payload) == 0 {
		return model.NewDSSError("Unable to extract JWS payload!")
	}
	digestCalculator.Update(payload)
	return nil
}

// writeSigDReferencedOctets ports the private writeSigDReferencedOctets(DSSMessageDigestCalculator,
// SigDMechanism).
func (b *JAdESTimestampMessageDigestBuilder) writeSigDReferencedOctets(digestCalculator *spi.DSSMessageDigestCalculator,
	sigDMechanism enumerations.SigDMechanism) error {
	switch sigDMechanism {
	case enumerations.SigDMechanism_HTTP_HEADERS:
		documentList := b.signature.SignedDocumentsByHTTPHeaderName()
		httpHeadersPayloadBuilder := NewHttpHeadersPayloadBuilder(documentList, true)
		sigDOctets, err := httpHeadersPayloadBuilder.Build()
		if err != nil {
			return err
		}
		digestCalculator.Update(sigDOctets)
		return nil

	case enumerations.SigDMechanism_OBJECT_ID_BY_URI, enumerations.SigDMechanism_OBJECT_ID_BY_URI_HASH:
		documentList := b.signature.SignedDocumentsForObjectIdByUriMechanism()
		return DSSJsonUtilsWriteDocumentsDigest(documentList, !b.signature.Jws().IsRfc7797UnencodedPayload(), digestCalculator)

	default:
		// Upstream logs "Unsupported SigDMechanism '{}' has been found!".
		return nil
	}
}

// SignatureTimestampMessageDigest implements timestamp.TimestampMessageDigestBuilder. Port of
// getSignatureTimestampMessageDigest().
func (b *JAdESTimestampMessageDigestBuilder) SignatureTimestampMessageDigest() model.DSSMessageDigest {
	messageDigest, err := b.signatureTimestampMessageDigest()
	if err != nil {
		// Upstream logs MESSAGE_IMPRINT_ERROR[_WITH_ID].
		return model.CreateEmptyDSSMessageDigest()
	}
	return messageDigest
}

// signatureTimestampMessageDigest is the fallible core of SignatureTimestampMessageDigest. Port
// of the try block of getSignatureTimestampMessageDigest().
//
// 5.3.4 The sigTst JSON object: the input of the message imprint computation for the time-stamp
// tokens encapsulated by sigTst JSON object shall be the base64url-encoded JWS Signature Value.
func (b *JAdESTimestampMessageDigestBuilder) signatureTimestampMessageDigest() (model.DSSMessageDigest, error) {
	signatureTimestampData := b.base64UrlEncodedSignatureValue()
	digest, err := spi.DSSUtilsDigest(b.digestAlgorithm, signatureTimestampData)
	if err != nil {
		return model.DSSMessageDigest{}, err
	}
	return model.NewDSSMessageDigestWithValue(b.digestAlgorithm, digest), nil
}

// TimestampX1MessageDigest implements timestamp.TimestampMessageDigestBuilder. Port of
// getTimestampX1MessageDigest().
func (b *JAdESTimestampMessageDigestBuilder) TimestampX1MessageDigest() model.DSSMessageDigest {
	messageDigest, err := b.timestampX1MessageDigest()
	if err != nil {
		// Upstream logs MESSAGE_IMPRINT_ERROR[_WITH_ID].
		return model.CreateEmptyDSSMessageDigest()
	}
	return messageDigest
}

// timestampX1MessageDigest is the fallible core of TimestampX1MessageDigest. Port of the try
// block of getTimestampX1MessageDigest() (A.1.5.1 The sigRTst JSON object).
func (b *JAdESTimestampMessageDigestBuilder) timestampX1MessageDigest() (model.DSSMessageDigest, error) {
	digestCalculator, err := spi.NewDSSMessageDigestCalculator(b.digestAlgorithm)
	if err != nil {
		return model.DSSMessageDigest{}, err
	}

	// 1) The value of the base64url-encoded JWS Signature Value.
	digestCalculator.Update(b.base64UrlEncodedSignatureValue())

	// 2) The character '.'.
	digestCalculator.UpdateByte('.')

	// 3) Those among sigTst/xRefs/rRefs/axRefs/arRefs that appear before sigRTst, in their order
	// of appearance within the etsiU array, base64url-encoded. NOTE: there is a difference in
	// processing base64url encoded values and clear incorporation.
	jws := b.signature.Jws()
	etsiU := DSSJsonUtilsEtsiU(jws)
	if DSSJsonUtilsCheckComponentsUnicity(etsiU) {
		etsiUHeader := b.signature.EtsiUHeader()
		for _, etsiUComponent := range etsiUHeader.Attributes() {
			if b.timestampAttribute != nil && b.timestampAttribute.Equals(&etsiUComponent.JAdESAttribute) {
				// the current timestamp is found, stop the iteration
				break
			}
			if jadesTMDBIsAllowedTypeEntry(etsiUComponent,
				JAdESHeaderParameterNamesSigTst, JAdESHeaderParameterNamesXRefs,
				JAdESHeaderParameterNamesRRefs, JAdESHeaderParameterNamesAxRefs,
				JAdESHeaderParameterNamesArRefs) {
				digestCalculator.Update(b.etsiUComponentValue(etsiUComponent, b.canonicalizationAlgorithm))
			}
		}
	} else {
		// Upstream logs "Unable to process 'etsiU' entries for a '{}' timestamp. The 'etsiU'
		// components shall have a common format (Strings or Objects)!".
	}

	return digestCalculator.MessageDigest(b.digestAlgorithm), nil
}

// jadesTMDBIsAllowedTypeEntry ports the private isAllowedTypeEntry(EtsiUComponent, String...).
func jadesTMDBIsAllowedTypeEntry(etsiUComponent *EtsiUComponent, allowedTypes ...string) bool {
	for _, allowedType := range allowedTypes {
		if allowedType == etsiUComponent.HeaderName() {
			return true
		}
	}
	return false
}

// TimestampX2MessageDigest implements timestamp.TimestampMessageDigestBuilder. Port of
// getTimestampX2MessageDigest().
func (b *JAdESTimestampMessageDigestBuilder) TimestampX2MessageDigest() model.DSSMessageDigest {
	messageDigest, err := b.timestampX2MessageDigest()
	if err != nil {
		// Upstream logs MESSAGE_IMPRINT_ERROR[_WITH_ID].
		return model.CreateEmptyDSSMessageDigest()
	}
	return messageDigest
}

// timestampX2MessageDigest is the fallible core of TimestampX2MessageDigest. Port of the try
// block of getTimestampX2MessageDigest() (A.1.5.2.2 The rfsTst JSON object).
func (b *JAdESTimestampMessageDigestBuilder) timestampX2MessageDigest() (model.DSSMessageDigest, error) {
	digestCalculator, err := spi.NewDSSMessageDigestCalculator(b.digestAlgorithm)
	if err != nil {
		return model.DSSMessageDigest{}, err
	}

	// The message imprint computation input shall be the concatenation of xRefs/rRefs/axRefs/
	// arRefs, base64url encoded, in their order of appearance within the etsiU array. NOTE: there
	// is a difference in processing base64url encoded values and clear incorporation.
	jws := b.signature.Jws()
	etsiU := DSSJsonUtilsEtsiU(jws)
	if DSSJsonUtilsCheckComponentsUnicity(etsiU) {
		etsiUHeader := b.signature.EtsiUHeader()
		for _, etsiUComponent := range etsiUHeader.Attributes() {
			if jadesTMDBIsAllowedTypeEntry(etsiUComponent,
				JAdESHeaderParameterNamesXRefs, JAdESHeaderParameterNamesRRefs,
				JAdESHeaderParameterNamesAxRefs, JAdESHeaderParameterNamesArRefs) {
				digestCalculator.Update(b.etsiUComponentValue(etsiUComponent, b.canonicalizationAlgorithm))
			}
		}
	} else {
		// Upstream logs "Unable to process 'etsiU' entries for an '{}' timestamp. ...".
	}

	return digestCalculator.MessageDigest(b.digestAlgorithm), nil
}

// ArchiveTimestampMessageDigest implements timestamp.TimestampMessageDigestBuilder. Port of
// getArchiveTimestampMessageDigest() (5.3.6.3.1 Processing).
func (b *JAdESTimestampMessageDigestBuilder) ArchiveTimestampMessageDigest() model.DSSMessageDigest {
	messageDigest, err := b.archiveTimestampMessageDigest()
	if err != nil {
		// Upstream logs MESSAGE_IMPRINT_ERROR[_WITH_ID].
		return model.CreateEmptyDSSMessageDigest()
	}
	return messageDigest
}

// archiveTimestampMessageDigest is the fallible core of ArchiveTimestampMessageDigest.
func (b *JAdESTimestampMessageDigestBuilder) archiveTimestampMessageDigest() (model.DSSMessageDigest, error) {
	jws := b.signature.Jws()

	digestCalculator, err := spi.NewDSSMessageDigestCalculator(b.digestAlgorithm)
	if err != nil {
		return model.DSSMessageDigest{}, err
	}

	// 1)/2) sigD-mechanism-aware JWS-payload-or-HttpHeaders octets, per writeSignedDataBinaries.
	if err := b.writeSignedDataBinaries(digestCalculator); err != nil {
		return model.DSSMessageDigest{}, err
	}

	// 3) The character '.'.
	digestCalculator.UpdateByte('.')

	// 4) The value of the JWS Protected Header, base64url encoded, followed by '.'.
	digestCalculator.Update([]byte(jws.EncodedHeader()))
	digestCalculator.UpdateByte('.')

	// 5) The value of the JAdES Signature Value, base64url encoded.
	digestCalculator.Update(b.base64UrlEncodedSignatureValue())

	// 6) The character '.'.
	digestCalculator.UpdateByte('.')

	// 7) Every etsiU array element up to (excluding) the timestamp itself, base64url encoded (or
	// canonicalized, when the etsiU array carries clear instances).
	etsiU := DSSJsonUtilsEtsiU(jws)
	if DSSJsonUtilsCheckComponentsUnicity(etsiU) {
		etsiUHeader := b.signature.EtsiUHeader()
		for _, etsiUComponent := range etsiUHeader.Attributes() {
			if b.timestampAttribute != nil && b.timestampAttribute.Equals(&etsiUComponent.JAdESAttribute) {
				// the timestamp is reached, stop the iteration
				break
			}
			digestCalculator.Update(b.etsiUComponentValue(etsiUComponent, b.canonicalizationAlgorithm))
		}
	} else {
		// Upstream logs "Unable to process 'etsiU' entries for an '{}' timestamp. ...".
	}

	return digestCalculator.MessageDigest(b.digestAlgorithm), nil
}

// base64UrlEncodedSignatureValue ports the private getBase64UrlEncodedSignatureValue().
func (b *JAdESTimestampMessageDigestBuilder) base64UrlEncodedSignatureValue() []byte {
	messageImprint := b.signature.Jws().EncodedSignature()
	return []byte(messageImprint)
}

// etsiUComponentValue ports the private getEtsiUComponentValue(EtsiUComponent, String).
func (b *JAdESTimestampMessageDigestBuilder) etsiUComponentValue(etsiUComponent *EtsiUComponent, canonicalizationMethod string) []byte {
	component := etsiUComponent.Component()
	if etsiUComponent.IsBase64UrlEncoded() {
		str, _ := component.(string)
		return []byte(str)
	}
	return jadesTMDBCanonicalizedValue(etsiUComponent.Value(), canonicalizationMethod)
}

// jadesTMDBCanonicalizedValue ports the private getCanonicalizedValue(Object, String). See the
// file header's "Canonicalization" note: this is a documented upstream limitation (canonicalization
// is not supported yet), reproduced verbatim rather than fixed.
func jadesTMDBCanonicalizedValue(jsonObject any, canonicalizationMethod string) []byte {
	// Upstream logs "Canonicalization is not supported in the current version. The message
	// imprint computation can lead to an unexpected result".
	// temporary solution
	return []byte(jose.JSON(jsonObject))
}
