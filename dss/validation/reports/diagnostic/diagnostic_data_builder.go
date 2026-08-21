// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/reports/diagnostic/DiagnosticDataBuilder.java (DSS 6.5.RC1).
//
// # Virtual dispatch
//
// Java's DiagnosticDataBuilder is abstract; SignedDocumentDiagnosticDataBuilder (this same
// manifest) overrides the protected linkSigningCertificateAndChains(Set<CertificateToken>) to a
// no-op, and DiagnosticDataBuilder.build() self-calls it - so the base's own build() needs
// virtual dispatch to reach the override when called through a SignedDocumentDiagnosticDataBuilder
// (or QWACCertificateDiagnosticDataBuilder) instance via super.build(). Every other protected/
// public method below is inherited unmodified by every subclass surveyed in this manifest and
// the forward-declared CAdES/PAdES/JAdES/ASiC/QWAC consumers, so it stays an ordinary method;
// only LinkSigningCertificateAndChains is collected into DiagnosticDataBuilderOverrides,
// following the AbstractSignatureIdentifierBuilder/DefaultDocumentAnalyzer precedent
// (Init<TypeName> registration, see PORTING.md's "Virtual dispatch" note).
//
// # Set<T> and Map<K,V> representations
//
// Java's Set<CertificateToken>/Set<RevocationToken<?>> usedCertificates/usedRevocations are the
// exact fields XmlDiagnosticDataFactory wires from the already-landed
// SignatureValidationContext.GetProcessedCertificates() []*model.CertificateToken /
// GetProcessedRevocations() []validation.AnyRevocationToken (spi/validation/validation_context.go)
// - already order-preserving slices, not Go maps - so this port keeps them as slices throughout
// (matching PORTING.md's "order-sensitive upstream iteration -> slice" rule) rather than
// resurrecting Set semantics through a map.
//
// Every Java HashMap the source uses purely as a get/put cache (xmlCertsMap, xmlRevocationsMap,
// xmlTrustedListsMap, xmlTrustSourceMap, referenceMap, certificateIdsMap, signingCertificateMap,
// tlInfoMap, loteInfoMap) stays a Go map: it is never ranged to produce output, only looked up
// by a token/identifier's DSSIdAsString(). The two exceptions - xmlOrphanCertificateTokensMap
// and xmlOrphanRevocationTokensMap, whose .values() feed buildXmlOrphanTokens()'s output lists -
// carry a companion insertion-order slice so the output is built by iterating that slice, never
// by ranging the map (S8C_BRIEF.md's "deterministic iteration only... NO map ranging into
// output" rule). The same technique is used, function-locally, everywhere else this file
// resolves a Java HashMap<Identifier,XmlTrustedList>-style local variable into an output list
// (buildXmlTrustedLists/buildXmlLoTEs and their helpers): a local ordered map is a
// (map[string]V, []string keys) pair built and drained through the keys slice.
//
// Java's HashMap iteration order (used, e.g., by trustedLists.addAll(mapTrustedLists.values()))
// is not itself insertion order - it is whatever bucket order the JVM's HashMap implementation
// assigns from Identifier.hashCode(). Reproducing that exact bucket order in Go is impractical
// (it would mean re-implementing java.util.HashMap's hashing/resizing algorithm bit for bit);
// this port instead uses first-seen insertion order, which is deterministic across Go runs for
// a given input (satisfying the brief's determinism rule) but is NOT guaranteed to byte-match
// the Java reference dump when a diagnostic data document carries more than one trusted list or
// list-of-trusted-entities reachable from the same certificate set. Flagged in the porter notes
// for the RPTDIAG byte-compare oracle to confirm against real multi-TL fixtures.
package diagnostic

import (
	"fmt"
	"sort"
	"time"

	"github.com/utain/esig/dss/diagnostic/jaxb"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/model/job"
	"github.com/utain/esig/dss/model/lote"
	"github.com/utain/esig/dss/model/tsl"
	"github.com/utain/esig/dss/model/x509/extension"
	"github.com/utain/esig/dss/model/x509/revocation"
	"github.com/utain/esig/dss/spi"
	spilote "github.com/utain/esig/dss/spi/lote"
	"github.com/utain/esig/dss/spi/validation"
	"github.com/utain/esig/dss/utils"
)

// DiagnosticDataBuilderOverrides declares the operations DiagnosticDataBuilder calls back into
// that at least one subclass in this port overrides. A concrete/intermediate builder registers
// itself through InitDiagnosticDataBuilder; every method a subclass does not override is
// supplied by the embedded DiagnosticDataBuilder through ordinary Go method promotion.
type DiagnosticDataBuilderOverrides interface {
	// LinkSigningCertificateAndChains links the certificates and their certificate chains. Port
	// of the protected linkSigningCertificateAndChains(Set<CertificateToken>).
	LinkSigningCertificateAndChains(certificates []*model.CertificateToken)
}

// DiagnosticDataBuilder contains a common code for DiagnosticData building.
type DiagnosticDataBuilder struct {
	// usedCertificates are the certificates used during the validation process. Port of the
	// protected Set<CertificateToken> usedCertificates.
	usedCertificates []*model.CertificateToken

	// usedRevocations are the revocation data used during the validation process. Port of the
	// protected Set<RevocationToken<?>> usedRevocations.
	usedRevocations []validation.AnyRevocationToken

	// allCertificateSources is the list of all certificate sources.
	allCertificateSources *spi.ListCertificateSource

	// validationDate is the validation time. The zero time.Time stands for Java's null.
	validationDate time.Time

	// tokenExtractionStrategy is the token extraction strategy.
	tokenExtractionStrategy enumerations.TokenExtractionStrategy

	// defaultDigestAlgorithm is the digest algorithm to use for digest computation.
	defaultDigestAlgorithm enumerations.DigestAlgorithm

	// identifierProvider generates ids for the tokens.
	identifierProvider model.TokenIdentifierProvider

	// xmlCertsMap is the cached map of certificates.
	xmlCertsMap map[string]*jaxb.XmlCertificate

	// xmlRevocationsMap is the cached map of revocation data.
	xmlRevocationsMap map[string]*jaxb.XmlRevocation

	// xmlTrustedListsMap is the cached map of trusted lists.
	xmlTrustedListsMap map[string]*jaxb.XmlTrustedList

	// xmlTrustSourceMap is the cached map of trusted sources.
	xmlTrustSourceMap map[string]*jaxb.XmlListOfTrustedEntities

	// xmlOrphanCertificateTokensMap is the cached map of orphan certificates.
	xmlOrphanCertificateTokensMap map[string]*jaxb.XmlOrphanCertificateToken
	// xmlOrphanCertificateTokensOrder is the first-seen insertion order of
	// xmlOrphanCertificateTokensMap's keys; see the file header's map-ranging note.
	xmlOrphanCertificateTokensOrder []string

	// xmlOrphanRevocationTokensMap is the cached map of orphan revocation data.
	xmlOrphanRevocationTokensMap map[string]*jaxb.XmlOrphanRevocationToken
	// xmlOrphanRevocationTokensOrder is the first-seen insertion order of
	// xmlOrphanRevocationTokensMap's keys; see the file header's map-ranging note.
	xmlOrphanRevocationTokensOrder []string

	// referenceMap maps reference ids to their related token ids (used to map references for
	// timestamped refs).
	referenceMap map[string]string

	// certificateIdsMap maps certificate id Strings to the related CertificateTokens.
	certificateIdsMap map[string]*model.CertificateToken

	// signingCertificateMap maps certificate id Strings to the related CertificateTokens for
	// signing certificates.
	signingCertificateMap map[string]*model.CertificateToken

	// tlInfoMap is the cached map of trusted lists with corresponding TLInfo.
	tlInfoMap map[string]*tsl.TLInfo

	// loteInfoMap is the cached map of lists of trusted entities with corresponding LoTEInfo.
	loteInfoMap map[string]*lote.LoTEInfo

	// overrides points back at the concrete/intermediate builder; see InitDiagnosticDataBuilder.
	overrides DiagnosticDataBuilderOverrides
}

// NewDiagnosticDataBuilder instantiates the object with default values. Port of the protected
// default constructor.
func NewDiagnosticDataBuilder() *DiagnosticDataBuilder {
	return &DiagnosticDataBuilder{
		allCertificateSources:         spi.NewListCertificateSource(),
		tokenExtractionStrategy:       enumerations.TokenExtractionStrategy_NONE,
		defaultDigestAlgorithm:        enumerations.DigestAlgorithm_SHA256,
		identifierProvider:            model.NewOriginalIdentifierProvider(),
		xmlCertsMap:                   map[string]*jaxb.XmlCertificate{},
		xmlRevocationsMap:             map[string]*jaxb.XmlRevocation{},
		xmlTrustedListsMap:            map[string]*jaxb.XmlTrustedList{},
		xmlTrustSourceMap:             map[string]*jaxb.XmlListOfTrustedEntities{},
		xmlOrphanCertificateTokensMap: map[string]*jaxb.XmlOrphanCertificateToken{},
		xmlOrphanRevocationTokensMap:  map[string]*jaxb.XmlOrphanRevocationToken{},
		referenceMap:                  map[string]string{},
		certificateIdsMap:             map[string]*model.CertificateToken{},
		signingCertificateMap:         map[string]*model.CertificateToken{},
		tlInfoMap:                     map[string]*tsl.TLInfo{},
		loteInfoMap:                   map[string]*lote.LoTEInfo{},
	}
}

// InitDiagnosticDataBuilder registers the concrete/intermediate builder with the base so it can
// dispatch to DiagnosticDataBuilderOverrides. Every constructor in the DiagnosticDataBuilder
// family must call this once (see PORTING.md's Init<TypeName> convention).
func (b *DiagnosticDataBuilder) InitDiagnosticDataBuilder(overrides DiagnosticDataBuilderOverrides) {
	b.overrides = overrides
}

func (b *DiagnosticDataBuilder) diagnosticDataBuilderOverrides() DiagnosticDataBuilderOverrides {
	if b.overrides == nil {
		panic("DiagnosticDataBuilder was not initialised: the concrete builder must call InitDiagnosticDataBuilder in its constructor")
	}
	return b.overrides
}

// LinkSigningCertificateAndChains links the certificates and their certificate chains. This is
// the base (default) body of DiagnosticDataBuilderOverrides.LinkSigningCertificateAndChains.
// Port of the protected linkSigningCertificateAndChains(Set<CertificateToken>).
func (b *DiagnosticDataBuilder) LinkSigningCertificateAndChains(certificates []*model.CertificateToken) {
	for _, certificateToken := range certificates {
		certificateToken = b.getProcessedCertificateToken(certificateToken)
		xmlCertificate := b.xmlCertsMap[certificateToken.DSSIDAsString()]
		if xmlCertificate.SigningCertificate == nil {
			xmlCertificate.SigningCertificate = b.getXmlSigningCertificateForToken(certificateToken)
			xmlCertificate.CertificateChain = b.certificateChainWrapper(b.GetXmlForCertificateChain(certificateToken))
			xmlCertificate.BasicSignature = b.GetXmlBasicSignature(certificateToken)
			xmlCertificate.IssuerEntityKey = b.getXmlIssuerEntityKey(certificateToken)
		}
	}
}

// UsedCertificates sets the used certificates. Port of usedCertificates(Set<CertificateToken>).
func (b *DiagnosticDataBuilder) UsedCertificates(usedCertificates []*model.CertificateToken) *DiagnosticDataBuilder {
	b.usedCertificates = usedCertificates
	return b
}

// UsedRevocations sets the used revocation data. Port of usedRevocations(Set<RevocationToken<?>>).
func (b *DiagnosticDataBuilder) UsedRevocations(usedRevocations []validation.AnyRevocationToken) *DiagnosticDataBuilder {
	b.usedRevocations = usedRevocations
	return b
}

// AllCertificateSources sets the ListCertificateSource containing all certificate sources used
// in the validator (including trusted certificate sources). Port of
// allCertificateSources(ListCertificateSource).
func (b *DiagnosticDataBuilder) AllCertificateSources(allCertificateSources *spi.ListCertificateSource) *DiagnosticDataBuilder {
	if allCertificateSources != nil && !allCertificateSources.ContainsTrustedCertSources() {
		// Port of LOG.warn(...): slf4j dropped per PORTING.md; the warning has no observable
		// effect and is not load-bearing.
	}
	b.allCertificateSources = allCertificateSources
	return b
}

// ValidationDate sets the validation date. Port of validationDate(Date).
func (b *DiagnosticDataBuilder) ValidationDate(validationDate time.Time) *DiagnosticDataBuilder {
	b.validationDate = validationDate
	return b
}

// TokenExtractionStrategy sets the TokenExtractionStrategy to follow for the token extraction.
// Port of tokenExtractionStrategy(TokenExtractionStrategy).
func (b *DiagnosticDataBuilder) TokenExtractionStrategy(tokenExtractionStrategy enumerations.TokenExtractionStrategy) *DiagnosticDataBuilder {
	b.tokenExtractionStrategy = tokenExtractionStrategy
	return b
}

// TokenIdentifierProvider sets the TokenIdentifierProvider for identifiers generation. Port of
// tokenIdentifierProvider(TokenIdentifierProvider).
func (b *DiagnosticDataBuilder) TokenIdentifierProvider(identifierProvider model.TokenIdentifierProvider) *DiagnosticDataBuilder {
	b.identifierProvider = identifierProvider
	return b
}

// GetTokenExtractionStrategy returns the TokenExtractionStrategy set via TokenExtractionStrategy.
// Cross-package accessor added during phase 8f un-gating for
// ASiCWithCAdESDiagnosticDataBuilder.buildDetachedXmlSignature() (out of this manifest), which
// needs to propagate this field into a freshly-built nested CAdESDiagnosticDataBuilder the way
// Java reads the protected tokenExtractionStrategy field directly - see
// SignedDocumentDiagnosticDataBuilder.GetDocumentCertificateSource's doc comment for the same
// cross-package-getter rationale. Purely additive; does not change TokenExtractionStrategy's
// existing fluent-setter behavior.
func (b *DiagnosticDataBuilder) GetTokenExtractionStrategy() enumerations.TokenExtractionStrategy {
	return b.tokenExtractionStrategy
}

// GetTokenIdentifierProvider returns the TokenIdentifierProvider set via TokenIdentifierProvider.
// See GetTokenExtractionStrategy's doc comment.
func (b *DiagnosticDataBuilder) GetTokenIdentifierProvider() model.TokenIdentifierProvider {
	return b.identifierProvider
}

// DefaultDigestAlgorithm sets the default DigestAlgorithm which will be used for tokens'
// DigestAlgoAndValue calculation. Port of defaultDigestAlgorithm(DigestAlgorithm).
func (b *DiagnosticDataBuilder) DefaultDigestAlgorithm(digestAlgorithm enumerations.DigestAlgorithm) *DiagnosticDataBuilder {
	b.defaultDigestAlgorithm = digestAlgorithm
	return b
}

// Build builds the XmlDiagnosticData. Port of build().
func (b *DiagnosticDataBuilder) Build() *jaxb.XmlDiagnosticData {
	diagnosticData := &jaxb.XmlDiagnosticData{}
	if !b.validationDate.IsZero() {
		diagnosticData.ValidationDate = jaxb.NewXSDateTime(b.validationDate)
	}

	xmlCertificates := b.buildXmlCertificates(b.usedCertificates)
	diagnosticData.UsedCertificates = &jaxb.UsedCertificatesWrapper{Items: xmlCertificates}

	xmlRevocations := b.buildXmlRevocations(b.usedRevocations)
	diagnosticData.UsedRevocations = &jaxb.UsedRevocationsWrapper{Items: xmlRevocations}

	b.diagnosticDataBuilderOverrides().LinkSigningCertificateAndChains(b.usedCertificates)
	b.linkCertificatesAndRevocations(b.usedCertificates)

	if b.isUseLoTEs() {
		lotes := b.buildXmlLoTEs(b.allCertificateSources)
		diagnosticData.ListsOfTrustedEntities = &jaxb.ListsOfTrustedEntitiesWrapper{Items: lotes}
		b.linkCertificatesAndTrustedEntities(b.usedCertificates)
	}
	if b.isUseTrustedLists() {
		trustedLists := b.buildXmlTrustedLists(b.allCertificateSources)
		diagnosticData.TrustedLists = &jaxb.TrustedListsWrapper{Items: trustedLists}
		b.linkCertificatesAndTrustServices(b.usedCertificates)
	}
	return diagnosticData
}

func (b *DiagnosticDataBuilder) isUseTrustedLists() bool {
	if !b.allCertificateSources.IsEmpty() {
		for _, certificateSource := range b.allCertificateSources.Sources() {
			if _, ok := certificateSource.(tsl.TrustPropertiesCertificateSource); ok {
				return true
			}
		}
	}
	return false
}

func (b *DiagnosticDataBuilder) isUseLoTEs() bool {
	if !b.allCertificateSources.IsEmpty() {
		for _, certificateSource := range b.allCertificateSources.Sources() {
			if _, ok := certificateSource.(*spilote.TrustedEntitiesCertificateSource); ok {
				return true
			}
		}
	}
	return false
}

func (b *DiagnosticDataBuilder) buildXmlCertificates(certificates []*model.CertificateToken) []*jaxb.XmlCertificate {
	builtCertificates := make([]*jaxb.XmlCertificate, 0)
	if utils.IsCollectionNotEmpty(certificates) {
		tokens := append([]*model.CertificateToken{}, certificates...)
		comparator := model.NewTokenComparator()
		sort.SliceStable(tokens, func(i, j int) bool {
			return comparator.Less(tokens[i], tokens[j])
		})
		for _, certificateToken := range tokens {
			id := certificateToken.DSSIDAsString()
			xmlCertificate, ok := b.xmlCertsMap[id]
			if !ok {
				xmlCertificate = b.BuildDetachedXmlCertificate(certificateToken)
				b.xmlCertsMap[id] = xmlCertificate
			}
			b.certificateIdsMap[certificateToken.DSSIDAsString()] = certificateToken
			builtCertificates = append(builtCertificates, xmlCertificate)
		}
	}
	return builtCertificates
}

func (b *DiagnosticDataBuilder) linkCertificatesAndTrustServices(certificates []*model.CertificateToken) {
	if utils.IsCollectionNotEmpty(certificates) {
		for _, certificateToken := range certificates {
			trustServiceProviders := NewXmlTrustServiceProviderBuilder(b.xmlCertsMap, b.xmlTrustedListsMap, b.tlInfoMap).
				Build(certificateToken, b.getRelatedTrustServices(certificateToken))
			if utils.IsCollectionNotEmpty(trustServiceProviders) {
				xmlCertificate := b.xmlCertsMap[certificateToken.DSSIDAsString()]
				xmlCertificate.TrustServiceProviders = &jaxb.TrustServiceProvidersWrapper{Items: trustServiceProviders}
			}
		}
	}
}

func (b *DiagnosticDataBuilder) getRelatedTrustServices(certToken *model.CertificateToken) map[*model.CertificateToken][]*tsl.TrustProperties {
	result := make(map[*model.CertificateToken][]*tsl.TrustProperties)
	for _, trustedSource := range b.allCertificateSources.Sources() {
		trustedCertSource, ok := trustedSource.(tsl.TrustPropertiesCertificateSource)
		if !ok {
			continue
		}
		processedTokens := make(map[string]bool)
		currentCertificate := certToken
		for currentCertificate != nil {
			trustServices := trustedCertSource.TrustServices(currentCertificate)
			if utils.IsCollectionNotEmpty(trustServices) {
				result[currentCertificate] = append(result[currentCertificate], trustServices...)
			}
			if currentCertificate.IsSelfSigned() || processedTokens[currentCertificate.DSSIDAsString()] {
				break
			}
			processedTokens[currentCertificate.DSSIDAsString()] = true

			issuerCertificate := b.getIssuerCertificate(currentCertificate)
			if issuerCertificate != nil && currentCertificate.IssuerEntityKey() != nil &&
				currentCertificate.IssuerEntityKey().Equals(issuerCertificate.EntityKey()) {
				currentCertificate = issuerCertificate
			} else {
				// avoid TrustProperties extraction for a not matching chain
				currentCertificate = nil
			}
		}
	}
	return result
}

func (b *DiagnosticDataBuilder) buildXmlRevocations(revocations []validation.AnyRevocationToken) []*jaxb.XmlRevocation {
	builtRevocations := make([]*jaxb.XmlRevocation, 0)
	if utils.IsCollectionNotEmpty(revocations) {
		tokens := append([]validation.AnyRevocationToken{}, revocations...)
		comparator := model.NewTokenComparator()
		sort.SliceStable(tokens, func(i, j int) bool {
			return comparator.Less(tokens[i], tokens[j])
		})
		uniqueIds := make(map[string]bool)
		for _, revocationToken := range tokens {
			id := revocationToken.DSSIDAsString()
			if uniqueIds[id] {
				continue
			}
			xmlRevocation := b.xmlRevocationsMap[id]
			if xmlRevocation == nil {
				xmlRevocation = b.BuildDetachedXmlRevocation(revocationToken)
				b.xmlRevocationsMap[id] = xmlRevocation
				builtRevocations = append(builtRevocations, xmlRevocation)
			}
			uniqueIds[id] = true
		}
	}
	return builtRevocations
}

func (b *DiagnosticDataBuilder) linkCertificatesAndRevocations(certificates []*model.CertificateToken) {
	if utils.IsCollectionNotEmpty(certificates) {
		for _, certificateToken := range certificates {
			xmlCertificate := b.xmlCertsMap[certificateToken.DSSIDAsString()]
			revocationsForCert := b.getRevocationsForCert(certificateToken)
			for _, revocationToken := range revocationsForCert {
				xmlRevocation := b.xmlRevocationsMap[revocationToken.DSSIDAsString()]
				xmlCertificateRevocation := &jaxb.XmlCertificateRevocation{}
				xmlCertificateRevocation.Revocation = xmlRevocation
				status := revocationToken.Status()
				xmlCertificateRevocation.Status = (*jaxb.CertificateStatusValue)(&status)
				if reason := revocationToken.Reason(); reason != "" {
					v := jaxb.RevocationReasonValue(reason)
					xmlCertificateRevocation.Reason = &v
				}
				if !revocationToken.RevocationDate().IsZero() {
					xmlCertificateRevocation.RevocationDate = jaxb.NewXSDateTime(revocationToken.RevocationDate())
				}
				xmlCertificate.Revocations = &jaxb.RevocationsWrapper{
					Items: append(xmlCertificate.Revocations.All(), xmlCertificateRevocation),
				}
			}
		}
	}
}

// BuildXmlOrphanTokens builds a list of XmlOrphanTokens. Port of buildXmlOrphanTokens().
func (b *DiagnosticDataBuilder) BuildXmlOrphanTokens() *jaxb.XmlOrphanTokens {
	if utils.IsMapNotEmpty(b.xmlOrphanCertificateTokensMap) || utils.IsMapNotEmpty(b.xmlOrphanRevocationTokensMap) {
		xmlOrphanTokens := &jaxb.XmlOrphanTokens{}
		for _, key := range b.xmlOrphanCertificateTokensOrder {
			xmlOrphanTokens.OrphanCertificate = append(xmlOrphanTokens.OrphanCertificate, b.xmlOrphanCertificateTokensMap[key])
		}
		for _, key := range b.xmlOrphanRevocationTokensOrder {
			xmlOrphanTokens.OrphanRevocation = append(xmlOrphanTokens.OrphanRevocation, b.xmlOrphanRevocationTokensMap[key])
		}
		return xmlOrphanTokens
	}
	return nil
}

func (b *DiagnosticDataBuilder) buildXmlTrustedLists(trustedCertificateSources *spi.ListCertificateSource) []*jaxb.XmlTrustedList {
	trustedLists := make([]*jaxb.XmlTrustedList, 0)

	mapTrustedLists := map[string]*jaxb.XmlTrustedList{}
	var mapTrustedListsOrder []string
	mapListOfTrustedLists := map[string]*jaxb.XmlTrustedList{}
	var mapListOfTrustedListsOrder []string

	for _, certificateSource := range trustedCertificateSources.Sources() {
		tlCertSource, ok := certificateSource.(tsl.TrustPropertiesCertificateSource)
		if !ok {
			continue
		}
		summary := tlCertSource.Summary()
		if summary != nil {
			b.mergeTrustedListsMap(mapTrustedLists, &mapTrustedListsOrder, b.getTrustedListsMap(tlCertSource, summary))
			b.mergeTrustedListsMap(mapListOfTrustedLists, &mapListOfTrustedListsOrder, b.getListOfTrustedListsMap(tlCertSource, summary))
		}
		// else: Port of LOG.warn("...TLValidationJob is not performed!"): slf4j dropped.
	}

	for _, key := range mapTrustedListsOrder {
		trustedLists = append(trustedLists, mapTrustedLists[key])
	}
	for _, key := range mapListOfTrustedListsOrder {
		trustedLists = append(trustedLists, mapListOfTrustedLists[key])
	}
	return trustedLists
}

func (b *DiagnosticDataBuilder) mergeTrustedListsMap(dst map[string]*jaxb.XmlTrustedList, order *[]string, src map[string]*jaxb.XmlTrustedList) {
	for key, value := range src {
		if _, exists := dst[key]; !exists {
			dst[key] = value
			*order = append(*order, key)
		}
	}
}

func (b *DiagnosticDataBuilder) getTrustedListsMap(tlCertSource tsl.TrustPropertiesCertificateSource,
	summary *tsl.TLValidationJobSummary) map[string]*jaxb.XmlTrustedList {
	mapTrustedLists := map[string]*jaxb.XmlTrustedList{}
	for _, tlID := range b.getTLIdentifiers(tlCertSource) {
		if _, exists := mapTrustedLists[tlID.AsXmlID()]; exists {
			continue
		}
		tlInfoByID := summary.TLInfoByID(tlID)
		if tlInfoByID != nil {
			mapTrustedLists[tlID.AsXmlID()] = b.getXmlTrustSourceListForTL(tlInfoByID)
		}
	}
	return mapTrustedLists
}

func (b *DiagnosticDataBuilder) getTLIdentifiers(tlCS tsl.TrustPropertiesCertificateSource) []model.Identifier {
	var tlIdentifiers []model.Identifier
	seen := map[string]bool{}
	for _, certificateToken := range b.usedCertificates {
		trustServices := tlCS.TrustServices(certificateToken)
		for _, trustProperties := range trustServices {
			tlInfo := trustProperties.TLInfo()
			if tlInfo != nil {
				id := tlInfo.DSSID()
				if !seen[id.AsXmlID()] {
					seen[id.AsXmlID()] = true
					tlIdentifiers = append(tlIdentifiers, id)
				}
			}
		}
	}
	return tlIdentifiers
}

func (b *DiagnosticDataBuilder) getListOfTrustedListsMap(tlCertSource tsl.TrustPropertiesCertificateSource,
	summary *tsl.TLValidationJobSummary) map[string]*jaxb.XmlTrustedList {
	mapListOfTrustedLists := map[string]*jaxb.XmlTrustedList{}
	for _, lotlID := range b.getLOTLIdentifiers(tlCertSource) {
		if _, exists := mapListOfTrustedLists[lotlID.AsXmlID()]; exists {
			continue
		}
		lotlInfoByID := summary.LOTLInfoByID(lotlID)
		if lotlInfoByID != nil {
			mapListOfTrustedLists[lotlID.AsXmlID()] = b.getXmlTrustSourceListForLOTL(lotlInfoByID)
		}
	}
	return mapListOfTrustedLists
}

func (b *DiagnosticDataBuilder) getLOTLIdentifiers(tlCS tsl.TrustPropertiesCertificateSource) []model.Identifier {
	var lotlIdentifiers []model.Identifier
	seen := map[string]bool{}
	for _, certificateToken := range b.usedCertificates {
		trustServices := tlCS.TrustServices(certificateToken)
		for _, trustProperties := range trustServices {
			lotlInfo := trustProperties.LOTLInfo()
			if lotlInfo != nil {
				id := lotlInfo.DSSID()
				if !seen[id.AsXmlID()] {
					seen[id.AsXmlID()] = true
					lotlIdentifiers = append(lotlIdentifiers, id)
				}
			}
		}
	}
	return lotlIdentifiers
}

// getXmlTrustSourceListForTL builds/looks up the XmlTrustedList for a TLInfo (Java's overload
// getXmlTrustSourceList(TLInfo)).
func (b *DiagnosticDataBuilder) getXmlTrustSourceListForTL(tlInfo *tsl.TLInfo) *jaxb.XmlTrustedList {
	id := tlInfo.DSSIDAsString()
	result, ok := b.xmlTrustedListsMap[id]
	if !ok {
		result = &jaxb.XmlTrustedList{}
		b.xmlTrustedListsMap[id] = result
	}
	// A TLInfo (not a LOTLInfo) never sets LOTL; matches Java's instanceof check.
	idStr := b.identifierProvider.IDAsString(tlInfo)
	result.Id = jaxb.NewCollapsedString(idStr)
	url := tlInfo.Url()
	result.Url = &url
	if tlInfo.Parent() != nil {
		result.Parent = b.getXmlTrustSourceListForLOTL(tlInfo.Parent())
	}
	parsingCacheInfo, parsingCacheInfoOK := tlInfo.TLParsingCacheInfo()
	if parsingCacheInfoOK {
		if parsingCacheInfo.TSLType() != nil {
			tp := parsingCacheInfo.TSLType().URI()
			result.Type = &tp
		}
		territory := parsingCacheInfo.Territory()
		result.CountryCode = &territory
		if !parsingCacheInfo.IssueDate().IsZero() {
			result.IssueDate = jaxb.NewXSDateTime(parsingCacheInfo.IssueDate())
		}
		if !parsingCacheInfo.NextUpdateDate().IsZero() {
			result.NextUpdate = jaxb.NewXSDateTime(parsingCacheInfo.NextUpdateDate())
		}
		if parsingCacheInfo.SequenceNumber() != nil {
			result.SequenceNumber = parsingCacheInfo.SequenceNumber()
		}
		if parsingCacheInfo.Version() != nil {
			result.Version = parsingCacheInfo.Version()
		}
		result.StructuralValidation = b.GetXmlStructuralValidation(parsingCacheInfo.StructureValidationMessages())
	}
	downloadCacheInfo := tlInfo.DownloadCacheInfo()
	if downloadCacheInfo != nil && !downloadCacheInfo.LastSuccessSynchronizationTime().IsZero() {
		result.LastLoading = jaxb.NewXSDateTime(downloadCacheInfo.LastSuccessSynchronizationTime())
	}
	validationCacheInfo := tlInfo.ValidationCacheInfo()
	if validationCacheInfo != nil {
		wellSigned := validationCacheInfo.IsValid()
		result.WellSigned = wellSigned
	}
	if tlInfo.OtherTSLPointer() != nil && tlInfo.OtherTSLPointer().Mra() != nil {
		mra := true
		result.Mra = &mra
	}
	b.tlInfoMap[id] = tlInfo
	return result
}

// getXmlTrustSourceListForLOTL builds/looks up the XmlTrustedList for a LOTLInfo (Java's
// getXmlTrustSourceList(TLInfo) called with a LOTLInfo, which is-a TLInfo in Java; Go has no
// upcast, so this is a distinct helper sharing the same body against LOTLInfo's TLInfo-shaped
// accessors).
func (b *DiagnosticDataBuilder) getXmlTrustSourceListForLOTL(lotlInfo *tsl.LOTLInfo) *jaxb.XmlTrustedList {
	id := lotlInfo.DSSIDAsString()
	result, ok := b.xmlTrustedListsMap[id]
	if !ok {
		result = &jaxb.XmlTrustedList{}
		b.xmlTrustedListsMap[id] = result
	}
	lotl := true
	result.LOTL = &lotl
	idStr := b.identifierProvider.IDAsString(lotlInfo)
	result.Id = jaxb.NewCollapsedString(idStr)
	url := lotlInfo.Url()
	result.Url = &url
	parsingCacheInfo, parsingCacheInfoOK := lotlInfo.TLParsingCacheInfo()
	if parsingCacheInfoOK {
		if parsingCacheInfo.TSLType() != nil {
			tp := parsingCacheInfo.TSLType().URI()
			result.Type = &tp
		}
		territory := parsingCacheInfo.Territory()
		result.CountryCode = &territory
		if !parsingCacheInfo.IssueDate().IsZero() {
			result.IssueDate = jaxb.NewXSDateTime(parsingCacheInfo.IssueDate())
		}
		if !parsingCacheInfo.NextUpdateDate().IsZero() {
			result.NextUpdate = jaxb.NewXSDateTime(parsingCacheInfo.NextUpdateDate())
		}
		if parsingCacheInfo.SequenceNumber() != nil {
			result.SequenceNumber = parsingCacheInfo.SequenceNumber()
		}
		if parsingCacheInfo.Version() != nil {
			result.Version = parsingCacheInfo.Version()
		}
		result.StructuralValidation = b.GetXmlStructuralValidation(parsingCacheInfo.StructureValidationMessages())
	}
	downloadCacheInfo := lotlInfo.DownloadCacheInfo()
	if downloadCacheInfo != nil && !downloadCacheInfo.LastSuccessSynchronizationTime().IsZero() {
		result.LastLoading = jaxb.NewXSDateTime(downloadCacheInfo.LastSuccessSynchronizationTime())
	}
	validationCacheInfo := lotlInfo.ValidationCacheInfo()
	if validationCacheInfo != nil {
		result.WellSigned = validationCacheInfo.IsValid()
	}
	// Java's tlInfoMap.put(id, tlInfo) stores the LOTLInfo itself (Java's LOTLInfo extends
	// TLInfo, so the Map<String, TLInfo> field accepts it upcast). model.tsl.LOTLInfo does not
	// embed/extend model.tsl.TLInfo in this port (Go has no inheritance and the two are
	// deliberately distinct types - see model/tsl/lotl_info.go), so a LOTLInfo cannot be stored
	// into this *TLInfo-typed map without data loss. tlInfoMap is only ever read back through
	// XmlTrustServiceProviderBuilder.getMRA() keyed by a TrustProperties' own TLInfo() (never
	// its LOTLInfo()), so omitting the entry here (rather than inserting a lossy placeholder)
	// is observably equivalent for every caller in this manifest; flagged for the tech lead in
	// case a future TSL/LOTL-owning phase needs a LOTL-keyed lookup here too.
	return result
}

func (b *DiagnosticDataBuilder) linkCertificatesAndTrustedEntities(certificates []*model.CertificateToken) {
	if utils.IsCollectionNotEmpty(certificates) {
		for _, certificateToken := range certificates {
			trustedEntities := NewXmlTrustedEntityBuilder(b.xmlCertsMap, b.xmlTrustSourceMap).
				Build(certificateToken, b.getRelatedTrustedProperties(certificateToken))
			if utils.IsCollectionNotEmpty(trustedEntities) {
				xmlCertificate := b.xmlCertsMap[certificateToken.DSSIDAsString()]
				xmlCertificate.TrustedEntities = &jaxb.TrustedEntitiesWrapper{Items: trustedEntities}
			}
		}
	}
}

func (b *DiagnosticDataBuilder) getRelatedTrustedProperties(certToken *model.CertificateToken) map[*model.CertificateToken][]*lote.TrustedProperties {
	result := make(map[*model.CertificateToken][]*lote.TrustedProperties)
	for _, trustedSource := range b.allCertificateSources.Sources() {
		trustedCertSource, ok := trustedSource.(*spilote.TrustedEntitiesCertificateSource)
		if !ok {
			continue
		}
		processedTokens := make(map[string]bool)
		currentCertificate := certToken
		for currentCertificate != nil {
			trustedProperties := trustedCertSource.TrustedProperties(currentCertificate)
			if utils.IsCollectionNotEmpty(trustedProperties) {
				result[currentCertificate] = append(result[currentCertificate], trustedProperties...)
			}
			if currentCertificate.IsSelfSigned() || processedTokens[currentCertificate.DSSIDAsString()] {
				break
			}
			processedTokens[currentCertificate.DSSIDAsString()] = true

			issuerCertificate := b.getIssuerCertificate(currentCertificate)
			if issuerCertificate != nil && currentCertificate.IssuerEntityKey() != nil &&
				currentCertificate.IssuerEntityKey().Equals(issuerCertificate.EntityKey()) {
				currentCertificate = issuerCertificate
			} else {
				currentCertificate = nil
			}
		}
	}
	return result
}

func (b *DiagnosticDataBuilder) buildXmlLoTEs(trustedCertificateSources *spi.ListCertificateSource) []*jaxb.XmlListOfTrustedEntities {
	trustSourceLists := make([]*jaxb.XmlListOfTrustedEntities, 0)

	mapLists := map[string]*jaxb.XmlListOfTrustedEntities{}
	var mapListsOrder []string
	mapListOfLists := map[string]*jaxb.XmlListOfTrustedEntities{}
	var mapListOfListsOrder []string

	for _, certificateSource := range trustedCertificateSources.Sources() {
		teCertSource, ok := certificateSource.(*spilote.TrustedEntitiesCertificateSource)
		if !ok {
			continue
		}
		summary := teCertSource.Summary()
		if summary != nil {
			b.mergeLoTEMap(mapLists, &mapListsOrder, b.getLoTEMap(teCertSource, summary))
			b.mergeLoTEMap(mapListOfLists, &mapListOfListsOrder, b.getListOfLoTEMap(teCertSource, summary))
		}
		// else: Port of LOG.warn("...LoTEValidationJob is not performed!"): slf4j dropped.
	}

	for _, key := range mapListsOrder {
		trustSourceLists = append(trustSourceLists, mapLists[key])
	}
	for _, key := range mapListOfListsOrder {
		trustSourceLists = append(trustSourceLists, mapListOfLists[key])
	}
	return trustSourceLists
}

func (b *DiagnosticDataBuilder) mergeLoTEMap(dst map[string]*jaxb.XmlListOfTrustedEntities, order *[]string, src map[string]*jaxb.XmlListOfTrustedEntities) {
	for key, value := range src {
		if _, exists := dst[key]; !exists {
			dst[key] = value
			*order = append(*order, key)
		}
	}
}

func (b *DiagnosticDataBuilder) getLoTEMap(teCertSource *spilote.TrustedEntitiesCertificateSource,
	summary *lote.LoTEValidationJobSummary) map[string]*jaxb.XmlListOfTrustedEntities {
	mapTrustedEntitiesLists := map[string]*jaxb.XmlListOfTrustedEntities{}
	for _, loteID := range b.getLOTEIdentifiers(teCertSource) {
		if _, exists := mapTrustedEntitiesLists[loteID.AsXmlID()]; exists {
			continue
		}
		listInfoByID := summary.LoTEInfoByID(loteID)
		if listInfoByID != nil {
			mapTrustedEntitiesLists[loteID.AsXmlID()] = b.getXmlTrustSourceListForLoTE(listInfoByID)
		}
	}
	return mapTrustedEntitiesLists
}

func (b *DiagnosticDataBuilder) getLOTEIdentifiers(teCertSource *spilote.TrustedEntitiesCertificateSource) []model.Identifier {
	var loteIdentifiers []model.Identifier
	seen := map[string]bool{}
	add := func(id model.Identifier) {
		if id == nil || seen[id.AsXmlID()] {
			return
		}
		seen[id.AsXmlID()] = true
		loteIdentifiers = append(loteIdentifiers, id)
	}
	for _, certificateToken := range b.usedCertificates {
		trustedServices := teCertSource.TrustedProperties(certificateToken)
		for _, trustedProperties := range trustedServices {
			listInfo := trustedProperties.LoTEInfo()
			if listInfo != nil {
				add(listInfo.DSSID())
				if listInfo.Parent() != nil {
					add(listInfo.Parent().DSSID())
				}
			}
		}
	}
	return loteIdentifiers
}

func (b *DiagnosticDataBuilder) getListOfLoTEMap(teCertSource *spilote.TrustedEntitiesCertificateSource,
	summary *lote.LoTEValidationJobSummary) map[string]*jaxb.XmlListOfTrustedEntities {
	mapListsOfTrustedEntitiesLists := map[string]*jaxb.XmlListOfTrustedEntities{}
	for _, loloteID := range b.getLoLoTEIdentifiers(teCertSource) {
		if _, exists := mapListsOfTrustedEntitiesLists[loloteID.AsXmlID()]; exists {
			continue
		}
		listOfListsInfoByID := summary.LoLoTEInfoByID(loloteID)
		if listOfListsInfoByID != nil {
			mapListsOfTrustedEntitiesLists[loloteID.AsXmlID()] = b.getXmlTrustSourceListForLoLoTE(listOfListsInfoByID)
		}
	}
	return mapListsOfTrustedEntitiesLists
}

func (b *DiagnosticDataBuilder) getLoLoTEIdentifiers(loteCS *spilote.TrustedEntitiesCertificateSource) []model.Identifier {
	var loloteIdentifiers []model.Identifier
	seen := map[string]bool{}
	for _, certificateToken := range b.usedCertificates {
		trustedServices := loteCS.TrustedProperties(certificateToken)
		for _, trustedProperties := range trustedServices {
			loloteInfo := trustedProperties.LoLoTEInfo()
			if loloteInfo != nil {
				id := loloteInfo.DSSID()
				if !seen[id.AsXmlID()] {
					seen[id.AsXmlID()] = true
					loloteIdentifiers = append(loloteIdentifiers, id)
				}
			}
		}
	}
	return loloteIdentifiers
}

// getXmlTrustSourceListForLoTE builds/looks up the XmlListOfTrustedEntities for a LoTEInfo
// (Java's overload getXmlTrustSourceList(LoTEInfo)).
func (b *DiagnosticDataBuilder) getXmlTrustSourceListForLoTE(loteInfo *lote.LoTEInfo) *jaxb.XmlListOfTrustedEntities {
	id := loteInfo.DSSIDAsString()
	result, ok := b.xmlTrustSourceMap[id]
	if !ok {
		result = &jaxb.XmlListOfTrustedEntities{}
		b.xmlTrustSourceMap[id] = result
	}
	idStr := b.identifierProvider.IDAsString(loteInfo)
	result.Id = jaxb.NewCollapsedString(idStr)
	url := loteInfo.Url()
	result.Url = &url
	if loteInfo.Parent() != nil {
		result.Parent = b.getXmlTrustSourceListForLoLoTE(loteInfo.Parent())
	}
	parsingCacheInfo := loteInfo.ParsingCacheInfo()
	if parsingCacheInfo != nil {
		if parsingCacheInfo.Type() != nil {
			tp := parsingCacheInfo.Type().URI()
			result.Type = &tp
		}
		territory := parsingCacheInfo.Territory()
		result.CountryCode = &territory
		if !parsingCacheInfo.IssueDate().IsZero() {
			result.IssueDate = jaxb.NewXSDateTime(parsingCacheInfo.IssueDate())
		}
		if !parsingCacheInfo.NextUpdateDate().IsZero() {
			result.NextUpdate = jaxb.NewXSDateTime(parsingCacheInfo.NextUpdateDate())
		}
		if parsingCacheInfo.SequenceNumber() != nil {
			result.SequenceNumber = parsingCacheInfo.SequenceNumber()
		}
		if parsingCacheInfo.Version() != nil {
			result.Version = parsingCacheInfo.Version()
		}
		result.StructuralValidation = b.GetXmlStructuralValidation(parsingCacheInfo.StructureValidationMessages())
	}
	downloadCacheInfo := loteInfo.DownloadCacheInfo()
	if downloadCacheInfo != nil && !downloadCacheInfo.LastSuccessSynchronizationTime().IsZero() {
		result.LastLoading = jaxb.NewXSDateTime(downloadCacheInfo.LastSuccessSynchronizationTime())
	}
	validationCacheInfo := loteInfo.ValidationCacheInfo()
	if validationCacheInfo != nil {
		wellSigned := validationCacheInfo.IsValid()
		result.WellSigned = wellSigned
	}
	b.loteInfoMap[id] = loteInfo
	return result
}

// getXmlTrustSourceListForLoLoTE builds/looks up the XmlListOfTrustedEntities for a LoLoTEInfo
// (Java's getXmlTrustSourceList(LoTEInfo) called with a LoLoTEInfo, which is-a LoTEInfo in
// Java; Go has no upcast, so this is a distinct helper against LoLoTEInfo's LoTEInfo-shaped
// accessors).
func (b *DiagnosticDataBuilder) getXmlTrustSourceListForLoLoTE(loloteInfo *lote.LoLoTEInfo) *jaxb.XmlListOfTrustedEntities {
	id := loloteInfo.DSSIDAsString()
	result, ok := b.xmlTrustSourceMap[id]
	if !ok {
		result = &jaxb.XmlListOfTrustedEntities{}
		b.xmlTrustSourceMap[id] = result
	}
	lolote := true
	result.LoLoTE = &lolote
	idStr := b.identifierProvider.IDAsString(loloteInfo)
	result.Id = jaxb.NewCollapsedString(idStr)
	url := loloteInfo.Url()
	result.Url = &url
	parsingCacheInfo := loloteInfo.ParsingCacheInfo()
	if parsingCacheInfo != nil {
		if parsingCacheInfo.Type() != nil {
			tp := parsingCacheInfo.Type().URI()
			result.Type = &tp
		}
		territory := parsingCacheInfo.Territory()
		result.CountryCode = &territory
		if !parsingCacheInfo.IssueDate().IsZero() {
			result.IssueDate = jaxb.NewXSDateTime(parsingCacheInfo.IssueDate())
		}
		if !parsingCacheInfo.NextUpdateDate().IsZero() {
			result.NextUpdate = jaxb.NewXSDateTime(parsingCacheInfo.NextUpdateDate())
		}
		if parsingCacheInfo.SequenceNumber() != nil {
			result.SequenceNumber = parsingCacheInfo.SequenceNumber()
		}
		if parsingCacheInfo.Version() != nil {
			result.Version = parsingCacheInfo.Version()
		}
		result.StructuralValidation = b.GetXmlStructuralValidation(parsingCacheInfo.StructureValidationMessages())
	}
	downloadCacheInfo := loloteInfo.DownloadCacheInfo()
	if downloadCacheInfo != nil && !downloadCacheInfo.LastSuccessSynchronizationTime().IsZero() {
		result.LastLoading = jaxb.NewXSDateTime(downloadCacheInfo.LastSuccessSynchronizationTime())
	}
	validationCacheInfo := loloteInfo.ValidationCacheInfo()
	if validationCacheInfo != nil {
		result.WellSigned = validationCacheInfo.IsValid()
	}
	// Same LoLoTEInfo/LoTEInfo type-mismatch as getXmlTrustSourceListForLOTL's LOTLInfo/TLInfo
	// case above: Java's loteInfoMap.put(id, loteInfo) upcasts the LoLoTEInfo into the
	// Map<String, LoTEInfo> field; this port keeps the two distinct, and loteInfoMap has no
	// reader in this manifest (it is populated but never consulted, matching the Java field's
	// own dead-write shape here), so the entry is simply omitted rather than inserted lossy.
	return result
}

// GetXmlStructuralValidation creates an XmlStructuralValidation for the given errorMessages.
// Port of the protected getXmlStructuralValidation(List<String>).
func (b *DiagnosticDataBuilder) GetXmlStructuralValidation(errorMessages []string) *jaxb.XmlStructuralValidation {
	xmlStructuralValidation := &jaxb.XmlStructuralValidation{}
	xmlStructuralValidation.Valid = utils.IsCollectionEmpty(errorMessages)
	if utils.IsCollectionNotEmpty(errorMessages) {
		xmlStructuralValidation.Message = append(xmlStructuralValidation.Message, errorMessages...)
	}
	return xmlStructuralValidation
}

// GetXmlSignerInfo creates an XmlSignerInfo from a SignerIdentifier. Port of the protected
// getXmlSignerInfo(SignerIdentifier).
func (b *DiagnosticDataBuilder) GetXmlSignerInfo(signerIdentifier *spi.SignerIdentifier) *jaxb.XmlSignerInfo {
	xmlSignerInfo := &jaxb.XmlSignerInfo{}
	if signerIdentifier.IssuerName() != nil {
		issuerName := signerIdentifier.IssuerName().String()
		xmlSignerInfo.IssuerName = &issuerName
	}
	xmlSignerInfo.SerialNumber = jaxb.NewBigInteger(signerIdentifier.SerialNumber())
	if signerIdentifier.Ski() != nil {
		ski := jaxb.Base64Binary(signerIdentifier.Ski())
		xmlSignerInfo.Ski = &ski
	}
	if signerIdentifier.IsCurrent() {
		current := signerIdentifier.IsCurrent()
		xmlSignerInfo.Current = &current
	}
	return xmlSignerInfo
}

func (b *DiagnosticDataBuilder) getXmlSignerInfoForResponderID(responderId *spi.ResponderId) *jaxb.XmlSignerInfo {
	xmlSignerInfo := &jaxb.XmlSignerInfo{}
	if responderId.X500Principal() != nil {
		issuerName := responderId.X500Principal().String()
		xmlSignerInfo.IssuerName = &issuerName
	}
	if responderId.Ski() != nil {
		ski := jaxb.Base64Binary(responderId.Ski())
		xmlSignerInfo.Ski = &ski
	}
	return xmlSignerInfo
}

// BuildDetachedXmlRevocation builds an XmlRevocation from the given RevocationToken. Port of
// the protected buildDetachedXmlRevocation(RevocationToken<?>).
func (b *DiagnosticDataBuilder) BuildDetachedXmlRevocation(revocationToken validation.AnyRevocationToken) *jaxb.XmlRevocation {
	xmlRevocation := &jaxb.XmlRevocation{}
	id := b.identifierProvider.IDAsString(revocationToken)
	xmlRevocation.Id = jaxb.NewCollapsedString(id)

	if revocationToken.IsInternal() {
		origin := jaxb.RevocationOriginValue(enumerations.RevocationOrigin_INPUT_DOCUMENT)
		xmlRevocation.Origin = &origin
	} else {
		origin := jaxb.RevocationOriginValue(revocationToken.ExternalOrigin())
		xmlRevocation.Origin = &origin
	}
	revType := jaxb.RevocationTypeValue(revocationToken.RevocationType())
	xmlRevocation.Type = &revType

	if !revocationToken.ProductionDate().IsZero() {
		xmlRevocation.ProductionDate = jaxb.NewXSDateTime(revocationToken.ProductionDate())
	}
	if !revocationToken.ThisUpdate().IsZero() {
		xmlRevocation.ThisUpdate = jaxb.NewXSDateTime(revocationToken.ThisUpdate())
	}
	if !revocationToken.NextUpdate().IsZero() {
		xmlRevocation.NextUpdate = jaxb.NewXSDateTime(revocationToken.NextUpdate())
	}
	xmlRevocation.CRLNumber = jaxb.NewBigInteger(revocationToken.CRLNumber())
	if !revocationToken.ExpiredCertsOnCRL().IsZero() {
		xmlRevocation.ExpiredCertsOnCRL = jaxb.NewXSDateTime(revocationToken.ExpiredCertsOnCRL())
	}
	if !revocationToken.ArchiveCutOff().IsZero() {
		xmlRevocation.ArchiveCutOff = jaxb.NewXSDateTime(revocationToken.ArchiveCutOff())
	}

	sourceURL := revocationToken.SourceURL()
	if utils.IsStringNotEmpty(sourceURL) { // not empty = online
		xmlRevocation.SourceAddress = &sourceURL
	}

	xmlRevocation.BasicSignature = b.GetXmlBasicSignature(revocationToken)

	xmlRevocation.SigningCertificate = b.getXmlSigningCertificateForTokenAndSource(revocationToken, revocationToken.CertificateSource())
	xmlRevocation.CertificateChain = b.certificateChainWrapper(b.getXmlForCertificateChainWithSource(revocationToken, revocationToken.CertificateSource()))

	certHashPresent := revocationToken.CertHashPresent()
	xmlRevocation.CertHashExtensionPresent = &certHashPresent
	certHashMatch := revocationToken.CertHashMatch()
	xmlRevocation.CertHashExtensionMatch = &certHashMatch

	if revocationToken.CertificateSource() != nil {
		// in case of OCSP token
		xmlRevocation.FoundCertificates = b.GetXmlFoundCertificatesForToken(revocationToken.DSSID(), revocationToken.CertificateSource().(foundCertificatesSource))
	}

	if b.tokenExtractionStrategy.IsRevocationData() {
		bin := jaxb.Base64Binary(revocationToken.Encoded())
		xmlRevocation.Base64Encoded = &bin
	} else {
		revocationDigest, err := revocationToken.Digest(b.defaultDigestAlgorithm)
		if err != nil {
			panic(err)
		}
		xmlRevocation.DigestAlgoAndValue = b.GetXmlDigestAlgoAndValueFor(b.defaultDigestAlgorithm, revocationDigest)
	}

	return xmlRevocation
}

// getXmlRevocationRefs returns a list of XmlRevocationRef for a token with tokenId. Port of the
// protected getXmlRevocationRefs(String, Map<RevocationRef<R>, Set<RevocationRefOrigin>>).
func (b *DiagnosticDataBuilder) getXmlRevocationRefs[R revocation.Revocation](tokenId string, refs []spi.RevocationRefOriginsEntry[R]) []*jaxb.XmlRevocationRef {
	xmlRevocationRefs := make([]*jaxb.XmlRevocationRef, 0)
	for _, entry := range refs {
		ref := entry.Reference
		origins := entry.Origins
		var xmlRef *jaxb.XmlRevocationRef
		if crlRef, ok := any(ref).(*spi.CRLRef); ok {
			xmlRef = b.GetXmlCRLRevocationRef(crlRef, origins)
		} else {
			xmlRef = b.GetXmlOCSPRevocationRef(any(ref).(*spi.OCSPRef), origins)
		}
		b.referenceMap[ref.DSSIDAsString()] = tokenId
		xmlRevocationRefs = append(xmlRevocationRefs, xmlRef)
	}
	return xmlRevocationRefs
}

// GetXmlCRLRevocationRef builds an XmlRevocationRef from a CRLRef. Port of the protected
// getXmlCRLRevocationRef(CRLRef, Set<RevocationRefOrigin>).
func (b *DiagnosticDataBuilder) GetXmlCRLRevocationRef(crlRef *spi.CRLRef, origins []enumerations.RevocationRefOrigin) *jaxb.XmlRevocationRef {
	xmlRevocationRef := &jaxb.XmlRevocationRef{}
	xmlRevocationRef.Origin = revocationRefOriginValues(origins)
	if !crlRef.Digest().IsEmpty() {
		xmlRevocationRef.DigestAlgoAndValue = b.GetXmlDigestAlgoAndValueForDigestValue(crlRef.Digest())
		if crlRef.CRLIssuer() != nil {
			issuer := crlRef.CRLIssuer().String()
			xmlRevocationRef.Issuer = &issuer
		}
		if !crlRef.CRLIssueTime().IsZero() {
			xmlRevocationRef.IssueTime = jaxb.NewXSDateTime(crlRef.CRLIssueTime())
		}
		if crlRef.CRLNumber() != nil {
			xmlRevocationRef.CRLNumber = jaxb.NewBigInteger(crlRef.CRLNumber())
		}
	}
	xmlRevocationRef.CRLNumber = jaxb.NewBigInteger(crlRef.CRLNumber())
	if crlRef.CRLURI() != "" {
		uri := crlRef.CRLURI()
		xmlRevocationRef.Uri = &uri
	}
	return xmlRevocationRef
}

// GetXmlOCSPRevocationRef builds an XmlRevocationRef from an OCSPRef. Port of the protected
// getXmlOCSPRevocationRef(OCSPRef, Set<RevocationRefOrigin>).
func (b *DiagnosticDataBuilder) GetXmlOCSPRevocationRef(ocspRef *spi.OCSPRef, origins []enumerations.RevocationRefOrigin) *jaxb.XmlRevocationRef {
	xmlRevocationRef := &jaxb.XmlRevocationRef{}
	xmlRevocationRef.Origin = revocationRefOriginValues(origins)
	if !ocspRef.Digest().IsEmpty() {
		xmlRevocationRef.DigestAlgoAndValue = b.GetXmlDigestAlgoAndValueForDigestValue(ocspRef.Digest())
	}
	if !ocspRef.ProducedAt().IsZero() {
		xmlRevocationRef.ProducedAt = jaxb.NewXSDateTime(ocspRef.ProducedAt())
	}
	responderId := ocspRef.ResponderId()
	if responderId != nil {
		xmlRevocationRef.ResponderId = b.getXmlSignerInfoForResponderID(responderId)
	}
	if ocspRef.URI() != "" {
		uri := ocspRef.URI()
		xmlRevocationRef.Uri = &uri
	}
	return xmlRevocationRef
}

func revocationRefOriginValues(origins []enumerations.RevocationRefOrigin) []jaxb.RevocationRefOriginValue {
	result := make([]jaxb.RevocationRefOriginValue, 0, len(origins))
	for _, o := range origins {
		result = append(result, jaxb.RevocationRefOriginValue(o))
	}
	return result
}

// GetXmlForCertificateChain returns a certificate chain for the token. Port of the protected
// getXmlForCertificateChain(Token).
func (b *DiagnosticDataBuilder) GetXmlForCertificateChain(token model.Token) []*jaxb.XmlChainItem {
	return b.getXmlForCertificateChainWithSource(token, nil)
}

// getXmlForCertificateChainWithSource returns a certificate chain for the token from the
// certificateSource. Port of the protected getXmlForCertificateChain(Token, CertificateSource).
func (b *DiagnosticDataBuilder) getXmlForCertificateChainWithSource(token model.Token, certificateSource tokenCertificateSourceRefs) []*jaxb.XmlChainItem {
	if token == nil {
		return nil
	}
	certChainTokens := make([]*jaxb.XmlChainItem, 0)
	var processedTokens []model.Token
	processedTokens = append(processedTokens, token)

	issuerToken := b.getIssuerCertificateFromSource(token, certificateSource)
	for issuerToken != nil {
		xmlChainItem := b.getXmlChainItem(issuerToken)
		if xmlChainItem != nil {
			certChainTokens = append(certChainTokens, xmlChainItem)
			if issuerToken.IsSelfSigned() || containsToken(processedTokens, issuerToken) {
				processedTokens = append(processedTokens, issuerToken)
				break
			}
			processedTokens = append(processedTokens, issuerToken)
			issuerToken = b.getIssuerCertificateFromSource(issuerToken, certificateSource)
		} else {
			// not validated -> break
			break
		}
	}

	b.ensureCertificateChain(token, certChainTokens, processedTokens)
	return certChainTokens
}

func containsToken(tokens []model.Token, token *model.CertificateToken) bool {
	for _, t := range tokens {
		if ct, ok := t.(*model.CertificateToken); ok && ct == token {
			return true
		}
	}
	return false
}

func (b *DiagnosticDataBuilder) ensureCertificateChain(token model.Token, certChain []*jaxb.XmlChainItem, processedTokens []model.Token) {
	if utils.IsCollectionNotEmpty(certChain) {
		certificate := b.xmlCertsMap[token.DSSIDAsString()]
		if certificate != nil {
			certificate.SigningCertificate = b.getXmlSigningCertificateFromXmlCertificate(certChain[0].Certificate)
			certificate.CertificateChain = b.certificateChainWrapper(b.getCertChainSinceIndex(certChain, 0))
			certificate.BasicSignature = b.GetXmlBasicSignature(token)
			certificate.IssuerEntityKey = b.getXmlIssuerEntityKey(token)
		}
		for i, chainItem := range certChain {
			chainCertificate := chainItem.Certificate
			if chainCertificate != nil && chainCertificate.SigningCertificate == nil && i+1 < len(certChain) {
				chainCertificate.SigningCertificate = b.getXmlSigningCertificateFromXmlCertificate(certChain[i+1].Certificate)
				chainCertificate.CertificateChain = b.certificateChainWrapper(b.getCertChainSinceIndex(certChain, i+1))
				chainCertificate.BasicSignature = b.GetXmlBasicSignature(processedTokens[i+1])
				chainCertificate.IssuerEntityKey = b.getXmlIssuerEntityKey(processedTokens[i+1])
			}
		}
	}
}

func (b *DiagnosticDataBuilder) getXmlSigningCertificateFromXmlCertificate(xmlCertificate *jaxb.XmlCertificate) *jaxb.XmlSigningCertificate {
	xmlSigningCertificate := &jaxb.XmlSigningCertificate{}
	xmlSigningCertificate.Certificate = xmlCertificate
	return xmlSigningCertificate
}

func (b *DiagnosticDataBuilder) getCertChainSinceIndex(certChain []*jaxb.XmlChainItem, index int) []*jaxb.XmlChainItem {
	result := make([]*jaxb.XmlChainItem, 0, len(certChain)-index)
	for i := index; i < len(certChain); i++ {
		result = append(result, certChain[i])
	}
	return result
}

// GetXmlForCertificateChainForValidity builds a certificate chain for a CertificateValidity.
// Port of the protected getXmlForCertificateChain(CertificateValidity, CertificateSource).
func (b *DiagnosticDataBuilder) GetXmlForCertificateChainForValidity(certificateValidity *spi.CertificateValidity,
	certificateSource tokenCertificateSourceRefs) []*jaxb.XmlChainItem {
	if certificateValidity == nil {
		return nil
	}
	signingCertificate := b.getSigningCertificate(certificateValidity)
	if signingCertificate == nil {
		return nil
	}
	signCertChainItem := b.getXmlChainItem(signingCertificate)
	if signCertChainItem == nil {
		return nil
	}
	certChainTokens := make([]*jaxb.XmlChainItem, 0)
	certChainTokens = append(certChainTokens, signCertChainItem)
	certChain := b.getXmlForCertificateChainWithSource(signingCertificate, certificateSource)
	if utils.IsCollectionNotEmpty(certChain) {
		for _, chainItem := range certChain {
			if chainItem.Certificate != nil && signingCertificate.DSSIDAsString() == chainItem.Certificate.Id.String() {
				break
			}
			certChainTokens = append(certChainTokens, chainItem)
		}
	}
	return certChainTokens
}

func (b *DiagnosticDataBuilder) getXmlChainItem(token *model.CertificateToken) *jaxb.XmlChainItem {
	xmlCertificate, ok := b.xmlCertsMap[token.DSSIDAsString()]
	if ok {
		chainItem := &jaxb.XmlChainItem{}
		chainItem.Certificate = xmlCertificate
		return chainItem
	}
	return nil
}

func (b *DiagnosticDataBuilder) getXmlSigningCertificateForToken(token model.Token) *jaxb.XmlSigningCertificate {
	return b.getXmlSigningCertificateForTokenAndSource(token, nil)
}

// GetXmlSigningCertificate creates the SigningCertificate element for the current token. Port
// of the protected getXmlSigningCertificate(Token, CertificateSource).
func (b *DiagnosticDataBuilder) getXmlSigningCertificateForTokenAndSource(token model.Token, certificateSource tokenCertificateSourceRefs) *jaxb.XmlSigningCertificate {
	xmlSignCertType := &jaxb.XmlSigningCertificate{}
	certificateByPubKey := b.getIssuerCertificateFromSource(token, certificateSource)
	if certificateByPubKey != nil {
		xmlSignCertType.Certificate = b.xmlCertsMap[certificateByPubKey.DSSIDAsString()]
		b.signingCertificateMap[token.DSSIDAsString()] = certificateByPubKey
	} else if token.PublicKeyOfTheSigner() != nil {
		bin := jaxb.Base64Binary(token.PublicKeyOfTheSigner().Encoded())
		xmlSignCertType.PublicKey = &bin
	} else {
		return nil
	}
	return xmlSignCertType
}

func (b *DiagnosticDataBuilder) getIssuerCertificate(token model.Token) *model.CertificateToken {
	return b.getIssuerCertificateFromSource(token, nil)
}

// tokenCertificateSourceRefs is the Go stand-in for Java's abstract CertificateSource
// parameter, restricted to the reference-lookup methods this file actually needs
// (getCertificates()); the concrete argument at every call site is either nil (Java's null) or
// a *spi.TokenCertificateSource-derived value (SignatureCertificateSource/
// RevocationCertificateSource-family) whose Certificates() method already satisfies this.
type tokenCertificateSourceRefs interface {
	Certificates() []*model.CertificateToken
}

func (b *DiagnosticDataBuilder) getIssuerCertificateFromSource(token model.Token, certificateSource tokenCertificateSourceRefs) *model.CertificateToken {
	if token == nil || token.PublicKeyOfTheSigner() == nil {
		return nil
	}

	var issuer *model.CertificateToken
	if certificateSource != nil {
		issuer = b.getBestCertificateFromCandidates(token, certificateSource.Certificates())
	}

	if issuer == nil {
		if cached, ok := b.signingCertificateMap[token.DSSIDAsString()]; ok {
			issuer = cached
		}
	}

	if issuer == nil {
		issuer = b.getBestCertificateFromCandidates(token, b.usedCertificates)
	}

	if issuer != nil {
		issuer = b.getProcessedCertificateToken(issuer)

		if _, exists := b.signingCertificateMap[token.DSSIDAsString()]; !exists {
			b.signingCertificateMap[token.DSSIDAsString()] = issuer
		}
	}

	return issuer
}

func (b *DiagnosticDataBuilder) getBestCertificateFromCandidates(token model.Token, candidates []*model.CertificateToken) *model.CertificateToken {
	return spi.NewTokenIssuerSelector(token, candidates).Issuer()
}

func (b *DiagnosticDataBuilder) getCertsWithPublicKey(publicKey *model.PublicKey, candidates []*model.CertificateToken) []*model.CertificateToken {
	founds := make([]*model.CertificateToken, 0)
	if publicKey != nil {
		for _, cert := range candidates {
			cert = b.getProcessedCertificateToken(cert)
			if publicKey.Equals(cert.PublicKey()) {
				founds = append(founds, cert)
				if b.allCertificateSources.IsTrusted(cert) {
					return []*model.CertificateToken{cert}
				}
			}
		}
	}
	return founds
}

func (b *DiagnosticDataBuilder) getProcessedCertificateToken(certificateToken *model.CertificateToken) *model.CertificateToken {
	processedCertificateToken, ok := b.certificateIdsMap[certificateToken.DSSIDAsString()]
	if !ok {
		processedCertificateToken = certificateToken
		b.certificateIdsMap[certificateToken.DSSIDAsString()] = certificateToken
	}
	return processedCertificateToken
}

// GetXmlSigningCertificateForIdentifier gets a signing certificate token for a token with
// tokenIdentifier. Port of the protected getXmlSigningCertificate(Identifier,
// CertificateValidity).
func (b *DiagnosticDataBuilder) GetXmlSigningCertificateForIdentifier(tokenIdentifier model.Identifier, certificateValidity *spi.CertificateValidity) *jaxb.XmlSigningCertificate {
	xmlSignCertType := &jaxb.XmlSigningCertificate{}
	signingCertificate := b.getSigningCertificate(certificateValidity)
	if signingCertificate != nil {
		xmlSignCertType.Certificate = b.xmlCertsMap[signingCertificate.DSSIDAsString()]
		b.signingCertificateMap[tokenIdentifier.AsXmlID()] = signingCertificate
	} else if certificateValidity.PublicKey() != nil {
		bin := jaxb.Base64Binary(certificateValidity.PublicKey().Encoded())
		xmlSignCertType.PublicKey = &bin
	} else if certificateValidity.SignerInfo() != nil {
		// TODO: add info to xsd
	}
	return xmlSignCertType
}

func (b *DiagnosticDataBuilder) getSigningCertificate(certificateValidity *spi.CertificateValidity) *model.CertificateToken {
	signingCertificateToken := certificateValidity.CertificateToken()
	if signingCertificateToken == nil && certificateValidity.PublicKey() != nil {
		signingCertificateToken = b.getCertificateByPubKey(certificateValidity.PublicKey())
	}
	if signingCertificateToken == nil && certificateValidity.SignerInfo() != nil {
		signingCertificateToken = b.getCertificateByCertificateIdentifier(certificateValidity.SignerInfo())
	}
	if signingCertificateToken != nil {
		signingCertificateToken = b.getProcessedCertificateToken(signingCertificateToken)
	}
	return signingCertificateToken
}

func (b *DiagnosticDataBuilder) getCertificateByPubKey(publicKey *model.PublicKey) *model.CertificateToken {
	if publicKey != nil {
		candidates := b.getCertsWithPublicKey(publicKey, b.usedCertificates)
		if utils.IsCollectionNotEmpty(candidates) {
			return candidates[0]
		}
	}
	return nil
}

func (b *DiagnosticDataBuilder) getCertificateByCertificateIdentifier(signerIdentifier *spi.SignerIdentifier) *model.CertificateToken {
	if signerIdentifier == nil {
		return nil
	}

	founds := make([]*model.CertificateToken, 0)
	for _, cert := range b.usedCertificates {
		related, err := signerIdentifier.IsRelatedToCertificate(cert)
		if err != nil {
			panic(err)
		}
		if related {
			founds = append(founds, cert)
			if b.allCertificateSources.IsTrusted(cert) {
				return cert
			}
		}
	}

	if utils.IsCollectionNotEmpty(founds) {
		return founds[0]
	}
	return nil
}

func (b *DiagnosticDataBuilder) getXmlDistinguishedName(x500PrincipalFormat, value string) *jaxb.XmlDistinguishedName {
	xmlDistinguishedName := &jaxb.XmlDistinguishedName{}
	xmlDistinguishedName.Format = &x500PrincipalFormat
	xmlDistinguishedName.Value = value
	return xmlDistinguishedName
}

func (b *DiagnosticDataBuilder) getCleanedUrls(urls []string) []string {
	cleanedUrls := make([]string, 0, len(urls))
	for _, url := range urls {
		cleanedUrls = append(cleanedUrls, b.getCleanedUrl(url))
	}
	return cleanedUrls
}

func (b *DiagnosticDataBuilder) getCleanedUrl(url string) string {
	return spi.DSSUtilsRemoveControlCharacters(url)
}

// GetXmlFoundCertificatesForSource returns found certificates from the source. Port of the
// protected getXmlFoundCertificates(TokenCertificateSource).
//
// The parameter is the foundCertificatesSource interface, not the concrete *spi.
// TokenCertificateSource: Java's parameter type is the ABSTRACT class TokenCertificateSource,
// so a call like getXmlFoundCertificates(ocspCertificateSource) dispatches
// ocspCertificateSource.getCertificateSourceType() virtually, reaching OCSPCertificateSource's
// override (OCSP_RESPONSE). A concrete Go struct parameter cannot reproduce that: a caller
// holding an *OCSPCertificateSource has no implicit conversion to *spi.TokenCertificateSource,
// so it would have to pass the ADDRESS OF THE EMBEDDED FIELD instead - which is a plain
// *spi.TokenCertificateSource value that has never heard of OCSPCertificateSource's override,
// so CertificateSourceType() resolves to the base's CertificateSourceType_OTHER and
// getXmlFoundCertificates's default/else branch's cast to signatureCertificateSourceRefs panics
// (found live via the phase 8f document-level harness on PAdES-LT.pdf: an orphan OCSP
// revocation identifier's certificate source hit exactly this). The interface parameter lets
// every caller pass the OUTER value it actually has (here, ocspCertificateSource itself),
// which correctly dispatches the override, matching Java.
func (b *DiagnosticDataBuilder) GetXmlFoundCertificatesForSource(certificateSource foundCertificatesSource) *jaxb.XmlFoundCertificates {
	return b.getXmlFoundCertificates(nil, certificateSource)
}

// GetXmlFoundCertificatesForToken returns found certificates from the source, given the owning
// token's identifier. Port of the protected getXmlFoundCertificates(Identifier,
// TokenCertificateSource) as used by revocation tokens, whose CertificateSource() is typed
// spi.RevocationCertificateSource (narrower than *spi.TokenCertificateSource); the reference
// certificate-source methods below are the ones this file needs (Certificates(),
// CertificateSourceType(), ReferencesForCertificateToken(), CertificateRefOrigins(),
// OrphanCertificateRefs()), which every concrete certificate source embeds
// spi.TokenCertificateSource for and therefore has - see foundCertificatesSource.
func (b *DiagnosticDataBuilder) GetXmlFoundCertificatesForToken(tokenIdentifier model.Identifier, certificateSource foundCertificatesSource) *jaxb.XmlFoundCertificates {
	return b.getXmlFoundCertificates(tokenIdentifier, certificateSource)
}

// foundCertificatesSource is the Go stand-in for Java's abstract TokenCertificateSource
// parameter of getXmlFoundCertificates/populateCertificateOriginMap/getXmlOrphanCertificate/
// etc: every method this file calls on a "TokenCertificateSource certificateSource" local,
// collected into one interface. *spi.SignatureCertificateSource, *spi.TokenCertificateSource
// and *spi.OCSPCertificateSource (via its embedded spi.RevocationCertificateSourceBase ->
// spi.TokenCertificateSource) all satisfy it structurally. See the file header's map/interface
// note: spi.RevocationCertificateSource (the static return type of RevocationToken.
// CertificateSource()) only re-exports spi.CertificateSource, not these
// spi.TokenCertificateSource-only members, so callers passing a revocation token's certificate
// source type-assert it to this interface first (see BuildDetachedXmlRevocation).
type foundCertificatesSource interface {
	CertificateSourceType() enumerations.CertificateSourceType
	Certificates() []*model.CertificateToken
	ReferencesForCertificateToken(certificateToken *model.CertificateToken) []*spi.CertificateRef
	CertificateRefOrigins(certificateRef *spi.CertificateRef) []enumerations.CertificateRefOrigin
	OrphanCertificateRefs() []*spi.CertificateRef
}

// signatureCertificateSourceRefs additionally exposes the members
// getXmlRelatedCertificates/getOrphanCertificates need only for the "else" (signature-shaped)
// branch, matching Java's cast to SignatureCertificateSource.
type signatureCertificateSourceRefs interface {
	foundCertificatesSource
	KeyInfoCertificates() []*model.CertificateToken
	SignedDataCertificates() []*model.CertificateToken
	CertificateValues() []*model.CertificateToken
	AttrAuthoritiesCertValues() []*model.CertificateToken
	TimeStampValidationDataCertValues() []*model.CertificateToken
	AnyValidationDataCertValues() []*model.CertificateToken
	DSSDictionaryCertValues() []*model.CertificateToken
	VRIDictionaryCertValues() []*model.CertificateToken
	UnprotectedHeaderCertificates() []*model.CertificateToken
}

func (b *DiagnosticDataBuilder) getXmlFoundCertificates(tokenIdentifier model.Identifier, certificateSource foundCertificatesSource) *jaxb.XmlFoundCertificates {
	xmlFoundCertificates := &jaxb.XmlFoundCertificates{}
	xmlFoundCertificates.RelatedCertificate = b.GetXmlRelatedCertificates(certificateSource)
	xmlRelatedCertificatesForOrphanReferences := b.GetXmlRelatedCertificateForOrphanReferences(certificateSource)
	for _, xmlRelatedCertificate := range xmlRelatedCertificatesForOrphanReferences {
		if !b.containsCertificate(xmlFoundCertificates.RelatedCertificate, xmlRelatedCertificate) {
			xmlFoundCertificates.RelatedCertificate = append(xmlFoundCertificates.RelatedCertificate, xmlRelatedCertificate)
		}
	}
	if tokenIdentifier != nil {
		signingCertificate := b.signingCertificateMap[tokenIdentifier.AsXmlID()]
		// safe null processing is implemented inside (to create orphan references, when needed)
		xmlFoundCertificates.OrphanCertificate = append(xmlFoundCertificates.OrphanCertificate, b.getOrphanCertificates(certificateSource, signingCertificate)...)
		xmlFoundCertificates.OrphanCertificate = append(xmlFoundCertificates.OrphanCertificate, b.getOrphanCertificateRefs(certificateSource, signingCertificate)...)
	}
	return xmlFoundCertificates
}

func (b *DiagnosticDataBuilder) containsCertificate(certificates []*jaxb.XmlRelatedCertificate, xmlRelatedCertificate *jaxb.XmlRelatedCertificate) bool {
	for _, c := range certificates {
		if xmlRelatedCertificate.Certificate.Id.String() == c.Certificate.Id.String() {
			return true
		}
	}
	return false
}

// GetXmlRelatedCertificates is the entry used directly by signature/timestamp/evidence-record
// callers whose certificateSource statically carries the extra SignatureCertificateSource
// accessors already (embedding satisfies signatureCertificateSourceRefs structurally); it
// dispatches on CertificateSourceType() exactly as Java's getXmlRelatedCertificates does. Port
// of the private getXmlRelatedCertificates(TokenCertificateSource).
func (b *DiagnosticDataBuilder) GetXmlRelatedCertificates(certificateSource foundCertificatesSource) []*jaxb.XmlRelatedCertificate {
	relatedCertificatesMap := map[string]*jaxb.XmlRelatedCertificate{}
	var order []string
	add := func(origin enumerations.CertificateOrigin, tokens []*model.CertificateToken) {
		b.populateCertificateOriginMap(relatedCertificatesMap, &order, origin, tokens, certificateSource)
	}

	switch certificateSource.CertificateSourceType() {
	case enumerations.CertificateSourceType_OCSP_RESPONSE:
		add(enumerations.CertificateOrigin_BASIC_OCSP_RESP, certificateSource.Certificates())
	case enumerations.CertificateSourceType_EVIDENCE_RECORD:
		add(enumerations.CertificateOrigin_EVIDENCE_RECORD, certificateSource.Certificates())
	case enumerations.CertificateSourceType_EAA:
		add(enumerations.CertificateOrigin_EAA, certificateSource.Certificates())
	default:
		signatureCertificateSource := certificateSource.(signatureCertificateSourceRefs)
		add(enumerations.CertificateOrigin_KEY_INFO, signatureCertificateSource.KeyInfoCertificates())
		add(enumerations.CertificateOrigin_SIGNED_DATA, signatureCertificateSource.SignedDataCertificates())
		add(enumerations.CertificateOrigin_CERTIFICATE_VALUES, signatureCertificateSource.CertificateValues())
		add(enumerations.CertificateOrigin_ATTR_AUTHORITIES_CERT_VALUES, signatureCertificateSource.AttrAuthoritiesCertValues())
		add(enumerations.CertificateOrigin_TIMESTAMP_VALIDATION_DATA, signatureCertificateSource.TimeStampValidationDataCertValues())
		add(enumerations.CertificateOrigin_ANY_VALIDATION_DATA, signatureCertificateSource.AnyValidationDataCertValues())
		add(enumerations.CertificateOrigin_DSS_DICTIONARY, signatureCertificateSource.DSSDictionaryCertValues())
		add(enumerations.CertificateOrigin_VRI_DICTIONARY, signatureCertificateSource.VRIDictionaryCertValues())
		add(enumerations.CertificateOrigin_UNPROTECTED_HEADER, signatureCertificateSource.UnprotectedHeaderCertificates())
	}

	result := make([]*jaxb.XmlRelatedCertificate, 0, len(order))
	for _, key := range order {
		result = append(result, relatedCertificatesMap[key])
	}
	return result
}

// PopulateCertificateOriginMap fills the certificates origins map with the given properties.
// Port of the protected populateCertificateOriginMap(Map<String,XmlRelatedCertificate>,
// CertificateOrigin, List<CertificateToken>, TokenCertificateSource).
func (b *DiagnosticDataBuilder) populateCertificateOriginMap(relatedCertificatesMap map[string]*jaxb.XmlRelatedCertificate, order *[]string,
	origin enumerations.CertificateOrigin, certificateTokens []*model.CertificateToken, certificateSource foundCertificatesSource) {
	for _, certificateToken := range certificateTokens {
		id := certificateToken.DSSIDAsString()
		if stored, ok := relatedCertificatesMap[id]; !ok {
			if _, known := b.xmlCertsMap[id]; known {
				xmlFoundCertificate := b.PopulateXmlRelatedCertificatesList(origin, certificateToken, certificateSource)
				relatedCertificatesMap[id] = xmlFoundCertificate
				*order = append(*order, id)
			}
		} else if !containsCertificateOrigin(stored.Origin, origin) {
			stored.Origin = append(stored.Origin, jaxb.CertificateOriginValue(origin))
		}
	}
}

func containsCertificateOrigin(origins []jaxb.CertificateOriginValue, origin enumerations.CertificateOrigin) bool {
	for _, o := range origins {
		if enumerations.CertificateOrigin(o) == origin {
			return true
		}
	}
	return false
}

// PopulateXmlRelatedCertificatesList builds an XmlRelatedCertificate. Port of the protected
// populateXmlRelatedCertificatesList(CertificateOrigin, CertificateToken, TokenCertificateSource).
func (b *DiagnosticDataBuilder) PopulateXmlRelatedCertificatesList(origin enumerations.CertificateOrigin, cert *model.CertificateToken,
	certificateSource foundCertificatesSource) *jaxb.XmlRelatedCertificate {
	xrc := &jaxb.XmlRelatedCertificate{}
	xrc.Origin = append(xrc.Origin, jaxb.CertificateOriginValue(origin))
	xrc.Certificate = b.xmlCertsMap[cert.DSSIDAsString()]
	referencesForCertificateToken := certificateSource.ReferencesForCertificateToken(cert)
	for _, certificateRef := range referencesForCertificateToken {
		for _, refOrigin := range certificateSource.CertificateRefOrigins(certificateRef) {
			xmlCertificateRef := b.GetXmlCertificateRef(certificateRef, refOrigin)
			b.VerifyAgainstCertificateToken(xmlCertificateRef, certificateRef, cert)
			xrc.CertificateRef = append(xrc.CertificateRef, xmlCertificateRef)
		}
		b.referenceMap[certificateRef.DSSIDAsString()] = cert.DSSIDAsString()
	}
	return xrc
}

// PopulateXmlRelatedCertificatesListInto builds an XmlRelatedCertificate and populates the
// relatedCertificates list. Port of the protected populateXmlRelatedCertificatesList(
// List<XmlRelatedCertificate>, TokenCertificateSource, CertificateToken, CertificateRef).
func (b *DiagnosticDataBuilder) PopulateXmlRelatedCertificatesListInto(relatedCertificates []*jaxb.XmlRelatedCertificate,
	certificateSource foundCertificatesSource, cert *model.CertificateToken, certificateRef *spi.CertificateRef) []*jaxb.XmlRelatedCertificate {
	xrc := b.getXmlRelatedCertificateWithId(relatedCertificates, b.identifierProvider.IDAsString(cert))
	if xrc == nil {
		xrc = &jaxb.XmlRelatedCertificate{}
		xrc.Certificate = b.xmlCertsMap[cert.DSSIDAsString()]
		relatedCertificates = append(relatedCertificates, xrc)
	}
	for _, refOrigin := range certificateSource.CertificateRefOrigins(certificateRef) {
		xmlCertificateRef := b.GetXmlCertificateRef(certificateRef, refOrigin)
		b.VerifyAgainstCertificateToken(xmlCertificateRef, certificateRef, cert)
		xrc.CertificateRef = append(xrc.CertificateRef, xmlCertificateRef)
	}
	b.referenceMap[certificateRef.DSSIDAsString()] = cert.DSSIDAsString()
	return relatedCertificates
}

func (b *DiagnosticDataBuilder) getXmlRelatedCertificateWithId(relatedCertificates []*jaxb.XmlRelatedCertificate, certId string) *jaxb.XmlRelatedCertificate {
	for _, relatedCertificate := range relatedCertificates {
		if certId == relatedCertificate.Certificate.Id.String() {
			return relatedCertificate
		}
	}
	return nil
}

// GetXmlCertificateRef builds an XmlCertificateRef from a CertificateRef. Port of the protected
// getXmlCertificateRef(CertificateRef, CertificateRefOrigin).
func (b *DiagnosticDataBuilder) GetXmlCertificateRef(ref *spi.CertificateRef, origin enumerations.CertificateRefOrigin) *jaxb.XmlCertificateRef {
	certificateRef := &jaxb.XmlCertificateRef{}
	signerIdentifier := ref.CertificateIdentifier()
	if signerIdentifier != nil {
		certificateRef.IssuerSerial = b.getXmlIssuerSerial(signerIdentifier)
	}
	refDigest := ref.CertDigest()
	responderId := ref.ResponderId()
	if !refDigest.IsEmpty() {
		certificateRef.DigestAlgoAndValue = b.GetXmlDigestAlgoAndValueFor(refDigest.Algorithm(), refDigest.Value())
	} else if signerIdentifier != nil {
		certificateRef.SerialInfo = b.GetXmlSignerInfo(signerIdentifier)
	} else if responderId != nil {
		certificateRef.SerialInfo = b.getXmlSignerInfoForResponderID(responderId)
	}
	if kid := ref.Kid(); kid != "" {
		certificateRef.KID = &kid
	}
	if x509Url := ref.X509Url(); x509Url != "" {
		certificateRef.X509Url = &x509Url
	}
	originVal := jaxb.CertificateRefOriginValue(origin)
	certificateRef.Origin = &originVal
	return certificateRef
}

func (b *DiagnosticDataBuilder) getOrphanCertificates(certificateSource foundCertificatesSource, signingCertificate *model.CertificateToken) []*jaxb.XmlOrphanCertificate {
	orphanCertificatesMap := map[string]*jaxb.XmlOrphanCertificate{}
	var order []string
	add := func(origin enumerations.CertificateOrigin, tokens []*model.CertificateToken) {
		b.populateOrphanCertificateOriginMap(orphanCertificatesMap, &order, origin, tokens, certificateSource, signingCertificate)
	}

	switch certificateSource.CertificateSourceType() {
	case enumerations.CertificateSourceType_OCSP_RESPONSE:
		add(enumerations.CertificateOrigin_BASIC_OCSP_RESP, certificateSource.Certificates())
	case enumerations.CertificateSourceType_EAA:
		add(enumerations.CertificateOrigin_EAA, certificateSource.Certificates())
	default:
		signatureCertificateSource := certificateSource.(signatureCertificateSourceRefs)
		add(enumerations.CertificateOrigin_KEY_INFO, signatureCertificateSource.KeyInfoCertificates())
		add(enumerations.CertificateOrigin_SIGNED_DATA, signatureCertificateSource.SignedDataCertificates())
		add(enumerations.CertificateOrigin_CERTIFICATE_VALUES, signatureCertificateSource.CertificateValues())
		add(enumerations.CertificateOrigin_ATTR_AUTHORITIES_CERT_VALUES, signatureCertificateSource.AttrAuthoritiesCertValues())
		add(enumerations.CertificateOrigin_TIMESTAMP_VALIDATION_DATA, signatureCertificateSource.TimeStampValidationDataCertValues())
		add(enumerations.CertificateOrigin_ANY_VALIDATION_DATA, signatureCertificateSource.AnyValidationDataCertValues())
		add(enumerations.CertificateOrigin_DSS_DICTIONARY, signatureCertificateSource.DSSDictionaryCertValues())
		add(enumerations.CertificateOrigin_VRI_DICTIONARY, signatureCertificateSource.VRIDictionaryCertValues())
		add(enumerations.CertificateOrigin_UNPROTECTED_HEADER, signatureCertificateSource.UnprotectedHeaderCertificates())
	}

	result := make([]*jaxb.XmlOrphanCertificate, 0, len(order))
	for _, key := range order {
		result = append(result, orphanCertificatesMap[key])
	}
	return result
}

// PopulateOrphanCertificateOriginMap fills the orphan certificate map with the given values.
// Port of the protected populateOrphanCertificateOriginMap(Map<String,XmlOrphanCertificate>,
// CertificateOrigin, List<CertificateToken>, TokenCertificateSource, CertificateToken).
func (b *DiagnosticDataBuilder) populateOrphanCertificateOriginMap(orphanCertificatesMap map[string]*jaxb.XmlOrphanCertificate, order *[]string,
	origin enumerations.CertificateOrigin, certificateTokens []*model.CertificateToken, certificateSource foundCertificatesSource,
	signingCertificate *model.CertificateToken) {
	for _, certificateToken := range certificateTokens {
		id := certificateToken.DSSIDAsString()
		if _, known := b.xmlCertsMap[id]; known {
			continue
		}
		if stored, ok := orphanCertificatesMap[id]; !ok {
			xmlOrphanCertificate := b.GetXmlOrphanCertificate(origin, certificateToken, certificateSource, signingCertificate)
			orphanCertificatesMap[id] = xmlOrphanCertificate
			*order = append(*order, id)
		} else if !containsCertificateOrigin(stored.Origin, origin) {
			stored.Origin = append(stored.Origin, jaxb.CertificateOriginValue(origin))
		}
	}
}

// GetXmlOrphanCertificate builds an XmlOrphanCertificateToken. Port of the protected
// getXmlOrphanCertificate(CertificateOrigin, CertificateToken, TokenCertificateSource,
// CertificateToken).
func (b *DiagnosticDataBuilder) GetXmlOrphanCertificate(origin enumerations.CertificateOrigin, certificateToken *model.CertificateToken,
	certificateSource foundCertificatesSource, signingCertificate *model.CertificateToken) *jaxb.XmlOrphanCertificate {
	xoc := &jaxb.XmlOrphanCertificate{}
	xoc.Origin = append(xoc.Origin, jaxb.CertificateOriginValue(origin))
	xoc.Token = b.BuildXmlOrphanCertificateToken(certificateToken)
	referencesForCertificateToken := certificateSource.ReferencesForCertificateToken(certificateToken)
	for _, certificateRef := range referencesForCertificateToken {
		for _, refOrigin := range certificateSource.CertificateRefOrigins(certificateRef) {
			xmlCertificateRef := b.GetXmlCertificateRef(certificateRef, refOrigin)
			b.VerifyAgainstCertificateToken(xmlCertificateRef, certificateRef, signingCertificate)
			xoc.CertificateRef = append(xoc.CertificateRef, xmlCertificateRef)
		}
		b.referenceMap[certificateRef.DSSIDAsString()] = certificateToken.DSSIDAsString()
	}
	return xoc
}

// IsKnownCertificate reports whether id (a CertificateToken.DSSIDAsString()) has already been
// recorded as a non-orphan XmlCertificate (i.e. is a key of the private xmlCertsMap cache).
// Cross-package accessor added during phase 8f un-gating: Java's PAdESDiagnosticDataBuilder.
// buildOrphanTokensFromDocumentSources() reads the protected xmlCertsMap field directly (Java
// `protected` grants cross-package subclass access DSS relies on here); Go embedding does not
// expose unexported fields to an embedding type in another package, so this getter is the
// narrowest surface that reproduces the same check. No existing behavior changes - purely
// additive.
func (b *DiagnosticDataBuilder) IsKnownCertificate(id string) bool {
	_, ok := b.xmlCertsMap[id]
	return ok
}

// IsKnownRevocation reports whether id (a revocation identifier's AsXmlID()) has already been
// recorded as a non-orphan XmlRevocation (i.e. is a key of the private xmlRevocationsMap cache).
// See IsKnownCertificate's doc comment for why this accessor exists.
func (b *DiagnosticDataBuilder) IsKnownRevocation(id string) bool {
	_, ok := b.xmlRevocationsMap[id]
	return ok
}

// BuildXmlOrphanCertificateToken builds an XmlOrphanCertificateToken from the given
// CertificateToken. Port of the protected buildXmlOrphanCertificateToken(CertificateToken).
func (b *DiagnosticDataBuilder) BuildXmlOrphanCertificateToken(certificateToken *model.CertificateToken) *jaxb.XmlOrphanCertificateToken {
	id := certificateToken.DSSIDAsString()
	orphanToken, ok := b.xmlOrphanCertificateTokensMap[id]
	if !ok {
		orphanToken = &jaxb.XmlOrphanCertificateToken{}
		encType := jaxb.XmlEncapsulationType_BINARIES
		orphanToken.EncapsulationType = &encType
		idStr := b.identifierProvider.IDAsString(certificateToken)
		orphanToken.Id = jaxb.NewCollapsedString(idStr)

		subject := certificateToken.Subject()
		orphanToken.SubjectDistinguishedName = append(orphanToken.SubjectDistinguishedName,
			b.getXmlDistinguishedName(x500PrincipalCanonical, subject.Canonical()))
		orphanToken.SubjectDistinguishedName = append(orphanToken.SubjectDistinguishedName,
			b.getXmlDistinguishedName(x500PrincipalRFC2253, subject.RFC2253()))

		issuer := certificateToken.Issuer()
		orphanToken.IssuerDistinguishedName = append(orphanToken.IssuerDistinguishedName,
			b.getXmlDistinguishedName(x500PrincipalCanonical, issuer.Canonical()))
		orphanToken.IssuerDistinguishedName = append(orphanToken.IssuerDistinguishedName,
			b.getXmlDistinguishedName(x500PrincipalRFC2253, issuer.RFC2253()))

		orphanToken.SerialNumber = jaxb.NewBigInteger(certificateToken.SerialNumber())

		orphanToken.NotAfter = jaxb.NewXSDateTime(certificateToken.NotAfter())
		orphanToken.NotBefore = jaxb.NewXSDateTime(certificateToken.NotBefore())

		entityKey := certificateToken.EntityKey().AsXmlID()
		orphanToken.EntityKey = &entityKey

		selfSigned := certificateToken.IsSelfSigned()
		orphanToken.SelfSigned = &selfSigned
		trusted := b.allCertificateSources.IsTrusted(certificateToken)
		orphanToken.Trusted = &trusted

		if b.tokenExtractionStrategy.IsCertificate() {
			bin := jaxb.Base64Binary(certificateToken.Encoded())
			orphanToken.Base64Encoded = &bin
		} else {
			certDigest, err := certificateToken.Digest(b.defaultDigestAlgorithm)
			if err != nil {
				panic(err)
			}
			orphanToken.DigestAlgoAndValue = b.GetXmlDigestAlgoAndValueFor(b.defaultDigestAlgorithm, certDigest)
		}

		b.xmlOrphanCertificateTokensMap[id] = orphanToken
		b.xmlOrphanCertificateTokensOrder = append(b.xmlOrphanCertificateTokensOrder, id)
	}
	return orphanToken
}

func (b *DiagnosticDataBuilder) getOrphanCertificateRefs(certificateSource foundCertificatesSource, signingCertificate *model.CertificateToken) []*jaxb.XmlOrphanCertificate {
	orphanCertificates := make([]*jaxb.XmlOrphanCertificate, 0)
	orphanCertificateRefs := certificateSource.OrphanCertificateRefs()
	for _, orphanCertificateRef := range orphanCertificateRefs {
		if utils.IsCollectionEmpty(b.GetUsedCertificatesByCertificateRef(orphanCertificateRef)) {
			orphanCertificates = append(orphanCertificates, b.createXmlOrphanCertificateFromRef(certificateSource, orphanCertificateRef, signingCertificate))
		}
	}
	return orphanCertificates
}

func (b *DiagnosticDataBuilder) createXmlOrphanCertificateFromRef(certificateSource foundCertificatesSource,
	orphanCertificateRef *spi.CertificateRef, signingCertificate *model.CertificateToken) *jaxb.XmlOrphanCertificate {
	orphanCertificate := &jaxb.XmlOrphanCertificate{}
	orphanCertificate.Token = b.getXmlOrphanCertificateTokenFromRef(orphanCertificateRef)
	for _, refOrigin := range certificateSource.CertificateRefOrigins(orphanCertificateRef) {
		xmlCertificateRef := b.GetXmlCertificateRef(orphanCertificateRef, refOrigin)
		b.VerifyAgainstCertificateToken(xmlCertificateRef, orphanCertificateRef, signingCertificate)
		orphanCertificate.CertificateRef = append(orphanCertificate.CertificateRef, xmlCertificateRef)
	}
	return orphanCertificate
}

func (b *DiagnosticDataBuilder) getXmlOrphanCertificateTokenFromRef(orphanCertificateRef *spi.CertificateRef) *jaxb.XmlOrphanCertificateToken {
	id := orphanCertificateRef.DSSIDAsString()
	orphanToken, ok := b.xmlOrphanCertificateTokensMap[id]
	if !ok {
		orphanToken = &jaxb.XmlOrphanCertificateToken{}
		encType := jaxb.XmlEncapsulationType_REFERENCE
		orphanToken.EncapsulationType = &encType
		idStr := b.identifierProvider.IDAsString(orphanCertificateRef)
		orphanToken.Id = jaxb.NewCollapsedString(idStr)
		if !orphanCertificateRef.CertDigest().IsEmpty() {
			orphanToken.DigestAlgoAndValue = b.GetXmlDigestAlgoAndValueForDigestValue(orphanCertificateRef.CertDigest())
		}
		b.xmlOrphanCertificateTokensMap[id] = orphanToken
		b.xmlOrphanCertificateTokensOrder = append(b.xmlOrphanCertificateTokensOrder, id)
	}
	return orphanToken
}

// GetXmlRelatedCertificateForOrphanReferences returns a list of XmlRelatedCertificates for
// orphan references within certificateSource. Port of the protected
// getXmlRelatedCertificateForOrphanReferences(TokenCertificateSource).
func (b *DiagnosticDataBuilder) GetXmlRelatedCertificateForOrphanReferences(certificateSource foundCertificatesSource) []*jaxb.XmlRelatedCertificate {
	relatedCertificates := make([]*jaxb.XmlRelatedCertificate, 0)
	for _, certificateRef := range certificateSource.OrphanCertificateRefs() {
		certificateTokens := b.GetUsedCertificatesByCertificateRef(certificateRef)
		if utils.IsCollectionNotEmpty(certificateTokens) {
			for _, certificateToken := range certificateTokens {
				relatedCertificates = b.PopulateXmlRelatedCertificatesListInto(relatedCertificates, certificateSource, certificateToken, certificateRef)
			}
		}
	}
	return relatedCertificates
}

// GetUsedCertificatesByCertificateRef returns used certificates matched by the certificateRef.
// Port of the protected getUsedCertificatesByCertificateRef(CertificateRef).
func (b *DiagnosticDataBuilder) GetUsedCertificatesByCertificateRef(certificateRef *spi.CertificateRef) []*model.CertificateToken {
	matcher := &spi.CertificateTokenRefMatcher{}
	tokensFromRefs := b.allCertificateSources.FindTokensFromCertRef(certificateRef)

	certificates := make([]*model.CertificateToken, 0)
	seen := map[string]bool{}
	for _, certificateToken := range b.usedCertificates {
		id := certificateToken.DSSIDAsString()
		if seen[id] {
			continue
		}
		if _, found := tokensFromRefs[id]; found || matcher.Match(certificateToken, certificateRef) {
			seen[id] = true
			certificates = append(certificates, certificateToken)
		}
	}
	return certificates
}

// VerifyAgainstCertificateToken verifies the reference against a certificate token. Port of the
// protected verifyAgainstCertificateToken(XmlCertificateRef, CertificateRef, CertificateToken).
func (b *DiagnosticDataBuilder) VerifyAgainstCertificateToken(xmlCertificateRef *jaxb.XmlCertificateRef, ref *spi.CertificateRef, signingCertificate *model.CertificateToken) {
	tokenRefMatcher := &spi.CertificateTokenRefMatcher{}
	digestAlgoAndValue := xmlCertificateRef.DigestAlgoAndValue
	if digestAlgoAndValue != nil {
		match := signingCertificate != nil && tokenRefMatcher.MatchByDigest(signingCertificate, ref)
		digestAlgoAndValue.Match = &match
	}
	issuerSerial := xmlCertificateRef.IssuerSerial
	if issuerSerial != nil {
		match := signingCertificate != nil && tokenRefMatcher.MatchByIssuerName(signingCertificate, ref) &&
			tokenRefMatcher.MatchBySerialNumber(signingCertificate, ref)
		issuerSerial.Match = &match
	}
}

func (b *DiagnosticDataBuilder) getXmlIssuerSerial(signerIdentifier *spi.SignerIdentifier) *jaxb.XmlIssuerSerial {
	xmlIssuerSerial := &jaxb.XmlIssuerSerial{}
	xmlIssuerSerial.Value = jaxb.Base64Binary(signerIdentifier.IssuerSerialEncoded())
	return xmlIssuerSerial
}

func (b *DiagnosticDataBuilder) getXmlIssuerEntityKey(token model.Token) *jaxb.XmlIssuerEntityKey {
	var issuerCertificate *model.CertificateToken
	if token.IsSelfSigned() {
		issuerCertificate = token.(*model.CertificateToken)
	} else {
		issuerCertificate = b.getIssuerCertificate(token)
	}
	if issuerCertificate != nil {
		xmlIssuerEntityKey := &jaxb.XmlIssuerEntityKey{}
		xmlIssuerEntityKey.Value = token.IssuerEntityKey().AsXmlID()
		var publicKeyOfSigner *model.PublicKey
		if token.IsSelfSigned() {
			publicKeyOfSigner = issuerCertificate.PublicKey()
		} else {
			publicKeyOfSigner = token.PublicKeyOfTheSigner()
		}
		xmlIssuerEntityKey.Key = issuerCertificate.PublicKey() != nil && publicKeyOfSigner != nil &&
			bytesEqual(publicKeyOfSigner.Encoded(), issuerCertificate.PublicKey().Encoded())
		xmlIssuerEntityKey.SubjectName = model.NewX500PrincipalHelper(token.IssuerX500Principal()).Equals(issuerCertificate.Subject())
		return xmlIssuerEntityKey
	}
	return nil
}

func bytesEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// GetXmlBasicSignature gets an XmlBasicSignature for a Token. Port of the protected
// getXmlBasicSignature(Token).
func (b *DiagnosticDataBuilder) GetXmlBasicSignature(token model.Token) *jaxb.XmlBasicSignature {
	xmlBasicSignatureType := &jaxb.XmlBasicSignature{}

	signatureAlgorithm := token.SignatureAlgorithm()
	if signatureAlgorithm != "" {
		encAlgo := jaxb.EncryptionAlgorithmValue(signatureAlgorithm.EncryptionAlgorithm())
		xmlBasicSignatureType.EncryptionAlgoUsedToSignThisToken = &encAlgo
		digestAlgo := jaxb.DigestAlgorithmValue(signatureAlgorithm.DigestAlgorithm())
		xmlBasicSignatureType.DigestAlgoUsedToSignThisToken = &digestAlgo
	}
	keyLength := spi.DSSPKUtilsStringPublicKeySizeOfToken(token)
	xmlBasicSignatureType.KeyLengthUsedToSignThisToken = &keyLength

	signatureValidity := token.SignatureValidity()
	if enumerations.SignatureValidity_NOT_EVALUATED != signatureValidity {
		signatureIntact := token.IsSignatureIntact()
		xmlBasicSignatureType.SignatureIntact = &signatureIntact
		signatureValid := token.IsValid()
		xmlBasicSignatureType.SignatureValid = &signatureValid
	}
	return xmlBasicSignatureType
}

// BuildDetachedXmlCertificate builds an XmlCertificate from the given CertificateToken. Port of
// the protected buildDetachedXmlCertificate(CertificateToken).
func (b *DiagnosticDataBuilder) BuildDetachedXmlCertificate(certToken *model.CertificateToken) *jaxb.XmlCertificate {
	xmlCert := &jaxb.XmlCertificate{}
	id := b.identifierProvider.IDAsString(certToken)
	xmlCert.Id = jaxb.NewCollapsedString(id)

	subject := certToken.Subject()
	xmlCert.SubjectDistinguishedName = append(xmlCert.SubjectDistinguishedName,
		b.getXmlDistinguishedName(x500PrincipalCanonical, subject.Canonical()))
	xmlCert.SubjectDistinguishedName = append(xmlCert.SubjectDistinguishedName,
		b.getXmlDistinguishedName(x500PrincipalRFC2253, subject.RFC2253()))

	issuer := certToken.Issuer()
	xmlCert.IssuerDistinguishedName = append(xmlCert.IssuerDistinguishedName,
		b.getXmlDistinguishedName(x500PrincipalCanonical, issuer.Canonical()))
	xmlCert.IssuerDistinguishedName = append(xmlCert.IssuerDistinguishedName,
		b.getXmlDistinguishedName(x500PrincipalRFC2253, issuer.RFC2253()))

	xmlCert.SerialNumber = jaxb.NewBigInteger(certToken.SerialNumber())

	if v := spi.DSSASN1UtilsExtractAttributeFromX500Principal(oidSubjectSerialNumber, subject); v != "" {
		xmlCert.SubjectSerialNumber = &v
	}
	if v := spi.DSSASN1UtilsExtractAttributeFromX500Principal(oidTitle, subject); v != "" {
		xmlCert.Title = &v
	}
	if v := spi.DSSASN1UtilsExtractAttributeFromX500Principal(oidCommonName, subject); v != "" {
		xmlCert.CommonName = &v
	}
	if v := spi.DSSASN1UtilsExtractAttributeFromX500Principal(oidLocality, subject); v != "" {
		xmlCert.Locality = &v
	}
	if v := spi.DSSASN1UtilsExtractAttributeFromX500Principal(oidState, subject); v != "" {
		xmlCert.State = &v
	}
	if v := spi.DSSASN1UtilsExtractAttributeFromX500Principal(oidCountryName, subject); v != "" {
		xmlCert.CountryName = &v
	}
	if v := spi.DSSASN1UtilsExtractAttributeFromX500Principal(oidOrganizationIdentifier, subject); v != "" {
		xmlCert.OrganizationIdentifier = &v
	}
	if v := spi.DSSASN1UtilsExtractAttributeFromX500Principal(oidOrganizationName, subject); v != "" {
		xmlCert.OrganizationName = &v
	}
	if v := spi.DSSASN1UtilsExtractAttributeFromX500Principal(oidOrganizationalUnit, subject); v != "" {
		xmlCert.OrganizationalUnit = &v
	}
	if v := spi.DSSASN1UtilsExtractAttributeFromX500Principal(oidGivenName, subject); v != "" {
		xmlCert.GivenName = &v
	}
	if v := spi.DSSASN1UtilsExtractAttributeFromX500Principal(oidSurname, subject); v != "" {
		xmlCert.Surname = &v
	}
	if v := spi.DSSASN1UtilsExtractAttributeFromX500Principal(oidPseudonym, subject); v != "" {
		xmlCert.Pseudonym = &v
	}
	if v := spi.DSSASN1UtilsExtractAttributeFromX500Principal(oidEmail, subject); v != "" {
		xmlCert.Email = &v
	}

	xmlCert.Sources = &jaxb.SourcesWrapper{Items: b.getXmlCertificateSources(certToken)}

	xmlCert.NotAfter = jaxb.NewXSDateTime(certToken.NotAfter())
	xmlCert.NotBefore = jaxb.NewXSDateTime(certToken.NotBefore())
	publicKey := certToken.PublicKey()
	xmlCert.PublicKeySize = spi.DSSPKUtilsPublicKeySize(publicKey)
	if encAlgo, err := enumerations.EncryptionAlgorithmForName(publicKey.Algorithm()); err == nil {
		v := jaxb.EncryptionAlgorithmValue(encAlgo)
		xmlCert.PublicKeyEncryptionAlgo = &v
	}
	entityKey := certToken.EntityKey().AsXmlID()
	xmlCert.EntityKey = &entityKey

	xmlCert.CertificateExtensions = &jaxb.CertificateExtensionsWrapper{Items: b.getXmlCertificateExtensions(certToken)}

	selfSigned := certToken.IsSelfSigned()
	xmlCert.SelfSigned = selfSigned
	xmlCert.Trusted = b.getXmlTrusted(certToken)

	if b.tokenExtractionStrategy.IsCertificate() {
		bin := jaxb.Base64Binary(certToken.Encoded())
		xmlCert.Base64Encoded = &bin
	} else {
		certDigest, err := certToken.Digest(b.defaultDigestAlgorithm)
		if err != nil {
			panic(err)
		}
		xmlCert.DigestAlgoAndValue = b.GetXmlDigestAlgoAndValueFor(b.defaultDigestAlgorithm, certDigest)
	}

	return xmlCert
}

func (b *DiagnosticDataBuilder) getXmlCertificateExtensions(token *model.CertificateToken) []jaxb.XmlCertificateExtensionItem {
	certificateExtensions, err := spi.CertificateExtensionsUtilsCertificateExtensions(token)
	if err != nil {
		panic(err)
	}

	xmlCertificateExtensions := make([]jaxb.XmlCertificateExtensionItem, 0)
	if certificateExtensions.AuthorityKeyIdentifier() != nil {
		xmlCertificateExtensions = append(xmlCertificateExtensions, b.getXmlAuthorityKeyIdentifier(certificateExtensions.AuthorityKeyIdentifier()))
	}
	if certificateExtensions.SubjectKeyIdentifier() != nil {
		xmlCertificateExtensions = append(xmlCertificateExtensions, b.getXmlSubjectKeyIdentifier(certificateExtensions.SubjectKeyIdentifier()))
	}
	if certificateExtensions.BasicConstraints() != nil {
		xmlCertificateExtensions = append(xmlCertificateExtensions, b.getXmlBasicConstraints(certificateExtensions.BasicConstraints()))
	}
	if certificateExtensions.KeyUsage() != nil {
		xmlCertificateExtensions = append(xmlCertificateExtensions, b.getXmlKeyUsages(certificateExtensions.KeyUsage()))
	}
	if certificateExtensions.CertificatePolicies() != nil {
		xmlCertificateExtensions = append(xmlCertificateExtensions, b.getXmlCertificatePolicies(certificateExtensions.CertificatePolicies()))
	}
	if certificateExtensions.SubjectAlternativeNames() != nil {
		xmlCertificateExtensions = append(xmlCertificateExtensions, b.getXmlSubjectAlternativeNames(certificateExtensions.SubjectAlternativeNames()))
	}
	if certificateExtensions.PolicyConstraints() != nil {
		xmlCertificateExtensions = append(xmlCertificateExtensions, b.getXmlPolicyConstraints(certificateExtensions.PolicyConstraints()))
	}
	if certificateExtensions.NameConstraints() != nil {
		xmlCertificateExtensions = append(xmlCertificateExtensions, b.getXmlNameConstraints(certificateExtensions.NameConstraints()))
	}
	if certificateExtensions.ExtendedKeyUsage() != nil {
		xmlCertificateExtensions = append(xmlCertificateExtensions, b.getXmlExtendedKeyUsages(certificateExtensions.ExtendedKeyUsage()))
	}
	if certificateExtensions.InhibitAnyPolicy() != nil {
		xmlCertificateExtensions = append(xmlCertificateExtensions, b.getXmlInhibitAnyPolicy(certificateExtensions.InhibitAnyPolicy()))
	}
	if certificateExtensions.AuthorityInformationAccess() != nil {
		xmlCertificateExtensions = append(xmlCertificateExtensions, b.getXmlAuthorityInformationAccess(certificateExtensions.AuthorityInformationAccess()))
	}
	if certificateExtensions.CRLDistributionPoints() != nil {
		xmlCertificateExtensions = append(xmlCertificateExtensions, b.getXmlCRLDistributionPoints(certificateExtensions.CRLDistributionPoints()))
	}
	if certificateExtensions.FreshestCRL() != nil {
		xmlCertificateExtensions = append(xmlCertificateExtensions, b.getXmlFreshestCRL(certificateExtensions.FreshestCRL()))
	}
	if certificateExtensions.OcspNoCheck() != nil {
		xmlCertificateExtensions = append(xmlCertificateExtensions, b.getXmlIdPkixOcspNoCheck(certificateExtensions.OcspNoCheck()))
	}
	if certificateExtensions.ValidityAssuredShortTerm() != nil {
		xmlCertificateExtensions = append(xmlCertificateExtensions, b.getXmlValAssuredShortTermCertificate(certificateExtensions.ValidityAssuredShortTerm()))
	}
	if certificateExtensions.NoRevAvail() != nil {
		xmlCertificateExtensions = append(xmlCertificateExtensions, b.getXmlNoRevAvail(certificateExtensions.NoRevAvail()))
	}
	if certificateExtensions.QcStatements() != nil {
		xmlCertificateExtensions = append(xmlCertificateExtensions, NewXmlQcStatementsBuilder().Build(certificateExtensions.QcStatements()))
	}
	if utils.IsCollectionNotEmpty(certificateExtensions.OtherExtensions()) {
		xmlCertificateExtensions = append(xmlCertificateExtensions, b.getXmlOtherCertificateExtensions(certificateExtensions.OtherExtensions())...)
	}

	return xmlCertificateExtensions
}

func (b *DiagnosticDataBuilder) getXmlKeyUsages(keyUsage *extension.KeyUsage) *jaxb.XmlKeyUsages {
	xmlKeyUsages := &jaxb.XmlKeyUsages{}
	b.fillXmlCertificateExtension(&xmlKeyUsages.XmlCertificateExtensionContent, &xmlKeyUsages.XmlCertificateExtensionAttrs, keyUsage)
	for _, bit := range keyUsage.KeyUsageBits() {
		xmlKeyUsages.KeyUsageBit = append(xmlKeyUsages.KeyUsageBit, jaxb.KeyUsageBitValue(bit))
	}
	return xmlKeyUsages
}

func (b *DiagnosticDataBuilder) getXmlExtendedKeyUsages(extendedKeyUsages *extension.ExtendedKeyUsages) *jaxb.XmlExtendedKeyUsages {
	xmlExtendedKeyUsages := &jaxb.XmlExtendedKeyUsages{}
	b.fillXmlCertificateExtension(&xmlExtendedKeyUsages.XmlCertificateExtensionContent, &xmlExtendedKeyUsages.XmlCertificateExtensionAttrs, extendedKeyUsages)
	xmlExtendedKeyUsages.ExtendedKeyUsageOid = b.getXmlOids(extendedKeyUsages.Oids())
	return xmlExtendedKeyUsages
}

func (b *DiagnosticDataBuilder) getXmlCertificatePolicies(certificatePolicies *extension.CertificatePolicies) *jaxb.XmlCertificatePolicies {
	xmlCertificatePolicies := &jaxb.XmlCertificatePolicies{}
	b.fillXmlCertificateExtension(&xmlCertificatePolicies.XmlCertificateExtensionContent, &xmlCertificatePolicies.XmlCertificateExtensionAttrs, certificatePolicies)
	xmlCertificatePolicies.CertificatePolicy = b.getXmlCertificatePolicyList(certificatePolicies.PolicyList())
	return xmlCertificatePolicies
}

func (b *DiagnosticDataBuilder) getXmlCertificatePolicyList(certificatePolicies []*extension.CertificatePolicy) []*jaxb.XmlCertificatePolicy {
	result := make([]*jaxb.XmlCertificatePolicy, 0, len(certificatePolicies))
	for _, cp := range certificatePolicies {
		xmlCP := &jaxb.XmlCertificatePolicy{}
		oid := cp.Oid()
		xmlCP.Value = oid
		if desc := model.OidRepositoryGetDescription(oid); desc != "" {
			xmlCP.Description = &desc
		}
		if cpsURL := b.getCleanedUrl(cp.CpsUrl()); cpsURL != "" {
			xmlCP.CpsUrl = &cpsURL
		}
		result = append(result, xmlCP)
	}
	return result
}

func (b *DiagnosticDataBuilder) getXmlSubjectAlternativeNames(subjectAlternativeNames *extension.SubjectAlternativeNames) *jaxb.XmlSubjectAlternativeNames {
	xmlSubjectAlternativeNames := &jaxb.XmlSubjectAlternativeNames{}
	b.fillXmlCertificateExtension(&xmlSubjectAlternativeNames.XmlCertificateExtensionContent, &xmlSubjectAlternativeNames.XmlCertificateExtensionAttrs, subjectAlternativeNames)
	xmlSubjectAlternativeNames.SubjectAlternativeName = b.getXmlGeneralNames(subjectAlternativeNames.GeneralNames())
	return xmlSubjectAlternativeNames
}

func (b *DiagnosticDataBuilder) getXmlGeneralNames(generalNames []*extension.GeneralName) []*jaxb.XmlGeneralName {
	result := make([]*jaxb.XmlGeneralName, 0, len(generalNames))
	for _, generalName := range generalNames {
		result = append(result, b.getXmlGeneralName(generalName))
	}
	return result
}

func (b *DiagnosticDataBuilder) getXmlGeneralName(generalName *extension.GeneralName) *jaxb.XmlGeneralName {
	xmlGeneralName := &jaxb.XmlGeneralName{}
	t := jaxb.GeneralNameTypeValue(generalName.GeneralNameType())
	xmlGeneralName.Type = &t
	xmlGeneralName.Value = generalName.Value()
	return xmlGeneralName
}

func (b *DiagnosticDataBuilder) getXmlBasicConstraints(basicConstraints *extension.BasicConstraints) *jaxb.XmlBasicConstraints {
	xmlBasicConstraints := &jaxb.XmlBasicConstraints{}
	b.fillXmlCertificateExtension(&xmlBasicConstraints.XmlCertificateExtensionContent, &xmlBasicConstraints.XmlCertificateExtensionAttrs, basicConstraints)
	xmlBasicConstraints.CA = basicConstraints.IsCa()
	if basicConstraints.PathLenConstraint() != -1 {
		v := basicConstraints.PathLenConstraint()
		xmlBasicConstraints.PathLenConstraint = &v
	}
	return xmlBasicConstraints
}

func (b *DiagnosticDataBuilder) getXmlPolicyConstraints(policyConstraints *extension.PolicyConstraints) *jaxb.XmlPolicyConstraints {
	xmlPolicyConstraints := &jaxb.XmlPolicyConstraints{}
	b.fillXmlCertificateExtension(&xmlPolicyConstraints.XmlCertificateExtensionContent, &xmlPolicyConstraints.XmlCertificateExtensionAttrs, policyConstraints)
	if policyConstraints.InhibitPolicyMapping() != -1 {
		v := policyConstraints.InhibitPolicyMapping()
		xmlPolicyConstraints.InhibitPolicyMapping = &v
	}
	if policyConstraints.RequireExplicitPolicy() != -1 {
		v := policyConstraints.RequireExplicitPolicy()
		xmlPolicyConstraints.RequireExplicitPolicy = &v
	}
	return xmlPolicyConstraints
}

func (b *DiagnosticDataBuilder) getXmlInhibitAnyPolicy(inhibitAnyPolicy *extension.InhibitAnyPolicy) *jaxb.XmlInhibitAnyPolicy {
	xmlInhibitAnyPolicy := &jaxb.XmlInhibitAnyPolicy{}
	b.fillXmlCertificateExtension(&xmlInhibitAnyPolicy.XmlCertificateExtensionContent, &xmlInhibitAnyPolicy.XmlCertificateExtensionAttrs, inhibitAnyPolicy)
	if inhibitAnyPolicy.Value() != -1 {
		v := inhibitAnyPolicy.Value()
		xmlInhibitAnyPolicy.Value = &v
	}
	return xmlInhibitAnyPolicy
}

func (b *DiagnosticDataBuilder) getXmlNameConstraints(nameConstraints *extension.NameConstraints) *jaxb.XmlNameConstraints {
	xmlNameConstraints := &jaxb.XmlNameConstraints{}
	b.fillXmlCertificateExtension(&xmlNameConstraints.XmlCertificateExtensionContent, &xmlNameConstraints.XmlCertificateExtensionAttrs, nameConstraints)
	if utils.IsCollectionNotEmpty(nameConstraints.PermittedSubtrees()) {
		xmlNameConstraints.PermittedSubtree = b.getXmlGeneralSubtrees(nameConstraints.PermittedSubtrees())
	}
	if utils.IsCollectionNotEmpty(nameConstraints.ExcludedSubtrees()) {
		xmlNameConstraints.ExcludedSubtree = b.getXmlGeneralSubtrees(nameConstraints.ExcludedSubtrees())
	}
	return xmlNameConstraints
}

func (b *DiagnosticDataBuilder) getXmlGeneralSubtrees(generalSubtrees []*extension.GeneralSubtree) []*jaxb.XmlGeneralSubtree {
	result := make([]*jaxb.XmlGeneralSubtree, 0, len(generalSubtrees))
	for _, generalSubtree := range generalSubtrees {
		result = append(result, b.getXmlGeneralSubtree(generalSubtree))
	}
	return result
}

func (b *DiagnosticDataBuilder) getXmlGeneralSubtree(generalSubtree *extension.GeneralSubtree) *jaxb.XmlGeneralSubtree {
	xmlGeneralSubtree := &jaxb.XmlGeneralSubtree{}
	t := jaxb.GeneralNameTypeValue(generalSubtree.GeneralNameType())
	xmlGeneralSubtree.Type = &t
	xmlGeneralSubtree.Value = generalSubtree.Value()
	xmlGeneralSubtree.Minimum = jaxb.NewBigInteger(generalSubtree.Minimum())
	xmlGeneralSubtree.Maximum = jaxb.NewBigInteger(generalSubtree.Maximum())
	return xmlGeneralSubtree
}

func (b *DiagnosticDataBuilder) getXmlCRLDistributionPoints(crlDistributionPoints *extension.CRLDistributionPoints) *jaxb.XmlCRLDistributionPoints {
	xmlCRLDistributionPoints := &jaxb.XmlCRLDistributionPoints{}
	b.fillXmlCertificateExtension(&xmlCRLDistributionPoints.XmlCertificateExtensionContent, &xmlCRLDistributionPoints.XmlCertificateExtensionAttrs, crlDistributionPoints)
	xmlCRLDistributionPoints.CrlUrl = b.getCleanedUrls(crlDistributionPoints.CrlUrls())
	return xmlCRLDistributionPoints
}

func (b *DiagnosticDataBuilder) getXmlFreshestCRL(freshestCRL *extension.FreshestCRL) *jaxb.XmlFreshestCRL {
	xmlFreshestCRL := &jaxb.XmlFreshestCRL{}
	b.fillXmlCertificateExtension(&xmlFreshestCRL.XmlCertificateExtensionContent, &xmlFreshestCRL.XmlCertificateExtensionAttrs, freshestCRL)
	xmlFreshestCRL.CrlUrl = b.getCleanedUrls(freshestCRL.CrlUrls())
	return xmlFreshestCRL
}

func (b *DiagnosticDataBuilder) getXmlAuthorityKeyIdentifier(aki *extension.AuthorityKeyIdentifier) *jaxb.XmlAuthorityKeyIdentifier {
	xmlAuthorityKeyIdentifier := &jaxb.XmlAuthorityKeyIdentifier{}
	b.fillXmlCertificateExtension(&xmlAuthorityKeyIdentifier.XmlCertificateExtensionContent, &xmlAuthorityKeyIdentifier.XmlCertificateExtensionAttrs, aki)
	if aki.KeyIdentifier() != nil {
		v := jaxb.Base64Binary(aki.KeyIdentifier())
		xmlAuthorityKeyIdentifier.KeyIdentifier = &v
	}
	if aki.AuthorityCertIssuerSerial() != nil {
		v := jaxb.Base64Binary(aki.AuthorityCertIssuerSerial())
		xmlAuthorityKeyIdentifier.AuthorityCertIssuerSerial = &v
	}
	return xmlAuthorityKeyIdentifier
}

func (b *DiagnosticDataBuilder) getXmlSubjectKeyIdentifier(ski *extension.SubjectKeyIdentifier) *jaxb.XmlSubjectKeyIdentifier {
	xmlSubjectKeyIdentifier := &jaxb.XmlSubjectKeyIdentifier{}
	b.fillXmlCertificateExtension(&xmlSubjectKeyIdentifier.XmlCertificateExtensionContent, &xmlSubjectKeyIdentifier.XmlCertificateExtensionAttrs, ski)
	if ski.Ski() != nil {
		v := jaxb.Base64Binary(ski.Ski())
		xmlSubjectKeyIdentifier.Ski = &v
	}
	return xmlSubjectKeyIdentifier
}

func (b *DiagnosticDataBuilder) getXmlAuthorityInformationAccess(aia *extension.AuthorityInformationAccess) *jaxb.XmlAuthorityInformationAccess {
	xmlAuthorityInformationAccess := &jaxb.XmlAuthorityInformationAccess{}
	b.fillXmlCertificateExtension(&xmlAuthorityInformationAccess.XmlCertificateExtensionContent, &xmlAuthorityInformationAccess.XmlCertificateExtensionAttrs, aia)
	xmlAuthorityInformationAccess.CaIssuersUrl = b.getCleanedUrls(aia.CaIssuers())
	xmlAuthorityInformationAccess.OcspUrl = b.getCleanedUrls(aia.Ocsp())
	return xmlAuthorityInformationAccess
}

func (b *DiagnosticDataBuilder) getXmlIdPkixOcspNoCheck(ocspNoCheck *extension.OCSPNoCheck) *jaxb.XmlIdPkixOcspNoCheck {
	xmlIdPkixOcspNoCheck := &jaxb.XmlIdPkixOcspNoCheck{}
	b.fillXmlCertificateExtension(&xmlIdPkixOcspNoCheck.XmlCertificateExtensionContent, &xmlIdPkixOcspNoCheck.XmlCertificateExtensionAttrs, ocspNoCheck)
	present := ocspNoCheck.IsOcspNoCheck()
	xmlIdPkixOcspNoCheck.Present = &present
	return xmlIdPkixOcspNoCheck
}

func (b *DiagnosticDataBuilder) getXmlValAssuredShortTermCertificate(valAssuredST *extension.ValidityAssuredShortTerm) *jaxb.XmlValAssuredShortTermCertificate {
	xmlValAssuredShortTermCertificate := &jaxb.XmlValAssuredShortTermCertificate{}
	b.fillXmlCertificateExtension(&xmlValAssuredShortTermCertificate.XmlCertificateExtensionContent, &xmlValAssuredShortTermCertificate.XmlCertificateExtensionAttrs, valAssuredST)
	present := valAssuredST.IsValAssuredSTCerts()
	xmlValAssuredShortTermCertificate.Present = &present
	return xmlValAssuredShortTermCertificate
}

func (b *DiagnosticDataBuilder) getXmlNoRevAvail(noRevAvail *extension.NoRevAvail) *jaxb.XmlNoRevAvail {
	xmlNoRevAvail := &jaxb.XmlNoRevAvail{}
	b.fillXmlCertificateExtension(&xmlNoRevAvail.XmlCertificateExtensionContent, &xmlNoRevAvail.XmlCertificateExtensionAttrs, noRevAvail)
	present := noRevAvail.IsNoRevAvail()
	xmlNoRevAvail.Present = &present
	return xmlNoRevAvail
}

func (b *DiagnosticDataBuilder) getXmlOtherCertificateExtensions(otherCertificateExtensions []*extension.CertificateExtension) []jaxb.XmlCertificateExtensionItem {
	result := make([]jaxb.XmlCertificateExtensionItem, 0, len(otherCertificateExtensions))
	for _, certificateExtension := range otherCertificateExtensions {
		xmlCertificateExtension := &jaxb.XmlCertificateExtension{}
		b.fillXmlCertificateExtension(&xmlCertificateExtension.XmlCertificateExtensionContent, &xmlCertificateExtension.XmlCertificateExtensionAttrs, certificateExtension)
		if certificateExtension.Octets() != nil {
			v := jaxb.Base64Binary(certificateExtension.Octets())
			xmlCertificateExtension.Octets = &v
		}
		result = append(result, xmlCertificateExtension)
	}
	return result
}

// certificateExtensionLike is the common surface of every model/x509/extension.* certificate
// extension type this file fills an XmlCertificateExtension(Content|Attrs) pair from.
type certificateExtensionLike interface {
	OID() string
	Description() string
	IsCritical() bool
}

func (b *DiagnosticDataBuilder) fillXmlCertificateExtension(content *jaxb.XmlCertificateExtensionContent, attrs *jaxb.XmlCertificateExtensionAttrs, certificateExtension certificateExtensionLike) {
	oid := certificateExtension.OID()
	attrs.OID = &oid
	// Java: xmlCertificateExtension.setDescription(certificateExtension.getDescription()),
	// whose null leaves the attribute absent. The generated member is a pointer, so an
	// empty description must map to nil - assigning &"" would emit description="".
	if description := certificateExtension.Description(); description != "" {
		attrs.Description = &description
	} else {
		attrs.Description = nil
	}
	critical := certificateExtension.IsCritical()
	attrs.Critical = &critical
}

func (b *DiagnosticDataBuilder) getXmlTrusted(certificateToken *model.CertificateToken) *jaxb.XmlTrusted {
	xmlTrusted := &jaxb.XmlTrusted{}
	if b.allCertificateSources.IsTrusted(certificateToken) {
		certificateTrustTime := b.getCertificateTrustTime(certificateToken)
		if certificateTrustTime != nil {
			xmlTrusted.Value = certificateTrustTime.IsTrusted()
			if !certificateTrustTime.StartDate().IsZero() {
				xmlTrusted.StartDate = jaxb.NewXSDateTime(certificateTrustTime.StartDate())
			}
			if !certificateTrustTime.EndDate().IsZero() {
				xmlTrusted.SunsetDate = jaxb.NewXSDateTime(certificateTrustTime.EndDate())
			}
		} else {
			xmlTrusted.Value = true
		}
	}
	return xmlTrusted
}

func (b *DiagnosticDataBuilder) getCertificateTrustTime(certificateToken *model.CertificateToken) *tsl.CertificateTrustTime {
	var certificateTrustTime *tsl.CertificateTrustTime
	for _, trustedSource := range b.allCertificateSources.Sources() {
		if !trustedSource.IsTrusted(certificateToken) {
			continue
		}
		if trustedSourceWithTime, ok := trustedSource.(tsl.TrustedCertificateSourceWithTime); ok {
			trustTime := trustedSourceWithTime.TrustTime(certificateToken)
			if certificateTrustTime == nil || !certificateTrustTime.IsTrusted() {
				certificateTrustTime = trustTime
			} else if trustTime != nil && trustTime.IsTrusted() {
				certificateTrustTime = certificateTrustTime.JointTrustTime(trustTime.StartDate(), trustTime.EndDate())
			}
		} else {
			certificateTrustTime = tsl.NewCertificateTrustTime(true)
		}
	}
	return certificateTrustTime
}

func (b *DiagnosticDataBuilder) getXmlCertificateSources(token *model.CertificateToken) []jaxb.CertificateSourceTypeValue {
	certificateSources := make([]enumerations.CertificateSourceType, 0)
	if b.allCertificateSources != nil {
		sourceTypes := b.allCertificateSources.CertificateSourceTypeOf(token)
		for sourceType := range sourceTypes {
			certificateSources = append(certificateSources, sourceType)
		}
		sort.Slice(certificateSources, func(i, j int) bool { return certificateSources[i] < certificateSources[j] })
	}
	if utils.IsCollectionEmpty(certificateSources) {
		certificateSources = append(certificateSources, enumerations.CertificateSourceType_UNKNOWN)
	}
	result := make([]jaxb.CertificateSourceTypeValue, 0, len(certificateSources))
	for _, cs := range certificateSources {
		result = append(result, jaxb.CertificateSourceTypeValue(cs))
	}
	return result
}

func (b *DiagnosticDataBuilder) getRevocationsForCert(certToken *model.CertificateToken) []validation.AnyRevocationToken {
	revocations := make([]validation.AnyRevocationToken, 0)
	if utils.IsCollectionNotEmpty(b.usedRevocations) {
		for _, revocationToken := range b.usedRevocations {
			if certToken.DSSIDAsString() == revocationToken.RelatedCertificateID() {
				revocations = append(revocations, revocationToken)
			}
		}
	}
	return revocations
}

func (b *DiagnosticDataBuilder) getXmlOids(oidList []string) []*jaxb.XmlOID {
	result := make([]*jaxb.XmlOID, 0)
	if utils.IsCollectionNotEmpty(oidList) {
		for _, oid := range oidList {
			xmlOID := &jaxb.XmlOID{}
			xmlOID.Value = oid
			if desc := model.OidRepositoryGetDescription(oid); desc != "" {
				xmlOID.Description = &desc
			}
			result = append(result, xmlOID)
		}
	}
	return result
}

// GetXmlDigestAlgoAndValueForDigestValue builds an XmlDigestAlgoAndValue for a Digest. Port of
// the protected getXmlDigestAlgoAndValue(Digest); model.Digest is a Go value type (its zero
// value stands for Java's null Digest), unlike the *model.CertificateExtension-family pointer
// types this file otherwise uses nil for.
func (b *DiagnosticDataBuilder) GetXmlDigestAlgoAndValueForDigestValue(digest model.Digest) *jaxb.XmlDigestAlgoAndValue {
	if digest.IsEmpty() {
		return b.GetXmlDigestAlgoAndValueFor("", nil)
	}
	return b.GetXmlDigestAlgoAndValueFor(digest.Algorithm(), digest.Value())
}

// GetXmlDigestAlgoAndValueFor builds an XmlDigestAlgoAndValue for a DigestAlgorithm and
// digestValue. Port of the protected getXmlDigestAlgoAndValue(DigestAlgorithm, byte[]).
func (b *DiagnosticDataBuilder) GetXmlDigestAlgoAndValueFor(digestAlgo enumerations.DigestAlgorithm, digestValue []byte) *jaxb.XmlDigestAlgoAndValue {
	xmlDigestAlgAndValue := &jaxb.XmlDigestAlgoAndValue{}
	if digestAlgo != "" {
		v := jaxb.DigestAlgorithmValue(digestAlgo)
		xmlDigestAlgAndValue.DigestMethod = &v
	}
	if digestValue == nil {
		digestValue = spi.DSSUtilsEmptyByteArray
	}
	bin := jaxb.Base64Binary(digestValue)
	xmlDigestAlgAndValue.DigestValue = &bin
	return xmlDigestAlgAndValue
}

// certificateChainWrapper wraps a possibly-nil chain slice into the jaxb wrapper type, matching
// Java's null List (no <CertificateChain/> element) vs a populated one.
func (b *DiagnosticDataBuilder) certificateChainWrapper(chain []*jaxb.XmlChainItem) *jaxb.CertificateChainWrapper {
	if chain == nil {
		return nil
	}
	return &jaxb.CertificateChainWrapper{Items: chain}
}

// X.500 principal format identifiers, matching javax.security.auth.x500.X500Principal's
// CANONICAL/RFC2253 constants.
const (
	x500PrincipalCanonical = "CANONICAL"
	x500PrincipalRFC2253   = "RFC2253"
)

var (
	oidSubjectSerialNumber    = mustBCStyleOID("2.5.4.5")
	oidTitle                  = mustBCStyleOID("2.5.4.12")
	oidCommonName             = mustBCStyleOID("2.5.4.3")
	oidLocality               = mustBCStyleOID("2.5.4.7")
	oidState                  = mustBCStyleOID("2.5.4.8")
	oidCountryName            = mustBCStyleOID("2.5.4.6")
	oidOrganizationIdentifier = mustBCStyleOID("2.5.4.97")
	oidOrganizationName       = mustBCStyleOID("2.5.4.10")
	oidOrganizationalUnit     = mustBCStyleOID("2.5.4.11")
	oidGivenName              = mustBCStyleOID("2.5.4.42")
	oidSurname                = mustBCStyleOID("2.5.4.4")
	oidPseudonym              = mustBCStyleOID("2.5.4.65")
	oidEmail                  = mustBCStyleOID("1.2.840.113549.1.9.1")
)

func mustBCStyleOID(dotted string) []int {
	parts := make([]int, 0)
	cur := 0
	has := false
	for _, r := range dotted {
		if r == '.' {
			parts = append(parts, cur)
			cur = 0
			has = false
			continue
		}
		cur = cur*10 + int(r-'0')
		has = true
	}
	if has {
		parts = append(parts, cur)
	}
	return parts
}

var _ = fmt.Sprintf
var _ job.InfoRecord
