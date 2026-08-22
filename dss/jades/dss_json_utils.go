// Ported from dss-jades/src/main/java/eu/europa/esig/dss/jades/DSSJsonUtils.java (DSS 6.5.RC1).
//
// Naming follows the DSSUtils/DSSASN1Utils precedent already in the tree: a static method of a
// utility class becomes <ClassName><MethodName> with a leading get dropped
// (getSigningInputBytes -> DSSJsonUtilsSigningInputBytes), and an overload keeps a suffix naming
// what distinguishes it (DSSUtilsObjectIdentifierValueWithQualifier).
//
// Java types map as follows, which is what makes the signatures below readable:
//
//	Map<String, Object> / LinkedHashMap  ->  *jose.Object  (insertion-ordered; see internal/jose)
//	org.json.simple JSONArray / List<?>  ->  []any
//	java.lang.Number                     ->  *jose.Number  (Long/BigInteger/Double are distinct)
//	java.lang.Boolean                    ->  *bool         (Java's null third state)
//	java.util.Date                       ->  time.Time     (absent is the zero Time)
//
// The specs module is reached through the dss/jades/specs package. Java's slf4j logging is not
// ported anywhere in this file.
package jades

import (
	"bytes"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/internal/jose"
	"github.com/ryftcore/dss-go/dss/jades/specs"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/spi/exception"
	"github.com/ryftcore/dss-go/dss/utils"
)

// DSSJsonUtilsMimeTypeApplicationPrefix is the MimeType application prefix. Port of
// MIME_TYPE_APPLICATION_PREFIX.
const DSSJsonUtilsMimeTypeApplicationPrefix = "application/"

// DSSJsonUtilsHTTPHeaderDigest is the HttpHeader defining the Digest value of a signed message
// body. Port of HTTP_HEADER_DIGEST.
const DSSJsonUtilsHTTPHeaderDigest = "Digest"

// DSSJsonUtilsContentEncodingBinary is the binary content encoding (RFC 2045). Port of
// CONTENT_ENCODING_BINARY.
const DSSJsonUtilsContentEncodingBinary = "binary"

// dssJsonUtilsProtectedCriticalHeaders contains protected header names that are supported and can
// be present in the critical ('crit') attribute. Port of the private static
// protectedCriticalHeaders set.
var dssJsonUtilsProtectedCriticalHeaders = map[string]struct{}{
	// JAdES TS 119 182-1 constraints
	JAdESHeaderParameterNamesSigT:    {},
	JAdESHeaderParameterNamesX5tO:    {},
	JAdESHeaderParameterNamesSigX5tS: {},
	JAdESHeaderParameterNamesSrCms:   {},
	JAdESHeaderParameterNamesSigPl:   {},
	JAdESHeaderParameterNamesSrAts:   {},
	JAdESHeaderParameterNamesAdoTst:  {},
	JAdESHeaderParameterNamesSigPid:  {},
	JAdESHeaderParameterNamesSigD:    {},
	// RFC 7519 'iat'
	JWTClaimNamesIat: {},
	JWTClaimNamesExp: {},
	// RFC 7797 'b64'
	jose.HeaderBase64URLEncodePayload: {},
}

// dssJsonUtilsCriticalHeaderExceptions contains the headers that MUST NOT be incorporated into a
// 'crit' header (RFC 7515, RFC 7518). Port of the private static criticalHeaderExceptions set.
var dssJsonUtilsCriticalHeaderExceptions = map[string]struct{}{
	// RFC 7515
	jose.HeaderAlgorithm:                       {},
	jose.HeaderJWKSetURL:                       {},
	jose.HeaderJWK:                             {},
	jose.HeaderKeyID:                           {},
	jose.HeaderX509URL:                         {},
	jose.HeaderX509CertificateChain:            {},
	jose.HeaderX509CertificateThumbprint:       {},
	jose.HeaderX509CertificateSHA256Thumbprint: {},
	jose.HeaderType:                            {},
	jose.HeaderContentType:                     {},
	jose.HeaderCritical:                        {},
	// RFC 7518
	jose.HeaderEphemeralPublicKey:   {},
	jose.HeaderAgreementPartyUInfo:  {},
	jose.HeaderAgreementPartyVInfo:  {},
	jose.HeaderInitializationVector: {},
	jose.HeaderAuthenticationTag:    {},
	jose.HeaderPBES2SaltInput:       {},
	jose.HeaderPBES2IterationCount:  {},
	jose.HeaderEncryptionMethod:     {},
	jose.HeaderZip:                  {},
}

// dssJsonUtilsRequiredCriticalHeaders contains the headers that are required to be present in the
// critical ('crit') attribute when used. Port of the private static requiredCriticalHeaders set.
var dssJsonUtilsRequiredCriticalHeaders = map[string]struct{}{
	// JAdES TS 119 182-1 constraints
	JAdESHeaderParameterNamesSigD: {},
	// RFC 7797 'b64'
	jose.HeaderBase64URLEncodePayload: {},
}

// DSSJsonUtilsAsciiBytes returns the ASCII-encoded array for str. Port of getAsciiBytes(String).
func DSSJsonUtilsAsciiBytes(str string) []byte {
	return jose.ASCIIBytes(str)
}

// DSSJsonUtilsToBase64Url returns a base64url-encoded string for binary. Port of
// toBase64Url(byte[]).
func DSSJsonUtilsToBase64Url(binary []byte) string {
	return jose.Base64URLEncode(binary)
}

// DSSJsonUtilsToBase64UrlDocument returns a base64url-encoded string for the document's octets.
// Port of toBase64Url(DSSDocument).
func DSSJsonUtilsToBase64UrlDocument(document model.DSSDocument) (string, error) {
	binaries, err := spi.DSSUtilsToByteArrayOfDocument(document)
	if err != nil {
		return "", err
	}
	return DSSJsonUtilsToBase64Url(binaries), nil
}

// DSSJsonUtilsToBase64UrlObject returns a base64url-encoded string of the JSON serialization of a
// JSON object or JSON array. Port of toBase64Url(Object).
//
// Upstream renders through JSONValue.toJSONString and then calls String.getBytes() with the
// platform default charset; jose.JSON already yields the same characters and Go strings are
// UTF-8, which is what every JVM DSS runs on uses. The difference would only show on a JVM
// started with a non-UTF-8 file.encoding, where upstream would produce a signature nobody else
// can verify.
func DSSJsonUtilsToBase64UrlObject(object any) string {
	return jose.Base64URLEncode([]byte(jose.JSON(object)))
}

// DSSJsonUtilsFromBase64Url returns the decoded binary for a base64url-encoded string. Port of
// fromBase64Url(String).
func DSSJsonUtilsFromBase64Url(base64URLEncoded string) []byte {
	return jose.Base64URLDecode(base64URLEncoded)
}

// DSSJsonUtilsIsBase64UrlEncoded checks whether str is base64url-encoded. Port of
// isBase64UrlEncoded(String).
//
// The decode call upstream makes first cannot fail (commons-codec accepts any byte array), so
// what actually decides the answer is the per-character alphabet check - which is why the
// lenient decoder in internal/jose must not be replaced by a strict one: doing so would not
// change this result, but it would change what fromBase64Url makes of the strings this accepts.
func DSSJsonUtilsIsBase64UrlEncoded(str string) bool {
	for i := 0; i < len(str); i++ {
		if !DSSJsonUtilsIsBase64UrlEncodedByte(str[i]) {
			return false
		}
	}
	return true
}

// DSSJsonUtilsIsBase64UrlEncodedByte checks whether the byte is in the base64url alphabet. Port
// of isBase64UrlEncoded(byte), which walks jose4j's URL_SAFE_ENCODE_TABLE.
func DSSJsonUtilsIsBase64UrlEncodedByte(b byte) bool {
	return jose.IsBase64URLCharacter(b)
}

// DSSJsonUtilsIsUrlSafePayload checks whether the payload is JWS URL safe. Port of
// isUrlSafePayload(String), see RFC 7797 section 5.2 "Unencoded JWS Compact Serialization
// Payload".
//
// Upstream writes the test as the regular expression "[^\P{Print}.]*", i.e. "every character is
// printable and is not a period". Java's \p{Print} is POSIX and ASCII-only, so it is exactly
// %x20-7E; the loop below says the same thing without a regexp engine, and agrees with
// DSSJsonUtilsIsUrlSafe on every byte except that this one also admits the period-free range.
func DSSJsonUtilsIsUrlSafePayload(payloadString string) bool {
	for i := 0; i < len(payloadString); i++ {
		c := payloadString[i]
		if c == '.' {
			return false
		}
		if c < 0x20 || c > 0x7e {
			return false
		}
	}
	return true
}

// DSSJsonUtilsIsUrlSafe checks whether the given byte is url-safe. Port of isUrlSafe(byte), see
// RFC 7797 section 5.2: the ranges %x20-2D and %x2F-7E, i.e. printable ASCII without the period.
func DSSJsonUtilsIsUrlSafe(b byte) bool {
	return 0x1f < b && b < 0x2e || 0x2e < b && b < 0x7f
}

// DSSJsonUtilsIsUtf8 checks whether the binaries contain a UTF-8 encoded string. Port of
// isUtf8(byte[]), whose CharsetDecoder rejects malformed and unmappable sequences.
func DSSJsonUtilsIsUtf8(binaries []byte) bool {
	return utf8.Valid(binaries)
}

// DSSJsonUtilsConcatenate concatenates the given strings with a '.' between, e.g. "xxx", "yyy",
// "zzz" to "xxx.yyy.zzz". Port of concatenate(String...).
func DSSJsonUtilsConcatenate(strings ...string) string {
	return jose.CompactSerialize(strings...)
}

// DSSJsonUtilsSupportedProtectedCriticalHeaders returns the set of supported protected critical
// headers. Port of getSupportedProtectedCriticalHeaders().
//
// A copy is returned rather than the package-level map: upstream hands out the very set it holds
// (a latent aliasing bug there), and a caller mutating the copy here cannot change what this
// package considers supported.
func DSSJsonUtilsSupportedProtectedCriticalHeaders() map[string]struct{} {
	out := make(map[string]struct{}, len(dssJsonUtilsProtectedCriticalHeaders))
	for name := range dssJsonUtilsProtectedCriticalHeaders {
		out[name] = struct{}{}
	}
	return out
}

// DSSJsonUtilsIsCriticalHeaderException checks whether headerName is a critical header exception,
// i.e. one that shall not be incorporated within a 'crit' header (RFC 7515). Port of
// isCriticalHeaderException(String).
func DSSJsonUtilsIsCriticalHeaderException(headerName string) bool {
	_, found := dssJsonUtilsCriticalHeaderExceptions[headerName]
	return found
}

// DSSJsonUtilsIsRequiredCriticalHeader checks whether headerName is required to be incorporated
// within a 'crit' header when used. Port of isRequiredCriticalHeader(String).
func DSSJsonUtilsIsRequiredCriticalHeader(headerName string) bool {
	_, found := dssJsonUtilsRequiredCriticalHeaders[headerName]
	return found
}

// DSSJsonUtilsDigest creates a Digest from a JSON structure carrying digAlg and digVal. Port of
// getDigest(Map), which swallows every exception and answers null; the zero Digest with ok=false
// stands in for that null.
func DSSJsonUtilsDigest(digestValueAndAlgo *jose.Object) (model.Digest, bool) {
	if digestValueAndAlgo.Size() == 0 {
		return model.Digest{}, false
	}
	digestAlgoURI := DSSJsonUtilsGetAsString(digestValueAndAlgo, JAdESHeaderParameterNamesDigAlg)
	digestValueBase64 := DSSJsonUtilsGetAsString(digestValueAndAlgo, JAdESHeaderParameterNamesDigVal)
	if utils.IsStringEmpty(digestAlgoURI) || utils.IsStringEmpty(digestValueBase64) {
		return model.Digest{}, false
	}
	digestAlgorithm, err := enumerations.DigestAlgorithmForJAdES(digestAlgoURI)
	if err != nil {
		// "Unable to extract Digest Algorithm and Value" - upstream logs and returns null.
		return model.Digest{}, false
	}
	return model.NewDigest(digestAlgorithm, DSSJsonUtilsFromBase64Url(digestValueBase64)), true
}

// DSSJsonUtilsOidObject creates an 'oid' JsonObject from an ObjectIdentifier, per TS 119 182-1
// clause 5.4.1 "The oId data type". Port of getOidObject(ObjectIdentifier).
//
// Upstream lets the delegate's Objects.requireNonNull raise an NPE when the identifier carries
// neither a URI nor an OID. That condition is driven by signature parameters, i.e. by input
// rather than by a programming mistake, so it is returned as an error here; the delegate below
// keeps the panic for the case where a caller passes "" directly, which is a programming
// mistake.
func DSSJsonUtilsOidObject(objectIdentifier enumerations.ObjectIdentifier) (*JsonObject, error) {
	uri := spi.DSSUtilsURIOrUrnOID(objectIdentifier)
	if uri == "" {
		return nil, fmt.Errorf("uri must be defined!")
	}
	return DSSJsonUtilsOidObjectFromURI(uri, objectIdentifier.Description(),
		objectIdentifier.DocumentationReferences()), nil
}

// DSSJsonUtilsOidObjectFromURI creates an 'oid' JsonObject per TS 119 182-1 clause 5.4.1. uri
// identifies the object and is REQUIRED; desc is the object description and is OPTIONAL; docRefs
// are URIs carrying additional information about the object and are OPTIONAL. Port of
// getOidObject(String, String, String[]).
//
// The member order is the order of the Java statements, and it is preserved by the
// LinkedHashMap upstream and by jose.Object here, because these bytes end up inside a signed
// header.
func DSSJsonUtilsOidObjectFromURI(uri, desc string, docRefs []string) *JsonObject {
	if uri == "" {
		panic("uri must be defined!")
	}

	oidParams := jose.NewObject()
	oidParams.Put(JAdESHeaderParameterNamesId, uri)
	if utils.IsStringNotEmpty(desc) {
		oidParams.Put(JAdESHeaderParameterNamesDesc, desc)
	}
	if utils.IsArrayNotEmpty(docRefs) {
		oidParams.Put(JAdESHeaderParameterNamesDocRefs, append([]string(nil), docRefs...))
	}

	return NewJsonObjectFromMap(oidParams)
}

// DSSJsonUtilsTstContainer creates a 'tstContainer' JsonObject per TS 119 182-1 clause 5.4.3.3
// "The tstContainer type". canonicalizationMethodURI is OPTIONAL - it shall not be present for
// content timestamps - and "" means absent. Port of getTstContainer(List, String).
func DSSJsonUtilsTstContainer(timestampBinaries []*model.TimestampBinary, canonicalizationMethodURI string) (*JsonObject, error) {
	if utils.IsCollectionEmpty(timestampBinaries) {
		return nil, fmt.Errorf("Impossible to create 'tstContainer'. List of TimestampBinaries cannot be null or empty!")
	}

	tstContainerParams := jose.NewObject()
	if canonicalizationMethodURI != "" {
		tstContainerParams.Put(JAdESHeaderParameterNamesCanonAlg, canonicalizationMethodURI)
	}
	tsTokens := make([]any, 0, len(timestampBinaries))
	for _, timestampBinary := range timestampBinaries {
		tsTokens = append(tsTokens, dssJsonUtilsTstToken(timestampBinary))
	}
	tstContainerParams.Put(JAdESHeaderParameterNamesTstTokens, tsTokens)

	return NewJsonObjectFromMap(tstContainerParams), nil
}

// dssJsonUtilsTstToken creates a 'tstToken' JsonObject per TS 119 182-1 clause 5.4.3.3. Port of
// the private getTstToken(TimestampBinary).
//
// Upstream backs this one with a HashMap rather than a LinkedHashMap. It holds a single member,
// so the choice is invisible in the output, but it is reproduced anyway rather than silently
// "corrected": if a later DSS release adds the 'type'/'encoding'/'specRef' members the comment
// below mentions, the order they come out in will be HashMap order, and this port will follow.
func dssJsonUtilsTstToken(timestampBinary *model.TimestampBinary) *JsonObject {
	if timestampBinary == nil {
		panic("timestampBinary cannot be null!")
	}

	tstTokenParams := jose.NewHashObject()
	// only RFC 3161 TimestampTokens are supported
	// 'type', 'encoding' and 'specRef' params are not need to be defined (see TS 119-182 ch. 5.4.3.3)
	tstTokenParams.Put(JAdESHeaderParameterNamesVal, utils.ToBase64(timestampBinary.Bytes()))

	return NewJsonObjectFromMap(tstTokenParams)
}

// DSSJsonUtilsConcatenateDSSDocuments concatenates document octets into a single byte array.
// isBase64URLEncoded defines whether the document octets shall be base64url-encoded. Port of
// concatenateDSSDocuments(List, boolean).
func DSSJsonUtilsConcatenateDSSDocuments(documents []model.DSSDocument, isBase64URLEncoded bool) ([]byte, error) {
	if utils.IsCollectionEmpty(documents) {
		return nil, fmt.Errorf("Unable to build a JWS Payload. Reason : the detached content is not provided!")
	}
	if len(documents) == 1 {
		return DSSJsonUtilsDocumentOctets(documents[0], isBase64URLEncoded)
	}

	var buf bytes.Buffer
	for _, document := range documents {
		octets, err := DSSJsonUtilsDocumentOctets(document, isBase64URLEncoded)
		if err != nil {
			return nil, model.NewDSSErrorMessageCause(
				fmt.Sprintf("Unable to build a JWS Payload. Reason : %s", err.Error()), err)
		}
		buf.Write(octets)
	}
	return buf.Bytes(), nil
}

// DSSJsonUtilsWriteDocumentsDigest writes the digest over the concatenated binaries of the
// provided documents. Port of writeDocumentsDigest(List, boolean, DSSMessageDigestCalculator).
func DSSJsonUtilsWriteDocumentsDigest(documents []model.DSSDocument, isBase64URLEncoded bool,
	digestCalculator *spi.DSSMessageDigestCalculator) error {
	if utils.IsCollectionEmpty(documents) {
		return fmt.Errorf("Unable to build a message-digest. Reason : the detached content is not provided!")
	}

	for _, document := range documents {
		octets, err := DSSJsonUtilsDocumentOctets(document, isBase64URLEncoded)
		if err != nil {
			return err
		}
		digestCalculator.Update(octets)
	}
	return nil
}

// DSSJsonUtilsDocumentOctets returns the binaries of the document to be used for payload
// computation. When isBase64URLEncoded is true it returns the base64url-encoded binaries, and
// otherwise the original octets. Port of getDocumentOctets(DSSDocument, boolean).
func DSSJsonUtilsDocumentOctets(document model.DSSDocument, isBase64URLEncoded bool) ([]byte, error) {
	octets, err := spi.DSSUtilsToByteArrayOfDocument(document)
	if err != nil {
		return nil, err
	}
	if isBase64URLEncoded {
		octets = []byte(DSSJsonUtilsToBase64Url(octets))
	}
	return octets, nil
}

// DSSJsonUtilsIsJsonDocument checks whether the provided document is a JSON document. Port of
// isJsonDocument(DSSDocument).
//
// The leading-'{' test before parsing is not an optimisation: it is what makes a compact JWS -
// which is also valid input to this library - answer false here rather than being handed to a
// JSON parser that would reject it anyway.
func DSSJsonUtilsIsJsonDocument(document model.DSSDocument) (bool, error) {
	if !DSSJsonUtilsIsAllowedSignatureDocumentType(document) {
		return false, nil
	}
	stream, err := document.OpenStream()
	if err != nil {
		return false, model.NewDSSErrorMessageCause(
			fmt.Sprintf("Cannot read the document. Reason : %s", err.Error()), err)
	}
	defer utils.CloseQuietly(stream)

	first := make([]byte, 1)
	n, err := stream.Read(first)
	if err != nil || n == 0 {
		// An empty document has no first character; upstream's InputStream.read() returns -1
		// and the branch is skipped.
		return false, nil
	}
	if first[0] != '{' {
		return false, nil
	}
	rest, err := utils.ToByteArray(stream)
	if err != nil {
		return false, model.NewDSSErrorMessageCause(
			fmt.Sprintf("Cannot read the document. Reason : %s", err.Error()), err)
	}
	content := append(first, rest...)
	if len(content) < 2 {
		return false, nil
	}
	if _, err := jose.ParseJSON(string(content)); err != nil {
		// "Unable to parse content as JSON" - upstream logs the JoseException and answers false.
		return false, nil
	}
	return true, nil
}

// DSSJsonUtilsIsAllowedSignatureDocumentType checks whether the signature document has a type
// whose bytes can be extracted. Port of isAllowedSignatureDocumentType(DSSDocument).
func DSSJsonUtilsIsAllowedSignatureDocumentType(document model.DSSDocument) bool {
	switch document.(type) {
	case *model.DigestDocument, *HTTPHeader:
		// "The provided document of class '{}' cannot be parsed."
		return false
	}
	return true
}

// DSSJsonUtilsEtsiU returns the 'etsiU' container with the unsigned properties, or an empty
// slice. Port of getEtsiU(JWS).
func DSSJsonUtilsEtsiU(jws *JWS) []any {
	unprotected := jws.Unprotected()
	if unprotected == nil {
		return nil
	}
	etsiU := unprotected.Value(JAdESHeaderParameterNamesEtsiU)
	list, ok := etsiU.([]any)
	if !ok {
		// "Unable to extract 'etsiU' header : the obtained entry is not an array!"
		return nil
	}
	return list
}

// DSSJsonUtilsUnsignedPropertiesWithHeaderName returns the unsigned 'etsiU' properties of
// etsiUHeader whose name is headerName. Port of
// getUnsignedPropertiesWithHeaderName(EtsiUHeader, String).
func DSSJsonUtilsUnsignedPropertiesWithHeaderName(etsiUHeader *EtsiUHeader, headerName string) []*EtsiUComponent {
	if !etsiUHeader.IsExist() {
		return nil
	}

	var componentsWithHeaderName []*EtsiUComponent
	for _, attribute := range etsiUHeader.Attributes() {
		if headerName == attribute.HeaderName() {
			componentsWithHeaderName = append(componentsWithHeaderName, attribute)
		}
	}
	return componentsWithHeaderName
}

// DSSJsonUtilsDate parses an IETF RFC 3339 dateTime string. Port of getDate(String).
//
// Deprecated: since DSS 6.5, use spi.DSSUtilsParseRFCDate instead. Kept because the manifest
// ports the class as it stands, deprecation and all.
func DSSJsonUtilsDate(dateTimeString string) time.Time {
	return spi.DSSUtilsParseRFCDate(dateTimeString)
}

// DSSJsonUtilsIssuerSerial parses the 'kid' header value as in IETF RFC 5035. Port of
// getIssuerSerial(String), which processes failures silently and answers null.
func DSSJsonUtilsIssuerSerial(value string) *spi.IssuerSerial {
	if utils.IsStringNotEmpty(value) && utils.IsBase64Encoded(value) {
		binary := utils.FromBase64(value)
		if spi.DSSASN1UtilsIsAsn1Encoded(binary) {
			return spi.DSSASN1UtilsIssuerSerial(binary)
		}
	}
	// process silently
	return nil
}

// DSSJsonUtilsExtractJAdESCounterSignature extracts a counter signature from a 'cSig' value with
// respect to the found format. Port of
// extractJAdESCounterSignature(EtsiUComponent, Signature).
func DSSJsonUtilsExtractJAdESCounterSignature(cSigAttribute *EtsiUComponent, masterSignature *Signature) (*Signature, error) {
	cSigObject := cSigAttribute.Value()

	var cSigValue string
	switch v := cSigObject.(type) {
	case string:
		cSigValue = v
	case *jose.Object:
		cSigValue = jose.JSON(v)
	case *JsonObject:
		cSigValue = v.ToJSONString()
	default:
		// "Unsupported entry of type 'cSig' found! The entry is skipped"
	}

	if utils.IsStringEmpty(cSigValue) {
		return nil, nil
	}

	cSigDocument := model.NewInMemoryDocument([]byte(cSigValue))

	factory := NewJWSDocumentAnalyzerFactory()
	if !factory.IsSupported(cSigDocument) {
		return nil, nil
	}
	documentAnalyzer := factory.Create(cSigDocument)
	signatures := documentAnalyzer.Signatures()

	/*
	 * 5.3.2 The cSig (counter signature) JSON object
	 *
	 * The cSig JSON object shall contain one counter signature of the JAdES signature where
	 * cSig is incorporated.
	 */
	if len(signatures) != 1 {
		// "{} counter signatures found in 'cSig' element. Only one is allowed!"
		return nil, nil
	}
	signature, ok := signatures[0].(*Signature) // only one is considered
	if !ok {
		return nil, nil
	}
	signature.SetMasterSignature(masterSignature)
	signature.SetMasterCSigComponent(cSigAttribute)
	signature.SetDetachedContents([]model.DSSDocument{
		model.NewInMemoryDocument(masterSignature.SignatureValue()),
	})
	return signature, nil
}

// DSSJsonUtilsValidateAgainstJAdESSchema validates a JWS against the JAdES schema
// (ETSI TS 119 182-1) and returns the validation errors, empty if none occurred. Port of
// validateAgainstJAdESSchema(JWS).
func DSSJsonUtilsValidateAgainstJAdESSchema(jws *JWS) []string {
	var errors []string

	headerJSON := jws.Headers().FullHeaderAsJSONString()
	errors = append(errors, specs.JAdESProtectedHeaderUtilsInstance().ValidateAgainstSchema(headerJSON)...)

	unprotected := jws.Unprotected()
	if unprotected.Size() != 0 {
		unprotectedJSON := jose.JSON(unprotected)
		errors = append(errors, specs.JAdESUnprotectedHeaderUtilsInstance().ValidateAgainstSchema(unprotectedJSON)...)

		if etsiUComponents, ok := unprotected.Value(JAdESHeaderParameterNamesEtsiU).([]any); ok {
			if DSSJsonUtilsAreAllBase64UrlComponents(etsiUComponents) {
				clearEtsiURepresentation := dssJsonUtilsClearEtsiURepresentation(unprotected)
				clearEtsiUJSON := jose.JSON(clearEtsiURepresentation)
				errors = append(errors, specs.JAdESUnprotectedHeaderUtilsInstance().ValidateAgainstSchema(clearEtsiUJSON)...)
			}
		}
	}

	return errors
}

// DSSJsonUtilsCheckComponentsUnicity checks whether all components have one type (strings or
// objects). Port of checkComponentsUnicity(List).
func DSSJsonUtilsCheckComponentsUnicity(components []any) bool {
	if utils.IsCollectionEmpty(components) {
		return true
	}
	stringFormat := DSSJsonUtilsIsStringFormat(components[0])
	for _, component := range components[1:] {
		if stringFormat != DSSJsonUtilsIsStringFormat(component) {
			return false
		}
	}
	return true
}

// DSSJsonUtilsIsStringFormat checks whether the object is a String. Port of isStringFormat(Object).
func DSSJsonUtilsIsStringFormat(object any) bool {
	_, ok := object.(string)
	return ok
}

// DSSJsonUtilsAreAllBase64UrlComponents checks whether all the components are base64url-encoded.
// Port of areAllBase64UrlComponents(List).
func DSSJsonUtilsAreAllBase64UrlComponents(components []any) bool {
	for _, component := range components {
		str, ok := component.(string)
		if !ok || !DSSJsonUtilsIsBase64UrlEncoded(str) {
			return false
		}
	}
	return true
}

// dssJsonUtilsClearEtsiURepresentation builds the clear-JSON view of a base64url-encoded 'etsiU'
// array, for schema validation. Port of the private getClearEtsiURepresentation(Map).
//
// Upstream backs the result with a HashMap; it holds one member, so the choice is invisible, and
// it is reproduced for the same reason as in dssJsonUtilsTstToken.
func dssJsonUtilsClearEtsiURepresentation(unprotected *jose.Object) *jose.Object {
	stringComponents, _ := unprotected.Value(JAdESHeaderParameterNamesEtsiU).([]any)
	clearComponents := make([]any, 0, len(stringComponents))
	for _, component := range stringComponents {
		parsed, _ := DSSJsonUtilsParseEtsiUComponent(component)
		clearComponents = append(clearComponents, parsed)
	}
	clearEtsiU := jose.NewHashObject()
	clearEtsiU.Put(JAdESHeaderParameterNamesEtsiU, clearComponents)
	return clearEtsiU
}

// DSSJsonUtilsParseEtsiUComponent parses an 'etsiU' component as it stands - base64url-encoded or
// clear JSON - and returns the resulting map. Port of parseEtsiUComponent(Object); the bool
// reports whether a map was produced, standing in for upstream's null return.
func DSSJsonUtilsParseEtsiUComponent(etsiUComponent any) (*jose.Object, bool) {
	switch v := etsiUComponent.(type) {
	case *jose.Object:
		if v.Size() != 1 {
			// "The 'etsiU' shall contain only one entry!"
			return nil, false
		}
		return v, true

	case *JsonObject:
		if v.Size() != 1 {
			return nil, false
		}
		return v.JSONObject(), true

	case string:
		if !DSSJsonUtilsIsBase64UrlEncoded(v) {
			// "A String component of 'etsiU' array shall be base64Url encoded!"
			return nil, false
		}
		itemBinaries := DSSJsonUtilsFromBase64Url(v)
		parsed, err := jose.ParseJSON(string(itemBinaries))
		if err != nil {
			// "An error occurred during 'etsiU' component parsing"
			return nil, false
		}
		return parsed, true

	default:
		// "A component of unsupported class found inside the 'etsiU' array!"
		return nil, false
	}
}

// DSSJsonUtilsParseSPDocSpecification builds an SpDocSpecification from the provided JSON object
// element. Port of parseSPDocSpecification(Object), which swallows exceptions and answers null.
func DSSJsonUtilsParseSPDocSpecification(spDocSpecificationObject any) *model.SpDocSpecification {
	spDSpec := DSSJsonUtilsToMap(spDocSpecificationObject, JAdESHeaderParameterNamesSpDspec)
	if spDSpec.Size() == 0 {
		// "The spDSpec element is empty!"
		return nil
	}

	spDocSpecification := model.NewSpDocSpecification()

	id := DSSJsonUtilsGetAsString(spDSpec, JAdESHeaderParameterNamesId)
	if utils.IsStringNotEmpty(id) {
		spDocSpecification.SetId(spi.DSSUtilsObjectIdentifierValue(id))
	}

	desc := DSSJsonUtilsGetAsString(spDSpec, JAdESHeaderParameterNamesDesc)
	if utils.IsStringNotEmpty(desc) {
		spDocSpecification.SetDescription(desc)
	}

	docRefsList := DSSJsonUtilsGetAsList(spDSpec, JAdESHeaderParameterNamesDocRefs)
	if utils.IsCollectionNotEmpty(docRefsList) {
		docRefs := make([]string, 0, len(docRefsList))
		for _, docRef := range docRefsList {
			str, _ := docRef.(string)
			docRefs = append(docRefs, str)
		}
		spDocSpecification.SetDocumentationReferences(docRefs...)
	}

	return spDocSpecification
}

// DSSJsonUtilsToJWSJsonSerializationObject converts a JWS to a JWSJsonSerializationObject. Port
// of toJWSJsonSerializationObject(JWS).
func DSSJsonUtilsToJWSJsonSerializationObject(jws *JWS) *JWSJsonSerializationObject {
	jwsJsonSerializationObject := NewJWSJsonSerializationObject()
	jwsJsonSerializationObject.AddSignature(jws)
	jwsJsonSerializationObject.SetPayload(jws.SignedPayload())
	return jwsJsonSerializationObject
}

// DSSJsonUtilsSigningInputBytes computes the signing input bytes for a JWS signature. Port of
// getSigningInputBytes(JWS), i.e. RFC 7797 section 3:
//
//	| "b64" | JWS Signing Input Formula                                  |
//	| true  | ASCII(BASE64URL(UTF8(JWS Protected Header)) || '.' ||      |
//	|       | BASE64URL(JWS Payload))                                    |
//	| false | ASCII(BASE64URL(UTF8(JWS Protected Header)) || '.') ||     |
//	|       | JWS Payload                                                |
//
// In the b64=false branch the payload is appended as raw bytes: as upstream's own comment puts
// it, "unencoded payload shall not be converted to a string, it can lead to a data corruption!"
func DSSJsonUtilsSigningInputBytes(jws *JWS) []byte {
	if !jws.IsRfc7797UnencodedPayload() {
		dataToBeSignedString := DSSJsonUtilsConcatenate(jws.EncodedHeader(), jws.EncodedPayload())
		return DSSJsonUtilsAsciiBytes(dataToBeSignedString)
	}

	var buf bytes.Buffer
	buf.Write(DSSJsonUtilsAsciiBytes(jws.EncodedHeader()))
	buf.WriteByte(0x2e) // ascii for "."
	payloadBytes := jws.UnverifiedPayloadBytes()
	if utils.IsArrayNotEmpty(payloadBytes) {
		buf.Write(payloadBytes)
	}
	return buf.Bytes()
}

// DSSJsonUtilsExtractJOSEHeaderMembersSet extracts the key set used within a JOSE header
// (protected + unprotected). Port of extractJOSEHeaderMembersSet(JWS).
func DSSJsonUtilsExtractJOSEHeaderMembersSet(jws *JWS) (map[string]struct{}, error) {
	signedHeaders, err := jose.ParseJSON(jws.Headers().FullHeaderAsJSONString())
	if err != nil {
		return nil, exception.NewIllegalInputExceptionWithCause(fmt.Sprintf(
			"Unable to extract key set from a JOSE header! Reason : %s", err.Error()), err)
	}
	joseHeaderMemberKeys := make(map[string]struct{})
	for _, key := range signedHeaders.Keys() {
		joseHeaderMemberKeys[key] = struct{}{}
	}
	if jws.Unprotected() != nil {
		for _, key := range jws.Unprotected().Keys() {
			joseHeaderMemberKeys[key] = struct{}{}
		}
	}
	return joseHeaderMemberKeys, nil
}

// DSSJsonUtilsGetAsBoolean gets the value under key as a Boolean. Port of getAsBoolean(Map, String).
func DSSJsonUtilsGetAsBoolean(m *jose.Object, key string) *bool {
	return DSSJsonUtilsToBoolean(m.Value(key), key)
}

// DSSJsonUtilsToBoolean safely converts an object to a Boolean if possible, and answers nil
// otherwise. headerName names the header the value came from, which upstream uses only to build
// its warning. Port of toBoolean(Object, String).
func DSSJsonUtilsToBoolean(object any, headerName string) *bool {
	if b, ok := object.(bool); ok {
		return &b
	}
	// "Unable to process '{}' header parameter. The Boolean type is expected!"
	_ = headerName
	return nil
}

// DSSJsonUtilsGetAsString gets the value under key as a String, or "". Port of
// getAsString(Map, String).
func DSSJsonUtilsGetAsString(m *jose.Object, key string) string {
	return DSSJsonUtilsToStringWithHeaderName(m.Value(key), key)
}

// DSSJsonUtilsToString safely converts an object to a String if possible, and answers "" (Java's
// Utils.EMPTY_STRING) otherwise. Port of toString(Object).
func DSSJsonUtilsToString(object any) string {
	return DSSJsonUtilsToStringWithHeaderName(object, "")
}

// DSSJsonUtilsToStringWithHeaderName safely converts an object to a String if possible, and
// answers "" otherwise. Port of toString(Object, String).
func DSSJsonUtilsToStringWithHeaderName(object any, headerName string) string {
	if s, ok := object.(string); ok {
		return s
	}
	// "Unable to process '{}' header parameter. The String type is expected!"
	_ = headerName
	return ""
}

// DSSJsonUtilsGetAsNumber gets the value under key as a Number, or nil. Port of
// getAsNumber(Map, String).
func DSSJsonUtilsGetAsNumber(m *jose.Object, key string) *jose.Number {
	return DSSJsonUtilsToNumberWithHeaderName(m.Value(key), key)
}

// DSSJsonUtilsToNumber safely converts an object to a Number if possible, and answers nil
// otherwise. Port of toNumber(Object).
func DSSJsonUtilsToNumber(object any) *jose.Number {
	return DSSJsonUtilsToNumberWithHeaderName(object, "")
}

// DSSJsonUtilsToNumberWithHeaderName safely converts an object to a Number if possible, and
// answers nil otherwise. Port of toNumber(Object, String).
func DSSJsonUtilsToNumberWithHeaderName(object any, headerName string) *jose.Number {
	if n, ok := object.(*jose.Number); ok {
		return n
	}
	// "Unable to process '{}' header parameter. The Number type is expected!"
	_ = headerName
	return nil
}

// DSSJsonUtilsGetAsMap gets the value under key as a JSON object, or an empty one. Port of
// getAsMap(Map, String).
func DSSJsonUtilsGetAsMap(m *jose.Object, key string) *jose.Object {
	return DSSJsonUtilsToMap(m.Value(key), key)
}

// DSSJsonUtilsToMapValue safely converts an object to a JSON object if possible, and answers an
// empty one otherwise. Port of toMap(Object).
func DSSJsonUtilsToMapValue(object any) *jose.Object {
	return DSSJsonUtilsToMap(object, "")
}

// DSSJsonUtilsToMap safely converts an object to a JSON object if possible, and answers an empty
// one (Java's Collections.emptyMap()) otherwise. Port of toMap(Object, String).
func DSSJsonUtilsToMap(object any, headerName string) *jose.Object {
	switch v := object.(type) {
	case *jose.Object:
		return v
	case *JsonObject:
		return v.JSONObject()
	}
	// "Unable to process '{}' header parameter. The JSON Object type is expected!"
	_ = headerName
	return jose.NewObject()
}

// DSSJsonUtilsGetAsList gets the value under key as a JSON array, or an empty one. Port of
// getAsList(Map, String).
func DSSJsonUtilsGetAsList(m *jose.Object, key string) []any {
	return DSSJsonUtilsToList(m.Value(key), key)
}

// DSSJsonUtilsToListValue safely converts an object to a JSON array if possible, and answers an
// empty one otherwise. Port of toList(Object).
func DSSJsonUtilsToListValue(object any) []any {
	return DSSJsonUtilsToList(object, "")
}

// DSSJsonUtilsToList safely converts an object to a JSON array if possible, and answers an empty
// one otherwise. Port of toList(Object, String).
func DSSJsonUtilsToList(object any, headerName string) []any {
	if list, ok := object.([]any); ok {
		return list
	}
	// "Unable to process '{}' header parameter. The JSON Array type is expected!"
	_ = headerName
	return nil
}

// DSSJsonUtilsToListOfStrings converts a list of objects to a list of Strings, skipping the
// entries that are not (or are empty) strings. Port of toListOfStrings(List).
func DSSJsonUtilsToListOfStrings(list []any) []string {
	listOfStrings := make([]string, 0, len(list))
	for _, item := range list {
		str := DSSJsonUtilsToString(item)
		if utils.IsStringNotEmpty(str) {
			listOfStrings = append(listOfStrings, str)
		}
		// else "An empty String entry within a JSON Object has been skipped."
	}
	return listOfStrings
}

// DSSJsonUtilsToListOfNumbers converts a list of objects to a list of Numbers, skipping the
// entries that are not numbers. Port of toListOfNumbers(List).
func DSSJsonUtilsToListOfNumbers(list []any) []*jose.Number {
	listOfNumbers := make([]*jose.Number, 0, len(list))
	for _, item := range list {
		if num := DSSJsonUtilsToNumber(item); num != nil {
			listOfNumbers = append(listOfNumbers, num)
		}
	}
	return listOfNumbers
}

// DSSJsonUtilsGetAsNumericDate gets the value under key as a Date. Port of
// getAsNumericDate(Map, String).
func DSSJsonUtilsGetAsNumericDate(m *jose.Object, key string) time.Time {
	return DSSJsonUtilsToNumericDateWithHeaderName(m.Value(key), key)
}

// DSSJsonUtilsToNumericDate safely converts an object to a Date if possible, and answers the zero
// Time otherwise. Port of toNumericDate(Object).
func DSSJsonUtilsToNumericDate(object any) time.Time {
	return DSSJsonUtilsToNumericDateWithHeaderName(object, "")
}

// DSSJsonUtilsToNumericDateWithHeaderName safely converts an object to a Date if possible, and
// answers the zero Time otherwise. Port of toNumericDate(Object, String).
//
// A NumericDate is seconds since the epoch, which upstream widens to milliseconds with
// getTimeValueInMilliseconds(number.longValue()) - so a fractional NumericDate loses its
// fraction, here as there.
func DSSJsonUtilsToNumericDateWithHeaderName(object any, headerName string) time.Time {
	if number, ok := object.(*jose.Number); ok && number != nil {
		timeValueInMilliseconds := spi.DSSUtilsTimeValueInMilliseconds(number.Int64())
		millis := float64(timeValueInMilliseconds)
		return spi.DSSUtilsDateFromMilliseconds(&millis)
	}
	// "Unable to process '{}' header parameter. The JSON Number type is expected!"
	_ = headerName
	return time.Time{}
}

// DSSJsonUtilsMimeTypeString returns a complete mime type string, adding the "application/"
// prefix when required. Port of getMimeTypeString(String).
func DSSJsonUtilsMimeTypeString(mimeType string) string {
	if utils.IsStringNotEmpty(mimeType) && !strings.Contains(mimeType, "/") {
		return DSSJsonUtilsMimeTypeApplicationPrefix + mimeType
	}
	return mimeType
}

// DSSJsonUtilsParseBase64UrlEncoded parses the provided base64url-encoded string and returns the
// corresponding object. The string shall conform to a JSON value specification but need not be a
// root element (a map). Port of parseBase64UrlEncoded(String).
func DSSJsonUtilsParseBase64UrlEncoded(base64URLEncodedString string) (any, error) {
	if !DSSJsonUtilsIsBase64UrlEncoded(base64URLEncodedString) {
		return nil, fmt.Errorf("Base64Url encoded string is expected.")
	}
	decodedString := string(DSSJsonUtilsFromBase64Url(base64URLEncodedString))
	value, err := DSSJsonUtilsParseJSONString(decodedString)
	if err != nil {
		return nil, model.NewDSSErrorMessageCause(
			fmt.Sprintf("An error occurred on decoding the string. Reason : %s", err.Error()), err)
	}
	return value, nil
}

// DSSJsonUtilsParseJSONString parses a plain string containing a JSON value. Port of
// parseJsonString(String).
//
// It goes through `new JSONParser().parse(...)` rather than JsonUtil.parseJson, so any JSON value
// is accepted at the root and objects come back in java.util.HashMap order; see
// jose.ParseJSONAny.
func DSSJsonUtilsParseJSONString(jsonString string) (any, error) {
	value, err := jose.ParseJSONAny(jsonString)
	if err != nil {
		return nil, model.NewDSSErrorMessageCause(
			fmt.Sprintf("An error occurred on parsing the string. Reason : %s", err.Error()), err)
	}
	return value, nil
}

// DSSJsonUtilsParseJSONStringToMap parses a JSON string to a JSON object. Port of
// parseJsonStringToMap(String).
func DSSJsonUtilsParseJSONStringToMap(jsonString string) (*jose.Object, error) {
	object, err := jose.ParseJSON(jsonString)
	if err != nil {
		return nil, model.NewDSSErrorMessageCause(
			fmt.Sprintf("Unable to parse string : %s", err.Error()), err)
	}
	return object, nil
}
