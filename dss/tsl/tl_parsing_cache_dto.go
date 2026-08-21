// Ported from dss-tsl-validation/src/main/java/eu/europa/esig/dss/tsl/dto/TLParsingCacheDTO.java (DSS 6.5.RC1).
//
// CROSS-CHUNK DEPENDENCY (see this batch's porter notes): this file embeds
// eu.europa.esig.dss.validation.job.dto.AbstractParsingCacheDTO, ported by the dss-validation-job
// chunk into Go package dss/validation/job, which had not landed when this file was written. The
// embedded name and its two constructors follow PORTING.md's "exported Go identifiers keep the
// Java name" rule literally (AbstractParsingCacheDTO / NewAbstractParsingCacheDTO /
// NewAbstractParsingCacheDTOFromCacheDTO for the AbstractCacheDTO copy constructor); if that
// chunk chose different spellings, the fix is a mechanical rename here.
//
// DEAD FIELD: Java re-declares `protected List<String> structureValidationMessages` on this class,
// shadowing the identically named field its superclass already declares. Nothing ever writes the
// shadowing copy (the inherited setter writes the superclass one, and the inherited getter reads
// the superclass one), so it is dropped rather than reproduced - a shadowed Go field would change
// which copy the promoted accessors see.
package tsl

import (
	"time"

	"github.com/utain/esig/dss/enumerations"
	tslmodel "github.com/utain/esig/dss/model/tsl"
	"github.com/utain/esig/dss/utils"
	"github.com/utain/esig/dss/validation/job"
)

// TLParsingCacheDTO is the DTO for a Trusted List parsing cache. It implements
// eu.europa.esig.dss.model.tsl.TLParsingInfoRecord (Go tslmodel.TLParsingInfoRecord).
//
// java.io.Serializable has no Go counterpart and is dropped.
type TLParsingCacheDTO struct {
	*job.AbstractParsingCacheDTO

	// tslType is the LOTL/TL TSLType.
	tslType enumerations.TSLType

	// sequenceNumber is the LOTL/TL sequence number. Nil stands in for Java's null Integer.
	sequenceNumber *int

	// version is the LOTL/TL version. Nil stands in for Java's null Integer.
	version *int

	// territory is the country (territory).
	territory string

	// issueDate is the issuance date of the LOTL/TL. The zero time.Time stands in for Java's
	// null Date.
	issueDate time.Time

	// nextUpdateDate is the next update date. The zero time.Time stands in for Java's null
	// Date.
	nextUpdateDate time.Time

	// distributionPoints holds the distribution points urls.
	distributionPoints []string

	// trustServiceProviders is the list of found trust service providers.
	trustServiceProviders []*tslmodel.TrustServiceProvider

	// lotlOtherPointers is the list of LOTL other pointers.
	lotlOtherPointers []*tslmodel.OtherTSLPointer

	// tlOtherPointers is the list of TL other pointers.
	tlOtherPointers []*tslmodel.OtherTSLPointer

	// pivotUrls is the list of pivot URLs.
	pivotUrls []string

	// signingCertificateAnnouncementUrl is the signing certificate announcement URL.
	signingCertificateAnnouncementUrl string
}

var _ tslmodel.TLParsingInfoRecord = (*TLParsingCacheDTO)(nil)

// NewTLParsingCacheDTO is the default constructor. Port of TLParsingCacheDTO().
func NewTLParsingCacheDTO() *TLParsingCacheDTO {
	return &TLParsingCacheDTO{AbstractParsingCacheDTO: job.NewAbstractParsingCacheDTO()}
}

// NewTLParsingCacheDTOFromCacheDTO copies the cache DTO. Port of
// TLParsingCacheDTO(AbstractCacheDTO).
func NewTLParsingCacheDTOFromCacheDTO(cacheDTO *job.AbstractCacheDTO) *TLParsingCacheDTO {
	return &TLParsingCacheDTO{AbstractParsingCacheDTO: job.NewAbstractParsingCacheDTOFrom(cacheDTO)}
}

// TSLType gets the TSLType. Port of getTSLType().
func (d *TLParsingCacheDTO) TSLType() enumerations.TSLType {
	return d.tslType
}

// SetTSLType sets the TSLType. Port of setTSLType(TSLType).
func (d *TLParsingCacheDTO) SetTSLType(tslType enumerations.TSLType) {
	d.tslType = tslType
}

// SequenceNumber gets the sequence number. Port of getSequenceNumber().
func (d *TLParsingCacheDTO) SequenceNumber() *int {
	return d.sequenceNumber
}

// SetSequenceNumber sets the sequence number. Port of setSequenceNumber(Integer).
func (d *TLParsingCacheDTO) SetSequenceNumber(sequenceNumber *int) {
	d.sequenceNumber = sequenceNumber
}

// Version gets the version. Port of getVersion().
func (d *TLParsingCacheDTO) Version() *int {
	return d.version
}

// SetVersion sets the version. Port of setVersion(Integer).
func (d *TLParsingCacheDTO) SetVersion(version *int) {
	d.version = version
}

// Territory gets the territory. Port of getTerritory().
func (d *TLParsingCacheDTO) Territory() string {
	return d.territory
}

// SetTerritory sets the territory. Port of setTerritory(String).
func (d *TLParsingCacheDTO) SetTerritory(territory string) {
	d.territory = territory
}

// IssueDate gets the issue date. Port of getIssueDate().
func (d *TLParsingCacheDTO) IssueDate() time.Time {
	return d.issueDate
}

// SetIssueDate sets the issue date. Port of setIssueDate(Date).
func (d *TLParsingCacheDTO) SetIssueDate(issueDate time.Time) {
	d.issueDate = issueDate
}

// NextUpdateDate gets the next update date. Port of getNextUpdateDate().
func (d *TLParsingCacheDTO) NextUpdateDate() time.Time {
	return d.nextUpdateDate
}

// SetNextUpdateDate sets the next update date. Port of setNextUpdateDate(Date).
func (d *TLParsingCacheDTO) SetNextUpdateDate(nextUpdateDate time.Time) {
	d.nextUpdateDate = nextUpdateDate
}

// DistributionPoints gets the list of distribution point urls. Port of getDistributionPoints().
func (d *TLParsingCacheDTO) DistributionPoints() []string {
	return d.distributionPoints
}

// SetDistributionPoints sets the list of distribution point urls. Port of
// setDistributionPoints(List).
func (d *TLParsingCacheDTO) SetDistributionPoints(distributionPoints []string) {
	d.distributionPoints = distributionPoints
}

// TrustServiceProviders gets the trust service providers. Port of getTrustServiceProviders().
func (d *TLParsingCacheDTO) TrustServiceProviders() []*tslmodel.TrustServiceProvider {
	return d.trustServiceProviders
}

// SetTrustServiceProviders sets the trust service providers. Port of
// setTrustServiceProviders(List).
func (d *TLParsingCacheDTO) SetTrustServiceProviders(trustServiceProviders []*tslmodel.TrustServiceProvider) {
	d.trustServiceProviders = trustServiceProviders
}

// LotlOtherPointers gets the LOTL other pointers. Port of getLotlOtherPointers().
func (d *TLParsingCacheDTO) LotlOtherPointers() []*tslmodel.OtherTSLPointer {
	return d.lotlOtherPointers
}

// SetLotlOtherPointers sets the LOTL other pointers. Port of setLotlOtherPointers(List).
func (d *TLParsingCacheDTO) SetLotlOtherPointers(lotlOtherPointers []*tslmodel.OtherTSLPointer) {
	d.lotlOtherPointers = lotlOtherPointers
}

// TlOtherPointers gets the TL other pointers. Port of getTlOtherPointers().
func (d *TLParsingCacheDTO) TlOtherPointers() []*tslmodel.OtherTSLPointer {
	return d.tlOtherPointers
}

// SetTlOtherPointers sets the TL other pointers. Port of setTlOtherPointers(List).
func (d *TLParsingCacheDTO) SetTlOtherPointers(tlOtherPointers []*tslmodel.OtherTSLPointer) {
	d.tlOtherPointers = tlOtherPointers
}

// PivotUrls gets the pivot URLs. Port of getPivotUrls().
func (d *TLParsingCacheDTO) PivotUrls() []string {
	return d.pivotUrls
}

// SetPivotUrls sets the pivot URLs. Port of setPivotUrls(List).
func (d *TLParsingCacheDTO) SetPivotUrls(pivotUrls []string) {
	d.pivotUrls = pivotUrls
}

// SigningCertificateAnnouncementUrl gets the signing certificate announcement URL. Port of
// getSigningCertificateAnnouncementUrl().
func (d *TLParsingCacheDTO) SigningCertificateAnnouncementUrl() string {
	return d.signingCertificateAnnouncementUrl
}

// SetSigningCertificateAnnouncementUrl sets the signing certificate announcement URL. Port of
// setSigningCertificateAnnouncementUrl(String).
func (d *TLParsingCacheDTO) SetSigningCertificateAnnouncementUrl(signingCertificateAnnouncementUrl string) {
	d.signingCertificateAnnouncementUrl = signingCertificateAnnouncementUrl
}

// TSPNumber returns the number of all TrustServiceProviders present in the TL. Port of
// getTSPNumber().
func (d *TLParsingCacheDTO) TSPNumber() int {
	if utils.IsCollectionNotEmpty(d.trustServiceProviders) {
		return len(d.trustServiceProviders)
	}
	return 0
}

// TSNumber returns the number of all TrustServices present in the TL. Port of getTSNumber().
func (d *TLParsingCacheDTO) TSNumber() int {
	tsNumber := 0
	if utils.IsCollectionNotEmpty(d.trustServiceProviders) {
		for _, tsp := range d.trustServiceProviders {
			tsNumber += len(tsp.Services())
		}
	}
	return tsNumber
}

// CertNumber returns the number of all CertificateTokens present in the TL. Port of
// getCertNumber().
func (d *TLParsingCacheDTO) CertNumber() int {
	certNumber := 0
	if utils.IsCollectionNotEmpty(d.trustServiceProviders) {
		for _, tsp := range d.trustServiceProviders {
			for _, trustService := range tsp.Services() {
				certNumber += len(trustService.Certificates())
			}
		}
	}
	return certNumber
}
