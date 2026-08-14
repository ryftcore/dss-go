//go:build phase8

// Extracted from ASiCWithCAdESService.isLtaExtensionPossible(ASiCContent)
// (dss-asic-cades/src/main/java/eu/europa/esig/dss/asic/cades/signature/
// ASiCWithCAdESService.java, DSS 6.5.RC1) during Phase 7/8 integration - see
// asic_with_cades_lta_validation_phase8.go's header for why this method is build-tag split.
package cades

import (
	"github.com/utain/esig/dss/asic"
	"github.com/utain/esig/dss/utils"
)

// isLtaExtensionPossible reports whether an LTA extension is possible: it is not when a
// signature does not have a signature-time-stamp, as it would make the further signature
// extension impossible as per 162-1. Ports the private isLtaExtensionPossible(ASiCContent).
func (s *ASiCWithCAdESService) isLtaExtensionPossible(asicContent *asic.ASiCContent) bool {
	containerValidator := NewASiCContainerWithCAdESAnalyzerFromContent(asicContent)
	containerValidator.SetCertificateVerifier(s.CertificateVerifier)

	signatures := containerValidator.Signatures()
	for _, signature := range signatures {
		if utils.IsCollectionEmpty(signature.SignatureTimestamps()) {
			// Upstream logs a warning naming the signature id here.
			return false
		}
	}
	return true
}
