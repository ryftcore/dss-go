package fc

import (
	"testing"

	"github.com/ryftcore/dss-go/dss/diagnostic"
	diagjaxb "github.com/ryftcore/dss-go/dss/diagnostic/jaxb"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

func newPDFRevision(permissions enumerations.CertificationPermission, formFill, annotation, undefined int) *diagnostic.PDFRevisionWrapper {
	objMods := &diagjaxb.XmlObjectModifications{
		SignatureOrFormFill: make([]*diagjaxb.XmlObjectModification, formFill),
		AnnotationChange:    make([]*diagjaxb.XmlObjectModification, annotation),
		Undefined:           make([]*diagjaxb.XmlObjectModification, undefined),
	}
	var docMDP *diagjaxb.XmlDocMDP
	if permissions != "" {
		p := diagjaxb.CertificationPermissionValue(permissions)
		docMDP = &diagjaxb.XmlDocMDP{Permissions: &p}
	}
	xmlRevision := &diagjaxb.XmlPDFRevision{
		PDFSignatureDictionary: &diagjaxb.XmlPDFSignatureDictionary{DocMDP: docMDP},
		ModificationDetection:  &diagjaxb.XmlModificationDetection{ObjectModifications: objMods},
	}
	return diagnostic.NewPDFRevisionWrapper(xmlRevision)
}

func TestDocMDPCheck_Process(t *testing.T) {
	tests := []struct {
		name        string
		permissions enumerations.CertificationPermission
		formFill    int
		annotation  int
		undefined   int
		want        bool
	}{
		{"happy: no permissions set, no modifications flag", "", 0, 0, 0, true},
		{"happy: NO_CHANGE_PERMITTED, no changes", enumerations.CertificationPermissionNoChangePermitted, 0, 0, 0, true},
		{"failure: NO_CHANGE_PERMITTED, form fill change", enumerations.CertificationPermissionNoChangePermitted, 1, 0, 0, false},
		{"happy: MINIMAL_CHANGES_PERMITTED, form fill only", enumerations.CertificationPermissionMinimalChangesPermitted, 1, 0, 0, true},
		{"failure: MINIMAL_CHANGES_PERMITTED, annotation change", enumerations.CertificationPermissionMinimalChangesPermitted, 0, 1, 0, false},
		{"happy: CHANGES_PERMITTED, form+annotation", enumerations.CertificationPermissionChangesPermitted, 1, 1, 0, true},
		{"failure: CHANGES_PERMITTED, undefined change", enumerations.CertificationPermissionChangesPermitted, 0, 0, 1, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewDocMDPCheck(testI18nProvider(t), newTestFCResult(),
				newPDFRevision(tt.permissions, tt.formFill, tt.annotation, tt.undefined),
				process.GetLevelRule(enumerations.LevelFail))
			if got := c.Process(); got != tt.want {
				t.Errorf("Process() = %v, want %v", got, tt.want)
			}
		})
	}
}
