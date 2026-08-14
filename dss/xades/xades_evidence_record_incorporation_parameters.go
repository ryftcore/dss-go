// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/evidencerecord/XAdESEvidenceRecordIncorporationParameters.java (DSS 6.5.RC1).
package xades

import (
	"github.com/utain/esig/dss/document"
	"github.com/utain/esig/dss/xades/definition"
	"github.com/utain/esig/dss/xml/common"
)

// XAdESEvidenceRecordIncorporationParameters holds parameters for an evidence record
// incorporation within a XAdES signature.
type XAdESEvidenceRecordIncorporationParameters struct {
	document.AbstractEvidenceRecordIncorporationParameters

	// xadesERNamespace is the XAdES 132-3 namespace definition for the evidence record element
	// incorporation.
	xadesERNamespace *common.DSSNamespace
}

// NewXAdESEvidenceRecordIncorporationParameters is the default constructor.
func NewXAdESEvidenceRecordIncorporationParameters() *XAdESEvidenceRecordIncorporationParameters {
	return &XAdESEvidenceRecordIncorporationParameters{
		AbstractEvidenceRecordIncorporationParameters: document.NewAbstractEvidenceRecordIncorporationParameters(),
		xadesERNamespace: definition.XAdESNamespace_XADES_EVIDENCERECORD_NAMESPACE,
	}
}

// XadesERNamespace gets a namespace for elements for the evidence record inclusion. Port of
// #getXadesERNamespace.
func (p *XAdESEvidenceRecordIncorporationParameters) XadesERNamespace() *common.DSSNamespace {
	return p.xadesERNamespace
}

// SetXadesERNamespace sets a namespace for elements for the evidence record inclusion.
// Default: xadesen:http://uri.etsi.org/19132/v1.1.1#
//
// Panics when xadesERNamespace is nil (Java's Objects.requireNonNull) or when its URI does not
// match the 132-3 definition (Java's IllegalArgumentException("The provided URI does not match
// the 132-3 definition!")). Port of #setXadesERNamespace.
func (p *XAdESEvidenceRecordIncorporationParameters) SetXadesERNamespace(xadesERNamespace *common.DSSNamespace) {
	if xadesERNamespace == nil {
		panic("xadesERNamespace cannot be null")
	}
	uri := xadesERNamespace.Uri()
	if definition.XAdESNamespace_XADES_EVIDENCERECORD_NAMESPACE.IsSameUri(uri) {
		p.xadesERNamespace = xadesERNamespace
	} else {
		panic("The provided URI does not match the 132-3 definition!")
	}
}
