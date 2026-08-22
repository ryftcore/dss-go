// Ported from DiagnosticData.xsd (DSS 6.5.RC1) via the JAXB classes generated into
// eu.europa.esig.dss.diagnostic.jaxb. The generated classes are grouped into
// schema-area files rather than one file per class; every Java class keeps its
// name and its exact field order so that encoding/xml reproduces the JAXB element
// sequence byte for byte.

package jaxb

// The wrapper types below reproduce JAXB's @XmlElementWrapper properties.
// A nil wrapper means the property was absent; a non-nil wrapper with no
// items reproduces JAXB's empty `<Wrapper/>` element, which a plain Go
// `parent>child` tag cannot express.

// AdditionalServiceInfoUrisWrapper wraps the <AdditionalServiceInfoUri> elements of the <AdditionalServiceInfoUris> property.
type AdditionalServiceInfoUrisWrapper struct {
	Items []string `xml:"AdditionalServiceInfoUri"`
}

// All returns the wrapped items, tolerating a nil wrapper.
func (w *AdditionalServiceInfoUrisWrapper) All() []string {
	if w == nil {
		return nil
	}
	return w.Items
}

// CapturedQualifiersWrapper wraps the <Qualifier> elements of the <CapturedQualifiers> property.
type CapturedQualifiersWrapper struct {
	Items []*XmlQualifier `xml:"Qualifier"`
}

// All returns the wrapped items, tolerating a nil wrapper.
func (w *CapturedQualifiersWrapper) All() []*XmlQualifier {
	if w == nil {
		return nil
	}
	return w.Items
}

// CertificateChainWrapper wraps the <ChainItem> elements of the <CertificateChain> property.
type CertificateChainWrapper struct {
	Items []*XmlChainItem `xml:"ChainItem"`
}

// All returns the wrapped items, tolerating a nil wrapper.
func (w *CertificateChainWrapper) All() []*XmlChainItem {
	if w == nil {
		return nil
	}
	return w.Items
}

// CertificateContentEquivalenceListWrapper wraps the <CertificateContentEquivalence> elements of the <CertificateContentEquivalenceList> property.
type CertificateContentEquivalenceListWrapper struct {
	Items []*XmlCertificateContentEquivalence `xml:"CertificateContentEquivalence"`
}

// All returns the wrapped items, tolerating a nil wrapper.
func (w *CertificateContentEquivalenceListWrapper) All() []*XmlCertificateContentEquivalence {
	if w == nil {
		return nil
	}
	return w.Items
}

// CommitmentTypeIndicationsWrapper wraps the <CommitmentTypeIndication> elements of the <CommitmentTypeIndications> property.
type CommitmentTypeIndicationsWrapper struct {
	Items []*XmlCommitmentTypeIndication `xml:"CommitmentTypeIndication"`
}

// All returns the wrapped items, tolerating a nil wrapper.
func (w *CommitmentTypeIndicationsWrapper) All() []*XmlCommitmentTypeIndication {
	if w == nil {
		return nil
	}
	return w.Items
}

// ContentFilesWrapper wraps the <ContentFile> elements of the <ContentFiles> property.
type ContentFilesWrapper struct {
	Items []string `xml:"ContentFile"`
}

// All returns the wrapped items, tolerating a nil wrapper.
func (w *ContentFilesWrapper) All() []string {
	if w == nil {
		return nil
	}
	return w.Items
}

// DataObjectReferencesWrapper wraps the <ObjectReference> elements of the <DataObjectReferences> property.
type DataObjectReferencesWrapper struct {
	Items []string `xml:"ObjectReference"`
}

// All returns the wrapped items, tolerating a nil wrapper.
func (w *DataObjectReferencesWrapper) All() []string {
	if w == nil {
		return nil
	}
	return w.Items
}

// DigestMatchersWrapper wraps the <DigestMatcher> elements of the <DigestMatchers> property.
type DigestMatchersWrapper struct {
	Items []*XmlDigestMatcher `xml:"DigestMatcher"`
}

// All returns the wrapped items, tolerating a nil wrapper.
func (w *DigestMatchersWrapper) All() []*XmlDigestMatcher {
	if w == nil {
		return nil
	}
	return w.Items
}

// DocumentationReferencesWrapper wraps the <DocumentationReference> elements of the <DocumentationReferences> property.
type DocumentationReferencesWrapper struct {
	Items []string `xml:"DocumentationReference"`
}

// All returns the wrapped items, tolerating a nil wrapper.
func (w *DocumentationReferencesWrapper) All() []string {
	if w == nil {
		return nil
	}
	return w.Items
}

// DocumentsWrapper wraps the <EAADocument> elements of the <Documents> property.
type DocumentsWrapper struct {
	Items []*XmlEAADocument `xml:"EAADocument"`
}

// All returns the wrapped items, tolerating a nil wrapper.
func (w *DocumentsWrapper) All() []*XmlEAADocument {
	if w == nil {
		return nil
	}
	return w.Items
}

// EAARevocationsWrapper wraps the <EAARevocationStatus> elements of the <EAARevocations> property.
type EAARevocationsWrapper struct {
	Items []*XmlEAARevocationStatus `xml:"EAARevocationStatus"`
}

// All returns the wrapped items, tolerating a nil wrapper.
func (w *EAARevocationsWrapper) All() []*XmlEAARevocationStatus {
	if w == nil {
		return nil
	}
	return w.Items
}

// EAAsWrapper wraps the <EAA> elements of the <EAAs> property.
type EAAsWrapper struct {
	Items []*XmlEAA `xml:"EAA"`
}

// All returns the wrapped items, tolerating a nil wrapper.
func (w *EAAsWrapper) All() []*XmlEAA {
	if w == nil {
		return nil
	}
	return w.Items
}

// EntriesWrapper wraps the <Entry> elements of the <Entries> property.
type EntriesWrapper struct {
	Items []string `xml:"Entry"`
}

// All returns the wrapped items, tolerating a nil wrapper.
func (w *EntriesWrapper) All() []string {
	if w == nil {
		return nil
	}
	return w.Items
}

// EvidenceRecordScopesWrapper wraps the <SignatureScope> elements of the <EvidenceRecordScopes> property.
type EvidenceRecordScopesWrapper struct {
	Items []*XmlSignatureScope `xml:"SignatureScope"`
}

// All returns the wrapped items, tolerating a nil wrapper.
func (w *EvidenceRecordScopesWrapper) All() []*XmlSignatureScope {
	if w == nil {
		return nil
	}
	return w.Items
}

// EvidenceRecordTimestampsWrapper wraps the <FoundTimestamp> elements of the <EvidenceRecordTimestamps> property.
type EvidenceRecordTimestampsWrapper struct {
	Items []*XmlFoundTimestamp `xml:"FoundTimestamp"`
}

// All returns the wrapped items, tolerating a nil wrapper.
func (w *EvidenceRecordTimestampsWrapper) All() []*XmlFoundTimestamp {
	if w == nil {
		return nil
	}
	return w.Items
}

// EvidenceRecordsWrapper wraps the <EvidenceRecord> elements of the <EvidenceRecords> property.
type EvidenceRecordsWrapper struct {
	Items []*XmlEvidenceRecord `xml:"EvidenceRecord"`
}

// All returns the wrapped items, tolerating a nil wrapper.
func (w *EvidenceRecordsWrapper) All() []*XmlEvidenceRecord {
	if w == nil {
		return nil
	}
	return w.Items
}

// FoundEvidenceRecordsWrapper wraps the <FoundEvidenceRecord> elements of the <FoundEvidenceRecords> property.
type FoundEvidenceRecordsWrapper struct {
	Items []*XmlFoundEvidenceRecord `xml:"FoundEvidenceRecord"`
}

// All returns the wrapped items, tolerating a nil wrapper.
func (w *FoundEvidenceRecordsWrapper) All() []*XmlFoundEvidenceRecord {
	if w == nil {
		return nil
	}
	return w.Items
}

// FoundTimestampsWrapper wraps the <FoundTimestamp> elements of the <FoundTimestamps> property.
type FoundTimestampsWrapper struct {
	Items []*XmlFoundTimestamp `xml:"FoundTimestamp"`
}

// All returns the wrapped items, tolerating a nil wrapper.
func (w *FoundTimestampsWrapper) All() []*XmlFoundTimestamp {
	if w == nil {
		return nil
	}
	return w.Items
}

// ListsOfTrustedEntitiesWrapper wraps the <ListOfTrustedEntities> elements of the <ListsOfTrustedEntities> property.
type ListsOfTrustedEntitiesWrapper struct {
	Items []*XmlListOfTrustedEntities `xml:"ListOfTrustedEntities"`
}

// All returns the wrapped items, tolerating a nil wrapper.
func (w *ListsOfTrustedEntitiesWrapper) All() []*XmlListOfTrustedEntities {
	if w == nil {
		return nil
	}
	return w.Items
}

// ManifestFilesWrapper wraps the <ManifestFile> elements of the <ManifestFiles> property.
type ManifestFilesWrapper struct {
	Items []*XmlManifestFile `xml:"ManifestFile"`
}

// All returns the wrapped items, tolerating a nil wrapper.
func (w *ManifestFilesWrapper) All() []*XmlManifestFile {
	if w == nil {
		return nil
	}
	return w.Items
}

// NamesWrapper wraps the <Name> elements of the <Names> property.
type NamesWrapper struct {
	Items []*XmlLangAndValue `xml:"Name"`
}

// All returns the wrapped items, tolerating a nil wrapper.
func (w *NamesWrapper) All() []*XmlLangAndValue {
	if w == nil {
		return nil
	}
	return w.Items
}

// ObjectReferencesWrapper wraps the <ObjectReference> elements of the <ObjectReferences> property.
type ObjectReferencesWrapper struct {
	Items []string `xml:"ObjectReference"`
}

// All returns the wrapped items, tolerating a nil wrapper.
func (w *ObjectReferencesWrapper) All() []string {
	if w == nil {
		return nil
	}
	return w.Items
}

// OriginalDocumentsWrapper wraps the <SignerData> elements of the <OriginalDocuments> property.
type OriginalDocumentsWrapper struct {
	Items []*XmlSignerData `xml:"SignerData"`
}

// All returns the wrapped items, tolerating a nil wrapper.
func (w *OriginalDocumentsWrapper) All() []*XmlSignerData {
	if w == nil {
		return nil
	}
	return w.Items
}

// OtherOIDsWrapper wraps the <OtherOID> elements of the <OtherOIDs> property.
type OtherOIDsWrapper struct {
	Items []*XmlOID `xml:"OtherOID"`
}

// All returns the wrapped items, tolerating a nil wrapper.
func (w *OtherOIDsWrapper) All() []*XmlOID {
	if w == nil {
		return nil
	}
	return w.Items
}

// QcCClegislationWrapper wraps the <CountryName> elements of the <QcCClegislation> property.
type QcCClegislationWrapper struct {
	Items []string `xml:"CountryName"`
}

// All returns the wrapped items, tolerating a nil wrapper.
func (w *QcCClegislationWrapper) All() []string {
	if w == nil {
		return nil
	}
	return w.Items
}

// QcEuPDSWrapper wraps the <PdsLocation> elements of the <QcEuPDS> property.
type QcEuPDSWrapper struct {
	Items []*XmlLangAndValue `xml:"PdsLocation"`
}

// All returns the wrapped items, tolerating a nil wrapper.
func (w *QcEuPDSWrapper) All() []*XmlLangAndValue {
	if w == nil {
		return nil
	}
	return w.Items
}

// QcQSCDlegislationWrapper wraps the <QSCDCountryName> elements of the <QcQSCDlegislation> property.
type QcQSCDlegislationWrapper struct {
	Items []string `xml:"QSCDCountryName"`
}

// All returns the wrapped items, tolerating a nil wrapper.
func (w *QcQSCDlegislationWrapper) All() []string {
	if w == nil {
		return nil
	}
	return w.Items
}

// QcTypesWrapper wraps the <QcType> elements of the <QcTypes> property.
type QcTypesWrapper struct {
	Items []*XmlOID `xml:"QcType"`
}

// All returns the wrapped items, tolerating a nil wrapper.
func (w *QcTypesWrapper) All() []*XmlOID {
	if w == nil {
		return nil
	}
	return w.Items
}

// RegistrationIdentifiersWrapper wraps the <RegistrationIdentifier> elements of the <RegistrationIdentifiers> property.
type RegistrationIdentifiersWrapper struct {
	Items []string `xml:"RegistrationIdentifier"`
}

// All returns the wrapped items, tolerating a nil wrapper.
func (w *RegistrationIdentifiersWrapper) All() []string {
	if w == nil {
		return nil
	}
	return w.Items
}

// RevocationsWrapper wraps the <CertificateRevocation> elements of the <Revocations> property.
type RevocationsWrapper struct {
	Items []*XmlCertificateRevocation `xml:"CertificateRevocation"`
}

// All returns the wrapped items, tolerating a nil wrapper.
func (w *RevocationsWrapper) All() []*XmlCertificateRevocation {
	if w == nil {
		return nil
	}
	return w.Items
}

// RolesOfPSPWrapper wraps the <RoleOfPSP> elements of the <RolesOfPSP> property.
type RolesOfPSPWrapper struct {
	Items []*XmlRoleOfPSP `xml:"RoleOfPSP"`
}

// All returns the wrapped items, tolerating a nil wrapper.
func (w *RolesOfPSPWrapper) All() []*XmlRoleOfPSP {
	if w == nil {
		return nil
	}
	return w.Items
}

// ServiceNamesWrapper wraps the <ServiceName> elements of the <ServiceNames> property.
type ServiceNamesWrapper struct {
	Items []*XmlLangAndValue `xml:"ServiceName"`
}

// All returns the wrapped items, tolerating a nil wrapper.
func (w *ServiceNamesWrapper) All() []*XmlLangAndValue {
	if w == nil {
		return nil
	}
	return w.Items
}

// ServiceSupplyPointsWrapper wraps the <ServiceSupplyPoint> elements of the <ServiceSupplyPoints> property.
type ServiceSupplyPointsWrapper struct {
	Items []string `xml:"ServiceSupplyPoint"`
}

// All returns the wrapped items, tolerating a nil wrapper.
func (w *ServiceSupplyPointsWrapper) All() []string {
	if w == nil {
		return nil
	}
	return w.Items
}

// SignatureScopesWrapper wraps the <SignatureScope> elements of the <SignatureScopes> property.
type SignatureScopesWrapper struct {
	Items []*XmlSignatureScope `xml:"SignatureScope"`
}

// All returns the wrapped items, tolerating a nil wrapper.
func (w *SignatureScopesWrapper) All() []*XmlSignatureScope {
	if w == nil {
		return nil
	}
	return w.Items
}

// SignaturesWrapper wraps the <Signature> elements of the <Signatures> property.
type SignaturesWrapper struct {
	Items []*XmlSignature `xml:"Signature"`
}

// All returns the wrapped items, tolerating a nil wrapper.
func (w *SignaturesWrapper) All() []*XmlSignature {
	if w == nil {
		return nil
	}
	return w.Items
}

// SignerInformationStoreWrapper wraps the <SignerInfo> elements of the <SignerInformationStore> property.
type SignerInformationStoreWrapper struct {
	Items []*XmlSignerInfo `xml:"SignerInfo"`
}

// All returns the wrapped items, tolerating a nil wrapper.
func (w *SignerInformationStoreWrapper) All() []*XmlSignerInfo {
	if w == nil {
		return nil
	}
	return w.Items
}

// SourcesWrapper wraps the <Source> elements of the <Sources> property.
type SourcesWrapper struct {
	Items []CertificateSourceTypeValue `xml:"Source"`
}

// All returns the wrapped items, tolerating a nil wrapper.
func (w *SourcesWrapper) All() []CertificateSourceTypeValue {
	if w == nil {
		return nil
	}
	return w.Items
}

// TSPNamesWrapper wraps the <TSPName> elements of the <TSPNames> property.
type TSPNamesWrapper struct {
	Items []*XmlLangAndValue `xml:"TSPName"`
}

// All returns the wrapped items, tolerating a nil wrapper.
func (w *TSPNamesWrapper) All() []*XmlLangAndValue {
	if w == nil {
		return nil
	}
	return w.Items
}

// TSPRegistrationIdentifiersWrapper wraps the <TSPRegistrationIdentifier> elements of the <TSPRegistrationIdentifiers> property.
type TSPRegistrationIdentifiersWrapper struct {
	Items []string `xml:"TSPRegistrationIdentifier"`
}

// All returns the wrapped items, tolerating a nil wrapper.
func (w *TSPRegistrationIdentifiersWrapper) All() []string {
	if w == nil {
		return nil
	}
	return w.Items
}

// TSPTradeNamesWrapper wraps the <TSPTradeName> elements of the <TSPTradeNames> property.
type TSPTradeNamesWrapper struct {
	Items []*XmlLangAndValue `xml:"TSPTradeName"`
}

// All returns the wrapped items, tolerating a nil wrapper.
func (w *TSPTradeNamesWrapper) All() []*XmlLangAndValue {
	if w == nil {
		return nil
	}
	return w.Items
}

// TimestampScopesWrapper wraps the <TimestampScope> elements of the <TimestampScopes> property.
type TimestampScopesWrapper struct {
	Items []*XmlSignatureScope `xml:"TimestampScope"`
}

// All returns the wrapped items, tolerating a nil wrapper.
func (w *TimestampScopesWrapper) All() []*XmlSignatureScope {
	if w == nil {
		return nil
	}
	return w.Items
}

// TimestampedObjectsWrapper wraps the <TimestampedObject> elements of the <TimestampedObjects> property.
type TimestampedObjectsWrapper struct {
	Items []*XmlTimestampedObject `xml:"TimestampedObject"`
}

// All returns the wrapped items, tolerating a nil wrapper.
func (w *TimestampedObjectsWrapper) All() []*XmlTimestampedObject {
	if w == nil {
		return nil
	}
	return w.Items
}

// TradeNamesWrapper wraps the <TradeName> elements of the <TradeNames> property.
type TradeNamesWrapper struct {
	Items []*XmlLangAndValue `xml:"TradeName"`
}

// All returns the wrapped items, tolerating a nil wrapper.
func (w *TradeNamesWrapper) All() []*XmlLangAndValue {
	if w == nil {
		return nil
	}
	return w.Items
}

// TransformationsWrapper wraps the <Transformation> elements of the <Transformations> property.
type TransformationsWrapper struct {
	Items []string `xml:"Transformation"`
}

// All returns the wrapped items, tolerating a nil wrapper.
func (w *TransformationsWrapper) All() []string {
	if w == nil {
		return nil
	}
	return w.Items
}

// TrustServiceProvidersWrapper wraps the <TrustServiceProvider> elements of the <TrustServiceProviders> property.
type TrustServiceProvidersWrapper struct {
	Items []*XmlTrustServiceProvider `xml:"TrustServiceProvider"`
}

// All returns the wrapped items, tolerating a nil wrapper.
func (w *TrustServiceProvidersWrapper) All() []*XmlTrustServiceProvider {
	if w == nil {
		return nil
	}
	return w.Items
}

// TrustServicesWrapper wraps the <TrustService> elements of the <TrustServices> property.
type TrustServicesWrapper struct {
	Items []*XmlTrustService `xml:"TrustService"`
}

// All returns the wrapped items, tolerating a nil wrapper.
func (w *TrustServicesWrapper) All() []*XmlTrustService {
	if w == nil {
		return nil
	}
	return w.Items
}

// TrustedEntitiesWrapper wraps the <TrustedEntity> elements of the <TrustedEntities> property.
type TrustedEntitiesWrapper struct {
	Items []*XmlTrustedEntity `xml:"TrustedEntity"`
}

// All returns the wrapped items, tolerating a nil wrapper.
func (w *TrustedEntitiesWrapper) All() []*XmlTrustedEntity {
	if w == nil {
		return nil
	}
	return w.Items
}

// TrustedEntityServicesWrapper wraps the <TrustedEntityService> elements of the <TrustedEntityServices> property.
type TrustedEntityServicesWrapper struct {
	Items []*XmlTrustedEntityService `xml:"TrustedEntityService"`
}

// All returns the wrapped items, tolerating a nil wrapper.
func (w *TrustedEntityServicesWrapper) All() []*XmlTrustedEntityService {
	if w == nil {
		return nil
	}
	return w.Items
}

// TrustedListsWrapper wraps the <TrustedList> elements of the <TrustedLists> property.
type TrustedListsWrapper struct {
	Items []*XmlTrustedList `xml:"TrustedList"`
}

// All returns the wrapped items, tolerating a nil wrapper.
func (w *TrustedListsWrapper) All() []*XmlTrustedList {
	if w == nil {
		return nil
	}
	return w.Items
}

// UsedCertificatesWrapper wraps the <Certificate> elements of the <UsedCertificates> property.
type UsedCertificatesWrapper struct {
	Items []*XmlCertificate `xml:"Certificate"`
}

// All returns the wrapped items, tolerating a nil wrapper.
func (w *UsedCertificatesWrapper) All() []*XmlCertificate {
	if w == nil {
		return nil
	}
	return w.Items
}

// UsedEAARevocationTokensWrapper wraps the <EAARevocationToken> elements of the <UsedEAARevocationTokens> property.
type UsedEAARevocationTokensWrapper struct {
	Items []*XmlEAARevocationToken `xml:"EAARevocationToken"`
}

// All returns the wrapped items, tolerating a nil wrapper.
func (w *UsedEAARevocationTokensWrapper) All() []*XmlEAARevocationToken {
	if w == nil {
		return nil
	}
	return w.Items
}

// UsedRevocationsWrapper wraps the <Revocation> elements of the <UsedRevocations> property.
type UsedRevocationsWrapper struct {
	Items []*XmlRevocation `xml:"Revocation"`
}

// All returns the wrapped items, tolerating a nil wrapper.
func (w *UsedRevocationsWrapper) All() []*XmlRevocation {
	if w == nil {
		return nil
	}
	return w.Items
}

// UsedTimestampsWrapper wraps the <Timestamp> elements of the <UsedTimestamps> property.
type UsedTimestampsWrapper struct {
	Items []*XmlTimestamp `xml:"Timestamp"`
}

// All returns the wrapped items, tolerating a nil wrapper.
func (w *UsedTimestampsWrapper) All() []*XmlTimestamp {
	if w == nil {
		return nil
	}
	return w.Items
}

// ValidationMessagesWrapper wraps the <Error> elements of the <ValidationMessages> property.
type ValidationMessagesWrapper struct {
	Items []string `xml:"Error"`
}

// All returns the wrapped items, tolerating a nil wrapper.
func (w *ValidationMessagesWrapper) All() []string {
	if w == nil {
		return nil
	}
	return w.Items
}
