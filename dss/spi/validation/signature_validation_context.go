// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/validation/SignatureValidationContext.java (DSS 6.5.RC1).
//
// This type is the central consumer of most other types in this package: CertificateVerifier,
// RevocationDataLoadingStrategyFactory, RevocationDataVerifier, TimestampTokenVerifier,
// TrustAnchorVerifier, ValidationContext (the interface this type implements), ValidationData,
// TokenStatus, RevocationFreshnessStatus and EvidenceRecord. Every exported method below is
// named to match those types' actual Java-derived Go names (get/is dropped, Set kept as SetX).
//
// TimestampTokenVerifier splits Java's 4 isAcceptable(...) overloads into 4 distinctly-named
// methods - IsAcceptable, IsAcceptableAt, IsAcceptableWithChain, IsAcceptableWithChainAt; the
// 3-argument overload used below calls IsAcceptableWithChainAt.
//
// Sets are ported as append-if-absent slices (mirroring Java's LinkedHashSet insertion-order
// iteration) rather than Go maps keyed by pointer, because CertificateToken/Token equality in
// Java is value-based (see model.CertificateToken.Equals), which pointer identity would not
// reproduce; see certificateSetContains/tokenSetContains helpers below.
package validation

import (
	"sort"

	"time"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/model/x509/revocation"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/spi/x509/aia"
	"github.com/ryftcore/dss-go/dss/utils"
)

// SignatureValidationContext is a "cache" for one validation request that contains every
// object retrieved so far, during the validation of a signature.
type SignatureValidationContext struct {
	// processedSignatures is the set of signatures to process.
	processedSignatures []AdvancedSignature
	// processedCertificates is the set of certificates to process.
	processedCertificates []*model.CertificateToken
	// processedRevocations is the set of revocation data to process.
	processedRevocations []AnyRevocationToken
	// processedTimestamps is the set of timestamps to process.
	processedTimestamps []*TimestampToken
	// processedEvidenceRecords is the set of evidence records to process.
	processedEvidenceRecords []EvidenceRecord

	// certificateVerifier is the CertificateVerifier to use.
	certificateVerifier CertificateVerifier

	// tokensToProcess maps tokens to whether they have been processed yet: absent means not
	// yet registered, false means registered but not yet verified, true means verified.
	// Java's Map<Token, Boolean> with a null value becomes a two-slice (key/state) pair here
	// since model.Token values are not comparable map keys across all implementations.
	tokensToProcess []signatureValidationContextTokenState

	// signCertificateSignaturesUsage maps signing certificate tokens to the corresponding
	// signatures (b-level creation).
	signCertificateSignaturesUsage []signatureValidationContextCertSignatures

	// timestampCertChainDates holds the usages of a timestamp's certificate tokens.
	timestampCertChainDates []signatureValidationContextCertDates

	// poeTimes maps token IDs to their corresponding POE times.
	poeTimes map[string][]*signatureValidationContextPOE

	// tokenIssuerMap caches tokens and their CertificateToken issuers.
	tokenIssuerMap []signatureValidationContextTokenIssuer

	// certificateChildrenMap caches parent CertificateTokens and their corresponding issued
	// certificates.
	certificateChildrenMap []signatureValidationContextCertChildren

	// externalRevocationTokensMap caches issuer CertificateTokens and their issued revocation
	// data (i.e. obtained online).
	externalRevocationTokensMap []signatureValidationContextExternalRevocations

	// documentCertificateSource holds certificates from the document.
	documentCertificateSource *spi.ListCertificateSource
	// documentCRLSource holds CRLs from the document.
	documentCRLSource *spi.ListRevocationSource[revocation.CRL]
	// documentOCSPSource holds OCSP from the document.
	documentOCSPSource *spi.ListRevocationSource[revocation.OCSP]
	// aiaCertificateSources holds certificates collected from AIA.
	aiaCertificateSources *spi.ListCertificateSource
	// revocationCertificateSources holds certificates collected from revocation tokens.
	revocationCertificateSources *spi.ListCertificateSource

	// aiaSource is used to access certificates by AIA.
	aiaSource aia.AIASource
	// remoteOCSPSource is the external OCSP source.
	remoteOCSPSource spi.RevocationSource[revocation.OCSP]
	// remoteCRLSource is the external CRL source.
	remoteCRLSource spi.RevocationSource[revocation.CRL]

	// revocationDataLoadingStrategyFactory builds a strategy deciding how to retrieve
	// revocation data (e.g. CRL or OCSP).
	revocationDataLoadingStrategyFactory RevocationDataLoadingStrategyFactory
	// revocationDataVerifier verifies the validity (i.e. consistency) of a revocation data.
	revocationDataVerifier *RevocationDataVerifier
	// revocationFallback defines whether a revocation data still shall be returned, when
	// validation of obtained revocation tokens failed.
	revocationFallback bool

	// timestampTokenVerifier verifies validity of a TimestampToken.
	timestampTokenVerifier *TimestampTokenVerifier
	// trustAnchorVerifier verifies whether a certificate is a trust anchor.
	trustAnchorVerifier *TrustAnchorVerifier

	// trustedCertSources are the external trusted certificate sources.
	trustedCertSources *spi.ListCertificateSource
	// adjunctCertSources are the external adjunct certificate sources.
	adjunctCertSources *spi.ListCertificateSource

	// checkRevocationForUntrustedChains sets the behavior to follow for revocation retrieving
	// in case of untrusted certificate chains.
	checkRevocationForUntrustedChains bool

	// currentTime is the time at what the validation is carried out.
	currentTime time.Time
}

// signatureValidationContextTokenState pairs a token with its Boolean processing state in
// Java's Map<Token, Boolean>; state == nil means "registered, not yet verified" (Java's null
// value), a pointer to true means verified.
type signatureValidationContextTokenState struct {
	token     model.Token
	processed *bool
}

type signatureValidationContextCertSignatures struct {
	cert       *model.CertificateToken
	signatures []AdvancedSignature
}

type signatureValidationContextCertDates struct {
	cert  *model.CertificateToken
	dates []time.Time
}

type signatureValidationContextTokenIssuer struct {
	token  model.Token
	issuer *model.CertificateToken
}

type signatureValidationContextCertChildren struct {
	parent   *model.CertificateToken
	children []*model.CertificateToken
}

type signatureValidationContextExternalRevocations struct {
	issuer     *model.CertificateToken
	revocation []AnyRevocationToken
}

// signatureValidationContextPOE defines a POE provided to the validation process or obtained
// from processed timestamps. Port of the private static nested class POE.
type signatureValidationContextPOE struct {
	// time is the POE time.
	time time.Time
	// timestampToken is the TimestampToken which provided the POE, when present.
	timestampToken *TimestampToken
}

// newSignatureValidationContextPOEFromTime instantiates a POE from a provided time. Port of
// POE(Date).
func newSignatureValidationContextPOEFromTime(t time.Time) *signatureValidationContextPOE {
	return &signatureValidationContextPOE{time: t}
}

// newSignatureValidationContextPOEFromTimestamp instantiates a POE from a TimestampToken. Port
// of POE(TimestampToken).
func newSignatureValidationContextPOEFromTimestamp(timestampToken *TimestampToken) *signatureValidationContextPOE {
	return &signatureValidationContextPOE{time: timestampToken.CreationDate(), timestampToken: timestampToken}
}

// NewSignatureValidationContext is the default constructor instantiating the object with null
// or empty values and the current time. Port of the no-arg constructor.
func NewSignatureValidationContext() *SignatureValidationContext {
	return NewSignatureValidationContextAtTime(time.Now())
}

// NewSignatureValidationContextAtTime instantiates the object with null or empty values and
// the provided validation time. Port of SignatureValidationContext(Date).
func NewSignatureValidationContextAtTime(validationTime time.Time) *SignatureValidationContext {
	return &SignatureValidationContext{
		documentCertificateSource:    spi.NewListCertificateSource(),
		documentCRLSource:            spi.NewListRevocationSource[revocation.CRL](),
		documentOCSPSource:           spi.NewListRevocationSource[revocation.OCSP](),
		aiaCertificateSources:        spi.NewListCertificateSource(),
		revocationCertificateSources: spi.NewListCertificateSource(),
		poeTimes:                     make(map[string][]*signatureValidationContextPOE),
		currentTime:                  validationTime,
	}
}

// Initialize sets the CertificateVerifier (e.g. using the TSL as list of trusted certificates).
// Port of initialize(CertificateVerifier); panics with Java's message when certificateVerifier
// is nil, mirroring Objects.requireNonNull.
func (c *SignatureValidationContext) Initialize(certificateVerifier CertificateVerifier) {
	if certificateVerifier == nil {
		panic("CertificateVerifier cannot be null!")
	}

	c.certificateVerifier = certificateVerifier
	c.remoteCRLSource = certificateVerifier.CrlSource()
	c.remoteOCSPSource = certificateVerifier.OcspSource()
	c.aiaSource = certificateVerifier.AIASource()
	c.adjunctCertSources = certificateVerifier.AdjunctCertSources()
	c.trustedCertSources = certificateVerifier.TrustedCertSources()
	c.checkRevocationForUntrustedChains = certificateVerifier.IsCheckRevocationForUntrustedChains()
	c.revocationDataLoadingStrategyFactory = certificateVerifier.RevocationDataLoadingStrategyFactory()
	c.revocationDataVerifier = certificateVerifier.RevocationDataVerifier()
	c.revocationFallback = certificateVerifier.IsRevocationFallback()
	c.timestampTokenVerifier = certificateVerifier.TimestampTokenVerifier()
	c.trustAnchorVerifier = certificateVerifier.TrustAnchorVerifier()
}

// getCertificateVerifier gets the CertificateVerifier instance. Panics with Java's message
// when Initialize has not been called yet.
func (c *SignatureValidationContext) getCertificateVerifier() CertificateVerifier {
	if c.certificateVerifier == nil {
		panic("CertificateVerifier shall be initialized! Please use #Initialize(CertificateVerifier) method.")
	}
	return c.certificateVerifier
}

// getRevocationDataVerifier returns an instance of RevocationDataVerifier, instantiating a
// default configuration from a default validation policy, if not defined.
func (c *SignatureValidationContext) getRevocationDataVerifier() *RevocationDataVerifier {
	if c.revocationDataVerifier == nil {
		c.revocationDataVerifier = NewDefaultRevocationDataVerifier()
	}
	if c.revocationDataVerifier.TrustAnchorVerifier() == nil {
		c.revocationDataVerifier.SetTrustAnchorVerifier(c.getTrustAnchorVerifier())
	}
	c.revocationDataVerifier.SetValidationContext(c)
	return c.revocationDataVerifier
}

func (c *SignatureValidationContext) getTimestampTokenVerifier() *TimestampTokenVerifier {
	if c.timestampTokenVerifier == nil {
		c.timestampTokenVerifier = NewDefaultTimestampTokenVerifier()
	}
	if c.timestampTokenVerifier.TrustAnchorVerifier() == nil {
		c.timestampTokenVerifier.SetTrustAnchorVerifier(c.getTrustAnchorVerifier())
	}
	if c.timestampTokenVerifier.RevocationDataVerifier() == nil {
		c.timestampTokenVerifier.SetRevocationDataVerifier(c.getRevocationDataVerifier())
	}
	return c.timestampTokenVerifier
}

func (c *SignatureValidationContext) getTrustAnchorVerifier() *TrustAnchorVerifier {
	if c.trustAnchorVerifier == nil {
		c.trustAnchorVerifier = NewDefaultTrustAnchorVerifier()
	}
	if c.trustAnchorVerifier.TrustedCertificateSource() == nil {
		c.trustAnchorVerifier.SetTrustedCertificateSource(c.trustedCertSources)
	}
	return c.trustAnchorVerifier
}

// AddSignatureForVerification registers a signature (and everything it references) for
// verification. Port of addSignatureForVerification(AdvancedSignature).
func (c *SignatureValidationContext) AddSignatureForVerification(signature AdvancedSignature) {
	if signature == nil {
		return
	}

	c.AddDocumentCertificateSource(signature.CertificateSource())
	c.AddDocumentCRLSource(signature.CRLSource())
	c.AddDocumentOCSPSource(signature.OCSPSource())
	c.registerPOETime(signature.ID(), c.currentTime)

	// Add resolved certificates
	signingCertificate := signature.SigningCertificateToken()
	if signingCertificate != nil {
		c.AddCertificateTokenForVerification(signingCertificate)
	} else {
		certificateValidities := signature.CandidatesForSigningCertificate().CertificateValidityList()
		if utils.IsCollectionNotEmpty(certificateValidities) {
			for _, certificateValidity := range certificateValidities {
				if certificateValidity.IsValid() && certificateValidity.CertificateToken() != nil {
					c.AddCertificateTokenForVerification(certificateValidity.CertificateToken())
				}
			}
		}
	}

	timestamps := signature.AllTimestamps()
	c.prepareTimestamps(timestamps)

	allEvidenceRecords := signature.AllEvidenceRecords()
	c.prepareEvidenceRecords(allEvidenceRecords)

	c.registerCertChainUsage(signature) // to be done after timestamp POE extraction

	c.processedSignatures = signatureValidationContextAppendSignature(c.processedSignatures, signature)

	counterSignatures := signature.CounterSignatures()
	c.prepareCounterSignatures(counterSignatures)
}

// AddDocumentCertificateSource adds a CertificateSource to the document certificate source.
// Port of addDocumentCertificateSource(CertificateSource).
func (c *SignatureValidationContext) AddDocumentCertificateSource(certificateSource spi.CertificateSource) {
	c.addCertificateSource(c.documentCertificateSource, certificateSource)
}

// AddDocumentCertificateSourceFromList adds every source of a ListCertificateSource to the
// document certificate source. Port of the addDocumentCertificateSource(ListCertificateSource)
// overload; Go has no overloading.
func (c *SignatureValidationContext) AddDocumentCertificateSourceFromList(listCertificateSource *spi.ListCertificateSource) {
	for _, certificateSource := range listCertificateSource.Sources() {
		c.AddDocumentCertificateSource(certificateSource)
	}
}

// addCertificateSource adds certificateSourceToAdd to the given listCertificateSource.
func (c *SignatureValidationContext) addCertificateSource(listCertificateSource *spi.ListCertificateSource, certificateSourceToAdd spi.CertificateSource) {
	if listCertificateSource.Add(certificateSourceToAdd) {
		// add all existing equivalent certificates for the validation
		allCertificateSources := c.GetAllCertificateSources()
		for _, certificateToken := range certificateSourceToAdd.Certificates() {
			c.addEquivalentCertificates(certificateToken, allCertificateSources)
		}
	}
}

func (c *SignatureValidationContext) addEquivalentCertificates(certificateToken *model.CertificateToken, certificateSource spi.CertificateSource) {
	equivalentCertificates := certificateSource.ByEntityKey(certificateToken.EntityKey())
	for _, equivalentCertificate := range equivalentCertificates {
		if certificateToken.DSSIDAsString() != equivalentCertificate.DSSIDAsString() {
			c.AddCertificateTokenForVerification(equivalentCertificate)
		}
	}
}

// AddDocumentCRLSource adds a CRL OfflineRevocationSource to the document CRL source. Port of
// addDocumentCRLSource(OfflineRevocationSource<CRL>).
func (c *SignatureValidationContext) AddDocumentCRLSource(crlSource spi.OfflineRevocationSource[revocation.CRL]) {
	c.documentCRLSource.Add(crlSource)
}

// AddDocumentCRLSourceFromList adds every source of a ListRevocationSource<CRL> to the
// document CRL source. Port of the addDocumentCRLSource(ListRevocationSource<CRL>) overload.
func (c *SignatureValidationContext) AddDocumentCRLSourceFromList(crlSource *spi.ListRevocationSource[revocation.CRL]) {
	c.documentCRLSource.AddAllFrom(crlSource)
}

// AddDocumentOCSPSource adds an OCSP OfflineRevocationSource to the document OCSP source. Port
// of addDocumentOCSPSource(OfflineRevocationSource<OCSP>).
func (c *SignatureValidationContext) AddDocumentOCSPSource(ocspSource spi.OfflineRevocationSource[revocation.OCSP]) {
	c.documentOCSPSource.Add(ocspSource)
}

// AddDocumentOCSPSourceFromList adds every source of a ListRevocationSource<OCSP> to the
// document OCSP source. Port of the addDocumentOCSPSource(ListRevocationSource<OCSP>) overload.
func (c *SignatureValidationContext) AddDocumentOCSPSourceFromList(ocspSource *spi.ListRevocationSource[revocation.OCSP]) {
	c.documentOCSPSource.AddAllFrom(ocspSource)
}

func (c *SignatureValidationContext) prepareTimestamps(timestampTokens []*TimestampToken) {
	if utils.IsCollectionNotEmpty(timestampTokens) {
		for _, timestampToken := range timestampTokens {
			c.AddTimestampTokenForVerification(timestampToken)
		}
	}
}

func (c *SignatureValidationContext) prepareEvidenceRecords(evidenceRecords []EvidenceRecord) {
	if utils.IsCollectionNotEmpty(evidenceRecords) {
		for _, evidenceRecord := range evidenceRecords {
			c.AddEvidenceRecordForVerification(evidenceRecord)
		}
	}
}

func (c *SignatureValidationContext) registerCertChainUsage(signature AdvancedSignature) {
	signingCertificate := signature.SigningCertificateToken()
	if signingCertificate != nil {
		entry := signatureValidationContextFindCertSignatures(c.signCertificateSignaturesUsage, signingCertificate)
		if entry == nil {
			c.signCertificateSignaturesUsage = append(c.signCertificateSignaturesUsage, signatureValidationContextCertSignatures{cert: signingCertificate})
			entry = &c.signCertificateSignaturesUsage[len(c.signCertificateSignaturesUsage)-1]
		}
		if !signatureValidationContextContainsSignature(entry.signatures, signature) {
			entry.signatures = append(entry.signatures, signature)
		}
	}
}

func (c *SignatureValidationContext) prepareCounterSignatures(counterSignatures []AdvancedSignature) {
	if utils.IsCollectionNotEmpty(counterSignatures) {
		for _, counterSignature := range counterSignatures {
			c.AddSignatureForVerification(counterSignature)
		}
	}
}

// GetCurrentTime returns the time at which the validation is carried out. Port of
// getCurrentTime().
func (c *SignatureValidationContext) GetCurrentTime() time.Time {
	return c.currentTime
}

// getNotYetVerifiedToken returns a token to verify, nil if there is no more token to verify.
func (c *SignatureValidationContext) getNotYetVerifiedToken() model.Token {
	revocationToken := c.getNotYetVerifiedRevocationToken()
	if revocationToken != nil {
		return revocationToken
	}
	for i := range c.tokensToProcess {
		entry := &c.tokensToProcess[i]
		if entry.processed == nil {
			trueVal := true
			entry.processed = &trueVal
			return entry.token
		}
	}
	return nil
}

// getNotYetVerifiedRevocationToken returns a revocation token to verify, nil if there is no
// more token to verify. As revocation tokens may be added dynamically, the method is executed
// as part of getNotYetVerifiedToken, with priority, in order to be able to identify a correct
// certificate chain.
func (c *SignatureValidationContext) getNotYetVerifiedRevocationToken() AnyRevocationToken {
	if utils.IsCollectionEmpty(c.processedRevocations) {
		return nil
	}
	revocationTokens := append([]AnyRevocationToken(nil), c.processedRevocations...)
	for _, revocationToken := range revocationTokens {
		if !c.isYetVerified(revocationToken) {
			return revocationToken
		}
	}
	return nil
}

// getNotYetVerifiedTimestamp returns a timestamp token to verify, nil if there is no more
// token to verify.
func (c *SignatureValidationContext) getNotYetVerifiedTimestamp() *TimestampToken {
	if utils.IsCollectionEmpty(c.processedTimestamps) {
		return nil
	}
	sortedTimestampTokens := append([]*TimestampToken(nil), c.processedTimestamps...)
	comparator := NewTimestampTokenComparator()
	sort.SliceStable(sortedTimestampTokens, func(i, j int) bool {
		return comparator.Compare(sortedTimestampTokens[i], sortedTimestampTokens[j]) < 0
	})
	// start processing from the freshest timestamp
	for i, j := 0, len(sortedTimestampTokens)-1; i < j; i, j = i+1, j-1 {
		sortedTimestampTokens[i], sortedTimestampTokens[j] = sortedTimestampTokens[j], sortedTimestampTokens[i]
	}
	for _, timestampToken := range sortedTimestampTokens {
		if !c.isYetVerified(timestampToken) {
			return timestampToken
		}
	}
	return nil
}

func (c *SignatureValidationContext) getNotYetVerifiedTokenFromChain(certChain []model.Token) model.Token {
	for _, token := range certChain {
		if !c.isYetVerified(token) {
			return token
		}
	}
	return nil
}

func (c *SignatureValidationContext) isYetVerified(token model.Token) bool {
	for i := range c.tokensToProcess {
		entry := &c.tokensToProcess[i]
		if entry.token == token {
			processed := entry.processed != nil && *entry.processed
			if !processed {
				trueVal := true
				entry.processed = &trueVal
				return false
			}
			return true
		}
	}
	trueVal := true
	c.tokensToProcess = append(c.tokensToProcess, signatureValidationContextTokenState{token: token, processed: &trueVal})
	return false
}

func (c *SignatureValidationContext) getOrderedCertificateChains() (map[*model.CertificateToken][]*model.CertificateToken, error) {
	order := spi.NewCertificateReorderer(c.processedCertificates)
	return order.OrderedCertificateChains()
}

// getCertChain builds the complete certificate chain from the given token.
func (c *SignatureValidationContext) getCertChain(token model.Token) []model.Token {
	var chain []model.Token
	issuerCertificateToken := token
	for {
		chain = append(chain, issuerCertificateToken)
		// getIssuerWithSource returns a concrete *model.CertificateToken; comparing that against
		// nil before it is assigned into the model.Token interface variable matters; the
		// interface's own dynamic type would otherwise stay *model.CertificateToken with a nil
		// value, so the loop's "issuerCertificateToken == nil" guard below would never see a nil
		// interface and signatureValidationContextContainsToken would call DSSIDAsString() on a
		// nil receiver.
		issuer := c.getIssuerWithSource(issuerCertificateToken, c.getTokenCertificateSource(token))
		if issuer == nil {
			break
		}
		issuerCertificateToken = issuer
		if signatureValidationContextContainsToken(chain, issuerCertificateToken) {
			break
		}
	}
	return chain
}

func (c *SignatureValidationContext) getTokenCertificateSource(token model.Token) spi.CertificateSource {
	switch t := token.(type) {
	case *spi.OCSPToken:
		return t.CertificateSource()
	case *TimestampToken:
		return t.CertificateSource()
	default:
		// other tokens do not have their own source
		return nil
	}
}

// getCertificateTokenChain computes the certificate chain for the given token without
// including the current token in the chain, when it is not an instance of CertificateToken.
func (c *SignatureValidationContext) getCertificateTokenChain(token model.Token) []*model.CertificateToken {
	var certificateChain []*model.CertificateToken
	for _, chainItem := range c.getCertChain(token) {
		if certificateToken, ok := chainItem.(*model.CertificateToken); ok {
			certificateChain = append(certificateChain, certificateToken)
		}
	}
	return certificateChain
}

func (c *SignatureValidationContext) getIssuer(token model.Token) *model.CertificateToken {
	return c.getIssuerWithSource(token, nil)
}

func (c *SignatureValidationContext) getIssuerWithSource(token model.Token, certificateSource spi.CertificateSource) *model.CertificateToken {
	// Return cached value
	issuerCertificateToken := c.getIssuerFromProcessedCertificates(token)
	if issuerCertificateToken != nil {
		// ensure equivalent certificates are processed
		if certificateSource != nil {
			c.addEquivalentCertificates(issuerCertificateToken, certificateSource)
		}
		return issuerCertificateToken
	}

	// See DSS-3720
	if certToken, ok := token.(*model.CertificateToken); ok && token.IsSelfSigned() {
		return certToken
	}

	// Find issuer candidates from a particular certificate source
	var candidates map[string]*model.CertificateToken

	// Avoid repeating over stateless sources
	if !signatureValidationContextTokenIssuerMapContainsKey(c.tokenIssuerMap, token) {
		if certificateSource == nil {
			// OCSP or Timestamp
			certificateSource = c.getTokenCertificateSource(token)
		}

		if certificateSource != nil {
			candidates = c.getIssuersFromSource(token, certificateSource)
		}

		// Find issuer candidates from document sources
		if utils.IsMapEmpty(candidates) {
			candidates = c.getIssuersFromSources(token, c.documentCertificateSource)
		}
	}

	// Find issuer candidates from all sources
	allCertificateSources := c.GetAllCertificateSources()
	if utils.IsMapEmpty(candidates) {
		candidates = c.getIssuersFromSources(token, allCertificateSources)
	}

	// Find issuer from provided certificate tokens
	if utils.IsMapEmpty(candidates) {
		candidates = make(map[string]*model.CertificateToken)
		for _, cert := range c.processedCertificates {
			candidates[cert.DSSIDAsString()] = cert
		}
		for _, cert := range c.documentCertificateSource.Certificates() {
			candidates[cert.DSSIDAsString()] = cert
		}
	}

	candidates = c.ensureCandidatesFromProcessedCertificates(candidates)

	issuerCertificateToken = spi.NewTokenIssuerSelector(token, signatureValidationContextMapValues(candidates)).Issuer()

	// Request AIA only when no issuer has been found yet
	if issuerCertificateToken == nil && c.aiaSource != nil {
		if certToken, ok := token.(*model.CertificateToken); ok && !signatureValidationContextTokenIssuerMapContainsKey(c.tokenIssuerMap, token) {
			aiaCertificateSource := aia.AIACertificateSourceForCertificateToken(certToken, c.aiaSource)
			issuerCertificateToken = aiaCertificateSource.IssuerFromAIA()
			c.addCertificateSource(c.aiaCertificateSources, aiaCertificateSource)
		}
	}

	if issuerCertificateToken == nil {
		if ocspToken, ok := token.(*spi.OCSPToken); ok {
			issuerCertificateToken = c.getOCSPIssuer(ocspToken, allCertificateSources)
		}
	}

	if issuerCertificateToken == nil {
		if timestampToken, ok := token.(*TimestampToken); ok {
			issuerCertificateToken = c.getTSACertificate(timestampToken, allCertificateSources)
		}
	}

	if issuerCertificateToken != nil {
		c.AddCertificateTokenForVerification(issuerCertificateToken)
	}

	// Cache the result (successful or unsuccessful)
	c.addToCacheMap(token, issuerCertificateToken)

	return issuerCertificateToken
}

// ensureCandidatesFromProcessedCertificates ensures the processed certificates are used on
// validation, in order to avoid cases when another certificate is being updated (not provided
// to the validation).
func (c *SignatureValidationContext) ensureCandidatesFromProcessedCertificates(candidates map[string]*model.CertificateToken) map[string]*model.CertificateToken {
	if utils.IsMapEmpty(candidates) {
		return map[string]*model.CertificateToken{}
	}
	result := make(map[string]*model.CertificateToken, len(candidates))
	for _, certificateToken := range candidates {
		resolved := certificateToken
		for _, processedCertificate := range c.processedCertificates {
			if certificateToken.Equals(processedCertificate) {
				resolved = processedCertificate
				break
			}
		}
		result[resolved.DSSIDAsString()] = resolved
	}
	return result
}

func (c *SignatureValidationContext) addToCacheMap(token model.Token, issuerCertificateToken *model.CertificateToken) {
	signatureValidationContextPutTokenIssuer(&c.tokenIssuerMap, token, issuerCertificateToken)

	if certificateToken, ok := token.(*model.CertificateToken); ok {
		entry := signatureValidationContextFindCertChildren(c.certificateChildrenMap, issuerCertificateToken)
		if entry == nil {
			c.certificateChildrenMap = append(c.certificateChildrenMap, signatureValidationContextCertChildren{parent: issuerCertificateToken})
			entry = &c.certificateChildrenMap[len(c.certificateChildrenMap)-1]
		}
		if !signatureValidationContextContainsCertificate(entry.children, certificateToken) {
			entry.children = append(entry.children, certificateToken)
		}
	}
}

func (c *SignatureValidationContext) getIssuerFromProcessedCertificates(token model.Token) *model.CertificateToken {
	issuerCertificateToken := signatureValidationContextGetTokenIssuer(c.tokenIssuerMap, token)
	// isSignedBy(...) check is required when a certificate is present in different sources in
	// order to instantiate a public key of the signer
	if issuerCertificateToken != nil && (token.PublicKeyOfTheSigner() != nil || token.IsSignedByToken(issuerCertificateToken)) {
		return issuerCertificateToken
	}
	return nil
}

// GetAllCertificateSources returns a merged ListCertificateSource of every certificate source
// known to the context. Port of getAllCertificateSources().
func (c *SignatureValidationContext) GetAllCertificateSources() *spi.ListCertificateSource {
	allCertificateSources := spi.NewListCertificateSource()
	allCertificateSources.AddAll(c.documentCertificateSource)
	allCertificateSources.AddAll(c.revocationCertificateSources)
	allCertificateSources.AddAll(c.aiaCertificateSources)
	allCertificateSources.AddAll(c.adjunctCertSources)
	allCertificateSources.AddAll(c.trustedCertSources)
	return allCertificateSources
}

// GetDocumentCertificateSource returns the document certificate source. Port of
// getDocumentCertificateSource().
func (c *SignatureValidationContext) GetDocumentCertificateSource() *spi.ListCertificateSource {
	return c.documentCertificateSource
}

// GetDocumentCRLSource returns the document CRL source. Port of getDocumentCRLSource().
func (c *SignatureValidationContext) GetDocumentCRLSource() *spi.ListRevocationSource[revocation.CRL] {
	return c.documentCRLSource
}

// GetDocumentOCSPSource returns the document OCSP source. Port of getDocumentOCSPSource().
func (c *SignatureValidationContext) GetDocumentOCSPSource() *spi.ListRevocationSource[revocation.OCSP] {
	return c.documentOCSPSource
}

func (c *SignatureValidationContext) getIssuersFromSources(token model.Token, allCertificateSources *spi.ListCertificateSource) map[string]*model.CertificateToken {
	if token.IssuerEntityKey() != nil {
		return allCertificateSources.ByEntityKey(token.IssuerEntityKey())
	} else if token.PublicKeyOfTheSigner() != nil {
		return allCertificateSources.ByPublicKey(token.PublicKeyOfTheSigner())
	} else if token.IssuerX500Principal() != nil {
		return allCertificateSources.BySubject(model.NewX500PrincipalHelper(token.IssuerX500Principal()))
	}
	return map[string]*model.CertificateToken{}
}

func (c *SignatureValidationContext) getIssuersFromSource(token model.Token, certificateSource spi.CertificateSource) map[string]*model.CertificateToken {
	if token.IssuerEntityKey() != nil {
		return certificateSource.ByEntityKey(token.IssuerEntityKey())
	} else if token.PublicKeyOfTheSigner() != nil {
		return certificateSource.ByPublicKey(token.PublicKeyOfTheSigner())
	} else if token.IssuerX500Principal() != nil {
		return certificateSource.BySubject(model.NewX500PrincipalHelper(token.IssuerX500Principal()))
	}
	return map[string]*model.CertificateToken{}
}

// certificateRefsSource narrows a RevocationCertificateSource down to the concrete sources
// (OCSPCertificateSource) that expose AllCertificateRefs; Java reaches it through a covariant
// override of getCertificateSource(), which Go cannot express (see OCSPToken.CertificateSource()).
type certificateRefsSource interface {
	AllCertificateRefs() []*spi.CertificateRef
}

func (c *SignatureValidationContext) getOCSPIssuer(token *spi.OCSPToken, allCertificateSources *spi.ListCertificateSource) *model.CertificateToken {
	var signingCertificateRefs []*spi.CertificateRef
	if refsSource, ok := token.CertificateSource().(certificateRefsSource); ok {
		signingCertificateRefs = refsSource.AllCertificateRefs()
	}
	if utils.CollectionSize(signingCertificateRefs) == 1 {
		signingCertificateRef := signingCertificateRefs[0]
		responderID := signingCertificateRef.ResponderId()
		if responderID != nil {
			issuerCandidates := make(map[string]*model.CertificateToken)
			if responderID.Ski() != nil {
				for id, cert := range allCertificateSources.BySki(responderID.Ski()) {
					issuerCandidates[id] = cert
				}
			}
			if responderID.X500Principal() != nil {
				for id, cert := range allCertificateSources.BySubject(model.NewX500PrincipalHelper(responderID.X500Principal())) {
					issuerCandidates[id] = cert
				}
			}
			return spi.NewTokenIssuerSelector(token, signatureValidationContextMapValues(issuerCandidates)).Issuer()
		}
	}
	return nil
}

func (c *SignatureValidationContext) getTSACertificate(timestamp *TimestampToken, allCertificateSources *spi.ListCertificateSource) *model.CertificateToken {
	candidatesForSigningCertificate := timestamp.CandidatesForSigningCertificate()
	theBestCandidate := candidatesForSigningCertificate.TheBestCandidate()
	if theBestCandidate != nil {
		issuerCandidates := make(map[string]*model.CertificateToken)
		timestampSigner := theBestCandidate.CertificateToken()
		if timestampSigner == nil {
			for id, cert := range allCertificateSources.BySignerIdentifier(theBestCandidate.SignerInfo()) {
				issuerCandidates[id] = cert
			}
		} else {
			issuerCandidates[timestampSigner.DSSIDAsString()] = timestampSigner
		}
		return spi.NewTokenIssuerSelector(timestamp, signatureValidationContextMapValues(issuerCandidates)).Issuer()
	}
	return nil
}

// addTokenForVerification adds a new token to the list of tokens to verify only if it was not
// already verified. Returns true if the token was not yet verified, false otherwise.
func (c *SignatureValidationContext) addTokenForVerification(token model.Token) bool {
	if token == nil {
		return false
	}

	for i := range c.tokensToProcess {
		if c.tokensToProcess[i].token == token {
			return false
		}
	}

	c.tokensToProcess = append(c.tokensToProcess, signatureValidationContextTokenState{token: token})
	c.registerPOETime(token.DSSIDAsString(), c.currentTime)
	return true
}

// AddRevocationTokenForVerification registers a revocation token for verification. Port of
// addRevocationTokenForVerification(RevocationToken<?>).
func (c *SignatureValidationContext) AddRevocationTokenForVerification(revocationToken AnyRevocationToken) {
	if c.addTokenForVerification(revocationToken) {
		revocationCertificateSource := revocationToken.CertificateSource()
		if revocationCertificateSource != nil {
			c.addCertificateSource(c.revocationCertificateSources, revocationCertificateSource)
		}

		issuerCertificateToken := revocationToken.IssuerCertificateToken()
		if issuerCertificateToken != nil {
			c.AddCertificateTokenForVerification(issuerCertificateToken)
		}

		c.processedRevocations = signatureValidationContextAppendRevocation(c.processedRevocations, revocationToken)
	}
}

// AddCertificateTokenForVerification registers a certificate token for verification. Port of
// addCertificateTokenForVerification(CertificateToken).
func (c *SignatureValidationContext) AddCertificateTokenForVerification(certificateToken *model.CertificateToken) {
	if c.addTokenForVerification(certificateToken) {
		c.processedCertificates = signatureValidationContextAppendCertificate(c.processedCertificates, certificateToken)
	}
}

// AddTimestampTokenForVerification registers a timestamp token for verification. Port of
// addTimestampTokenForVerification(TimestampToken).
func (c *SignatureValidationContext) AddTimestampTokenForVerification(timestampToken *TimestampToken) {
	if c.addTokenForVerification(timestampToken) {
		c.AddDocumentCertificateSource(timestampToken.CertificateSource())
		c.AddDocumentCRLSource(timestampToken.CRLSource())
		c.AddDocumentOCSPSource(timestampToken.OCSPSource())

		certificateValidities := timestampToken.CandidatesForSigningCertificate().CertificateValidityList()
		if utils.IsCollectionNotEmpty(certificateValidities) {
			for _, certificateValidity := range certificateValidities {
				if certificateValidity.IsValid() && certificateValidity.CertificateToken() != nil {
					c.AddCertificateTokenForVerification(certificateValidity.CertificateToken())
				}
			}
		}

		c.registerTimestampUsageDate(timestampToken)

		c.processedTimestamps = signatureValidationContextAppendTimestamp(c.processedTimestamps, timestampToken)
	}
}

func (c *SignatureValidationContext) registerTimestampUsageDate(timestampToken *TimestampToken) {
	tsaCertificate := c.getTSACertificate(timestampToken, c.GetAllCertificateSources())
	if tsaCertificate == nil {
		return
	}

	tsaCertificateChain := signatureValidationContextToCertificateTokenChain(c.getCertChain(timestampToken))
	usageDate := timestampToken.CreationDate()
	for _, cert := range tsaCertificateChain {
		entry := signatureValidationContextFindCertDates(c.timestampCertChainDates, cert)
		if entry == nil {
			c.timestampCertChainDates = append(c.timestampCertChainDates, signatureValidationContextCertDates{cert: cert})
			entry = &c.timestampCertChainDates[len(c.timestampCertChainDates)-1]
		}
		if !signatureValidationContextContainsTime(entry.dates, usageDate) {
			entry.dates = append(entry.dates, usageDate)
		}
	}
}

func (c *SignatureValidationContext) registerTimestampPOE(timestampToken *TimestampToken) {
	if c.isTimestampValid(timestampToken) {
		for _, timestampedReference := range timestampToken.TimestampedReferences() {
			c.registerPOETimestamp(timestampedReference.ObjectId(), timestampToken)
		}
	}
}

// isTimestampValid verifies whether a timestampToken is valid and can be used as a valid POE
// for covered objects.
func (c *SignatureValidationContext) isTimestampValid(timestampToken *TimestampToken) bool {
	certificateTokenChain := c.getCertificateTokenChain(timestampToken)
	lowestPOETime := c.getLowestPOETimeForToken(timestampToken)
	return c.getTimestampTokenVerifier().IsAcceptableWithChainAt(timestampToken, certificateTokenChain, lowestPOETime)
}

func (c *SignatureValidationContext) registerPOETimestamp(tokenID string, timestampToken *TimestampToken) {
	c.poeTimes[tokenID] = append(c.poeTimes[tokenID], newSignatureValidationContextPOEFromTimestamp(timestampToken))
}

func (c *SignatureValidationContext) registerPOETime(tokenID string, poeTime time.Time) {
	c.poeTimes[tokenID] = append(c.poeTimes[tokenID], newSignatureValidationContextPOEFromTime(poeTime))
}

// AddEvidenceRecordForVerification registers an evidence record for verification. Port of
// addEvidenceRecordForVerification(EvidenceRecord).
func (c *SignatureValidationContext) AddEvidenceRecordForVerification(evidenceRecord EvidenceRecord) {
	c.AddDocumentCertificateSource(evidenceRecord.CertificateSource())
	c.AddDocumentCRLSource(evidenceRecord.CRLSource())
	c.AddDocumentOCSPSource(evidenceRecord.OCSPSource())
	c.prepareTimestamps(evidenceRecord.Timestamps())

	c.processedEvidenceRecords = signatureValidationContextAppendEvidenceRecord(c.processedEvidenceRecords, evidenceRecord)
}

// Validate runs the validation of every not-yet-verified token registered in the context. Port
// of validate().
func (c *SignatureValidationContext) Validate() {
	timestampToken := c.getNotYetVerifiedTimestamp()
	for timestampToken != nil {
		c.validateTimestamp(timestampToken)
		timestampToken = c.getNotYetVerifiedTimestamp()
	}

	token := c.getNotYetVerifiedToken()
	for token != nil {
		c.validateToken(token)
		token = c.getNotYetVerifiedToken()
	}
}

func (c *SignatureValidationContext) validateTimestamp(timestampToken *TimestampToken) {
	certChain := c.getCertChain(timestampToken)
	token := c.getNotYetVerifiedTokenFromChain(certChain)
	for token != nil {
		c.validateToken(token)
		token = c.getNotYetVerifiedTokenFromChain(certChain)
	}
	c.registerTimestampPOE(timestampToken) // POE is extracted after TST validation
}

func (c *SignatureValidationContext) validateToken(token model.Token) {
	// extract the certificate chain and add missing tokens for verification
	certChain := c.getCertChain(token)
	if utils.CollectionSize(certChain) > 1 { // ensure certificate chain is processed
		certChainToken := c.getNotYetVerifiedTokenFromChain(certChain)
		if certChainToken != nil {
			c.validateToken(certChainToken)
		}
	}
	if certificateToken, ok := token.(*model.CertificateToken); ok {
		c.findRevocationData(certificateToken, certChain)
	}
}

func (c *SignatureValidationContext) validateTokenIfNeeded(token model.Token) {
	c.addTokenForVerification(token) // ensure the token is added to the validation context
	if !c.isYetVerified(token) {
		c.validateToken(token)
	}
}

// findRevocationData retrieves the revocation data from signature (if exists) or from the
// online sources. The issuer certificate must be provided, the underlying library needs it to
// build the request.
func (c *SignatureValidationContext) findRevocationData(certToken *model.CertificateToken, certChain []model.Token) []AnyRevocationToken {
	if c.isRevocationDataNotRequired(certToken, c.getLowestPOETimeForToken(certToken)) {
		return nil
	}

	issuerToken := c.getIssuer(certToken)
	if issuerToken == nil {
		return nil
	}

	var revocations []AnyRevocationToken

	// ALL Embedded revocation data
	crlTokens, err := c.documentCRLSource.RevocationTokens(certToken, issuerToken)
	if err == nil {
		for _, revocationToken := range crlTokens {
			revocations = signatureValidationContextAppendRevocation(revocations, revocationToken)
			c.AddRevocationTokenForVerification(revocationToken)
		}
	}

	ocspTokens, err := c.documentOCSPSource.RevocationTokens(certToken, issuerToken)
	if err == nil {
		for _, revocationToken := range ocspTokens {
			revocations = signatureValidationContextAppendRevocation(revocations, revocationToken)
			c.AddRevocationTokenForVerification(revocationToken)
			c.AddDocumentCertificateSource(revocationToken.CertificateSource()) // applicable only for OCSP
		}
	}

	// add processed revocation tokens
	for _, revocationToken := range c.getRelatedRevocationTokens(certToken) {
		revocations = signatureValidationContextAppendRevocation(revocations, revocationToken)
	}

	externalRevocationTokens := c.getExternalRevocationTokens(certToken, issuerToken)
	for _, revocationToken := range externalRevocationTokens {
		revocations = signatureValidationContextAppendRevocation(revocations, revocationToken)
	}

	if (c.remoteOCSPSource != nil || c.remoteCRLSource != nil) &&
		(utils.IsCollectionEmpty(revocations) || (utils.IsCollectionEmpty(externalRevocationTokens) && c.isRevocationDataRefreshNeeded(certToken, revocations))) {
		if c.checkRevocationForUntrustedChains || c.containsTrustAnchor(certChain) {
			trustAnchor, _ := c.getFirstTrustAnchor(certChain).(*model.CertificateToken)

			// Fetch OCSP or CRL from online sources
			onlineRevocationToken := c.getRevocationToken(certToken, issuerToken, trustAnchor)

			// Check if the obtained revocation is not yet present
			if onlineRevocationToken != nil && !signatureValidationContextContainsRevocation(revocations, onlineRevocationToken) {
				revocations = append(revocations, onlineRevocationToken)
				c.AddRevocationTokenForVerification(onlineRevocationToken)
				c.linkRevocationToOtherCertificates(onlineRevocationToken, issuerToken)
			}
		}
	}

	return revocations
}

func (c *SignatureValidationContext) getExternalRevocationTokens(certToken, issuerCertificateToken *model.CertificateToken) []AnyRevocationToken {
	var result []AnyRevocationToken
	if issuerCertificateToken != nil {
		entry := signatureValidationContextFindExternalRevocations(c.externalRevocationTokensMap, issuerCertificateToken)
		if entry != nil {
			for _, revocationToken := range entry.revocation {
				if crlToken, ok := revocationToken.(*spi.CRLToken); ok && !certToken.Equals(revocationToken.RelatedCertificate()) {
					newCRLToken, err := spi.NewCRLToken(certToken, crlToken.CrlValidity())
					if err == nil {
						newCRLToken.SetExternalOrigin(crlToken.ExternalOrigin())
						newCRLToken.SetSourceURL(crlToken.SourceURL())
						c.AddRevocationTokenForVerification(newCRLToken)
						result = append(result, newCRLToken)
					}
				}
			}
		}
	}
	return result
}

// GetRevocationData returns the revocation data related to the given certificate. Port of
// getRevocationData(CertificateToken).
func (c *SignatureValidationContext) GetRevocationData(certificateToken *model.CertificateToken) []AnyRevocationToken {
	c.validateTokenIfNeeded(certificateToken)
	var result []AnyRevocationToken
	for _, revocationToken := range c.processedRevocations {
		if utils.AreStringsEqual(certificateToken.DSSIDAsString(), revocationToken.RelatedCertificateID()) {
			result = append(result, revocationToken)
		}
	}
	return result
}

func (c *SignatureValidationContext) containsTrustAnchor(certChain []model.Token) bool {
	return c.getFirstTrustAnchor(certChain) != nil
}

func (c *SignatureValidationContext) getFirstTrustAnchor(certChain []model.Token) model.Token {
	if utils.IsCollectionNotEmpty(certChain) {
		for _, token := range certChain {
			if c.isTrustedAtUsageTime(token) {
				return token
			}
		}
	}
	return nil
}

func (c *SignatureValidationContext) containsTrustAnchorAtTime(certChain []model.Token, controlTime time.Time) bool {
	if utils.IsCollectionNotEmpty(certChain) {
		for _, token := range certChain {
			if c.isTrustedAtTime(token, controlTime, "") {
				return true
			}
		}
	}
	return false
}

func (c *SignatureValidationContext) linkRevocationToOtherCertificates(revocationToken AnyRevocationToken, issuerCertificateToken *model.CertificateToken) {
	entry := signatureValidationContextFindExternalRevocations(c.externalRevocationTokensMap, issuerCertificateToken)
	if entry == nil {
		c.externalRevocationTokensMap = append(c.externalRevocationTokensMap, signatureValidationContextExternalRevocations{issuer: issuerCertificateToken})
		entry = &c.externalRevocationTokensMap[len(c.externalRevocationTokensMap)-1]
	}
	entry.revocation = append(entry.revocation, revocationToken)
}

func (c *SignatureValidationContext) getRevocationToken(certificateToken, issuerCertificate, trustAnchor *model.CertificateToken) AnyRevocationToken {
	// configure the CompositeRevocationSource
	var currentOCSPSource spi.RevocationSource[revocation.OCSP]
	var currentCRLSource spi.RevocationSource[revocation.CRL]
	if !c.trustedCertSources.IsEmpty() && trustAnchor != nil {
		currentOCSPSource = c.instantiateOCSPWithTrustServices(trustAnchor)
		currentCRLSource = c.instantiateCRLWithTrustServices(trustAnchor)
	} else {
		currentOCSPSource = c.remoteOCSPSource
		currentCRLSource = c.remoteCRLSource
	}

	// fetch the data
	revocationDataLoadingStrategy := c.revocationDataLoadingStrategyFactory.Create()
	revocationDataLoadingStrategy.setCrlSource(currentCRLSource)
	revocationDataLoadingStrategy.setOcspSource(currentOCSPSource)
	revocationDataLoadingStrategy.setRevocationDataVerifier(c.getRevocationDataVerifier())
	revocationDataLoadingStrategy.setFallbackEnabled(c.revocationFallback)
	return revocationDataLoadingStrategy.RevocationToken(certificateToken, issuerCertificate)
}

func (c *SignatureValidationContext) instantiateOCSPWithTrustServices(trustAnchor *model.CertificateToken) spi.RevocationSource[revocation.OCSP] {
	alternativeOCSPUrls := c.getAlternativeOCSPUrls(trustAnchor)
	if support, ok := c.remoteOCSPSource.(spi.RevocationSourceAlternateUrlsSupport[revocation.OCSP]); ok && utils.IsCollectionNotEmpty(alternativeOCSPUrls) {
		return spi.NewAlternateUrlsSourceAdapter[revocation.OCSP](support, alternativeOCSPUrls)
	}
	return c.remoteOCSPSource
}

func (c *SignatureValidationContext) instantiateCRLWithTrustServices(trustAnchor *model.CertificateToken) spi.RevocationSource[revocation.CRL] {
	alternativeCRLUrls := c.getAlternativeCRLUrls(trustAnchor)
	if support, ok := c.remoteCRLSource.(spi.RevocationSourceAlternateUrlsSupport[revocation.CRL]); ok && utils.IsCollectionNotEmpty(alternativeCRLUrls) {
		return spi.NewAlternateUrlsSourceAdapter[revocation.CRL](support, alternativeCRLUrls)
	}
	return c.remoteCRLSource
}

func (c *SignatureValidationContext) getAlternativeOCSPUrls(trustAnchor *model.CertificateToken) []string {
	var alternativeOCSPUrls []string
	for _, certificateSource := range c.trustedCertSources.Sources() {
		if trustedCertSource, ok := certificateSource.(spi.TrustedCertificateSource); ok {
			alternativeOCSPUrls = append(alternativeOCSPUrls, trustedCertSource.AlternativeOCSPUrls(trustAnchor)...)
		}
	}
	return alternativeOCSPUrls
}

func (c *SignatureValidationContext) getAlternativeCRLUrls(trustAnchor *model.CertificateToken) []string {
	var alternativeCRLUrls []string
	for _, certificateSource := range c.trustedCertSources.Sources() {
		if trustedCertSource, ok := certificateSource.(spi.TrustedCertificateSource); ok {
			alternativeCRLUrls = append(alternativeCRLUrls, trustedCertSource.AlternativeCRLUrls(trustAnchor)...)
		}
	}
	return alternativeCRLUrls
}

// CheckAllRequiredRevocationDataPresent reports whether every processed certificate has the
// required revocation data present. Port of checkAllRequiredRevocationDataPresent().
func (c *SignatureValidationContext) CheckAllRequiredRevocationDataPresent() bool {
	status, _ := c.allRequiredRevocationDataPresent()
	return status.IsEmpty()
}

// allRequiredRevocationDataPresent returns the status of the required-revocation-data-present
// check.
func (c *SignatureValidationContext) allRequiredRevocationDataPresent() (*TokenStatus, error) {
	status := NewTokenStatus()
	orderedCertificateChains, err := c.getOrderedCertificateChains()
	if err != nil {
		return status, err
	}
	for _, orderedCertChain := range orderedCertificateChains {
		c.checkRevocationForCertificateChainAgainstBestSignatureTime(orderedCertChain, time.Time{}, status, "")
	}
	if !status.IsEmpty() {
		status.SetMessage("Revocation data is missing for one or more certificate(s).")
	}
	return status, nil
}

func (c *SignatureValidationContext) checkRevocationForCertificateChainAgainstBestSignatureTime(certificates []*model.CertificateToken, bestSignatureTime time.Time, status *TokenStatus, context enumerations.Context) {
	hasBestSignatureTime := !bestSignatureTime.IsZero()
	for _, certificateToken := range certificates {
		if c.isSelfSignedOrTrustedAtTime(certificateToken, bestSignatureTime) {
			// break on the first trusted entry
			break
		} else if c.isRevocationDataNotRequired(certificateToken, bestSignatureTime) {
			// skip the revocation check for OCSP certs if no check is specified
			continue
		}

		found := false
		var earliestNextUpdate time.Time

		relatedRevocationTokens := c.getRelatedRevocationTokens(certificateToken)
		for _, revocationToken := range relatedRevocationTokens {
			if !hasBestSignatureTime || c.isRevocationFresh(revocationToken, bestSignatureTime, context) {
				found = true
				break
			} else {
				nextUpdate := revocationToken.NextUpdate()
				if !nextUpdate.IsZero() && (earliestNextUpdate.IsZero() || earliestNextUpdate.After(nextUpdate)) && c.currentTime.Before(nextUpdate) {
					earliestNextUpdate = nextUpdate
				}
			}
		}

		if !found {
			if !c.getCertificateVerifier().IsCheckRevocationForUntrustedChains() && !c.containsTrustAnchorAtTime(signatureValidationContextTokensFromCertificates(certificates), bestSignatureTime) {
				status.AddRelatedTokenAndErrorMessage(certificateToken, "Revocation data is skipped for untrusted certificate chain!")
			} else if utils.IsCollectionEmpty(relatedRevocationTokens) || !hasBestSignatureTime {
				// simple revocation presence check
				status.AddRelatedTokenAndErrorMessage(certificateToken, "No revocation data found for certificate!")
			} else if !earliestNextUpdate.IsZero() {
				status.AddRelatedTokenAndErrorMessage(certificateToken, "No revocation data found after the best signature time ["+
					spi.DSSUtilsFormatDateToRFC(bestSignatureTime)+"]! The nextUpdate available after : ["+
					spi.DSSUtilsFormatDateToRFC(earliestNextUpdate)+"]")
			} else {
				status.AddRelatedTokenAndErrorMessage(certificateToken, "No revocation data found after the best signature time ["+
					spi.DSSUtilsFormatDateToRFC(bestSignatureTime)+"]!")
			}

			if freshnessStatus, ok := any(status).(*RevocationFreshnessStatus); ok {
				if utils.IsCollectionNotEmpty(relatedRevocationTokens) && signatureValidationContextNoNextUpdateDefined(relatedRevocationTokens) && earliestNextUpdate.IsZero() {
					// Define next update based on Timestamp time, when no NextUpdate is defined
					lowestPOETime := c.getLowestPOETimeForToken(certificateToken)
					if !lowestPOETime.IsZero() {
						earliestNextUpdate = lowestPOETime.Add(time.Second) // last usage + 1s
					}
				}
				if !earliestNextUpdate.IsZero() {
					freshnessStatus.AddTokenAndRevocationNextUpdateTime(certificateToken, earliestNextUpdate)
				}
			}
		}
	}
}

func signatureValidationContextNoNextUpdateDefined(relatedRevocationTokens []AnyRevocationToken) bool {
	for _, revocationToken := range relatedRevocationTokens {
		if !revocationToken.NextUpdate().IsZero() {
			return false
		}
	}
	return true
}

// CheckAllPOECoveredByRevocationData reports whether every POE is covered by fresh revocation
// data. Port of checkAllPOECoveredByRevocationData().
func (c *SignatureValidationContext) CheckAllPOECoveredByRevocationData() bool {
	status, _ := c.allPOECoveredByRevocationData()
	return status.IsEmpty()
}

// allPOECoveredByRevocationData returns the status of the POE-covered-by-revocation-data check.
func (c *SignatureValidationContext) allPOECoveredByRevocationData() (*RevocationFreshnessStatus, error) {
	status := NewRevocationFreshnessStatus()
	orderedCertificateChains, err := c.getOrderedCertificateChains()
	if err != nil {
		return status, err
	}
	for firstChainCertificate, chain := range orderedCertificateChains {
		lastCertUsageDate := c.getLatestTimestampUsageDate(firstChainCertificate)
		if !lastCertUsageDate.IsZero() {
			c.checkRevocationForCertificateChainAgainstBestSignatureTime(chain, lastCertUsageDate, &status.TokenStatus, enumerations.ContextTimestamp)
		}
	}
	if !status.IsEmpty() {
		status.SetMessage("Revocation data is missing for one or more POE(s).")
	}
	return status, nil
}

// CheckAllTimestampsValid reports whether every processed timestamp is valid. Port of
// checkAllTimestampsValid().
func (c *SignatureValidationContext) CheckAllTimestampsValid() bool {
	return c.allTimestampsValid().IsEmpty()
}

// allTimestampsValid returns the status of the all-timestamps-valid check.
func (c *SignatureValidationContext) allTimestampsValid() *TokenStatus {
	status := NewTokenStatus()
	for _, timestampToken := range c.processedTimestamps {
		if !timestampToken.IsValid() {
			status.AddRelatedTokenAndErrorMessage(timestampToken, "Signature is not intact!")
		}
	}
	if !status.IsEmpty() {
		status.SetMessage("Broken timestamp(s) detected.")
	}
	return status
}

// CheckCertificateNotRevoked reports whether the given certificate is not revoked. Port of
// checkCertificateNotRevoked(CertificateToken).
func (c *SignatureValidationContext) CheckCertificateNotRevoked(certificateToken *model.CertificateToken) bool {
	return c.certificateNotRevoked(certificateToken).IsEmpty()
}

// certificateNotRevoked returns the status of the certificate-not-revoked check.
func (c *SignatureValidationContext) certificateNotRevoked(certificateToken *model.CertificateToken) *TokenStatus {
	status := NewTokenStatus()
	c.checkCertificateIsNotRevokedRecursively(certificateToken, c.poeTimes[certificateToken.DSSIDAsString()], status)
	if !status.IsEmpty() {
		status.SetMessage("Revoked/Suspended certificate(s) detected.")
	}
	return status
}

// CheckAllSignatureCertificatesNotRevoked reports whether every signature's signing
// certificate is not revoked. Port of checkAllSignatureCertificatesNotRevoked().
func (c *SignatureValidationContext) CheckAllSignatureCertificatesNotRevoked() bool {
	return c.allSignatureCertificatesNotRevoked().IsEmpty()
}

// allSignatureCertificatesNotRevoked returns the status of the
// all-signature-certificates-not-revoked check.
func (c *SignatureValidationContext) allSignatureCertificatesNotRevoked() *TokenStatus {
	status := NewTokenStatus()
	for _, signature := range c.processedSignatures {
		c.checkSignatureCertificatesNotRevoked(signature, status)
	}
	if !status.IsEmpty() {
		status.SetMessage("Revoked/Suspended certificate(s) detected.")
	}
	return status
}

func (c *SignatureValidationContext) checkSignatureCertificatesNotRevoked(signature AdvancedSignature, status *TokenStatus) {
	signingCertificate := signature.SigningCertificateToken()
	if signingCertificate != nil {
		c.checkCertificateIsNotRevokedRecursively(signingCertificate, c.poeTimes[signature.ID()], status)
	}
}

func (c *SignatureValidationContext) checkCertificateIsNotRevokedRecursively(certificateToken *model.CertificateToken, poeTimeList []*signatureValidationContextPOE, status *TokenStatus) bool {
	lowestPOETime := c.getLowestPOETime(poeTimeList)
	if c.isSelfSignedOrTrustedAtTime(certificateToken, lowestPOETime) {
		return true
	} else if !c.isRevocationDataNotRequired(certificateToken, lowestPOETime) {
		relatedRevocationTokens := c.getRelatedRevocationTokens(certificateToken)
		// check only available revocation data in order to not duplicate
		// CheckAllRequiredRevocationDataPresent()
		if utils.IsCollectionNotEmpty(relatedRevocationTokens) {
			// check if there is a best-signature-time before the revocation date
			for _, revocationToken := range relatedRevocationTokens {
				if !c.getRevocationDataVerifier().CheckCertificateNotRevoked(revocationToken, lowestPOETime) {
					if status != nil {
						status.AddRelatedTokenAndErrorMessage(certificateToken, "Certificate is revoked/suspended!")
					}
					return false
				}
			}
		}
	}

	issuer := c.getIssuer(certificateToken)
	if issuer != nil {
		return c.checkCertificateIsNotRevokedRecursively(issuer, poeTimeList, status)
	}
	return true
}

func (c *SignatureValidationContext) isRevocationDataNotRequired(certToken *model.CertificateToken, controlTime time.Time) bool {
	return c.getRevocationDataVerifier().IsRevocationDataSkip(certToken, controlTime)
}

func (c *SignatureValidationContext) isSelfSignedOrTrustedAtTime(certToken *model.CertificateToken, controlTime time.Time) bool {
	return c.isSelfSigned(certToken) || c.isTrustedAtTime(certToken, controlTime, "")
}

func (c *SignatureValidationContext) isSelfSigned(certToken model.Token) bool {
	return certToken.IsSelfSigned()
}

func (c *SignatureValidationContext) getRelatedRevocationTokens(certificateToken *model.CertificateToken) []AnyRevocationToken {
	return c.GetRevocationData(certificateToken)
}

func (c *SignatureValidationContext) isRevocationDataRefreshNeeded(certToken *model.CertificateToken, revocations []AnyRevocationToken) bool {
	var context enumerations.Context
	// get best-signature-time for b-level certificate chain
	refreshNeededAfterTime := c.getLatestBestSignatureTime(certToken)
	if !refreshNeededAfterTime.IsZero() {
		context = enumerations.ContextSignature
	}
	// get last usage dates for the same timestamp certificate chain
	lastTimestampUsageTime := c.getLatestTimestampUsageDate(certToken)
	if !lastTimestampUsageTime.IsZero() {
		if context == "" {
			context = enumerations.ContextTimestamp
		}
	}
	// return best POE for other cases
	if refreshNeededAfterTime.IsZero() {
		// shall not return zero
		refreshNeededAfterTime = c.getLowestPOETimeForToken(certToken)
		if context == "" {
			context = enumerations.ContextRevocation
		}
	}
	freshRevocationDataFound := false
	for _, revocationToken := range revocations {
		certificateTokenChain := signatureValidationContextToCertificateTokenChain(c.getCertChain(revocationToken))
		if utils.IsCollectionEmpty(certificateTokenChain) {
			continue
		}

		issuerCertificateToken := certificateTokenChain[0]
		if c.isRevocationFresh(revocationToken, refreshNeededAfterTime, context) &&
			c.isRevocationIssuedAfterLastTimestampUsage(revocationToken, lastTimestampUsageTime, context) &&
			enumerations.RevocationReasonCertificateHold != revocationToken.Reason() &&
			c.isRevocationAcceptable(revocationToken, issuerCertificateToken, c.getLowestPOETimeForToken(issuerCertificateToken)) &&
			c.hasValidPOE(revocationToken, certToken, issuerCertificateToken) {
			freshRevocationDataFound = true
			break
		}
	}

	return !freshRevocationDataFound
}

func (c *SignatureValidationContext) getLatestBestSignatureTime(certificateToken *model.CertificateToken) time.Time {
	var latestPOETime time.Time
	for _, bestSignatureTime := range c.getBestSignatureTimes(certificateToken) {
		if latestPOETime.IsZero() || bestSignatureTime.After(latestPOETime) {
			latestPOETime = bestSignatureTime
		}
	}
	return latestPOETime
}

func (c *SignatureValidationContext) getBestSignatureTimes(certificateToken *model.CertificateToken) []time.Time {
	signatures := c.getSignaturesIssuedByCertificateOrItsChildren(certificateToken, nil)
	if utils.IsCollectionEmpty(signatures) {
		return nil
	}
	var bestSignatureTimes []time.Time
	for _, signature := range signatures {
		poeTime := c.getLowestPOETime(c.poeTimes[signature.ID()])
		if !poeTime.IsZero() {
			if !signatureValidationContextContainsTime(bestSignatureTimes, poeTime) {
				bestSignatureTimes = append(bestSignatureTimes, poeTime)
			}
		}
	}
	return bestSignatureTimes
}

func (c *SignatureValidationContext) getSignaturesIssuedByCertificateOrItsChildren(certificateToken *model.CertificateToken, processedCertificates []*model.CertificateToken) []AdvancedSignature {
	var signatures []AdvancedSignature
	entry := signatureValidationContextFindCertSignatures(c.signCertificateSignaturesUsage, certificateToken)
	if entry != nil {
		for _, signature := range entry.signatures {
			if !signatureValidationContextContainsSignature(signatures, signature) {
				signatures = append(signatures, signature)
			}
		}
	}
	processedCertificates = append(processedCertificates, certificateToken)

	childrenEntry := signatureValidationContextFindCertChildren(c.certificateChildrenMap, certificateToken)
	if childrenEntry != nil {
		for _, certKid := range childrenEntry.children {
			if !signatureValidationContextContainsCertificate(processedCertificates, certKid) {
				for _, signature := range c.getSignaturesIssuedByCertificateOrItsChildren(certKid, processedCertificates) {
					if !signatureValidationContextContainsSignature(signatures, signature) {
						signatures = append(signatures, signature)
					}
				}
			}
		}
	}
	return signatures
}

func (c *SignatureValidationContext) getLatestTimestampUsageDate(certificateToken *model.CertificateToken) time.Time {
	entry := signatureValidationContextFindCertDates(c.timestampCertChainDates, certificateToken)
	if entry == nil {
		return time.Time{}
	}
	return signatureValidationContextGetLatestTime(entry.dates)
}

func signatureValidationContextGetLatestTime(dates []time.Time) time.Time {
	var latestTime time.Time
	for _, date := range dates {
		if latestTime.IsZero() || date.After(latestTime) {
			latestTime = date
		}
	}
	return latestTime
}

func (c *SignatureValidationContext) getLowestPOETimeForToken(token model.Token) time.Time {
	return c.getLowestPOETime(c.poeTimes[token.DSSIDAsString()])
}

func (c *SignatureValidationContext) getLowestPOETime(poeList []*signatureValidationContextPOE) time.Time {
	return c.getLowestPOE(poeList).time
}

func (c *SignatureValidationContext) getLowestPOE(poeList []*signatureValidationContextPOE) *signatureValidationContextPOE {
	if utils.IsCollectionEmpty(poeList) {
		panic("POE shall be defined before accessing the 'poeTimes' list!")
	}
	lowestPOE := poeList[0]
	for _, poe := range poeList[1:] {
		if poe.time.Before(lowestPOE.time) {
			lowestPOE = poe
		}
	}
	return lowestPOE
}

func (c *SignatureValidationContext) isRevocationFresh(revocationToken AnyRevocationToken, refreshNeededAfterTime time.Time, context enumerations.Context) bool {
	return c.getRevocationDataVerifier().IsRevocationDataFresh(revocationToken, refreshNeededAfterTime, context)
}

func (c *SignatureValidationContext) isRevocationIssuedAfterLastTimestampUsage(revocationToken AnyRevocationToken, lastTimestampUsage time.Time, context enumerations.Context) bool {
	if lastTimestampUsage.IsZero() {
		return true
	}
	return c.getRevocationDataVerifier().IsRevocationDataFresh(revocationToken, lastTimestampUsage, context)
}

func (c *SignatureValidationContext) isRevocationAcceptable(revocation AnyRevocationToken, issuerCertificateToken *model.CertificateToken, controlTime time.Time) bool {
	return c.getRevocationDataVerifier().IsAcceptableForChain(revocation, issuerCertificateToken, c.getCertificateTokenChain(issuerCertificateToken), controlTime)
}

func (c *SignatureValidationContext) hasValidPOE(revocationToken AnyRevocationToken, relatedCertToken, issuerCertToken *model.CertificateToken) bool {
	if !revocationToken.NextUpdate().IsZero() && !c.hasPOEAfterThisUpdateAndBeforeNextUpdate(revocationToken) {
		return false
	}
	// useful for short-life certificates (i.e. ocsp responder)
	if issuerCertToken != nil && !c.isTrustedAtUsageTimeWithContext(issuerCertToken, enumerations.ContextRevocation) && !c.hasPOEInTheValidityRange(issuerCertToken) {
		return false
	}
	return true
}

func (c *SignatureValidationContext) hasPOEAfterThisUpdateAndBeforeNextUpdate(revocation AnyRevocationToken) bool {
	poeTimeList := c.poeTimes[revocation.DSSIDAsString()]
	if utils.IsCollectionNotEmpty(poeTimeList) {
		for _, poeTime := range poeTimeList {
			if c.getRevocationDataVerifier().IsAfterThisUpdateAndBeforeNextUpdate(revocation, poeTime.time) {
				return true
			}
		}
	}
	return false
}

func (c *SignatureValidationContext) hasPOEInTheValidityRange(certificateToken *model.CertificateToken) bool {
	poeTimeList := c.poeTimes[certificateToken.DSSIDAsString()]
	if utils.IsCollectionNotEmpty(poeTimeList) {
		for _, poeTime := range poeTimeList {
			if certificateToken.IsValidOn(poeTime.time) {
				return true
			}
		}
	}
	return false
}

// CheckAllSignatureCertificateHaveFreshRevocationData reports whether every signature's
// certificate chain has fresh revocation data. Port of
// checkAllSignatureCertificateHaveFreshRevocationData().
func (c *SignatureValidationContext) CheckAllSignatureCertificateHaveFreshRevocationData() bool {
	return c.allSignatureCertificateHaveFreshRevocationData().IsEmpty()
}

// allSignatureCertificateHaveFreshRevocationData returns the status of the
// all-signature-certificate-have-fresh-revocation-data check.
func (c *SignatureValidationContext) allSignatureCertificateHaveFreshRevocationData() *RevocationFreshnessStatus {
	status := NewRevocationFreshnessStatus()
	for _, signature := range c.processedSignatures {
		c.checkAtLeastOneRevocationDataPresentAfterBestSignatureTime(signature, status)
	}
	if !status.IsEmpty() {
		status.SetMessage("Fresh revocation data is missing for one or more certificate(s).")
	}
	return status
}

func (c *SignatureValidationContext) checkAtLeastOneRevocationDataPresentAfterBestSignatureTime(signature AdvancedSignature, status *RevocationFreshnessStatus) {
	signingCertificateToken := signature.SigningCertificateToken()
	orderedCertificateChains, err := c.getOrderedCertificateChains()
	if err != nil {
		return
	}
	for firstChainCertificate, chain := range orderedCertificateChains {
		if signingCertificateToken != nil && firstChainCertificate.Equals(signingCertificateToken) {
			bestSignatureTime := c.getEarliestTimestampTime()
			c.checkRevocationForCertificateChainAgainstBestSignatureTime(chain, bestSignatureTime, &status.TokenStatus, enumerations.ContextSignature)
		}
	}
}

func (c *SignatureValidationContext) getEarliestTimestampTime() time.Time {
	var earliestDate time.Time
	for _, timestamp := range c.GetProcessedTimestamps() {
		if timestamp.TimeStampType().CoversSignature() {
			timestampTime := timestamp.CreationDate()
			if earliestDate.IsZero() || timestampTime.Before(earliestDate) {
				earliestDate = timestampTime
			}
		}
	}
	return earliestDate
}

// CheckAllSignaturesNotExpired reports whether every processed signature's signing certificate
// has not expired without a POE. Port of checkAllSignaturesNotExpired().
func (c *SignatureValidationContext) CheckAllSignaturesNotExpired() bool {
	return c.allSignaturesNotExpired().IsEmpty()
}

// allSignaturesNotExpired returns the status of the all-signatures-not-expired check.
func (c *SignatureValidationContext) allSignaturesNotExpired() *SignatureStatus {
	status := NewSignatureStatus()
	for _, signature := range c.processedSignatures {
		c.checkSignatureNotExpired(signature, status)
	}
	if !status.IsEmpty() {
		status.SetMessage("Expired signature found.")
	}
	return status
}

// CheckCertificateNotExpired reports whether the given certificate has not expired without a
// POE. Port of checkCertificateNotExpired(CertificateToken).
func (c *SignatureValidationContext) CheckCertificateNotExpired(certificateToken *model.CertificateToken) bool {
	return c.certificateNotExpired(certificateToken).IsEmpty()
}

// certificateNotExpired returns the status of the certificate-not-expired check.
func (c *SignatureValidationContext) certificateNotExpired(certificateToken *model.CertificateToken) *TokenStatus {
	status := NewTokenStatus()
	c.checkCertificateNotExpired(certificateToken, status)
	if !status.IsEmpty() {
		status.SetMessage("Expired certificate found.")
	}
	return status
}

func (c *SignatureValidationContext) checkSignatureNotExpired(signature AdvancedSignature, status *SignatureStatus) {
	signingCertificate := signature.SigningCertificateToken()
	if signingCertificate != nil {
		signatureNotExpired := c.verifyCertificateTokenNotExpired(signingCertificate, c.poeTimes[signature.ID()])
		if !signatureNotExpired {
			status.AddRelatedTokenAndErrorMessage(signature, "The signing certificate has expired and there is no POE during its validity range : ["+
				spi.DSSUtilsFormatDateToRFC(signingCertificate.NotBefore())+" - "+spi.DSSUtilsFormatDateToRFC(signingCertificate.NotAfter())+"]!")
		}
	}
}

func (c *SignatureValidationContext) checkCertificateNotExpired(certificateToken *model.CertificateToken, status *TokenStatus) {
	certificateNotExpired := c.verifyCertificateTokenNotExpired(certificateToken, c.poeTimes[certificateToken.DSSIDAsString()])
	if !certificateNotExpired {
		status.AddRelatedTokenAndErrorMessage(certificateToken, "The signing certificate has expired and there is no POE during its validity range : ["+
			spi.DSSUtilsFormatDateToRFC(certificateToken.NotBefore())+" - "+spi.DSSUtilsFormatDateToRFC(certificateToken.NotAfter())+"]!")
	}
}

func (c *SignatureValidationContext) verifyCertificateTokenNotExpired(certificateToken *model.CertificateToken, poeTimeList []*signatureValidationContextPOE) bool {
	if utils.IsCollectionNotEmpty(poeTimeList) && !certificateToken.NotAfter().IsZero() {
		for _, poeTime := range poeTimeList {
			if !poeTime.time.IsZero() && !poeTime.time.After(certificateToken.NotAfter()) {
				return true
			}
		}
	}
	return false
}

// CheckAllSignaturesAreYetValid reports whether every processed signature's signing
// certificate was already valid at the earliest available POE. Port of
// checkAllSignaturesAreYetValid().
func (c *SignatureValidationContext) CheckAllSignaturesAreYetValid() bool {
	return c.allSignaturesAreYetValid().IsEmpty()
}

// allSignaturesAreYetValid returns the status of the all-signatures-are-yet-valid check.
func (c *SignatureValidationContext) allSignaturesAreYetValid() *SignatureStatus {
	status := NewSignatureStatus()
	for _, signature := range c.processedSignatures {
		c.checkSignatureIsYetValid(signature, status)
	}
	if !status.IsEmpty() {
		status.SetMessage("Not yet valid signature found.")
	}
	return status
}

// CheckCertificateIsYetValid reports whether the given certificate was already valid at the
// earliest available POE. Port of checkCertificateIsYetValid(CertificateToken).
func (c *SignatureValidationContext) CheckCertificateIsYetValid(certificateToken *model.CertificateToken) bool {
	return c.certificateNotExpired(certificateToken).IsEmpty()
}

// certificateIsYetValid returns the status of the certificate-yet-valid check.
func (c *SignatureValidationContext) certificateIsYetValid(certificateToken *model.CertificateToken) *TokenStatus {
	status := NewTokenStatus()
	c.checkCertificateIsYetValid(certificateToken, status)
	if !status.IsEmpty() {
		status.SetMessage("Not yet valid certificate found.")
	}
	return status
}

func (c *SignatureValidationContext) checkSignatureIsYetValid(signature AdvancedSignature, status *SignatureStatus) {
	signingCertificate := signature.SigningCertificateToken()
	if signingCertificate != nil {
		signatureNotExpired := c.verifyCertificateTokenIsYetValid(signingCertificate, c.poeTimes[signature.ID()])
		if !signatureNotExpired {
			status.AddRelatedTokenAndErrorMessage(signature, "The signing certificate with validity range ["+
				spi.DSSUtilsFormatDateToRFC(signingCertificate.NotBefore())+" - "+spi.DSSUtilsFormatDateToRFC(signingCertificate.NotAfter())+"] is not yet valid at signing time!")
		}
	}
}

func (c *SignatureValidationContext) checkCertificateIsYetValid(certificateToken *model.CertificateToken, status *TokenStatus) {
	certificateNotExpired := c.verifyCertificateTokenIsYetValid(certificateToken, c.poeTimes[certificateToken.DSSIDAsString()])
	if !certificateNotExpired {
		status.AddRelatedTokenAndErrorMessage(certificateToken, "The signing certificate with validity range ["+
			spi.DSSUtilsFormatDateToRFC(certificateToken.NotBefore())+" - "+spi.DSSUtilsFormatDateToRFC(certificateToken.NotAfter())+"] is not yet valid at signing time!")
	}
}

func (c *SignatureValidationContext) verifyCertificateTokenIsYetValid(certificateToken *model.CertificateToken, poeTimeList []*signatureValidationContextPOE) bool {
	if utils.IsCollectionNotEmpty(poeTimeList) && !certificateToken.NotAfter().IsZero() {
		for _, poeTime := range poeTimeList {
			if !poeTime.time.IsZero() && !poeTime.time.Before(certificateToken.NotBefore()) {
				return true
			}
		}
	}
	return false
}

// GetProcessedSignatures returns the unmodifiable set of processed signatures. Port of
// getProcessedSignatures().
func (c *SignatureValidationContext) GetProcessedSignatures() []AdvancedSignature {
	return append([]AdvancedSignature(nil), c.processedSignatures...)
}

// GetProcessedCertificates returns the unmodifiable set of processed certificates. Port of
// getProcessedCertificates().
func (c *SignatureValidationContext) GetProcessedCertificates() []*model.CertificateToken {
	return append([]*model.CertificateToken(nil), c.processedCertificates...)
}

// GetProcessedRevocations returns the unmodifiable set of processed revocations. Port of
// getProcessedRevocations().
func (c *SignatureValidationContext) GetProcessedRevocations() []AnyRevocationToken {
	return append([]AnyRevocationToken(nil), c.processedRevocations...)
}

// GetProcessedTimestamps returns the unmodifiable set of processed timestamps. Port of
// getProcessedTimestamps().
func (c *SignatureValidationContext) GetProcessedTimestamps() []*TimestampToken {
	return append([]*TimestampToken(nil), c.processedTimestamps...)
}

// GetProcessedEvidenceRecords returns the unmodifiable set of processed evidence records. Port
// of getProcessedEvidenceRecords().
func (c *SignatureValidationContext) GetProcessedEvidenceRecords() []EvidenceRecord {
	return append([]EvidenceRecord(nil), c.processedEvidenceRecords...)
}

func (c *SignatureValidationContext) isTrustedAtUsageTime(token model.Token) bool {
	return c.isTrustedAtUsageTimeWithContext(token, "")
}

func (c *SignatureValidationContext) isTrustedAtUsageTimeWithContext(token model.Token, context enumerations.Context) bool {
	certificateToken, ok := token.(*model.CertificateToken)
	if !ok {
		return false
	}
	bestSignatureTimes := c.getBestSignatureTimes(certificateToken)
	if utils.IsCollectionNotEmpty(bestSignatureTimes) {
		for _, date := range bestSignatureTimes {
			if c.isTrustedAtTime(token, date, context) {
				return true
			}
		}
		return false
	}
	lowestPOETime := c.getLowestPOETimeForToken(token)
	return c.isTrustedAtTime(token, lowestPOETime, context)
}

func (c *SignatureValidationContext) isTrustedAtTime(token model.Token, controlTime time.Time, context enumerations.Context) bool {
	certificateToken, ok := token.(*model.CertificateToken)
	if !ok {
		return false
	}
	return c.getTrustAnchorVerifier().IsTrustedAtTime(certificateToken, controlTime, context)
}

// GetValidationData returns the validation data (certificates + revocation) for the given
// signature. Port of getValidationData(AdvancedSignature).
func (c *SignatureValidationContext) GetValidationData(signature AdvancedSignature) *ValidationData {
	return c.getValidationDataForCertificate(signature.SigningCertificateToken())
}

// GetValidationDataForTimestamp returns the validation data (certificates + revocation) for
// the given timestamp. Port of the getValidationData(TimestampToken) overload.
func (c *SignatureValidationContext) GetValidationDataForTimestamp(timestampToken *TimestampToken) *ValidationData {
	return c.getValidationDataForCertificate(c.getIssuer(timestampToken))
}

func (c *SignatureValidationContext) getValidationDataForCertificate(certificateToken *model.CertificateToken) *ValidationData {
	validationData := NewValidationData()
	if certificateToken != nil {
		c.populateValidationDataRecursively(certificateToken, validationData)
	}
	return validationData
}

func (c *SignatureValidationContext) populateValidationDataRecursively(token model.Token, validationData *ValidationData) {
	added := validationData.AddToken(token)
	if added {
		if certificateToken, ok := token.(*model.CertificateToken); ok {
			for _, revocationToken := range c.getRelatedRevocationTokens(certificateToken) {
				c.populateValidationDataRecursively(revocationToken, validationData)
			}
		}
		issuerToken := c.getIssuer(token)
		if issuerToken != nil {
			c.populateValidationDataRecursively(issuerToken, validationData)
		}
	}
}

// --- append-if-absent / lookup helpers over Java's value-equal Set<T>/Map<K,V> collections ---

func signatureValidationContextAppendSignature(signatures []AdvancedSignature, signature AdvancedSignature) []AdvancedSignature {
	if signatureValidationContextContainsSignature(signatures, signature) {
		return signatures
	}
	return append(signatures, signature)
}

func signatureValidationContextContainsSignature(signatures []AdvancedSignature, signature AdvancedSignature) bool {
	for _, s := range signatures {
		if s.ID() == signature.ID() {
			return true
		}
	}
	return false
}

func signatureValidationContextAppendCertificate(certificates []*model.CertificateToken, certificateToken *model.CertificateToken) []*model.CertificateToken {
	if signatureValidationContextContainsCertificate(certificates, certificateToken) {
		return certificates
	}
	return append(certificates, certificateToken)
}

func signatureValidationContextContainsCertificate(certificates []*model.CertificateToken, certificateToken *model.CertificateToken) bool {
	for _, c := range certificates {
		if c.Equals(certificateToken) {
			return true
		}
	}
	return false
}

func signatureValidationContextAppendRevocation(revocations []AnyRevocationToken, revocationToken AnyRevocationToken) []AnyRevocationToken {
	if signatureValidationContextContainsRevocation(revocations, revocationToken) {
		return revocations
	}
	return append(revocations, revocationToken)
}

// signatureValidationContextContainsRevocation is the membership test of Java's
// Set<RevocationToken<?>> processedRevocations, and therefore has to reproduce
// RevocationToken#equals - which compares the DSS Id *and* the related certificate:
//
//	if (!getDSSId().equals(other.getDSSId())) return false;
//	if (relatedCertificate == null) return other.relatedCertificate == null;
//	else return relatedCertificate.equals(other.relatedCertificate);
//
// The related-certificate half is load-bearing, not incidental: one CRL commonly covers several
// certificates of the same chain, and upstream deliberately keeps one RevocationToken per
// (revocation, certificate) pair so that getRevocationData(certificate) - which matches on
// relatedCertificateId - finds it for each of them. Comparing DSS Ids alone collapsed those into
// a single entry, so every certificate but the first silently lost its revocation data and
// checkAllRequiredRevocationDataPresent() reported an LT-incomplete signature (observed on
// CAdESDoubleLTA.p7m: the intermediate good-ca kept no CRL, which downgraded the detected level
// from CAdES-BASELINE-LTA to CAdES-BASELINE-T).
func signatureValidationContextContainsRevocation(revocations []AnyRevocationToken, revocationToken AnyRevocationToken) bool {
	for _, r := range revocations {
		if r.DSSIDAsString() == revocationToken.DSSIDAsString() &&
			r.RelatedCertificateID() == revocationToken.RelatedCertificateID() {
			return true
		}
	}
	return false
}

func signatureValidationContextAppendTimestamp(timestamps []*TimestampToken, timestampToken *TimestampToken) []*TimestampToken {
	for _, t := range timestamps {
		if t.DSSIDAsString() == timestampToken.DSSIDAsString() {
			return timestamps
		}
	}
	return append(timestamps, timestampToken)
}

func signatureValidationContextAppendEvidenceRecord(evidenceRecords []EvidenceRecord, evidenceRecord EvidenceRecord) []EvidenceRecord {
	for _, e := range evidenceRecords {
		if e == evidenceRecord {
			return evidenceRecords
		}
	}
	return append(evidenceRecords, evidenceRecord)
}

func signatureValidationContextContainsToken(tokens []model.Token, token model.Token) bool {
	for _, t := range tokens {
		if t.DSSIDAsString() == token.DSSIDAsString() {
			return true
		}
	}
	return false
}

func signatureValidationContextContainsTime(times []time.Time, t time.Time) bool {
	for _, existing := range times {
		if existing.Equal(t) {
			return true
		}
	}
	return false
}

func signatureValidationContextToCertificateTokenChain(tokens []model.Token) []*model.CertificateToken {
	var chain []*model.CertificateToken
	for _, token := range tokens {
		if certificateToken, ok := token.(*model.CertificateToken); ok {
			chain = append(chain, certificateToken)
		}
	}
	return chain
}

func signatureValidationContextTokensFromCertificates(certificates []*model.CertificateToken) []model.Token {
	tokens := make([]model.Token, len(certificates))
	for i, c := range certificates {
		tokens[i] = c
	}
	return tokens
}

func signatureValidationContextMapValues(m map[string]*model.CertificateToken) []*model.CertificateToken {
	values := make([]*model.CertificateToken, 0, len(m))
	for _, v := range m {
		values = append(values, v)
	}
	return values
}

func signatureValidationContextFindCertSignatures(entries []signatureValidationContextCertSignatures, cert *model.CertificateToken) *signatureValidationContextCertSignatures {
	for i := range entries {
		if entries[i].cert.Equals(cert) {
			return &entries[i]
		}
	}
	return nil
}

func signatureValidationContextFindCertDates(entries []signatureValidationContextCertDates, cert *model.CertificateToken) *signatureValidationContextCertDates {
	for i := range entries {
		if entries[i].cert.Equals(cert) {
			return &entries[i]
		}
	}
	return nil
}

func signatureValidationContextFindCertChildren(entries []signatureValidationContextCertChildren, parent *model.CertificateToken) *signatureValidationContextCertChildren {
	if parent == nil {
		return nil
	}
	for i := range entries {
		if entries[i].parent != nil && entries[i].parent.Equals(parent) {
			return &entries[i]
		}
	}
	return nil
}

func signatureValidationContextFindExternalRevocations(entries []signatureValidationContextExternalRevocations, issuer *model.CertificateToken) *signatureValidationContextExternalRevocations {
	for i := range entries {
		if entries[i].issuer.Equals(issuer) {
			return &entries[i]
		}
	}
	return nil
}

func signatureValidationContextTokenIssuerMapContainsKey(entries []signatureValidationContextTokenIssuer, token model.Token) bool {
	for _, entry := range entries {
		if entry.token.DSSIDAsString() == token.DSSIDAsString() {
			return true
		}
	}
	return false
}

func signatureValidationContextGetTokenIssuer(entries []signatureValidationContextTokenIssuer, token model.Token) *model.CertificateToken {
	for _, entry := range entries {
		if entry.token.DSSIDAsString() == token.DSSIDAsString() {
			return entry.issuer
		}
	}
	return nil
}

func signatureValidationContextPutTokenIssuer(entries *[]signatureValidationContextTokenIssuer, token model.Token, issuer *model.CertificateToken) {
	for i := range *entries {
		if (*entries)[i].token.DSSIDAsString() == token.DSSIDAsString() {
			(*entries)[i].issuer = issuer
			return
		}
	}
	*entries = append(*entries, signatureValidationContextTokenIssuer{token: token, issuer: issuer})
}
