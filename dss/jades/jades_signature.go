// Ported from dss-jades/src/main/java/eu/europa/esig/dss/jades/validation/JAdESSignature.java
// (DSS 6.5.RC1).
//
// EtsiUHeader/EtsiUComponent, TimestampSource and SignatureScopeFinder consume
// *Signature's API as implemented below. Three of this file's own methods deliberately
// deviate from conventions used elsewhere in this port to match what those siblings expect:
// Jws() (not JWS()), SigDMechanism() returning *enumerations.SigDMechanism (a pointer, not the
// plain-value+"" sentinel convention used elsewhere), and OriginalDocuments() returning
// ([]model.DSSDocument, error) (Java's unchecked exception surfaced as an error here rather than
// a panic).
//
// slf4j logging is dropped; every LOG.warn/LOG.debug call site is called out in the surrounding
// comment instead, except for one try/catch with no realistic Go panic source (see
// getSignaturePolicyStore's doc comment).
package jades

import (
	"bytes"
	"fmt"
	"time"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/internal/jose"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/model/scope"
	"github.com/ryftcore/dss-go/dss/model/signature"
	"github.com/ryftcore/dss-go/dss/model/x509/revocation"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/spi/validation"
	"github.com/ryftcore/dss-go/dss/utils"
)

// JAdESSignature represents the JAdES signature. Port of the class JAdESSignature, extending
// validation.DefaultAdvancedSignature.
//
// serialVersionUID and java.io.Serializable are dropped (no Go counterpart).
type Signature struct {
	validation.DefaultAdvancedSignature

	// jws is the JWS signature object. Port of the private final JWS jws field.
	jws *JWS

	// isDetached defines if the validating signature is detached. Port of the private final
	// boolean isDetached field.
	isDetached bool

	// masterCSigComponent is the 'cSig' object embedding the current signature.
	// NOTE: used for counter signatures only. Port of the private EtsiUComponent
	// masterCSigComponent field.
	masterCSigComponent *EtsiUComponent

	// etsiUHeader is the list of unsigned properties embedded into the 'etsiU' array. Port of
	// the private EtsiUHeader etsiUHeader field.
	etsiUHeader *EtsiUHeader

	// jadesCachedCryptoVerification is this port's idempotency guard for CheckSignatureIntegrity,
	// a replacement for reading the (inaccessible, cross-package-private)
	// signatureCryptographicVerification field directly the way Java's own method does for its
	// "already computed" early-return check. It is set to the same
	// *signature.CryptographicVerification value handed to
	// SetSignatureCryptographicVerification, so the two never disagree - the same pattern
	// cades_signature.go's cachedCryptoVerification field documents.
	jadesCachedCryptoVerification *signature.CryptographicVerification
}

// NewJAdESSignature is the default constructor. Port of the public JAdESSignature(JWS)
// constructor.
func NewJAdESSignature(jws *JWS) *Signature {
	s := &Signature{
		DefaultAdvancedSignature: validation.NewDefaultAdvancedSignatureBase(),
		jws:                      jws,
		isDetached:               utils.IsArrayEmpty(jws.UnverifiedPayloadBytes()),
	}
	s.InitDefaultAdvancedSignature(s)
	return s
}

// Jws gets the associated JWS. Port of the public JWS getJws().
//
// Named Jws (not JWS) to match the surface other, already-landed sibling files in this package
// depend on (e.g. jades_timestamp_source.go, jades_timestamp_message_digest_builder.go,
// jades_diagnostic_data_builder.go, jades_counter_signature_builder.go).
func (s *Signature) Jws() *JWS {
	return s.jws
}

// SignatureForm specifies the format of the signature. Port of getSignatureForm().
func (s *Signature) SignatureForm() enumerations.SignatureForm {
	return enumerations.SignatureFormJAdES
}

// SignatureAlgorithm retrieves the signature algorithm (or cipher) used for generating the
// signature. Port of getSignatureAlgorithm().
func (s *Signature) SignatureAlgorithm() enumerations.SignatureAlgorithm {
	signatureAlgorithm := enumerations.SignatureAlgorithmForJWADefault(s.jws.AlgorithmHeaderValue(), "")
	if signatureAlgorithm == "" {
		// Upstream logs "SignatureAlgorithm '{}' is not supported!".
	} else if enumerations.EncryptionAlgorithmEDDSA == signatureAlgorithm.EncryptionAlgorithm() {
		signatureAlgorithm = spi.DSSUtilsEdDSASignatureAlgorithm(s.SignatureValue())
	}
	return signatureAlgorithm
}

// SigningTime returns the signing time included within the signature, or nil. Port of
// getSigningTime().
func (s *Signature) SigningTime() *time.Time {
	iat := s.jws.ProtectedHeaderValueAsNumber(JWTClaimNamesIat)
	sigT := s.jws.ProtectedHeaderValueAsString(JAdESHeaderParameterNamesSigT)
	if iat != nil && utils.IsStringNotEmpty(sigT) {
		// Upstream logs "Unable to extract claimed signing-time: Conflict between 'iat' and
		// 'sigT' header parameters! Only one shall be present.".
		return nil
	} else if iat != nil {
		numericDate := DSSJsonUtilsToNumericDate(iat)
		return &numericDate
	} else if utils.IsStringNotEmpty(sigT) {
		parsed := spi.DSSUtilsParseRFCDate(sigT)
		return &parsed
	}
	// Upstream logs "Unable to extract claimed signing-time: No signing-time identifying header
	// was found.".
	return nil
}

// ExpirationTime returns a value of 'exp' protected header, when defined. The field is used
// within ETSI TS 119 411-5 TLS Certificate Binding signatures and contains the expiry date of
// the binding. The maximum effective expiry time is whichever is soonest of this field, the
// longest-lived TLS certificate identified in the sigD member payload, or the notAfter time of
// the signing certificate. The value shall be encoded as specified in IETF RFC 7519. Port of the
// public Date getExpirationTime().
func (s *Signature) ExpirationTime() time.Time {
	exp := s.jws.ProtectedHeaderValueAsNumber(JWTClaimNamesExp)
	return DSSJsonUtilsToNumericDate(exp)
}

// IsDetachedSignature checks if the JAdES Signature is detached (payload is not present within
// the signature structure). Port of the public boolean isDetachedSignature().
func (s *Signature) IsDetachedSignature() bool {
	return s.isDetached
}

// MasterCSigComponent gets the 'cSig' component embedding the current signature. Port of the
// public EtsiUComponent getMasterCSigComponent().
func (s *Signature) MasterCSigComponent() *EtsiUComponent {
	return s.masterCSigComponent
}

// SetMasterCSigComponent sets the 'cSig' component embedding the current signature. Port of the
// public void setMasterCSigComponent(EtsiUComponent).
func (s *Signature) SetMasterCSigComponent(masterCSigComponent *EtsiUComponent) {
	s.masterCSigComponent = masterCSigComponent
}

// CertificateSource gets a certificate source which contains ALL certificates embedded in the
// signature. Port of getCertificateSource().
func (s *Signature) CertificateSource() *spi.SignatureCertificateSource {
	if s.OfflineCertificateSource() == nil {
		jadesCertificateSource := NewJAdESCertificateSource(s.jws, s.EtsiUHeader())
		s.SetOfflineCertificateSource(&jadesCertificateSource.SignatureCertificateSource)
	}
	return s.OfflineCertificateSource()
}

// CRLSource gets a CRL source which contains ALL CRLs embedded in the signature. Port of
// getCRLSource().
func (s *Signature) CRLSource() spi.OfflineRevocationSource[revocation.CRL] {
	if s.SignatureCRLSource() == nil {
		s.SetSignatureCRLSource(NewJAdESCRLSource(s.EtsiUHeader()))
	}
	return s.SignatureCRLSource()
}

// OCSPSource gets an OCSP source which contains ALL OCSP responses embedded in the signature.
// Port of getOCSPSource().
func (s *Signature) OCSPSource() spi.OfflineRevocationSource[revocation.OCSP] {
	if s.SignatureOCSPSource() == nil {
		s.SetSignatureOCSPSource(NewJAdESOCSPSource(s.EtsiUHeader()))
	}
	return s.SignatureOCSPSource()
}

// TimestampSource gets a Signature Timestamp source which contains ALL timestamps embedded in
// the signature. Port of getTimestampSource(), covariant in Java (returns JAdESTimestampSource);
// Go interface satisfaction needs the exact validation.TimestampSource return type, so JAdES
// callers needing the concrete type assert on the result.
func (s *Signature) TimestampSource() validation.TimestampSource {
	if s.SignatureTimestampSource() == nil {
		s.SetSignatureTimestampSource(NewJAdESTimestampSource(s))
	}
	return s.SignatureTimestampSource()
}

// ProductionPlace returns information about the place where the signature was
// generated. Port of getSignatureProductionPlace().
func (s *Signature) SignatureProductionPlace() *signature.ProductionPlace {
	signaturePlace := s.jws.ProtectedHeaderValueAsMap(JAdESHeaderParameterNamesSigPl)
	if signaturePlace.Size() != 0 {
		result := signature.NewSignatureProductionPlace()
		result.SetCity(DSSJsonUtilsGetAsString(signaturePlace, JAdESHeaderParameterNamesAddressLocality))
		result.SetStreetAddress(DSSJsonUtilsGetAsString(signaturePlace, JAdESHeaderParameterNamesStreetAddress))
		result.SetPostOfficeBoxNumber(DSSJsonUtilsGetAsString(signaturePlace, JAdESHeaderParameterNamesPostOfficeBoxNumber))
		result.SetPostalCode(DSSJsonUtilsGetAsString(signaturePlace, JAdESHeaderParameterNamesPostalCode))
		result.SetStateOrProvince(DSSJsonUtilsGetAsString(signaturePlace, JAdESHeaderParameterNamesAddressRegion))
		result.SetCountryName(DSSJsonUtilsGetAsString(signaturePlace, JAdESHeaderParameterNamesAddressCountry))
		return result
	}
	return nil
}

// SignaturePolicyStore returns the Signature Policy Store from the signature. Port of
// getSignaturePolicyStore().
//
// Java wraps this whole method in try/catch(Exception), logging "Cannot read signature policy
// store". None of the operations below (map/string extraction, base64 decoding) has a Go panic
// source to guard against, so the try/catch is not reproduced with a recover().
func (s *Signature) SignaturePolicyStore() *model.SignaturePolicyStore {
	sigPStMap := s.unsignedPropertyAsMap(JAdESHeaderParameterNamesSigPSt)
	if sigPStMap.Size() == 0 {
		return nil
	}
	signaturePolicyStore := model.NewSignaturePolicyStore()

	sigPolDocBase64 := DSSJsonUtilsGetAsString(sigPStMap, JAdESHeaderParameterNamesSigPolDoc)
	if utils.IsStringNotEmpty(sigPolDocBase64) {
		policyContent := model.NewInMemoryDocument(utils.FromBase64(sigPolDocBase64))
		signaturePolicyStore.SetSignaturePolicyContent(policyContent)
	}

	sigPolLocalURI := DSSJsonUtilsGetAsString(sigPStMap, JAdESHeaderParameterNamesSigPolLocalURI)
	if utils.IsStringNotEmpty(sigPolLocalURI) {
		signaturePolicyStore.SetSigPolDocLocalURI(sigPolLocalURI)
	}

	spDSpec := sigPStMap.Value(JAdESHeaderParameterNamesSpDspec)
	if spDSpec != nil {
		spDocSpecification := DSSJsonUtilsParseSPDocSpecification(spDSpec)
		signaturePolicyStore.SetSpDocSpecification(spDocSpecification)
	}

	return signaturePolicyStore
}

// CommitmentTypeIndications obtains the information concerning commitment type indication linked
// to the signature. Port of getCommitmentTypeIndications().
func (s *Signature) CommitmentTypeIndications() []*signature.CommitmentTypeIndication {
	var result []*signature.CommitmentTypeIndication
	signedCommitments := s.jws.ProtectedHeaderValueAsList(JAdESHeaderParameterNamesSrCms)
	for _, signedCommitment := range signedCommitments {
		signedCommitmentMap := DSSJsonUtilsToMapValue(signedCommitment)
		if signedCommitmentMap.Size() != 0 {
			commIdMap := DSSJsonUtilsGetAsMap(signedCommitmentMap, JAdESHeaderParameterNamesCommId)
			if commIdMap.Size() != 0 {
				uri := DSSJsonUtilsGetAsString(commIdMap, JAdESHeaderParameterNamesId)
				uri = spi.DSSUtilsObjectIdentifierValue(uri)
				if utils.IsStringNotBlank(uri) {
					commitmentTypeIndication := signature.NewCommitmentTypeIndication(uri)
					desc := DSSJsonUtilsGetAsString(commIdMap, JAdESHeaderParameterNamesDesc)
					commitmentTypeIndication.SetDescription(desc)
					docRefs := DSSJsonUtilsGetAsList(commIdMap, JAdESHeaderParameterNamesDocRefs)
					commitmentTypeIndication.SetDocumentReferences(DSSJsonUtilsToListOfStrings(docRefs))
					result = append(result, commitmentTypeIndication)
				} else {
					// Upstream logs "Id parameter in the OID with the value '{}' is not
					// conformant! The entry is skipped.".
				}
			}
		}
	}
	return result
}

// ContentType returns the content type of the signature, not applicable for JAdES (see
// TS 119 102-2 v1.4.1). Port of getContentType().
func (s *Signature) ContentType() string {
	return ""
}

// MimeType returns the MIME type of the signature. Port of getMimeType().
//
// TS 119 102-2 v1.4.1: The MimeType element shall contain the value of the cty header
// parameter, prefixed with the string "application/" when this prefix has been omitted in the
// cty header parameter. NOTE: the sigD header parameter has one member that contains
// information of the format and type of the constituents of the JWS Payload.
func (s *Signature) MimeType() string {
	value := s.jws.ContentTypeHeaderValue()
	if utils.IsStringEmpty(value) {
		// sigD: return the first one when present
		ctys := s.signedDataContentTypeList()
		if utils.IsCollectionNotEmpty(ctys) {
			value = ctys[0]
		}
	}
	if utils.IsStringNotEmpty(value) {
		return DSSJsonUtilsMimeTypeString(value)
	}
	return ""
}

// SignatureType returns value of the "typ" header parameter, declaring the media type of the
// JWS, when present. Port of getSignatureType().
func (s *Signature) SignatureType() string {
	value := s.jws.ProtectedHeaderValueAsString(jose.HeaderType)
	if utils.IsStringNotEmpty(value) {
		return DSSJsonUtilsMimeTypeString(value)
	}
	return ""
}

// CertifiedSignerRoles returns the certified roles of the signer. Port of
// getCertifiedSignerRoles().
func (s *Signature) CertifiedSignerRoles() []*signature.SignerRole {
	var result []*signature.SignerRole
	signerAttributes := s.signerAttributes()
	if signerAttributes.Size() != 0 {
		certified := DSSJsonUtilsGetAsList(signerAttributes, JAdESHeaderParameterNamesCertified)
		for _, certifiedItem := range certified {
			certifiedVal := jadesSignatureCertifiedVal(certifiedItem)
			if utils.IsStringNotEmpty(certifiedVal) {
				result = append(result, signature.NewSignerRole(certifiedVal, enumerations.EndorsementTypeCertified))
			}
		}
	}
	return result
}

// jadesSignatureCertifiedVal ports the private getCertifiedVal(Object).
func jadesSignatureCertifiedVal(certifiedItem any) string {
	certifiedItemMap := DSSJsonUtilsToMap(certifiedItem, JAdESHeaderParameterNamesCertifiedAttrs)

	x509AttrCert := DSSJsonUtilsGetAsMap(certifiedItemMap, JAdESHeaderParameterNamesX509AttrCert)
	if x509AttrCert.Size() != 0 {
		return DSSJsonUtilsGetAsString(x509AttrCert, JAdESHeaderParameterNamesVal)
	}

	otherAttrCert := DSSJsonUtilsGetAsMap(certifiedItemMap, JAdESHeaderParameterNamesOtherAttrCert)
	if otherAttrCert.Size() != 0 {
		// Upstream logs "Unsupported {} found".
		return ""
	}

	// Upstream logs "One of types {} or {} is expected in {}".
	return ""
}

// ClaimedSignerRoles returns the claimed roles of the signer. Port of getClaimedSignerRoles().
func (s *Signature) ClaimedSignerRoles() []*signature.SignerRole {
	signerAttributes := s.signerAttributes()
	if signerAttributes.Size() != 0 {
		claimed := DSSJsonUtilsGetAsList(signerAttributes, JAdESHeaderParameterNamesClaimed)
		if utils.IsCollectionNotEmpty(claimed) {
			return jadesSignatureQArraySignerRoles(claimed, enumerations.EndorsementTypeClaimed)
		}
	}
	return nil
}

// SignedAssertions returns the list of embedded signed assertions. Port of
// getSignedAssertions().
func (s *Signature) SignedAssertions() []*signature.SignerRole {
	signerAttributes := s.signerAttributes()
	if signerAttributes.Size() != 0 {
		signedAssertions := DSSJsonUtilsGetAsList(signerAttributes, JAdESHeaderParameterNamesSignedAssertions)
		if utils.IsCollectionNotEmpty(signedAssertions) {
			return jadesSignatureQArraySignerRoles(signedAssertions, enumerations.EndorsementTypeSigned)
		}
	}
	return nil
}

// jadesSignatureQArraySignerRoles ports the private getQArraySignerRoles(List, EndorsementType).
func jadesSignatureQArraySignerRoles(qArrays []any, category enumerations.EndorsementType) []*signature.SignerRole {
	var result []*signature.SignerRole
	for _, qArray := range qArrays {
		qArrayMap := DSSJsonUtilsToMapValue(qArray)
		vals := DSSJsonUtilsGetAsList(qArrayMap, JAdESHeaderParameterNamesQVals)
		for _, val := range vals {
			result = append(result, signature.NewSignerRole(jadesSignatureValueToString(val), category))
		}
	}
	return result
}

// jadesSignatureValueToString mirrors Java's polymorphic Object#toString(), which
// getQArraySignerRoles's `val.toString()` relies on for whatever JSON value type a 'qVals' entry
// happens to hold - unlike DSSJsonUtilsToString(Object), the port of a DIFFERENT, String-only DSS
// utility method (DSSJsonUtils.toString(Object)) that answers "" for anything else.
func jadesSignatureValueToString(val any) string {
	switch v := val.(type) {
	case nil:
		return "null"
	case string:
		return v
	case bool:
		if v {
			return "true"
		}
		return "false"
	case *jose.Number:
		return v.String()
	default:
		return jose.JSON(v)
	}
}

// signerAttributes ports the private getSignerAttributes().
func (s *Signature) signerAttributes() *jose.Object {
	return s.jws.ProtectedHeaderValueAsMap(JAdESHeaderParameterNamesSrAts)
}

// CounterSignatures returns a list of counter signatures applied to this signature. Port of
// getCounterSignatures().
func (s *Signature) CounterSignatures() []validation.AdvancedSignature {
	if s.CachedCounterSignatures() != nil {
		return s.CachedCounterSignatures()
	}
	var counterSignatures []validation.AdvancedSignature

	etsiUComponents := s.EtsiUHeader().Attributes()
	for _, etsiUComponent := range etsiUComponents {
		if JAdESHeaderParameterNamesCSig == etsiUComponent.HeaderName() {
			counterSignature, err := DSSJsonUtilsExtractJAdESCounterSignature(etsiUComponent, s)
			if err != nil {
				// Upstream lets extractJAdESCounterSignature swallow its own failures and
				// return null; a parse error here is treated the same way (skip this entry).
				continue
			}
			if counterSignature != nil {
				counterSignature.SetFilename(s.Filename())
				counterSignatures = append(counterSignatures, counterSignature)
			}
		}
	}

	s.SetCachedCounterSignatures(counterSignatures)
	return counterSignatures
}

// DAIdentifier is not applicable for JAdES. Port of getDAIdentifier().
func (s *Signature) DAIdentifier() string {
	return ""
}

// BuildSignaturePolicy extracts a signature policy from a signature and builds the object. Port
// of the protected buildSignaturePolicy().
func (s *Signature) BuildSignaturePolicy() *signature.Policy {
	sigPolicy := s.jws.ProtectedHeaderValueAsMap(JAdESHeaderParameterNamesSigPid)
	if sigPolicy.Size() == 0 {
		return nil
	}
	policyId := DSSJsonUtilsGetAsMap(sigPolicy, JAdESHeaderParameterNamesId)
	if policyId.Size() == 0 {
		return nil
	}

	id := DSSJsonUtilsGetAsString(policyId, JAdESHeaderParameterNamesId)
	signaturePolicy := signature.NewSignaturePolicyWithIdentifier(spi.DSSUtilsObjectIdentifierValue(id))
	desc := DSSJsonUtilsGetAsString(policyId, JAdESHeaderParameterNamesDesc)
	signaturePolicy.SetDescription(desc)
	docRefs := DSSJsonUtilsGetAsList(policyId, JAdESHeaderParameterNamesDocRefs)
	signaturePolicy.SetDocumentationReferences(DSSJsonUtilsToListOfStrings(docRefs))

	if digest, ok := DSSJsonUtilsDigest(sigPolicy); ok {
		signaturePolicy.SetDigest(digest)
	}

	qualifiers := DSSJsonUtilsGetAsList(sigPolicy, JAdESHeaderParameterNamesSigPQuals)
	if utils.IsCollectionNotEmpty(qualifiers) {
		signaturePolicy.SetURI(jadesSignatureSPUri(qualifiers))
		signaturePolicy.SetUserNotice(jadesSignatureSPUserNotice(qualifiers))
		signaturePolicy.SetDocSpecification(jadesSignatureSPDSpec(qualifiers))
	}

	digPSp := DSSJsonUtilsGetAsBoolean(sigPolicy, JAdESHeaderParameterNamesDigPSp)
	if digPSp != nil {
		signaturePolicy.SetHashAsInTechnicalSpecification(*digPSp)
	}

	return signaturePolicy
}

// jadesSignatureSPUri ports the private getSPUri(List).
func jadesSignatureSPUri(qualifiers []any) string {
	for _, qualifier := range qualifiers {
		qualifierMap := DSSJsonUtilsToMap(qualifier, JAdESHeaderParameterNamesSigPQual)
		if qualifierMap.Size() != 0 {
			spUri := DSSJsonUtilsGetAsString(qualifierMap, JAdESHeaderParameterNamesSpURI)
			if utils.IsStringNotEmpty(spUri) {
				return spUri
			}
		}
	}
	return ""
}

// jadesSignatureSPUserNotice ports the private getSPUserNotice(List).
//
// Java wraps the body in try/catch(Exception), logging "Unable to build SPUserNotice qualifier"
// and returning null; nothing in the translated body below has a Go panic source, so the
// try/catch is not reproduced with a recover() (see the file header note).
func jadesSignatureSPUserNotice(qualifiers []any) *model.UserNotice {
	for _, qualifier := range qualifiers {
		qualifierMap := DSSJsonUtilsToMap(qualifier, JAdESHeaderParameterNamesSigPQual)
		if qualifierMap.Size() != 0 {
			spUserNotice := DSSJsonUtilsGetAsMap(qualifierMap, JAdESHeaderParameterNamesSpUserNotice)
			if spUserNotice.Size() != 0 {
				userNotice := model.NewUserNotice()

				noticeRef := DSSJsonUtilsGetAsMap(spUserNotice, JAdESHeaderParameterNamesNoticeRef)
				if noticeRef.Size() != 0 {
					organization := DSSJsonUtilsGetAsString(noticeRef, JAdESHeaderParameterNamesOrgantization)
					if utils.IsStringNotBlank(organization) {
						userNotice.SetOrganization(organization)
					}

					noticeNumbers := DSSJsonUtilsGetAsList(noticeRef, JAdESHeaderParameterNamesNoticeNumbers)
					if utils.IsCollectionNotEmpty(noticeNumbers) {
						numbers := DSSJsonUtilsToListOfNumbers(noticeNumbers)
						ints := make([]int, len(numbers))
						for i, number := range numbers {
							ints[i] = int(number.Int64())
						}
						userNotice.SetNoticeNumbers(ints...)
					}
				}
				explText := DSSJsonUtilsGetAsString(spUserNotice, JAdESHeaderParameterNamesExplText)
				if utils.IsStringNotBlank(explText) {
					userNotice.SetExplicitText(explText)
				}
				return userNotice
			}
		}
	}
	return nil
}

// jadesSignatureSPDSpec ports the private getSPDSpec(List).
func jadesSignatureSPDSpec(qualifiers []any) *model.SpDocSpecification {
	for _, qualifier := range qualifiers {
		qualifierMap := DSSJsonUtilsToMap(qualifier, JAdESHeaderParameterNamesSigPQual)
		if qualifierMap.Size() != 0 {
			spDSpec := qualifierMap.Value(JAdESHeaderParameterNamesSpDspec)
			if spDSpec != nil {
				return DSSJsonUtilsParseSPDocSpecification(spDSpec)
			}
		}
	}
	return nil
}

// SignatureValue gets the SignatureValue bytes. Port of getSignatureValue().
func (s *Signature) SignatureValue() []byte {
	return s.jws.SignatureValue()
}

// EtsiUHeader returns unsigned properties embedded into the 'etsiU' array. Port of the public
// EtsiUHeader getEtsiUHeader().
func (s *Signature) EtsiUHeader() *EtsiUHeader {
	if s.etsiUHeader == nil {
		s.etsiUHeader = NewJAdESEtsiUHeader(s.jws)
	}
	return s.etsiUHeader
}

// BuildSignatureDigestReference builds a new DigestReference according to the
// applicable signature format rules. Port of buildSignatureDigestReference(DigestAlgorithm).
//
// TODO: no definition available in ETSI TS 119 442 - V1.1.1 (upstream's own TODO, reproduced).
func (s *Signature) BuildSignatureDigestReference(digestAlgorithm enumerations.DigestAlgorithm) *signature.DigestReference {
	encodedHeader := s.jws.EncodedHeader()
	var payload string
	if s.jws.IsRfc7797UnencodedPayload() {
		payload = s.jws.UnverifiedPayload()
	} else {
		payload = s.jws.EncodedPayload()
	}
	encodedSignature := s.jws.EncodedSignature()
	signatureReferenceBytes := []byte(DSSJsonUtilsConcatenate(encodedHeader, payload, encodedSignature))
	digestValue, err := spi.DSSUtilsDigest(digestAlgorithm, signatureReferenceBytes)
	if err != nil {
		panic(err)
	}
	return signature.NewSignatureDigestReference(model.NewDigest(digestAlgorithm, digestValue))
}

// DataToBeSignedRepresentation returns the DTBSR. Port of getDataToBeSignedRepresentation().
func (s *Signature) DataToBeSignedRepresentation() model.Digest {
	referenceValidations := s.ReferenceValidations()
	for _, referenceValidation := range referenceValidations {
		if enumerations.DigestMatcherTypeJWSSigningInput == referenceValidation.Type() {
			if referenceValidation.IsFound() {
				return referenceValidation.Digest()
			}
			return model.Digest{}
		}
	}
	// shall not happen
	panic(model.NewDSSError("JWS_SIGNING_INPUT is not found! Unable to compute DTBSR."))
}

// SignatureIdentifierBuilder returns a builder to define and build a signature Id. Port of the
// protected getSignatureIdentifierBuilder().
func (s *Signature) SignatureIdentifierBuilder() validation.SignatureIdentifierBuilder {
	return NewJAdESSignatureIdentifierBuilder(s)
}

// CheckSignatureIntegrity verifies the signature integrity; checks if the signed content has not
// been tampered with. Port of checkSignatureIntegrity().
func (s *Signature) CheckSignatureIntegrity() {
	if s.jadesCachedCryptoVerification != nil {
		return
	}

	verification := signature.NewSignatureCryptographicVerification()
	s.jadesCachedCryptoVerification = verification
	s.SetSignatureCryptographicVerification(verification)

	refsFound := false
	refsIntact := false

	referenceValidations := s.ReferenceValidations()

	if utils.IsCollectionNotEmpty(referenceValidations) {
		refsFound = true
		refsIntact = true

		for _, referenceValidation := range referenceValidations {
			if enumerations.DigestMatcherTypeJWSSigningInput == referenceValidation.Type() {
				verification.SetSignatureIntact(referenceValidation.IsIntact())
				for _, errorMessage := range referenceValidation.ErrorMessages() {
					verification.SetErrorMessage(errorMessage)
				}
			}
			refsFound = refsFound && referenceValidation.IsFound()
			refsIntact = refsIntact && referenceValidation.IsIntact()
		}
	}

	verification.SetReferenceDataFound(refsFound)
	verification.SetReferenceDataIntact(refsIntact)
}

// ReferenceValidations returns reference validations for the signature. Port of the public
// getReferenceValidations(). Left entirely abstract by DefaultAdvancedSignature (see that type's
// cachedReferenceValidations field doc comment); CachedReferenceValidations/
// SetCachedReferenceValidations are used as the cache slot.
func (s *Signature) ReferenceValidations() []*model.ReferenceValidation {
	if s.CachedReferenceValidations() == nil {
		var referenceValidations []*model.ReferenceValidation

		signingInputReferenceValidation := s.signingInputReferenceValidation()
		referenceValidations = append(referenceValidations, signingInputReferenceValidation)

		if s.IsDetachedSignature() {
			detachedReferenceValidations := s.detachedReferenceValidations()
			if utils.IsCollectionNotEmpty(detachedReferenceValidations) {
				referenceValidations = append(referenceValidations, detachedReferenceValidations...)
			}
		}

		if s.IsCounterSignature() {
			referenceValidations = append(referenceValidations, s.counterSignatureReferenceValidation())
		}

		if s.IsKeyBindingSignature() {
			referenceValidations = append(referenceValidations, s.keyBindingSignatureReferenceValidation())
		}

		s.SetCachedReferenceValidations(referenceValidations)
	}
	return s.CachedReferenceValidations()
}

// signingInputReferenceValidation ports the private getSigningInputReferenceValidation().
//
// Java wraps the whole method body in try/catch(Exception), logging "The validation of signed
// input failed!" on any failure and returning the (possibly partially populated)
// signatureValueReferenceValidation regardless. The inner try/catch around the payload
// resolution is preserved (jadesSignaturePayload below); the outer one has no remaining Go panic
// source once that inner recovery is in place (every DSSJsonUtils/spi call downstream returns
// values, not panics, in this port) and is accordingly not reproduced with a second recover().
func (s *Signature) signingInputReferenceValidation() *model.ReferenceValidation {
	signatureValueReferenceValidation := model.NewReferenceValidation()
	signatureValueReferenceValidation.SetType(enumerations.DigestMatcherTypeJWSSigningInput)

	encodedHeader := s.jws.EncodedHeader()
	if utils.IsStringEmpty(encodedHeader) {
		return signatureValueReferenceValidation
	}

	s.jadesSignaturePayload(signatureValueReferenceValidation)

	signatureAlgorithm := s.SignatureAlgorithm()
	if signatureAlgorithm != "" {
		dataToSign := DSSJsonUtilsSigningInputBytes(s.jws)
		digestAlgorithm := signatureAlgorithm.DigestAlgorithm()
		digestValue, err := spi.DSSUtilsDigest(digestAlgorithm, dataToSign)
		if err != nil {
			// Upstream logs "The validation of signed input failed! Reason : {}".
			return signatureValueReferenceValidation
		}
		signatureValueReferenceValidation.SetDigest(model.NewDigest(digestAlgorithm, digestValue))

		s.jws.SetDoKeyValidation(false) // restrict on key size,...

		candidatesForSigningCertificate := s.CandidatesForSigningCertificate()

		signingCertificateValidator := NewJAdESSignatureIntegrityValidator(s.jws)
		certificateValidity := signingCertificateValidator.Validate(candidatesForSigningCertificate)
		if certificateValidity != nil {
			_ = candidatesForSigningCertificate.SetTheCertificateValidity(certificateValidity)
		}

		errorMessages := signingCertificateValidator.ErrorMessages()
		signatureValueReferenceValidation.SetErrorMessages(errorMessages)
		signatureValueReferenceValidation.SetIntact(certificateValidity != nil)
	}

	return signatureValueReferenceValidation
}

// jadesSignaturePayload ports the inner try body of getSigningInputReferenceValidation() that
// resolves and sets the JWS payload for a detached signature, catching any failure the same way
// (a dropped LOG.warn "Unable to determine a JWS payload").
func (s *Signature) jadesSignaturePayload(signatureValueReferenceValidation *model.ReferenceValidation) {
	defer func() {
		if recover() != nil {
			// Upstream logs "Unable to determine a JWS payload. Reason : {}".
		}
	}()

	sigDMechanism := s.SigDMechanism()
	detachedContentPresent := utils.IsCollectionNotEmpty(s.DetachedContents())
	switch {
	case !s.IsDetachedSignature():
		// not detached
		signatureValueReferenceValidation.SetFound(true)

	case sigDMechanism == nil && detachedContentPresent:
		// simple detached signature
		payload := s.incorporatedPayload()
		s.jws.SetPayloadOctets(payload)
		signatureValueReferenceValidation.SetFound(len(s.DetachedContents()) == 1)

	case sigDMechanism != nil && enumerations.SigDMechanismHTTPHeaders == *sigDMechanism:
		// detached with HTTP_HEADERS mechanism
		payload := s.payloadForHttpHeadersMechanism()
		s.jws.SetPayloadOctets(payload)
		signatureValueReferenceValidation.SetFound(payload != nil)

	case sigDMechanism != nil && enumerations.SigDMechanismObjectIDByURI == *sigDMechanism:
		// detached with OBJECT_ID_BY_URI mechanism
		signedDataUriList := s.signedDataUriList()
		payload := s.payloadForObjectIdByUriMechanism(signedDataUriList)
		s.jws.SetPayloadOctets(payload)
		signatureValueReferenceValidation.SetFound(payload != nil)
		signatureValueReferenceValidation.SetDataObjectReferences(signedDataUriList)

	case sigDMechanism != nil && enumerations.SigDMechanismObjectIDByURIHash == *sigDMechanism:
		// the sigD itself is signed with OBJECT_ID_BY_URI_HASH mechanism
		signatureValueReferenceValidation.SetFound(true)

	default:
		// otherwise original content is not found
		// Upstream logs "The payload is not found! The detached content must be provided!".
	}
}

// Kid gets Kid value when present. Port of the public String getKid().
func (s *Signature) Kid() string {
	return s.jws.KeyIDHeaderValue()
}

// detachedReferenceValidations ports the private getDetachedReferenceValidations().
func (s *Signature) detachedReferenceValidations() []*model.ReferenceValidation {
	sigDMechanism := s.SigDMechanism()
	if sigDMechanism != nil {
		switch *sigDMechanism {
		case enumerations.SigDMechanismHTTPHeaders, enumerations.SigDMechanismObjectIDByURI:
			// the documents are added to the payload, not possible to extract separate
			// reference validations
		case enumerations.SigDMechanismObjectIDByURIHash:
			return s.referenceValidationsByUriHashMechanism()
		default:
			// Upstream logs "The SigDMechanism '{}' is not supported!".
		}
	}
	return nil
}

// SigDMechanism returns a mechanism used in 'sigD' to cover a detached content, or nil when
// absent. Port of the public SigDMechanism getSigDMechanism().
//
// Returns a pointer (not the plain enumerations.SigDMechanism value every other enum accessor in
// this port uses, with "" as the null sentinel) to match the surface
// jades_timestamp_message_digest_builder.go depends on.
func (s *Signature) SigDMechanism() *enumerations.SigDMechanism {
	signatureDetached := s.jws.ProtectedHeaderValueAsMap(JAdESHeaderParameterNamesSigD)
	if signatureDetached.Size() != 0 {
		mechanismUri := DSSJsonUtilsGetAsString(signatureDetached, JAdESHeaderParameterNamesMId)
		sigDMechanism := enumerations.SigDMechanismForJAdESUri(mechanismUri)
		if sigDMechanism == "" {
			// Upstream logs "The sigDMechanism with uri '{}' is not supported!".
			return nil
		}
		return &sigDMechanism
	}
	return nil
}

// incorporatedPayload ports the private getIncorporatedPayload().
func (s *Signature) incorporatedPayload() []byte {
	payload, err := DSSJsonUtilsDocumentOctets(s.DetachedContents()[0], !s.jws.IsRfc7797UnencodedPayload())
	if err != nil {
		panic(err)
	}
	return payload
}

// payloadForHttpHeadersMechanism ports the private getPayloadForHttpHeadersMechanism().
//
// Panics with the Java message when the detached contents are missing (IllegalArgumentException).
func (s *Signature) payloadForHttpHeadersMechanism() []byte {
	if utils.IsCollectionEmpty(s.DetachedContents()) {
		panic("The detached contents shall be provided for validating a detached signature!")
	}

	/*
	 * Case-insensitive, see TS 119 182-1 "5.2.8.2 Mechanism HttpHeaders":
	 *
	 * For this referencing mechanism, the contents of the pars member shall be an array of
	 * lowercased names of HTTP header fields, each one with the semantics and syntax specified
	 * in clause 2.1.3 of draft-cavage-http-signatures-10: "Signing HTTP Messages" [17].
	 */
	documentsByUri := s.SignedDocumentsByHTTPHeaderName()
	httpHeadersPayloadBuilder := NewHttpHeadersPayloadBuilder(documentsByUri, false)

	payload, err := httpHeadersPayloadBuilder.Build()
	if err != nil {
		panic(err)
	}
	return payload
}

// SignedDocumentsByHTTPHeaderName returns a list of signed documents by the list of URIs present
// in 'sigD'. Keeps the original order according to 'pars' dictionary content. Used in
// HTTPHeaders detached signature mechanism. Port of the public
// List<DSSDocument> getSignedDocumentsByHTTPHeaderName().
//
// Panics with the Java message when a named signed document is not found (IllegalArgumentException).
func (s *Signature) SignedDocumentsByHTTPHeaderName() []model.DSSDocument {
	signedDataUriList := s.signedDataUriList()

	if utils.IsCollectionEmpty(s.DetachedContents()) {
		// Upstream logs "Detached content is not provided!".
		return nil
	}

	if len(signedDataUriList) == 1 && len(s.DetachedContents()) == 1 {
		return s.DetachedContents()
	}

	var signedDocuments []model.DSSDocument
	for _, signedDataName := range signedDataUriList {
		found := false
		for _, document := range s.DetachedContents() {
			if utils.AreStringsEqualIgnoreCase(signedDataName, document.Name()) {
				found = true
				signedDocuments = append(signedDocuments, document)
				// do not break - same name docs possible
			}
		}
		if !found {
			panic("The detached content for a signed data with name '" + signedDataName + "' has not been found!")
		}
	}

	return signedDocuments
}

// payloadForObjectIdByUriMechanism ports the private
// getPayloadForObjectIdByUriMechanism(List<String>).
//
// Panics with the Java message when the detached contents are missing (IllegalArgumentException).
func (s *Signature) payloadForObjectIdByUriMechanism(signedDataUriList []string) []byte {
	if utils.IsCollectionEmpty(s.DetachedContents()) {
		panic("The detached contents shall be provided for validating a detached signature!")
	}

	signedDocumentsByUri := s.signedDocumentsForUris(signedDataUriList)
	payload, err := DSSJsonUtilsConcatenateDSSDocuments(signedDocumentsByUri, !s.jws.IsRfc7797UnencodedPayload())
	if err != nil {
		panic(err)
	}
	return payload
}

// SignedDocumentsForObjectIdByUriMechanism returns a list of documents for ObjectIdByUrl or
// ObjectIdByUriHash mechanisms. Keeps the original order according to 'pars' dictionary content.
// Port of the public List<DSSDocument> getSignedDocumentsForObjectIdByUriMechanism().
func (s *Signature) SignedDocumentsForObjectIdByUriMechanism() []model.DSSDocument {
	signedDataUriList := s.signedDataUriList()
	return s.signedDocumentsForUris(signedDataUriList)
}

// signedDocumentsForUris ports the private getSignedDocumentsForUris(List<String>).
//
// Panics with the Java message when a named signed document is not found (IllegalArgumentException).
func (s *Signature) signedDocumentsForUris(signedDataUriList []string) []model.DSSDocument {
	var signedDocumentsByUri []model.DSSDocument
	if len(signedDataUriList) == 1 && len(s.DetachedContents()) == 1 {
		signedDocumentsByUri = []model.DSSDocument{s.DetachedContents()[0]}

	} else if utils.IsCollectionNotEmpty(signedDataUriList) {
		for _, signedDataName := range signedDataUriList {
			detachedDocumentByName := jadesSignatureDetachedDocumentByName(signedDataName, s.DetachedContents())
			if detachedDocumentByName != nil {
				signedDocumentsByUri = append(signedDocumentsByUri, detachedDocumentByName)
			} else {
				panic("The detached content for a signed data with name '" + signedDataName + "' has not been found!")
			}
		}
	}
	return signedDocumentsByUri
}

// referenceValidationsByUriHashMechanism ports the private
// getReferenceValidationsByUriHashMechanism().
func (s *Signature) referenceValidationsByUriHashMechanism() []*model.ReferenceValidation {
	detachedDocuments := s.DetachedContents()

	if utils.IsCollectionEmpty(s.DetachedContents()) {
		// Upstream logs "The detached content is not provided! Validation of '{}' is not
		// possible.".
		detachedDocuments = nil
		// continue in order to extract signed data references
	}

	signedDataHashMap, signedDataOrder := s.signedDataUriHashMap()
	if len(signedDataHashMap) == 0 {
		// Upstream logs "The SignedData has not been found or incorrect for detached content.".
		emptyReference := model.NewReferenceValidation()
		emptyReference.SetType(enumerations.DigestMatcherTypeSigDEntry)
		return []*model.ReferenceValidation{emptyReference}
	}

	digestAlgorithm := s.digestAlgorithmForDetachedContent()
	if digestAlgorithm == "" {
		// Upstream logs "The DigestAlgorithm has not been found for the detached content.".
	}

	var detachedReferenceValidations []*model.ReferenceValidation

	for _, signedDataName := range signedDataOrder {
		expectedDigestString := signedDataHashMap[signedDataName]

		referenceValidation := model.NewReferenceValidation()
		referenceValidation.SetType(enumerations.DigestMatcherTypeSigDEntry)
		referenceValidation.SetUri(signedDataName)

		expectedDigest := DSSJsonUtilsFromBase64Url(expectedDigestString)
		if digestAlgorithm != "" {
			referenceValidation.SetDigest(model.NewDigest(digestAlgorithm, expectedDigest))
		}

		var detachedDocument model.DSSDocument
		if len(signedDataHashMap) == 1 && len(detachedDocuments) == 1 {
			detachedDocument = detachedDocuments[0]
		} else {
			detachedDocument = s.detachedDocumentByDigest(digestAlgorithm, expectedDigest, signedDataName, detachedDocuments)
			if detachedDocument == nil {
				detachedDocument = jadesSignatureDetachedDocumentByName(signedDataName, detachedDocuments)
			}
		}

		if detachedDocument != nil {
			referenceValidation.SetFound(true)
			referenceValidation.SetDocument(detachedDocument)
			if digestAlgorithm != "" && s.isDocumentDigestMatch(detachedDocument, digestAlgorithm, expectedDigest, signedDataName) {
				referenceValidation.SetIntact(true)
			}
		} else {
			// Upstream logs "A detached document for the '{}' header with name '{}' has not
			// been found!".
		}

		detachedReferenceValidations = append(detachedReferenceValidations, referenceValidation)
	}

	if utils.IsCollectionEmpty(detachedReferenceValidations) {
		// add an empty reference if none found
		referenceValidation := model.NewReferenceValidation()
		referenceValidation.SetType(enumerations.DigestMatcherTypeSigDEntry)
		detachedReferenceValidations = append(detachedReferenceValidations, referenceValidation)
	}

	return detachedReferenceValidations
}

// digestAlgorithmForDetachedContent ports the private getDigestAlgorithmForDetachedContent().
func (s *Signature) digestAlgorithmForDetachedContent() enumerations.DigestAlgorithm {
	signatureDetached := s.jws.ProtectedHeaderValueAsMap(JAdESHeaderParameterNamesSigD)
	if signatureDetached.Size() != 0 {
		digestAlgoUri := DSSJsonUtilsGetAsString(signatureDetached, JAdESHeaderParameterNamesHashM)
		digestAlgorithm, err := enumerations.DigestAlgorithmForJAdES(digestAlgoUri)
		if err != nil {
			// Upstream logs "Unable to extract DigestAlgorithm for '{}' element. Reason : {}".
			return ""
		}
		return digestAlgorithm
	}
	return ""
}

// detachedDocumentByDigest ports the private getDetachedDocumentByDigest(DigestAlgorithm,
// byte[], String, List<DSSDocument>).
func (s *Signature) detachedDocumentByDigest(digestAlgorithm enumerations.DigestAlgorithm, expectedDigest []byte,
	signedDataName string, detachedContent []model.DSSDocument) model.DSSDocument {
	if digestAlgorithm == "" || expectedDigest == nil {
		return nil
	}
	for _, detachedDocument := range detachedContent {
		if s.isDocumentDigestMatch(detachedDocument, digestAlgorithm, expectedDigest, signedDataName) {
			return detachedDocument
		}
	}
	return nil
}

// jadesSignatureDetachedDocumentByName ports the private getDetachedDocumentByName(String,
// List<DSSDocument>).
func jadesSignatureDetachedDocumentByName(documentName string, detachedContent []model.DSSDocument) model.DSSDocument {
	documentName = spi.DSSUtilsDecodeURI(documentName)
	return spi.DSSUtilsDocumentWithName(detachedContent, documentName)
}

// signedDataUriHashMap ports the private getSignedDataUriHashMap(), returning both the map
// (Java's Map<String, String>) and its insertion order (Java's LinkedHashMap iteration order),
// since Go maps do not preserve one.
func (s *Signature) signedDataUriHashMap() (map[string]string, []string) {
	signedDataHashMap := make(map[string]string)
	var order []string

	signedDataUriList := s.signedDataUriList()
	signedDataHashList := s.signedDataHashList()
	if len(signedDataUriList) != len(signedDataHashList) {
		// Upstream logs "The size of 'pars' and 'hashV' dictionaries does not match! See '5.2.8
		// The sigD header parameter'.".
		return signedDataHashMap, order
	}

	for i := 0; i < len(signedDataUriList); i++ {
		if _, exists := signedDataHashMap[signedDataUriList[i]]; !exists {
			order = append(order, signedDataUriList[i])
		}
		signedDataHashMap[signedDataUriList[i]] = signedDataHashList[i]
	}
	return signedDataHashMap, order
}

// signedDataUriList ports the private getSignedDataUriList().
func (s *Signature) signedDataUriList() []string {
	signatureDetached := s.jws.ProtectedHeaderValueAsMap(JAdESHeaderParameterNamesSigD)
	if signatureDetached.Size() != 0 {
		pars := DSSJsonUtilsGetAsList(signatureDetached, JAdESHeaderParameterNamesPars)
		return DSSJsonUtilsToListOfStrings(pars)
	}
	return nil
}

// signedDataHashList ports the private getSignedDataHashList().
func (s *Signature) signedDataHashList() []string {
	signatureDetached := s.jws.ProtectedHeaderValueAsMap(JAdESHeaderParameterNamesSigD)
	if signatureDetached.Size() != 0 {
		pars := DSSJsonUtilsGetAsList(signatureDetached, JAdESHeaderParameterNamesHashV)
		return DSSJsonUtilsToListOfStrings(pars)
	}
	return nil
}

// signedDataContentTypeList ports the private getSignedDataContentTypeList().
func (s *Signature) signedDataContentTypeList() []string {
	signatureDetached := s.jws.ProtectedHeaderValueAsMap(JAdESHeaderParameterNamesSigD)
	if signatureDetached.Size() != 0 {
		ctys := DSSJsonUtilsGetAsList(signatureDetached, JAdESHeaderParameterNamesCtys)
		return DSSJsonUtilsToListOfStrings(ctys)
	}
	return nil
}

// isDocumentDigestMatch ports the private isDocumentDigestMatch(DSSDocument, DigestAlgorithm,
// byte[], String).
//
// Panics with the underlying error when the document's digest cannot be computed: the Java
// getDigestValue/toBase64Url calls this reaches are not guarded by a try/catch here, so an
// unchecked exception would propagate out of this (and, transitively, getReferenceValidations())
// uncaught.
func (s *Signature) isDocumentDigestMatch(document model.DSSDocument, digestAlgorithm enumerations.DigestAlgorithm,
	expectedDigest []byte, signedDataName string) bool {
	_, isDigestDocument := document.(*model.DigestDocument)

	var computedDigestValue []byte
	if s.jws.IsRfc7797UnencodedPayload() || isDigestDocument {
		value, err := document.DigestValue(digestAlgorithm)
		if err != nil {
			panic(err)
		}
		computedDigestValue = value
	} else {
		base64UrlEncodedDocument, err := DSSJsonUtilsToBase64UrlDocument(document)
		if err != nil {
			panic(err)
		}
		computedDigestValue, err = spi.DSSUtilsDigest(digestAlgorithm, []byte(base64UrlEncodedDocument))
		if err != nil {
			panic(err)
		}
	}

	if bytes.Equal(expectedDigest, computedDigestValue) {
		return true
	}
	// Upstream logs "The computed digest '{}' from a document with name '{}' does not match one
	// provided on the sigD : {}!" at WARN (matching name) or DEBUG (mismatching name) level.
	_ = signedDataName
	return false
}

// counterSignatureReferenceValidation ports the private getCounterSignatureReferenceValidation().
func (s *Signature) counterSignatureReferenceValidation() *model.ReferenceValidation {
	referenceValidation := model.NewReferenceValidation()
	referenceValidation.SetType(enumerations.DigestMatcherTypeCounterSignedSignatureValue)

	masterSignature, _ := s.MasterSignature().(*Signature)
	if masterSignature != nil {
		signatureValue := masterSignature.jws.SignatureValue()
		if utils.IsArrayNotEmpty(signatureValue) {
			referenceValidation.SetFound(true)
		}

		unverifiedPayloadBytes := s.jws.UnverifiedPayloadBytes()
		if utils.IsArrayNotEmpty(unverifiedPayloadBytes) {
			intact := bytes.Equal(signatureValue, unverifiedPayloadBytes)
			if !intact {
				// Upstream logs "The payload of a cSig with Id '{}' does not match the
				// signature value of its master signature!".
			}
			referenceValidation.SetIntact(intact)
		} else {
			// nothing to compare against for an attached signature
			referenceValidation.SetIntact(true)
		}
	}

	return referenceValidation
}

// keyBindingSignatureReferenceValidation ports the private
// getKeyBindingSignatureReferenceValidation().
//
// Panics with the underlying error when the key binding input digest cannot be computed: the
// Java DSSDocument#getDigest call this reaches is not guarded by a try/catch here.
func (s *Signature) keyBindingSignatureReferenceValidation() *model.ReferenceValidation {
	referenceValidation := model.NewReferenceValidation()
	referenceValidation.SetType(enumerations.DigestMatcherTypeEAAKeyBinding)

	sdHash := s.sdHash()
	if sdHash != nil {
		if len(s.DetachedContents()) == 1 {
			referenceValidation.SetFound(true)

			digestAlgorithm := s.sdAlg()
			if digestAlgorithm != "" {
				referenceValidation.SetDigest(model.NewDigest(digestAlgorithm, sdHash))

				kbInputDigest, err := s.DetachedContents()[0].Digest(digestAlgorithm)
				if err != nil {
					panic(err)
				}
				intact := bytes.Equal(sdHash, kbInputDigest.Value())
				if !intact {
					// Upstream logs "The sd_hash present within key binding signature does not
					// match the hash over the computed key binding input! Found : {}, Computed :
					// {}".
				}
				referenceValidation.SetIntact(intact)
			}

		} else {
			// Upstream logs "No detached content was found for the key binding signature
			// verification!".
		}
	}

	return referenceValidation
}

// sdHash ports the private getSdHash().
func (s *Signature) sdHash() []byte {
	payload, err := s.jws.DecodedPayload()
	if err != nil {
		panic(err)
	}
	sdHashB64Url := DSSJsonUtilsGetAsString(payload, "sd_hash")
	if sdHashB64Url != "" && DSSJsonUtilsIsBase64UrlEncoded(sdHashB64Url) {
		return DSSJsonUtilsFromBase64Url(sdHashB64Url)
	}
	// Upstream logs "A base64url-encoded sd_hash header shall be present within SD-JWT key
	// binding signature payload!".
	return nil
}

// sdAlg ports the private getSdAlg().
//
// Panics with the Java message when there are no EAA signatures (IllegalStateException).
func (s *Signature) sdAlg() enumerations.DigestAlgorithm {
	eaaSignatures := s.EAA().Signatures()
	if utils.IsCollectionEmpty(eaaSignatures) {
		panic("EAA signatures cannot be null or empty!")
	}
	eaaSignature := eaaSignatures[0].(*Signature)
	payload, err := eaaSignature.jws.DecodedPayload()
	if err != nil {
		panic(err)
	}
	sdAlgId := DSSJsonUtilsGetAsString(payload, "_sd_alg")
	if sdAlgId == "" {
		// Upstream logs "No _sd_alg header found within the SD-JWT payload!".
	}
	digestAlgorithm, err := enumerations.DigestAlgorithmForSdJwtId(sdAlgId)
	if err != nil {
		// Upstream logs "Unable to find a corresponding DigestAlgorithm for SD-JWT claim for
		// value '{}'!".
		return ""
	}
	return digestAlgorithm
}

// unsignedPropertyAsMap ports the private getUnsignedPropertyAsMap(String).
func (s *Signature) unsignedPropertyAsMap(headerName string) *jose.Object {
	unsignedPropertiesWithHeaderName := DSSJsonUtilsUnsignedPropertiesWithHeaderName(s.EtsiUHeader(), headerName)
	if utils.IsCollectionNotEmpty(unsignedPropertiesWithHeaderName) {
		// return the first occurrence
		return DSSJsonUtilsToMap(unsignedPropertiesWithHeaderName[0].Value(), headerName)
	}
	return nil
}

// OriginalDocuments returns a list of original documents signed by the signature. Port of the
// public List<DSSDocument> getOriginalDocuments().
//
// Java lets SignedDocumentsByHTTPHeaderName/SignedDocumentsForObjectIdByUriMechanism's unchecked
// IllegalArgumentException (a named detached document not found) propagate out of this method;
// this port instead reports it as an error, matching the (documents, error) contract
// abstract_jws_document_analyzer.go relies on and its own caller's err != nil handling.
func (s *Signature) OriginalDocuments() (documents []model.DSSDocument, err error) {
	defer func() {
		if r := recover(); r != nil {
			documents = nil
			if e, ok := r.(error); ok {
				err = e
			} else {
				err = model.NewDSSError(fmt.Sprint(r))
			}
		}
	}()

	if s.IsDetachedSignature() {
		var originalDocuments []model.DSSDocument

		referenceValidations := s.ReferenceValidations()
		for _, referenceValidation := range referenceValidations {
			if enumerations.DigestMatcherTypeSigDEntry == referenceValidation.Type() && referenceValidation.IsIntact() {
				detachedDocument := referenceValidation.Document()
				if detachedDocument != nil {
					originalDocuments = append(originalDocuments, detachedDocument)
				}
			}
		}

		if utils.IsCollectionEmpty(originalDocuments) {
			// check if the signature of an old detached format
			signatureCryptographicVerification := s.SignatureCryptographicVerification()
			if signatureCryptographicVerification.IsSignatureIntact() {
				sigDMechanism := s.SigDMechanism()
				if len(s.DetachedContents()) == 1 {
					return []model.DSSDocument{s.DetachedContents()[0]}, nil

				} else if sigDMechanism != nil && enumerations.SigDMechanismHTTPHeaders == *sigDMechanism {
					return s.SignedDocumentsByHTTPHeaderName(), nil

				} else if sigDMechanism != nil && enumerations.SigDMechanismObjectIDByURI == *sigDMechanism {
					return s.SignedDocumentsForObjectIdByUriMechanism(), nil
				}
			}
		}

		return originalDocuments, nil

	}
	payloadBytes := s.jws.UnverifiedPayloadBytes()
	return []model.DSSDocument{model.NewInMemoryDocument(payloadBytes)}, nil
}

// DataFoundUpToLevel returns the level up to which the signature has been found conformant.
// Port of getDataFoundUpToLevel().
func (s *Signature) DataFoundUpToLevel() enumerations.SignatureLevel {
	if !s.HasAdESProfile() {
		return enumerations.SignatureLevelJSONNotETSI
	}
	if !s.HasBProfile() {
		return enumerations.SignatureLevelJAdES
	}
	if !s.HasTProfile() {
		return enumerations.SignatureLevelJAdESBaselineB
	}
	if s.HasLTProfile() {
		if s.HasLTAProfile() {
			return enumerations.SignatureLevelJAdESBaselineLTA
		}
		return enumerations.SignatureLevelJAdESBaselineLT
	}
	return enumerations.SignatureLevelJAdESBaselineT
}

// CreateBaselineRequirementsChecker instantiates a BaselineRequirementsChecker according to the
// signature format. Port of the protected createBaselineRequirementsChecker(CertificateVerifier).
func (s *Signature) CreateBaselineRequirementsChecker(certificateVerifier validation.CertificateVerifier) validation.BaselineRequirementsCheckerContract {
	return NewJAdESBaselineRequirementsChecker(s, certificateVerifier)
}

// ValidateStructure processes the structure validation of the signature. Port of the protected
// validateStructure().
func (s *Signature) ValidateStructure() []string {
	validationErrors := DSSJsonUtilsValidateAgainstJAdESSchema(s.jws)
	if utils.IsCollectionNotEmpty(validationErrors) {
		// Upstream logs "Error(s) occurred during the JSON schema validation : {}".
	}
	return validationErrors
}

// FindSignatureScopes finds signature scopes. Port of the protected findSignatureScopes().
func (s *Signature) FindSignatureScopes() []scope.SignatureScope {
	return NewJAdESSignatureScopeFinder().FindSignatureScope(s)
}

// AddExternalTimestamp is not supported for JAdES. Port of addExternalTimestamp(TimestampToken).
//
// Panics with the Java message (UnsupportedOperationException).
func (s *Signature) AddExternalTimestamp(timestamp *validation.TimestampToken) {
	panic("The method addExternalTimestamp(timestamp) is not supported for JAdES!")
}

// compile-time assertion: a Signature satisfies its own overrides contract.
var _ validation.DefaultAdvancedSignatureOverrides = (*Signature)(nil)
