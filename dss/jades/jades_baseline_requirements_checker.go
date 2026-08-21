// Ported from
// dss-jades/src/main/java/eu/europa/esig/dss/jades/validation/JAdESBaselineRequirementsChecker.java
// (DSS 6.5.RC1).
//
// Performs checks according to TS 119 182-1 v1.1.1 "6.3 Requirements on JAdES components and
// services".
//
// slf4j logging is dropped per PORTING.md; every LOG.warn/LOG.debug call site is called out in
// the surrounding comment instead.
//
// critRequirements calls DSSJsonUtilsExtractJOSEHeaderMembersSet, which - unlike its Java
// counterpart's undeclared (unchecked) failure - returns an explicit error in this port. Java
// lets the underlying JoseException propagate out of hasAdESProfile()/hasBaselineBProfile()
// uncaught (neither has a throws clause nor a surrounding try/catch), so a parse failure here
// panics with the same message, reproducing that unchecked propagation - none of the
// BaselineRequirementsCheckerContract methods this file implements has an error return to use
// instead.
package jades

import (
	"time"

	"github.com/ryftcore/dss-go/dss/internal/jose"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/spi/validation"
	"github.com/ryftcore/dss-go/dss/utils"
)

// jadesBaselineRequirementsCheckerSigTObsolescenceDate is the 2025-07-15T00:00:00Z date, see
// TS 119 182-1. Port of the private static final Date SIG_T_OBSOLESCENCE_DATE field.
var jadesBaselineRequirementsCheckerSigTObsolescenceDate = time.Date(2025, time.July, 15, 0, 0, 0, 0, time.UTC)

// JAdESBaselineRequirementsChecker checks conformance of a JAdES signature to the requested
// baseline format. Port of the class JAdESBaselineRequirementsChecker, extending
// validation.BaselineRequirementsChecker[*JAdESSignature].
type JAdESBaselineRequirementsChecker struct {
	validation.BaselineRequirementsChecker[*JAdESSignature]
}

// NewJAdESBaselineRequirementsChecker is the default constructor. Port of the public
// JAdESBaselineRequirementsChecker(JAdESSignature, CertificateVerifier) constructor.
func NewJAdESBaselineRequirementsChecker(signature *JAdESSignature, offlineCertificateVerifier validation.CertificateVerifier) *JAdESBaselineRequirementsChecker {
	checker := &JAdESBaselineRequirementsChecker{
		BaselineRequirementsChecker: validation.NewBaselineRequirementsCheckerBaseWithVerifier[*JAdESSignature](signature, offlineCertificateVerifier),
	}
	checker.InitBaselineRequirementsChecker(checker)
	return checker
}

// HasAdESProfile checks if the signature is conformant to the corresponding AdES profile. Port
// of the hasAdESProfile() override.
func (b *JAdESBaselineRequirementsChecker) HasAdESProfile() bool {
	signature := b.Signature()
	jws := signature.Jws()

	// 5.1.3 The cty (content type) header parameter
	if signature.IsCounterSignature() && utils.IsStringNotEmpty(jws.ProtectedHeaderValueAsString(jose.HeaderContentType)) {
		// Upstream logs "cty header shall not be present for a JAdES counter signature!".
		return false
	}

	// 5.1.7 The x5t#S256 (X.509 Certificate SHA-256 Thumbprint) header parameter
	certHeaders := 0
	if utils.IsStringNotEmpty(jws.ProtectedHeaderValueAsString(jose.HeaderX509CertificateSHA256Thumbprint)) {
		certHeaders++
	}
	if utils.IsCollectionNotEmpty(jws.ProtectedHeaderValueAsList(jose.HeaderX509CertificateChain)) {
		certHeaders++
	}
	if utils.IsCollectionNotEmpty(jws.ProtectedHeaderValueAsList(JAdESHeaderParameterNamesSigX5tS)) {
		certHeaders++
	}
	if jws.ProtectedHeaderValueAsMap(JAdESHeaderParameterNamesX5tO).Size() != 0 {
		certHeaders++
	}
	if certHeaders == 0 {
		// Upstream logs "At least one of x5t#256, x5c, sigX5ts, x5t#o headers shall be present
		// for JAdES signature!".
		return false
	}

	// 5.1.9 The crit (critical) header parameter
	if !b.critRequirements(jws, "JAdES") {
		// validation errors returned inside
		return false
	}

	// 5.2.1 The sigT (claimed signing time) header parameter
	sigT := jws.ProtectedHeaderValueAsString(JAdESHeaderParameterNamesSigT)
	signingTime := spi.DSSUtilsParseRFCDate(sigT)
	if utils.IsStringNotEmpty(sigT) && !signingTime.IsZero() && !signingTime.Before(jadesBaselineRequirementsCheckerSigTObsolescenceDate) {
		// Upstream logs "sigT header shall not be present for JAdES signature produced starting
		// at 2025-07-15T00:00:00Z!".
		return false
	}
	return true
}

// HasBaselineBProfile checks if the signature has a corresponding BASELINE-B profile. Port of
// the hasBaselineBProfile() override.
func (b *JAdESBaselineRequirementsChecker) HasBaselineBProfile() bool {
	signature := b.Signature()
	jws := signature.Jws()
	etsiUHeader := signature.EtsiUHeader()

	// alg (Cardinality == 1)
	if utils.IsStringEmpty(jws.ProtectedHeaderValueAsString(jose.HeaderAlgorithm)) {
		// Upstream logs "alg header shall be present for JAdES-BASELINE-B signature
		// (cardinality == 1)!".
		return false
	}
	// cty (Conditional presence)
	if signature.IsCounterSignature() && utils.IsStringNotEmpty(jws.ProtectedHeaderValueAsString(jose.HeaderContentType)) {
		// Upstream logs "cty header shall not be present for a JAdES-BASELINE-B counter
		// signature!".
		return false
	}
	// verify 'crit' as of RFC 7515 and ETSI TS 119 182-1
	if !b.critRequirements(jws, "JAdES-BASELINE-B") {
		// validation errors returned inside
		return false
	}
	// sigT (Cardinality == 1)
	if !b.signingTimeRequirement(jws) {
		return false
	}
	// x5t#256 / x5t#o / sigX5ts (Cardinality == 1)
	certHeaders := 0
	if utils.IsStringNotEmpty(jws.ProtectedHeaderValueAsString(jose.HeaderX509CertificateSHA256Thumbprint)) {
		certHeaders++
	}
	if jws.ProtectedHeaderValueAsMap(JAdESHeaderParameterNamesX5tO).Size() != 0 {
		certHeaders++
	}
	if utils.IsCollectionNotEmpty(jws.ProtectedHeaderValueAsList(JAdESHeaderParameterNamesSigX5tS)) {
		certHeaders++
	}
	if certHeaders != 1 {
		// Upstream logs "One and only one of x5t#256, x5t#o, sigX5ts headers shall be present
		// for JAdES-BASELINE-B signature (cardinality == 1)!".
		return false
	}
	// sigPSt (Cardinality 0 or 1)
	if len(DSSJsonUtilsUnsignedPropertiesWithHeaderName(etsiUHeader, JAdESHeaderParameterNamesSigPSt)) > 1 {
		// Upstream logs "Only one sigPSt header shall be present for JAdES-BASELINE-B signature
		// (cardinality 0 or 1)!".
		return false
	}
	// Additional requirement (b)
	if !b.IsSignaturePolicyIdentifierHashPresent() && signature.SignaturePolicyStore() != nil {
		// Upstream logs "sigPSt header shall not be incorporated for JAdES-BASELINE-B signature
		// with not defined sigPId/hashAV (requirement (b))!".
		return false
	}
	return true
}

// critRequirements ports the private critRequirements(JWS, String).
func (b *JAdESBaselineRequirementsChecker) critRequirements(jws *JWS, profile string) bool {
	var critList []any

	// crit (conditional presence, required only for some elements)
	crit := jws.Headers().Value(jose.HeaderCritical)
	if crit != nil {
		// crit shall be an array (List)
		list, ok := crit.([]any)
		if !ok {
			// Upstream logs "crit header shall be an instance of json array type for a {}
			// signature!".
			return false
		}
		// crit cannot be empty
		critList = list
		if utils.IsCollectionEmpty(critList) {
			// Upstream logs "crit header shall not be empty for a {} signature (see RFC
			// 7515)!".
			return false
		}
		// A plain map[any]struct{} would panic on an unhashable entry (a malformed crit array
		// member that is itself a JSON array or object); Java's HashSet<Object> tolerates that
		// case via identity-based equals()/hashCode(), so duplicate detection here uses an O(n^2)
		// equality scan instead - crit arrays are always small in practice.
		var duplicates []any
		for i, entry := range critList {
			duplicate := false
			for _, earlier := range critList[:i] {
				if jadesBaselineRequirementsCheckerEqual(entry, earlier) {
					duplicate = true
					break
				}
			}
			if duplicate {
				duplicates = append(duplicates, entry)
			}
		}
		if utils.IsCollectionNotEmpty(duplicates) {
			// Upstream logs "crit header shall not contain duplicates for a {} signature (see
			// RFC 7515)! Found duplicates : '{}'".
			return false
		}
	}

	keySet, err := DSSJsonUtilsExtractJOSEHeaderMembersSet(jws)
	if err != nil {
		panic(err)
	}
	for key := range keySet {
		// critical headers shall not be present within crit
		if DSSJsonUtilsIsCriticalHeaderException(key) {
			if jadesBaselineRequirementsCheckerContains(critList, key) {
				// Upstream logs "crit header shall not contain headers listed in RFC 7515 or
				// RFC 7518 for a {} signature (see RFC 7515)! Found header : '{}'".
				return false
			}
		} else if DSSJsonUtilsIsRequiredCriticalHeader(key) {
			if crit == nil {
				// Upstream logs "crit header shall be present when '{}' header is present in a
				// signature for a {} signature!".
				return false
			} else if !jadesBaselineRequirementsCheckerContains(critList, key) {
				// Upstream logs "crit header shall contain '{}' header when present in a
				// signature for a {} signature!".
				return false
			}
		}
	}
	for _, critEntry := range critList {
		// crit shall contain String entries
		critString, ok := critEntry.(string)
		if !ok {
			// Upstream logs "An entry of crit header shall be an instance of String type for a
			// {} signature!".
			return false
		}
		// crit shall not contain not-used entries
		if _, ok := keySet[critString]; !ok {
			// Upstream logs "crit header can contain only entries used within a signed header
			// for a {} signature (see RFC 7515)! Found header : '{}'".
			return false
		}
		// Conforming implementations must reject input containing critical extensions that are
		// not understood or cannot be processed.
		_, supported := DSSJsonUtilsSupportedProtectedCriticalHeaders()[critString]
		if !supported && JAdESHeaderParameterNamesEtsiU != critString {
			// Upstream logs "crit header shall not contain a header that cannot be understood
			// and processed for a {} signature (see RFC 7515)! Found header : '{}'".
			return false
		}
	}
	return true
}

// jadesBaselineRequirementsCheckerEqual compares two crit-list entries the way Java's
// Object#equals() would, without risking Go's == panic on an unhashable/uncomparable dynamic
// type (a malformed crit entry that is itself a JSON array or object): entries of a comparable
// dynamic type compare by ==, anything else compares unequal (matching Java's default
// identity-based Object#equals() for two distinct such values, which is the only case this
// port's *jose.Object/[]any parser can ever produce for repeated malformed entries).
func jadesBaselineRequirementsCheckerEqual(a, b any) (equal bool) {
	defer func() {
		if recover() != nil {
			equal = false
		}
	}()
	return a == b
}

// jadesBaselineRequirementsCheckerContains ports List#contains(Object) for the crit list.
func jadesBaselineRequirementsCheckerContains(list []any, key string) bool {
	for _, entry := range list {
		if entry == key {
			return true
		}
	}
	return false
}

// signingTimeRequirement ports the private signingTimeRequirement(JWS).
func (b *JAdESBaselineRequirementsChecker) signingTimeRequirement(jws *JWS) bool {
	/*
	 * a) Requirements for iat and sigT. Before 2025-07-15T00:00:00Z the generator should include
	 *    the iat header parameter for indicating the claimed signing time in new JAdES signatures
	 *    and should not include the iat header parameter for indicating the claimed signing time
	 *    in new JAdES signatures. Starting at 2025-07-15T00:00:00Z the generator shall include the
	 *    iat header parameter for indicating the claimed signing time in new JAdES signatures.
	 */
	iat := jws.ProtectedHeaderValueAsNumber(JWTClaimNamesIat)
	sigT := jws.ProtectedHeaderValueAsString(JAdESHeaderParameterNamesSigT)
	signingTime := b.Signature().SigningTime()

	// iat or sigT (Cardinality == 1)
	if iat == nil && utils.IsStringEmpty(sigT) {
		// Upstream logs "Either iat header or sigT header (for signatures before
		// 2025-07-15T00:00:00Z) shall be present for JAdES-BASELINE-B signature (cardinality ==
		// 1)!".
		return false

	} else if signingTime == nil {
		// Upstream logs "Invalid date format extracted from {} header parameter for
		// JAdES-BASELINE-B signature (cardinality == 1)!".
		return false

	} else if iat == nil {
		if signingTime.Before(jadesBaselineRequirementsCheckerSigTObsolescenceDate) {
			// Upstream logs "iat header should be present for JAdES-BASELINE-B signature
			// produced before 2025-07-15T00:00:00Z (cardinality == 0 or 1)!" at DEBUG level.
		} else {
			// Upstream logs "iat header shall be present for JAdES-BASELINE-B signature
			// produced starting at 2025-07-15T00:00:00Z (cardinality == 1)!".
			return false
		}

	} else if utils.IsStringNotEmpty(sigT) {
		// Upstream logs "Both iat and sigT headers are not allowed for JAdES-BASELINE-B
		// signature (cardinality == 1)!".
		return false
	}

	return true
}

// HasBaselineTProfile checks if the signature has a corresponding BASELINE-T profile. Port of
// the hasBaselineTProfile() override.
func (b *JAdESBaselineRequirementsChecker) HasBaselineTProfile() bool {
	if !b.MinimalTRequirement() {
		return false
	}
	etsiUHeader := b.Signature().EtsiUHeader()
	// Additional requirement (c)
	for _, etsiUComponent := range DSSJsonUtilsUnsignedPropertiesWithHeaderName(etsiUHeader, JAdESHeaderParameterNamesSigTst) {
		sigTst := DSSJsonUtilsToMap(etsiUComponent.Value(), JAdESHeaderParameterNamesSigTst)
		tstTokens := DSSJsonUtilsGetAsList(sigTst, JAdESHeaderParameterNamesTstTokens)
		if len(tstTokens) != 1 {
			// Upstream logs "sigTst shall contain only one electronic timestamp for
			// JAdES-BASELINE-T signature (requirement (c))!".
			return false
		}
	}
	// Additional requirement (d)
	if !b.SignatureTimestampsCreatedBeforeSignCertExpiration() {
		// Upstream logs "sigTst shall be created before expiration of the signing-certificate
		// for JAdES-BASELINE-T signature (requirement (d))!".
		return false
	}
	return true
}

// HasBaselineLTProfile checks if the signature has a corresponding BASELINE-LT profile. Port of
// the hasBaselineLTProfile() override.
func (b *JAdESBaselineRequirementsChecker) HasBaselineLTProfile() bool {
	if !b.MinimalLTRequirement() {
		return false
	}
	etsiUHeader := b.Signature().EtsiUHeader()
	// xRefs (Cardinality == 0)
	if len(DSSJsonUtilsUnsignedPropertiesWithHeaderName(etsiUHeader, JAdESHeaderParameterNamesXRefs)) > 0 {
		// Upstream logs "xRefs header shall not be present for JAdES-BASELINE-LT signature
		// (cardinality == 0)!".
		return false
	}
	// axRefs (Cardinality == 0)
	if len(DSSJsonUtilsUnsignedPropertiesWithHeaderName(etsiUHeader, JAdESHeaderParameterNamesAxRefs)) > 0 {
		// Upstream logs "axRefs header shall not be present for JAdES-BASELINE-LT signature
		// (cardinality == 0)!".
		return false
	}
	// rRefs (Cardinality == 0)
	if len(DSSJsonUtilsUnsignedPropertiesWithHeaderName(etsiUHeader, JAdESHeaderParameterNamesRRefs)) > 0 {
		// Upstream logs "rRefs header shall not be present for JAdES-BASELINE-LT signature
		// (cardinality == 0)!".
		return false
	}
	// arRefs (Cardinality == 0)
	if len(DSSJsonUtilsUnsignedPropertiesWithHeaderName(etsiUHeader, JAdESHeaderParameterNamesArRefs)) > 0 {
		// Upstream logs "arRefs header shall not be present for JAdES-BASELINE-LT signature
		// (cardinality == 0)!".
		return false
	}
	// sigRTst (Cardinality == 0)
	if len(DSSJsonUtilsUnsignedPropertiesWithHeaderName(etsiUHeader, JAdESHeaderParameterNamesSigRTst)) > 0 {
		// Upstream logs "sigRTst header shall not be present for JAdES-BASELINE-LT signature
		// (cardinality == 0)!".
		return false
	}
	// rfsTst (Cardinality == 0)
	if len(DSSJsonUtilsUnsignedPropertiesWithHeaderName(etsiUHeader, JAdESHeaderParameterNamesRfsTst)) > 0 {
		// Upstream logs "rfsTst header shall not be present for JAdES-BASELINE-LT signature
		// (cardinality == 0)!".
		return false
	}
	return true
}

// ContainsLTLevelCertificates verifies whether the signature contains some of the LT-/XL level
// attributes. Port of the protected containsLTLevelCertificates() override, implementing
// validation.BaselineRequirementsCheckerOverrides.
func (b *JAdESBaselineRequirementsChecker) ContainsLTLevelCertificates() bool {
	return b.containsCertificateValues() || b.containsTstOrAnyValDataCertificates()
}

func (b *JAdESBaselineRequirementsChecker) containsCertificateValues() bool {
	etsiUHeader := b.Signature().EtsiUHeader()
	return len(DSSJsonUtilsUnsignedPropertiesWithHeaderName(etsiUHeader, JAdESHeaderParameterNamesXVals))+
		len(DSSJsonUtilsUnsignedPropertiesWithHeaderName(etsiUHeader, JAdESHeaderParameterNamesAxVals)) != 0
}

func (b *JAdESBaselineRequirementsChecker) containsTstOrAnyValDataCertificates() bool {
	etsiUHeader := b.Signature().EtsiUHeader()
	var validationDataHeaders []*EtsiUComponent
	validationDataHeaders = append(validationDataHeaders, DSSJsonUtilsUnsignedPropertiesWithHeaderName(etsiUHeader, JAdESHeaderParameterNamesTstVD)...)
	validationDataHeaders = append(validationDataHeaders, DSSJsonUtilsUnsignedPropertiesWithHeaderName(etsiUHeader, JAdESHeaderParameterNamesAnyValData)...)
	for _, etsiUComponent := range validationDataHeaders {
		valDataMap, ok := etsiUComponent.Value().(*jose.Object)
		if ok && valDataMap.Value(JAdESHeaderParameterNamesXVals) != nil {
			return true
		}
	}
	return false
}

// HasBaselineLTAProfile checks if the signature has a corresponding BASELINE-LTA profile. Port
// of the hasBaselineLTAProfile() override.
func (b *JAdESBaselineRequirementsChecker) HasBaselineLTAProfile() bool {
	return b.MinimalLTARequirement()
}

// compile-time assertion: a JAdESBaselineRequirementsChecker satisfies its own overrides contract.
var _ validation.BaselineRequirementsCheckerOverrides = (*JAdESBaselineRequirementsChecker)(nil)
