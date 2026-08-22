// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/evidencerecord/XAdESEvidenceRecordIncorporationParameters.java (DSS 6.5.RC1).
package xades

import (
	"github.com/ryftcore/dss-go/dss/document"
	"github.com/ryftcore/dss-go/dss/xades/definition"
	"github.com/ryftcore/dss-go/dss/xml/common"
)

// EvidenceRecordIncorporationParameters holds parameters for an evidence record
// incorporation within a XAdES signature.
type EvidenceRecordIncorporationParameters struct {
	document.AbstractEvidenceRecordIncorporationParameters

	// xadesERNamespace is the XAdES 132-3 namespace definition for the evidence record element
	// incorporation.
	xadesERNamespace *common.DSSNamespace
}

// NewXAdESEvidenceRecordIncorporationParameters is the default constructor.
func NewXAdESEvidenceRecordIncorporationParameters() *EvidenceRecordIncorporationParameters {
	return &EvidenceRecordIncorporationParameters{
		AbstractEvidenceRecordIncorporationParameters: document.NewAbstractEvidenceRecordIncorporationParameters(),
		xadesERNamespace: definition.XAdESNamespaceXAdESEvidencerecordNamespace,
	}
}

// XadesERNamespace gets a namespace for elements for the evidence record inclusion. Port of
// #getXadesERNamespace.
func (p *EvidenceRecordIncorporationParameters) XadesERNamespace() *common.DSSNamespace {
	return p.xadesERNamespace
}

// SetXadesERNamespace sets a namespace for elements for the evidence record inclusion.
// Default: xadesen:http://uri.etsi.org/19132/v1.1.1#
//
// Panics when xadesERNamespace is nil (Java's Objects.requireNonNull) or when its URI does not
// match the 132-3 definition (Java's IllegalArgumentException("The provided URI does not match
// the 132-3 definition!")). Port of #setXadesERNamespace.
func (p *EvidenceRecordIncorporationParameters) SetXadesERNamespace(xadesERNamespace *common.DSSNamespace) {
	if xadesERNamespace == nil {
		panic("xadesERNamespace cannot be null")
	}
	uri := xadesERNamespace.Uri()
	if definition.XAdESNamespaceXAdESEvidencerecordNamespace.IsSameUri(uri) {
		p.xadesERNamespace = xadesERNamespace
	} else {
		panic("The provided URI does not match the 132-3 definition!")
	}
}
