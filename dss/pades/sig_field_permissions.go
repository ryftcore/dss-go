// Ported from dss-pades/src/main/java/eu/europa/esig/dss/pdf/SigFieldPermissions.java (DSS 6.5.RC1).
//
// eu.europa.esig.dss.pdf is implemented by this file and others (see pdf_object.go's header).
// Its shape (SetAction/SetFields/SetCertificationPermission plus their getter counterparts) is
// used by pades_utils.go, pdf_signature_field.go, pdf_signature_dictionary.go and
// pades_diagnostic_data_builder.go.
package pades

import "github.com/ryftcore/dss-go/dss/enumerations"

// SigFieldPermissions defines a list of restrictions imposed to a PDF document's modifications
// by the current signature/field. Port of the SigFieldPermissions class.
type SigFieldPermissions struct {
	// action indicates the set of fields that should be locked.
	action enumerations.PdfLockAction

	// fields contains a set of fields.
	fields []string

	// certificationPermission is the access permissions (optional).
	certificationPermission enumerations.CertificationPermission
}

// NewSigFieldPermissions instantiates an object with null-equivalent (zero) values.
// Port of the default constructor.
func NewSigFieldPermissions() *SigFieldPermissions {
	return &SigFieldPermissions{}
}

// Action gets the defined action. Port of #getAction.
func (p *SigFieldPermissions) Action() enumerations.PdfLockAction {
	return p.action
}

// SetAction sets the action. Port of #setAction.
func (p *SigFieldPermissions) SetAction(action enumerations.PdfLockAction) {
	p.action = action
}

// Fields gets a list of field names. Port of #getFields.
func (p *SigFieldPermissions) Fields() []string {
	return p.fields
}

// SetFields sets a list of field names. Port of #setFields.
func (p *SigFieldPermissions) SetFields(fields []string) {
	p.fields = fields
}

// CertificationPermission gets the CertificationPermission. Port of #getCertificationPermission.
func (p *SigFieldPermissions) CertificationPermission() enumerations.CertificationPermission {
	return p.certificationPermission
}

// SetCertificationPermission sets the CertificationPermission. Port of #setCertificationPermission.
func (p *SigFieldPermissions) SetCertificationPermission(certificationPermission enumerations.CertificationPermission) {
	p.certificationPermission = certificationPermission
}
