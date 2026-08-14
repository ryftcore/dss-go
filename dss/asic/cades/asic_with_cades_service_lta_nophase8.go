//go:build !phase8

// Non-phase8 counterpart of asic_with_cades_service_lta_phase8.go - see
// asic_with_cades_lta_validation_phase8.go's header for why this method is build-tag split.
// Determining LTA-extension possibility requires the dss/validation engine (Phase 8), which is
// not linked into a non-phase8 build, so this stub panics rather than silently mis-deciding.
package cades

import "github.com/utain/esig/dss/asic"

// isLtaExtensionPossible is unavailable without the dss/validation engine (build with -tags
// phase8 once Phase 8 lands). See asic_with_cades_service_lta_phase8.go.
func (s *ASiCWithCAdESService) isLtaExtensionPossible(asicContent *asic.ASiCContent) bool {
	panic("ASiC-CAdES timestamping of an already-signed/timestamped container requires the dss/validation engine (Phase 8, not yet ported); rebuild with -tags phase8")
}
