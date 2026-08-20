// Ported from dss-tsl-validation/src/main/java/eu/europa/esig/dss/tsl/parsing/ParsingUtils.java (DSS 6.5.RC1).
package tsl

import (
	tslmodel "github.com/utain/esig/dss/model/tsl"
	"github.com/utain/esig/dss/spi"
	"github.com/utain/esig/dss/utils"
)

// ParsingUtilsXMLLOTLPointer extracts the XML LOTL Pointer from the parsing cache of a pivot,
// answering nil when the pivot does not carry exactly one. Port of the static
// getXMLLOTLPointer(TLParsingCacheDTO).
//
// The two slf4j warnings upstream logs ("The provided parsing cache DTO is null or does not
// exist!" and "Unable to find the XML LOTL Pointer in the pivot (nb occurrences : {}). Must be one
// occurence!") are dropped, per PORTING.md; the nbLOTLPointersInPivot counter existed only to feed
// the second one.
func ParsingUtilsXMLLOTLPointer(parsingCacheDTO *TLParsingCacheDTO) *tslmodel.OtherTSLPointer {
	if parsingCacheDTO != nil && parsingCacheDTO.IsResultExist() {
		lotlOtherPointers := parsingCacheDTO.LotlOtherPointers()
		if utils.CollectionSize(lotlOtherPointers) == 1 {
			return lotlOtherPointers[0]
		}
	}
	return nil
}

// ParsingUtilsLOTLAnnouncedCertificateSource extracts the SDIs present in an OtherTSLPointer into
// a CertificateSource. Port of the static getLOTLAnnouncedCertificateSource(OtherTSLPointer).
func ParsingUtilsLOTLAnnouncedCertificateSource(currentLOTLPointer *tslmodel.OtherTSLPointer) spi.CertificateSource {
	certificateSource := spi.NewCommonCertificateSource()
	for _, certificate := range currentLOTLPointer.SdiCertificates() {
		certificateSource.AddCertificate(certificate)
	}
	return &certificateSource
}
