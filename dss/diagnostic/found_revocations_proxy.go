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

// GetRelatedRevocationData returns a list of related revocation data. Port of
// getRelatedRevocationData().
func (p *FoundRevocationsProxy) GetRelatedRevocationData() []*RelatedRevocationWrapper {
	var revocationWrappers []*RelatedRevocationWrapper
	for _, relatedRevocation := range p.getFoundRevocations().RelatedRevocations {
		revocationWrappers = append(revocationWrappers, NewRelatedRevocationWrapper(relatedRevocation))
	}
	return revocationWrappers
}

// GetOrphanRevocationData returns a list of orphan revocations. Port of
// getOrphanRevocationData().
func (p *FoundRevocationsProxy) GetOrphanRevocationData() []*OrphanRevocationWrapper {
	var orphanTokens []*OrphanRevocationWrapper
	for _, orphanRevocation := range p.getFoundRevocations().OrphanRevocations {
		orphanTokens = append(orphanTokens, NewOrphanRevocationWrapper(orphanRevocation))
	}
	return orphanTokens
}

// GetRelatedRevocationsByOrigin returns a list of all RelatedRevocationWrapper used for the
// signature validation process with the given originType. Port of
// getRelatedRevocationsByOrigin(RevocationOrigin).
func (p *FoundRevocationsProxy) GetRelatedRevocationsByOrigin(originType enumerations.RevocationOrigin) []*RelatedRevocationWrapper {
	var revocationWrappers []*RelatedRevocationWrapper
	for _, relatedRevocation := range p.GetRelatedRevocationData() {
		if slices.Contains(relatedRevocation.GetOrigins(), originType) {
			revocationWrappers = append(revocationWrappers, relatedRevocation)
		}
	}
	return revocationWrappers
}

// GetOrphanRevocationsByOrigin returns a list of all OrphanRevocationWrapper used for the
// signature validation process with the given originType. Port of
// getOrphanRevocationsByOrigin(RevocationOrigin).
func (p *FoundRevocationsProxy) GetOrphanRevocationsByOrigin(originType enumerations.RevocationOrigin) []*OrphanRevocationWrapper {
	var revocationWrappers []*OrphanRevocationWrapper
	for _, orphanRevocation := range p.GetOrphanRevocationData() {
		if slices.Contains(orphanRevocation.GetOrigins(), originType) {
			revocationWrappers = append(revocationWrappers, orphanRevocation)
		}
	}
	return revocationWrappers
}

// GetRelatedRevocationsByRefOrigin returns a list of all RelatedRevocationWrapper used for the
// signature validation process with the given revocation reference origin. Port of
// getRelatedRevocationsByRefOrigin(RevocationRefOrigin).
func (p *FoundRevocationsProxy) GetRelatedRevocationsByRefOrigin(refOrigin enumerations.RevocationRefOrigin) []*RelatedRevocationWrapper {
	var revocationWrappers []*RelatedRevocationWrapper
	for _, relatedRevocation := range p.GetRelatedRevocationData() {
		for _, revocationRef := range relatedRevocation.GetReferences() {
			if slices.Contains(revocationRef.GetOrigins(), refOrigin) {
				revocationWrappers = append(revocationWrappers, relatedRevocation)
				break
			}
		}
	}
	return revocationWrappers
}

// GetOrphanRevocationsByRefOrigin returns a list of all OrphanRevocationWrapper used for the
// signature validation process with the given reference origin. Port of
// getOrphanRevocationsByRefOrigin(RevocationRefOrigin).
func (p *FoundRevocationsProxy) GetOrphanRevocationsByRefOrigin(refOrigin enumerations.RevocationRefOrigin) []*OrphanRevocationWrapper {
	var revocationWrappers []*OrphanRevocationWrapper
	for _, orphanRevocation := range p.GetOrphanRevocationData() {
		for _, refWrapper := range orphanRevocation.GetReferences() {
			if slices.Contains(refWrapper.GetOrigins(), refOrigin) {
				revocationWrappers = append(revocationWrappers, orphanRevocation)
				break
			}
		}
	}
	return revocationWrappers
}

// GetRelatedRevocationsByType returns a list of all RelatedRevocationWrapper used for the
// signature validation process with the given type. Port of
// getRelatedRevocationsByType(RevocationType).
func (p *FoundRevocationsProxy) GetRelatedRevocationsByType(revocationType enumerations.RevocationType) []*RelatedRevocationWrapper {
	var revocationWrappers []*RelatedRevocationWrapper
	for _, relatedRevocation := range p.GetRelatedRevocationData() {
		if revocationType == relatedRevocation.GetRevocationType() {
			revocationWrappers = append(revocationWrappers, relatedRevocation)
		}
	}
	return revocationWrappers
}

// GetOrphanRevocationsByType returns a list of all OrphanRevocationWrapper found in the
// signature, but not used during the validation process with the given type. Port of
// getOrphanRevocationsByType(RevocationType).
func (p *FoundRevocationsProxy) GetOrphanRevocationsByType(revocationType enumerations.RevocationType) []*OrphanRevocationWrapper {
	var revocationWrappers []*OrphanRevocationWrapper
	for _, orphanRevocation := range p.GetOrphanRevocationData() {
		if revocationType == orphanRevocation.GetRevocationType() {
			revocationWrappers = append(revocationWrappers, orphanRevocation)
		}
	}
	return revocationWrappers
}

// GetRelatedRevocationRefs returns a list of all found references for related revocations.
// Port of getRelatedRevocationRefs().
func (p *FoundRevocationsProxy) GetRelatedRevocationRefs() []*RevocationRefWrapper {
	var revocationRefs []*RevocationRefWrapper
	for _, revocationWrapper := range p.GetRelatedRevocationData() {
		revocationRefs = append(revocationRefs, revocationWrapper.GetReferences()...)
	}
	return revocationRefs
}

// GetOrphanRevocationRefs returns a list of all found references for orphan revocations. Port
// of getOrphanRevocationRefs().
func (p *FoundRevocationsProxy) GetOrphanRevocationRefs() []*RevocationRefWrapper {
	var revocationRefs []*RevocationRefWrapper
	for _, revocationWrapper := range p.GetOrphanRevocationData() {
		revocationRefs = append(revocationRefs, revocationWrapper.GetReferences()...)
	}
	return revocationRefs
}

// GetRelatedRevocationsByTypeAndOrigin returns a list of related revocation data by the given
// origin and type. Port of getRelatedRevocationsByTypeAndOrigin(RevocationType,
// RevocationOrigin).
func (p *FoundRevocationsProxy) GetRelatedRevocationsByTypeAndOrigin(revocationType enumerations.RevocationType, origin enumerations.RevocationOrigin) []*RelatedRevocationWrapper {
	var allRevocations []*RelatedRevocationWrapper
	for _, revocationWrapper := range p.GetRelatedRevocationsByOrigin(origin) {
		if revocationType == revocationWrapper.GetRevocationType() {
			allRevocations = append(allRevocations, revocationWrapper)
		}
	}
	return allRevocations
}

// GetOrphanRevocationsByTypeAndOrigin returns a list of orphan revocation data by the given
// origin and type. Port of getOrphanRevocationsByTypeAndOrigin(RevocationType,
// RevocationOrigin).
func (p *FoundRevocationsProxy) GetOrphanRevocationsByTypeAndOrigin(revocationType enumerations.RevocationType, origin enumerations.RevocationOrigin) []*OrphanRevocationWrapper {
	var allRevocations []*OrphanRevocationWrapper
	for _, revocationWrapper := range p.GetOrphanRevocationsByOrigin(origin) {
		if revocationType == revocationWrapper.GetRevocationType() {
			allRevocations = append(allRevocations, revocationWrapper)
		}
	}
	return allRevocations
}

// GetRelatedRevocationsByTypeAndRefOrigin returns a list of related revocation data by the
// given reference origin and type. Port of getRelatedRevocationsByTypeAndRefOrigin(
// RevocationType, RevocationRefOrigin).
func (p *FoundRevocationsProxy) GetRelatedRevocationsByTypeAndRefOrigin(revocationType enumerations.RevocationType, refOrigin enumerations.RevocationRefOrigin) []*RelatedRevocationWrapper {
	var allRevocations []*RelatedRevocationWrapper
	for _, revocationWrapper := range p.GetRelatedRevocationsByRefOrigin(refOrigin) {
		if revocationType == revocationWrapper.GetRevocationType() {
			allRevocations = append(allRevocations, revocationWrapper)
		}
	}
	return allRevocations
}

// GetOrphanRevocationsByTypeAndRefOrigin returns a list of orphan revocation data by the given
// reference origin and type. Port of getOrphanRevocationsByTypeAndRefOrigin(RevocationType,
// RevocationRefOrigin).
func (p *FoundRevocationsProxy) GetOrphanRevocationsByTypeAndRefOrigin(revocationType enumerations.RevocationType, refOrigin enumerations.RevocationRefOrigin) []*OrphanRevocationWrapper {
	var allRevocations []*OrphanRevocationWrapper
	for _, revocationWrapper := range p.GetOrphanRevocationsByRefOrigin(refOrigin) {
		if revocationType == revocationWrapper.GetRevocationType() {
			allRevocations = append(allRevocations, revocationWrapper)
		}
	}
	return allRevocations
}
