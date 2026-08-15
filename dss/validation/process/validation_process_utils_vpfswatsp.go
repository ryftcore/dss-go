// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/ValidationProcessUtils.java (DSS 6.5.RC1).
//
// This is the vpfswatsp-dependent half of ValidationProcessUtils: the two public
// methods that take a POEExtraction, plus the private helper one of them calls.
// It was gated behind the "phase8e" build tag while
// eu.europa.esig.dss.validation.process.vpfswatsp was unported; phase 8e ported
// that package, so the gate is gone and the file is a plain member of the
// package again. It is kept as a separate file (rather than merged back into
// validation_process_utils.go) so that the one-Go-file-per-Java-class rule stays
// legible: both files carry the same "Ported from ValidationProcessUtils.java"
// header, and the split is documented in both.
//
// The POEExtraction parameter is the interface declared below rather than
// vpfswatsp.POEExtraction itself: vpfswatsp imports this package (Chain,
// ChainItem, the rest of ValidationProcessUtils), so this package cannot import
// vpfswatsp back. *vpfswatsp.POEExtraction satisfies it structurally.
package process

import (
	"time"

	"github.com/utain/esig/dss/detailedreport/jaxb"
	"github.com/utain/esig/dss/diagnostic"
	"github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/utils"
)

// POEExtraction is the set of proofs of existence the three functions below
// consult. It is the Go stand-in for the parameter type
// eu.europa.esig.dss.validation.process.vpfswatsp.POEExtraction, declared here
// as an interface because that package imports this one (see the file header);
// *vpfswatsp.POEExtraction is its only implementation.
type POEExtraction interface {
	// IsPOEExists returns true if a POE exists for the given token id at (or
	// before) the control time. Port of isPOEExists(String, Date).
	IsPOEExists(tokenId string, controlTime time.Time) bool
	// IsPOEExistInRange checks if a POE exists for the token with the given id
	// within the validity range between notBefore and notAfter inclusively.
	// Port of isPOEExistInRange(String, Date, Date).
	IsPOEExistInRange(tokenId string, notBefore, notAfter *time.Time) bool
}

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
