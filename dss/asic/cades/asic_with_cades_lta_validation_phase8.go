//go:build phase8

// Extracted from ASiCWithCAdESLevelBaselineLTA.extend(ASiCContent, DigestAlgorithm,
// DigestAlgorithm) (dss-asic-cades/src/main/java/eu/europa/esig/dss/asic/cades/signature/
// ASiCWithCAdESLevelBaselineLTA.java, DSS 6.5.RC1) during Phase 7/8 integration.
//
// INTEGRATOR NOTE: re-validating the container's existing signatures and detached timestamps
// to gather fresh validation data for an archive-timestamp renewal requires the dss/validation
// engine (Phase 8, not yet ported at the time this Phase 7 chunk landed). This file carries
// that one analyzer-dependent step, gated the same way as the other analyzer/validator files in
// this tree (see asic_container_with_cades_analyzer.go); asic_with_cades_lta_validation_nophase8.go
// carries the same method's signature for the default (non-phase8) build. Once Phase 8 lands,
// this split can be folded back into asic_with_cades_level_baseline_lta.go.
package cades

import (
	"github.com/utain/esig/dss/asic"
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/spi/validation/executor"
)

// extendLastArchiveTimestampWithValidationData re-validates the container's signatures and
// detached timestamps, gathers the resulting validation data, and returns a new version of
// lastTimestamp extended (via CMS) with that data. Ports the analyzer-dependent portion of the
// private extend(ASiCContent, DigestAlgorithm, DigestAlgorithm).
func (e *ASiCWithCAdESLevelBaselineLTA) extendLastArchiveTimestampWithValidationData(asicContent *asic.ASiCContent, lastTimestamp model.DSSDocument) model.DSSDocument {
	containerValidator := NewASiCContainerWithCAdESAnalyzerFromContent(asicContent)
	containerValidator.SetCertificateVerifier(e.CertificateVerifier)
	containerValidator.SetValidationContextExecutor(executor.CompleteValidationContextExecutorInstance)

	allSignatures := containerValidator.GetAllSignatures()
	detachedTimestamps := containerValidator.DetachedTimestamps()

	validationDataContainer, err := containerValidator.GetValidationDataWithTimestamps(allSignatures, detachedTimestamps)
	if err != nil {
		panic(err)
	}
	allValidationData := validationDataContainer.AllValidationData()

	// ensure the validation data is not duplicated
	for _, signature := range allSignatures {
		allValidationData.ExcludeCertificateTokens(signature.CompleteCertificateSource().Certificates())
		allValidationData.ExcludeCRLTokens(asicCAdESLTARevocationIdentifiers(signature.CompleteCRLSource().AllRevocationBinaries()))
		allValidationData.ExcludeOCSPTokens(asicCAdESLTARevocationIdentifiers(signature.CompleteOCSPSource().AllRevocationBinaries()))
	}
	for _, timestampToken := range detachedTimestamps {
		allValidationData.ExcludeCertificateTokens(timestampToken.CertificateSource().Certificates())
		allValidationData.ExcludeCRLTokens(asicCAdESLTARevocationIdentifiers(timestampToken.CRLSource().AllRevocationBinaries()))
		allValidationData.ExcludeOCSPTokens(asicCAdESLTARevocationIdentifiers(timestampToken.OCSPSource().AllRevocationBinaries()))
	}

	// a newer version of the timestamp must be created
	return e.extendTimestamp(lastTimestamp, allValidationData)
}
