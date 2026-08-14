// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/validation/timestamp/DetachedTimestampSource.java (DSS 6.5.RC1).
package timestamp

import (
	"github.com/utain/esig/dss/model/x509/revocation"
	"github.com/utain/esig/dss/spi"
	"github.com/utain/esig/dss/spi/validation"
)

// DetachedTimestampSource performs processing of detached timestamps.
type DetachedTimestampSource struct {
	AbstractTimestampSource

	// certificateSource is the merged certificate source from timestamps.
	certificateSource *spi.ListCertificateSource

	// crlSource is the merged CRL source.
	crlSource *spi.ListRevocationSource[revocation.CRL]

	// ocspSource is the merged OCSP source.
	ocspSource *spi.ListRevocationSource[revocation.OCSP]

	// detachedTimestamps is a list of detached timestamps.
	detachedTimestamps []*validation.TimestampToken

	// detachedEvidenceRecords contains the list of evidence records detached from the
	// time-stamp document.
	detachedEvidenceRecords []validation.EvidenceRecord
}

// NewDetachedTimestampSource instantiates an object with empty resources.
// Port of the default constructor.
func NewDetachedTimestampSource() *DetachedTimestampSource {
	return &DetachedTimestampSource{
		certificateSource: spi.NewListCertificateSource(),
		crlSource:         spi.NewListRevocationSource[revocation.CRL](),
		ocspSource:        spi.NewListRevocationSource[revocation.OCSP](),
	}
}

// NewDetachedTimestampSourceWithTimestamp instantiates a list of time-stamps with the given
// TimestampToken. Port of the DetachedTimestampSource(TimestampToken) constructor.
func NewDetachedTimestampSourceWithTimestamp(timestampToken *validation.TimestampToken) *DetachedTimestampSource {
	source := NewDetachedTimestampSource()
	source.detachedTimestamps = append(source.detachedTimestamps, timestampToken)
	return source
}

// DetachedTimestamps returns a list of processed detached timestamps.
// Port of getDetachedTimestamps().
func (d *DetachedTimestampSource) DetachedTimestamps() []*validation.TimestampToken {
	return d.detachedTimestamps
}

// AddExternalTimestamp adds the external timestamp to the source.
// Port of addExternalTimestamp(TimestampToken).
func (d *DetachedTimestampSource) AddExternalTimestamp(timestamp *validation.TimestampToken) error {
	if err := d.processExternalTimestamp(timestamp); err != nil {
		return err
	}
	d.detachedTimestamps = append(d.detachedTimestamps, timestamp)
	return nil
}

// processExternalTimestamp ports the private processExternalTimestamp(TimestampToken).
func (d *DetachedTimestampSource) processExternalTimestamp(externalTimestamp *validation.TimestampToken) error {
	d.populateSources(externalTimestamp)
	manifestReferences, err := d.getManifestReferences(externalTimestamp)
	if err != nil {
		return err
	}
	timestampAddReferences(externalTimestamp, manifestReferences)
	return nil
}

// getManifestReferences ports the private getManifestReferences(TimestampToken).
func (d *DetachedTimestampSource) getManifestReferences(externalTimestamp *validation.TimestampToken) (
	[]*validation.TimestampedReference, error) {
	result := []*validation.TimestampedReference{}
	manifestFile := externalTimestamp.ManifestFile()
	if manifestFile != nil {
		for _, timestampToken := range d.detachedTimestamps {
			if manifestFile.IsDocumentCovered(timestampToken.Filename()) {
				fromTimestamp, err := ReferencesFromTimestamp(timestampToken, d.certificateSource, d.crlSource, d.ocspSource)
				if err != nil {
					return nil, err
				}
				addReferences(&result, fromTimestamp)
			}
		}
	}
	return result, nil
}

// AddExternalEvidenceRecord adds the external evidence record to the source.
// Port of addExternalEvidenceRecord(EvidenceRecord).
func (d *DetachedTimestampSource) AddExternalEvidenceRecord(evidenceRecord validation.EvidenceRecord) error {
	if err := d.processExternalEvidenceRecord(evidenceRecord); err != nil {
		return err
	}
	d.detachedEvidenceRecords = append(d.detachedEvidenceRecords, evidenceRecord)
	return nil
}

// processExternalEvidenceRecord ports the private processExternalEvidenceRecord(EvidenceRecord).
func (d *DetachedTimestampSource) processExternalEvidenceRecord(evidenceRecord validation.EvidenceRecord) error {
	if err := d.addEncapsulatedReferencesFromTimestamps(evidenceRecord, d.DetachedTimestamps()); err != nil {
		return err
	}
	ProcessEvidenceRecordTimestamps(evidenceRecord)
	ProcessEmbeddedEvidenceRecords(evidenceRecord)
	for _, timestampToken := range evidenceRecord.Timestamps() {
		d.populateSources(timestampToken)
	}
	return nil
}

// populateSources ports the private populateSources(TimestampToken).
func (d *DetachedTimestampSource) populateSources(timestampToken *validation.TimestampToken) {
	d.certificateSource.Add(timestampToken.CertificateSource())
	d.crlSource.Add(timestampToken.CRLSource())
	d.ocspSource.Add(timestampToken.OCSPSource())
}

// addEncapsulatedReferencesFromTimestamps ports the private
// addEncapsulatedReferencesFromTimestamps(EvidenceRecord, List).
func (d *DetachedTimestampSource) addEncapsulatedReferencesFromTimestamps(evidenceRecord validation.EvidenceRecord,
	timestampTokens []*validation.TimestampToken) error {
	for _, timestampToken := range timestampTokens {
		if d.isCoveredTimestamp(evidenceRecord, timestampToken) {
			fromTimestamp, err := ReferencesFromTimestamp(timestampToken, d.certificateSource, d.crlSource, d.ocspSource)
			if err != nil {
				return err
			}
			evidenceRecord.SetTimestampedReferences(mergeReferences(evidenceRecord.TimestampedReferences(), fromTimestamp))
		}
	}
	return nil
}

// isCoveredTimestamp ports the private isCoveredTimestamp(EvidenceRecord, TimestampToken).
func (d *DetachedTimestampSource) isCoveredTimestamp(evidenceRecord validation.EvidenceRecord, timestampToken *validation.TimestampToken) bool {
	manifestFile := evidenceRecord.ManifestFile()
	if manifestFile != nil {
		for _, manifestEntry := range manifestFile.Entries() {
			if timestampToken.Filename() != "" && timestampToken.Filename() == manifestEntry.Uri() {
				return true
			}
		}
		return false
	}
	return true
}
