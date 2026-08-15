//go:build !eaa

// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/BasicBuildingBlocks.java (DSS 6.5.RC1) -
// !eaa stub counterpart of basic_building_blocks_eaa.go; see that file's
// header for STRUCTURAL FIX B's rationale.
//
// Without the "eaa" build tag, none of the four EAA-specific constructors
// (fc.NewEAAFormatChecking, fc.NewEAARevocationFormatChecking,
// sav.NewEAAAcceptanceValidation, sav.NewEAARevocationTokenAcceptanceValidation)
// are compiled into the binary at all, so this build can never validate an
// EAA or EAA_REVOCATION token to begin with - upstream Java dispatches by the
// document/token type actually being validated, and a build without the EAA
// module present simply never constructs an EAAWrapper/EAARevocationTokenWrapper
// token in the first place. These two methods therefore return nil
// unconditionally, mirroring "never construct those" rather than reimplementing
// any EAA logic: they exist only so the two Context_EAA(_REVOCATION) branches
// in basic_building_blocks.go have something to call under both tag states.
package blocks

import (
	"github.com/utain/esig/dss/detailedreport/jaxb"
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
