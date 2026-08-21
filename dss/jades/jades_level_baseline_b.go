// Ported from dss-jades/src/main/java/eu/europa/esig/dss/jades/signature/JAdESLevelBaselineB.java (DSS 6.5.RC1).
//
// This class builds the JWS protected header of a JAdES signature according to ETSI TS 119 182-1.
// The serialized protected header is what gets SIGNED, so the byte-for-byte contract of this file
// is the whole point of the chunk:
//
//   - Java accumulates the header in a LinkedHashMap and jose4j serializes it with
//     JsonUtil.toJson, i.e. in INSERTION order. A Go map would randomize that and change the
//     signed bytes, so `signedProperties` is a jose.Object (the insertion-ordered counterpart of
//     LinkedHashMap) and every incorporate* method below runs, and inserts, in exactly the Java
//     statement order of getSignedProperties(): alg, cty, kid, x5u, x5t#S256 / x5t#o, x5c, typ,
//     exp, b64, iat / sigT, srCms, sigPl, srAts, adoTst, sigPId, sigD, crit - with crit last,
//     because it is computed from the keys already present.
//   - Every nested object Java builds as `new LinkedHashMap<>()` wrapped in `new JsonObject(map)`
//     is built here as a jose.NewObject() filled by Put in the same statement order and wrapped
//     with NewJsonObjectFromMap. The ONE site where upstream uses the no-argument `new
//     JsonObject()` - the non-JSON CommitmentTypeQualifier - keeps NewJsonObject(), which wraps a
//     HashMap-ordered object, because that is what upstream serializes there.
//   - org.jose4j.json.internal.json_simple.JSONArray becomes []any, which the ordered-JSON writer
//     of internal/jose serializes exactly as json_simple serializes a Collection. Java's one
//     `int[]` value (UserNotice noticeNumbers) is widened to []any of ints, which json_simple
//     writes identically ("[1,2,3]") through its int[] branch.
//
// The org.jose4j.jwx.HeaderParameterNames constants come from internal/jose, whose Header*
// constants carry those verbatim values.
//
// # Errors
//
// Objects.requireNonNull becomes a panic carrying the Java message; every Java throw
// (IllegalArgumentException, DSSException, UnsupportedOperationException, IllegalInputException)
// becomes a returned error, so each incorporate* method and the two public entry points carry an
// error. slf4j logging is dropped.
package jades

import (
	"errors"
	"fmt"
	"strings"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/internal/jose"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/spi/exception"
	"github.com/ryftcore/dss-go/dss/spi/validation"
	"github.com/ryftcore/dss-go/dss/utils"
)

// JAdESLevelBaselineB builds a JOSE header according to TS 119-182.
type JAdESLevelBaselineB struct {
	// certificateVerifier is the CertificateVerifier to use.
	certificateVerifier validation.CertificateVerifier

	// parameters holds the signature parameters.
	parameters *JAdESSignatureParameters

	// documentsToSign is the list of documents to sign.
	documentsToSign []model.DSSDocument

	// signedProperties is the JOSE header map representation. Java's LinkedHashMap becomes an
	// insertion-ordered jose.Object: see the file header.
	signedProperties *jose.Object
}

// NewJAdESLevelBaselineB is the default constructor.
// Port of JAdESLevelBaselineB(CertificateVerifier, JAdESSignatureParameters, List<DSSDocument>).
func NewJAdESLevelBaselineB(certificateVerifier validation.CertificateVerifier,
	parameters *JAdESSignatureParameters, documentsToSign []model.DSSDocument) (*JAdESLevelBaselineB, error) {
	if certificateVerifier == nil {
		panic("certificateVerifier must not be null!")
	}
	// Upstream repeats the same null-check with the "signatureParameters must be defined!"
	// message against certificateVerifier (a known upstream copy-paste); reproduced verbatim.
	if certificateVerifier == nil {
		panic("signatureParameters must be defined!")
	}
	if utils.IsCollectionEmpty(documentsToSign) {
		return nil, errors.New("Documents to sign must be provided!")
	}
	return &JAdESLevelBaselineB{
		certificateVerifier: certificateVerifier,
		parameters:          parameters,
		documentsToSign:     documentsToSign,
		signedProperties:    jose.NewObject(),
	}, nil
}

// SignedProperties returns the map representing the signed header of a signature.
// Port of #getSignedProperties.
func (b *JAdESLevelBaselineB) SignedProperties() (*jose.Object, error) {
	// RFC 7515 headers
	if err := b.IncorporateSignatureAlgorithm(); err != nil {
		return nil, err
	}
	if err := b.IncorporateContentType(); err != nil {
		return nil, err
	}
	b.IncorporateKeyIdentifier()
	b.IncorporateSigningCertificateUri()
	if err := b.IncorporateSigningCertificate(); err != nil {
		return nil, err
	}
	if err := b.IncorporateCertificateChain(); err != nil {
		return nil, err
	}
	if err := b.IncorporateType(); err != nil {
		return nil, err
	}

	// RFC 7519 headers
	b.incorporateExpirationTime()

	// RFC 7797
	if err := b.IncorporateB64(); err != nil {
		return nil, err
	}

	// TS 119-182 headers
	if err := b.IncorporateSigningTime(); err != nil {
		return nil, err
	}
	b.IncorporateX509CertificateDigests()
	if err := b.IncorporateSignerCommitments(); err != nil {
		return nil, err
	}
	b.IncorporateSignatureProductionPlace()
	b.IncorporateSignerAttributes()
	if err := b.IncorporateContentTimestamps(); err != nil {
		return nil, err
	}
	if err := b.IncorporateSignaturePolicy(); err != nil {
		return nil, err
	}
	if err := b.IncorporateDetachedContents(); err != nil {
		return nil, err
	}

	// must be executed the last
	b.IncorporateCritical()

	return b.signedProperties, nil
}

// IncorporateSignatureAlgorithm incorporates 5.1.2 the alg (X.509 URL) header parameter.
// Port of the protected #incorporateSignatureAlgorithm.
func (b *JAdESLevelBaselineB) IncorporateSignatureAlgorithm() error {
	id := b.parameters.SignatureAlgorithm().JWAID()
	if utils.IsStringNotEmpty(id) {
		b.AddHeader(jose.HeaderAlgorithm, id)
		return nil
	}
	return fmt.Errorf("The defined signature algorithm '%s' is not supported!",
		b.parameters.SignatureAlgorithm())
}

// IncorporateContentType incorporates 5.1.3 the cty (content type) header parameter.
// Port of the protected #incorporateContentType.
func (b *JAdESLevelBaselineB) IncorporateContentType() error {
	if enumerations.SignaturePackaging_DETACHED == b.parameters.SignaturePackaging() &&
		b.parameters.ContentType() == "" {
		// SHOULD NOT be used for detached signatures (see EN 119-182 ch.5.1.3)
		return nil
	}
	mimeTypeString := b.parameters.ContentType()
	if mimeTypeString == "" {
		mimeType := b.documentsToSign[0].MimeType()
		if mimeType != nil {
			mimeTypeString = mimeType.MimeTypeString()
		}
	}
	if mimeTypeString != "" {
		conformant, err := b.rfc7515ConformantMimeTypeString(mimeTypeString)
		if err != nil {
			return err
		}
		b.AddHeader(jose.HeaderContentType, conformant)
	}
	return nil
}

// rfc7515ConformantMimeTypeString ports the private getRFC7515ConformantMimeTypeString.
func (b *JAdESLevelBaselineB) rfc7515ConformantMimeTypeString(mimeTypeString string) (string, error) {
	/*
	 * RFC 7515 :
	 * To keep messages compact in common situations, it is RECOMMENDED that
	 * producers omit an "application/" prefix of a media type value in a
	 * "cty" Header Parameter when no other '/' appears in the media type
	 * value.
	 */
	shortMimeTypeString, err := spi.DSSUtilsStripFirstLeadingOccurrence(mimeTypeString,
		DSSJsonUtilsMimeTypeApplicationPrefix)
	if err != nil {
		return "", err
	}
	if !strings.Contains(shortMimeTypeString, "/") {
		return shortMimeTypeString, nil
	}
	// return original if contains other '/'
	return mimeTypeString, nil
}

// IncorporateKeyIdentifier incorporates 5.1.4 the kid (key identifier) header parameter.
// Port of the protected #incorporateKeyIdentifier.
func (b *JAdESLevelBaselineB) IncorporateKeyIdentifier() {
	if b.parameters.IsIncludeKeyIdentifier() {
		kid := b.parameters.KeyIdentifier()
		if kid == "" && b.parameters.SigningCertificate() != nil {
			issuerSerial := spi.DSSUtilsGenerateKid(b.parameters.SigningCertificate())
			kid = utils.ToBase64(issuerSerial)
		}
		if kid != "" {
			b.AddHeader(jose.HeaderKeyID, kid)
		}
	}
}

// IncorporateSigningCertificateUri incorporates 5.1.5 the x5u (X.509 URL) header parameter.
// Port of the protected #incorporateSigningCertificateUri.
func (b *JAdESLevelBaselineB) IncorporateSigningCertificateUri() {
	x509Url := b.parameters.X509Url()
	if utils.IsStringNotEmpty(x509Url) {
		b.AddHeader(jose.HeaderX509URL, x509Url)
	}
}

// IncorporateSigningCertificate incorporates 5.1.7 the x5t#S256 (X.509 Certificate SHA-256
// Thumbprint) header parameter or 5.2.2 the x5t#o (X509 certificate digest) header parameter.
// Port of the protected #incorporateSigningCertificate.
func (b *JAdESLevelBaselineB) IncorporateSigningCertificate() error {
	signingCertificate := b.parameters.SigningCertificate()
	if signingCertificate == nil {
		return nil
	}

	signingCertificateDigestMethod := b.parameters.SigningCertificateDigestMethod()
	if enumerations.DigestAlgorithm_SHA256 == signingCertificateDigestMethod {
		return b.IncorporateSigningCertificateSha256Thumbprint(signingCertificate)
	}
	return b.IncorporateSigningCertificateOtherDigestReference(signingCertificate,
		signingCertificateDigestMethod)
}

// IncorporateSigningCertificateSha256Thumbprint incorporates 5.1.7 the x5t#S256 (X.509
// Certificate SHA-256 Thumbprint) header parameter.
// Port of the protected #incorporateSigningCertificateSha256Thumbprint.
//
// Upstream delegates to org.jose4j.keys.X509Util.x5tS256(X509Certificate), whose whole body is
// base64url(SHA-256(certificate.getEncoded())) - there is no jose4j state or configuration
// involved, so it is spelled out here over CertificateToken#getDigest(SHA256), which digests
// exactly that DER encoding. Java's returned error channel: getDigest can fail, which Java's
// version cannot express, so the signature carries an error.
func (b *JAdESLevelBaselineB) IncorporateSigningCertificateSha256Thumbprint(
	signingCertificate *model.CertificateToken) error {
	thumbprint, err := signingCertificate.Digest(enumerations.DigestAlgorithm_SHA256)
	if err != nil {
		return err
	}
	x5tS256 := DSSJsonUtilsToBase64Url(thumbprint)
	b.AddHeader(jose.HeaderX509CertificateSHA256Thumbprint, x5tS256)
	return nil
}

// IncorporateCertificateChain incorporates 5.1.8 the x5c (X.509 Certificate Chain) header
// parameter. Port of the protected #incorporateCertificateChain.
func (b *JAdESLevelBaselineB) IncorporateCertificateChain() error {
	if !b.parameters.IsIncludeCertificateChain() || b.parameters.SigningCertificate() == nil {
		return nil
	}

	certificateSelector := spi.NewBaselineBCertificateSelector(b.parameters.SigningCertificate(),
		b.parameters.CertificateChain()).
		SetTrustAnchorBPPolicy(b.parameters.BLevel().IsTrustAnchorBPPolicy()).
		SetTrustedCertificateSource(b.certificateVerifier.TrustedCertSources())
	certificates, err := certificateSelector.Certificates()
	if err != nil {
		return err
	}

	base64Certificates := make([]any, 0, len(certificates))
	for _, certificateToken := range certificates {
		base64Certificates = append(base64Certificates, utils.ToBase64(certificateToken.Encoded()))
	}
	b.AddHeader(jose.HeaderX509CertificateChain, base64Certificates)
	return nil
}

// IncorporateCritical incorporates 5.1.9 the crit (critical) header parameter.
// Port of the protected #incorporateCritical.
func (b *JAdESLevelBaselineB) IncorporateCritical() {
	/*
	 * ETSI TS 119 182-1, 5.1.9	The crit (critical) header parameter
	 *
	 * If the JAdES signature includes the sigD header parameter,
	 * the crit header parameter shall also be present and "sigD" shall be one of its JSON array
	 * elements.
	 */
	criticalHeaderNames := make([]any, 0)
	for _, header := range b.signedProperties.Keys() {
		if DSSJsonUtilsIsRequiredCriticalHeader(header) {
			criticalHeaderNames = append(criticalHeaderNames, header)
		}
	}
	if utils.IsCollectionNotEmpty(criticalHeaderNames) {
		b.AddHeader(jose.HeaderCritical, criticalHeaderNames)
	}
}

// IncorporateType incorporates RFC 7515 : 4.1.9. "typ" (Type) Header Parameter.
// Port of the protected #incorporateType.
func (b *JAdESLevelBaselineB) IncorporateType() error {
	if !b.parameters.IsIncludeSignatureType() {
		return nil
	}

	signatureType := b.parameters.SignatureType()
	if utils.IsStringEmpty(signatureType) {
		/*
		 * RFC 7515 : 4.1.9. "typ" (Type) Header Parameter
		 *
		 * The "typ" value "JOSE" can be used by applications to indicate that
		 * this object is a JWS or JWE using the JWS Compact Serialization or
		 * the JWE Compact Serialization.  The "typ" value "JOSE+JSON" can be
		 * used by applications to indicate that this object is a JWS or JWE
		 * using the JWS JSON Serialization or the JWE JSON Serialization.
		 */

		var signatureMimeType enumerations.MimeType
		switch b.parameters.JwsSerializationType() {
		case enumerations.JWSSerializationType_COMPACT_SERIALIZATION:
			signatureMimeType = enumerations.MimeTypeEnum_JOSE
		case enumerations.JWSSerializationType_JSON_SERIALIZATION,
			enumerations.JWSSerializationType_FLATTENED_JSON_SERIALIZATION:
			signatureMimeType = enumerations.MimeTypeEnum_JOSE_JSON
		default:
			return fmt.Errorf("The given JWS serialization type '%s' is not supported!",
				b.parameters.JwsSerializationType())
		}
		signatureType = signatureMimeType.MimeTypeString()
	}

	headerType, err := b.rfc7515ConformantMimeTypeString(signatureType)
	if err != nil {
		return err
	}
	b.AddHeader(jose.HeaderType, headerType)
	return nil
}

// IncorporateB64 incorporates the RFC 7797 Unencoded Payload Option.
// Port of the protected #incorporateB64.
func (b *JAdESLevelBaselineB) IncorporateB64() error {
	// incorporate only with FALSE value
	if !b.parameters.IsBase64UrlEncodedPayload() {
		if err := b.assertPayloadEncodingValid(); err != nil {
			return err
		}
		b.AddHeader(jose.HeaderBase64URLEncodePayload, b.parameters.IsBase64UrlEncodedPayload())
	}
	return nil
}

// assertPayloadEncodingValid ports the private assertPayloadEncodingValid.
func (b *JAdESLevelBaselineB) assertPayloadEncodingValid() error {
	payloadBytes, err := b.PayloadBytes()
	if err != nil {
		return err
	}
	// see RFC 7797 (only for compact format not detached payload shall be uri-safe)
	if !b.parameters.IsBase64UrlEncodedPayload() &&
		enumerations.SignaturePackaging_DETACHED != b.parameters.SignaturePackaging() &&
		utils.IsArrayNotEmpty(payloadBytes) {

		switch b.parameters.JwsSerializationType() {
		/*
		 * RFC 7797 ch. "5. Unencoded Payload Content Restrictions"
		 */
		case enumerations.JWSSerializationType_COMPACT_SERIALIZATION:
			if !DSSJsonUtilsIsUrlSafePayload(string(payloadBytes)) {
				return exception.NewIllegalInputException("The payload contains not URL-safe characters! " +
					"With Unencoded Payload ('b64' = false) only ASCII characters in ranges " +
					"%x20-2D and %x2F-7E are allowed for a COMPACT_SERIALIZATION!")
			}
		case enumerations.JWSSerializationType_FLATTENED_JSON_SERIALIZATION,
			enumerations.JWSSerializationType_JSON_SERIALIZATION:
			if !DSSJsonUtilsIsUtf8(payloadBytes) {
				return exception.NewIllegalInputException("The payload contains not valid content! " +
					"With Unencoded Payload ('b64' = false) only UTF-8 characters are allowed!")
			}
		default:
			return fmt.Errorf("The JWSSerializationType '%s' is not supported!",
				b.parameters.JwsSerializationType())
		}
	}
	return nil
}

// IncorporateSigningTime incorporates 5.1.11 iat or 5.2.1 sigT (claimed signing time) header
// parameter. Port of the protected #incorporateSigningTime.
func (b *JAdESLevelBaselineB) IncorporateSigningTime() error {
	signingDate := b.parameters.BLevel().SigningDate()
	switch b.parameters.JadesSigningTimeType() {
	case JAdESSigningTimeType_IAT:
		signedTimeInSeconds := spi.DSSUtilsTimeValueInSeconds(signingDate.UnixMilli())
		b.AddHeader(JWTClaimNamesIat, signedTimeInSeconds)
	case JAdESSigningTimeType_SIG_T:
		stringSigningTime := spi.DSSUtilsFormatDateToRFC(*signingDate)
		b.AddHeader(JAdESHeaderParameterNamesSigT, stringSigningTime)
	case JAdESSigningTimeType_NONE:
		// No signing time header to incorporate
	default:
		return fmt.Errorf("The JAdESSigningTimeType '%s' is not supported!",
			b.parameters.JadesSigningTimeType())
	}
	return nil
}

// IncorporateSigningCertificateOtherDigestReference incorporates 5.2.2.2 the x5t#o (X509
// certificate digest) header parameter.
// Port of the protected #incorporateSigningCertificateOtherDigestReference.
func (b *JAdESLevelBaselineB) IncorporateSigningCertificateOtherDigestReference(
	signingCertificate *model.CertificateToken, digestAlgorithm enumerations.DigestAlgorithm) error {
	digestValue, err := signingCertificate.Digest(digestAlgorithm)
	if err != nil {
		return err
	}

	x5toParams := jose.NewObject()
	x5toParams.Put(JAdESHeaderParameterNamesDigAlg, digestAlgorithm.JAdESID())
	x5toParams.Put(JAdESHeaderParameterNamesDigVal, DSSJsonUtilsToBase64Url(digestValue))

	b.AddHeader(JAdESHeaderParameterNamesX5tO, NewJsonObjectFromMap(x5toParams))
	return nil
}

// IncorporateX509CertificateDigests incorporates 5.2.2.3 the sigX5ts (X509 certificates digests).
// Port of the protected #incorporateX509CertificateDigests.
func (b *JAdESLevelBaselineB) IncorporateX509CertificateDigests() {
	// addition of multiple signing certificate references are not supported in DSS
}

// IncorporateSignerCommitments incorporates 5.2.3 the srCms (signer commitments) header
// parameter. Port of the protected #incorporateSignerCommitments.
func (b *JAdESLevelBaselineB) IncorporateSignerCommitments() error {
	if utils.IsCollectionEmpty(b.parameters.BLevel().CommitmentTypeIndications()) {
		return nil
	}

	srCms := make([]any, 0)

	for _, commitmentType := range b.parameters.BLevel().CommitmentTypeIndications() {
		if utils.IsStringEmpty(commitmentType.URI()) && utils.IsStringEmpty(commitmentType.OID()) {
			return errors.New("Either URI or OID shall be defined for CommitmentType signed attribute in JAdES!")
		}

		srCmParams := jose.NewObject()

		// Only simple Oid form is supported
		oidObject, err := DSSJsonUtilsOidObject(commitmentType)
		if err != nil {
			return err
		}
		srCmParams.Put(JAdESHeaderParameterNamesCommId, oidObject)

		commQuals, err := b.commitmentQualifiers(commitmentType)
		if err != nil {
			return err
		}
		if utils.IsCollectionNotEmpty(commQuals) {
			srCmParams.Put(JAdESHeaderParameterNamesCommQuals, commQuals)
		}

		srCms = append(srCms, NewJsonObjectFromMap(srCmParams))
	}

	b.AddHeader(JAdESHeaderParameterNamesSrCms, srCms)
	return nil
}

// commitmentQualifiers ports the private getCommitmentQualifiers. Java returns
// List<JsonObject>, which the caller stores as a plain List value (not a JSONArray); the []any
// used here serializes identically.
func (b *JAdESLevelBaselineB) commitmentQualifiers(
	commitmentType enumerations.CommitmentType) ([]any, error) {
	commQuals := make([]any, 0)
	commonCommitmentType, ok := commitmentType.(*model.CommonCommitmentType)
	if !ok {
		return commQuals, nil
	}
	commitmentQualifiers := commonCommitmentType.CommitmentTypeQualifiers()
	if !utils.IsArrayNotEmpty(commitmentQualifiers) {
		return commQuals, nil
	}
	for _, commitmentQualifier := range commitmentQualifiers {
		if commitmentQualifier == nil {
			panic("CommitmentTypeQualifier cannot be null!")
		}
		content := commitmentQualifier.Content()
		if content == nil {
			return nil, errors.New("CommitmentTypeQualifier content cannot be null!")
		}

		isJsonDocument, err := DSSJsonUtilsIsJsonDocument(content)
		if err != nil {
			return nil, err
		}
		if isJsonDocument {
			binaries, err := spi.DSSUtilsToByteArrayOfDocument(content)
			if err != nil {
				return nil, err
			}
			object, err := DSSJsonUtilsParseJSONStringToMap(string(binaries))
			if err != nil {
				return nil, fmt.Errorf("Unable to parse JSON Commitment Type Qualifier : %s", err.Error())
			}
			commQuals = append(commQuals, NewJsonObjectFromMap(object))

		} else {
			// Upstream logs "None JSON encoded CommitmentTypeQualifier has been provided.
			// Incorporate as JSONObject."
			binaries, err := spi.DSSUtilsToByteArrayOfDocument(content)
			if err != nil {
				return nil, err
			}
			jsonObject := NewJsonObject()
			jsonObject.Put(JAdESHeaderParameterNamesVal, string(binaries))
			commQuals = append(commQuals, jsonObject)
		}
	}
	return commQuals, nil
}

// IncorporateSignatureProductionPlace incorporates 5.2.4 the sigPl (signature production place)
// header parameter. Port of the protected #incorporateSignatureProductionPlace.
func (b *JAdESLevelBaselineB) IncorporateSignatureProductionPlace() {
	signerProductionPlace := b.parameters.BLevel().SignerLocation()
	if signerProductionPlace == nil || signerProductionPlace.IsEmpty() {
		return
	}

	city := signerProductionPlace.Locality()
	streetAddress := signerProductionPlace.StreetAddress()
	stateOrProvince := signerProductionPlace.StateOrProvince()
	postOfficeBoxNumber := signerProductionPlace.PostOfficeBoxNumber()
	postalCode := signerProductionPlace.PostalCode()
	country := signerProductionPlace.Country()

	sigPlaceMap := jose.NewObject()

	// Java tests each member against null; the Go port's SignerLocation answers "" for an unset
	// member, which is the same "absent" state - the type has no null String to distinguish.
	if country != "" {
		sigPlaceMap.Put(JAdESHeaderParameterNamesAddressCountry, country)
	}
	if city != "" {
		sigPlaceMap.Put(JAdESHeaderParameterNamesAddressLocality, city)
	}
	if stateOrProvince != "" {
		sigPlaceMap.Put(JAdESHeaderParameterNamesAddressRegion, stateOrProvince)
	}
	if postOfficeBoxNumber != "" {
		sigPlaceMap.Put(JAdESHeaderParameterNamesPostOfficeBoxNumber, postOfficeBoxNumber)
	}
	if postalCode != "" {
		sigPlaceMap.Put(JAdESHeaderParameterNamesPostalCode, postalCode)
	}
	if streetAddress != "" {
		sigPlaceMap.Put(JAdESHeaderParameterNamesStreetAddress, streetAddress)
	}

	b.AddHeader(JAdESHeaderParameterNamesSigPl, NewJsonObjectFromMap(sigPlaceMap))
}

// IncorporateSignerAttributes incorporates 5.2.5 the srAts (signer attributes) header parameter.
// Port of the protected #incorporateSignerAttributes.
func (b *JAdESLevelBaselineB) IncorporateSignerAttributes() {
	srAtsParams := jose.NewObject()

	// TODO : certified are not supported
	// srAtsParams.Put(JAdESHeaderParameterNamesCertified, certifiedList)

	claimedSignerRoles := b.parameters.BLevel().ClaimedSignerRoles()
	if utils.IsCollectionNotEmpty(claimedSignerRoles) {
		srAtsParams.Put(JAdESHeaderParameterNamesClaimed, b.qArray(claimedSignerRoles))
	}

	signedAssertions := b.parameters.BLevel().SignedAssertions()
	if utils.IsCollectionNotEmpty(signedAssertions) {
		srAtsParams.Put(JAdESHeaderParameterNamesSignedAssertions, b.qArray(signedAssertions))
	}

	if srAtsParams.Size() != 0 {
		srAtsParamsObject := NewJsonObjectFromMap(srAtsParams)
		b.AddHeader(JAdESHeaderParameterNamesSrAts, srAtsParamsObject)
	}
}

// qArray ports the private getQArray.
func (b *JAdESLevelBaselineB) qArray(qArrayVals []string) []any {
	qArrays := make([]any, 0)

	/*
	 * Each instance of this type shall be a JSON array whose elements are JSON
	 * objects. Each JSON object shall contain three members, namely:
	 */
	qArrayMap := jose.NewObject()

	/*
	 * a) The mediaType member, which shall contain a string identifying the type of
	 * the signed assertions or the claimed attributes present in qVals member,
	 * and shall meet the requirements specified in clause 8.4 of
	 * draft-handrews-json-schema-validation-01 [20].
	 */

	/*
	 * RFC 2046 "4.1.3. Plain Subtype"
	 *
	 * The simplest and most important subtype of "text" is "plain". This indicates
	 * plain text that does not contain any formatting commands or directives. Plain
	 * text is intended to be displayed "as-is", that is, no interpretation of
	 * embedded formatting commands, font attribute specifications, processing
	 * instructions, interpretation directives, or content markup should be
	 * necessary for proper display.
	 */
	qArrayMap.Put(JAdESHeaderParameterNamesMediaType, enumerations.MimeTypeEnum_TEXT.MimeTypeString())

	/*
	 * b) The encoding member, which shall contain a string identifying the encoding
	 * of the signed assertions or the claimed attributes present in qVals member,
	 * and shall meet the requirements specified in clause 8.3 of
	 * draft-handrews-json-schema-validation-01 [20].
	 */

	/*
	 * RFC 2045 "2.9. Binary Data"
	 *
	 * "Binary data" refers to data where any sequence of octets whatsoever is
	 * allowed.
	 */
	qArrayMap.Put(JAdESHeaderParameterNamesEncoding, DSSJsonUtilsContentEncodingBinary)

	/*
	 * c) The qVals member, which shall be a JSON array of at least one item. The
	 * elements of qVals JSON array shall be the values of the signed assertions or
	 * the claimed attributes encoded as indicated within the encoding member.
	 */
	values := make([]any, 0, len(qArrayVals))
	for _, value := range qArrayVals {
		values = append(values, value)
	}
	qArrayMap.Put(JAdESHeaderParameterNamesQVals, values)

	qArray := NewJsonObjectFromMap(qArrayMap)
	qArrays = append(qArrays, qArray)

	return qArrays
}

// IncorporateContentTimestamps incorporates 5.2.6 the adoTst (signed data time-stamp) header
// parameter. Port of the protected #incorporateContentTimestamps.
func (b *JAdESLevelBaselineB) IncorporateContentTimestamps() error {
	if utils.IsCollectionEmpty(b.parameters.ContentTimestamps()) {
		return nil
	}

	// canonicalization shall be null for content timestamps (see 5.2.6)
	contentTimestampBinaries := jadesLevelBaselineBToTimestampBinaries(b.parameters.ContentTimestamps())
	tstContainer, err := DSSJsonUtilsTstContainer(contentTimestampBinaries, "")
	if err != nil {
		return err
	}
	b.AddHeader(JAdESHeaderParameterNamesAdoTst, tstContainer)
	return nil
}

// jadesLevelBaselineBToTimestampBinaries ports the private toTimestampBinaries.
func jadesLevelBaselineBToTimestampBinaries(
	timestampTokens []*validation.TimestampToken) []*model.TimestampBinary {
	if utils.IsCollectionEmpty(timestampTokens) {
		return []*model.TimestampBinary{}
	}
	timestampBinaries := make([]*model.TimestampBinary, 0, len(timestampTokens))
	for _, timestampToken := range timestampTokens {
		timestampBinaries = append(timestampBinaries, model.NewTimestampBinary(timestampToken.Encoded()))
	}
	return timestampBinaries
}

// IncorporateSignaturePolicy incorporates 5.2.7 the sigPId (signature policy identifier) header
// parameter. Port of the protected #incorporateSignaturePolicy.
func (b *JAdESLevelBaselineB) IncorporateSignaturePolicy() error {
	signaturePolicy := b.parameters.BLevel().SignaturePolicy()
	if signaturePolicy == nil || signaturePolicy.IsEmpty() {
		return nil
	}
	if err := jadesLevelBaselineBAssertSignaturePolicyValid(signaturePolicy); err != nil {
		return err
	}

	sigPIdParams := jose.NewObject()

	signaturePolicyId := signaturePolicy.Id()
	oid := DSSJsonUtilsOidObjectFromURI(signaturePolicyId, signaturePolicy.Description(),
		signaturePolicy.DocumentationReferences())
	sigPIdParams.Put(JAdESHeaderParameterNamesId, oid)

	if signaturePolicy.DigestAlgorithm() != "" && signaturePolicy.DigestValue() != nil {
		sigPIdParams.Put(JAdESHeaderParameterNamesDigAlg, signaturePolicy.DigestAlgorithm().JAdESID())
		sigPIdParams.Put(JAdESHeaderParameterNamesDigVal,
			DSSJsonUtilsToBase64Url(signaturePolicy.DigestValue()))
	}

	/*
	 * The hashPSp digPSp member shall be a boolean. When present and set to "true",
	 * it shall indicate that the digest of the signature policy document has been
	 * computed as specified in a technical specification. Absence of this member
	 * shall be considered as if present and set to "false". If this member is
	 * present and set to "true", then the qualifier spDSpec qualifier shall be
	 * present and shall identify the aforementioned technical specification.
	 */
	if signaturePolicy.IsHashAsInTechnicalSpecification() {
		sigPIdParams.Put(JAdESHeaderParameterNamesDigPSp, signaturePolicy.IsHashAsInTechnicalSpecification())
	}

	if signaturePolicy.IsSPQualifierPresent() {
		signaturePolicyQualifiers := jadesLevelBaselineBSignaturePolicyQualifiers(signaturePolicy)
		sigPIdParams.Put(JAdESHeaderParameterNamesSigPQuals, signaturePolicyQualifiers)
	}

	b.AddHeader(JAdESHeaderParameterNamesSigPid, NewJsonObjectFromMap(sigPIdParams))
	return nil
}

// jadesLevelBaselineBAssertSignaturePolicyValid ports the private assertSignaturePolicyValid.
func jadesLevelBaselineBAssertSignaturePolicyValid(signaturePolicy *model.Policy) error {
	if utils.IsStringEmpty(signaturePolicy.Id()) {
		// see TS 119-182 ch. 5.2.7.1 Semantics and syntax ('id' is required)
		return errors.New("Implicit policy is not allowed in JAdES! The signaturePolicyId attribute is required!")
	}
	if signaturePolicy.IsHashAsInTechnicalSpecification() &&
		(signaturePolicy.SpDocSpecification() == nil ||
			utils.IsStringEmpty(signaturePolicy.SpDocSpecification().Id())) {
		return errors.New("SpDocSpecification shall be defined when DigestAsInTechnicalSpecification is set to true!")
	}
	return nil
}

// jadesLevelBaselineBSignaturePolicyQualifiers ports the private getSignaturePolicyQualifiers.
func jadesLevelBaselineBSignaturePolicyQualifiers(signaturePolicy *model.Policy) []any {
	sigPQualifiers := make([]any, 0)
	/*
	 * NOTE: Intermediate objects are created in order to allow multiple instances of the same
	 * qualifiers
	 *
	 * TS 119-182 ch. 5.2.7.1 Semantics and syntax:
	 * The sigPQuals member may contain one or more qualifiers of the same type.
	 */
	spuri := signaturePolicy.Spuri()
	if utils.IsStringNotEmpty(spuri) {
		qualifier := jose.NewObject()
		qualifier.Put(JAdESHeaderParameterNamesSpURI, spuri)
		sigPQualifiers = append(sigPQualifiers, NewJsonObjectFromMap(qualifier))
	}

	userNotice := signaturePolicy.UserNotice()
	if userNotice != nil && !userNotice.IsEmpty() {
		spUserNotice := jose.NewObject()

		organization := userNotice.Organization()
		noticeNumbers := userNotice.NoticeNumbers()
		if utils.IsStringNotEmpty(organization) && noticeNumbers != nil && len(noticeNumbers) > 0 {
			noticeRef := jose.NewObject()
			noticeRef.Put(JAdESHeaderParameterNamesOrgantization, organization)
			// Java stores the raw int[]; json_simple writes it through its int[] branch, which
			// produces the same bytes as a list of the same numbers.
			numbers := make([]any, 0, len(noticeNumbers))
			for _, noticeNumber := range noticeNumbers {
				numbers = append(numbers, noticeNumber)
			}
			noticeRef.Put(JAdESHeaderParameterNamesNoticeNumbers, numbers)
			spUserNotice.Put(JAdESHeaderParameterNamesNoticeRef, NewJsonObjectFromMap(noticeRef))
		}

		explicitText := userNotice.ExplicitText()
		if utils.IsStringNotEmpty(explicitText) {
			spUserNotice.Put(JAdESHeaderParameterNamesExplText, explicitText)
		}

		qualifier := jose.NewObject()
		qualifier.Put(JAdESHeaderParameterNamesSpUserNotice, NewJsonObjectFromMap(spUserNotice))
		sigPQualifiers = append(sigPQualifiers, NewJsonObjectFromMap(qualifier))
	}

	spDocSpecification := signaturePolicy.SpDocSpecification()
	if spDocSpecification != nil && utils.IsStringNotEmpty(spDocSpecification.Id()) {
		spDSpec := DSSJsonUtilsOidObjectFromURI(spDocSpecification.Id(),
			spDocSpecification.Description(), spDocSpecification.DocumentationReferences())

		qualifier := jose.NewObject()
		qualifier.Put(JAdESHeaderParameterNamesSpDspec, spDSpec)
		sigPQualifiers = append(sigPQualifiers, NewJsonObjectFromMap(qualifier))
	}

	return sigPQualifiers
}

// IncorporateDetachedContents incorporates 5.2.8 the sigD header parameter.
// Port of the protected #incorporateDetachedContents.
func (b *JAdESLevelBaselineB) IncorporateDetachedContents() error {
	if enumerations.SignaturePackaging_DETACHED != b.parameters.SignaturePackaging() {
		return nil
	}
	if err := b.assertDetachedContentValid(); err != nil {
		return err
	}

	var sigDParams *jose.Object
	var err error
	switch b.parameters.SigDMechanism() {
	case enumerations.SigDMechanism_HTTP_HEADERS:
		// 5.2.8.2 Mechanism HttpHeaders
		if err = b.assertHttpHeadersConfigurationValid(); err != nil {
			return err
		}
		sigDParams = b.sigDForHttpHeadersMechanism(b.documentsToSign)
	case enumerations.SigDMechanism_OBJECT_ID_BY_URI:
		// 5.2.8.3.2 Mechanism ObjectIdByURI
		sigDParams, err = b.sigDForObjectIdByUriMechanism(b.documentsToSign)
	case enumerations.SigDMechanism_OBJECT_ID_BY_URI_HASH:
		// 5.2.8.3.3 Mechanism ObjectIdByURIHash
		sigDParams, err = b.sigDForObjectIdByUriHashMechanism(b.documentsToSign)
	case enumerations.SigDMechanism_NO_SIG_D:
		// do not incorporate the SigD
		return nil
	default:
		return fmt.Errorf("The 'sigD' mechanism '%s' is not supported!", b.parameters.SigDMechanism())
	}
	if err != nil {
		return err
	}

	b.AddHeader(JAdESHeaderParameterNamesSigD, NewJsonObjectFromMap(sigDParams))
	return nil
}

// assertDetachedContentValid ports the private assertDetachedContentValid.
func (b *JAdESLevelBaselineB) assertDetachedContentValid() error {
	sigDMechanism := b.parameters.SigDMechanism()
	if sigDMechanism == "" {
		return errors.New("The SigDMechanism is not defined for a detached signature! " +
			"Please use JAdESSignatureParameters.setSigDMechanism(sigDMechanism) method.")
	}
	if enumerations.SigDMechanism_NO_SIG_D == sigDMechanism {
		if utils.CollectionSize(b.documentsToSign) > 1 {
			return fmt.Errorf("Only one detached document is allowed with '%s' mechanism!",
				enumerations.SigDMechanism_NO_SIG_D)
		}
		return nil
	}

	documentNames := make([]string, 0)
	for _, document := range b.documentsToSign {
		if utils.IsStringEmpty(document.Name()) {
			return errors.New("The signed document must have names for a detached JAdES signature!")
		}
		if enumerations.SigDMechanism_HTTP_HEADERS != sigDMechanism {
			for _, name := range documentNames {
				if name == document.Name() {
					return fmt.Errorf("The documents to be signed shall have different names! "+
						"The name '%s' appears multiple times.", document.Name())
				}
			}
		}
		documentNames = append(documentNames, document.Name())
	}
	return nil
}

// assertHttpHeadersConfigurationValid ports the private assertHttpHeadersConfigurationValid.
func (b *JAdESLevelBaselineB) assertHttpHeadersConfigurationValid() error {
	/*
	 * 5.1.10 The b64 header parameter
	 *
	 * If the sigD header parameter is present with its member set to
	 * "http://uri.etsi.org/19182/HttpHeaders" then the b64 header parameter shall
	 * be present and set to "false".
	 */
	if enumerations.SigDMechanism_HTTP_HEADERS == b.parameters.SigDMechanism() &&
		b.parameters.IsBase64UrlEncodedPayload() {
		return fmt.Errorf("'%s' SigD Mechanism can be used only with non-base64url encoded payload! "+
			"Set JAdESSignatureParameters.setBase64UrlEncodedPayload(false).",
			enumerations.SigDMechanism_HTTP_HEADERS.JAdESUri())
	}
	return nil
}

// sigDForHttpHeadersMechanism ports the private getSigDForHttpHeadersMechanism.
func (b *JAdESLevelBaselineB) sigDForHttpHeadersMechanism(detachedContents []model.DSSDocument) *jose.Object {
	sigDParams := jose.NewObject()

	sigDParams.Put(JAdESHeaderParameterNamesMId, enumerations.SigDMechanism_HTTP_HEADERS.JAdESUri())
	sigDParams.Put(JAdESHeaderParameterNamesPars, jadesLevelBaselineBHttpHeaderNames(detachedContents))

	return sigDParams
}

// sigDForObjectIdByUriMechanism ports the private getSigDForObjectIdByUriMechanism.
func (b *JAdESLevelBaselineB) sigDForObjectIdByUriMechanism(
	detachedContents []model.DSSDocument) (*jose.Object, error) {
	sigDParams := jose.NewObject()

	sigDParams.Put(JAdESHeaderParameterNamesMId, enumerations.SigDMechanism_OBJECT_ID_BY_URI.JAdESUri())
	sigDParams.Put(JAdESHeaderParameterNamesPars, jadesLevelBaselineBSignedDataReferences(detachedContents))

	ctys, err := b.signedDataMimeTypesIfPresent(detachedContents)
	if err != nil {
		return nil, err
	}
	sigDParams.Put(JAdESHeaderParameterNamesCtys, ctys)

	return sigDParams, nil
}

// sigDForObjectIdByUriHashMechanism ports the private getSigDForObjectIdByUriHashMechanism.
func (b *JAdESLevelBaselineB) sigDForObjectIdByUriHashMechanism(
	detachedContents []model.DSSDocument) (*jose.Object, error) {
	sigDParams := jose.NewObject()

	sigDParams.Put(JAdESHeaderParameterNamesMId, enumerations.SigDMechanism_OBJECT_ID_BY_URI_HASH.JAdESUri())
	sigDParams.Put(JAdESHeaderParameterNamesPars, jadesLevelBaselineBSignedDataReferences(detachedContents))

	digestAlgorithm := b.referenceDigestAlgorithmOrDefault()
	sigDParams.Put(JAdESHeaderParameterNamesHashM, digestAlgorithm.JAdESID())
	hashV, err := b.signedDataDigests(detachedContents, digestAlgorithm)
	if err != nil {
		return nil, err
	}
	sigDParams.Put(JAdESHeaderParameterNamesHashV, hashV)

	ctys, err := b.signedDataMimeTypesIfPresent(detachedContents)
	if err != nil {
		return nil, err
	}
	sigDParams.Put(JAdESHeaderParameterNamesCtys, ctys)

	return sigDParams, nil
}

// jadesLevelBaselineBSignedDataReferences ports the private getSignedDataReferences.
func jadesLevelBaselineBSignedDataReferences(detachedContents []model.DSSDocument) []any {
	references := make([]any, 0, len(detachedContents))
	for _, document := range detachedContents {
		references = append(references, document.Name())
	}
	return references
}

// referenceDigestAlgorithmOrDefault ports the private getReferenceDigestAlgorithmOrDefault.
func (b *JAdESLevelBaselineB) referenceDigestAlgorithmOrDefault() enumerations.DigestAlgorithm {
	if b.parameters.ReferenceDigestAlgorithm() != "" {
		return b.parameters.ReferenceDigestAlgorithm()
	}
	return b.parameters.DigestAlgorithm()
}

// signedDataDigests ports the private getSignedDataDigests.
func (b *JAdESLevelBaselineB) signedDataDigests(detachedContents []model.DSSDocument,
	digestAlgorithm enumerations.DigestAlgorithm) ([]any, error) {
	/*
	 * The hashV member shall be a non-empty array of strings. Each element of the
	 * array shall contain:
	 */
	digests := make([]any, 0, len(detachedContents))
	for _, document := range detachedContents {
		var docDigest []byte
		var err error
		_, isDigestDocument := document.(*model.DigestDocument)
		/*
		 * 1) The base64url-encoded digest value of the data object referenced by the
		 * parameter value (...) if the b64 header parameter is present and set to
		 * "false".
		 */
		if !b.parameters.IsBase64UrlEncodedPayload() || isDigestDocument {
			docDigest, err = document.DigestValue(digestAlgorithm)
			if err != nil {
				return nil, err
			}
		} else {
			/*
			 * 2) The base64url-encoded digest value of the base64url-encoded data object
			 * referenced by the parameter value (...) if the b64 header parameter is absent
			 * or it is present and set to "true".
			 */
			base64urlDocumentContent, err := DSSJsonUtilsToBase64UrlDocument(document)
			if err != nil {
				return nil, err
			}
			docDigest, err = spi.DSSUtilsDigest(digestAlgorithm, []byte(base64urlDocumentContent))
			if err != nil {
				return nil, err
			}
		}
		digests = append(digests, DSSJsonUtilsToBase64Url(docDigest)) // base64Url digest
	}
	return digests, nil
}

// signedDataMimeTypesIfPresent returns a 'ctys' array for the given documents.
// Port of the private getSignedDataMimeTypesIfPresent.
func (b *JAdESLevelBaselineB) signedDataMimeTypesIfPresent(
	detachedContents []model.DSSDocument) ([]any, error) {
	mimeTypes := make([]any, 0, len(detachedContents))
	for _, document := range detachedContents {
		mimeType := document.MimeType()
		if mimeType == nil {
			mimeType = enumerations.MimeTypeEnum_BINARY
		}
		rfc7515MimeType, err := b.rfc7515ConformantMimeTypeString(mimeType.MimeTypeString())
		if err != nil {
			return nil, err
		}
		mimeTypes = append(mimeTypes, rfc7515MimeType)
	}
	return mimeTypes, nil
}

// jadesLevelBaselineBHttpHeaderNames returns the list of HTTP message field names being included
// into 'sigD' for the HttpHeaders mechanism. Port of the private getHttpHeaderNames.
func jadesLevelBaselineBHttpHeaderNames(detachedContents []model.DSSDocument) []any {
	/*
	 * TS 119 182-1 "5.2.8.2 Mechanism HttpHeaders" :
	 *
	 * For this referencing mechanism, the contents of the pars member
	 * shall be an array of lowercased names of HTTP header fields, each one
	 * with the semantics and syntax specified in clause
	 * 2.1.3 of draft-cavage-http-signatures-10: "Signing HTTP Messages" [17].
	 */
	httpHeaderNames := make([]any, 0)

	for _, document := range detachedContents {
		// Java's `document instanceof HTTPHeader` also matches the HTTPHeaderDigest subclass, so
		// the 'Digest' header contributes its lowercased name to 'pars' like any other; Go's
		// type assertion does not follow embedding, hence the explicit two-case switch.
		switch document.(type) {
		case *HTTPHeader, *HTTPHeaderDigest:
		default:
			continue
		}
		headerName := utils.LowerCase(document.Name())
		found := false
		for _, existing := range httpHeaderNames {
			if existing == headerName {
				found = true
				break
			}
		}
		if !found {
			httpHeaderNames = append(httpHeaderNames, headerName)
		}
	}

	return httpHeaderNames
}

// incorporateExpirationTime incorporates RFC 7519 : 4.1.4. "exp" (Expiration Time) Claim.
// Port of the private incorporateExpirationTime.
func (b *JAdESLevelBaselineB) incorporateExpirationTime() {
	if b.parameters.ExpirationTime() != nil {
		expirationTimeInSeconds := spi.DSSUtilsTimeValueInSeconds(b.parameters.ExpirationTime().UnixMilli())
		b.AddHeader(JWTClaimNamesExp, expirationTimeInSeconds)
	}
}

// AddHeader adds a new header to the signedProperties map.
// Port of the protected #addHeader.
func (b *JAdESLevelBaselineB) AddHeader(headerName string, value any) {
	b.signedProperties.Put(headerName, value)
}

// PayloadBytes returns the JWS payload for the given signature parameters.
// Port of #getPayloadBytes.
func (b *JAdESLevelBaselineB) PayloadBytes() ([]byte, error) {
	if enumerations.SignaturePackaging_DETACHED != b.parameters.SignaturePackaging() ||
		enumerations.SigDMechanism_NO_SIG_D == b.parameters.SigDMechanism() {
		return b.incorporatedPayload()

	} else if enumerations.SigDMechanism_HTTP_HEADERS == b.parameters.SigDMechanism() {
		return b.payloadForHttpHeadersMechanism()

	} else if enumerations.SigDMechanism_OBJECT_ID_BY_URI == b.parameters.SigDMechanism() {
		return b.payloadForObjectIdByUriMechanism()

	} else if enumerations.SigDMechanism_OBJECT_ID_BY_URI_HASH == b.parameters.SigDMechanism() {
		/*
		 * 5.2.8.3.3 Mechanism ObjectIdByURIHash
		 *
		 * When using this mechanism, the JWS Payload shall contribute as an empty
		 * stream to the computation of the JWS Signature Value.
		 */
		return spi.DSSUtilsEmptyByteArray, nil
	}
	return nil, errors.New("The configured signature format is not supported!")
}

// incorporatedPayload ports the private getIncorporatedPayload.
func (b *JAdESLevelBaselineB) incorporatedPayload() ([]byte, error) {
	return DSSJsonUtilsDocumentOctets(b.documentsToSign[0], b.parameters.IsBase64UrlEncodedPayload())
}

// payloadForHttpHeadersMechanism ports the private getPayloadForHttpHeadersMechanism.
func (b *JAdESLevelBaselineB) payloadForHttpHeadersMechanism() ([]byte, error) {
	httpHeadersPayloadBuilder := NewHttpHeadersPayloadBuilder(b.documentsToSign, false)
	return httpHeadersPayloadBuilder.Build()
}

// payloadForObjectIdByUriMechanism ports the private getPayloadForObjectIdByUriMechanism.
func (b *JAdESLevelBaselineB) payloadForObjectIdByUriMechanism() ([]byte, error) {
	// NOTE: base64url encoding is processed by JWS
	return DSSJsonUtilsConcatenateDSSDocuments(b.documentsToSign, b.parameters.IsBase64UrlEncodedPayload())
}
