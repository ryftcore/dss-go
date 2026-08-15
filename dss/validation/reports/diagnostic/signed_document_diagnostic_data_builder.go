// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/reports/diagnostic/SignedDocumentDiagnosticDataBuilder.java (DSS 6.5.RC1).
//
// # Virtual dispatch
//
// BuildDetachedXmlSignature, BuildDetachedXmlTimestamp and AssertConfigurationValid are called
// back into from Build() (Java: buildXmlSignatures()/getXmlSignature(),
// buildXmlTimestamps(), and the top of build() respectively) and are overridden by
// QWACCertificateDiagnosticDataBuilder (this manifest, BuildDetachedXmlSignature and
// AssertConfigurationValid) and by the forward-declared CAdES/PAdES/JAdES/ASiC
// DiagnosticDataBuilders (out of this manifest, BuildDetachedXmlSignature and
// BuildDetachedXmlTimestamp) - so all three are collected into
// SignedDocumentDiagnosticDataBuilderOverrides, following the same Init<TypeName> pattern as
// DiagnosticDataBuilderOverrides.
//
// # Revocation-source wildcard types
//
// Java's OfflineRevocationSource<R> parameter is a raw abstract-class reference shared by
// AdvancedSignature.getCRLSource()/getOCSPSource() and EvidenceRecord.getCRLSource()/
// getOCSPSource(). The already-landed Go ports type these two call sites differently -
// AdvancedSignature returns the narrow spi.OfflineRevocationSource[R] interface,
// EvidenceRecord returns the concrete *spi.OfflineCRLSourceBase/*spi.OfflineOCSPSourceBase -
// and neither exposes the full spi.OfflineRevocationSourceBase[R] method set this file needs
// (UniqueRevocationTokensWithOrigins, FindRefsAndOriginsFor*, OrphanRevocationReferencesWithOrigins,
// the per-origin binary getters). Every concrete revocation source in this codebase embeds
// OfflineRevocationSourceBase[R] regardless, so offlineRevocationSourceRefs below is a local,
// wider interface every call site type-asserts its narrower static type down to.
//
// Java's OfflineRevocationSource.getAllRevocationBinariesWithOrigins() (Map<..., Set<
// RevocationOrigin>>) has no direct Go counterpart: the port only exposes per-origin binary
// getters for the 8 origins collectable from a signature/evidence-record's own structure
// (CMS_SIGNED_DATA, REVOCATION_VALUES, ATTRIBUTE_REVOCATION_VALUES, TIMESTAMP_VALIDATION_DATA,
// ANY_VALIDATION_DATA, DSS_DICTIONARY, VRI_DICTIONARY, ADBE_REVOCATION_INFO_ARCHIVAL), not the
// online/document origins (INPUT_DOCUMENT, EXTERNAL, CACHED, EVIDENCE_RECORD). allBinariesWithOrigins
// below reconstructs the map from those 8 getters; it is complete for the embedded-in-signature
// case this builder's only caller (getXmlOrphanRevocations, over a signature/timestamp/evidence
// record's OWN CRL/OCSP source) exercises, but is flagged here since it cannot reproduce an
// online-fetched orphan binary's origin. See the porter notes for confirmation against the
// RPTDIAG byte-compare oracle.
package diagnostic

import (
	"time"

	"github.com/utain/esig/dss/crlparser"
	"github.com/utain/esig/dss/diagnostic/jaxb"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/model/scope"
	"github.com/utain/esig/dss/model/signature"
	"github.com/utain/esig/dss/model/x509/revocation"
	"github.com/utain/esig/dss/spi"
	"github.com/utain/esig/dss/spi/validation"
	"github.com/utain/esig/dss/utils"
)

// SignedDocumentDiagnosticDataBuilderOverrides declares the operations
// SignedDocumentDiagnosticDataBuilder calls back into that at least one subclass in this port
// (or its forward-declared out-of-manifest CAdES/PAdES/JAdES/ASiC consumers) overrides.
type SignedDocumentDiagnosticDataBuilderOverrides interface {
	// BuildDetachedXmlSignature builds the XmlSignature. Port of the public
	// buildDetachedXmlSignature(AdvancedSignature).
	BuildDetachedXmlSignature(signature validation.AdvancedSignature) *jaxb.XmlSignature
	// BuildDetachedXmlTimestamp builds the XmlTimestamp from a TimestampToken. Port of the
	// protected buildDetachedXmlTimestamp(TimestampToken).
	BuildDetachedXmlTimestamp(timestampToken *validation.TimestampToken) *jaxb.XmlTimestamp
	// AssertConfigurationValid verifies the configuration is valid to build a DiagnosticData.
	// Port of the protected assertConfigurationValid().
	AssertConfigurationValid()
}

// SignedDocumentDiagnosticDataBuilder is the common builder for DiagnosticData creation from a
// signed/timestamped document.
type SignedDocumentDiagnosticDataBuilder struct {
	DiagnosticDataBuilder

	// signedDocument is the signed document.
	signedDocument model.DSSDocument

	// signatures is the collection of signatures.
	signatures []validation.AdvancedSignature

	// usedTimestamps is the collection of timestamp tokens.
	usedTimestamps []*validation.TimestampToken

	// evidenceRecords is the collection of evidence records.
	evidenceRecords []validation.EvidenceRecord

	// documentCertificateSource is the list of all certificate sources extracted from a
	// validating document (signature(s), timestamp(s)).
	documentCertificateSource *spi.ListCertificateSource

	// documentCRLSource is the list of all CRL revocation sources extracted from a validating
	// document (signature(s), timestamp(s)).
	documentCRLSource *spi.ListRevocationSource[revocation.CRL]

	// documentOCSPSource is the list of all OCSP revocation sources extracted from a validating
	// document (signature(s), timestamp(s)).
	documentOCSPSource *spi.ListRevocationSource[revocation.OCSP]

	// xmlSignaturesMap is the cached map of signatures.
	xmlSignaturesMap map[string]*jaxb.XmlSignature

	// xmlTimestampsMap is the cached map of timestamps.
	xmlTimestampsMap map[string]*jaxb.XmlTimestamp

	// xmlEvidenceRecordMap is the cached map of evidence records.
	xmlEvidenceRecordMap map[string]*jaxb.XmlEvidenceRecord

	// xmlSignedDataMap is the cached map of original signed data.
	xmlSignedDataMap map[string]*jaxb.XmlSignerData

	// overrides points back at the concrete/intermediate builder; see
	// InitSignedDocumentDiagnosticDataBuilder.
	overrides SignedDocumentDiagnosticDataBuilderOverrides
}

// NewSignedDocumentDiagnosticDataBuilder instantiates the object with null values and empty
// maps. Port of the public default constructor.
func NewSignedDocumentDiagnosticDataBuilder() *SignedDocumentDiagnosticDataBuilder {
	b := &SignedDocumentDiagnosticDataBuilder{
		DiagnosticDataBuilder:     *NewDiagnosticDataBuilder(),
		documentCertificateSource: spi.NewListCertificateSource(),
		documentCRLSource:         spi.NewListRevocationSource[revocation.CRL](),
		documentOCSPSource:        spi.NewListRevocationSource[revocation.OCSP](),
		xmlSignaturesMap:          map[string]*jaxb.XmlSignature{},
		xmlTimestampsMap:          map[string]*jaxb.XmlTimestamp{},
		xmlEvidenceRecordMap:      map[string]*jaxb.XmlEvidenceRecord{},
		xmlSignedDataMap:          map[string]*jaxb.XmlSignerData{},
	}
	b.InitSignedDocumentDiagnosticDataBuilder(b)
	return b
}

// InitSignedDocumentDiagnosticDataBuilder registers the concrete/intermediate builder with the
// base so it can dispatch to SignedDocumentDiagnosticDataBuilderOverrides, and re-registers it
// as the DiagnosticDataBuilderOverrides (LinkSigningCertificateAndChains) since this type
// overrides that method too. Every constructor for a type in this family must call this once.
func (b *SignedDocumentDiagnosticDataBuilder) InitSignedDocumentDiagnosticDataBuilder(overrides SignedDocumentDiagnosticDataBuilderOverrides) {
	b.overrides = overrides
	b.DiagnosticDataBuilder.InitDiagnosticDataBuilder(b)
}

func (b *SignedDocumentDiagnosticDataBuilder) signedDocumentDiagnosticDataBuilderOverrides() SignedDocumentDiagnosticDataBuilderOverrides {
	if b.overrides == nil {
		panic("SignedDocumentDiagnosticDataBuilder was not initialised: the concrete builder must call InitSignedDocumentDiagnosticDataBuilder in its constructor")
	}
	return b.overrides
}

// LinkSigningCertificateAndChains overrides DiagnosticDataBuilder's default (the certificate
// chain is built based on provided tokens instead). Port of the protected
// @Override linkSigningCertificateAndChains(Set<CertificateToken>).
func (b *SignedDocumentDiagnosticDataBuilder) LinkSigningCertificateAndChains(certificates []*model.CertificateToken) {
	// skip (certificate chain is build based on provided tokens)
}

// UsedCertificates re-declares the fluent setter with the SignedDocumentDiagnosticDataBuilder
// return type. Port of the covariant-return @Override usedCertificates(Set<CertificateToken>).
func (b *SignedDocumentDiagnosticDataBuilder) UsedCertificates(usedCertificates []*model.CertificateToken) *SignedDocumentDiagnosticDataBuilder {
	b.DiagnosticDataBuilder.UsedCertificates(usedCertificates)
	return b
}

// UsedRevocations re-declares the fluent setter with the SignedDocumentDiagnosticDataBuilder
// return type. Port of the covariant-return @Override usedRevocations(Set<RevocationToken<?>>).
func (b *SignedDocumentDiagnosticDataBuilder) UsedRevocations(usedRevocations []validation.AnyRevocationToken) *SignedDocumentDiagnosticDataBuilder {
	b.DiagnosticDataBuilder.UsedRevocations(usedRevocations)
	return b
}

// AllCertificateSources re-declares the fluent setter with the SignedDocumentDiagnosticDataBuilder
// return type. Port of the covariant-return @Override allCertificateSources(ListCertificateSource).
func (b *SignedDocumentDiagnosticDataBuilder) AllCertificateSources(allCertificateSources *spi.ListCertificateSource) *SignedDocumentDiagnosticDataBuilder {
	b.DiagnosticDataBuilder.AllCertificateSources(allCertificateSources)
	return b
}

// ValidationDate re-declares the fluent setter with the SignedDocumentDiagnosticDataBuilder
// return type. Port of the covariant-return @Override validationDate(Date).
func (b *SignedDocumentDiagnosticDataBuilder) ValidationDate(validationDate time.Time) *SignedDocumentDiagnosticDataBuilder {
	b.DiagnosticDataBuilder.ValidationDate(validationDate)
	return b
}

// TokenExtractionStrategy re-declares the fluent setter with the
// SignedDocumentDiagnosticDataBuilder return type. Port of the covariant-return @Override
// tokenExtractionStrategy(TokenExtractionStrategy).
func (b *SignedDocumentDiagnosticDataBuilder) TokenExtractionStrategy(tokenExtractionStrategy enumerations.TokenExtractionStrategy) *SignedDocumentDiagnosticDataBuilder {
	b.DiagnosticDataBuilder.TokenExtractionStrategy(tokenExtractionStrategy)
	return b
}

// DefaultDigestAlgorithm re-declares the fluent setter with the SignedDocumentDiagnosticDataBuilder
// return type. Port of the covariant-return @Override defaultDigestAlgorithm(DigestAlgorithm).
func (b *SignedDocumentDiagnosticDataBuilder) DefaultDigestAlgorithm(digestAlgorithm enumerations.DigestAlgorithm) *SignedDocumentDiagnosticDataBuilder {
	b.DiagnosticDataBuilder.DefaultDigestAlgorithm(digestAlgorithm)
	return b
}

// TokenIdentifierProvider re-declares the fluent setter (not overridden in Java, but exposed
// here at the concrete-return level so XmlDiagnosticDataFactory's chain compiles cleanly; see
// PORTING.md's covariant-return precedent).
func (b *SignedDocumentDiagnosticDataBuilder) TokenIdentifierProvider(identifierProvider model.TokenIdentifierProvider) *SignedDocumentDiagnosticDataBuilder {
	b.DiagnosticDataBuilder.TokenIdentifierProvider(identifierProvider)
	return b
}

// Document sets the document which is analysed. Port of document(DSSDocument).
func (b *SignedDocumentDiagnosticDataBuilder) Document(signedDocument model.DSSDocument) *SignedDocumentDiagnosticDataBuilder {
	b.signedDocument = signedDocument
	return b
}

// FoundSignatures sets the found signatures. Port of foundSignatures(Collection<AdvancedSignature>).
func (b *SignedDocumentDiagnosticDataBuilder) FoundSignatures(signatures []validation.AdvancedSignature) *SignedDocumentDiagnosticDataBuilder {
	b.signatures = signatures
	return b
}

// UsedTimestamps sets the timestamps. Port of usedTimestamps(Collection<TimestampToken>).
func (b *SignedDocumentDiagnosticDataBuilder) UsedTimestamps(usedTimestamps []*validation.TimestampToken) *SignedDocumentDiagnosticDataBuilder {
	b.usedTimestamps = usedTimestamps
	return b
}

// FoundEvidenceRecords sets the evidence records. Port of
// foundEvidenceRecords(Collection<EvidenceRecord>).
func (b *SignedDocumentDiagnosticDataBuilder) FoundEvidenceRecords(evidenceRecords []validation.EvidenceRecord) *SignedDocumentDiagnosticDataBuilder {
	b.evidenceRecords = evidenceRecords
	return b
}

// DocumentCertificateSource sets a document Certificate Source containing all sources extracted
// from the provided signature(s)/timestamp(s). Port of
// documentCertificateSource(ListCertificateSource).
func (b *SignedDocumentDiagnosticDataBuilder) DocumentCertificateSource(documentCertificateSource *spi.ListCertificateSource) *SignedDocumentDiagnosticDataBuilder {
	b.documentCertificateSource = documentCertificateSource
	return b
}

// DocumentCRLSource sets a document CRL Source containing all sources extracted from the
// provided signature(s)/timestamp(s). Port of documentCRLSource(ListRevocationSource<CRL>).
func (b *SignedDocumentDiagnosticDataBuilder) DocumentCRLSource(documentCRLSource *spi.ListRevocationSource[revocation.CRL]) *SignedDocumentDiagnosticDataBuilder {
	b.documentCRLSource = documentCRLSource
	return b
}

// DocumentOCSPSource sets a document OCSP Source containing all sources extracted from the
// provided signature(s)/timestamp(s). Port of documentOCSPSource(ListRevocationSource<OCSP>).
func (b *SignedDocumentDiagnosticDataBuilder) DocumentOCSPSource(documentOCSPSource *spi.ListRevocationSource[revocation.OCSP]) *SignedDocumentDiagnosticDataBuilder {
	b.documentOCSPSource = documentOCSPSource
	return b
}

// Build builds the XmlDiagnosticData. Port of the public @Override build().
func (b *SignedDocumentDiagnosticDataBuilder) Build() *jaxb.XmlDiagnosticData {
	overrides := b.signedDocumentDiagnosticDataBuilderOverrides()
	overrides.AssertConfigurationValid()

	diagnosticData := b.DiagnosticDataBuilder.Build() // fill certificates and revocation data
	if b.signedDocument != nil {
		documentName := b.removeSpecialCharsForXml(b.signedDocument.Name())
		diagnosticData.DocumentName = &documentName
	}

	// collect original signer documents
	xmlSignerData := b.buildXmlSignerDataList(b.signatures, b.usedTimestamps, b.evidenceRecords)
	diagnosticData.OriginalDocuments = &jaxb.OriginalDocumentsWrapper{Items: xmlSignerData}

	if utils.IsCollectionNotEmpty(b.signatures) {
		xmlSignatures := b.buildXmlSignatures(b.signatures)
		diagnosticData.Signatures = &jaxb.SignaturesWrapper{Items: xmlSignatures}
		b.attachCounterSignatures(b.signatures)
	}

	if utils.IsCollectionNotEmpty(b.usedTimestamps) {
		builtTimestamps := b.buildXmlTimestamps(b.usedTimestamps)
		diagnosticData.UsedTimestamps = &jaxb.UsedTimestampsWrapper{Items: builtTimestamps}
		b.linkSignaturesAndTimestamps(b.signatures)
	}

	if utils.IsCollectionNotEmpty(b.evidenceRecords) {
		builtEvidenceRecords := b.buildXmlEvidenceRecords(b.evidenceRecords)
		diagnosticData.EvidenceRecords = &jaxb.EvidenceRecordsWrapper{Items: builtEvidenceRecords}
		b.linkSignaturesAndEvidenceRecords(b.signatures)
		b.linkTimestampsAndEvidenceRecords(b.usedTimestamps)
		b.linkEvidenceRecordsAndTimestamps(b.evidenceRecords)
	}

	// link the rest certificates
	b.DiagnosticDataBuilder.LinkSigningCertificateAndChains(b.usedCertificates)

	diagnosticData.OrphanTokens = b.BuildXmlOrphanTokens()

	// timestamped objects must be linked after building of orphan tokens
	if utils.IsCollectionNotEmpty(b.usedTimestamps) {
		b.linkTimestampsAndTimestampsObjects(b.usedTimestamps)
	}
	if utils.IsCollectionNotEmpty(b.evidenceRecords) {
		b.linkEvidenceRecordsAndTimestampsObjects(b.evidenceRecords)
	}

	return diagnosticData
}

// AssertConfigurationValid verifies whether the configuration is valid in order to build a
// Diagnostic Data. Port of the protected assertConfigurationValid(): Java's
// Objects.requireNonNull panics with the same message.
func (b *SignedDocumentDiagnosticDataBuilder) AssertConfigurationValid() {
	if b.signedDocument == nil {
		panic("signedDocument shall be provided! Use 'document()' method.")
	}
}

func (b *SignedDocumentDiagnosticDataBuilder) removeSpecialCharsForXml(text string) string {
	if utils.IsStringNotEmpty(text) {
		return replaceAll(text, "&", "")
	}
	return ""
}

func replaceAll(s, old, new string) string {
	result := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		if len(old) > 0 && i+len(old) <= len(s) && s[i:i+len(old)] == old {
			result = append(result, new...)
			i += len(old) - 1
		} else {
			result = append(result, s[i])
		}
	}
	return string(result)
}

func (b *SignedDocumentDiagnosticDataBuilder) buildXmlSignerDataList(signatures []validation.AdvancedSignature,
	timestamps []*validation.TimestampToken, evidenceRecords []validation.EvidenceRecord) []*jaxb.XmlSignerData {
	signerDataList := make([]*jaxb.XmlSignerData, 0)
	if utils.IsCollectionNotEmpty(signatures) {
		for _, sig := range signatures {
			signerDataList = append(signerDataList, b.buildXmlSignerData(sig.SignatureScopes(), nil)...)
		}
	}
	if utils.IsCollectionNotEmpty(timestamps) {
		for _, timestampToken := range timestamps {
			signerDataList = append(signerDataList, b.buildXmlSignerData(timestampToken.TimestampScopes(), nil)...)
		}
	}
	if utils.IsCollectionNotEmpty(evidenceRecords) {
		for _, evidenceRecord := range evidenceRecords {
			signerDataList = append(signerDataList, b.buildXmlSignerData(evidenceRecord.EvidenceRecordScopes(), nil)...)
		}
	}
	return signerDataList
}

func (b *SignedDocumentDiagnosticDataBuilder) buildXmlSignerData(signatureScopes []scope.SignatureScope, parentSignerData *jaxb.XmlSignerData) []*jaxb.XmlSignerData {
	result := make([]*jaxb.XmlSignerData, 0)
	if utils.IsCollectionNotEmpty(signatureScopes) {
		for _, signatureScope := range signatureScopes {
			if b.xmlSignedDataMap[signatureScope.DSSIDAsString()] == nil {
				xmlSignerData := b.getOrBuildXmlSignerData(signatureScope)
				if parentSignerData != nil {
					xmlSignerData.Parent = parentSignerData
				}
				result = append(result, xmlSignerData)
				if utils.IsCollectionNotEmpty(signatureScope.Children()) {
					result = append(result, b.buildXmlSignerData(signatureScope.Children(), xmlSignerData)...)
				}
			}
		}
	}
	return result
}

func (b *SignedDocumentDiagnosticDataBuilder) getOrBuildXmlSignerData(signatureScope scope.SignatureScope) *jaxb.XmlSignerData {
	id := signatureScope.DSSIDAsString()
	xmlSignerData, ok := b.xmlSignedDataMap[id]
	if !ok {
		xmlSignerData = b.getXmlSignerData(signatureScope)
		b.xmlSignedDataMap[id] = xmlSignerData
	}
	return xmlSignerData
}

func (b *SignedDocumentDiagnosticDataBuilder) getXmlSignerData(signatureScope scope.SignatureScope) *jaxb.XmlSignerData {
	xmlSignedData := &jaxb.XmlSignerData{}
	idStr := b.identifierProvider.IDAsString(signatureScope)
	xmlSignedData.Id = jaxb.NewCollapsedString(idStr)
	digest, err := signatureScope.Digest(b.defaultDigestAlgorithm)
	if err != nil {
		panic(err)
	}
	xmlSignedData.DigestAlgoAndValue = b.GetXmlDigestAlgoAndValueForDigestValue(digest)
	referencedName := signatureScope.Name(b.identifierProvider)
	xmlSignedData.ReferencedName = &referencedName
	return xmlSignedData
}

func (b *SignedDocumentDiagnosticDataBuilder) buildXmlSignatures(signatures []validation.AdvancedSignature) []*jaxb.XmlSignature {
	builtSignatures := make([]*jaxb.XmlSignature, 0)
	for _, advancedSignature := range signatures {
		id := advancedSignature.ID()
		xmlSignature, ok := b.xmlSignaturesMap[id]
		if !ok {
			xmlSignature = b.getXmlSignature(advancedSignature)
			builtSignatures = append(builtSignatures, xmlSignature)
		}
	}
	return builtSignatures
}

func (b *SignedDocumentDiagnosticDataBuilder) attachCounterSignatures(signatures []validation.AdvancedSignature) {
	for _, advancedSignature := range signatures {
		if advancedSignature.IsCounterSignature() {
			currentSignature := b.xmlSignaturesMap[advancedSignature.ID()]
			// attach master
			masterSignature := advancedSignature.MasterSignature()
			xmlMasterSignature := b.xmlSignaturesMap[masterSignature.ID()]
			currentSignature.CounterSignature = boolPtr(true)
			currentSignature.Parent = xmlMasterSignature
		}
	}
}

func boolPtr(v bool) *bool { return &v }

func (b *SignedDocumentDiagnosticDataBuilder) getXmlSignature(sig validation.AdvancedSignature) *jaxb.XmlSignature {
	xmlSignature := b.signedDocumentDiagnosticDataBuilderOverrides().BuildDetachedXmlSignature(sig)
	b.checkDuplicatesSignature(xmlSignature, sig)

	b.setXmlSigningCertificate(xmlSignature, sig)
	b.setXmlPolicy(xmlSignature, sig)

	xmlSignature.FoundCertificates = b.GetXmlFoundCertificatesForToken(sig.DSSID(), (foundCertificatesSource)(sig.CertificateSource()))
	xmlSignature.FoundRevocations = b.getXmlFoundRevocations(
		sig.CRLSource().(revocationSourceRefs[revocation.CRL]), sig.OCSPSource().(revocationSourceRefs[revocation.OCSP]))
	xmlSignature.SignatureScopes = &jaxb.SignatureScopesWrapper{Items: b.getXmlSignatureScopes(sig.SignatureScopes())}

	b.xmlSignaturesMap[sig.ID()] = xmlSignature

	return xmlSignature
}

func (b *SignedDocumentDiagnosticDataBuilder) checkDuplicatesSignature(xmlSignature *jaxb.XmlSignature, sig validation.AdvancedSignature) {
	if b.hasDuplicateSignature(sig) {
		duplicated := true
		xmlSignature.Duplicated = &duplicated
	}
}

func (b *SignedDocumentDiagnosticDataBuilder) hasDuplicateSignature(currentSignature validation.AdvancedSignature) bool {
	// NOTE: with DSS-3847 we introduce a stricter verification of the signature duplication,
	// verifying by DSS and DA identifiers across all signature files, as well as a SignatureDigestReference
	for _, sig := range b.signatures {
		if currentSignature != sig {
			sameID := currentSignature.ID() == sig.ID()
			sameDA := currentSignature.DAIdentifier() != "" && currentSignature.DAIdentifier() == sig.DAIdentifier()
			sdr1 := currentSignature.SignatureDigestReference(b.defaultDigestAlgorithm)
			sdr2 := sig.SignatureDigestReference(b.defaultDigestAlgorithm)
			sameDigestRef := sdr1 != nil && sdr2 != nil && sdr1.Equals(sdr2)
			if sameID || sameDA || sameDigestRef {
				return true
			}
		}
	}
	return false
}

func (b *SignedDocumentDiagnosticDataBuilder) setXmlSigningCertificate(xmlSignature *jaxb.XmlSignature, sig validation.AdvancedSignature) {
	candidatesForSigningCertificate := sig.CandidatesForSigningCertificate()
	theCertificateValidity := candidatesForSigningCertificate.TheCertificateValidity()
	var signingCertificatePublicKey *model.PublicKey
	if theCertificateValidity != nil {
		xmlSignature.SigningCertificate = b.GetXmlSigningCertificateForIdentifier(sig.DSSID(), theCertificateValidity)
		xmlSignature.CertificateChain = b.certificateChainWrapper(b.GetXmlForCertificateChainForValidity(theCertificateValidity, sig.CertificateSource()))
		signingCertificatePublicKey = theCertificateValidity.PublicKey()
	}

	xmlSignature.BasicSignature = b.GetXmlBasicSignatureForSignature(sig, signingCertificatePublicKey)
	xmlSignature.DigestMatchers = &jaxb.DigestMatchersWrapper{Items: b.getXmlDigestMatchersForSignature(sig)}
}

func (b *SignedDocumentDiagnosticDataBuilder) setXmlPolicy(xmlSignature *jaxb.XmlSignature, sig validation.AdvancedSignature) {
	if sig.SignaturePolicy() != nil {
		policyBuilder := b.getPolicyBuilder(sig)
		xmlSignature.Policy = policyBuilder.Build()
		xmlSignature.SignaturePolicyStore = policyBuilder.BuildSignaturePolicyStore()
	}
}

// BuildDetachedXmlSignature builds the XmlSignature. Port of the public
// buildDetachedXmlSignature(AdvancedSignature). This is the base (default) body of
// SignedDocumentDiagnosticDataBuilderOverrides.BuildDetachedXmlSignature.
func (b *SignedDocumentDiagnosticDataBuilder) BuildDetachedXmlSignature(sig validation.AdvancedSignature) *jaxb.XmlSignature {
	xmlSignature := &jaxb.XmlSignature{}
	filename := b.removeSpecialCharsForXml(sig.Filename())
	xmlSignature.SignatureFilename = &filename

	id := b.identifierProvider.IDAsString(sig)
	xmlSignature.Id = jaxb.NewCollapsedString(id)
	if daID := sig.DAIdentifier(); daID != "" {
		xmlSignature.DAIdentifier = &daID
	}
	if signingTime := sig.SigningTime(); signingTime != nil {
		xmlSignature.ClaimedSigningTime = jaxb.NewXSDateTime(*signingTime)
	}
	xmlSignature.StructuralValidation = b.GetXmlStructuralValidationForSignature(sig)
	signatureFormat := jaxb.SignatureLevelValue(sig.DataFoundUpToLevel())
	xmlSignature.SignatureFormat = &signatureFormat

	xmlSignature.SignatureProductionPlace = b.getXmlSignatureProductionPlace(sig.SignatureProductionPlace())
	xmlSignature.CommitmentTypeIndications = &jaxb.CommitmentTypeIndicationsWrapper{Items: b.getXmlCommitmentTypeIndications(sig.CommitmentTypeIndications())}
	xmlSignature.SignerRole = b.getXmlSignerRoles(sig.SignerRoles())

	if signatureType := sig.SignatureType(); signatureType != "" {
		xmlSignature.SignatureType = &signatureType
	}
	if contentType := sig.ContentType(); contentType != "" {
		xmlSignature.ContentType = &contentType
	}
	if mimeType := sig.MimeType(); mimeType != "" {
		xmlSignature.MimeType = &mimeType
	}

	xmlSignature.SignatureDigestReference = b.getXmlSignatureDigestReference(sig)

	xmlSignature.DataToBeSignedRepresentation = b.getXmlDataToBeSignedRepresentation(sig)
	xmlSignature.SignerDocumentRepresentations = b.getXmlSignerDocumentRepresentations(sig)

	if sv := sig.SignatureValue(); sv != nil {
		bin := jaxb.Base64Binary(sv)
		xmlSignature.SignatureValue = &bin
	}

	return xmlSignature
}

// GetXmlStructuralValidationForSignature gets the structural validation result of an advanced
// signature. Port of the protected getXmlStructuralValidation(AdvancedSignature).
func (b *SignedDocumentDiagnosticDataBuilder) GetXmlStructuralValidationForSignature(sig validation.AdvancedSignature) *jaxb.XmlStructuralValidation {
	return b.GetXmlStructuralValidation(sig.StructureValidationResult())
}

func (b *SignedDocumentDiagnosticDataBuilder) getXmlSignatureProductionPlace(place *signature.SignatureProductionPlace) *jaxb.XmlSignatureProductionPlace {
	if place != nil {
		xmlSignatureProductionPlace := &jaxb.XmlSignatureProductionPlace{}
		xmlSignatureProductionPlace.CountryName = b.emptyToNilPtr(place.CountryName())
		xmlSignatureProductionPlace.StateOrProvince = b.emptyToNilPtr(place.StateOrProvince())
		xmlSignatureProductionPlace.PostOfficeBoxNumber = b.emptyToNilPtr(place.PostOfficeBoxNumber())
		xmlSignatureProductionPlace.PostalCode = b.emptyToNilPtr(place.PostalCode())
		xmlSignatureProductionPlace.StreetAddress = b.emptyToNilPtr(place.StreetAddress())
		xmlSignatureProductionPlace.City = b.emptyToNilPtr(place.City())
		if utils.IsCollectionNotEmpty(place.PostalAddress()) {
			xmlSignatureProductionPlace.PostalAddress = place.PostalAddress()
		}
		return xmlSignatureProductionPlace
	}
	return nil
}

// EmptyToNil returns nil if text is empty, or the original text otherwise. Port of the
// protected emptyToNull(String).
func (b *SignedDocumentDiagnosticDataBuilder) EmptyToNil(text string) string {
	return text
}

func (b *SignedDocumentDiagnosticDataBuilder) emptyToNilPtr(text string) *string {
	if utils.IsStringEmpty(text) {
		return nil
	}
	return &text
}

func (b *SignedDocumentDiagnosticDataBuilder) getXmlCommitmentTypeIndications(commitmentTypeIndications []*signature.CommitmentTypeIndication) []*jaxb.XmlCommitmentTypeIndication {
	if utils.IsCollectionNotEmpty(commitmentTypeIndications) {
		result := make([]*jaxb.XmlCommitmentTypeIndication, 0, len(commitmentTypeIndications))
		for _, cti := range commitmentTypeIndications {
			result = append(result, b.getXmlCommitmentTypeIndication(cti))
		}
		return result
	}
	return []*jaxb.XmlCommitmentTypeIndication{}
}

func (b *SignedDocumentDiagnosticDataBuilder) getXmlCommitmentTypeIndication(cti *signature.CommitmentTypeIndication) *jaxb.XmlCommitmentTypeIndication {
	xmlCommitmentTypeIndication := &jaxb.XmlCommitmentTypeIndication{}
	identifier := cti.Identifier()
	xmlCommitmentTypeIndication.Identifier = &identifier
	description := cti.Description()
	xmlCommitmentTypeIndication.Description = &description
	xmlCommitmentTypeIndication.DocumentationReferences = &jaxb.DocumentationReferencesWrapper{Items: cti.DocumentReferences()}
	if cti.IsAllDataSignedObjects() {
		v := cti.IsAllDataSignedObjects()
		xmlCommitmentTypeIndication.AllDataSignedObjects = &v
	} else {
		xmlCommitmentTypeIndication.ObjectReferences = &jaxb.ObjectReferencesWrapper{Items: cti.ObjectReferences()}
	}
	return xmlCommitmentTypeIndication
}

func (b *SignedDocumentDiagnosticDataBuilder) getXmlSignerRoles(signerRoles []*signature.SignerRole) []*jaxb.XmlSignerRole {
	xmlSignerRoles := make([]*jaxb.XmlSignerRole, 0)
	if utils.IsCollectionNotEmpty(signerRoles) {
		for _, signerRole := range signerRoles {
			xmlSignerRole := &jaxb.XmlSignerRole{}
			role := signerRole.Role()
			xmlSignerRole.Role = &role
			category := jaxb.EndorsementTypeValue(signerRole.Category())
			xmlSignerRole.Category = &category
			if !signerRole.NotBefore().IsZero() {
				xmlSignerRole.NotBefore = jaxb.NewXSDateTime(signerRole.NotBefore())
			}
			if !signerRole.NotAfter().IsZero() {
				xmlSignerRole.NotAfter = jaxb.NewXSDateTime(signerRole.NotAfter())
			}
			xmlSignerRoles = append(xmlSignerRoles, xmlSignerRole)
		}
	}
	return xmlSignerRoles
}

// GetXmlBasicSignatureForSignature gets an XmlBasicSignature for a signature. Port of the
// protected getXmlBasicSignature(AdvancedSignature, PublicKey).
func (b *SignedDocumentDiagnosticDataBuilder) GetXmlBasicSignatureForSignature(sig validation.AdvancedSignature, signingCertificatePublicKey *model.PublicKey) *jaxb.XmlBasicSignature {
	xmlBasicSignature := &jaxb.XmlBasicSignature{}
	if encAlg := sig.EncryptionAlgorithm(); encAlg != "" {
		v := jaxb.EncryptionAlgorithmValue(encAlg)
		xmlBasicSignature.EncryptionAlgoUsedToSignThisToken = &v
	}
	keyLength := spi.DSSPKUtilsStringPublicKeySize(signingCertificatePublicKey)
	xmlBasicSignature.KeyLengthUsedToSignThisToken = &keyLength
	if digestAlg := sig.DigestAlgorithm(); digestAlg != "" {
		v := jaxb.DigestAlgorithmValue(digestAlg)
		xmlBasicSignature.DigestAlgoUsedToSignThisToken = &v
	}

	scv := sig.SignatureCryptographicVerification()
	signatureIntact := scv.IsSignatureIntact()
	xmlBasicSignature.SignatureIntact = &signatureIntact
	signatureValid := scv.IsSignatureValid()
	xmlBasicSignature.SignatureValid = &signatureValid
	return xmlBasicSignature
}

func (b *SignedDocumentDiagnosticDataBuilder) getXmlDigestMatchersForSignature(sig validation.AdvancedSignature) []*jaxb.XmlDigestMatcher {
	return b.getXmlDigestMatchers(sig.ReferenceValidations(), sig.DetachedContents())
}

func (b *SignedDocumentDiagnosticDataBuilder) getXmlDigestMatchers(referenceValidations []*model.ReferenceValidation, detachedContents []model.DSSDocument) []*jaxb.XmlDigestMatcher {
	refs := make([]*jaxb.XmlDigestMatcher, 0)
	if utils.IsCollectionNotEmpty(referenceValidations) {
		for _, referenceValidation := range referenceValidations {
			refs = append(refs, b.getXmlDigestMatcher(referenceValidation))
			dependentValidations := referenceValidation.DependentValidations()
			if utils.IsCollectionNotEmpty(dependentValidations) &&
				(utils.IsCollectionNotEmpty(detachedContents) || b.isAtLeastOneFound(dependentValidations)) {
				for _, dependentValidation := range referenceValidation.DependentValidations() {
					refs = append(refs, b.getXmlDigestMatcher(dependentValidation))
				}
			}
		}
	}
	return refs
}

func (b *SignedDocumentDiagnosticDataBuilder) getTimestampReferenceValidationDigestMatchers(referenceValidations []*model.ReferenceValidation) []*jaxb.XmlDigestMatcher {
	refs := make([]*jaxb.XmlDigestMatcher, 0)
	if utils.IsCollectionNotEmpty(referenceValidations) {
		for _, referenceValidation := range referenceValidations {
			refs = append(refs, b.getXmlDigestMatcher(referenceValidation))
		}
	}
	return refs
}

func (b *SignedDocumentDiagnosticDataBuilder) getXmlDigestMatcher(referenceValidation *model.ReferenceValidation) *jaxb.XmlDigestMatcher {
	ref := &jaxb.XmlDigestMatcher{}
	t := jaxb.DigestMatcherTypeValue(referenceValidation.Type())
	ref.Type = &t
	if id := referenceValidation.Id(); id != "" {
		ref.Id = &id
	}
	if uri := referenceValidation.Uri(); uri != "" {
		ref.Uri = &uri
	}
	if refs := referenceValidation.DataObjectReferences(); refs != nil {
		ref.DataObjectReferences = &jaxb.DataObjectReferencesWrapper{Items: refs}
	}
	if referenceValidation.Document() != nil {
		docName := referenceValidation.Document().Name()
		ref.DocumentName = &docName
	}
	digest := referenceValidation.Digest()
	if !digest.IsEmpty() {
		v := jaxb.Base64Binary(digest.Value())
		ref.DigestValue = &v
		dm := jaxb.DigestAlgorithmValue(digest.Algorithm())
		ref.DigestMethod = &dm
	}
	dataFound := referenceValidation.IsFound()
	ref.DataFound = dataFound
	dataIntact := referenceValidation.IsIntact()
	ref.DataIntact = dataIntact
	if referenceValidation.IsDuplicated() {
		v := referenceValidation.IsDuplicated()
		ref.Duplicated = &v
	}
	return ref
}

func (b *SignedDocumentDiagnosticDataBuilder) isAtLeastOneFound(referenceValidations []*model.ReferenceValidation) bool {
	for _, referenceValidation := range referenceValidations {
		if referenceValidation.IsFound() {
			return true
		}
	}
	return false
}

func (b *SignedDocumentDiagnosticDataBuilder) getPolicyBuilder(sig validation.AdvancedSignature) *XmlPolicyBuilder {
	signaturePolicy := sig.SignaturePolicy()
	signaturePolicyStore := sig.SignaturePolicyStore()

	xmlPolicyBuilder := NewXmlPolicyBuilder(signaturePolicy)
	xmlPolicyBuilder.SetSignaturePolicyStore(signaturePolicyStore)
	return xmlPolicyBuilder
}

func (b *SignedDocumentDiagnosticDataBuilder) getXmlSignatureDigestReference(sig validation.AdvancedSignature) *jaxb.XmlSignatureDigestReference {
	signatureDigestReference := sig.SignatureDigestReference(b.defaultDigestAlgorithm)
	if signatureDigestReference != nil {
		xmlDigestReference := &jaxb.XmlSignatureDigestReference{}
		canonicalizationMethod := signatureDigestReference.CanonicalizationMethod()
		xmlDigestReference.CanonicalizationMethod = &canonicalizationMethod
		dm := jaxb.DigestAlgorithmValue(signatureDigestReference.DigestAlgorithm())
		xmlDigestReference.DigestMethod = &dm
		v := jaxb.Base64Binary(signatureDigestReference.DigestValue())
		xmlDigestReference.DigestValue = &v
		return xmlDigestReference
	}
	return nil
}

func (b *SignedDocumentDiagnosticDataBuilder) getXmlDataToBeSignedRepresentation(sig validation.AdvancedSignature) *jaxb.XmlDigestAlgoAndValue {
	dtbsr := sig.DataToBeSignedRepresentation()
	if !dtbsr.IsEmpty() {
		return b.GetXmlDigestAlgoAndValueForDigestValue(dtbsr)
	}
	return nil
}

func (b *SignedDocumentDiagnosticDataBuilder) getXmlSignerDocumentRepresentations(sig validation.AdvancedSignature) *jaxb.XmlSignerDocumentRepresentations {
	if utils.IsCollectionEmpty(sig.DetachedContents()) {
		return nil
	}
	signerDocumentRepresentation := &jaxb.XmlSignerDocumentRepresentations{}
	signerDocumentRepresentation.DocHashOnly = sig.IsDocHashOnlyValidation()
	signerDocumentRepresentation.HashOnly = sig.IsHashOnlyValidation()
	return signerDocumentRepresentation
}

func (b *SignedDocumentDiagnosticDataBuilder) getXmlFoundRevocations(crlSource revocationSourceRefs[revocation.CRL], ocspSource revocationSourceRefs[revocation.OCSP]) *jaxb.XmlFoundRevocations {
	foundRevocations := &jaxb.XmlFoundRevocations{}
	foundRevocations.RelatedRevocation = b.getXmlRelatedRevocations(crlSource, ocspSource)
	foundRevocations.OrphanRevocation = b.getXmlOrphanRevocations(crlSource, ocspSource)
	foundRevocations.OrphanRevocation = append(foundRevocations.OrphanRevocation, b.getXmlOrphanRevocationRefs(crlSource, ocspSource)...)
	return foundRevocations
}

func (b *SignedDocumentDiagnosticDataBuilder) getXmlRelatedRevocations(crlSource revocationSourceRefs[revocation.CRL], ocspSource revocationSourceRefs[revocation.OCSP]) []*jaxb.XmlRelatedRevocation {
	xmlRelatedRevocations := make([]*jaxb.XmlRelatedRevocation, 0)
	xmlRelatedRevocations = append(xmlRelatedRevocations, addRelatedRevocations(b, crlSource)...)
	xmlRelatedRevocations = append(xmlRelatedRevocations, addRelatedRevocations(b, ocspSource)...)
	return xmlRelatedRevocations
}

func addRelatedRevocations[R revocation.Revocation](b *SignedDocumentDiagnosticDataBuilder, source revocationSourceRefs[R]) []*jaxb.XmlRelatedRevocation {
	result := make([]*jaxb.XmlRelatedRevocation, 0)
	for _, entry := range source.UniqueRevocationTokensWithOrigins() {
		token := entry.Token
		id := token.DSSIDAsString()
		xmlRevocation, ok := b.xmlRevocationsMap[id]
		if !ok {
			continue
		}
		xmlRelatedRevocation := &jaxb.XmlRelatedRevocation{}
		xmlRelatedRevocation.Revocation = xmlRevocation
		t := jaxb.RevocationTypeValue(token.RevocationType())
		xmlRelatedRevocation.Type = &t
		xmlRelatedRevocation.Origin = revocationOriginValues(entry.Origins)
		xmlRelatedRevocation.RevocationRef = GetXmlRevocationRefs(&b.DiagnosticDataBuilder, xmlRevocation.Id.String(), source.FindRefsAndOriginsForRevocationToken(token))
		result = append(result, xmlRelatedRevocation)
	}
	return result
}

func revocationOriginValues(origins []enumerations.RevocationOrigin) []jaxb.RevocationOriginValue {
	result := make([]jaxb.RevocationOriginValue, 0, len(origins))
	for _, o := range origins {
		result = append(result, jaxb.RevocationOriginValue(o))
	}
	return result
}

func (b *SignedDocumentDiagnosticDataBuilder) getXmlOrphanRevocations(crlSource revocationSourceRefs[revocation.CRL], ocspSource revocationSourceRefs[revocation.OCSP]) []*jaxb.XmlOrphanRevocation {
	xmlOrphanRevocations := make([]*jaxb.XmlOrphanRevocation, 0)
	xmlOrphanRevocations = append(xmlOrphanRevocations, addOrphanRevocations(b, crlSource)...)
	xmlOrphanRevocations = append(xmlOrphanRevocations, addOrphanRevocations(b, ocspSource)...)
	return xmlOrphanRevocations
}

func addOrphanRevocations[R revocation.Revocation](b *SignedDocumentDiagnosticDataBuilder, source revocationSourceRefs[R]) []*jaxb.XmlOrphanRevocation {
	result := make([]*jaxb.XmlOrphanRevocation, 0)
	allBinariesWithOrigins := allRevocationBinariesWithOrigins(source)
	for _, entry := range allBinariesWithOrigins {
		token := entry.binary
		tokenID := token.AsXmlID()
		if _, exists := b.xmlRevocationsMap[tokenID]; exists {
			continue
		}
		xmlOrphanRevocation := getXmlOrphanRevocation[R](b, token, entry.origins)
		xmlOrphanRevocation.RevocationRef = append(xmlOrphanRevocation.RevocationRef, GetXmlRevocationRefs(&b.DiagnosticDataBuilder, tokenID, source.FindRefsAndOriginsForBinary(token))...)
		result = append(result, xmlOrphanRevocation)
	}
	return result
}

// revocationSourceRefs is the local widening interface every OfflineRevocationSourceBase[R]-
// derived value (whatever its narrower static type at the call site: the
// spi.OfflineRevocationSource[R] interface AdvancedSignature.CRLSource()/OCSPSource() return,
// the concrete *spi.OfflineCRLSourceBase/*spi.OfflineOCSPSourceBase EvidenceRecord returns, or
// TimestampToken's TimestampCRLSource/TimestampOCSPSource) is type-asserted down to at every
// call site in this file; see the file header's revocation-source note.
type revocationSourceRefs[R revocation.Revocation] interface {
	UniqueRevocationTokensWithOrigins() []spi.RevocationTokenOriginsEntry[R]
	FindRefsAndOriginsForRevocationToken(token spi.RevocationToken[R]) []spi.RevocationRefOriginsEntry[R]
	OrphanRevocationReferencesWithOrigins() []spi.RevocationRefOriginsEntry[R]
	FindRefsAndOriginsForBinary(identifier spi.EncapsulatedRevocationTokenIdentifier[R]) []spi.RevocationRefOriginsEntry[R]
	AllRevocationBinaries() []spi.EncapsulatedRevocationTokenIdentifier[R]
	CMSSignedDataRevocationBinaries() []spi.EncapsulatedRevocationTokenIdentifier[R]
	RevocationValuesBinaries() []spi.EncapsulatedRevocationTokenIdentifier[R]
	AttributeRevocationValuesBinaries() []spi.EncapsulatedRevocationTokenIdentifier[R]
	TimestampValidationDataBinaries() []spi.EncapsulatedRevocationTokenIdentifier[R]
	AnyValidationDataBinaries() []spi.EncapsulatedRevocationTokenIdentifier[R]
	DSSDictionaryBinaries() []spi.EncapsulatedRevocationTokenIdentifier[R]
	VRIDictionaryBinaries() []spi.EncapsulatedRevocationTokenIdentifier[R]
	ADBERevocationValuesBinaries() []spi.EncapsulatedRevocationTokenIdentifier[R]
}

// revocationBinaryWithOrigins pairs a revocation binary with the origins it has been found
// with, standing in for a Java Map<EncapsulatedRevocationTokenIdentifier<R>,
// Set<RevocationOrigin>> entry; see the file header's map-reconstruction note.
type revocationBinaryWithOrigins[R revocation.Revocation] struct {
	binary  spi.EncapsulatedRevocationTokenIdentifier[R]
	origins []enumerations.RevocationOrigin
}

// allRevocationBinariesWithOrigins reconstructs Java's getAllRevocationBinariesWithOrigins()
// from the 8 per-origin binary getters spi.OfflineRevocationSourceBase[R] exposes; see the file
// header's fidelity note.
func allRevocationBinariesWithOrigins[R revocation.Revocation](source revocationSourceRefs[R]) []revocationBinaryWithOrigins[R] {
	origins := []struct {
		origin  enumerations.RevocationOrigin
		binarys []spi.EncapsulatedRevocationTokenIdentifier[R]
	}{
		{enumerations.RevocationOrigin_CMS_SIGNED_DATA, source.CMSSignedDataRevocationBinaries()},
		{enumerations.RevocationOrigin_REVOCATION_VALUES, source.RevocationValuesBinaries()},
		{enumerations.RevocationOrigin_ATTRIBUTE_REVOCATION_VALUES, source.AttributeRevocationValuesBinaries()},
		{enumerations.RevocationOrigin_TIMESTAMP_VALIDATION_DATA, source.TimestampValidationDataBinaries()},
		{enumerations.RevocationOrigin_ANY_VALIDATION_DATA, source.AnyValidationDataBinaries()},
		{enumerations.RevocationOrigin_DSS_DICTIONARY, source.DSSDictionaryBinaries()},
		{enumerations.RevocationOrigin_VRI_DICTIONARY, source.VRIDictionaryBinaries()},
		{enumerations.RevocationOrigin_ADBE_REVOCATION_INFO_ARCHIVAL, source.ADBERevocationValuesBinaries()},
	}

	byID := map[string]*revocationBinaryWithOrigins[R]{}
	var order []string
	for _, group := range origins {
		for _, binary := range group.binarys {
			id := binary.AsXmlID()
			entry, ok := byID[id]
			if !ok {
				entry = &revocationBinaryWithOrigins[R]{binary: binary}
				byID[id] = entry
				order = append(order, id)
			}
			entry.origins = append(entry.origins, group.origin)
		}
	}

	result := make([]revocationBinaryWithOrigins[R], 0, len(order))
	for _, id := range order {
		result = append(result, *byID[id])
	}
	return result
}

func getXmlOrphanRevocation[R revocation.Revocation](b *SignedDocumentDiagnosticDataBuilder, token spi.EncapsulatedRevocationTokenIdentifier[R], origins []enumerations.RevocationOrigin) *jaxb.XmlOrphanRevocation {
	xmlOrphanRevocation := &jaxb.XmlOrphanRevocation{}
	var revType enumerations.RevocationType
	if _, ok := any(token).(*crlparser.CRLBinary); ok {
		revType = enumerations.RevocationType_CRL
	} else {
		revType = enumerations.RevocationType_OCSP
	}
	t := jaxb.RevocationTypeValue(revType)
	xmlOrphanRevocation.Type = &t
	xmlOrphanRevocation.Origin = revocationOriginValues(origins)
	xmlOrphanRevocation.Token = createOrphanTokenFromRevocationIdentifier[R](b, token)
	return xmlOrphanRevocation
}

// CreateOrphanTokenFromRevocationIdentifier creates an orphan revocation token from an
// EncapsulatedRevocationTokenIdentifier. Port of the protected
// createOrphanTokenFromRevocationIdentifier(EncapsulatedRevocationTokenIdentifier).
//
// Go has no generic methods, only generic free functions (see PORTING.md), so this is one; R is
// inferred from the identifier at every call site.
func CreateOrphanTokenFromRevocationIdentifier[R revocation.Revocation](b *SignedDocumentDiagnosticDataBuilder, revocationIdentifier spi.EncapsulatedRevocationTokenIdentifier[R]) *jaxb.XmlOrphanRevocationToken {
	return createOrphanTokenFromRevocationIdentifier[R](b, revocationIdentifier)
}

func createOrphanTokenFromRevocationIdentifier[R revocation.Revocation](b *SignedDocumentDiagnosticDataBuilder, revocationIdentifier spi.EncapsulatedRevocationTokenIdentifier[R]) *jaxb.XmlOrphanRevocationToken {
	id := revocationIdentifier.AsXmlID()
	orphanToken, ok := b.xmlOrphanRevocationTokensMap[id]
	if ok {
		return orphanToken
	}
	orphanToken = &jaxb.XmlOrphanRevocationToken{}
	encType := jaxb.XmlEncapsulationType_BINARIES
	orphanToken.EncapsulationType = &encType
	idStr := b.identifierProvider.IDAsString(revocationIdentifier)
	orphanToken.Id = jaxb.NewCollapsedString(idStr)
	if binRef, ok := any(revocationIdentifier).(revocationBinaryRef); ok && b.tokenExtractionStrategy.IsRevocationData() {
		bin := jaxb.Base64Binary(binRef.Binaries())
		orphanToken.Base64Encoded = &bin
	} else {
		digestValue, err := revocationIdentifier.DigestValue(b.defaultDigestAlgorithm)
		if err != nil {
			panic(err)
		}
		orphanToken.DigestAlgoAndValue = b.GetXmlDigestAlgoAndValueFor(b.defaultDigestAlgorithm, digestValue)
	}

	if crlBin, ok := any(revocationIdentifier).(*crlparser.CRLBinary); ok {
		_ = crlBin
		t := jaxb.RevocationTypeValue(enumerations.RevocationType_CRL)
		orphanToken.RevocationType = &t
	} else if ocspBin, ok := any(revocationIdentifier).(*spi.OCSPResponseBinary); ok {
		t := jaxb.RevocationTypeValue(enumerations.RevocationType_OCSP)
		orphanToken.RevocationType = &t
		ocspCertificateSource, err := spi.NewOCSPCertificateSource(ocspBin.BasicOCSPResp())
		if err != nil {
			panic(err)
		}
		b.GetXmlFoundCertificatesForSource(&ocspCertificateSource.TokenCertificateSource) // create from OCSP Certificate Source
	}
	b.xmlOrphanRevocationTokensMap[id] = orphanToken
	b.xmlOrphanRevocationTokensOrder = append(b.xmlOrphanRevocationTokensOrder, id)
	return orphanToken
}

// revocationBinaryRef is the local widening interface for the AsXmlID/DigestValue/Binaries
// surface every concrete EncapsulatedRevocationTokenIdentifier[R] (model.
// EncapsulatedRevocationTokenIdentifier[R]-embedding: crlparser.CRLBinary, spi.OCSPResponseBinary)
// exposes via promotion, wider than the spi.EncapsulatedRevocationTokenIdentifier[R] interface.
type revocationBinaryRef interface {
	Binaries() []byte
}

func (b *SignedDocumentDiagnosticDataBuilder) getXmlOrphanRevocationRefs(crlSource revocationSourceRefs[revocation.CRL], ocspSource revocationSourceRefs[revocation.OCSP]) []*jaxb.XmlOrphanRevocation {
	xmlOrphanRevocationRefs := make([]*jaxb.XmlOrphanRevocation, 0)
	xmlOrphanRevocationRefs = append(xmlOrphanRevocationRefs, addOrphanRevocationRefs(b, crlSource, b.documentCRLSource)...)
	xmlOrphanRevocationRefs = append(xmlOrphanRevocationRefs, addOrphanRevocationRefs(b, ocspSource, b.documentOCSPSource)...)
	return xmlOrphanRevocationRefs
}

func addOrphanRevocationRefs[R revocation.Revocation](b *SignedDocumentDiagnosticDataBuilder, source revocationSourceRefs[R], allSources *spi.ListRevocationSource[R]) []*jaxb.XmlOrphanRevocation {
	result := make([]*jaxb.XmlOrphanRevocation, 0)
	for _, entry := range source.OrphanRevocationReferencesWithOrigins() {
		ref := entry.Reference
		if allSources.IsOrphan(ref) && sourceDoesNotContainOrphanBinaries(b, source, ref) {
			result = append(result, createOrphanRevocationFromRef[R](b, ref, entry.Origins))
		}
	}
	return result
}

func sourceDoesNotContainOrphanBinaries[R revocation.Revocation](b *SignedDocumentDiagnosticDataBuilder, source revocationSourceRefs[R], ref spi.RevocationRef[R]) bool {
	tokenID, ok := b.referenceMap[ref.DSSIDAsString()]
	if !ok {
		return true
	}
	for _, revocationIdentifier := range source.AllRevocationBinaries() {
		if tokenID == revocationIdentifier.AsXmlID() {
			return false
		}
	}
	return true
}

func createOrphanRevocationFromRef[R revocation.Revocation](b *SignedDocumentDiagnosticDataBuilder, ref spi.RevocationRef[R], origins []enumerations.RevocationRefOrigin) *jaxb.XmlOrphanRevocation {
	xmlOrphanRevocation := &jaxb.XmlOrphanRevocation{}

	orphanToken := &jaxb.XmlOrphanRevocationToken{}
	encType := jaxb.XmlEncapsulationType_REFERENCE
	orphanToken.EncapsulationType = &encType
	idStr := b.identifierProvider.IDAsString(ref)
	orphanToken.Id = jaxb.NewCollapsedString(idStr)
	if !ref.Digest().IsEmpty() {
		orphanToken.DigestAlgoAndValue = b.GetXmlDigestAlgoAndValueForDigestValue(ref.Digest())
	}
	b.xmlOrphanRevocationTokensMap[ref.DSSIDAsString()] = orphanToken
	b.xmlOrphanRevocationTokensOrder = append(b.xmlOrphanRevocationTokensOrder, ref.DSSIDAsString())

	xmlOrphanRevocation.Token = orphanToken
	if crlRef, ok := any(ref).(*spi.CRLRef); ok {
		t := jaxb.RevocationTypeValue(enumerations.RevocationType_CRL)
		orphanToken.RevocationType = &t
		xt := jaxb.RevocationTypeValue(enumerations.RevocationType_CRL)
		xmlOrphanRevocation.Type = &xt
		xmlOrphanRevocation.RevocationRef = append(xmlOrphanRevocation.RevocationRef, b.GetXmlCRLRevocationRef(crlRef, origins))
	} else {
		t := jaxb.RevocationTypeValue(enumerations.RevocationType_OCSP)
		orphanToken.RevocationType = &t
		xt := jaxb.RevocationTypeValue(enumerations.RevocationType_OCSP)
		xmlOrphanRevocation.Type = &xt
		xmlOrphanRevocation.RevocationRef = append(xmlOrphanRevocation.RevocationRef, b.GetXmlOCSPRevocationRef(any(ref).(*spi.OCSPRef), origins))
	}
	return xmlOrphanRevocation
}

func (b *SignedDocumentDiagnosticDataBuilder) getXmlSignatureScopes(scopes []scope.SignatureScope) []*jaxb.XmlSignatureScope {
	xmlScopes := make([]*jaxb.XmlSignatureScope, 0)
	if utils.IsCollectionNotEmpty(scopes) {
		for _, signatureScope := range scopes {
			xmlScopes = append(xmlScopes, b.getXmlSignatureScope(signatureScope))
			if utils.IsCollectionNotEmpty(signatureScope.Children()) {
				xmlScopes = append(xmlScopes, b.getXmlSignatureScopes(signatureScope.Children())...)
			}
		}
	}
	return xmlScopes
}

func (b *SignedDocumentDiagnosticDataBuilder) getXmlSignatureScope(sc scope.SignatureScope) *jaxb.XmlSignatureScope {
	xmlSignatureScope := &jaxb.XmlSignatureScope{}
	name := sc.Name(b.identifierProvider)
	xmlSignatureScope.Name = &name
	t := jaxb.SignatureScopeTypeValue(sc.Type())
	xmlSignatureScope.Scope = &t
	description := sc.Description(b.identifierProvider)
	xmlSignatureScope.Description = &description
	if transformations := sc.Transformations(); transformations != nil {
		xmlSignatureScope.Transformations = &jaxb.TransformationsWrapper{Items: transformations}
	}
	xmlSignatureScope.SignerData = b.xmlSignedDataMap[sc.DSSIDAsString()]
	return xmlSignatureScope
}

func (b *SignedDocumentDiagnosticDataBuilder) buildXmlEvidenceRecords(evidenceRecords []validation.EvidenceRecord) []*jaxb.XmlEvidenceRecord {
	xmlEvidenceRecords := make([]*jaxb.XmlEvidenceRecord, 0)
	if utils.IsCollectionNotEmpty(evidenceRecords) {
		for _, evidenceRecord := range evidenceRecords {
			id := evidenceRecord.Id()
			xmlEvidenceRecord, ok := b.xmlEvidenceRecordMap[id]
			if !ok {
				xmlEvidenceRecord = b.buildXmlEvidenceRecord(evidenceRecord)
				xmlEvidenceRecords = append(xmlEvidenceRecords, xmlEvidenceRecord)
			}
		}
	}
	return xmlEvidenceRecords
}

func (b *SignedDocumentDiagnosticDataBuilder) buildXmlEvidenceRecord(evidenceRecord validation.EvidenceRecord) *jaxb.XmlEvidenceRecord {
	xmlEvidenceRecord := &jaxb.XmlEvidenceRecord{}
	b.checkDuplicatesEvidenceRecord(xmlEvidenceRecord, evidenceRecord)

	idStr := b.identifierProvider.IDAsString(evidenceRecord)
	xmlEvidenceRecord.Id = jaxb.NewCollapsedString(idStr)
	if filename := evidenceRecord.Filename(); filename != "" {
		xmlEvidenceRecord.DocumentName = &filename
	}
	t := jaxb.EvidenceRecordTypeEnumValue(evidenceRecord.EvidenceRecordType())
	xmlEvidenceRecord.Type = &t
	origin := jaxb.EvidenceRecordOriginValue(evidenceRecord.Origin())
	xmlEvidenceRecord.Origin = &origin
	incorporationType := jaxb.EvidenceRecordIncorporationTypeValue(evidenceRecord.IncorporationType())
	xmlEvidenceRecord.IncorporationType = &incorporationType
	xmlEvidenceRecord.StructuralValidation = b.getXmlStructuralValidationForEvidenceRecord(evidenceRecord)
	xmlEvidenceRecord.DigestMatchers = &jaxb.DigestMatchersWrapper{Items: b.getXmlDigestMatchersForEvidenceRecord(evidenceRecord)}
	xmlEvidenceRecord.EvidenceRecordScopes = &jaxb.EvidenceRecordScopesWrapper{Items: b.getXmlSignatureScopes(evidenceRecord.EvidenceRecordScopes())}
	xmlEvidenceRecord.FoundCertificates = b.GetXmlFoundCertificatesForSource(evidenceRecord.CertificateSource())
	xmlEvidenceRecord.FoundRevocations = b.getXmlFoundRevocations(evidenceRecord.CRLSource(), evidenceRecord.OCSPSource())

	encoded := evidenceRecord.Encoded()
	if b.tokenExtractionStrategy.IsEvidenceRecord() {
		bin := jaxb.Base64Binary(encoded)
		xmlEvidenceRecord.Base64Encoded = &bin
	} else {
		digest, err := spi.DSSUtilsDigest(b.defaultDigestAlgorithm, encoded)
		if err != nil {
			panic(err)
		}
		xmlEvidenceRecord.DigestAlgoAndValue = b.GetXmlDigestAlgoAndValueFor(b.defaultDigestAlgorithm, digest)
	}

	b.xmlEvidenceRecordMap[evidenceRecord.Id()] = xmlEvidenceRecord

	return xmlEvidenceRecord
}

func (b *SignedDocumentDiagnosticDataBuilder) checkDuplicatesEvidenceRecord(xmlEvidenceRecord *jaxb.XmlEvidenceRecord, evidenceRecord validation.EvidenceRecord) {
	if b.hasDuplicateEvidenceRecord(evidenceRecord) {
		duplicated := true
		xmlEvidenceRecord.Duplicated = &duplicated
	}
}

func (b *SignedDocumentDiagnosticDataBuilder) hasDuplicateEvidenceRecord(currentEvidenceRecord validation.EvidenceRecord) bool {
	for _, evidenceRecord := range b.evidenceRecords {
		if currentEvidenceRecord != evidenceRecord && currentEvidenceRecord.Id() == evidenceRecord.Id() {
			return true
		}
	}
	return false
}

func (b *SignedDocumentDiagnosticDataBuilder) getXmlDigestMatchersForEvidenceRecord(evidenceRecord validation.EvidenceRecord) []*jaxb.XmlDigestMatcher {
	return b.getXmlDigestMatchers(evidenceRecord.ReferenceValidation(), nil)
}

func (b *SignedDocumentDiagnosticDataBuilder) linkSignaturesAndEvidenceRecords(signatures []validation.AdvancedSignature) {
	for _, sig := range signatures {
		xmlSignature := b.xmlSignaturesMap[sig.ID()]
		xmlSignature.FoundEvidenceRecords = &jaxb.FoundEvidenceRecordsWrapper{Items: b.getXmlSignatureEvidenceRecords(sig)}
	}
	for _, evidenceRecord := range b.evidenceRecords {
		if evidenceRecord.IsEmbedded() {
			xmlEvidenceRecord := b.xmlEvidenceRecordMap[evidenceRecord.Id()]
			embedded := evidenceRecord.IsEmbedded()
			xmlEvidenceRecord.Embedded = &embedded

			xmlSignature := b.xmlSignaturesMap[evidenceRecord.MasterSignature().ID()]
			xmlEvidenceRecord.Parent = xmlSignature
		}
	}
}

func (b *SignedDocumentDiagnosticDataBuilder) getXmlSignatureEvidenceRecords(sig validation.AdvancedSignature) []*jaxb.XmlFoundEvidenceRecord {
	foundEvidenceRecords := make([]*jaxb.XmlFoundEvidenceRecord, 0)
	for _, evidenceRecord := range sig.AllEvidenceRecords() {
		foundEvidenceRecord := &jaxb.XmlFoundEvidenceRecord{}
		foundEvidenceRecord.EvidenceRecord = b.xmlEvidenceRecordMap[evidenceRecord.Id()]
		foundEvidenceRecords = append(foundEvidenceRecords, foundEvidenceRecord)
	}
	return foundEvidenceRecords
}

func (b *SignedDocumentDiagnosticDataBuilder) linkTimestampsAndEvidenceRecords(timestampTokens []*validation.TimestampToken) {
	for _, timestampToken := range timestampTokens {
		xmlTimestamp := b.xmlTimestampsMap[timestampToken.DSSIDAsString()]
		xmlTimestamp.FoundEvidenceRecords = &jaxb.FoundEvidenceRecordsWrapper{Items: b.getXmlTimestampEvidenceRecords(timestampToken)}
	}
}

func (b *SignedDocumentDiagnosticDataBuilder) getXmlTimestampEvidenceRecords(timestampToken *validation.TimestampToken) []*jaxb.XmlFoundEvidenceRecord {
	foundEvidenceRecords := make([]*jaxb.XmlFoundEvidenceRecord, 0)
	for _, evidenceRecord := range timestampToken.DetachedEvidenceRecords() {
		foundEvidenceRecord := &jaxb.XmlFoundEvidenceRecord{}
		foundEvidenceRecord.EvidenceRecord = b.xmlEvidenceRecordMap[evidenceRecord.Id()]
		foundEvidenceRecords = append(foundEvidenceRecords, foundEvidenceRecord)
	}
	return foundEvidenceRecords
}

func (b *SignedDocumentDiagnosticDataBuilder) linkEvidenceRecordsAndTimestamps(evidenceRecords []validation.EvidenceRecord) {
	for _, evidenceRecord := range evidenceRecords {
		currentEvidenceRecord := b.xmlEvidenceRecordMap[evidenceRecord.Id()]
		// attach timestamps
		currentEvidenceRecord.EvidenceRecordTimestamps = &jaxb.EvidenceRecordTimestampsWrapper{Items: b.getXmlEvidenceRecordTimestamps(evidenceRecord)}
	}
}

func (b *SignedDocumentDiagnosticDataBuilder) getXmlEvidenceRecordTimestamps(evidenceRecord validation.EvidenceRecord) []*jaxb.XmlFoundTimestamp {
	foundTimestamps := make([]*jaxb.XmlFoundTimestamp, 0)
	for _, timestampToken := range evidenceRecord.Timestamps() {
		foundTimestamp := &jaxb.XmlFoundTimestamp{}
		foundTimestamp.Timestamp = b.xmlTimestampsMap[timestampToken.DSSIDAsString()]
		foundTimestamps = append(foundTimestamps, foundTimestamp)
	}
	return foundTimestamps
}

func (b *SignedDocumentDiagnosticDataBuilder) getXmlStructuralValidationForEvidenceRecord(evidenceRecord validation.EvidenceRecord) *jaxb.XmlStructuralValidation {
	return b.GetXmlStructuralValidation(evidenceRecord.StructureValidationResult())
}

func (b *SignedDocumentDiagnosticDataBuilder) buildXmlTimestamps(timestamps []*validation.TimestampToken) []*jaxb.XmlTimestamp {
	xmlTimestampsList := make([]*jaxb.XmlTimestamp, 0)
	if utils.IsCollectionNotEmpty(timestamps) {
		tokens := append([]*validation.TimestampToken{}, timestamps...)
		comparator := validation.NewTimestampTokenComparator()
		sortTimestamps(tokens, comparator)
		for _, timestampToken := range tokens {
			id := timestampToken.DSSIDAsString()
			xmlTimestamp, ok := b.xmlTimestampsMap[id]
			if !ok {
				xmlTimestamp = b.signedDocumentDiagnosticDataBuilderOverrides().BuildDetachedXmlTimestamp(timestampToken)
				xmlTimestampsList = append(xmlTimestampsList, xmlTimestamp)
			}
		}
	}
	return xmlTimestampsList
}

func sortTimestamps(tokens []*validation.TimestampToken, comparator validation.TimestampTokenComparator) {
	for i := 1; i < len(tokens); i++ {
		for j := i; j > 0 && comparator.Compare(tokens[j-1], tokens[j]) > 0; j-- {
			tokens[j-1], tokens[j] = tokens[j], tokens[j-1]
		}
	}
}

// BuildDetachedXmlTimestamp builds an XmlTimestamp from a TimestampToken. Port of the protected
// buildDetachedXmlTimestamp(TimestampToken). This is the base (default) body of
// SignedDocumentDiagnosticDataBuilderOverrides.BuildDetachedXmlTimestamp.
func (b *SignedDocumentDiagnosticDataBuilder) BuildDetachedXmlTimestamp(timestampToken *validation.TimestampToken) *jaxb.XmlTimestamp {
	xmlTimestampToken := &jaxb.XmlTimestamp{}
	b.checkDuplicatesTimestamp(xmlTimestampToken, timestampToken)

	idStr := b.identifierProvider.IDAsString(timestampToken)
	xmlTimestampToken.Id = jaxb.NewCollapsedString(idStr)
	tsType := jaxb.TimestampTypeValue(timestampToken.TimeStampType())
	xmlTimestampToken.Type = &tsType
	// property is defined only for archival timestamps
	if archiveTimestampType := timestampToken.ArchiveTimestampType(); archiveTimestampType != "" {
		v := jaxb.ArchiveTimestampTypeValue(archiveTimestampType)
		xmlTimestampToken.ArchiveTimestampType = &v
	}
	if ertType := timestampToken.EvidenceRecordTimestampType(); ertType != "" {
		v := jaxb.EvidenceRecordTimestampTypeValue(ertType)
		xmlTimestampToken.EvidenceRecordTimestampType = &v
	}

	xmlTimestampToken.ProductionTime = jaxb.NewXSDateTime(timestampToken.GenerationTime())
	if filename := timestampToken.Filename(); filename != "" {
		xmlTimestampToken.TimestampFilename = &filename
	}
	xmlTimestampToken.DigestMatcher = append(xmlTimestampToken.DigestMatcher, b.getXmlDigestMatchersForTimestamp(timestampToken)...)
	xmlTimestampToken.BasicSignature = b.GetXmlBasicSignature(timestampToken)
	if infos := b.GetXmlSignerInformationStore(timestampToken.SignerInformationStoreInfos()); infos != nil {
		xmlTimestampToken.SignerInformationStore = &jaxb.SignerInformationStoreWrapper{Items: infos}
	}
	xmlTimestampToken.TSAGeneralName = b.getXmlTSAGeneralName(timestampToken)

	candidatesForSigningCertificate := timestampToken.CandidatesForSigningCertificate()
	theCertificateValidity := candidatesForSigningCertificate.TheCertificateValidity()
	if theCertificateValidity != nil {
		xmlTimestampToken.SigningCertificate = b.GetXmlSigningCertificateForIdentifier(timestampToken.DSSID(), theCertificateValidity)
		xmlTimestampToken.CertificateChain = b.certificateChainWrapper(b.GetXmlForCertificateChainForValidity(theCertificateValidity, timestampToken.CertificateSource()))
	}

	xmlTimestampToken.FoundCertificates = b.GetXmlFoundCertificatesForToken(timestampToken.DSSID(), (foundCertificatesSource)(timestampToken.CertificateSource()))
	xmlTimestampToken.FoundRevocations = b.getXmlFoundRevocations(timestampToken.CRLSource(), timestampToken.OCSPSource())

	if utils.IsCollectionNotEmpty(timestampToken.TimestampScopes()) {
		xmlTimestampToken.TimestampScopes = &jaxb.TimestampScopesWrapper{Items: b.getXmlSignatureScopes(timestampToken.TimestampScopes())}
	}

	if b.tokenExtractionStrategy.IsTimestamp() {
		bin := jaxb.Base64Binary(timestampToken.Encoded())
		xmlTimestampToken.Base64Encoded = &bin
	} else {
		tstDigest, err := timestampToken.Digest(b.defaultDigestAlgorithm)
		if err != nil {
			panic(err)
		}
		xmlTimestampToken.DigestAlgoAndValue = b.GetXmlDigestAlgoAndValueFor(b.defaultDigestAlgorithm, tstDigest)
	}

	b.xmlTimestampsMap[timestampToken.DSSIDAsString()] = xmlTimestampToken

	return xmlTimestampToken
}

func (b *SignedDocumentDiagnosticDataBuilder) checkDuplicatesTimestamp(xmlTimestamp *jaxb.XmlTimestamp, timestampToken *validation.TimestampToken) {
	if b.hasDuplicateTimestamp(timestampToken) {
		duplicated := true
		xmlTimestamp.Duplicated = &duplicated
	}
}

func (b *SignedDocumentDiagnosticDataBuilder) hasDuplicateTimestamp(currentTimestampToken *validation.TimestampToken) bool {
	for _, timestampToken := range b.usedTimestamps {
		if currentTimestampToken != timestampToken && currentTimestampToken.DSSIDAsString() == timestampToken.DSSIDAsString() {
			return true
		}
	}
	return false
}

func (b *SignedDocumentDiagnosticDataBuilder) getXmlDigestMatchersForTimestamp(timestampToken *validation.TimestampToken) []*jaxb.XmlDigestMatcher {
	digestMatchers := make([]*jaxb.XmlDigestMatcher, 0)
	digestMatchers = append(digestMatchers, b.getImprintDigestMatcher(timestampToken))
	digestMatchers = append(digestMatchers, b.getManifestEntriesDigestMatchers(timestampToken.ManifestFile())...)
	digestMatchers = append(digestMatchers, b.getTimestampReferenceValidationDigestMatchers(timestampToken.ReferenceValidations())...)
	return digestMatchers
}

func (b *SignedDocumentDiagnosticDataBuilder) getImprintDigestMatcher(timestampToken *validation.TimestampToken) *jaxb.XmlDigestMatcher {
	digestMatcher := &jaxb.XmlDigestMatcher{}
	t := jaxb.DigestMatcherTypeValue(enumerations.DigestMatcherType_MESSAGE_IMPRINT)
	digestMatcher.Type = &t
	messageImprint := timestampToken.MessageImprint()
	if !messageImprint.IsEmpty() {
		dm := jaxb.DigestAlgorithmValue(messageImprint.Algorithm())
		digestMatcher.DigestMethod = &dm
		v := jaxb.Base64Binary(messageImprint.Value())
		digestMatcher.DigestValue = &v
	}
	dataFound := timestampToken.IsMessageImprintDataFound()
	digestMatcher.DataFound = dataFound
	dataIntact := timestampToken.IsMessageImprintDataIntact()
	digestMatcher.DataIntact = dataIntact
	manifestFile := timestampToken.ManifestFile()
	if manifestFile != nil {
		filename := manifestFile.Filename()
		digestMatcher.DocumentName = &filename
	}
	return digestMatcher
}

func (b *SignedDocumentDiagnosticDataBuilder) getManifestEntriesDigestMatchers(manifestFile *model.ManifestFile) []*jaxb.XmlDigestMatcher {
	digestMatchers := make([]*jaxb.XmlDigestMatcher, 0)
	if manifestFile != nil && utils.IsCollectionNotEmpty(manifestFile.Entries()) {
		for _, entry := range manifestFile.Entries() {
			digestMatcher := &jaxb.XmlDigestMatcher{}
			t := jaxb.DigestMatcherTypeValue(enumerations.DigestMatcherType_MANIFEST_ENTRY)
			digestMatcher.Type = &t
			digest := entry.Digest()
			if !digest.IsEmpty() {
				dm := jaxb.DigestAlgorithmValue(digest.Algorithm())
				digestMatcher.DigestMethod = &dm
				v := jaxb.Base64Binary(digest.Value())
				digestMatcher.DigestValue = &v
			}
			dataFound := entry.IsFound()
			digestMatcher.DataFound = dataFound
			dataIntact := entry.IsIntact()
			digestMatcher.DataIntact = dataIntact
			if uri := entry.Uri(); uri != "" {
				digestMatcher.Uri = &uri
			}
			if entry.Document() != nil {
				docName := entry.Document().Name()
				digestMatcher.DocumentName = &docName
			}

			digestMatchers = append(digestMatchers, digestMatcher)
		}
	}
	return digestMatchers
}

// GetXmlSignerInformationStore builds a list of XmlSignerInfo from SignerIdentifiers. Port of
// the protected getXmlSignerInformationStore(Set<SignerIdentifier>).
func (b *SignedDocumentDiagnosticDataBuilder) GetXmlSignerInformationStore(signerIdentifiers []*spi.SignerIdentifier) []*jaxb.XmlSignerInfo {
	if utils.IsCollectionNotEmpty(signerIdentifiers) {
		signerInfos := make([]*jaxb.XmlSignerInfo, 0, len(signerIdentifiers))
		for _, signerIdentifier := range signerIdentifiers {
			signerInfos = append(signerInfos, b.GetXmlSignerInfo(signerIdentifier))
		}
		return signerInfos
	}
	return nil
}

func (b *SignedDocumentDiagnosticDataBuilder) getXmlTSAGeneralName(timestampToken *validation.TimestampToken) *jaxb.XmlTSAGeneralName {
	tstInfoTsa := timestampToken.TSTInfoTsa()
	if tstInfoTsa != nil {
		xmlTSAGeneralName := &jaxb.XmlTSAGeneralName{}

		x500PrincipalHelper := model.NewX500PrincipalHelper(tstInfoTsa)
		xmlTSAGeneralName.Value = x500PrincipalHelper.RFC2253()

		issuerX500Principal := timestampToken.IssuerX500Principal()
		if issuerX500Principal != nil {
			xmlTSAGeneralName.ContentMatch = spi.DSSASN1UtilsX500PrincipalAreEquals(tstInfoTsa, issuerX500Principal)
			xmlTSAGeneralName.OrderMatch = tstInfoTsa.Equals(issuerX500Principal)
		}

		return xmlTSAGeneralName
	}
	return nil
}

func (b *SignedDocumentDiagnosticDataBuilder) linkSignaturesAndTimestamps(signatures []validation.AdvancedSignature) {
	for _, advancedSignature := range signatures {
		currentSignature := b.xmlSignaturesMap[advancedSignature.ID()]
		// attach timestamps
		currentSignature.FoundTimestamps = &jaxb.FoundTimestampsWrapper{Items: b.getXmlFoundTimestamps(advancedSignature)}
	}
}

func (b *SignedDocumentDiagnosticDataBuilder) getXmlFoundTimestamps(sig validation.AdvancedSignature) []*jaxb.XmlFoundTimestamp {
	foundTimestamps := make([]*jaxb.XmlFoundTimestamp, 0)
	for _, timestampToken := range sig.AllTimestamps() {
		foundTimestamp := &jaxb.XmlFoundTimestamp{}
		foundTimestamp.Timestamp = b.xmlTimestampsMap[timestampToken.DSSIDAsString()]
		foundTimestamps = append(foundTimestamps, foundTimestamp)
	}
	return foundTimestamps
}

func (b *SignedDocumentDiagnosticDataBuilder) linkTimestampsAndTimestampsObjects(timestamps []*validation.TimestampToken) {
	for _, timestampToken := range timestamps {
		xmlTimestampToken := b.xmlTimestampsMap[timestampToken.DSSIDAsString()]
		xmlTimestampToken.TimestampedObjects = b.getXmlTimestampedObjects(timestampToken.TimestampedReferences())
	}
}

func (b *SignedDocumentDiagnosticDataBuilder) linkEvidenceRecordsAndTimestampsObjects(evidenceRecords []validation.EvidenceRecord) {
	for _, evidenceRecord := range evidenceRecords {
		xmlEvidenceRecord := b.xmlEvidenceRecordMap[evidenceRecord.Id()]
		xmlEvidenceRecord.TimestampedObjects = b.getXmlTimestampedObjects(evidenceRecord.TimestampedReferences())
	}
}

func (b *SignedDocumentDiagnosticDataBuilder) getXmlTimestampedObjects(timestampReferences []*validation.TimestampedReference) *jaxb.TimestampedObjectsWrapper {
	if utils.IsCollectionNotEmpty(timestampReferences) {
		objects := make([]*jaxb.XmlTimestampedObject, 0)
		addedTokenIds := map[string]bool{}
		for _, timestampReference := range timestampReferences {
			id := timestampReference.ObjectId()

			timestampedObject := b.createXmlTimestampedObject(timestampReference)
			if timestampedObject.Token == nil {
				panic("Token with Id '" + id + "' not found")
			}
			id = timestampedObject.Token.ID // can change in case of ref
			if addedTokenIds[id] {
				// skip the ref if it was added before
				continue
			}
			addedTokenIds[id] = true

			objects = append(objects, timestampedObject)
		}
		return &jaxb.TimestampedObjectsWrapper{Items: objects}
	}
	return nil
}

func (b *SignedDocumentDiagnosticDataBuilder) createXmlTimestampedObject(timestampReference *validation.TimestampedReference) *jaxb.XmlTimestampedObject {
	timestampedObj := &jaxb.XmlTimestampedObject{}
	category := timestampReference.Category()
	c := jaxb.TimestampedObjectTypeValue(category)
	timestampedObj.Category = &c

	objectID := timestampReference.ObjectId()

	switch category {
	case enumerations.TimestampedObjectType_SIGNATURE:
		timestampedObj.Token = jaxb.NewXmlTokenRef(b.xmlSignaturesMap[objectID])
		return timestampedObj

	case enumerations.TimestampedObjectType_CERTIFICATE:
		if !b.isUsedCertificate(objectID) {
			relatedCertificateID, ok := b.referenceMap[objectID]
			if ok {
				objectID = relatedCertificateID
				if !b.isUsedCertificate(objectID) {
					break // break to create an orphan token
				}
			} else {
				break
			}
		}
		timestampedObj.Token = jaxb.NewXmlTokenRef(b.xmlCertsMap[objectID])
		return timestampedObj

	case enumerations.TimestampedObjectType_REVOCATION:
		if !b.isUsedRevocation(objectID) {
			relatedRevocationID, ok := b.referenceMap[objectID]
			if ok {
				objectID = relatedRevocationID
				if !b.isUsedRevocation(objectID) {
					break // break to create an orphan token
				}
			} else {
				break
			}
		}
		timestampedObj.Token = jaxb.NewXmlTokenRef(b.xmlRevocationsMap[objectID])
		return timestampedObj

	case enumerations.TimestampedObjectType_TIMESTAMP:
		timestampedObj.Token = jaxb.NewXmlTokenRef(b.xmlTimestampsMap[objectID])
		return timestampedObj

	case enumerations.TimestampedObjectType_EVIDENCE_RECORD:
		timestampedObj.Token = jaxb.NewXmlTokenRef(b.xmlEvidenceRecordMap[objectID])
		return timestampedObj

	case enumerations.TimestampedObjectType_SIGNED_DATA:
		timestampedObj.Token = jaxb.NewXmlTokenRef(b.xmlSignedDataMap[objectID])
		return timestampedObj

	default:
		panic("Unsupported category '" + string(category) + "'")
	}

	switch enumerations.TimestampedObjectType(timestampedObj.Category.TimestampedObjectType()) {
	case enumerations.TimestampedObjectType_CERTIFICATE:
		timestampedObj.Token = jaxb.NewXmlTokenRef(b.xmlOrphanCertificateTokensMap[objectID])
		orphanCategory := jaxb.TimestampedObjectTypeValue(enumerations.TimestampedObjectType_ORPHAN_CERTIFICATE)
		timestampedObj.Category = &orphanCategory

	case enumerations.TimestampedObjectType_REVOCATION:
		timestampedObj.Token = jaxb.NewXmlTokenRef(b.xmlOrphanRevocationTokensMap[objectID])
		orphanCategory := jaxb.TimestampedObjectTypeValue(enumerations.TimestampedObjectType_ORPHAN_REVOCATION)
		timestampedObj.Category = &orphanCategory

	default:
		panic("The type of object [" + string(timestampedObj.Category.TimestampedObjectType()) + "] is not supported for Orphan Tokens!")
	}

	return timestampedObj
}

func (b *SignedDocumentDiagnosticDataBuilder) isUsedCertificate(tokenID string) bool {
	for _, token := range b.usedCertificates {
		if token.DSSIDAsString() == tokenID {
			return true
		}
	}
	return false
}

func (b *SignedDocumentDiagnosticDataBuilder) isUsedRevocation(tokenID string) bool {
	for _, token := range b.usedRevocations {
		if token.DSSIDAsString() == tokenID {
			return true
		}
	}
	return false
}
