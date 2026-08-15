//go:build phase8e

// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/ValidationProcessUtils.java (DSS 6.5.RC1).
//
// This is the vpfswatsp-dependent half of ValidationProcessUtils: the two public
// methods that take a POEExtraction, plus the private helper one of them calls.
// eu.europa.esig.dss.validation.process.vpfswatsp (POEExtraction) is ported in
// phase 8e, so these are gated behind the "phase8e" build tag - the same
// tag-split the ASiC chunk uses for its phase8 forward dependencies - while the
// rest of the class compiles now (validation_process_utils.go).
//
// FORWARD DEPENDENCY: POEExtraction is assumed to be ported into this same Go
// package tree as process.POEExtraction with the shape the calls below make:
//
//	type POEExtraction interface {
//	    IsPOEExists(tokenId string, controlTime time.Time) bool
//	    IsPOEExistInRange(tokenId string, notBefore, notAfter *time.Time) bool
//	}
package process

import (
	"time"

	"github.com/utain/esig/dss/detailedreport/jaxb"
	"github.com/utain/esig/dss/diagnostic"
	"github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/utils"
)

// GetLatestAcceptableRevocationData returns a revocation data used for basic
// signature validation. Port of
// getLatestAcceptableRevocationData(TokenProxy, CertificateWrapper, Collection, Date, Map, POEExtraction).
func GetLatestAcceptableRevocationData(token diagnostic.TokenProxy, certificate *diagnostic.CertificateWrapper,
	revocationData []*diagnostic.CertificateRevocationWrapper, controlTime time.Time,
	bbbs map[string]*jaxb.XmlBasicBuildingBlocks, poe POEExtraction) *diagnostic.CertificateRevocationWrapper {
	var latestRevocationData *diagnostic.CertificateRevocationWrapper
	if poe.IsPOEExists(certificate.Id(), controlTime) {
		for _, revocationWrapper := range revocationData {
			revocationBBB := bbbs[revocationWrapper.Id()]
			if IsAllowedBasicRevocationDataValidation(revocationBBB.Conclusion) &&
				IsRevocationDataAcceptable(bbbs[token.Id()], certificate, &revocationWrapper.RevocationWrapper) &&
				revocationWrapper.ThisUpdate() != nil && revocationWrapper.ThisUpdate().Before(controlTime) &&
				poe.IsPOEExists(revocationWrapper.Id(), controlTime) &&
				(latestRevocationData == nil || (revocationWrapper.ProductionDate() != nil &&
					latestRevocationData.ProductionDate().Before(*revocationWrapper.ProductionDate()))) {
				latestRevocationData = revocationWrapper
			}
		}
	}
	return latestRevocationData
}

// GetAcceptableRevocationDataForPSVIfExistOrReturnAll verifies if there is an
// acceptable revocation data according to rules defined in 5.6.2.4 step 1) and
// returns a list of the revocation data. If none of the revocation data found,
// the method returns all the available revocation data. Port of
// getAcceptableRevocationDataForPSVIfExistOrReturnAll(TokenProxy, CertificateWrapper, Date, Map, POEExtraction, LevelRule).
func GetAcceptableRevocationDataForPSVIfExistOrReturnAll(token diagnostic.TokenProxy,
	certificate *diagnostic.CertificateWrapper, currentTime time.Time,
	bbbs map[string]*jaxb.XmlBasicBuildingBlocks, poe POEExtraction,
	revocationIssuerSunsetDateConstraint policy.LevelRule) []*diagnostic.CertificateRevocationWrapper {
	revocationWrappers := filterRevocationDataForPastSignatureValidation(
		token, certificate, currentTime, bbbs, poe, revocationIssuerSunsetDateConstraint)
	if utils.IsCollectionNotEmpty(revocationWrappers) {
		return revocationWrappers
	} else {
		return certificate.CertificateRevocationData()
	}
}

// filterRevocationDataForPastSignatureValidation filters revocation data for a
// signing certificate token according to rules defined in 5.6.2.4 step 1). Port
// of the private
// filterRevocationDataForPastSignatureValidation(TokenProxy, CertificateWrapper, Date, Map, POEExtraction, LevelRule).
func filterRevocationDataForPastSignatureValidation(token diagnostic.TokenProxy,
	certificate *diagnostic.CertificateWrapper, currentTime time.Time,
	bbbs map[string]*jaxb.XmlBasicBuildingBlocks, poe POEExtraction,
	revocationIssuerSunsetDateConstraint policy.LevelRule) []*diagnostic.CertificateRevocationWrapper {
	certificateRevocations := make([]*diagnostic.CertificateRevocationWrapper, 0)

	for _, certificateRevocation := range certificate.CertificateRevocationData() {
		revocationBBB := bbbs[certificateRevocation.Id()]
		revocationIssuer := certificateRevocation.SigningCertificate()

		if IsAllowedBasicRevocationDataValidation(revocationBBB.Conclusion) &&
			IsRevocationDataAcceptable(bbbs[token.Id()], certificate, &certificateRevocation.RevocationWrapper) &&
			revocationIssuer != nil && (IsTrustAnchor(revocationIssuer, currentTime, revocationIssuerSunsetDateConstraint) ||
			poe.IsPOEExistInRange(revocationIssuer.Id(), revocationIssuer.NotBefore(), revocationIssuer.NotAfter())) {
			certificateRevocations = append(certificateRevocations, certificateRevocation)
		}
	}
	return certificateRevocations
}
