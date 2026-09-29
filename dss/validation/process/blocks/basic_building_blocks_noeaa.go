//go:build !eaa

// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/BasicBuildingBlocks.java (DSS 6.5.RC1) -
// !eaa stub counterpart of basic_building_blocks_eaa.go.
//
// Without the "eaa" build tag, none of the four EAA-specific constructors
// (fc.NewEAAFormatChecking, fc.NewEAARevocationFormatChecking,
// sav.NewEAAAcceptanceValidation, sav.NewEAARevocationTokenAcceptanceValidation)
// are compiled into the binary at all, so this build cannot run the EAA-specific
// format-checking and acceptance-validation blocks of an EAA or EAA_REVOCATION token.
// These two methods therefore return nil unconditionally, rather than
// reimplementing any EAA logic: they exist only so the two ContextEAA /
// ContextEAARevocation branches in basic_building_blocks.go have something to call
// under both tag states.
//
// NOTE: this does NOT mean an EAAWrapper is never constructed in a build without the
// tag - diagnostic data parsing (dss/diagnostic) and the EAA presentation executor and
// report builder (dss/validation/executor) are untagged, so diagnostic data holding an
// EAA is accepted here, and the basic building blocks of its EAA token then simply carry
// no FC block. The signature/certificate executors report on such data without the EAA
// blocks (see docs/compatibility/known-gaps.md, "EAA"); the EAA presentation process
// (dss/validation/process/eaa) needs the FC block and fails with an explicit "requires
// the 'eaa' build tag" error instead of proceeding without the format-checking gate.
package blocks

import (
	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
)

// executeEAAFormatChecking is the !eaa stub of the eaa-tagged method of the
// same name: no EAA format-checking module is present in this build.
func (b *BasicBuildingBlocks) executeEAAFormatChecking() *jaxb.XmlFC {
	return nil
}

// eaaAcceptanceValidationBlock is the !eaa stub of the eaa-tagged method of
// the same name: no EAA acceptance-validation module is present in this
// build.
func (b *BasicBuildingBlocks) eaaAcceptanceValidationBlock(xmlAOV *jaxb.XmlAOV) savBlock {
	return nil
}
