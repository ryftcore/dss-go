//go:build eaa

// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/BasicBuildingBlocks.java (DSS 6.5.RC1) -
// EAA-dependent slice only.
//
// fc.NewEAAFormatChecking,
// fc.NewEAARevocationFormatChecking, sav.NewEAAAcceptanceValidation and
// sav.NewEAARevocationTokenAcceptanceValidation are themselves gated behind
// //go:build eaa in their own packages (fc/eaa_format_checking.go,
// fc/eaa_revocation_format_checking.go, sav/eaa_acceptance_validation.go,
// sav/eaa_revocation_token_acceptance_validation.go - eu.europa.esig.dss.validation.process.eaa
// is deferred, so those files carry the same "eaa" feature tag).
// basic_building_blocks.go is un-tagged, so its two EAA-context
// branches that call into those four constructors are split into this
// eaa-tagged file - and its basic_building_blocks_noeaa.go !eaa counterpart -
// so the un-tagged, default `go build ./...` never references a symbol that
// only exists under -tags eaa.
package blocks

import (
	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/validation/process/bbb/fc"
	"github.com/ryftcore/dss-go/dss/validation/process/bbb/sav"
)

// executeEAAFormatChecking dispatches the Context_EAA / Context_EAA_REVOCATION
// branches of the private executeFormatChecking() that call
// fc.NewEAAFormatChecking / fc.NewEAARevocationFormatChecking.
func (b *BasicBuildingBlocks) executeEAAFormatChecking() *jaxb.XmlFC {
	if enumerations.Context_EAA == b.context {
		block := fc.NewEAAFormatChecking(b.i18nProvider, b.diagnosticData,
			b.token.(*diagnostic.EAAWrapper), b.context, b.policy)
		return block.Execute()
	} else if enumerations.Context_EAA_REVOCATION == b.context {
		block := fc.NewEAARevocationFormatChecking(b.i18nProvider, b.diagnosticData,
			b.token.(*diagnostic.EAARevocationTokenWrapper), b.context, b.policy)
		return block.Execute()
	}
	return nil
}

// eaaAcceptanceValidationBlock dispatches the Context_EAA / Context_EAA_REVOCATION
// branches of the private executeSignatureAcceptanceValidation(XmlAOV) that
// call sav.NewEAAAcceptanceValidation / sav.NewEAARevocationTokenAcceptanceValidation.
func (b *BasicBuildingBlocks) eaaAcceptanceValidationBlock(xmlAOV *jaxb.XmlAOV) savBlock {
	if enumerations.Context_EAA == b.context {
		return sav.NewEAAAcceptanceValidation(b.i18nProvider, b.currentTime,
			b.token.(*diagnostic.EAAWrapper), b.bbbs, xmlAOV, b.policy)
	} else if enumerations.Context_EAA_REVOCATION == b.context {
		return sav.NewEAARevocationTokenAcceptanceValidation(b.i18nProvider, b.currentTime,
			b.token.(*diagnostic.EAARevocationTokenWrapper), xmlAOV, b.policy)
	}
	return nil
}
