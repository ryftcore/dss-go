// Ported from dss-diagnostic-jaxb/src/main/java/eu/europa/esig/dss/diagnostic/FoundCertificatesProxy.java (DSS 6.5.RC1).
package diagnostic

import (
	"slices"

	"github.com/utain/esig/dss/diagnostic/jaxb"
	"github.com/utain/esig/dss/enumerations"
)

// FoundCertificatesProxy handles method logic to process XmlFoundCertificates and returns
// wrappers.
type FoundCertificatesProxy struct {
	// foundCertificates is the wrapped XmlFoundCertificates.
	foundCertificates *jaxb.XmlFoundCertificates
}

// NewFoundCertificatesProxy is the default constructor.
func NewFoundCertificatesProxy(foundCertificates *jaxb.XmlFoundCertificates) *FoundCertificatesProxy {
	return &FoundCertificatesProxy{foundCertificates: foundCertificates}
}

// getFoundCertificates lazily instantiates the wrapped object. Port of the private
// getFoundCertificates().
func (p *FoundCertificatesProxy) getFoundCertificates() *jaxb.XmlFoundCertificates {
	if p.foundCertificates == nil {
		p.foundCertificates = &jaxb.XmlFoundCertificates{}
	}
	return p.foundCertificates
}

// GetRelatedCertificates returns a list of related certificates. Port of
// getRelatedCertificates().
func (p *FoundCertificatesProxy) GetRelatedCertificates() []*RelatedCertificateWrapper {
	var certificateWrappers []*RelatedCertificateWrapper
	for _, relatedCertificate := range p.getFoundCertificates().RelatedCertificates {
		certificateWrappers = append(certificateWrappers, NewRelatedCertificateWrapper(relatedCertificate))
	}
	return certificateWrappers
}

// GetOrphanCertificates returns a list of orphan certificates. Port of getOrphanCertificates().
func (p *FoundCertificatesProxy) GetOrphanCertificates() []*OrphanCertificateWrapper {
	var orphanTokens []*OrphanCertificateWrapper
	for _, orphanCertificate := range p.getFoundCertificates().OrphanCertificates {
		orphanTokens = append(orphanTokens, NewOrphanCertificateWrapper(orphanCertificate))
	}
	return orphanTokens
}

// GetRelatedCertificatesByOrigin returns a list of found related RelatedCertificateWrapper
// with the given origin. Port of getRelatedCertificatesByOrigin(CertificateOrigin).
func (p *FoundCertificatesProxy) GetRelatedCertificatesByOrigin(origin enumerations.CertificateOrigin) []*RelatedCertificateWrapper {
	var certificateWrappers []*RelatedCertificateWrapper
	for _, relatedCertificate := range p.GetRelatedCertificates() {
		if slices.Contains(relatedCertificate.GetOrigins(), origin) {
			certificateWrappers = append(certificateWrappers, relatedCertificate)
		}
	}
	return certificateWrappers
}

// GetOrphanCertificatesByOrigin returns a list of found OrphanCertificateTokenWrapper with the
// given origin. Port of getOrphanCertificatesByOrigin(CertificateOrigin).
func (p *FoundCertificatesProxy) GetOrphanCertificatesByOrigin(origin enumerations.CertificateOrigin) []*OrphanCertificateWrapper {
	var orphanCertificatesWrappers []*OrphanCertificateWrapper
	for _, orphanCertificate := range p.GetOrphanCertificates() {
		if slices.Contains(orphanCertificate.GetOrigins(), origin) {
			orphanCertificatesWrappers = append(orphanCertificatesWrappers, orphanCertificate)
		}
	}
	return orphanCertificatesWrappers
}

// GetRelatedCertificatesByRefOrigin returns a list of found RelatedCertificateWrapper with the
// given reference origin. Port of getRelatedCertificatesByRefOrigin(CertificateRefOrigin).
func (p *FoundCertificatesProxy) GetRelatedCertificatesByRefOrigin(refOrigin enumerations.CertificateRefOrigin) []*RelatedCertificateWrapper {
	var certificateWrappers []*RelatedCertificateWrapper
	for _, relatedCertificate := range p.GetRelatedCertificates() {
		for _, certificateRef := range relatedCertificate.GetReferences() {
			if refOrigin == certificateRef.Origin() {
				certificateWrappers = append(certificateWrappers, relatedCertificate)
				break
			}
		}
	}
	return certificateWrappers
}

// GetOrphanCertificatesByRefOrigin returns a list of found OrphanCertificateTokenWrapper with
// the given reference origin. Port of getOrphanCertificatesByRefOrigin(CertificateRefOrigin).
func (p *FoundCertificatesProxy) GetOrphanCertificatesByRefOrigin(refOrigin enumerations.CertificateRefOrigin) []*OrphanCertificateWrapper {
	var orphanCertificatesWrappers []*OrphanCertificateWrapper
	for _, orphanCertificate := range p.GetOrphanCertificates() {
		for _, certificateRef := range orphanCertificate.GetReferences() {
			if refOrigin == certificateRef.Origin() {
				orphanCertificatesWrappers = append(orphanCertificatesWrappers, orphanCertificate)
				break
			}
		}
	}
	return orphanCertificatesWrappers
}

// GetRelatedCertificateRefs returns a list of all found references for related certificates.
// Port of getRelatedCertificateRefs().
func (p *FoundCertificatesProxy) GetRelatedCertificateRefs() []*CertificateRefWrapper {
	var certificateRefs []*CertificateRefWrapper
	for _, certificateWrapper := range p.GetRelatedCertificates() {
		certificateRefs = append(certificateRefs, certificateWrapper.GetReferences()...)
	}
	return certificateRefs
}

// GetOrphanCertificateRefs returns a list of all found references for orphan certificates.
// Port of getOrphanCertificateRefs().
func (p *FoundCertificatesProxy) GetOrphanCertificateRefs() []*CertificateRefWrapper {
	var certificateRefs []*CertificateRefWrapper
	for _, certificateWrapper := range p.GetOrphanCertificates() {
		certificateRefs = append(certificateRefs, certificateWrapper.GetReferences()...)
	}
	return certificateRefs
}

// GetRelatedCertificateRefsByRefOrigin returns a list of related certificate references by the
// given certificate reference origin. Port of
// getRelatedCertificateRefsByRefOrigin(CertificateRefOrigin).
func (p *FoundCertificatesProxy) GetRelatedCertificateRefsByRefOrigin(refOrigin enumerations.CertificateRefOrigin) []*CertificateRefWrapper {
	var certificateRefs []*CertificateRefWrapper
	for _, ref := range p.GetRelatedCertificateRefs() {
		if refOrigin == ref.Origin() {
			certificateRefs = append(certificateRefs, ref)
		}
	}
	return certificateRefs
}

// GetOrphanCertificateRefsByRefOrigin returns a list of orphan certificate references by the
// given certificate reference origin. Port of
// getOrphanCertificateRefsByRefOrigin(CertificateRefOrigin).
func (p *FoundCertificatesProxy) GetOrphanCertificateRefsByRefOrigin(refOrigin enumerations.CertificateRefOrigin) []*CertificateRefWrapper {
	var certificateRefs []*CertificateRefWrapper
	for _, ref := range p.GetOrphanCertificateRefs() {
		if refOrigin == ref.Origin() {
			certificateRefs = append(certificateRefs, ref)
		}
	}
	return certificateRefs
}
