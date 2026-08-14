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

// RelatedCertificates returns a list of related certificates. Port of
// getRelatedCertificates().
func (p *FoundCertificatesProxy) RelatedCertificates() []*RelatedCertificateWrapper {
	var certificateWrappers []*RelatedCertificateWrapper
	for _, relatedCertificate := range p.getFoundCertificates().RelatedCertificates {
		certificateWrappers = append(certificateWrappers, NewRelatedCertificateWrapper(relatedCertificate))
	}
	return certificateWrappers
}

// OrphanCertificates returns a list of orphan certificates. Port of getOrphanCertificates().
func (p *FoundCertificatesProxy) OrphanCertificates() []*OrphanCertificateWrapper {
	var orphanTokens []*OrphanCertificateWrapper
	for _, orphanCertificate := range p.getFoundCertificates().OrphanCertificates {
		orphanTokens = append(orphanTokens, NewOrphanCertificateWrapper(orphanCertificate))
	}
	return orphanTokens
}

// RelatedCertificatesByOrigin returns a list of found related RelatedCertificateWrapper
// with the given origin. Port of getRelatedCertificatesByOrigin(CertificateOrigin).
func (p *FoundCertificatesProxy) RelatedCertificatesByOrigin(origin enumerations.CertificateOrigin) []*RelatedCertificateWrapper {
	var certificateWrappers []*RelatedCertificateWrapper
	for _, relatedCertificate := range p.RelatedCertificates() {
		if slices.Contains(relatedCertificate.Origins(), origin) {
			certificateWrappers = append(certificateWrappers, relatedCertificate)
		}
	}
	return certificateWrappers
}

// OrphanCertificatesByOrigin returns a list of found OrphanCertificateTokenWrapper with the
// given origin. Port of getOrphanCertificatesByOrigin(CertificateOrigin).
func (p *FoundCertificatesProxy) OrphanCertificatesByOrigin(origin enumerations.CertificateOrigin) []*OrphanCertificateWrapper {
	var orphanCertificatesWrappers []*OrphanCertificateWrapper
	for _, orphanCertificate := range p.OrphanCertificates() {
		if slices.Contains(orphanCertificate.Origins(), origin) {
			orphanCertificatesWrappers = append(orphanCertificatesWrappers, orphanCertificate)
		}
	}
	return orphanCertificatesWrappers
}

// RelatedCertificatesByRefOrigin returns a list of found RelatedCertificateWrapper with the
// given reference origin. Port of getRelatedCertificatesByRefOrigin(CertificateRefOrigin).
func (p *FoundCertificatesProxy) RelatedCertificatesByRefOrigin(refOrigin enumerations.CertificateRefOrigin) []*RelatedCertificateWrapper {
	var certificateWrappers []*RelatedCertificateWrapper
	for _, relatedCertificate := range p.RelatedCertificates() {
		for _, certificateRef := range relatedCertificate.References() {
			if refOrigin == certificateRef.Origin() {
				certificateWrappers = append(certificateWrappers, relatedCertificate)
				break
			}
		}
	}
	return certificateWrappers
}

// OrphanCertificatesByRefOrigin returns a list of found OrphanCertificateTokenWrapper with
// the given reference origin. Port of getOrphanCertificatesByRefOrigin(CertificateRefOrigin).
func (p *FoundCertificatesProxy) OrphanCertificatesByRefOrigin(refOrigin enumerations.CertificateRefOrigin) []*OrphanCertificateWrapper {
	var orphanCertificatesWrappers []*OrphanCertificateWrapper
	for _, orphanCertificate := range p.OrphanCertificates() {
		for _, certificateRef := range orphanCertificate.References() {
			if refOrigin == certificateRef.Origin() {
				orphanCertificatesWrappers = append(orphanCertificatesWrappers, orphanCertificate)
				break
			}
		}
	}
	return orphanCertificatesWrappers
}

// RelatedCertificateRefs returns a list of all found references for related certificates.
// Port of getRelatedCertificateRefs().
func (p *FoundCertificatesProxy) RelatedCertificateRefs() []*CertificateRefWrapper {
	var certificateRefs []*CertificateRefWrapper
	for _, certificateWrapper := range p.RelatedCertificates() {
		certificateRefs = append(certificateRefs, certificateWrapper.References()...)
	}
	return certificateRefs
}

// OrphanCertificateRefs returns a list of all found references for orphan certificates.
// Port of getOrphanCertificateRefs().
func (p *FoundCertificatesProxy) OrphanCertificateRefs() []*CertificateRefWrapper {
	var certificateRefs []*CertificateRefWrapper
	for _, certificateWrapper := range p.OrphanCertificates() {
		certificateRefs = append(certificateRefs, certificateWrapper.References()...)
	}
	return certificateRefs
}

// RelatedCertificateRefsByRefOrigin returns a list of related certificate references by the
// given certificate reference origin. Port of
// getRelatedCertificateRefsByRefOrigin(CertificateRefOrigin).
func (p *FoundCertificatesProxy) RelatedCertificateRefsByRefOrigin(refOrigin enumerations.CertificateRefOrigin) []*CertificateRefWrapper {
	var certificateRefs []*CertificateRefWrapper
	for _, ref := range p.RelatedCertificateRefs() {
		if refOrigin == ref.Origin() {
			certificateRefs = append(certificateRefs, ref)
		}
	}
	return certificateRefs
}

// OrphanCertificateRefsByRefOrigin returns a list of orphan certificate references by the
// given certificate reference origin. Port of
// getOrphanCertificateRefsByRefOrigin(CertificateRefOrigin).
func (p *FoundCertificatesProxy) OrphanCertificateRefsByRefOrigin(refOrigin enumerations.CertificateRefOrigin) []*CertificateRefWrapper {
	var certificateRefs []*CertificateRefWrapper
	for _, ref := range p.OrphanCertificateRefs() {
		if refOrigin == ref.Origin() {
			certificateRefs = append(certificateRefs, ref)
		}
	}
	return certificateRefs
}
