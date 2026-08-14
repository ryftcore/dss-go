//go:build !phase8

// Non-phase8 counterpart of asic_with_cades_lta_validation_phase8.go - see that file's header
// for why this step is build-tag split. Renewing an existing archive timestamp's validation
// data requires the dss/validation engine (Phase 8), which is not linked into a non-phase8
// build, so this stub panics rather than silently producing an under-validated timestamp.
package cades

import (
	"github.com/utain/esig/dss/asic"
	"github.com/utain/esig/dss/model"
)

// extendLastArchiveTimestampWithValidationData is unavailable without the dss/validation engine
// (build with -tags phase8 once Phase 8 lands). See asic_with_cades_lta_validation_phase8.go.
func (e *ASiCWithCAdESLevelBaselineLTA) extendLastArchiveTimestampWithValidationData(asicContent *asic.ASiCContent, lastTimestamp model.DSSDocument) model.DSSDocument {
	panic("ASiC-CAdES LTA extension of an existing archive timestamp requires the dss/validation engine (Phase 8, not yet ported); rebuild with -tags phase8")
}
