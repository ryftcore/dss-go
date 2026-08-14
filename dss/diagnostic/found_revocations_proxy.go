// Ported from dss-diagnostic-jaxb/src/main/java/eu/europa/esig/dss/diagnostic/FoundRevocationsProxy.java (DSS 6.5.RC1).
package diagnostic

import (
	"slices"

	"github.com/utain/esig/dss/diagnostic/jaxb"
	"github.com/utain/esig/dss/enumerations"
)

// FoundRevocationsProxy handles method logic to process XmlFoundRevocations and returns
// wrappers.
type FoundRevocationsProxy struct {
	// foundRevocations is the wrapped XmlFoundRevocations.
	foundRevocations *jaxb.XmlFoundRevocations
}

// NewFoundRevocationsProxy is the default constructor.
func NewFoundRevocationsProxy(foundRevocations *jaxb.XmlFoundRevocations) *FoundRevocationsProxy {
	return &FoundRevocationsProxy{foundRevocations: foundRevocations}
}

// getFoundRevocations lazily instantiates the wrapped object. Port of the private
// getFoundRevocations().
func (p *FoundRevocationsProxy) getFoundRevocations() *jaxb.XmlFoundRevocations {
	if p.foundRevocations == nil {
		p.foundRevocations = &jaxb.XmlFoundRevocations{}
	}
	return p.foundRevocations
}

// RelatedRevocationData returns a list of related revocation data. Port of
// getRelatedRevocationData().
func (p *FoundRevocationsProxy) RelatedRevocationData() []*RelatedRevocationWrapper {
	var revocationWrappers []*RelatedRevocationWrapper
	for _, relatedRevocation := range p.getFoundRevocations().RelatedRevocation {
		revocationWrappers = append(revocationWrappers, NewRelatedRevocationWrapper(relatedRevocation))
	}
	return revocationWrappers
}

// OrphanRevocationData returns a list of orphan revocations. Port of
// getOrphanRevocationData().
func (p *FoundRevocationsProxy) OrphanRevocationData() []*OrphanRevocationWrapper {
	var orphanTokens []*OrphanRevocationWrapper
	for _, orphanRevocation := range p.getFoundRevocations().OrphanRevocation {
		orphanTokens = append(orphanTokens, NewOrphanRevocationWrapper(orphanRevocation))
	}
	return orphanTokens
}

// RelatedRevocationsByOrigin returns a list of all RelatedRevocationWrapper used for the
// signature validation process with the given originType. Port of
// getRelatedRevocationsByOrigin(RevocationOrigin).
func (p *FoundRevocationsProxy) RelatedRevocationsByOrigin(originType enumerations.RevocationOrigin) []*RelatedRevocationWrapper {
	var revocationWrappers []*RelatedRevocationWrapper
	for _, relatedRevocation := range p.RelatedRevocationData() {
		if slices.Contains(relatedRevocation.Origins(), originType) {
			revocationWrappers = append(revocationWrappers, relatedRevocation)
		}
	}
	return revocationWrappers
}

// OrphanRevocationsByOrigin returns a list of all OrphanRevocationWrapper used for the
// signature validation process with the given originType. Port of
// getOrphanRevocationsByOrigin(RevocationOrigin).
func (p *FoundRevocationsProxy) OrphanRevocationsByOrigin(originType enumerations.RevocationOrigin) []*OrphanRevocationWrapper {
	var revocationWrappers []*OrphanRevocationWrapper
	for _, orphanRevocation := range p.OrphanRevocationData() {
		if slices.Contains(orphanRevocation.Origins(), originType) {
			revocationWrappers = append(revocationWrappers, orphanRevocation)
		}
	}
	return revocationWrappers
}

// RelatedRevocationsByRefOrigin returns a list of all RelatedRevocationWrapper used for the
// signature validation process with the given revocation reference origin. Port of
// getRelatedRevocationsByRefOrigin(RevocationRefOrigin).
func (p *FoundRevocationsProxy) RelatedRevocationsByRefOrigin(refOrigin enumerations.RevocationRefOrigin) []*RelatedRevocationWrapper {
	var revocationWrappers []*RelatedRevocationWrapper
	for _, relatedRevocation := range p.RelatedRevocationData() {
		for _, revocationRef := range relatedRevocation.References() {
			if slices.Contains(revocationRef.Origins(), refOrigin) {
				revocationWrappers = append(revocationWrappers, relatedRevocation)
				break
			}
		}
	}
	return revocationWrappers
}

// OrphanRevocationsByRefOrigin returns a list of all OrphanRevocationWrapper used for the
// signature validation process with the given reference origin. Port of
// getOrphanRevocationsByRefOrigin(RevocationRefOrigin).
func (p *FoundRevocationsProxy) OrphanRevocationsByRefOrigin(refOrigin enumerations.RevocationRefOrigin) []*OrphanRevocationWrapper {
	var revocationWrappers []*OrphanRevocationWrapper
	for _, orphanRevocation := range p.OrphanRevocationData() {
		for _, refWrapper := range orphanRevocation.References() {
			if slices.Contains(refWrapper.Origins(), refOrigin) {
				revocationWrappers = append(revocationWrappers, orphanRevocation)
				break
			}
		}
	}
	return revocationWrappers
}

// RelatedRevocationsByType returns a list of all RelatedRevocationWrapper used for the
// signature validation process with the given type. Port of
// getRelatedRevocationsByType(RevocationType).
func (p *FoundRevocationsProxy) RelatedRevocationsByType(revocationType enumerations.RevocationType) []*RelatedRevocationWrapper {
	var revocationWrappers []*RelatedRevocationWrapper
	for _, relatedRevocation := range p.RelatedRevocationData() {
		if revocationType == relatedRevocation.RevocationType() {
			revocationWrappers = append(revocationWrappers, relatedRevocation)
		}
	}
	return revocationWrappers
}

// OrphanRevocationsByType returns a list of all OrphanRevocationWrapper found in the
// signature, but not used during the validation process with the given type. Port of
// getOrphanRevocationsByType(RevocationType).
func (p *FoundRevocationsProxy) OrphanRevocationsByType(revocationType enumerations.RevocationType) []*OrphanRevocationWrapper {
	var revocationWrappers []*OrphanRevocationWrapper
	for _, orphanRevocation := range p.OrphanRevocationData() {
		if revocationType == orphanRevocation.RevocationType() {
			revocationWrappers = append(revocationWrappers, orphanRevocation)
		}
	}
	return revocationWrappers
}

// RelatedRevocationRefs returns a list of all found references for related revocations.
// Port of getRelatedRevocationRefs().
func (p *FoundRevocationsProxy) RelatedRevocationRefs() []*RevocationRefWrapper {
	var revocationRefs []*RevocationRefWrapper
	for _, revocationWrapper := range p.RelatedRevocationData() {
		revocationRefs = append(revocationRefs, revocationWrapper.References()...)
	}
	return revocationRefs
}

// OrphanRevocationRefs returns a list of all found references for orphan revocations. Port
// of getOrphanRevocationRefs().
func (p *FoundRevocationsProxy) OrphanRevocationRefs() []*RevocationRefWrapper {
	var revocationRefs []*RevocationRefWrapper
	for _, revocationWrapper := range p.OrphanRevocationData() {
		revocationRefs = append(revocationRefs, revocationWrapper.References()...)
	}
	return revocationRefs
}

// RelatedRevocationsByTypeAndOrigin returns a list of related revocation data by the given
// origin and type. Port of getRelatedRevocationsByTypeAndOrigin(RevocationType,
// RevocationOrigin).
func (p *FoundRevocationsProxy) RelatedRevocationsByTypeAndOrigin(revocationType enumerations.RevocationType, origin enumerations.RevocationOrigin) []*RelatedRevocationWrapper {
	var allRevocations []*RelatedRevocationWrapper
	for _, revocationWrapper := range p.RelatedRevocationsByOrigin(origin) {
		if revocationType == revocationWrapper.RevocationType() {
			allRevocations = append(allRevocations, revocationWrapper)
		}
	}
	return allRevocations
}

// OrphanRevocationsByTypeAndOrigin returns a list of orphan revocation data by the given
// origin and type. Port of getOrphanRevocationsByTypeAndOrigin(RevocationType,
// RevocationOrigin).
func (p *FoundRevocationsProxy) OrphanRevocationsByTypeAndOrigin(revocationType enumerations.RevocationType, origin enumerations.RevocationOrigin) []*OrphanRevocationWrapper {
	var allRevocations []*OrphanRevocationWrapper
	for _, revocationWrapper := range p.OrphanRevocationsByOrigin(origin) {
		if revocationType == revocationWrapper.RevocationType() {
			allRevocations = append(allRevocations, revocationWrapper)
		}
	}
	return allRevocations
}

// RelatedRevocationsByTypeAndRefOrigin returns a list of related revocation data by the
// given reference origin and type. Port of getRelatedRevocationsByTypeAndRefOrigin(
// RevocationType, RevocationRefOrigin).
func (p *FoundRevocationsProxy) RelatedRevocationsByTypeAndRefOrigin(revocationType enumerations.RevocationType, refOrigin enumerations.RevocationRefOrigin) []*RelatedRevocationWrapper {
	var allRevocations []*RelatedRevocationWrapper
	for _, revocationWrapper := range p.RelatedRevocationsByRefOrigin(refOrigin) {
		if revocationType == revocationWrapper.RevocationType() {
			allRevocations = append(allRevocations, revocationWrapper)
		}
	}
	return allRevocations
}

// OrphanRevocationsByTypeAndRefOrigin returns a list of orphan revocation data by the given
// reference origin and type. Port of getOrphanRevocationsByTypeAndRefOrigin(RevocationType,
// RevocationRefOrigin).
func (p *FoundRevocationsProxy) OrphanRevocationsByTypeAndRefOrigin(revocationType enumerations.RevocationType, refOrigin enumerations.RevocationRefOrigin) []*OrphanRevocationWrapper {
	var allRevocations []*OrphanRevocationWrapper
	for _, revocationWrapper := range p.OrphanRevocationsByRefOrigin(refOrigin) {
		if revocationType == revocationWrapper.RevocationType() {
			allRevocations = append(allRevocations, revocationWrapper)
		}
	}
	return allRevocations
}
